package inflights_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// inflights is a minimal struct to test the freeTo method
type inflights struct {
	buffer []uint64
	start  int
	count  int
	size   int
}

func (in *inflights) freeTo(to uint64) {
	if in.count == 0 || to < in.buffer[in.start] {
		// out of the left side of the window
		return
	}

	idx := in.start
	var i int
	for i = 0; i < in.count; i++ {
		if to < in.buffer[idx] { // found the first large inflight
			break
		}

		// increase index and maybe rotate
		size := in.size
		if idx++; idx >= size {
			idx -= size
		}
	}
	// free i inflights and set new start index
	in.count -= i
	in.start = idx
	if in.count == 0 {
		// inflights is empty, reset the start index so that we don't grow the
		// buffer unnecessarily.
		in.start = 0
	}
}

func TestFreeTo_EmptyInflights(t *testing.T) {
	in := &inflights{
		buffer: []uint64{1, 2, 3},
		start:  0,
		count:  0,
		size:   3,
	}

	originalStart := in.start
	originalCount := in.count

	in.freeTo(5)

	assert.Equal(t, originalStart, in.start, "start should not change when count is 0")
	assert.Equal(t, originalCount, in.count, "count should remain 0")
}

func TestFreeTo_ToLessThanFirstElement(t *testing.T) {
	in := &inflights{
		buffer: []uint64{10, 20, 30},
		start:  0,
		count:  3,
		size:   3,
	}

	originalStart := in.start
	originalCount := in.count

	in.freeTo(5)

	assert.Equal(t, originalStart, in.start, "start should not change when to < first element")
	assert.Equal(t, originalCount, in.count, "count should not change when to < first element")
}

func TestFreeTo_FreeAllElements(t *testing.T) {
	in := &inflights{
		buffer: []uint64{1, 2, 3, 4, 5},
		start:  0,
		count:  5,
		size:   5,
	}

	in.freeTo(10)

	assert.Equal(t, 0, in.count, "all elements should be freed")
	assert.Equal(t, 0, in.start, "start should be reset to 0 when empty")
}

func TestFreeTo_FreePartialElementsFromStart(t *testing.T) {
	in := &inflights{
		buffer: []uint64{1, 2, 3, 4, 5},
		start:  0,
		count:  5,
		size:   5,
	}

	in.freeTo(3)

	assert.Equal(t, 2, in.count, "should free first 3 elements")
	assert.Equal(t, 3, in.start, "start should move to index 3")
}

func TestFreeTo_FreePartialElementsWithWrapAround(t *testing.T) {
	in := &inflights{
		buffer: []uint64{5, 6, 1, 2, 3},
		start:  2, // points to value 1
		count:  5,
		size:   5,
	}

	in.freeTo(3)

	assert.Equal(t, 2, in.count, "should free elements 1,2,3")
	assert.Equal(t, 0, in.start, "start should wrap around to index 0")
}

func TestFreeTo_FreeExactlyToBoundary(t *testing.T) {
	in := &inflights{
		buffer: []uint64{1, 2, 3, 4, 5},
		start:  0,
		count:  5,
		size:   5,
	}

	in.freeTo(5)

	assert.Equal(t, 0, in.count, "should free all elements when to equals max")
	assert.Equal(t, 0, in.start, "start should reset to 0 when empty")
}

func TestFreeTo_FreeNoneWhenToBetweenElements(t *testing.T) {
	in := &inflights{
		buffer: []uint64{1, 3, 5, 7, 9},
		start:  0,
		count:  5,
		size:   5,
	}

	in.freeTo(0)

	assert.Equal(t, 5, in.count, "count should remain unchanged")
	assert.Equal(t, 0, in.start, "start should remain unchanged")
}

func TestFreeTo_FreePartialWithToBetweenElements(t *testing.T) {
	in := &inflights{
		buffer: []uint64{1, 3, 5, 7, 9},
		start:  0,
		count:  5,
		size:   5,
	}

	in.freeTo(4)

	assert.Equal(t, 3, in.count, "should free elements 1 and 3")
	assert.Equal(t, 2, in.start, "start should move to index 2 (value 5)")
}

func TestFreeTo_SingleElementFree(t *testing.T) {
	in := &inflights{
		buffer: []uint64{42},
		start:  0,
		count:  1,
		size:   1,
	}

	in.freeTo(42)

	assert.Equal(t, 0, in.count, "single element should be freed")
	assert.Equal(t, 0, in.start, "start should reset to 0")
}

func TestFreeTo_SingleElementNotFreed(t *testing.T) {
	in := &inflights{
		buffer: []uint64{42},
		start:  0,
		count:  1,
		size:   1,
	}

	in.freeTo(10)

	assert.Equal(t, 1, in.count, "single element should not be freed when to < element")
	assert.Equal(t, 0, in.start, "start should remain unchanged")
}

func TestFreeTo_ComplexWrapAroundScenario(t *testing.T) {
	in := &inflights{
		buffer: []uint64{8, 9, 1, 2, 3, 4, 5, 6, 7},
		start:  2, // points to value 1
		count:  9,
		size:   9,
	}

	in.freeTo(5)

	assert.Equal(t, 4, in.count, "should free elements 1,2,3,4,5")
	assert.Equal(t, 7, in.start, "start should move to index 7 (value 6)")
}

func TestFreeTo_ToLessThanAllWithWrapAround(t *testing.T) {
	in := &inflights{
		buffer: []uint64{5, 6, 7, 1, 2},
		start:  3, // points to value 1
		count:  5,
		size:   5,
	}

	in.freeTo(0)

	assert.Equal(t, 5, in.count, "no elements should be freed")
	assert.Equal(t, 3, in.start, "start should remain unchanged")
}

func TestFreeTo_BufferFullWrapAround(t *testing.T) {
	in := &inflights{
		buffer: []uint64{7, 8, 1, 2, 3, 4, 5, 6},
		start:  2, // points to value 1
		count:  8,
		size:   8,
	}

	in.freeTo(4)

	assert.Equal(t, 4, in.count, "should free elements 1,2,3,4")
	assert.Equal(t, 6, in.start, "start should move to index 6 (value 5)")
}