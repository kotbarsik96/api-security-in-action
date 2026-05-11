package auth

import (
	"bytes"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	buf := bytes.NewBufferString(password)
	hashedPasswordBytes, err := bcrypt.GenerateFromPassword(buf.Bytes(), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedPasswordBytes), nil
}
