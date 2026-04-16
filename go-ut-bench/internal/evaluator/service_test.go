package evaluator

import "testing"

func TestParsePytestCounts(t *testing.T) {
	out := "1 failed, 14 passed in 0.09s"
	passed, total := parsePytestCounts(out)
	if passed == nil || total == nil {
		t.Fatalf("expected counts parsed")
	}
	if *passed != 14 || *total != 15 {
		t.Fatalf("unexpected parsed counts: passed=%d total=%d", *passed, *total)
	}
}

func TestParsePytestCountsPassedOnly(t *testing.T) {
	out := "13 passed in 0.05s"
	passed, total := parsePytestCounts(out)
	if passed == nil || total == nil {
		t.Fatalf("expected counts parsed")
	}
	if *passed != 13 || *total != 13 {
		t.Fatalf("unexpected parsed counts: passed=%d total=%d", *passed, *total)
	}
}
