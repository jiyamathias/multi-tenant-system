package controller

import (
	"errors"
)

var (
	// ErrRecordNotFound if the record is not found in the database
	ErrRecordNotFound = errors.New("record not found")
	// ErrEmailAlreadyExists when an email already exists in the DB
	ErrEmailAlreadyExists = errors.New("user with this email already exists")
	// ErrUserDoesNotExist when an email already exists in the DB
	ErrUserDoesNotExist = errors.New("user with this email does not exist")
	// ErrIncorrectLoginDetails when user auth details are incorrect
	ErrIncorrectLoginDetails = errors.New("incorrect email or password")
	// ErrWithToken issues parsing to token
	ErrWithToken = errors.New("error occurred with reset token")
	// ErrTransactionID when theres an error getting transaction with ID
	ErrTransactionID = errors.New("error getting transaction with ID")
	// MismatchedTransactionType when the transaction type from the webhook and that in the transaction model mismatch
	MismatchedTransactionType = errors.New("mismatch transaction type")
)
