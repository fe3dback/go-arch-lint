package models

type CodeReference struct {
	Pointer  Reference
	LineFrom int
	LineTo   int
}

func NewCodeReference(pointer Reference, lineFrom, lineTo int) CodeReference {
	return CodeReference{
		Pointer:  pointer,
		LineFrom: lineFrom,
		LineTo:   lineTo,
	}
}

func NewCodeReferenceRelative(pointer Reference, relTop, relBottom int) CodeReference {
	return NewCodeReference(
		pointer,
		pointer.Line-relTop,
		pointer.Line+relBottom,
	)
}
