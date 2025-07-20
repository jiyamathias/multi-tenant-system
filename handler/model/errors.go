package model

import (
	"errors"
	"fmt"
)

var (
	// ErrEmptyValidIDs if array of ids has no valid id
	ErrEmptyValidIDs = errors.New("no valid id has been submitted")

	// ErrIncompleteDetails when user details is incomplete
	ErrIncompleteDetails = errors.New("incomplete detail please fill the missing information")
	// ErrIncompleteLoginDetails for sign in
	ErrIncompleteLoginDetails = errors.New("invalid sign in details please complete form to sign in")

	// ErrIncompleteDetailsTier when user details is incomplete
	ErrIncompleteDetailsTier = errors.New("incomplete user details - tier")

	// ErrInvalidToken when user details is incomplete
	ErrInvalidToken = errors.New("incorrect Details Kindly enter the 6 digit code (otp) sent to your email")
)

// ErrDynamicInvalidUUID used used to return an invalid UUID with custom error message
func ErrDynamicInvalidUUID(uuidType string) error {
	return fmt.Errorf("the %s is not a valid uuid", uuidType)
}
