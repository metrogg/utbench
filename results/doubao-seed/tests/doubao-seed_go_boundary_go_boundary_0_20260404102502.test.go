package progress

import "testing"

type ProgressState int

const (
	ProgressStateReplicate ProgressState = iota
	ProgressStateProbe
	ProgressStateSnapshot
)

type Progress struct {
	State        ProgressState
	Match        uint64
	Next         uint64
	resumeCalled bool
}

func (pr *Progress) resume() {
	pr.resumeCalled = true
}

func min(a, b uint64) uint64 {
	if a < b {
		return a
	}
	return b
}

// Function under test (included for test runnability)
func (pr *Progress) maybeDecrTo(rejected, last uint64) bool {
	if pr.State == ProgressStateReplicate {
		if rejected <= pr.Match {
			return false
		}
		pr.Next = pr.Match + 1
		return true
	}

	if pr.Next-1 != rejected {
		return false
	}

	if pr.Next = min(rejected, last+1); pr.Next < 1 {
		pr.Next = 1
	}
	pr.resume()
	return true
}

func TestProgress_maybeDecrTo(t *testing.T) {
	tests := []struct {
		name             string
		setupPr          func() *Progress
		rejected         uint64
		last             uint64
		wantReturn       bool
		wantNext         uint64
		wantResumeCalled bool
	}{
		{
			name: "ReplicateState_RejectedLessThanMatch_ReturnsFalse",
			setupPr: func() *Progress {
				return &Progress{State: ProgressStateReplicate, Match: 5, Next: 10}
			},
			rejected:         3,
			last:             0,
			wantReturn:       false,
			wantNext:         10,
			wantResumeCalled: false,
		},
		{
			name: "ReplicateState_RejectedEqualsMatch_ReturnsFalse",
			setupPr: func() *Progress {
				return &Progress{State: ProgressStateReplicate, Match: 5, Next: 10}
			},
			rejected:         5,
			last:             0,
			wantReturn:       false,
			wantNext:         10,
			wantResumeCalled: false,
		},
		{
			name: "ReplicateState_RejectedGreaterThanMatch_ReturnsTrueAdjustsNext",
			setupPr: func() *Progress {
				return &Progress{State: ProgressStateReplicate, Match: 5, Next: 10}
			},
			rejected:         6,
			last:             0,
			wantReturn:       true,
			wantNext:         6,
			wantResumeCalled: false,
		},
		{
			name: "NonReplicateState_RejectedNotMatchNextMinus1_ReturnsFalse",
			setupPr: func() *Progress {
				return &Progress{State: ProgressStateProbe, Next: 10}
			},
			rejected:         8,
			last:             10,
			wantReturn:       false,
			wantNext:         10,
			wantResumeCalled: false,
		},
		{
			name: "NonReplicateState_RejectedMatchNextMinus1_RejectedLessThanLastPlus1_SetNextToRejected",
			setupPr: func() *Progress {
				return &Progress{State: ProgressStateProbe, Next: 10}
			},
			rejected:         9,
			last:             10,
			wantReturn:       true,
			wantNext:         9,
			wantResumeCalled: true,
		},
		{
			name: "NonReplicateState_RejectedMatchNextMinus1_RejectedGreaterThanLastPlus1_SetNextToLastPlus1",
			setupPr: func() *Progress {
				return &Progress{State: ProgressStateSnapshot, Next: 10}
			},
			rejected:         9,
			last:             7,
			wantReturn:       true,
			wantNext:         8,
			wantResumeCalled: true,
		},
		{
			name: "NonReplicateState_AdjustedNextLessThan1_ClampTo1",
			setupPr: func() *Progress {
				return &Progress{State: ProgressStateProbe, Next: 1}
			},
			rejected:         0,
			last:             0,
			wantReturn:       true,
			wantNext:         1,
			wantResumeCalled: true,
		},
		{
			name: "NonReplicateState_AdjustedNextEquals1_NoClamp",
			setupPr: func() *Progress {
				return &Progress{State: ProgressStateSnapshot, Next: 2}
			},
			rejected:         1,
			last:             0,
			wantReturn:       true,
			wantNext:         1,
			wantResumeCalled: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pr := tt.setupPr()
			gotReturn := pr.maybeDecrTo(tt.rejected, tt.last)
			if gotReturn != tt.wantReturn {
				t.Errorf("maybeDecrTo() returned %v, want %v", gotReturn, tt.wantReturn)
			}
			if pr.Next != tt.wantNext {
				t.Errorf("maybeDecrTo() Next = %v, want %v", pr.Next, tt.wantNext)
			}
			if pr.resumeCalled != tt.wantResumeCalled {
				t.Errorf("maybeDecrTo() resume called = %v, want %v", pr.resumeCalled, tt.wantResumeCalled)
			}
		})
	}
}