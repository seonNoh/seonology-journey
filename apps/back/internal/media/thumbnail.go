package media

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"

	"golang.org/x/image/draw"
)

const thumbnailSize = 300

// GenerateThumbnail reads an image from r, produces a 300x300 (max) JPEG thumbnail.
// The aspect ratio is preserved; the image is scaled to fit within 300x300.
func GenerateThumbnail(r io.Reader, mimeType string) ([]byte, error) {
	img, err := decodeImage(r, mimeType)
	if err != nil {
		return nil, fmt.Errorf("thumbnail decode: %w", err)
	}

	thumb := resizeFit(img, thumbnailSize, thumbnailSize)

	var buf bytes.Buffer
	if err := jpeg.Encode(&buf, thumb, &jpeg.Options{Quality: 80}); err != nil {
		return nil, fmt.Errorf("thumbnail encode: %w", err)
	}
	return buf.Bytes(), nil
}

func decodeImage(r io.Reader, mimeType string) (image.Image, error) {
	switch mimeType {
	case "image/jpeg", "image/jpg":
		return jpeg.Decode(r)
	case "image/png":
		return png.Decode(r)
	default:
		// fallback: try generic decode
		img, _, err := image.Decode(r)
		return img, err
	}
}

// resizeFit scales img to fit within maxW x maxH preserving aspect ratio.
func resizeFit(img image.Image, maxW, maxH int) image.Image {
	bounds := img.Bounds()
	srcW := bounds.Dx()
	srcH := bounds.Dy()

	if srcW <= maxW && srcH <= maxH {
		return img
	}

	ratio := float64(srcW) / float64(srcH)
	newW, newH := maxW, maxH
	if ratio > 1 {
		newH = int(float64(maxW) / ratio)
	} else {
		newW = int(float64(maxH) * ratio)
	}

	dst := image.NewRGBA(image.Rect(0, 0, newW, newH))
	draw.CatmullRom.Scale(dst, dst.Bounds(), img, bounds, draw.Over, nil)
	return dst
}
