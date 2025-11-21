package minid

import (
	"sort"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDiffToStringAndStringToDiff(t *testing.T) {
	tests := []struct {
		name    string
		diff    uint64
		maxDiff uint64
		length  int
		want    string
	}{
		{
			name:    "zero",
			diff:    0,
			maxDiff: maxUnixDiff,
			length:  6,
			want:    "111111",
		},
		{
			name:    "one",
			diff:    1,
			maxDiff: maxUnixDiff,
			length:  6,
			want:    "111112",
		},
		{
			name:    "base",
			diff:    61,
			maxDiff: maxUnixDiff,
			length:  6,
			want:    "111121",
		},
		{
			name:    "large number",
			diff:    12345,
			maxDiff: maxUnixDiff,
			length:  6,
			want:    "1114KO",
		},
		{
			name:    "large number",
			diff:    12345,
			maxDiff: maxUnixDiff,
			length:  6,
			want:    "1114KO",
		},
		{
			name:    "max unix diff",
			diff:    uint64(time.Date(3657, 8, 13, 15, 6, 0, 0, time.UTC).Unix() - time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC).Unix()),
			maxDiff: maxUnixDiff,
			length:  6,
			want:    "zzzzzz",
		},
		{
			name:    "zero-milli",
			diff:    0,
			maxDiff: maxUnixMilliDiff,
			length:  8,
			want:    "11111111",
		},
		{
			name:    "max unix milli diff",
			diff:    uint64(time.Date(8099, 12, 18, 15, 23, 17, 280_000_000, time.UTC).UnixMilli() - epochMilli),
			maxDiff: maxUnixMilliDiff,
			length:  8,
			want:    "zzzzzzzz",
		},
		{
			name:    "max unix micro diff",
			diff:    uint64(time.Date(2395, 7, 29, 21, 54, 52, 834_140_000, time.UTC).UnixMicro() - epochMicro),
			maxDiff: maxUnixMicroDiff,
			length:  9,
			want:    "zzzzzzzzz",
		},
		{
			name:    "max unix nano diff",
			diff:    uint64(time.Date(2317, 4, 12, 23, 47, 16, 854_775_808, time.UTC).UnixNano() - epochNano),
			maxDiff: maxUnixNanoDiff,
			length:  11,
			want:    "Dvik2W96uj9",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := diffToString(tt.diff, tt.maxDiff, tt.length)
			require.Len(t, got, tt.length)

			back := stringToDiff(got)

			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.diff, back)
		})
	}
}

func TestSortability(t *testing.T) {
	var seqs []string

	atSeconds := uint64(1000)
	for range 100 {
		atSeconds += 10
		seqs = append(seqs, diffToString(atSeconds, maxUnixDiff, 6))
	}

	// Check if they are sorted lexicographically
	if !sort.StringsAreSorted(seqs) {
		t.Errorf("Generated sequences are not sorted lexicographically")
		for i := 0; i < len(seqs)-1; i++ {
			if seqs[i] > seqs[i+1] {
				t.Logf("Unsorted pair at index %d: %s > %s", i, seqs[i], seqs[i+1])
			}
		}
	}
}

func TestRandom(t *testing.T) {
	seqs := Random(100, 3)

	seen := make(map[string]struct{}, len(seqs))
	for _, s := range seqs {
		if _, exists := seen[s.String()]; exists {
			t.Errorf("Duplicate value found in seqs: %s", s)
		}
		seen[s.String()] = struct{}{}
		if len(s) != 3 {
			t.Errorf("Sequence length is not 3: %s (got %d)", s, len(s))
		}
	}
}

func TestRandomUnix(t *testing.T) {
	seqs := RandomUnix(100, 3)

	seen := make(map[string]struct{}, len(seqs))
	for _, s := range seqs {
		if _, exists := seen[s.String()]; exists {
			t.Errorf("Duplicate value found in seqs: %s", s.String())
		}
		seen[s.String()] = struct{}{}
		if len(s.String()) != 9 {
			t.Errorf("Sequence length is not 9: %s (got %d)", s.String(), len(s.String()))
		}
	}
}
