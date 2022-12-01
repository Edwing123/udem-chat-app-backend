package codes

import "errors"

var (
	ErrUserNameAlreadyExists     = errors.New("user_name_already_exists")
	ErrUserNameEmpty             = errors.New("user_name_empty")
	ErrUserPasswordEmpty         = errors.New("user_password_empty")
	ErrUserBirtdateEmpty         = errors.New("user_birthdate_empty")
	ErrUserProfilePictureIdEmpty = errors.New("user_profile_picture_id_empty")
	ErrPasswordsMismatch         = errors.New("current_password_new_password_mismatch")

	ErrNoRecords           = errors.New("no_records")
	ErrLoginFail           = errors.New("login_fail")
	ErrDatabaseFail        = errors.New("database_fail")
	ErrPasswordHashingFail = errors.New("password_hashing_fail")

	ErrNotLoggedIn    = errors.New("not_logged_in")
	ErrServerInternal = errors.New("server_internal")
	ErrClient         = errors.New("client_error")
)
