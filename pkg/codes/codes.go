package codes

func NewCode(description string) Code {
	return Code(description)
}

// Represents an error code.
type Code string

// It implements the error interface.
func (c Code) Error() string {
	return string(c)
}

var (
	ErrUserNameAlreadyExists     = NewCode("user_name_already_exists")
	ErrUserNameEmpty             = NewCode("user_name_empty")
	ErrUserPasswordEmpty         = NewCode("user_password_empty")
	ErrUserBirtdateEmpty         = NewCode("user_birthdate_empty")
	ErrUserProfilePictureIdEmpty = NewCode("user_profile_picture_id_empty")
	ErrPasswordsMismatch         = NewCode("current_password_new_password_mismatch")

	ErrNoRecords           = NewCode("no_records")
	ErrLoginFail           = NewCode("login_fail")
	ErrDatabaseFail        = NewCode("database_fail")
	ErrPasswordHashingFail = NewCode("password_hashing_fail")
	ErrNoUpdatesToPerform  = NewCode("no_updates_to_perform")

	ErrNotLoggedIn           = NewCode("not_logged_in")
	ErrServerInternal        = NewCode("server_internal")
	ErrClient                = NewCode("client_error")
	ErrPasswordNotValid      = NewCode("password_not_valid")
	ErrAuthRequired          = NewCode("auth_required")
	ErrProfileImageTooBig    = NewCode("profile_image_too_big")
	ErrImageTypeNotSupported = NewCode("image_type_not_supported")
)
