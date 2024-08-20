package container

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/fe3dback/go-arch-lint-sdk/pkg/ptr"
)

func getPointerToOne() *int {
	return once(func() *int {
		return ptr.Ref(1)
	})
}

func getPointerToTwo() *int {
	return once(func() *int {
		return ptr.Ref(2)
	})
}

func Test_doOnce_ReturnsCachedInstance(t *testing.T) {
	pointerToOne := getPointerToOne()
	pointerToTwo := getPointerToTwo()
	anotherPointerToOne := getPointerToOne()

	require.Equal(t, pointerToOne, anotherPointerToOne)
	require.NotEqual(t, pointerToOne, pointerToTwo)
}

func Test_doOnce_RaceCondition(t *testing.T) {
	const cycles = 1000

	var mux sync.Mutex
	var wg sync.WaitGroup
	wg.Add(cycles)

	originalPointerToOne := getPointerToOne()
	originalPointerToTwo := getPointerToTwo()

	successCycles := 0

	for i := 0; i < cycles; i++ {
		i := i

		go func() {
			defer wg.Done()

			if i%2 == 0 {
				pointerToOne := getPointerToOne()
				require.Equal(t, originalPointerToOne, pointerToOne)
			} else {
				pointerToTwo := getPointerToTwo()
				require.Equal(t, originalPointerToTwo, pointerToTwo)
			}

			mux.Lock()
			successCycles += 1
			mux.Unlock()
		}()
	}

	wg.Wait()
	require.Equal(t, cycles, successCycles)
}
