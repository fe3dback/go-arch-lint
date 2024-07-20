//go:build tools

package tools

//go:generate go build -o ../bin/mockgen go.uber.org/mock/mockgen

import (
	_ "go.uber.org/mock/mockgen"
)
