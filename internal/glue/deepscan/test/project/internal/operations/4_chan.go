package operations

import "context"

type (
	message interface {
		Bytes(ctx context.Context) []byte
		Headers() []string
	}

	chFetcher interface {
		Fetch()
	}

	messageImpl struct {
	}
)

func InvisibleFetchFrom4() chan message {
	// invisible: code not depend on interface
	return nil
}

func InvisibleWriteTo4(buffer chan<- message) {
	// invisible: code not depend on interface
	buffer <- messageImpl{}
}

func InvisibleWriteToSpread4(buffers ...chan<- message) {
	// invisible: code not depend on interface
	for _, buff := range buffers {
		buff <- messageImpl{}
	}
}

func InvisibleWriteToSlice4(buffers []chan<- message) {
	// invisible: code not depend on interface
	for _, buff := range buffers {
		buff <- messageImpl{}
	}
}

func VisibleReadFrom4(buffer <-chan chFetcher) {
	// visible: code depend on interface
	for m := range buffer {
		_ = m.Fetch
	}
}

func VisibleReadFromSpread4(buffers ...<-chan chFetcher) {
	// visible: code depend on interface
	for _, buff := range buffers {
		for fetcher := range buff {
			fetcher.Fetch()
		}
	}
}

func VisibleReadFromSlice4(buffers []<-chan chFetcher) {
	// visible: code depend on interface
	for _, buff := range buffers {
		for fetcher := range buff {
			fetcher.Fetch()
		}
	}
}

func VisibleReadFromArray4(buffers [5]<-chan chFetcher) {
	// visible: code depend on interface
	for _, buff := range buffers {
		for fetcher := range buff {
			fetcher.Fetch()
		}
	}
}

func VisiblePotentiallyReadFrom4(buffer chan chFetcher) {
	// visible: code possible depend on interface
	// potentially we can read from it
	_ = buffer
}

func VisiblePotentiallyReadFromSpread4(buffer ...chan chFetcher) {
	// visible: code possible depend on interface
	// potentially we can read from it
	_ = buffer
}

func VisiblePotentiallyReadFromArray4(buffer [5]chan chFetcher) {
	// visible: code possible depend on interface
	// potentially we can read from it
	_ = buffer
}

func (m messageImpl) Bytes(ctx context.Context) []byte {
	return nil
}

func (m messageImpl) Headers() []string {
	return nil
}
