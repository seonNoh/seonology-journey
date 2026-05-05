package media

import (
	"bytes"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"testing"
)

func createTestJPEG(w, h int) []byte {
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	for y := range h {
		for x := range w {
			img.Set(x, y, color.RGBA{R: uint8(x % 256), G: uint8(y % 256), B: 128, A: 255})
		}
	}
	var buf bytes.Buffer
	jpeg.Encode(&buf, img, nil)
	return buf.Bytes()
}

func createTestPNG(w, h int) []byte {
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	for y := range h {
		for x := range w {
			img.Set(x, y, color.RGBA{R: 200, G: 100, B: 50, A: 255})
		}
	}
	var buf bytes.Buffer
	png.Encode(&buf, img)
	return buf.Bytes()
}

func TestGenerateThumbnail_JPEG(t *testing.T) {
	t.Parallel()
	data := createTestJPEG(800, 600)
	thumb, err := GenerateThumbnail(bytes.NewReader(data), "image/jpeg")
	if err != nil {
		t.Fatalf("GenerateThumbnail: %v", err)
	}
	if len(thumb) == 0 {
		t.Fatal("thumbnail empty")
	}

	// Decode and check dimensions
	img, err := jpeg.Decode(bytes.NewReader(thumb))
	if err != nil {
		t.Fatalf("decode thumb: %v", err)
	}
	bounds := img.Bounds()
	if bounds.Dx() > thumbnailSize || bounds.Dy() > thumbnailSize {
		t.Errorf("thumbnail too large: %dx%d", bounds.Dx(), bounds.Dy())
	}
	// 800x600 → 300x225 (aspect preserved)
	if bounds.Dx() != 300 {
		t.Errorf("expected width 300, got %d", bounds.Dx())
	}
}

func TestGenerateThumbnail_PNG(t *testing.T) {
	t.Parallel()
	data := createTestPNG(400, 800)
	thumb, err := GenerateThumbnail(bytes.NewReader(data), "image/png")
	if err != nil {
		t.Fatalf("GenerateThumbnail: %v", err)
	}
	img, err := jpeg.Decode(bytes.NewReader(thumb))
	if err != nil {
		t.Fatalf("decode: %v", err)
	}
	bounds := img.Bounds()
	// 400x800 → 150x300 (aspect preserved)
	if bounds.Dy() != 300 {
		t.Errorf("expected height 300, got %d", bounds.Dy())
	}
}

func TestGenerateThumbnail_SmallImage(t *testing.T) {
	t.Parallel()
	data := createTestJPEG(100, 100)
	thumb, err := GenerateThumbnail(bytes.NewReader(data), "image/jpeg")
	if err != nil {
		t.Fatalf("GenerateThumbnail: %v", err)
	}
	img, err := jpeg.Decode(bytes.NewReader(thumb))
	if err != nil {
		t.Fatalf("decode: %v", err)
	}
	bounds := img.Bounds()
	// Should not upscale
	if bounds.Dx() != 100 || bounds.Dy() != 100 {
		t.Errorf("small image should not be resized: %dx%d", bounds.Dx(), bounds.Dy())
	}
}

func TestExtractEXIF_NoEXIF(t *testing.T) {
	t.Parallel()
	data := createTestJPEG(100, 100)
	result, err := ExtractEXIF(bytes.NewReader(data))
	if err != nil {
		t.Fatalf("ExtractEXIF: %v", err)
	}
	if result.HasGPS {
		t.Error("expected no GPS")
	}
	if result.HasTime {
		t.Error("expected no time")
	}
}

func TestExtractEXIF_InvalidData(t *testing.T) {
	t.Parallel()
	result, err := ExtractEXIF(bytes.NewReader([]byte("not an image")))
	if err != nil {
		t.Fatalf("ExtractEXIF should not error on invalid data: %v", err)
	}
	if result.HasGPS || result.HasTime {
		t.Error("expected empty result for invalid data")
	}
}

func TestThumbnailKey(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"trips/trip-1/media-1/photo.jpg", "thumb/trip-1/media-1/photo.jpg"},
		{"trips/abc/def/image.png", "thumb/abc/def/image.png"},
		{"singlepart", "thumb/singlepart"},
	}
	for _, tt := range tests {
		got := thumbnailKey(tt.input)
		if got != tt.want {
			t.Errorf("thumbnailKey(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}
