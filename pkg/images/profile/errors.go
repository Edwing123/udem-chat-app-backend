package profile

import (
	"github.com/Edwing123/udem-chat-app/pkg/codes"
)

var (
	ErrImageTypeNotSupported = codes.NewCode("image_type_not_supported")
	ErrImageProcessFail      = codes.NewCode("image_process_fail")
	ErrCannotGetImageSize    = codes.NewCode("cannot_get_image_size")
	ErrImageConvertionFail   = codes.NewCode(("image_convertion_fail"))
	ErrImageWriteFail        = codes.NewCode(("image_write_fail"))
	ErrImageArchiveFail      = codes.NewCode(("image_archive_fail"))
)
