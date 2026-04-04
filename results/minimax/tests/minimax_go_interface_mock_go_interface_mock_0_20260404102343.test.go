package yourpackage

import "testing"

func TestInflightsFreeTo(t *testing.T) {
	t.Run("empty inflights returns early", func(t *testing.T) {
		in := &inflights{
			buffer: make([]uint64, 10),
			size:   10,
			count:  0,
			start:  0,
		}
		in.freeTo(5)
		// Nothing should change
		if in.count != 0 || in.start != 0 {
			t.Errorf("expected count=0, start=0, got count=%d, start=%d", in.count, in.start)
		}
	})

	t.Run("to less than first element returns early", func(t *testing.T) {
		in := &inflights{
			buffer: []uint64{5, 6, 7},
			size:   10,
			count:  3,
			start:  0,
		}
		in.freeTo(4)
		// Nothing should change since 4 < 5
		if in.count != 3 || in.start != 0 {
			t.Errorf("expected count=3, start=0, got count=%d, start=%d", in.count, in.start)
		}
	})

	t.Run("frees some elements", func(t *testing.T) {
		in := &inflights{
			buffer: []uint64{1, 2, 5, 6, 7},
			size:   10,
			count:  5,
			start:  0,
		}
		in.freeTo(3)
		// Should free elements 1 and 2, keep 5, 6, 7
		if in.count != 3 || in.buffer[0] != 5 {
			t.Errorf("expected count=3, first element=5, got count=%d, first=%d", in.count, in.buffer[0])
		}
	})

	t.Run("frees all elements", func(t *testing.T) {
		in := &inflights{
			buffer: []uint64{1, 2, 3},
			size:   10,
			count:  3,
			start:  0,
		}
		in.freeTo(10)
		// Should free all elements
		if in.count != 0 {
			t.Errorf("expected count=0, got count=%d", in.count)
		}
	})

	t.Run("single element freed", func(t *testing.T) {
		in := &inflights{
			buffer: []uint64{5},
			size:   10,
			count:  1,
			start:  0,
		}
		in.freeTo(6)
		// Should free the single element
		if in.count != 0 {
			t.Errorf("expected count=0, got count=%d", in.count)
		}
	})

	t.Run("wrap around case", func(t *testing.T) {
		// Simulate a circular buffer where start is in the middle
		buffer := make([]uint64, 5)
		buffer[3] = 1
		buffer[4] = 2
		buffer[0] = 5
		buffer[1] = 6
		
		in := &inflights{
			buffer: buffer,
			size:   5,
			count:  4,
			start:  3, // start is in the middle
		}
		in.freeTo(3)
		// Should free elements 1 and 2, keep 5 and 6
		if in.count != 2 {
			t.Errorf("expected count=2, got count=%d", in.count)
		}
		if in.start != 0 {
			t.Errorf("expected start=0, got start=%d", in.start)
		}
	})
}