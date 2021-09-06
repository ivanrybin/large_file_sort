package main

import (
	"github.com/ivanrybin/large_file_sort/pkg/gen"
	"github.com/spf13/pflag"
	"log"
	"os"
)

func main() {
	var stringsCnt int
	pflag.IntVarP(&stringsCnt, "count", "c", 0, "strings count")

	var stringsMaxLen int
	pflag.IntVarP(&stringsMaxLen, "max-len", "l", 0, "strings max length")

	var output string
	pflag.StringVarP(&output, "output", "o", "", "output path")

	var repeatedReversedAlphabetic bool
	pflag.BoolVarP(&repeatedReversedAlphabetic, "alpha", "a", false, "alphabetic reversed")

	pflag.Parse()

	if output == "" {
		log.Fatal("no output file was provided")
	}
	if stringsCnt <= 0 {
		log.Fatal("strings count must be > 0")
	}

	outFile, err := os.Create(output)
	if err != nil {
		log.Fatal(err)
	}
	defer func() { _ = outFile.Close() }()

	if repeatedReversedAlphabetic {

		if err = gen.RepeatedReversedAlphabeticStrings(stringsCnt, outFile); err != nil {
			log.Fatal(err)
		}

	} else {

		if err = gen.RndStrings(stringsCnt, stringsMaxLen, outFile); err != nil {
			log.Fatal(err)
		}

	}
}
