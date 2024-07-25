//go:build none

package tests

import "time"

type someCode struct {
	hello string
	world time.Time
}

var test = "world"

func newSomeCode() *someCode {
	return &someCode{
		hello: "world!123",
		world: time.Now(),
	}
}
