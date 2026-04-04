package inflights

import (
	"testing"
)

// inflights is a circular buffer to track messages in flight.
// This struct definition is inferred from the usage in the source code.
type inflights struct {
	buffer []uint64
	start  int
	count  int
	size   int
}

// freeTo frees the inflights up to the given index (excluding the index itself if it matches).
// Source code provided in the prompt.
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

func TestInflights_freeTo(t *testing.T) {
	tests := []struct {
		name     string
		initial  inflights
		to       uint64
		expected inflights
	}{
		{
			name: "Empty buffer",
			initial: inflights{
				buffer: make([]uint64, 3),
				size:   3,
			},
			to: 10,
			expected: inflights{
				buffer: make([]uint64, 3),
				size:   3,
			},
		},
		{
			name: "To is less than start element",
			initial: inflights{
				buffer: []uint64{10, 20, 30},
				start:  0,
				count:  3,
				size:   3,
			},
			to: 5,
			expected: inflights{
				buffer: []uint64{10, 20, 30},
				start:  0,
				count:  3,
				size:   3,
			},
		},
		{
			name: "Free one element",
			initial: inflights{
				buffer: []uint64{10, 20, 30},
				start:  0,
				count:  3,
				size:   3,
			},
			to: 10,
			expected: inflights{
				buffer: []uint64{10, 20, 30},
				start:  1,
				count:  2,
				size:   3,
			},
		},
		{
			name: "Free multiple elements",
			initial: inflights{
				buffer: []uint64{10, 20, 30},
				start:  0,
				count:  3,
				size:   3,
			},
			to: 20,
			expected: inflights{
				buffer的后 []uint64{10, 20, 30},
				start:  2,
				count:  1,
				size:   3,
			},
		},
		{
			name: "Free all elements",
			initial: inflights{
				buffer: []uint64{10, 20, 30},
				start:  0,
				count:  3,
				size:   3,
			},
			to: 100,
			expected: inflights{
				buffer: []uint64{10, 20, 30},
				start:  0,
				count:  0,
				size:   3,
			},
		},
		{
			name: "Circular buffer - free partial",
			initial: inflights{
				buffer: []uint64{0, 0, 30, 40, 50},
				start:  2,
				count:  3,
				size:   5,
			},
			to: 45,
			expected: inflights{
				buffer: []uint64{0, 0, 30, 40, 50},
				start:  4,
				count:  1,
				size:   5,
			},
		},
		{
			name: "Circular buffer - free all",
			initial: inflights{
				buffer: []uint64{0, 0, 30, 40, 50},
				start:  2,
				count:  3,
				size:   5,
			},
			to: 100,
			expected: inflights{
				buffer: []uint64{0, 0, 30, 40, 50},
				start:  0,
				count:  0,
				size:   5,
			},
		},
		{
			name: "Circular buffer - wrap around logic",
			initial: inflights{
				buffer: []uint64{20, 30, 10},
				start:  2,
				count:  3,
				size:   3,
			},
			to: 25,
			expected: inflights{
				buffer: []uint64{20, 30, 10},
				start:  1,
				count:  1,
				size:   3,
			},
		},
		{
			name: "Circular buffer - wrap around and free all",
			initial: inflights{
				buffer: []uint64{20, 30, 10},
				start:  2,
				count:  3,
				size:   3,
			},
			to: 100,
			expected: inflights{
				buffer: []uint64{20, 30, 10},
				start:  0,
				count:  0,
				size:   3,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Make a copy of initial to ensure we don't modify the test case definition
			in := tt.initial
			in.freeTo(tt.to)

			if in.count != tt.expected.count {
				t.Errorf("count = %d, want %d", in.count, tt.expected.count)
			}
			if in.start != tt.expected.start {
				t.Errorf("start = %d, want %d", in.start, tt.expected.start)
			}
			if in.size != tt.expected.size {
				t.Errorf("size = %d, want %d", in.size, tt.expected.size)
			}
			// buffer content is not modified by freeTo, but we check length just in case
			if len(in.buffer) != len(tt.expected.buffer) {
				t.Errorf("len(buffer) = %d, want %d", len(in.buffer), len(tt.expected.buffer))
			}
		})
	}
}