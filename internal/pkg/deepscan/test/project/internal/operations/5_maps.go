package operations

type (
	processor interface {
		Process()
	}
)

func MapSimple5(processors map[string]processor) {
	_ = processors
}

func MapSlice5(processors map[string][]processor) {
	_ = processors
}

func MapChanInvisible5(processors map[string]chan<- processor) {
	// invisible: not depend on interface
	for _, c := range processors {
		c <- nil
	}
}

func MapChanVisible5(processors map[string]<-chan processor) {
	// visible: depend on interface
	for _, c := range processors {
		for p := range c {
			p.Process()
		}
	}
}

func MapChanPotentionallyVisible5(processors map[string]chan processor) {
	// visible: depend on interface
	_ = processors
}

func MapInside5(processors map[string]map[int]processor) {
	_ = processors
}
