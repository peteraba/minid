package minid

import (
	"errors"
	"fmt"
	"math"
	"math/rand/v2"
	"slices"
	"strings"
	"time"
)

const maxRetries = 100

var duplicateDetector = make(map[string]struct{})

var (
	epoch = time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)

	epochUnix  = epoch.Unix()
	epochMilli = epoch.UnixMilli()
	epochMicro = epoch.UnixMicro()
	epochNano  = epoch.UnixNano()

	letters    = []rune("123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz")
	base       = uint64(len(letters))
	lettersMap = [128]int8{}

	maxUnixDiff      = uint64(math.Pow(float64(len(letters)), 6)) - 1
	maxUnixMilliDiff = uint64(math.Pow(float64(len(letters)), 8))
	maxUnixMicroDiff = uint64(math.Pow(float64(len(letters)), 9))
	maxUnixNanoDiff  = uint64(math.MaxUint64)
)

func init() {
	for i := range lettersMap {
		lettersMap[i] = -1
	}
	for i, r := range letters {
		lettersMap[r] = int8(i)
	}
}

var (
	errTimeTravelDetected = errors.New("the relative timestamp is negative")
	errDiffTooLarge       = errors.New("time diff out of range")
	errTooManyRetries     = errors.New("too many retries")
)

type RandomType string

const (
	RandomTypeRandom    RandomType = "r"
	RandomTypeUnix      RandomType = "s"
	RandomTypeUnixMilli RandomType = "m"
	RandomTypeUnixMicro RandomType = "u"
	RandomTypeUnixNano  RandomType = "n"
)

// Minid is a single ID.
type Minid string

// String returns the ID as a string.
func (m Minid) String() string {
	return string(m)
}

// Bytes returns the ID as a byte slice.
func (m Minid) Bytes() []byte {
	s := string(m)
	n := len(s)
	if n == 0 {
		return []byte{}
	}

	// Calculate byte length
	// Pack 6 bits per char.
	// n chars -> ceil(n * 6 / 8) bytes
	numBytes := (n*6 + 7) / 8
	b := make([]byte, numBytes)

	var bitBuf uint64
	var bitCount int
	byteIdx := 0

	for _, r := range s {
		var idx uint64
		if int(r) < len(lettersMap) {
			val := lettersMap[r]
			if val >= 0 {
				idx = uint64(val)
			} else {
				// Invalid character, treat as 0 ('1')
				idx = 0
			}
		} else {
			idx = 0
		}

		// Add 6 bits to buffer
		bitBuf = (bitBuf << 6) | idx
		bitCount += 6

		for bitCount >= 8 {
			// Extract top 8 bits
			bitCount -= 8
			b[byteIdx] = byte(bitBuf >> bitCount)
			byteIdx++
			// Mask to keep only valid bits in buf
			bitBuf &= (1 << bitCount) - 1
		}
	}

	// Handle remaining bits
	if bitCount > 0 {
		// We have bitCount bits at the bottom of bitBuf.
		// Shift them to the top of the byte.
		b[byteIdx] = byte(bitBuf << (8 - bitCount))

		// Handle ambiguity padding.
		// Ambiguity arises when len(b) % 3 == 0.
		// If len(b) % 3 == 0, it could be 3-char case (with slack) or 4-char case (full).
		// If 3-char case (slack): bitCount should be 2 (18 bits used = 2 full + 2 bits).
		// Remaining bits = 6.
		if (byteIdx+1)%3 == 0 {
			// We must pad the remaining 6 bits with 1s.
			remainingBits := 8 - bitCount
			mask := byte((1 << remainingBits) - 1)
			b[byteIdx] |= mask
		}
	}

	return b
}

// FromBytes reconstructs a Minid from its byte representation.
func FromBytes(b []byte) (Minid, error) {
	nBytes := len(b)
	if nBytes == 0 {
		return "", nil
	}

	nChars := (nBytes * 8) / 6

	if nBytes%3 == 0 {
		// Check for padding in the last byte
		// In 3-char case, last byte has 6 bits of padding (1s).
		// 1s are in the LSB positions.
		if b[nBytes-1]&0x3F == 0x3F {
			nChars--
		}
	}

	sb := strings.Builder{}
	sb.Grow(nChars)

	var bitBuf uint64
	var bitCount int

	// We process byte by byte, but stop when we have extracted nChars
	charsExtracted := 0

	for i := 0; i < nBytes && charsExtracted < nChars; i++ {
		bitBuf = (bitBuf << 8) | uint64(b[i])
		bitCount += 8

		for bitCount >= 6 && charsExtracted < nChars {
			// Extract top 6 bits
			val := (bitBuf >> (bitCount - 6)) & 0x3F
			bitCount -= 6

			if val >= uint64(len(letters)) {
				return "", fmt.Errorf("invalid character index: %d", val)
			}
			sb.WriteByte(byte(letters[val]))
			charsExtracted++
		}
		// Keep remaining bits in bitBuf
		bitBuf &= (1 << bitCount) - 1
	}

	return Minid(sb.String()), nil
}

// Minids is a slice of Minid.
type Minids []Minid

// StringSlice returns a slice of strings representing the IDs.
func (m Minids) StringSlice() []string {
	slice := make([]string, 0, len(m))
	for _, seq := range m {
		slice = append(slice, seq.String())
	}

	return slice
}

// Print prints the IDs to the console.
func (m Minids) Print() {
	for _, seq := range m {
		fmt.Println(string(seq))
	}
}

// Sort sorts the IDs lexicographically.
func (m Minids) Sort() {
	slices.Sort(m)
}

// stringToDiff converts a minid-string to a diff.
func stringToDiff(s string) uint64 {
	// Calculate maximum safe length: largest n where base^n <= math.MaxUint64 + 1
	// We add 1 to account for the fact that we need to round up
	maxSafeLength := int(math.Ceil(math.Log(float64(math.MaxUint64)+1) / math.Log(float64(base))))
	if len(s) > maxSafeLength {
		panic(fmt.Errorf("string %q (length %d) exceeds maximum safe length %d for uint64 conversion", s, len(s), maxSafeLength))
	}

	num := uint64(0)
	for _, c := range s {
		var index uint64
		if int(c) < len(lettersMap) && lettersMap[c] != -1 {
			index = uint64(lettersMap[c])
		}
		// Check for overflow before multiplication
		// num*base + index would overflow if num > (math.MaxUint64 - index) / base
		if num > (math.MaxUint64-index)/base {
			panic(fmt.Errorf("string %q would overflow uint64 when converted to diff", s))
		}
		num = num*base + index
	}
	return num
}

// diffToString converts a diff to a minid-string.
func diffToString(diff, maxDiff uint64, length int) string {
	if diff > maxDiff {
		panic(errDiffTooLarge)
	}

	var result []rune

	// Generate minid-letters (little-endian)
	for diff > 0 {
		result = append(result, letters[diff%base])
		diff /= base
	}

	// Pad with '1' (letters[0]) to fixed length of 6
	for len(result) < length {
		result = append(result, letters[0])
	}

	// Reverse the result to get minid-letters (big-endian)
	for i, j := 0, len(result)-1; i < j; i, j = i+1, j-1 {
		result[i], result[j] = result[j], result[i]
	}

	return string(result)
}

// randSeq generates a random sequence of letters.
func randSeq(n int, retryCount int) string {
	if retryCount > maxRetries {
		panic(errTooManyRetries)
	}

	reTry := true
	strBuilder := strings.Builder{}
	for range n {
		letter := letters[rand.IntN(len(letters))]
		if reTry && letter > '9' {
			reTry = false
		}
		strBuilder.WriteRune(letter)
	}

	str := strBuilder.String()
	if _, ok := duplicateDetector[str]; ok {
		return randSeq(n, retryCount+1)
	}

	duplicateDetector[str] = struct{}{}

	if reTry {
		return randSeq(n, retryCount+1)
	}

	return str
}

// Random generates random IDs.
func Random(count, randLength int) Minids {
	var seqs = make([]Minid, 0, count)
	for range count {
		seqs = append(seqs, Minid(randSeq(randLength, 0)))
	}

	return Minids(seqs)
}

// RandomUnix generates random IDs with Unix second precision.
func RandomUnix(count, randLength int) Minids {
	var seqs = make(Minids, 0, count)
	for range count {
		t := time.Now()
		if t.Before(epoch) {
			panic(errTimeTravelDetected)
		}

		ts := time.Now().Unix() - epochUnix
		seqs = append(seqs, Minid(fmt.Sprintf("%s%s", diffToString(uint64(ts), maxUnixDiff, 6), randSeq(randLength, 0))))
	}

	return seqs
}

// RandomUnixMilli generates random IDs with Unix millisecond precision.
func RandomUnixMilli(count, randLength int) Minids {
	var seqs = make(Minids, 0, count)
	for range count {
		t := time.Now()
		if t.Before(epoch) {
			panic(errTimeTravelDetected)
		}

		ts := t.UnixMilli() - epochMilli
		seqs = append(seqs, Minid(fmt.Sprintf("%s%s", diffToString(uint64(ts), maxUnixMilliDiff, 8), randSeq(randLength, 0))))
	}

	return seqs
}

// RandomUnixMicro generates random IDs with Unix microsecond precision.
func RandomUnixMicro(count, randLength int) Minids {
	var seqs = make(Minids, 0, count)
	for range count {
		t := time.Now()
		if t.Before(epoch) {
			panic(errTimeTravelDetected)
		}

		ts := time.Now().UnixMicro() - epochMicro
		seqs = append(seqs, Minid(fmt.Sprintf("%s%s", diffToString(uint64(ts), maxUnixMicroDiff, 10), randSeq(randLength, 0))))
	}

	return seqs
}

// RandomNano generates random IDs with Unix nanosecond precision.
func RandomNano(count, randLength int) Minids {
	var seqs = make(Minids, 0, count)
	for range count {
		t := time.Now()
		if t.Before(epoch) {
			panic(errTimeTravelDetected)
		}

		ts := t.UnixNano() - epochNano
		seqs = append(seqs, Minid(fmt.Sprintf("%s%s", diffToString(uint64(ts), maxUnixNanoDiff, 11), randSeq(randLength, 0))))
	}

	return seqs
}
