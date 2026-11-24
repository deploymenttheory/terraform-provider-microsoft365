package main

import (
	"fmt"
	"io"
	"strconv"
	"strings"
	"text/tabwriter"
)

// summary holds the package-wide counts printed above each report.
type summary struct {
	Packages        int
	Total           int
	Kiota           int
	Framework       int
	Helpers         int
	StateMappers    int
	SingleRef       int
	Short           int
	ShortSingleRef  int
	FewCallers      int
	ShortFewCallers int
	NameFails       int
	Pass            int
}

// summarize counts every non-standard function, then the helpers against each part of the helper rule.
// State mappers are counted on their own and left out of the other figures, which keeps them comparable with the #4929 baseline.
func summarize(packages int, funcs []Func) summary {
	s := summary{Packages: packages}
	for _, f := range funcs {
		if f.Kind == kindStateMapper {
			s.StateMappers++
			continue
		}
		s.Total++
		short := f.Lines <= maxShortLines
		if f.Refs <= 1 {
			s.SingleRef++
			if short {
				s.ShortSingleRef++
			}
		}
		if short {
			s.Short++
		}
		switch f.Kind {
		case kindKiota:
			s.Kiota++
			continue
		case kindFramework:
			s.Framework++
			continue
		}
		s.Helpers++
		if len(f.CRUDCallers) < minCRUDCallers {
			s.FewCallers++
			if short {
				s.ShortFewCallers++
			}
		}
		if f.NameCheck != nameOK {
			s.NameFails++
		}
		if f.Rule() == "pass" {
			s.Pass++
		}
	}
	return s
}

// lines renders the summary as plain sentences shared by the text and Markdown reports.
func (s summary) lines() []string {
	callers := fmt.Sprintf("fewer than %d CRUD callers %d (%d of them also short)", minCRUDCallers, s.FewCallers, s.ShortFewCallers)
	return []string{
		fmt.Sprintf("Packages scanned: %d", s.Packages),
		fmt.Sprintf("Non-standard functions: %d (Kiota interface methods %d, framework validator methods %d, helpers %d)", s.Total, s.Kiota, s.Framework, s.Helpers),
		fmt.Sprintf("All non-standard functions: referenced at most once %d; %d lines or fewer %d; both %d", s.SingleRef, maxShortLines, s.Short, s.ShortSingleRef),
		fmt.Sprintf("Helpers: %s; name fails convention %d; meet the helper rule %d", callers, s.NameFails, s.Pass),
		fmt.Sprintf("State mappers (%s, exempt from the caller and line rules): %d", stateMapperName, s.StateMappers),
	}
}

// crudList renders the CRUD callers as method names, or "-" when there are none.
func crudList(f Func) string {
	if len(f.CRUDCallers) == 0 {
		return "-"
	}
	methods := make([]string, len(f.CRUDCallers))
	for i, caller := range f.CRUDCallers {
		_, method, _ := strings.Cut(caller, ".")
		methods[i] = method
	}
	return strings.Join(methods, ",")
}

// displayName prefixes methods with their receiver type.
func displayName(f Func) string {
	if f.Receiver == "" {
		return f.Name
	}
	return f.Receiver + "." + f.Name
}

// writeText prints the summary and an aligned table of helpers for a CI job log.
func writeText(w io.Writer, s summary, funcs []Func) error {
	var b strings.Builder
	b.WriteString("Helper report (report-only, see docs/development/guide.md#helpers)\n")
	for _, line := range s.lines() {
		b.WriteString(line + "\n")
	}
	b.WriteString("\n")
	tw := tabwriter.NewWriter(&b, 0, 4, 2, ' ', 0)
	fmt.Fprintln(tw, "PACKAGE\tFILE\tFUNC\tLINES\tREFS\tCRUD CALLERS\tNAME\tRULE")
	for _, f := range funcs {
		if !inTable(f) {
			continue
		}
		fmt.Fprintf(tw, "%s\t%s\t%s\t%d\t%d\t%d %s\t%s\t%s\n", f.Package, f.File, displayName(f), f.Lines, f.Refs, len(f.CRUDCallers), crudList(f), f.NameCheck, f.Rule())
	}
	if err := tw.Flush(); err != nil {
		return fmt.Errorf("format table: %w", err)
	}
	return writeString(w, b.String())
}

// writeTSV prints every non-standard function, excluded kinds included, as tab-separated rows.
func writeTSV(w io.Writer, funcs []Func) error {
	var b strings.Builder
	b.WriteString("package\tfile\tfunc\treceiver\tkind\tlines\trefs\tcrud_callers\tcrud_methods\tname\trule\n")
	for _, f := range funcs {
		row := []string{
			f.Package, f.File, f.Name, f.Receiver, f.Kind,
			strconv.Itoa(f.Lines), strconv.Itoa(f.Refs), strconv.Itoa(len(f.CRUDCallers)), crudList(f), f.NameCheck, f.Rule(),
		}
		b.WriteString(strings.Join(row, "\t") + "\n")
	}
	return writeString(w, b.String())
}

// writeMarkdown prints the summary and helper table for a GitHub step summary.
func writeMarkdown(w io.Writer, s summary, funcs []Func) error {
	var b strings.Builder
	b.WriteString("## Helper report\n\nReport-only. The rules are in `docs/development/guide.md` under Helpers.\n\n")
	for _, line := range s.lines() {
		b.WriteString("- " + line + "\n")
	}
	b.WriteString("\n| Package | File | Func | Lines | Refs | CRUD callers | Name | Rule |\n| --- | --- | --- | ---: | ---: | --- | --- | --- |\n")
	for _, f := range funcs {
		if !inTable(f) {
			continue
		}
		fmt.Fprintf(&b, "| %s | %s | `%s` | %d | %d | %d %s | %s | %s |\n", f.Package, f.File, displayName(f), f.Lines, f.Refs, len(f.CRUDCallers), crudList(f), f.NameCheck, f.Rule())
	}
	return writeString(w, b.String())
}

// writeString writes s to w and wraps any error.
func writeString(w io.Writer, s string) error {
	if _, err := io.WriteString(w, s); err != nil {
		return fmt.Errorf("write report: %w", err)
	}
	return nil
}

// inTable reports whether a function appears in the text and Markdown tables: helpers and the state mapper.
func inTable(f Func) bool {
	return f.Kind == kindHelper || f.Kind == kindStateMapper
}
