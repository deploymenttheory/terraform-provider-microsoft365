// Command helperscan reports package-local helpers against the helper rule in docs/development/guide.md.
//
// For each function in a package that is not a CRUD, schema or metadata method or constructResource, it
// reports the declaration's line span, how many places in the package reference it, which CRUD methods (Create, Read, Update, Delete,
// ImportState, ModifyPlan) reference it directly, and whether its name follows the verb + Graph entity + qualifier convention. Kiota
// interface methods and framework validator methods are counted but kept out of the helper table. The package's Graph-to-state mapper,
// MapRemoteStateToTerraform (MapRemoteStateToDataSource in a package under internal/services/datasources), is listed as a
// state-mapper: its file is checked, but the caller and line rules don't apply to it.
//
// Usage:
//
//	go run ./tools/helperscan [-format text|markdown|tsv] <package dir>...
//
// The report is informational: the command exits 0 whatever it finds, and only fails on bad usage.
package main

import (
	"flag"
	"fmt"
	"io"
	"os"
)

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

// run scans the package directories in args and writes the report, returning the process exit code.
func run(args []string, stdout, stderr io.Writer) int {
	flags := flag.NewFlagSet("helperscan", flag.ContinueOnError)
	flags.SetOutput(stderr)
	format := flags.String("format", "text", "output format: text, markdown or tsv")
	if err := flags.Parse(args); err != nil {
		return 2
	}
	if *format != "text" && *format != "markdown" && *format != "tsv" {
		fmt.Fprintf(stderr, "helperscan: unknown format %q\n", *format)
		return 2
	}
	if flags.NArg() == 0 {
		fmt.Fprintln(stderr, "usage: helperscan [-format text|markdown|tsv] <package dir>...")
		return 2
	}

	var funcs []Func
	scanned := 0
	for _, dir := range flags.Args() {
		found, err := scanDir(dir)
		if err != nil {
			fmt.Fprintf(stderr, "helperscan: skipping %s: %v\n", dir, err)
			continue
		}
		scanned++
		funcs = append(funcs, found...)
	}

	var err error
	switch *format {
	case "tsv":
		err = writeTSV(stdout, funcs)
	case "markdown":
		err = writeMarkdown(stdout, summarize(scanned, funcs), funcs)
	default:
		err = writeText(stdout, summarize(scanned, funcs), funcs)
	}
	if err != nil {
		fmt.Fprintf(stderr, "helperscan: %v\n", err)
	}
	return 0
}
