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
			name: "replicate state with stale rejection (rejected <= match)",
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
			name: "replicate state with valid rejection (rejected > match)",
			progress: &Progress{
				State: ProgressStateReplicate,
				Match: 10,
				Next:  15,
			},
			rejected: 12,
			last:     20,
			want:     true,
			wantNext: 11,
			wantMatch: 10,
		},
		{
			name: "replicate state with rejection equal to match",
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
			name: "replicate state with rejection at boundary (match + 1)",
			progress: &Progress{
				State: ProgressStateReplicate,
				Match: 10,
				Next:  15,
			},
			rejected: 11,
			last:     20,
			want:     true,
			wantNext: 11,
			wantMatch: 10,
		},

		// Non-Replicate state cases
		{
			name: "non-replicate state with stale rejection (rejected != next-1)",
			progress: &Progress{
				State: ProgressStateProbe,
				Match: 5,
				Next:  10,
			},
			rejected: 8,
			last:     20,
			want:     false,
			wantNext: 10,
			wantMatch: 5,
		},
		{
			name: "non-replicate state with valid rejection (rejected == next-1)",
			progress: &Progress{
				State: ProgressStateProbe,
				Match: 5,
				Next:  10,
			},
			rejected: 9,
			last:     20,
			want:     true,
			wantNext: 9,
			wantMatch: 5,
		},
		{
			name: "non-replicate state with rejection and last smaller",
			progress: &Progress{
				State: ProgressStateProbe,
				Match: 5,
				Next:  10,
			},
			rejected: 9,
			last:     7,
			want:     true,
			wantNext: 8,
			wantMatch: 5,
		},
		{
			name: "non-replicate state with rejection and last larger",
			progress: &Progress{
				State: ProgressStateProbe,
				Match: 5,
				Next:  10,
			},
			rejected: 9,
			last:     15,
			want:     true,
			wantNext: 9,
			wantMatch: 5,
		},
		{
			name: "non-replicate state with rejection causing next to go below 1",
			progress: &Progress{
				State: ProgressStateProbe,
				Match: 0,
				Next:  1,
			},
			rejected: 0,
			last:     0,
			want:     true,
			wantNext: 1,
			wantMatch: 0,
		},
		{
			name: "non-replicate state with rejection at boundary (next-1 == rejected)",
			progress: &Progress{
				State: ProgressStateProbe,
				Match: 0,
				Next:  1,
			},
			rejected: 0,
			last:     5,
			want:     true,
			wantNext: 1,
			wantMatch: 0,
		},
		{
			name: "non-replicate state with rejection and last causing next to be 1",
			progress: &Progress{
				State: ProgressStateProbe,
				Match: 0,
				Next:  5,
			},
			rejected: 4,
			last:     0,
			want:     true,
			wantNext: 1,
			wantMatch: 0,
		},
		{
			name: "non-replicate state with rejection equal to last+1",
			progress: &Progress{
				State: ProgressStateProbe,
				Match: 5,
				Next:  10,
			},
			rejected: 9,
			last:     8,
			want:     true,
			wantNext: 9,
			wantMatch: 5,
		},
		{
			name: "non-replicate state with rejection less than last+1",
			progress: &Progress{
				State: ProgressStateProbe,
				Match: 5,
				Next:  10,
			},
			rejected: 9,
			last:     10,
			want:     true,
			wantNext: 9,
			wantMatch: 5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a copy of progress to avoid test pollution
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
				t.Errorf("maybeDecrTo() progress.Next = %v, want %v", progress.Next, tt.wantNext)
			}
			if progress.Match != tt.wantMatch {
				t.Errorf("maybeDecrTo() progress.Match = %v, want %v", progress.Match, tt.wantMatch)
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