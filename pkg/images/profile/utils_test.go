package profile

import (
	"testing"

	"github.com/h2non/bimg"
)

func TestPercentageCropToPixelsCrop(t *testing.T) {
	// 1000x1000 image.
	size := bimg.ImageSize{
		Width:  1000,
		Height: 1000,
	}

	// Crop the image to a 50%x50% square in the center.
	percentageCrop := Crop{
		Width:  50,
		Height: 50,
		X:      25,
		Y:      25,
	}

	pixelsCrop := percentageCropToPixelsCrop(percentageCrop, size)

	if pixelsCrop.Width != 500 || pixelsCrop.Height != 500 {
		t.Errorf(
			"expected (width=%d, height=%d), got (width=%d, height=%d)",
			500, 500,
			pixelsCrop.Width, pixelsCrop.Height,
		)
	}

	if pixelsCrop.X != 250 || pixelsCrop.Y != 250 {
		t.Errorf(
			"expected (x=%d, y=%d), got (x=%d, y=%d)",
			500, 500,
			pixelsCrop.X, pixelsCrop.Y,
		)
	}
}
