package profile

import "errors"

var (
	ErrImageTypeNotSupported = errors.New("image type not supported")
	ErrImageProccesingFail   = errors.New("image processing fail")
	ErrCannotGetImageSize    = errors.New("cannot get image size")
	ErrImageConvertionFail   = errors.New(("image convertion fail"))
)
