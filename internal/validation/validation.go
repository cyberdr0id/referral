// Package validation provides functions for requests data validating.
package validation

import (
	"errors"
	"regexp"
)

var (
	// errInvalidLength presents an error when user enters data with invalid length.
	errInvalidLength = errors.New("input parameter has invalid length")

	// errParameterRequired presents an error when user didn't fill some information.
	errParameterRequired = errors.New("input parameter is required")

	// errInvalidFile presents an error when user load file with invalid name or wrong extension.
	errInvalidFile = errors.New("invalid format or name of input file")

	// errInvalidName presents an error when user send candidate with invalid name/surname.
	errInvalidName = errors.New("input name didn't match to the desired format")
)

// CheckAuthorizationRequestData validates user authorization data.
func CheckAuthorizationRequestData(name, password string) error {
	if name == "" || password == "" {
		return errParameterRequired
	}

	if len(name) < 6 || len(name) > 18 {
		return errInvalidLength
	}

	if len(password) < 6 || len(password) > 18 {
		return errInvalidLength
	}

	return nil
}

// CheckCvRequestData checks correctness of user input data after CV sending.
func CheckCvRequestData(fileName, candidateName, candidateSurname string) error {
	if len(candidateName) == 0 || len(candidateSurname) == 0 || len(fileName) == 0 {
		return errParameterRequired
	}

	isRightFile, _ := regexp.MatchString("([a-zA-Z0-9\\s_\\.\\-\\(\\):])+(.PDF|.pdf)$", fileName)
	if !isRightFile {
		return errInvalidFile
	}

	isValidName, _ := regexp.MatchString("(^[A-Za-zА-Яа-я]{2,16})?([ ]{0,1})([A-Za-zА-Яа-я]{2,16})?", candidateName)
	if !isValidName {
		return errInvalidName
	}

	isValidSurname, _ := regexp.MatchString("(^[A-Za-zА-Яа-я]{2,16})?([ ]{0,1})([A-Za-zА-Яа-я]{2,16})?", candidateSurname)
	if !isValidSurname {
		return errInvalidName
	}

	return nil
}
