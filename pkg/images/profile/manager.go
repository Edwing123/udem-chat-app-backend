package profile

import (
	"os"
	"path"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/h2non/bimg"
	"golang.org/x/exp/slog"
)

// Manager.
// TODO: add more documentation about its purpose.
type Manager struct {
	// The file system directory where
	// the profile manager will store
	// the images.
	rootDir string

	logger *slog.Logger
}

// InitDirs creates the directories where
// the images will be stored.
//
// The directories hierarchy will be:
//   - rootDir
//   - rootDir/active (images that can be served to requests)
//   - rootDir/archive (images are that archived )
//   - rootDir/original (images before cropping, resizing)
func (pm *Manager) InitDirs() error {
	rootDir := pm.rootDir

	for _, format := range []string{"jpeg", "png", "webp"} {
		err := os.MkdirAll(
			path.Join(
				rootDir,
				activeDir,
				format,
			),
			0o755,
		)
		if err != nil {
			return err
		}
	}

	for _, format := range []string{"jpeg", "png", "webp"} {
		err := os.MkdirAll(
			path.Join(
				rootDir,
				archiveDir,
				format,
			),
			0o755,
		)
		if err != nil {
			return err
		}
	}

	err := os.MkdirAll(
		path.Join(
			rootDir,
			originalDir,
		),
		0o755,
	)
	if err != nil {
		return err
	}

	return nil
}

const (
	activeDir   = "active"
	archiveDir  = "archive"
	originalDir = "original"
)

var (
	defaultProcessOptions = bimg.Options{
		StripMetadata: true,
		Quality:       5,
		Lossless:      true,
	}
)

// Creates a new profile image cropped by the provided crop.
// it returns the UUID of the image and a nil error if everything
// goes sucessful, otherwise an empty string and a non-nil error is returned.
func (pm *Manager) New(image Image, crop Crop) (string, error) {
	if !IsImageTypeSupported(image.Type) {
		return "", ErrImageTypeNotSupported
	}

	originalImageBuffer, err := pm.process(image.Buffer, defaultProcessOptions)
	if err != nil {
		pm.logger.Error("Process original image", err, "imageType", image.Type)
		return "", ErrImageProcessFail
	}

	// Create a unique id for the image.
	imageId := uuid.New().String()

	err = pm.save(imageId, "original", originalImageBuffer)
	if err != nil {
		pm.logger.Error("Save original image", err, "imageType", image.Type, "imageId", imageId)
		return "", ErrImageWriteFail
	}

	// Crop the image.
	bImage := bimg.NewImage(originalImageBuffer)

	imageSize, err := bImage.Size()
	if err != nil {
		pm.logger.Error("Get image size", err, "imageType", image.Type, "imageId", imageId)
		return "", ErrCannotGetImageSize
	}

	// The provided crop is assumed to be in percentages,
	// so a convertion is performed to get the crop
	// in pixels based on the image size.
	cropInPixels := percentageCropToPixelsCrop(crop, imageSize)

	// The crop has a 1:1 aspect ratio, that's why
	// the width is set as both the width and height.
	croppedImageBuffer, err := bImage.Process(bimg.Options{
		Width:  cropInPixels.Width,
		Height: cropInPixels.Width,
		Left:   crop.X,
		Top:    cropInPixels.Y,
		Crop:   true,
	})
	if err != nil {
		pm.logger.Error("Process image crop", err, "imageType", image.Type, "imageId", imageId)
		return "", ErrImageProcessFail
	}

	// Resize the image.
	croppedAndResizedImageBuffer, err := bimg.NewImage(croppedImageBuffer).Resize(400, 400)
	if err != nil {
		pm.logger.Error("Process image resize", err, "imageType", image.Type, "imageId", imageId)
		return "", ErrImageProcessFail
	}

	// Convert the cropped and resized image to the remaining formats.
	images := []Image{
		{
			Type:   image.Type,
			Buffer: croppedAndResizedImageBuffer,
		},
	}

	convertedImages, err := pm.convert(images[0])
	if err != nil {
		return "", err
	}

	images = append(images, convertedImages...)

	// Save images to activeDir.
	for _, image := range images {
		typeName := bimg.ImageTypeName(image.Type)

		err := pm.save(imageId, path.Join(activeDir, typeName), image.Buffer)
		if err != nil {
			pm.logger.Error("Save image", err, "imageType", image.Type, "imageId", imageId)
			return "", err
		}
	}

	return imageId, nil
}

// Converts the provided image to the remaining formats.
func (pm *Manager) convert(image Image) ([]Image, error) {
	var imagesType []bimg.ImageType

	switch image.Type {
	// In case it's already in jpg, convert to webp and png.
	case bimg.JPEG:
		imagesType = append(imagesType, bimg.WEBP, bimg.PNG)

	// In case it's already in webp, convert to jpg and png.
	case bimg.WEBP:
		imagesType = append(imagesType, bimg.JPEG, bimg.PNG)

	// In case it's already in png, convert to jpg and webp.
	case bimg.PNG:
		imagesType = append(imagesType, bimg.JPEG, bimg.WEBP)
	}

	images := []Image{}

	for _, format := range imagesType {
		bimg := bimg.NewImage(image.Buffer)

		convertedImageBuffer, err := bimg.Convert(format)
		if err != nil {
			pm.logger.Error("Convert image", err, "imageType", format)
			return nil, ErrImageConvertionFail
		}

		images = append(images, Image{
			Type:   format,
			Buffer: convertedImageBuffer,
		})
	}

	return images, nil
}

// Process the image buffer with the provided options.
func (pm *Manager) process(buffer []byte, options bimg.Options) ([]byte, error) {
	image := bimg.NewImage(buffer)

	imageBuffer, err := image.Process(options)
	if err != nil {
		return nil, err
	}

	return imageBuffer, nil
}

// Saves the image buffer with name id to directory dir.
func (pm *Manager) save(id string, dir string, buffer []byte) error {
	fileName := path.Join(pm.rootDir, dir, id)

	err := bimg.Write(fileName, buffer)
	if err != nil {
		return err
	}

	return nil
}

// Move images identified by id from
// the active dir to the archive dir.
func (pm *Manager) Archive(id string) error {
	for _, format := range []string{"jpeg", "webp", "png"} {
		err := os.Rename(
			path.Join(
				pm.rootDir,
				activeDir,
				format,
				id,
			),
			path.Join(
				pm.rootDir,
				archiveDir,
				format,
				id,
			),
		)

		if err != nil {
			pm.logger.Error("Archive image", err, "imageId", id)
			return err
		}
	}

	return nil
}

// Handler for serving images to requests.
func (pm *Manager) ServeImage(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return fiber.ErrNotFound
	}

	// If it's not present, default to jpeg.
	imageType := c.Query("type", "jpeg")

	// If it's represent, check if the type is supported.
	if !IsImageTypeNameSupported(imageType) {
		return fiber.ErrBadRequest
	}

	// Get image.
	return c.SendFile(
		path.Join(
			pm.rootDir,
			activeDir,
			imageType,
			id,
		),
	)
}

// Image represents the buffer and type
// of an image to be proccesed.
type Image struct {
	Type   bimg.ImageType
	Buffer []byte
}

// Crop represent the values used
// to crop the profile image.
type Crop struct {
	Width  int
	Height int
	X      int
	Y      int
}
