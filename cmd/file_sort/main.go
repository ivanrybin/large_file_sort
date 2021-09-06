package main

import (
	"github.com/ivanrybin/large_file_sort/pkg/sort"
	"github.com/spf13/pflag"
	"log"
)

func main() {
	var inputPath string
	pflag.StringVarP(&inputPath, "input", "i", "", "input file path")

	var outputPath string
	pflag.StringVarP(&outputPath, "output", "o", "", "output file path")
	pflag.Parse()

	if inputPath == "" {
		log.Fatal("no input path was provided")
	}
	if outputPath == "" {
		log.Fatal("no output path was provided")
	}

	if err := sort.Sort(inputPath, outputPath); err != nil {
		log.Fatal(err)
	}
}
