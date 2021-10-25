// Package validation provides functions for requests data validating.
package validation

import (
	"errors"
	"regexp"
)

var (
	// errInvalidLength presents an error when user input data with invalid length.
	errInvalidLength = errors.New("input parameter has invalid length")

	// errParameterRequired presents an error when user didn't fill some information.
	errParameterRequired = errors.New("input parameter is required")

	// errInvalidFile presents an error when user load file with invalid name or wrong extension.
	errInvalidFile = errors.New("invalid format or name of input file")
)

// CheckAuthorizationRequestData validates user authorization data.
func CheckAuthorizationRequestData(name, password string) error {
	if len(name) < 6 || len(name) > 18 {
		return errInvalidLength
	}

	if len(password) < 6 || len(password) > 18 {
		return errInvalidLength
	}

	if name == "" || password == "" {
		return errParameterRequired
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

	return nil
}
