package media

import (
	"io"
	"time"

	"github.com/rwcarlsen/goexif/exif"
)

// EXIFData holds extracted GPS and timestamp info from image EXIF.
type EXIFData struct {
	Latitude  float64
	Longitude float64
	TakenAt   time.Time
	HasGPS    bool
	HasTime   bool
}

// ExtractEXIF reads EXIF data from an image reader and extracts GPS + DateTime.
// Returns a zero EXIFData if no EXIF is found (not an error).
func ExtractEXIF(r io.Reader) (EXIFData, error) {
	x, err := exif.Decode(r)
	if err != nil {
		// No EXIF is not an error condition.
		return EXIFData{}, nil
	}

	var data EXIFData

	lat, lon, err := x.LatLong()
	if err == nil {
		data.Latitude = lat
		data.Longitude = lon
		data.HasGPS = true
	}

	t, err := x.DateTime()
	if err == nil {
		data.TakenAt = t
		data.HasTime = true
	}

	return data, nil
}
