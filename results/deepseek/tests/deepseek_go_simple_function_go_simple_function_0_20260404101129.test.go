package inflights_test

import (
	"testing"
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

func TestFreeTo(t *testing.T) {
	tests := []struct {
		name     string
		buffer   []uint64
		start    int
		count    int
		size     int
		to       uint64
		expected struct {
			count int
			start int
		}
	}{
		// Normal cases
		{
			name:   "free all when to is greater than all elements",
			buffer: []uint64{1, 2, 3, 4, 5},
			start:  0,
			count:  5,
			size:   5,
			to:     10,
			expected: struct {
				count int
				start int
			}{count: 0, start: 0},
		},
		{
			name:   "free first three elements",
			buffer: []uint64{1, 2, 3, 4, 5},
			start:  0,
			count:  5,
			size:   5,
			to:     3,
			expected: struct {
				count int
				start int
			}{count: 2, start: 3},
		},
		{
			name:   "free none when to is less than first element",
			buffer: []uint64{5, 6, 7, 8, 9},
			start:  0,
			count:  5,
			size:   5,
			to:     3,
			expected: struct {
				count int
				start int
			}{count: 5, start: 0},
		},
		{
			name:   "free exactly to boundary element",
			buffer: []uint64{1, 2, 3, 4, 5},
			start:  0,
			count:  5,
			size:   5,
			to:     3,
			expected: struct {
				count int
				start int
			}{count: 2, start: 3},
		},
		// Boundary conditions
		{
			name:   "empty buffer",
			buffer: []uint64{},
			start:  0,
			count:  0,
			size:   0,
			to:     5,
			expected: struct {
				count int
				start int
			}{count: 0, start: 0},
		},
		{
			name:   "single element buffer free all",
			buffer: []uint64{10},
			start:  0,
			count:  1,
			size:   1,
			to:     10,
			expected: struct {
				count int
				start int
			}{count: 0, start: 0},
		},
		{
			name:   "single element buffer free none",
			buffer: []uint64{10},
			start:  0,
			count:  1,
			size:   1,
			to:     5,
			expected: struct {
				count int
				start int
			}{count: 1, start: 0},
		},
		// Circular buffer cases
		{
			name:   "circular buffer free across boundary",
			buffer: []uint64{5, 6, 1, 2, 3},
			start:  2,
			count:  5,
			size:   5,
			to:     2,
			expected: struct {
				count int
				start int
			}{count: 2, start: 4},
		},
		{
			name:   "circular buffer free all across boundary",
			buffer: []uint64{5, 6, 1, 2, 3},
			start:  2,
			count:  5,
			size:   5,
			to:     6,
			expected: struct {
				count int
				start int
			}{count: 0, start: 0},
		},
		{
			name:   "circular buffer free none when to less than start",
			buffer: []uint64{5, 6, 1, 2, 3},
			start:  2,
			count:  5,
			size:   5,
			to:     0,
			expected: struct {
				count int
				start int
			}{count: 5, start: 2},
		},
		// Edge cases
		{
			name:   "to equals first element",
			buffer: []uint64{5, 6, 7, 8, 9},
			start:  0,
			count:  5,
			size:   5,
			to:     5,
			expected: struct {
				count int
				start int
			}{count: 4, start: 1},
		},
		{
			name:   "partial free with wrap-around index calculation",
			buffer: []uint64{8, 9, 5, 6, 7},
			start:  2,
			count:  5,
			size:   5,
			to:     6,
			expected: struct {
				count int
				start int
			}{count: 2, start: 4},
		},
		{
			name:   "free all elements exactly",
			buffer: []uint64{1, 2, 3, 4, 5},
			start:  0,
			count:  5,
			size:   5,
			to:     5,
			expected: struct {
				count int
				start int
			}{count: 0, start: 0},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			in := &inflights{
				buffer: tt.buffer,
				start:  tt.start,
				count:  tt.count,
				size:   tt.size,
			}

			in.freeTo(tt.to)

			if in.count != tt.expected.count {
				t.Errorf("count = %d, expected %d", in.count, tt.expected.count)
			}
			if in.start != tt.expected.start {
				t.Errorf("start = %d, expected %d", in.start, tt.expected.start)
			}
		})
	}
}