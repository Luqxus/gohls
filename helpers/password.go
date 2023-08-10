package helpers

import "golang.org/x/crypto/bcrypt"

func HashPassword(givenPassword string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(givenPassword), 14)

	return string(hash), err
}

func VerifyPassword(foundPassword, givenPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(foundPassword), []byte(givenPassword))
}
