package models

import "fmt"

type (
	ErrorWithNotices struct {
		OverallMessage string
		Notices        []Notice
	}

	Notice struct {
		Message   string
		Reference Reference
	}
)

func NewErrorWithNotices(overallMessage string, notices []Notice) *ErrorWithNotices {
	return &ErrorWithNotices{
		OverallMessage: overallMessage,
		Notices:        notices,
	}
}

func (en ErrorWithNotices) Error() string {
	return fmt.Sprintf("%s (has %d notices)", en.OverallMessage, len(en.Notices))
}
