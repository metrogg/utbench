package inflight

import "testing"

type inflights struct {
	count  int
	start  int
	size   int
	buffer []uint64
}

func (in *inflights) freeTo(to uint64) {
	if in.count == 0 || to < in.buffer[in.start] {
		return
	}

	idx := in.start
	var i int
	for i = 0; i < in.count; i++ {
		if to < in.buffer[idx] {
			break
		}

		size := in.size
		if idx++; idx >= size {
			idx -= size
		}
	}
	in.count -= i
	in.start = idx
	if in.count == 0 {
		in.start = 0
	}
}

func TestInflights_FreeTo(t *testing.T) {
	testCases := []struct {
		name          string
		initial       inflights
		to            uint64
		expectedCount int
		expectedStart int
	}{
		{
			name: "empty inflights no change",
			initial: inflights{
				count:  0,
				start:  2,
				size:   5,
				buffer: []uint64{1, 2, 3, 4, 5},
			},
			to:            10,
			expectedCount: 0,
			expectedStart: 2,
		},
		{
			name: "to less than first inflight no change",
			initial: inflights{
				count:  3,
				start:  0,
				size:   4,
				buffer: []uint64{2, 3, 4, 0},
			},
			to:            1,
			expectedCount: 3,
			expectedStart: 0,
		},
		{
			name: "free partial elements no wrap",
			initial: inflights{
				count:  4,
				start:  0,
				size:   4,
				buffer: []uint64{1, 2, 3, 4},
			},
			to:            2,
			expectedCount: 2,
			expectedStart: 2,
		},
		{
			name: "free all elements no wrap reset start",
			initial: inflights{
				count:  3,
				start:  1,
				size:   4,
				buffer: []uint64{0, 1, 2, 3},
			},
			to:            3,
			expectedCount: 0,
			expectedStart: 0,
		},
		{
			name: "free partial elements with wrap",
			initial: inflights{
				count:  3,
				start:  3,
				size:   4,
				buffer: []uint64{3, 4, 0, 2},
			},
			to:            3,
			expectedCount: 1,
			expectedStart: 1,
		},
		{
			name: "free all elements with wrap reset start",
			initial: inflights{
				count:  3,
				start:  3,
				size:   4,
				buffer: []uint64{3, 4, 0, 2},
			},
			to:            4,
			expectedCount: 0,
			expectedStart: 0,
		},
		{
			name: "to larger than all elements free all",
			initial: inflights{
				count:  4,
				start:  0,
				size:   5,
				buffer: []uint64{1, 3, 5, 7, 9},
			},
			to:            100,
			expectedCount: 0,
			expectedStart: 0,
		},
		{
			name: "free exactly first element",
			initial: inflights{
				count:  3,
				start:  0,
				size:   3,
				buffer: []uint64{5, 6, 7},
			},
			to:            5,
			expectedCount: 2,
			expectedStart: 1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			in := tc.initial
			in.freeTo(tc.to)
			if in.count != tc.expectedCount {
				t.Errorf("count mismatch: expected %d, got %d", tc.expectedCount, in.count)
			}
			if in.start != tc.expectedStart {
				t.Errorf("start mismatch: expected %d, got %d", tc.expectedStart, in.start)
			}
		})
	}
}