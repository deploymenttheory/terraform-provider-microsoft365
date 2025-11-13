package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/deploymenttheory/terraform-provider-microsoft365/scripts/go/nightly_tests/coveragemerger/merger"
)

func main() {
	var (
		inputDir  string
		outputFile string
	)

	flag.StringVar(&inputDir, "input", "", "Input directory containing coverage files")
	flag.StringVar(&outputFile, "output", "", "Output merged coverage file path")
	flag.Parse()

	if inputDir == "" || outputFile == "" {
		fmt.Fprintln(os.Stderr, "Error: both -input and -output flags are required")
		flag.Usage()
		os.Exit(1)
	}

	m := merger.New(inputDir, outputFile)

	if err := m.Merge(); err != nil {
		fmt.Fprintf(os.Stderr, "Error merging coverage files: %v\n", err)
		os.Exit(1)
	}
}
