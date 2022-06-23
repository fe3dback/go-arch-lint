package operations

type (
	rcvProcessor struct {
	}

	rcvRepository interface {
		SayHello()
	}
)

func (rcv *rcvProcessor) DoWork3(r rcvRepository) {
	r.SayHello()
}

func (rcv rcvProcessor) DoCopyWork3(r rcvRepository) {
	r.SayHello()
}
