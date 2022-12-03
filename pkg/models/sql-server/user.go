package sqlserver

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Edwing123/udem-chat-app/pkg/models"
	"github.com/Edwing123/udem-chat-app/pkg/validations/hashing"
	mssql "github.com/microsoft/go-mssqldb"
	"golang.org/x/exp/slog"
)

type UserManager struct {
	db     *sql.DB
	logger *slog.Logger
}

func isUserNameExistsError(err error) bool {
	var sqlErr mssql.Error
	_ = errors.As(err, &sqlErr)
	return strings.Contains(sqlErr.Message, "Unique_User_Name")
}

func isValidBirthdateFormat(birthdate string) bool {
	_, err := time.Parse(models.UserBirthdateFormat, birthdate)
	return err == nil
}

func (um *UserManager) New(user models.User) error {
	var err error

	// Validate user input.
	switch true {
	case user.Name == "":
		err = models.ErrUserNameEmpty

	case len(user.Name) > models.UserNameMaxLength:
		err = models.ErrUserNameExceedsMaxLength

	case user.Password == "":
		err = models.ErrUserPasswordEmpty

	case user.Birthdate == "":
		err = models.ErrUserBirthdateEmpty

	case !isValidBirthdateFormat(user.Birthdate):
		err = models.ErrUserBirthdateBadFormat
	}

	if err != nil {
		return err
	}

	// Hash the password.
	hashedPassword, err := hashing.HashPassword([]byte(user.Password))
	if err != nil {
		um.logger.Error("Hash password", err)
		return hashing.ErrPasswordHashingFail
	}

	_, err = um.db.ExecContext(
		rootCtx,
		insertUser,
		sql.Named(userName, user.Name),
		sql.Named(userPassword, string(hashedPassword)),
		sql.Named(userBirthdate, user.Birthdate),
	)
	if err != nil {
		if isUserNameExistsError(err) {
			err = models.ErrUserNameExists
		} else {
			um.logger.Error("New user", err, "name", user.Name, "birthdate", user.Birthdate)
			err = models.ErrDatabaseServerFail
		}

		return err
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
	var nullableImageId sql.NullString

	err := row.Scan(
		&user.Id,
		&user.Name,
		&user.Birthdate,
		&nullableImageId,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return user, models.ErrNoRecords
		}

		um.logger.Error("Get user", err, "userId", id)
		return user, models.ErrDatabaseServerFail
	}

	if nullableImageId.Valid {
		user.ProfilePictureId = nullableImageId.String
	}

	return user, nil
}

func (um *UserManager) Login(user models.User) (int, error) {
	var err error

	// Validate user input.
	switch true {
	case user.Name == "":
		err = models.ErrUserNameEmpty

	case len(user.Name) > models.UserNameMaxLength:
		err = models.ErrUserNameExceedsMaxLength

	case user.Password == "":
		err = models.ErrUserPasswordEmpty
	}

	var userId int
	var hashedPassword string

	row := um.db.QueryRowContext(
		rootCtx,
		getUserIdAndPasswordByName,
		sql.Named(userName, user.Name),
	)
	err = row.Scan(
		&userId,
		&hashedPassword,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, models.ErrLoginFail
		}

		um.logger.Error("Login user", err, "name", user.Name)
		return 0, models.ErrDatabaseServerFail
	}

	isPasswordValid := hashing.VerifyPassword([]byte(hashedPassword), []byte(user.Password))
	if !isPasswordValid {
		return 0, models.ErrLoginFail
	}

	return userId, nil
}

func (um *UserManager) Update(id int, user models.User) (models.User, string, error) {
	// Only update non-empty fields.
	fieldsToUpdate := []string{}
	values := []any{}

	if user.Name != "" {
		if len(user.Name) > models.UserNameMaxLength {
			return models.User{}, "", models.ErrUserNameExceedsMaxLength
		}

		fieldsToUpdate = append(fieldsToUpdate, fmt.Sprintf("%s = @%s", userName, userName))
		values = append(values, sql.Named(userName, user.Name))
	}

	if user.Birthdate != "" {
		if !isValidBirthdateFormat(user.Birthdate) {
			return models.User{}, "", models.ErrUserBirthdateBadFormat
		}

		fieldsToUpdate = append(fieldsToUpdate, fmt.Sprintf("%s = @%s", userBirthdate, userBirthdate))
		values = append(values, sql.Named(userBirthdate, user.Birthdate))
	}

	if user.ProfilePictureId != "" {
		if len(user.ProfilePictureId) > models.UserProfilePictureIdLength {
			return models.User{}, "", models.ErrUserProfilePictureIdNotValidLength
		}

		fieldsToUpdate = append(fieldsToUpdate, fmt.Sprintf("%s = @%s", userProfilePictureId, userProfilePictureId))
		values = append(values, sql.Named(userProfilePictureId, user.ProfilePictureId))
	}

	// If all update-able fields are empty, then return an error to notify
	// that no updates were performed.
	if len(fieldsToUpdate) == 0 {
		return models.User{}, "", models.ErrNoUpdates
	}

	// Otherwise, build the query with the columns
	// that will be updated.
	query := fmt.Sprintf(
		`UPDATE [User] SET %s WHERE [Id] = @Id;`,
		strings.Join(fieldsToUpdate, ","),
	)

	values = append(values, sql.Named(userId, id))

	tx, err := um.db.BeginTx(rootCtx, &sql.TxOptions{})
	if err != nil {
		um.logger.Error("Update user - begin transaction", err, "details")
		return models.User{}, "", models.ErrDatabaseServerFail
	}
	defer tx.Rollback()

	// Get current profile picture id only
	// if it's meant to be updated, that is,
	// the value for user.ProfilePictureId is
	// not empty.
	var oldImageId string

	if user.ProfilePictureId != "" {
		row := tx.QueryRowContext(rootCtx, getUserProfilePictureIdById, sql.Named(userId, id))

		var nullableImageId sql.NullString

		err = row.Scan(&nullableImageId)
		if err != nil {
			um.logger.Error("Update user - select current picture id", err)
			return models.User{}, "", models.ErrDatabaseServerFail
		}

		if nullableImageId.Valid {
			oldImageId = nullableImageId.String
		}
	}

	_, err = tx.ExecContext(
		rootCtx,
		query,
		values...,
	)
	if err != nil {
		if isUserNameExistsError(err) {
			err = models.ErrUserNameExists
		} else {
			err = models.ErrDatabaseServerFail
			um.logger.Error(
				"Update user", err,
				"userId", id,
				"name", user.Name,
				"birthdate", user.Birthdate,
				"profilePictureId", user.ProfilePictureId,
			)
		}

		return models.User{}, "", err
	}

	err = tx.Commit()
	if err != nil {
		um.logger.Error("Update user - close transaction", err)
		return models.User{}, "", models.ErrDatabaseServerFail
	}

	return user, oldImageId, nil
}

func (um *UserManager) ChangePassword(id int, currentPass, newPass string) error {
	var hashedPassword string

	if currentPass == "" || newPass == "" {
		return models.ErrUserPasswordEmpty
	}

	row := um.db.QueryRowContext(
		rootCtx,
		getUserPasswordById,
		sql.Named(userId, id),
	)

	err := row.Scan(&hashedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.ErrNoRecords
		}

		um.logger.Error("Change user password", err, "userId", id)
		return models.ErrDatabaseServerFail
	}

	isValidPassword := hashing.VerifyPassword([]byte(hashedPassword), []byte(currentPass))
	if !isValidPassword {
		return models.ErrPasswordMismatch
	}

	newHashedPass, err := hashing.HashPassword([]byte(newPass))
	if err != nil {
		return hashing.ErrPasswordHashingFail
	}

	_, err = um.db.ExecContext(
		rootCtx,
		updateUserPassword,
		sql.Named(userPassword, string(newHashedPass)),
		sql.Named(userId, id),
	)
	if err != nil {
		return models.ErrDatabaseServerFail
	}

	return nil
}
