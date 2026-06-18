package security

import (
	"golang.org/x/crypto/bcrypt"
)

func GeneratePassHash(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return "", err
	}

	return string(hashed), nil
}

func IsMatch(hashed, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password))
}
