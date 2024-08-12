package xpath

import "fmt"

func getType(x any) string {
	return fmt.Sprintf("%T", x)
}
