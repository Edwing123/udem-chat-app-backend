package hashing

import "golang.org/x/crypto/bcrypt"

// HashPassword creates the bcrypt hash of the password at the given cost.
// On success it returns the hashed password and a nil error, and on failure,
// it returns a nil hashed password and non-nil error.
func HashPassword(password []byte, cost int) ([]byte, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword(password, cost)
	if err != nil {
		return nil, err
	}

	return hashedPassword, nil
}

// VerifyPassword compares a bcrypt hashed password with its
// possible plaintext equivalent. Returns true on success, or false on failure.
func VerifyPassword(hashedPassword, password []byte) bool {
	err := bcrypt.CompareHashAndPassword(
		hashedPassword,
		password,
	)

	return err == nil
}
