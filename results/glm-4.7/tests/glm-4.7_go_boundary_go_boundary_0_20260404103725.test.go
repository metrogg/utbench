package progress

import (
	"math"
	"testing"
)

// Constants for Progress State
const (
	ProgressStateProbe     = 0
	ProgressStateReplicate = 1
	ProgressStateSnapshot  = 2
)

// Progress struct definition to support the method under test
type Progress struct {
	State uint64
	Match uint64
	Next  uint64
	// resumed is used to mock the behavior of resume() and verify it was called
	resumed bool
}

// resume is a helper method mocked for testing
func (pr *Progress) resume() {
	pr.resumed = true
}

// maybeDecrTo is the function under test
func (pr *Progress) maybeDecrTo(rejected, last uint64) bool {
	if pr.State == ProgressStateReplicate {
		// the rejection must be stale if the progress has matched and "rejected"
		// is smaller than "match".
		if rejected <= pr.Match {
			return false
		}
		// directly decrease next to match + 1
		pr.Next = pr.Match + 1
		return true
	}

	// the rejection must be stale if "rejected" does not match next - 1
	if pr.Next-1 != rejected {
		return false
	}

	if pr.Next = min(rejected, last+1); pr.Next < 1 {
		pr.Next = 1
	}
	pr.resume()
	return true
}

func TestMaybeDecrTo(t *testing.T) {
	tests := []struct {
		name         string
		state        uint64
		match        uint64
		next         uint64
		rejected     uint64
		last         uint64
		wantOk       bool
		wantNext     uint64
		wantResumed  bool
	}{
		{
			name:    "ReplicateState: rejected is stale (less than match)",
			state:   ProgressStateReplicate,
			match:   10,
			next:    20,
			rejected: 5,
			last:    100,
			wantOk:  false,
			wantNext: 20,
		},
		{
			name:    "ReplicateState: rejected is stale (equal to match)",
			state:   ProgressStateReplicate,
			match:   10,
			next:    20,
			rejected: 10,
			last:    100,
			wantOk:  false,
			wantNext: 20,
		},
		{
			name:    "ReplicateState: rejected is fresh, decrease Next to Match + 1",
			state:   ProgressStateReplicate,
			match:   10,
			next:    20,
			rejected: 15,
			last:    100,
			wantOk:  true,
			wantNext: 11,
		},
		{
			name:    "ProbeState: rejected is stale (does not match Next - 1)",
			state:   ProgressStateProbe,
			match:   0,
			next:    10,
			rejected: 5,
			last:    100,
			wantOk:  false,
			wantNext: 10,
		},
		{
			name:        "ProbeState: rejected matches Next - 1, rejected < last+1",
			state:       ProgressStateProbe,
			matchNext:   10,
			next:        10,
			rejected:    9,
			last:        100,
			wantOk:      true,
			wantNext:    9,
			wantResumed: true,
		},
		{
			name:        "ProbeState: rejected matches Next - 1, rejected > last+1",
			state:       ProgressStateProbe,
			match:       0,
			next:        20,
			rejected:    19,
			last:        15,
			wantOk:      true,
			wantNext:    16,
			wantResumed: true,
		},
		{
			name:        "ProbeState: rejected matches Next - 1, rejected == last+1",
			state:       ProgressStateProbe,
			match:       0,
			next:        20,
			rejected:    19,
			last:        18,
			wantOk:      true,
			wantNext:    19,
			wantResumed: true,
		},
		{
			name:        "Boundary: min result is 0, Next should be set to 1",
			state:       ProgressStateProbe,
			match:       0,
			next:        1, // Next-1 = 0
			rejected:    0,
			last:        math.MaxUint64, // last+1 overflows to 0
			wantOk:      true,
			wantNext:    1,
			wantResumed: true,
		},
		{
			name:        "Boundary: Next is 1, rejected is 0, last is 0",
			state:       ProgressStateProbe,
			match:       0,
			next:        1,
			rejected:    0,
			last:        0,
			wantOk:      true,
			wantNext:    1,
			wantResumed: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, state(t *testing.T) {
			pr := &Progress{
				State: tt.state,
				Match: tt.match,
				Next:  tt.next,
			}

			got := pr.maybeDecrTo(tt.rejected, tt.last)

			if got != tt.wantOk {
				t.Errorf("maybeDecrTo() returned = %v, want %v", got, tt.wantOk)
			}
			if pr.Next != tt.wantNext {
				t.Errorf("Next = %v, want %v", pr.Next, tt.wantNext)
			}
			if pr.resumed != tt.wantResumed {
				t.Errorf("resume() called = %v, want %v", pr.resumed, tt.wantResumed)
			}
		})
	}
}