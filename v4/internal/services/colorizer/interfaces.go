package colorizer

//go:generate ../../../bin/mockgen -source=$GOFILE -destination=mocks/mocks.go -package=mocks

type (
	colorizer interface {
		Red(in string) (out string)
		Green(in string) (out string)
		Yellow(in string) (out string)
		Blue(in string) (out string)
		Magenta(in string) (out string)
		Cyan(in string) (out string)
		Gray(in string) (out string)
	}
)
