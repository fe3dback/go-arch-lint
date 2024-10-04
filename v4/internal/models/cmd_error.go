package models

import (
	"errors"

	"github.com/fe3dback/go-arch-lint-sdk/arch"
)

type (
	CmdStdoutErrorOut struct {
		OverallNote string        `json:"OverallNote"`
		Errors      []arch.Notice `json:"Errors"`
	}
)

type (
	// UserLandError will be rendered to ascii or json,
	// so it will not be printed again as simple stdout fmt.Println
	// when linter crashes
	// but statusCode=1 will be set anyway.
	UserLandError struct {
		err error
	}
)

// ProjectNotMatchSpecificationsError this is very special error
// marker, that work absolutely like err=nil, but it will set
// exitCode to 1. This useful when checked project has some linter
// notices.
var ProjectNotMatchSpecificationsError = errors.New("ProjectNotMatchSpecificationsError")

func NewUserLandError(err error) *UserLandError {
	return &UserLandError{
		err: err,
	}
}

func (ule *UserLandError) Error() string {
	return ule.err.Error()
}
