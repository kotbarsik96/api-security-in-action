package apierrors

import "errors"

var LoginIsTaken = errors.New("Login is taken")
