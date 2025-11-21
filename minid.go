package minid

import (
	"errors"
	"fmt"
	"math"
	"math/rand/v2"
	"sort"
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
	return nil
}

// FromBytes reconstructs a Minid from its byte representation.
func FromBytes(b []byte) (Minid, error) {
	return "", nil
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
	sort.Slice(m, func(i, j int) bool {
		return m[i] < m[j]
	})
}

// stringToNum converts a minid-string to a number.
func stringToNum(s string) uint64 {
	num := uint64(0)
	for _, c := range s {
		var index uint64
		if int(c) < len(lettersMap) && lettersMap[c] != -1 {
			index = uint64(lettersMap[c])
		}
		num = num*base + index
	}
	return num
}

// numToString converts a number to a minid-string.
func numToString(diff, maxDiff uint64, length int) string {
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
		seqs = append(seqs, Minid(fmt.Sprintf("%s%s", numToString(uint64(ts), maxUnixDiff, 6), randSeq(randLength, 0))))
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
		seqs = append(seqs, Minid(fmt.Sprintf("%s%s", numToString(uint64(ts), maxUnixMilliDiff, 8), randSeq(randLength, 0))))
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
		seqs = append(seqs, Minid(fmt.Sprintf("%s%s", numToString(uint64(ts), maxUnixMicroDiff, 10), randSeq(randLength, 0))))
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
		seqs = append(seqs, Minid(fmt.Sprintf("%s%s", numToString(uint64(ts), maxUnixNanoDiff, 11), randSeq(randLength, 0))))
	}

	return seqs
}
