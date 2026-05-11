package apierrors

import "errors"

var SpaceNotFound = errors.New("Space not found")

var InvalidIdParam = errors.New("Invalid id param: must be non-negative number")
