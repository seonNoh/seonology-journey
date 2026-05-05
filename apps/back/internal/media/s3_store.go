package media

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// S3Presigner implements Presigner using real AWS S3 presigned URLs.
type S3Presigner struct {
	client *s3.PresignClient
	bucket string
	ttl    time.Duration
}

// NewS3Presigner creates a presigner backed by S3.
func NewS3Presigner(s3Client *s3.Client, bucket string, ttl time.Duration) *S3Presigner {
	return &S3Presigner{
		client: s3.NewPresignClient(s3Client),
		bucket: bucket,
		ttl:    ttl,
	}
}

func (p *S3Presigner) UploadURL(ctx context.Context, key, contentType string, size int64) (string, time.Time, error) {
	input := &s3.PutObjectInput{
		Bucket:        aws.String(p.bucket),
		Key:           aws.String(key),
		ContentType:   aws.String(contentType),
		ContentLength: aws.Int64(size),
	}
	presigned, err := p.client.PresignPutObject(ctx, input, func(opts *s3.PresignOptions) {
		opts.Expires = p.ttl
	})
	if err != nil {
		return "", time.Time{}, fmt.Errorf("s3 presign upload: %w", err)
	}
	return presigned.URL, time.Now().Add(p.ttl), nil
}

func (p *S3Presigner) DownloadURL(ctx context.Context, key string) (string, time.Time, error) {
	input := &s3.GetObjectInput{
		Bucket: aws.String(p.bucket),
		Key:    aws.String(key),
	}
	presigned, err := p.client.PresignGetObject(ctx, input, func(opts *s3.PresignOptions) {
		opts.Expires = p.ttl
	})
	if err != nil {
		return "", time.Time{}, fmt.Errorf("s3 presign download: %w", err)
	}
	return presigned.URL, time.Now().Add(p.ttl), nil
}

// S3ObjectStore implements object deletion on S3.
type S3ObjectStore struct {
	client *s3.Client
	bucket string
}

// NewS3ObjectStore creates an S3 object store.
func NewS3ObjectStore(client *s3.Client, bucket string) *S3ObjectStore {
	return &S3ObjectStore{client: client, bucket: bucket}
}

// DeleteObject removes an object from S3.
func (s *S3ObjectStore) DeleteObject(ctx context.Context, key string) error {
	_, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return fmt.Errorf("s3 delete %s: %w", key, err)
	}
	return nil
}
