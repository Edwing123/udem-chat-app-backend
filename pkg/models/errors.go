package models

import "errors"

var (
	ErrUserNameAlreadyExists     = errors.New("user name already exists")
	ErrUserNameEmpty             = errors.New("user name is empty")
	ErrUserPasswordEmpty         = errors.New("user password is empty")
	ErrUserBirtdateEmpty         = errors.New("user birthdate is empty")
	ErrUserProfilePictureIdEmpty = errors.New("user profile picture id is empty")
	ErrPasswordsMismatch         = errors.New("current password and new password mismatch")

	ErrNoRecords           = errors.New("no records found")
	ErrLoginFail           = errors.New("login fail")
	ErrDatabaseFail        = errors.New("database fail")
	ErrPasswordHashingFail = errors.New("password hashing fail")
)
