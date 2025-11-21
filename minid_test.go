package minid

import (
	"fmt"
	"math/rand/v2"
	"sort"
	"strings"
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

func TestStringToDiffOverflow(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		shouldPanic bool
	}{
		{
			name:        "valid with many z's",
			input:       "zzzzzzzzzz",
			shouldPanic: false,
		},
		{
			name:        "overflow with many z's",
			input:       "zzzzzzzzzzz",
			shouldPanic: true,
		},
		{
			name:        "valid with many a's",
			input:       "aaaaaaaaaa",
			shouldPanic: false,
		},
		{
			name:        "overflow with too many a's",
			input:       "aaaaaaaaaaa",
			shouldPanic: true,
		},
		{
			name:        "valid with many 1's",
			input:       "11111111111",
			shouldPanic: false,
		},
		{
			name:        "overflow with too many 1's",
			input:       "111111111111",
			shouldPanic: true,
		},
		{
			name:        "valid max value",
			input:       "zzzzzz",
			shouldPanic: false,
		},
		{
			name:        "valid normal value",
			input:       "111111",
			shouldPanic: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.shouldPanic {
				require.Panics(t, func() {
					stringToDiff(tt.input)
				}, "stringToDiff should panic for input %q", tt.input)
			} else {
				require.NotPanics(t, func() {
					stringToDiff(tt.input)
				}, "stringToDiff should not panic for input %q", tt.input)
			}
		})
	}
}

func TestBytesAndFromBytes(t *testing.T) {
	tests := []string{
		"",
		"1",
		"z",
		"11",
		"1z",
		"zz",
		"111",
		"zzz",
		"1111",
		"zzzz",
		"11111",
		"123",
	}

	// Add some random valid strings
	letters := []rune("123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz")
	for i := range 100 {
		l := i % 20
		sb := strings.Builder{}
		for j := 0; j < l; j++ {
			sb.WriteRune(letters[rand.IntN(len(letters))])
		}
		tests = append(tests, sb.String())
	}

	for _, s := range tests {
		t.Run(fmt.Sprintf("len-%d-%s", len(s), s), func(t *testing.T) {
			m := Minid(s)
			b := m.Bytes()

			// Verify determinism
			b2 := m.Bytes()
			assert.Equal(t, b, b2)

			// Verify round trip
			m2, err := FromBytes(b)
			require.NoError(t, err)
			assert.Equal(t, m, m2)

			// Verify efficiency (roughly)
			// Length should be ceil(len(s) * 6 / 8)
			expectedLen := (len(s)*6 + 7) / 8
			assert.Equal(t, expectedLen, len(b), "Byte length mismatch for %s", s)
		})
	}
}

func BenchmarkRandom4(b *testing.B) {
	for b.Loop() {
		Random(1, 4)
	}
}

func BenchmarkRandom6(b *testing.B) {
	for b.Loop() {
		Random(1, 6)
	}
}

func BenchmarkRandomUnix(b *testing.B) {
	for b.Loop() {
		RandomUnix(1, 6)
	}
}

func BenchmarkRandomUnixMilli(b *testing.B) {
	for b.Loop() {
		RandomUnixMilli(1, 6)
	}
}

func BenchmarkRandomUnixMicro(b *testing.B) {
	for b.Loop() {
		RandomUnixMicro(1, 6)
	}
}

func BenchmarkRandomNano(b *testing.B) {
	for b.Loop() {
		RandomNano(1, 6)
	}
}
