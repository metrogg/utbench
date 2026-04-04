package inflight

import "testing"

type inflights struct {
	buffer []uint64
	start  int
	count  int
	size   int
}

func TestInflights_FreeTo(t *testing.T) {
	testCases := []struct {
		name          string
		initInflight  func() *inflights
		to            uint64
		expectedCnt   int
		expectedStart int
	}{
		{
			name: "empty inflight returns early no change",
			initInflight: func() *inflights {
				return &inflights{
					count:  0,
					start:  2,
					size:   5,
					buffer: make([]uint64, 5),
				}
			},
			to:            10,
			expectedCnt:   0,
			expectedStart: 2,
		},
		{
			name: "to smaller than first element returns early no change",
			initInflight: func() *inflights {
				return &inflights{
					buffer: []uint64{1, 2, 3, 0, 0},
					start:  0,
					count:  3,
					size:   5,
				}
			},
			to:            0,
			expectedCnt:   3,
			expectedStart: 0,
		},
		{
			name: "free partial elements no buffer wrap",
			initInflight: func() *inflights {
				return &inflights{
					buffer: []uint64{1, 2, 3, 4, 5},
					start:  0,
					count:  5,
					size:   5,
				}
			},
			to:            3,
			expectedCnt:   2,
			expectedStart: 3,
		},
		{
			name: "free all elements no wrap resets start to 0",
			initInflight: func() *inflights {
				return &inflights{
					buffer: []uint64{1, 2, 3, 4, 5},
					start:  0,
					count:  5,
					size:   5,
				}
			},
			to:            5,
			expectedCnt:   0,
			expectedStart: 0,
		},
		{
			name: "free partial elements with buffer wrap",
			initInflight: func() *inflights {
				return &inflights{
					buffer: []uint64{2, 3, 0, 1},
					start:  3,
					count:  3,
					size:   4,
				}
			},
			to:            2,
			expectedCnt:   1,
			expectedStart: 1,
		},
		{
			name: "free all elements with wrap resets start to 0",
			initInflight: func() *inflights {
				return &inflights{
					buffer: []uint64{2, 3, 0, 1},
					start:  3,
					count:  3,
					size:   4,
				}
			},
			to:            3,
			expectedCnt:   0,
			expectedStart: 0,
		},
		{
			name: "single element to equals value frees all",
			initInflight: func() *inflights {
				return &inflights{
					buffer: []uint64{5, 0},
					start:  0,
					count:  1,
					size:   2,
				}
			},
			to:            5,
			expectedCnt:   0,
			expectedStart: 0,
		},
		{
			name: "single element to smaller than value no change",
			initInflight: func() *inflights {
				return &inflights{
					buffer: []uint64{5, 0},
					start:  0,
					count:  1,
					size:   2,
				}
			},
			to:            4,
			expectedCnt:   1,
			expectedStart: 0,
		},
		{
			name: "to larger than all elements frees all",
			initInflight: func() *inflights {
				return &inflights{
					buffer: []uint64{1, 3, 5, 7},
					start:  0,
					count:  4,
					size:   4,
				}
			},
			to:            100,
			expectedCnt:   0,
			expectedStart: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			in := tc.initInflight()
			in.freeTo(tc.to)
			if in.count != tc.expectedCnt {
				t.Errorf("count mismatch: expected %d, got %d", tc.expectedCnt, in.count)
			}
			if in.start != tc.expectedStart {
				t.Errorf("start mismatch: expected %d, got %d", tc.expectedStart, in.start)
			}
		})
	}
}