package main

import (
	"bytes"
	"errors"
	"path/filepath"
	"strings"
	"testing"
)

func sampleFuncs(t *testing.T) []Func {
	t.Helper()
	funcs, err := scanDir(filepath.Join("testdata", "sample"))
	if err != nil {
		t.Fatal(err)
	}
	return funcs
}

func TestSummarize(t *testing.T) {
	got := summarize(1, sampleFuncs(t))
	want := summary{Packages: 1, Total: 14, Kiota: 4, Framework: 3, Helpers: 7, SingleRef: 10, Short: 12, ShortSingleRef: 10, FewCallers: 3, ShortFewCallers: 3, NameFails: 4, Pass: 2}
	if got != want {
		t.Errorf("summarize = %+v\nwant        %+v", got, want)
	}
}

func TestWriteText(t *testing.T) {
	funcs := sampleFuncs(t)
	var out bytes.Buffer
	if err := writeText(&out, summarize(1, funcs), funcs); err != nil {
		t.Fatal(err)
	}
	text := out.String()
	for _, want := range []string{"Packages scanned: 1", "PACKAGE", "constructWidgetRequest", "3 Read,Create,Update", "SampleResource.refreshWidget", "fails: callers,lines,name"} {
		if !strings.Contains(text, want) {
			t.Errorf("text report missing %q:\n%s", want, text)
		}
	}
	if strings.Contains(text, "Serialize") || strings.Contains(text, "ValidateString") {
		t.Errorf("text report lists excluded functions:\n%s", text)
	}
}

func TestWriteMarkdown(t *testing.T) {
	funcs := sampleFuncs(t)
	var out bytes.Buffer
	if err := writeMarkdown(&out, summarize(1, funcs), funcs); err != nil {
		t.Fatal(err)
	}
	md := out.String()
	for _, want := range []string{"## Helper report", "- Packages scanned: 1", "| sample | construct.go | `constructWidgetRequest` | 11 | 3 | 3 Read,Create,Update | ok | pass |"} {
		if !strings.Contains(md, want) {
			t.Errorf("markdown report missing %q:\n%s", want, md)
		}
	}
}

func TestWriteTSV(t *testing.T) {
	funcs := sampleFuncs(t)
	var out bytes.Buffer
	if err := writeTSV(&out, funcs); err != nil {
		t.Fatal(err)
	}
	rows := strings.Split(strings.TrimSpace(out.String()), "\n")
	if len(rows) != len(funcs)+1 {
		t.Fatalf("got %d rows, want %d", len(rows), len(funcs)+1)
	}
	if !strings.Contains(out.String(), "sample\tmodel.go\tSerialize\twidgetResponse\tkiota\t3\t0\t0\t-\t\texcluded\n") {
		t.Errorf("TSV missing the excluded Kiota row:\n%s", out.String())
	}
}

type failingWriter struct{}

var errWrite = errors.New("write failed")

func (failingWriter) Write([]byte) (int, error) { return 0, errWrite }

func TestWritersReportWriteErrors(t *testing.T) {
	funcs := sampleFuncs(t)
	s := summarize(1, funcs)
	errs := map[string]error{
		"text":     writeText(failingWriter{}, s, funcs),
		"markdown": writeMarkdown(failingWriter{}, s, funcs),
		"tsv":      writeTSV(failingWriter{}, funcs),
	}
	for name, err := range errs {
		if !errors.Is(err, errWrite) {
			t.Errorf("%s: got %v, want the write error", name, err)
		}
	}
}
