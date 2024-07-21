package models

type ReferencedError struct {
	err error
	ref Reference
}

func NewReferencedError(err error, ref Reference) ReferencedError {
	return ReferencedError{
		err: err,
		ref: ref,
	}
}

func (re ReferencedError) Reference() Reference {
	return re.ref
}

func (re ReferencedError) Error() string {
	return re.err.Error()
}
