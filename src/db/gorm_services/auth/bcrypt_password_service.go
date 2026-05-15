package auth

import (
	"bytes"

	"golang.org/x/crypto/bcrypt"
)

type BcryptPasswordService struct {
}

func NewBcryptPasswordService() *BcryptPasswordService {
	return &BcryptPasswordService{}
}

func (s *BcryptPasswordService) CreateHash(password string) (string, error) {
	buf := bytes.NewBufferString(password)
	hashedPasswordBytes, err := bcrypt.GenerateFromPassword(buf.Bytes(), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedPasswordBytes), nil
}

func (s *BcryptPasswordService) CompareHashAndPassword(hashedPassword, password []byte) bool {
	err := bcrypt.CompareHashAndPassword(hashedPassword, password)
	return err == nil
}
