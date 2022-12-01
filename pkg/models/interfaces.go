package models

type UserManager interface {
	New(user User) error
	Get(id int) (User, error)
	Login(user User) (int, error)
	Update(user User) error
	ChangePassword(id int, currentPass, newPass string) error
}
