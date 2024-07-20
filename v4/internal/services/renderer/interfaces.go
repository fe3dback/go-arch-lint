package renderer

//go:generate ../../../bin/mockgen -source=$GOFILE -destination=mocks/mocks.go -package=mocks

type typeRenderer interface {
	Render(model any) (string, error)
}
