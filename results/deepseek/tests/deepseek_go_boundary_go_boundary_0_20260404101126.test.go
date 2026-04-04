package raft

import (
	"testing"
)

func TestProgress_maybeDecrTo(t *testing.T) {
	tests := []struct {
		name     string
		progress *Progress
		rejected uint64
		last     uint64
		want     bool
		wantNext uint64
		wantMatch uint64
	}{
		// ProgressStateReplicate cases
		{
			name: "Replicate state with stale rejection (rejected <= match)",
			progress: &Progress{
				State: ProgressStateReplicate,
				Match: 10,
				Next:  15,
			},
			rejected: 8,
			last:     20,
			want:     false,
			wantNext: 15,
			wantMatch: 10,
		},
		{
			name: "Replicate state with valid rejection (rejected > match)",
			progress: &Progress{
				State: ProgressStateReplicate,
				Match: 10,
				Next:  15,
			},
			rejected: 12,
			last:     20,
			want:     true,
			wantNext: 11, // match + 1
			wantMatch: 10,
		},
		{
			name: "Replicate state with rejected equal to match",
			progress: &Progress{
				State: ProgressStateReplicate,
				Match: 10,
				Next:  15,
			},
			rejected: 10,
			last:     20,
			want:     false,
			wantNext: 15,
			wantMatch: 10,
		},
		{
			name: "Replicate state with match at zero",
			progress: &Progress{
				State: ProgressStateReplicate,
				Match: 0,
				Next:  5,
			},
			rejected: 3,
			last:     10,
			want:     true,
			wantNext: 1, // 0 + 1
			wantMatch: 0,
		},

		// Non-Replicate state cases (ProgressStateProbe or ProgressStateSnapshot)
		{
			name: "Non-replicate state with valid rejection (next-1 == rejected)",
			progress: &Progress{
				State: ProgressStateProbe,
				Match: 5,
				Next:  10,
			},
			rejected: 9,
			last:     15,
			want:     true,
			wantNext: 9, // min(9, 15+1) = min(9, 16) = 9
			wantMatch: 5,
		},
		{
			name: "Non-replicate state with stale rejection (next-1 != rejected)",
			progress: &Progress{
				State: ProgressStateProbe,
				Match: 5,
				Next:  10,
			},
			rejected: 8,
			last:     15,
			want:     false,
			wantNext: 10,
			wantMatch: 5,
		},
		{
			name: "Non-replicate state with rejected smaller than last+1",
			progress: &Progress{
				State: ProgressStateProbe,
				Match: 5,
				Next:  10,
			},
			rejected: 9,
			last:     7,
			want:     true,
			wantNext: 8, // min(9, 7+1) = min(9, 8) = 8
			wantMatch: 5,
		},
		{
			name: "Non-replicate state with next adjustment to 1 when result < 1",
			progress: &Progress{
				State: ProgressStateProbe,
				Match: 5,
				Next:  2,
			},
			rejected: 1,
			last:     0,
			want:     true,
			wantNext: 1, // min(1, 0+1) = min(1, 1) = 1
			wantMatch: 5,
		},
		{
			name: "Non-replicate state with next at 1 and rejected at 0",
			progress: &Progress{
				State: ProgressStateProbe,
				Match: 0,
				Next:  1,
			},
			rejected: 0,
			last:     0,
			want:     true,
			wantNext: 1, // min(0, 0+1) = min(0, 1) = 0, then adjusted to 1
			wantMatch: 0,
		},
		{
			name: "Non-replicate state with large values",
			progress: &Progress{
				State: ProgressStateSnapshot,
				Match: 1000,
				Next:  2000,
			},
			rejected: 1999,
			last:     3000,
			want:     true,
			wantNext: 1999, // min(1999, 3000+1) = min(1999, 3001) = 1999
			wantMatch: 1000,
		},
		{
			name: "Non-replicate state with next at max uint64",
			progress: &Progress{
				State: ProgressStateProbe,
				Match: 100,
				Next:  ^uint64(0),
			},
			rejected: ^uint64(0) - 1,
			last:     ^uint64(0) - 2,
			want:     true,
			wantNext: ^uint64(0) - 1, // min(max-1, max-2+1) = min(max-1, max-1) = max-1
			wantMatch: 100,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a copy to avoid test pollution
			progress := &Progress{
				State: tt.progress.State,
				Match: tt.progress.Match,
				Next:  tt.progress.Next,
			}
			
			got := progress.maybeDecrTo(tt.rejected, tt.last)
			
			if got != tt.want {
				t.Errorf("maybeDecrTo() = %v, want %v", got, tt.want)
			}
			if progress.Next != tt.wantNext {
				t.Errorf("Next = %v, want %v", progress.Next, tt.wantNext)
			}
			if progress.Match != tt.wantMatch {
				t.Errorf("Match = %v, want %v", progress.Match, tt.wantMatch)
			}
		})
	}
}

// Helper function for min operation
func min(a, b uint64) uint64 {
	if a < b {
		return a
	}
	return b
}