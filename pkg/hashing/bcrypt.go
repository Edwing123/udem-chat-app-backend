package hashing

import "golang.org/x/crypto/bcrypt"

// HashPassword returns the bcrypt hash of the password at the given cost.
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
