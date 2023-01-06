package main

import (
	"errors"
)

var (
	ErrNoFiles          = errors.New("No such file present")
	ErrInvalidColumn    = errors.New("No such column present")
	ErrNotNumber        = errors.New("Data is not numeric")
	ErrInvalidOperation = errors.New("Invalid operation requested")
)
