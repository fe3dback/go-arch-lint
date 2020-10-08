package models

type UserSpaceEmptyError struct{}

func (u UserSpaceEmptyError) Error() string {
	return ""
}

func NewUserSpaceEmptyError() UserSpaceEmptyError {
	return UserSpaceEmptyError{}
}
