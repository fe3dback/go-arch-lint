package models

type (
	UserSpaceError struct {
		msg string
	}

	ReferableErr struct {
		original  error
		reference Reference
	}
)

func (u UserSpaceError) Error() string {
	return u.msg
}

func (r ReferableErr) Error() string {
	return r.original.Error()
}

func (r ReferableErr) Reference() Reference {
	return r.reference
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

func (r ReferableErr) Is(err error) bool {
	if err == nil {
		return false
	}

	if _, ok := err.(ReferableErr); ok {
		return true
	}

	return false
}

func NewUserSpaceError(msg string) UserSpaceError {
	return UserSpaceError{
		msg: msg,
	}
}

func NewReferableErr(err error, ref Reference) ReferableErr {
	return ReferableErr{
		original:  err,
		reference: ref,
	}
}
