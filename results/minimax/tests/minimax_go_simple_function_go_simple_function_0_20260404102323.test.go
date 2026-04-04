package yourpackage

import (
	"testing"
)

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

func TestFreeTo_EmptyBuffer(t *testing.T) {
	in := &inflights{
		buffer: []uint64{1, 2, 3},
		start:  0,
		count:  0,
		size:   3,
	}

	in.freeTo(100)

	if in.count != 0 {
		t.Errorf("expected count to be 0, got %d", in.count)
	}
	if in.start != 0 {
		t.Errorf("expected start to be 0, got %d", in.start)
	}
}

func TestFreeTo_ToLessThanFirstElement(t *testing.T) {
	in := &inflights{
		buffer: []uint64{10, 20, 30},
		start:  0,
		count:  3,
		size:   3,
	}

	in.freeTo(5)

	if in.count != 3 {
		t.Errorf("expected count to remain 3, got %d", in.count)
	}
	if in.start != 0 {
		t.Errorf("expected start to remain 0, got %d", in.start)
	}
}

func TestFreeTo_FreeAllElements(t *testing.T) {
	in := &inflights{
		buffer: []uint64{10, 20, 30},
		start:  0,
		count:  3,
		size:   3,
	}

	in.freeTo(30)

	if in.count != 0 {
		t.Errorf("expected count to be 0, got %d", in.count)
	}
	if in.start != 0 {
		t.Errorf("expected start to be reset to 0, got %d", in.start)
	}
}

func TestFreeTo_FreeAllElementsWithToGreaterThanMax(t *testing.T) {
	in := &inflights{
		buffer: []uint64{10, 20, 30},
		start:  0,
		count:  3,
		size:   3,
	}

	in.freeTo(100)

	if in.count != 0 {
		t.Errorf("expected count to be 0, got %d", in.count)
	}
}

func TestFreeTo_FreePartialElements(t *testing.T) {
	in := &inflights{
		buffer: []uint64{10, 20, 30, 40, 50},
		start:  0,
		count:  5,
		size:   5,
	}

	in.freeTo(25)

	if in.count != 3 {
		t.Errorf("expected count to be 3, got %d", in.count)
	}
	if in.start != 2 {
		t.Errorf("expected start to be 2, got %d", in.start)
	}
}

func TestFreeTo_FreeFirstElementOnly(t *testing.T) {
	in := &inflights{
		buffer: []uint64{10, 20, 30},
		start:  0,
		count:  3,
		size:   3,
	}

	in.freeTo(10)

	if in.count != 2 {
		t.Errorf("expected count to be 2, got %d", in.count)
	}
	if in.start != 1 {
		t.Errorf("expected start to be 1, got %d", in.start)
	}
}

func TestFreeTo_CircularBufferRotation(t *testing.T) {
	in := &inflights{
		buffer: []uint64{30, 40, 10, 20},
		start:  2,
		count:  2,
		size:   4,
	}

	in.freeTo(20)

	if in.count != 0 {
		t.Errorf("expected count to be 0, got %d", in.count)
	}
	if in.start != 0 {
		t.Errorf("expected start to be reset to 0, got %d", in.start)
	}
}

func TestFreeTo_CircularBufferPartialFree(t *testing.T) {
	in := &inflights{
		buffer: []uint64{30, 40, 10, 20, 50, 60},
		start:  2,
		count:  4,
		size:   6,
	}

	in.freeTo(25)

	if in.count != 2 {
		t.Errorf("expected count to be 2, got %d", in.count)
	}
	if in.start != 4 {
		t.Errorf("expected start to be 4, got %d", in.start)
	}
}

func TestFreeTo_BufferSizeOne(t *testing.T) {
	in := &inflights{
		buffer: []uint64{42},
		start:  0,
		count:  1,
		size:   1,
	}

	in.freeTo(42)

	if in.count != 0 {
		t.Errorf("expected count to be 0, got %d", in.count)
	}
}

func TestFreeTo_BufferSizeOne_NotFreed(t *testing.T) {
	in := &inflights{
		buffer: []uint64{42},
		start:  0,
		count:  1,
		size:   1,
	}

	in.freeTo(10)

	if in.count != 1 {
		t.Errorf("expected count to be 1, got %d", in.count)
	}
}

func TestFreeTo_LargeCircularBuffer(t *testing.T) {
	size := 10
	buffer := make([]uint64, size)
	for i := 0; i < size; i++ {
		buffer[i] = uint64((i + 5) * 10)
	}

	in := &inflights{
		buffer: buffer,
		start:  5,
		count:  5,
		size:   size,
	}

	in.freeTo(80)

	if in.count != 2 {
		t.Errorf("expected count to be 2, got %d", in.count)
	}
	if in.start != 8 {
		t.Errorf("expected start to be 8, got %d", in.start)
	}
}

func TestFreeTo_ToEqualToLastElement(t *testing.T) {
	in := &inflights{
		buffer: []uint64{10, 20, 30},
		start:  0,
		count:  3,
		size:   3,
	}

	in.freeTo(20)

	if in.count != 1 {
		t.Errorf("expected count to be 1, got %d", in.count)
	}
	if in.start != 2 {
		t.Errorf("expected start to be 2, got %d", in.start)
	}
}

func TestFreeTo_AllElementsSameValue(t *testing.T) {
	in := &inflights{
		buffer: []uint64{50, 50, 50, 50},
		start:  0,
		count:  4,
		size:   4,
	}

	in.freeTo(50)

	if in.count != 0 {
		t.Errorf("expected count to be 0, got %d", in.count)
	}
}

func TestFreeTo_CountBecomesZeroResetsStart(t *testing.T) {
	in := &inflights{
		buffer: []uint64{10, 20, 30},
		start:  1,
		count:  2,
		size:   3,
	}

	in.freeTo(30)

	if in.count != 0 {
		t.Errorf("expected count to be 0, got %d", in.count)
	}
	if in.start != 0 {
		t.Errorf("expected start to be reset to 0, got %d", in.start)
	}
}