package auth

import "fmt"

var PasswordMinLength = 8

var PasswordMaxLength = 200

func ValidatePassword(password string) error {
	l := len(password)

	if l < PasswordMinLength {
		return fmt.Errorf("Min password length is %v", PasswordMinLength)
	}

	if l > PasswordMaxLength {
		return fmt.Errorf("Max password length is %v", PasswordMaxLength)
	}

	return nil
}
