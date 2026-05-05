package media

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"strings"

	journeyv1 "github.com/seonNoh/seonology-journey/proto/gen/go/journey/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// ObjectStore abstracts S3 operations needed by media processing.
type ObjectStore interface {
	// GetObject downloads an object by key.
	GetObject(ctx context.Context, key string) (io.ReadCloser, error)
	// PutObject uploads data to the given key with the specified content type.
	PutObject(ctx context.Context, key string, data []byte, contentType string) error
	// DeleteObject removes an object by key.
	DeleteObject(ctx context.Context, key string) error
}

// ProcessUpload downloads the uploaded original from S3, extracts EXIF data,
// generates a thumbnail, uploads it, and updates the media metadata.
// This should be called after ConfirmUpload.
func (s *Service) ProcessUpload(ctx context.Context, mediaID string) error {
	m, err := s.repo.Get(ctx, mediaID)
	if err != nil {
		return err
	}
	if s.store == nil {
		slog.Warn("media: object store not configured, skipping post-processing")
		return nil
	}

	// Download original
	rc, err := s.store.GetObject(ctx, m.GetS3Key())
	if err != nil {
		return fmt.Errorf("media: download original: %w", err)
	}
	defer rc.Close()

	// Read into buffer for multi-pass (EXIF + thumbnail)
	data, err := io.ReadAll(rc)
	if err != nil {
		return fmt.Errorf("media: read original: %w", err)
	}

	// Extract EXIF
	exifData, _ := ExtractEXIF(bytes.NewReader(data))
	if exifData.HasGPS {
		m.Location = &journeyv1.GeoPoint{
			Latitude:  exifData.Latitude,
			Longitude: exifData.Longitude,
		}
	}
	if exifData.HasTime {
		m.TakenAt = timestamppb.New(exifData.TakenAt)
	}

	// Generate thumbnail (only for images)
	mimeType := m.GetMimeType()
	if isImage(mimeType) {
		thumbData, err := GenerateThumbnail(bytes.NewReader(data), mimeType)
		if err != nil {
			slog.Warn("media: thumbnail generation failed", "mediaId", mediaID, "error", err)
		} else {
			thumbKey := thumbnailKey(m.GetS3Key())
			if err := s.store.PutObject(ctx, thumbKey, thumbData, "image/jpeg"); err != nil {
				slog.Warn("media: thumbnail upload failed", "mediaId", mediaID, "error", err)
			} else {
				m.ThumbnailS3Key = thumbKey
			}
		}
	}

	// Update metadata
	m.MimeType = mimeType
	m.Size = int64(len(data))
	if m.Audit != nil {
		m.Audit.UpdatedAt = timestamppb.New(s.now().UTC())
	}
	return s.repo.Save(ctx, m)
}

func thumbnailKey(originalKey string) string {
	// trips/tripId/mediaId/filename.jpg → thumb/tripId/mediaId/filename.jpg
	parts := strings.SplitN(originalKey, "/", 2)
	if len(parts) == 2 {
		return "thumb/" + parts[1]
	}
	return "thumb/" + originalKey
}

func isImage(mime string) bool {
	return strings.HasPrefix(mime, "image/")
}
