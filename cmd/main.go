package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/peteraba/minid"
)

type RandomType string

const (
	RandomTypeRandom    RandomType = "r"
	RandomTypeUnix      RandomType = "s"
	RandomTypeUnixMilli RandomType = "ms"
	RandomTypeUnixMicro RandomType = "us"
	RandomTypeUnixNano  RandomType = "ns"
)

func main() {
	var err error

	randLengthFlag := flag.Int("randLength", 0, "Length of the random suffix")
	randLengthShort := flag.Int("rl", 0, "Length of the random suffix (short)")
	flag.Parse()

	randomType := RandomTypeRandom
	count := 1
	randLength := 4

	// Get remaining positional arguments after flags
	args := flag.Args()
	arg1 := ""
	if len(args) > 0 {
		arg1 = args[0]
	}

	switch arg1 {
	case "s", "ms", "us", "ns":
		randomType = RandomType(arg1)
		randLength = 3
	case "r", "":
	default:
		count, err = strconv.Atoi(arg1)
		if err != nil {
			os.Exit(1)
		}
	}

	if len(args) > 1 {
		count, err = strconv.Atoi(args[1])
		if err != nil {
			os.Exit(1)
		}
	}

	// Use flag value if set, otherwise use default based on arguments
	if *randLengthFlag > 0 {
		randLength = *randLengthFlag
	} else if *randLengthShort > 0 {
		randLength = *randLengthShort
	}

	switch randomType {
	case RandomTypeRandom:
		seqs := minid.Random(count, randLength)
		printSeqs(seqs)
	case RandomTypeUnix:
		seqs := minid.RandomUnix(count, randLength)
		printSeqs(seqs)
	case RandomTypeUnixMilli:
		seqs := minid.RandomUnixMilli(count, randLength)
		printSeqs(seqs)
	case RandomTypeUnixMicro:
		seqs := minid.RandomUnixMicro(count, randLength)
		printSeqs(seqs)
	case RandomTypeUnixNano:
		seqs := minid.RandomNano(count, randLength)
		printSeqs(seqs)

	default:
		fmt.Printf("Invalid random type: %s\n", randomType)
		os.Exit(1)
	}
}

func printSeqs(seqs []string) {
	for _, seq := range seqs {
		fmt.Println(seq)
	}
}
