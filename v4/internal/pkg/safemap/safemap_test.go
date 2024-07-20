package safemap_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/fe3dback/go-arch-lint/v4/internal/pkg/safemap"
)

func TestMap_Set(t *testing.T) {
	sm := safemap.New[int, string]()

	// check key not found
	data, exist := sm.Get(5)
	assert.False(t, exist)
	assert.Equal(t, "", data)

	// we add some key, check that key exist now
	sm.Set(5, "some-value")
	data, exist = sm.Get(5)
	assert.True(t, exist)
	assert.Equal(t, "some-value", data)

	// then change, check that value is changed
	sm.Set(5, "another-value")
	data, exist = sm.Get(5)
	assert.True(t, exist)
	assert.Equal(t, "another-value", data)

	// check iterator
	totalKeys := 0
	keysAggregatedSum := 0
	valuesAggregatedRow := ""

	sm.Iterate(func(key int, value string) {
		totalKeys++
		keysAggregatedSum += key
		valuesAggregatedRow = fmt.Sprintf("%s,%s", valuesAggregatedRow, value)
	})

	assert.Equal(t, 1, totalKeys)
	assert.Equal(t, 5, keysAggregatedSum)
	assert.Equal(t, ",another-value", valuesAggregatedRow)
}
