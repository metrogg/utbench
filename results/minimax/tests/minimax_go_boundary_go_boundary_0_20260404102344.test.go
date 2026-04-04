package raft

import (
	"testing"
)

type Progress struct {
	State StateType
	Match uint64
	Next  uint64
}

func TestProgress_maybeDecrTo_ProgressStateReplicate(t *testing.T) {
	tests := []struct {
		name     string
		pr       *Progress
		rejected uint64
		want     bool
		wantNext uint64
	}{
		{
			name: "rejected less than match should return false",
			pr: &Progress{
				State: ProgressStateReplicate,
				Match: 10,
				Next:  11,
			},
			rejected: 5,
			want:     false,
		},
		{
			name: "rejected equal to match should return false",
			pr: &Progress{
				State: ProgressStateReplicate,
				Match: 10,
				Next:  11,
			},
			rejected: 10,
			want:     false,
		},
		{
			name: "rejected greater than match should return true and update Next",
			pr: &Progress{
				State: ProgressStateReplicate,
				Match: 10,
				Next:  11,
			},
			rejected: 15,
			want:     true,
			wantNext: 11,
		},
	}