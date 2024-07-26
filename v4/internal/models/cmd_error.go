package models

type (
	CmdStdoutErrorOut struct {
		OverallNote string         `json:"OverallNote"`
		Errors      []StdoutNotice `json:"Errors"`
	}

	StdoutNotice struct {
		Notice
		Preview string `json:"-"`
	}
)

// UserLandError will be rendered to ascii or json,
// so it will not be printed again as simple stdout fmt.Println
// when linter crashes
// but statusCode=1 will be set anyway.
type UserLandError struct {
	err error
}

func NewUserLandError(err error) *UserLandError {
	return &UserLandError{
		err: err,
	}
}

func (ule *UserLandError) Error() string {
	return ule.err.Error()
}
