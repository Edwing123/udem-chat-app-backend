package models

import "github.com/Edwing123/udem-chat-app/pkg/codes"

var (
	// User errors.
	ErrUserNameExists                     = codes.NewCode("user_name_exists")
	ErrUserNameEmpty                      = codes.NewCode("user_name_empty")
	ErrUserPasswordEmpty                  = codes.NewCode("user_password_empty")
	ErrUserBirthdateEmpty                 = codes.NewCode("user_birthdate_empty")
	ErrUserBirthdateBadFormat             = codes.NewCode("user_birthdate_bad_format")
	ErrUserProfilePictureIdExists         = codes.NewCode("user_profile_picture_id_exists")
	ErrUserNameExceedsMaxLength           = codes.NewCode("user_name_exceeds_max_length")
	ErrUserPasswordNotValidLength         = codes.NewCode("user_password_not_valid_length")
	ErrUserProfilePictureIdNotValidLength = codes.NewCode("user_profile_picture_id_not_valid_length")

	// Authentication and password change errors.
	ErrPasswordMismatch = codes.NewCode("password_mismatch")
	ErrLoginFail        = codes.NewCode("login_fail")

	// Generic database errors.
	ErrNoRecords          = codes.NewCode("no_records")
	ErrDatabaseServerFail = codes.NewCode("database_server_fail")
	ErrNoUpdates          = codes.NewCode("no_updates_performed")
)
