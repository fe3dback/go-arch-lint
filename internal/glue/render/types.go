package render

type (
	ColorPrinter interface {
		Red(in string) (out string)
		Green(in string) (out string)
		Yellow(in string) (out string)
		Blue(in string) (out string)
		Magenta(in string) (out string)
		Cyan(in string) (out string)
	}
)
