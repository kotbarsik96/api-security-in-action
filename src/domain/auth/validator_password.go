package auth

import (
	"fmt"
	"regexp"
)

var LoginRegex = regexp.MustCompile(`^[\p{L}\p{M}\p{N}\p{Z}\p{P}]+$`)

var LoginMaxLength = 100

func ValidateLogin(login string) error {
	l := len(login)
	if l > LoginMaxLength {
		return fmt.Errorf("Login must have less than %v symbols", LoginMaxLength)
	}

	if match := LoginRegex.MatchString(login); !match {
		return fmt.Errorf("Login has invalid characters")
	}

	return nil
}
