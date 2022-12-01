package profile

import (
	"github.com/h2non/bimg"
)

// Returns true if imageType is a supported,
// otherwise it returns false.
//
// The supported image types are: jpg, webp and png.
func IsImageTypeSupported(imageType bimg.ImageType) bool {
	return imageType == bimg.JPEG || imageType == bimg.WEBP || imageType == bimg.PNG
}

// Returns true if imageType is a supported,
// otherwise it returns false.
//
// The supported image types are: jpg, webp and png.
func IsImageTypeNameSupported(imageType string) bool {
	return imageType == "jpeg" || imageType == "webp" || imageType == "png"
}

// Creates a new profile manager which will save
// profile images on the file system under the provided dir.
func New(dir string) Manager {
	return Manager{
		rootDir: dir,
	}
}

func percentageCropToPixelsCrop(percentageCrop Crop, size bimg.ImageSize) Crop {
	croppedWidth := size.Width * percentageCrop.Width / 100
	croppedHeight := size.Height * percentageCrop.Height / 100
	xInPixels := size.Width * percentageCrop.X / 100
	yInPixels := size.Height * percentageCrop.Y / 100

	return Crop{
		Width:  croppedWidth,
		Height: croppedHeight,
		X:      xInPixels,
		Y:      yInPixels,
	}
}
