package helper

import "errors"

var (
	// ErrRecordNotFound if the record is not found in the database
	ErrRecordNotFound = errors.New("record not found")
	// ErrRecordCreatingFailed if error occurred while trying to insert record into the database
	ErrRecordCreatingFailed = errors.New("record failed to insert")
	// ErrRecordUpdateFailed if error occurred in attempt to update the row
	ErrRecordUpdateFailed = errors.New("record update failed")
	// ErrDeleteFailed if error occurred in an attempt to delete a record from the database
	ErrDeleteFailed = errors.New("failed to delete record")
	// ErrInvalidResponse when the response cannot be interpreted
	ErrInvalidResponse = errors.New("invalid response")
	// ErrServiceUnsupported when service is currently unsupported by the provider
	ErrServiceUnsupported = errors.New("service currently unsupported")
	// ErrSomeFieldsMissing some  required fields are missinf
	ErrSomeFieldsMissing = errors.New("some fields are missing")
	// ErrDefault is a default error
	ErrDefault = errors.New("an error occured")
	// ErrUserIDParamsMissing for user id params missing
	ErrUserIDParamsMissing = errors.New("user id params missing")
	// ErrIDMissing for id missing
	ErrIDMissing = errors.New("id is missing")
	// ErrCreatingAcctNumber for generating random account number
	ErrCreatingAcctNumber = errors.New("error genarating account number")
)
