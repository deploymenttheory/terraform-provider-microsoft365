package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/deploymenttheory/terraform-provider-microsoft365/scripts/go/nightly_tests/testrunner/runner"
)

func main() {
	cfg := parseFlags()

	if err := cfg.Validate(); err != nil {
		fmt.Fprintf(os.Stderr, "Configuration error: %v\n", err)
		flag.Usage()
		os.Exit(1)
	}

	r := runner.New(cfg, runner.NewEnvCredentialProvider())

	if err := r.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Test execution failed: %v\n", err)
		os.Exit(1)
	}
}

func parseFlags() *runner.Config {
	cfg := &runner.Config{}

	flag.StringVar(&cfg.TestType, "type", "", "Type of test: provider-core, resources, datasources")
	flag.StringVar(&cfg.Service, "service", "", "Service name (required for resources/datasources)")
	flag.StringVar(&cfg.CoverageOutput, "coverage", "", "Coverage output file path")
	flag.BoolVar(&cfg.Verbose, "verbose", false, "Enable verbose output")
	flag.Parse()

	return cfg
}
