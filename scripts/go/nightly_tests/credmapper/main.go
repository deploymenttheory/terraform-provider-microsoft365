package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/deploymenttheory/terraform-provider-microsoft365/scripts/go/nightly_tests/credmapper/mapper"
)

func main() {
	var service string
	flag.StringVar(&service, "service", "", "Service name to map credentials for")
	flag.Parse()

	if service == "" {
		fmt.Fprintln(os.Stderr, "Error: -service flag is required")
		flag.Usage()
		os.Exit(1)
	}

	m := mapper.New(mapper.NewEnvExporter())

	if err := m.Map(service); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
