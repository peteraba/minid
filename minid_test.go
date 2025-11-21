package minid

import (
	"sort"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNumToStringAndStringToNum(t *testing.T) {
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
			diff:    uint64(time.Date(8099, 12, 18, 15, 23, 17, 280_000_000, time.UTC).UnixMilli() - startAtUnixMilli),
			maxDiff: maxUnixMilliDiff,
			length:  8,
			want:    "zzzzzzzz",
		},
		{
			name:    "max unix micro diff",
			diff:    uint64(time.Date(2395, 7, 29, 21, 54, 52, 834_140_000, time.UTC).UnixMicro() - startAtUnixMicro),
			maxDiff: maxUnixMicroDiff,
			length:  9,
			want:    "zzzzzzzzz",
		},
		{
			name:    "max unix nano diff",
			diff:    uint64(time.Date(2317, 4, 12, 23, 47, 16, 854_775_808, time.UTC).UnixNano() - startAtUnixNano),
			maxDiff: maxUnixNanoDiff,
			length:  11,
			want:    "Dvik2W96uj9",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := numToString(tt.diff, tt.maxDiff, tt.length)
			require.Len(t, got, tt.length)

			back := stringToNum(got)

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
		seqs = append(seqs, numToString(atSeconds, maxUnixDiff, 6))
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

func TestBytes(t *testing.T) {
	tests := []struct {
		name    string
		minid   Minid
		wantLen int
	}{
		{
			name:    "zero",
			minid:   Minid("111111"),
			wantLen: 8,
		},
		{
			name:    "one",
			minid:   Minid("111112"),
			wantLen: 8,
		},
		{
			name:    "random",
			minid:   Minid("aB3x"),
			wantLen: 8,
		},
		{
			name:    "unix",
			minid:   Minid("132f3bXSZ"),
			wantLen: 8,
		},
		{
			name:    "max unix",
			minid:   Minid("zzzzzz"),
			wantLen: 8,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bytes := tt.minid.Bytes()
			require.Len(t, bytes, tt.wantLen, "Bytes() should return exactly 8 bytes")
		})
	}
}

func TestFromBytes(t *testing.T) {
	tests := []struct {
		name    string
		minid   Minid
		wantErr bool
	}{
		{
			name:    "zero",
			minid:   Minid("111111"),
			wantErr: false,
		},
		{
			name:    "one",
			minid:   Minid("111112"),
			wantErr: false,
		},
		{
			name:    "random",
			minid:   Minid("aB3x"),
			wantErr: false,
		},
		{
			name:    "unix",
			minid:   Minid("132f3bXSZ"),
			wantErr: false,
		},
		{
			name:    "max unix",
			minid:   Minid("zzzzzz"),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bytes := tt.minid.Bytes()
			decoded, err := FromBytes(bytes)

			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			// Verify the numeric value is preserved (even if string representation differs)
			assert.Equal(t, tt.minid.Uint64(), decoded.Uint64(),
				"Round-trip conversion should preserve numeric value")
		})
	}

	// Test error cases
	t.Run("invalid length - too short", func(t *testing.T) {
		_, err := FromBytes([]byte{1, 2, 3})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "must be exactly 8 bytes")
	})

	t.Run("invalid length - too long", func(t *testing.T) {
		_, err := FromBytes(make([]byte, 16))
		require.Error(t, err)
		assert.Contains(t, err.Error(), "must be exactly 8 bytes")
	})
}

func TestBytesRoundTrip(t *testing.T) {
	// Test various Minid types to ensure round-trip works
	testCases := []struct {
		name  string
		minid Minid
	}{
		{"random", Minid("aB3x")},
		{"unix", RandomUnix(1, 3)[0]},
		{"unixMilli", RandomUnixMilli(1, 3)[0]},
		{"unixMicro", RandomUnixMicro(1, 3)[0]},
		{"unixNano", RandomNano(1, 3)[0]},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			originalNum := tc.minid.Uint64()
			bytes := tc.minid.Bytes()
			decoded, err := FromBytes(bytes)

			require.NoError(t, err)
			assert.Equal(t, originalNum, decoded.Uint64(),
				"Round-trip conversion should preserve numeric value")
		})
	}
}
