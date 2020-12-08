package models

type UserSpaceError struct {
	msg string
}

func (u UserSpaceError) Error() string {
	return u.msg
}

func (u UserSpaceError) Is(err error) bool {
	if err == nil {
		return false
	}

	if _, ok := err.(UserSpaceError); ok {
		return true
	}

	return false
}

func NewUserSpaceError(msg string) UserSpaceError {
	return UserSpaceError{
		msg: msg,
	}
}
