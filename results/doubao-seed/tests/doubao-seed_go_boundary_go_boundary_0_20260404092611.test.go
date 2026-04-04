package progress

import "testing"

type ProgressState int

const (
	ProgressStateReplicate ProgressState = 1
	ProgressStateProbe     ProgressState = 2
)

type Progress struct {
	State        ProgressState
	Match        uint64
	Next         uint64
	ResumeCalled bool
}

func (pr *Progress) resume() {
	pr.ResumeCalled = true
}

func min(a, b uint64) uint64 {
	if a < b {
		return a
	}
	return b
}

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
		name         string
		initialState ProgressState
		initialMatch uint64
		initialNext  uint64
		rejected     uint64
		last         uint64
		wantRet      bool
		wantNext     uint64
		wantResume   bool
	}{
		{
			name:         "ReplicateState_RejectedLessThanMatch_ReturnsFalse",
			initialState: ProgressStateReplicate,
			initialMatch: 5,
			initialNext:  10,
			rejected:     3,
			last:         0,
			wantRet:      false,
			wantNext:     10,
			wantResume:   false,
		},
		{
			name:         "ReplicateState_RejectedEqualToMatch_ReturnsFalse",
			initialState: ProgressStateReplicate,
			initialMatch: 5,
			initialNext:  10,
			rejected:     5,
			last:         0,
			wantRet:      false,
			wantNext:     10,
			wantResume:   false,
		},
		{
			name:         "ReplicateState_RejectedGreaterThanMatch_ReturnsTrueAdjustsNext",
			initialState: ProgressStateReplicate,
			initialMatch: 5,
			initialNext:  10,
			rejected:     6,
			last:         0,
			wantRet:      true,
			wantNext:     6,
			wantResume:   false,
		},
		{
			name:         "NonReplicateState_RejectedNotMatchNextMinus1_ReturnsFalse",
			initialState: ProgressStateProbe,
			initialMatch: 3,
			initialNext:  10,
			rejected:     8,
			last:         10,
			wantRet:      false,
			wantNext:     10,
			wantResume:   false,
		},
		{
			name:         "NonReplicateState_RejectedMatchesNextMinus1_RejectedLTLastPlus1",
			initialState: ProgressStateProbe,
			initialMatch: 3,
			initialNext:  10,
			rejected:     9,
			last:         10,
			wantRet:      true,
			wantNext:     9,
			wantResume:   true,
		},
		{
			name:         "NonReplicateState_RejectedMatchesNextMinus1_LastPlus1LTRejected",
			initialState: ProgressStateProbe,
			initialMatch: 3,
			initialNext:  10,
			rejected:     9,
			last:         7,
			wantRet:      true,
			wantNext:     8,
			wantResume:   true,
		},
		{
			name:         "NonReplicateState_RejectedMatchesNextMinus1_MinLessThan1_SetTo1",
			initialState: ProgressStateProbe,
			initialMatch: 0,
			initialNext:  1,
			rejected:     0,
			last:         0,
			wantRet:      true,
			wantNext:     1,
			wantResume:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pr := &Progress{
				State: tt.initialState,
				Match: tt.initialMatch,
				Next:  tt.initialNext,
			}
			gotRet := pr.maybeDecrTo(tt.rejected, tt.last)
			if gotRet != tt.wantRet {
				t.Errorf("maybeDecrTo() returned %v, want %v", gotRet, tt.wantRet)
			}
			if pr.Next != tt.wantNext {
				t.Errorf("maybeDecrTo() Next = %d, want %d", pr.Next, tt.wantNext)
			}
			if pr.ResumeCalled != tt.wantResume {
				t.Errorf("maybeDecrTo() ResumeCalled = %v, want %v", pr.ResumeCalled, tt.wantResume)
			}
			if pr.Match != tt.initialMatch {
				t.Errorf("maybeDecrTo() Match changed unexpectedly from %d to %d", tt.initialMatch, pr.Match)
			}
		})
	}
}