package utils

import "golang.org/x/crypto/bcrypt"

// HashPassword hashes a plain text password using bcrypt.
// It returns the hashed password string or an error if hashing fails.
func HashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return "", err
	}

	return string(hashedBytes), nil
}

// short hand property, string on right side is applied to the hash type too, hence its a string
func VerifyPassword(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
