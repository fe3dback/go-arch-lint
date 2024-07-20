package renderer

//go:generate ../../../bin/mockgen -source=$GOFILE -destination=mocks/mocks.go -package=mocks

type asciiRenderer interface {
	Render(model any) (string, error)
}

type jsonRenderer interface {
	Render(model any, format bool) (string, error)
}
