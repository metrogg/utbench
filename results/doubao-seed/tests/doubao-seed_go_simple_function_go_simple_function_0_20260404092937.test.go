package inflights

import "testing"

type inflights struct {
	count  int
	start  int
	size   int
	buffer []uint64
}

func TestInflights_FreeTo(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		in        *inflights
		freeTo    uint64
		wantCnt   int
		wantStart int
	}{
		{
			name: "count is zero, no changes",
			in: &inflights{
				count:  0,
				start:  2,
				size:   5,
				buffer: []uint64{1, 2, 3, 4, 5},
			},
			freeTo:    10,
			wantCnt:   0,
			wantStart: 2,
		},
		{
			name: "to smaller than first element, no changes",
			in: &inflights{
				count:  3,
				start:  0,
				size:   5,
				buffer: []uint64{5, 6, 7, 0, 0},
			},
			freeTo:    3,
			wantCnt:   3,
			wantStart: 0,
		},
		{
			name: "free first 2 elements, no wrap around",
			in: &inflights{
				count:  4,
				start:  0,
				size:   5,
				buffer: []uint64{1, 2, 3, 4, 0},
			},
			freeTo:    2,
			wantCnt:   2,
			wantStart: 2,
		},
		{
			name: "free all elements, start reset to 0",
			in: &inflights{
				count:  4,
				start:  0,
				size:   4,
				buffer: []uint64{1, 2, 3, 4},
			},
			freeTo:    4,
			wantCnt:   0,
			wantStart: 0,
		},
		{
			name: "free elements across wrap boundary",
			in: &inflights{
				count:  4,
				start:  3,
				size:   5,
				buffer: []uint64{7, 8, 0, 5, 6},
			},
			freeTo:    7,
			wantCnt:   1,
			wantStart: 1,
		},
		{
			name: "free exactly 1 element",
			in: &inflights{
				count:  2,
				start:  0,
				size:   2,
				buffer: []uint64{3, 5},
			},
			freeTo:    3,
			wantCnt:   1,
			wantStart: 1,
		},
		{
			name: "free all elements with non-zero initial start",
			in: &inflights{
				count:  2,
				start:  2,
				size:   5,
				buffer: []uint64{0, 0, 1, 2, 3},
			},
			freeTo:    2,
			wantCnt:   0,
			wantStart: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.in.freeTo(tt.freeTo)
			if tt.in.count != tt.wantCnt {
				t.Errorf("count mismatch: got %d, want %d", tt.in.count, tt.wantCnt)
			}
			if tt.in.start != tt.wantStart {
				t.Errorf("start mismatch: got %d, want %d", tt.in.start, tt.wantStart)
			}
		})
	}
}

// Method under test, included for standalone compilation (remove if testing against actual source)
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