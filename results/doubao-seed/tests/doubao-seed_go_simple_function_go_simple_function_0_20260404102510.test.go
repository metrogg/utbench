package inflight

import "testing"

type inflights struct {
	buffer []uint64
	start  int
	count  int
	size   int
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

func (in *inflights) getElements() []uint64 {
	res := make([]uint64, 0, in.count)
	idx := in.start
	for i := 0; i < in.count; i++ {
		res = append(res, in.buffer[idx])
		idx++
		if idx >= in.size {
			idx -= in.size
		}
	}
	return res
}

func TestInflightsFreeTo(t *testing.T) {
	tests := []struct {
		name      string
		in        *inflights
		to        uint64
		wantCnt   int
		wantStart int
		wantElems []uint64
	}{
		{
			name: "count zero no change",
			in: &inflights{
				count:  0,
				start:  2,
				size:   5,
				buffer: make([]uint64, 5),
			},
			to:        10,
			wantCnt:   0,
			wantStart: 2,
			wantElems: []uint64{},
		},
		{
			name: "to smaller than first element no change",
			in: &inflights{
				buffer: []uint64{2, 3, 4, 0, 0},
				start:  0,
				count:  3,
				size:   5,
			},
			to:        1,
			wantCnt:   3,
			wantStart: 0,
			wantElems: []uint64{2, 3, 4},
		},
		{
			name: "free partial elements no wrap",
			in: &inflights{
				buffer: []uint64{1, 2, 3, 4, 5},
				start:  0,
				count:  5,
				size:   5,
			},
			to:        3,
			wantCnt:   2,
			wantStart: 3,
			wantElems: []uint64{4, 5},
		},
		{
			name: "free all elements no wrap",
			in: &inflights{
				buffer: []uint64{1, 2, 3, 4, 5},
				start:  0,
				count:  5,
				size:   5,
			},
			to:        5,
			wantCnt:   0,
			wantStart: 0,
			wantElems: []uint64{},
		},
		{
			name: "free partial elements with wrap",
			in: &inflights{
				buffer: []uint64{4, 5, 1, 2, 3},
				start:  2,
				count:  5,
				size:   5,
			},
			to:        3,
			wantCnt:   2,
			wantStart: 0,
			wantElems: []uint64{4, 5},
		},
		{
			name: "free exactly first element",
			in: &inflights{
				buffer: []uint64{1, 2, 3, 0, 0},
				start:  0,
				count:  3,
				size:   5,
			},
			to:        1,
			wantCnt:   2,
			wantStart: 1,
			wantElems: []uint64{2, 3},
		},
		{
			name: "to larger than all elements free all",
			in: &inflights{
				buffer: []uint64{2, 4, 6, 8, 0},
				start:  0,
				count:  4,
				size:   5,
			},
			to:        100,
			wantCnt:   0,
			wantStart: 0,
			wantElems: []uint64{},
		},
		{
			name: "free up to middle element with wrap",
			in: &inflights{
				buffer: []uint64{3, 4, 1, 2, 0},
				start:  2,
				count:  4,
				size:   5,
			},
			to:        2,
			wantCnt:   2,
			wantStart: 0,
			wantElems: []uint64{3, 4},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.in.freeTo(tt.to)
			if tt.in.count != tt.wantCnt {
				t.Errorf("count mismatch: got %d, want %d", tt.in.count, tt.wantCnt)
			}
			if tt.in.start != tt.wantStart {
				t.Errorf("start mismatch: got %d, want %d", tt.in.start, tt.wantStart)
			}
			elems := tt.in.getElements()
			if len(elems) != len(tt.wantElems) {
				t.Errorf("elements length mismatch: got %d, want %d", len(elems), len(tt.wantElems))
				return
			}
			for i := range elems {
				if elems[i] != tt.wantElems[i] {
					t.Errorf("element at %d mismatch: got %d, want %d", i, elems[i], tt.wantElems[i])
				}
			}
		})
	}
}