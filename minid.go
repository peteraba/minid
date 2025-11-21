package minid

import (
	"errors"
	"fmt"
	"math"
	"math/rand/v2"
	"strings"
	"time"
)

const maxRetries = 100

var duplicateDetector = make(map[string]struct{})

var (
	start            = time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	startAtUnix      = start.Unix()
	startAtUnixMilli = start.UnixMilli()
	startAtUnixMicro = start.UnixMicro()
	startAtUnixNano  = start.UnixNano()
	letters          = []rune("123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz")
	base             = uint64(len(letters))
	maxUnixDiff      = uint64(math.Pow(float64(len(letters)), 6)) - 1
	maxUnixMilliDiff = uint64(math.Pow(float64(len(letters)), 8))
	maxUnixMicroDiff = uint64(math.Pow(float64(len(letters)), 9))
	maxUnixNanoDiff  = uint64(math.MaxUint64)
)

var (
	errNegativeRelativeTimestamp              = errors.New("relative timestamp is negative")
	errDiffTooLarge                           = errors.New("diff is too large")
	errFailedToGenerateSequenceWithoutNumbers = errors.New("failed to generate a sequence without numbers")
)

func numToString(diff, maxDiff uint64, length int) string {
	if diff > maxDiff {
		panic(errDiffTooLarge)
	}

	var result []rune

	// Generate base61 digits (little-endian)
	for diff > 0 {
		result = append(result, letters[diff%base])
		diff /= base
	}

	// Pad with '1' (letters[0]) to fixed length of 6
	for len(result) < length {
		result = append(result, letters[0])
	}

	// Reverse the result to get big-endian
	for i, j := 0, len(result)-1; i < j; i, j = i+1, j-1 {
		result[i], result[j] = result[j], result[i]
	}

	return string(result)
}

func randSeq(n int, retryCount int) string {
	if retryCount > maxRetries {
		panic(errFailedToGenerateSequenceWithoutNumbers)
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

func Random(count, randLength int) []string {
	var seqs = make([]string, 0, count)
	for range count {
		seqs = append(seqs, randSeq(randLength, 0))
	}

	return seqs
}

func RandomUnix(count, randLength int) []string {
	var seqs = make([]string, 0, count)
	for range count {
		t := time.Now()
		if t.Before(start) {
			panic(errNegativeRelativeTimestamp)
		}

		ts := time.Now().Unix() - startAtUnix
		seqs = append(seqs, fmt.Sprintf("%s%s", numToString(uint64(ts), maxUnixDiff, 6), randSeq(randLength, 0)))
	}

	return seqs
}

func RandomUnixMilli(count, randLength int) []string {
	var seqs = make([]string, 0, count)
	for range count {
		t := time.Now()
		if t.Before(start) {
			panic(errNegativeRelativeTimestamp)
		}

		ts := t.UnixMilli() - startAtUnixMilli
		seqs = append(seqs, fmt.Sprintf("%s%s", numToString(uint64(ts), maxUnixMilliDiff, 8), randSeq(randLength, 0)))
	}

	return seqs
}

func RandomUnixMicro(count, randLength int) []string {
	var seqs = make([]string, 0, count)
	for range count {
		t := time.Now()
		if t.Before(start) {
			panic(errNegativeRelativeTimestamp)
		}

		ts := time.Now().UnixMicro() - startAtUnixMicro
		seqs = append(seqs, fmt.Sprintf("%s%s", numToString(uint64(ts), maxUnixMicroDiff, 10), randSeq(randLength, 0)))
	}

	return seqs
}

func RandomNano(count, randLength int) []string {
	var seqs = make([]string, 0, count)
	for range count {
		t := time.Now()
		if t.Before(start) {
			panic(errNegativeRelativeTimestamp)
		}

		ts := t.UnixNano() - startAtUnixNano
		seqs = append(seqs, fmt.Sprintf("%s%s", numToString(uint64(ts), maxUnixNanoDiff, 11), randSeq(randLength, 0)))
	}

	return seqs
}
