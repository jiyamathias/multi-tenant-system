package storage

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
	// ErrPasswordIncorrect when the password check failed because it is incorrect
	ErrPasswordIncorrect = errors.New("password is incorrect")
	// ErrPinIncorrect when the password check failed because it is incorrect
	ErrPinIncorrect = errors.New("pin is incorrect")
	// ErrSetPinFailed if cannot set customer pin
	ErrSetPinFailed = errors.New("unable to set user password")
	// ErrEmptyResult when result from database query is empty
	ErrEmptyResult = errors.New("the result is empty")
	// ErrDuplicateRecord when unique error occurs as a result of attempt trying to insert duplicated into the db
	ErrDuplicateRecord = errors.New("record already exist, duplicate record")
	//ErrUnauthorizedAccess if error occurred while comparing the designated role sent to the role required to perform a certain action
	ErrUnauthorizedAccess = errors.New("you have no access to perform task")
	//ErrInvalidDownloadType if an invalid download type is selected while trying to download reports or portfolios
	ErrInvalidDownloadType = errors.New("download type is invalid")
	//ErrGeneric is used to report error that we dont necessarily need to tell the user specifies.
	ErrGeneric = errors.New("request failed, try again later")
	//ErrAlreadyTrashed if a row appears to be already deleted
	ErrAlreadyTrashed = errors.New("record is already trash-ed")
	//ErrRatePlan if you want to rate more than one plan in a month
	ErrRatePlan = errors.New("you can only rate one plan in a month")
	// ErrInvalidDateDuration if the date put is invalid
	ErrInvalidDateDuration = errors.New("invalid period value")
)
