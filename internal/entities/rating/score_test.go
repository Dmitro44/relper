package rating

import (
	"math"
	"testing"
)

const testZ = 1.96

func TestWilsonNoVotes(t *testing.T) {
	if got := Wilson(0, 0, testZ); got != 0 {
		t.Fatalf("expected 0, got %v", got)
	}
}

func TestWilsonMoreVotesRankHigher(t *testing.T) {
	few := Wilson(10, 0, testZ)
	many := Wilson(900, 100, testZ)

	if !(few < many) {
		t.Fatalf("expected Wilson(10,0)=%v < Wilson(900,100)=%v", few, many)
	}
}

func TestWilsonAllNegative(t *testing.T) {
	for _, tc := range [][2]int64{{-1, -1}, {-10, 0}, {0, -10}} {
		if got := Wilson(tc[0], tc[1], testZ); got != 0 {
			t.Fatalf("Wilson(%d,%d) = %v, want 0", tc[0], tc[1], got)
		}
	}
}

func TestWilsonSingleVote(t *testing.T) {
	got := Wilson(1, 0, testZ)
	if got < 0 || got > 1 {
		t.Fatalf("expected value in [0,1], got %v", got)
	}
}

func TestWilsonAlwaysInRange(t *testing.T) {
	for up := int64(0); up <= 20; up++ {
		for down := int64(0); down <= 20; down++ {
			got := Wilson(up, down, testZ)
			if math.IsNaN(got) || got < 0 || got > 1 {
				t.Fatalf("Wilson(%d,%d) = %v, expected [0,1]", up, down, got)
			}
		}
	}
}

func TestWilsonMonotonicInUp(t *testing.T) {
	for n := int64(1); n <= 50; n++ {
		prev := -1.0
		for up := int64(0); up <= n; up++ {
			got := Wilson(up, n-up, testZ)
			if got < prev-1e-12 {
				t.Fatalf("Wilson not monotonic: n=%d up=%d prev=%v got=%v", n, up, prev, got)
			}
			prev = got
		}
	}
}
