package sqlserver

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/Edwing123/udem-chat-app/pkg/codes"
	"github.com/Edwing123/udem-chat-app/pkg/models"
	"github.com/Edwing123/udem-chat-app/pkg/validations/hashing"
)

type UserManager struct {
	db *sql.DB
}

func (um *UserManager) UserNameAlreadyExists(name string) error {
	row := um.db.QueryRowContext(
		rootCtx,
		getUserByName,
		sql.Named(userName, name),
	)

	var existingName string
	err := row.Scan(&existingName)

	// If the err is not nil, that means there's alredy a user
	// with the same name.
	if err == nil {
		return codes.ErrUserNameAlreadyExists
	}

	return nil
}
func (um *UserManager) New(user models.User) error {
	if user.Name == "" {
		return codes.ErrUserNameEmpty
	}

	if user.Password == "" {
		return codes.ErrUserPasswordEmpty
	}

	if user.Birthdate == "" {
		return codes.ErrUserBirtdateEmpty
	}

	// Let's check if there is already a user with the same name.
	err := um.UserNameAlreadyExists(user.Name)
	if err != nil {
		return err
	}

	// Hash the password.
	hashedPassword, err := hashing.HashPassword([]byte(user.Password))
	if err != nil {
		return codes.ErrPasswordHashingFail
	}

	_, err = um.db.ExecContext(
		rootCtx,
		insertUser,
		sql.Named(userName, user.Name),
		sql.Named(userPassword, string(hashedPassword)),
		sql.Named(userBirthdate, user.Birthdate),
	)
	if err != nil {
		return codes.ErrDatabaseFail
	}

	return nil
}

func (um *UserManager) Get(id int) (models.User, error) {
	row := um.db.QueryRowContext(
		rootCtx,
		getUserById,
		sql.Named(userId, id),
	)

	var user models.User

	err := row.Scan(
		&user.Id,
		&user.Name,
		&user.Birthdate,
		&user.ProfilePictureId,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return user, codes.ErrNoRecords
		}

		return user, codes.ErrDatabaseFail
	}

	return user, nil
}

func (um *UserManager) Login(user models.User) (int, error) {
	if user.Name == "" || user.Password == "" {
		return 0, codes.ErrLoginFail
	}

	var id int
	var hashedPassword string

	row := um.db.QueryRowContext(
		rootCtx,
		getUserIdAndPasswordByName,
		sql.Named(userName, user.Name),
	)

	err := row.Scan(
		&id,
		&hashedPassword,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, codes.ErrLoginFail
		}

		return 0, codes.ErrDatabaseFail
	}

	isPasswordValid := hashing.VerifyPassword([]byte(hashedPassword), []byte(user.Password))
	if !isPasswordValid {
		return 0, codes.ErrLoginFail
	}

	return id, nil
}

func (um *UserManager) Update(id int, user models.User) error {
	log.Printf("%+v\n", user)

	// Only update non-empty fields.
	fieldsToUpdate := []string{}
	values := []any{}

	if user.Name != "" {
		fieldsToUpdate = append(fieldsToUpdate, fmt.Sprintf("%s = @%s", userName, userName))
		values = append(values, sql.Named(userName, user.Name))

		// Let's check if there is already a user with the same name.
		err := um.UserNameAlreadyExists(user.Name)
		if err != nil {
			return err
		}
	}

	if user.Birthdate != "" {
		fieldsToUpdate = append(fieldsToUpdate, fmt.Sprintf("%s = @%s", userBirthdate, userBirthdate))
		values = append(values, sql.Named(userBirthdate, user.Birthdate))
	}

	if user.ProfilePictureId != "" {
		fieldsToUpdate = append(fieldsToUpdate, fmt.Sprintf("%s = @%s", userProfilePictureId, userProfilePictureId))
		values = append(values, sql.Named(userProfilePictureId, user.ProfilePictureId))
	}

	// If all update-able fields are empty.
	if len(fieldsToUpdate) < 1 {
		return codes.ErrNoUpdatesToPerform
	}

	query := fmt.Sprintf(
		`UPDATE [User] SET %s WHERE [Id] = @Id;`,
		strings.Join(fieldsToUpdate, ","),
	)

	values = append(values, sql.Named(userId, id))

	_, err := um.db.ExecContext(
		rootCtx,
		query,
		values...,
	)
	if err != nil {
		return codes.ErrDatabaseFail
	}

	return nil
}

func (um *UserManager) ChangePassword(id int, currentPass, newPass string) error {
	var hashedPassword string

	row := um.db.QueryRowContext(
		rootCtx,
		getUserPasswordById,
		sql.Named(userId, id),
	)

	err := row.Scan(&hashedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return codes.ErrNoRecords
		}

		return codes.ErrDatabaseFail
	}

	isValidPassword := hashing.VerifyPassword([]byte(hashedPassword), []byte(currentPass))
	if !isValidPassword {
		return codes.ErrPasswordsMismatch
	}

	newPassHashed, err := hashing.HashPassword([]byte(newPass))
	if err != nil {
		return codes.ErrPasswordHashingFail
	}

	_, err = um.db.ExecContext(
		rootCtx,
		updateUserPassword,
		sql.Named(userPassword, string(newPassHashed)),
	)
	if err != nil {
		return codes.ErrDatabaseFail
	}

	return nil
}
