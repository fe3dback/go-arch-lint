package models

import "fmt"

type UserSpaceError struct {
	msg string
}

func (u UserSpaceError) Error() string {
	return fmt.Sprintf("cmd: %s", u.msg)
}

func NewUserSpaceError(msg string) UserSpaceError {
	return UserSpaceError{
		msg: msg,
	}
}
