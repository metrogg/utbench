package inflights

import (
	"testing"
)

// Mocking the struct definition to ensure the test is self-contained and runnable.
type inflights struct {
	buffer []uint64
	start  int
	count  int
	size   int
}

// The function under test.
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
StartIdxCheck:
	if in.count == 0 {
		// inflights is empty, reset the start index so that we don't grow the
		// buffer unnecessarily.
		in.start = 0
	}
}

func TestInflights_freeTo(t *testing.T) {
	tests := []struct {
		name     string
		initial  *inflights
		to       uint64
		wantCnt  int
		wantStart int
		wantFirstVal uint64 // Expected value at buffer[start] if count > 0
	}{
		{
			name: "empty inflights",
			initial: &inflights{
				buffer: make([]uint64, 10),
				start:  0,
				count:  0,
				size:   10,
			},
			to:       5,
			wantCnt:  0,
			wantStart: 0,
		},
		{
			name: "to is out of left side of window",
			initial: &inflights{
				buffer: []uint64{10, 20, 30, 0, 0},
				start:  0,
				count:  3,
				size:   5,
			},
			to:       5,
			wantCnt:  3,
			wantStart: 0,
			wantFirstVal: 10,
		},
		{
			name: "free some elements (linear buffer)",
			initial: &inflights{
				buffer: []uint64{1, 2, 3, 4, 0},
				start:  0,
				count:  4,
				size:   5,
			},
			to:       2, // Frees 1, 2. Keeps 3, 4.
			wantCnt:  2,
			wantStart: 2, // Index of 3
			wantFirstVal: 3,
		},
		{
			name: "free all elements (linear buffer)",
			initial: &inflights{
				buffer: []uint64{1, 2, 3, 0, 0},
				start:  0,
				count:  3,
				size:   5,
			},
			to:       100,
			wantCnt:  0,
			wantStart: 0, // Reset to 0 when empty
		},
		{
			name: "free some elements (circular buffer)",
			// Buffer: [5, 6, 1, 2, 3, 4]
			// Start: 2 (points to 1)
			// Count: 6
			// Logical: 1, 2, 3, 4, 5, 6
			initial: &inflights{
				buffer: []uint64{5, 6, 1, 2, 3, 4},
				start:  2,
				count:  6,
				size:   6,
			},
			to:       2, // Frees 1, 2. Keeps 3, 4, 5, 6.
			wantCnt:  4,
			wantStart: 4, // Index of 3
			wantFirstVal: 3,
[...]