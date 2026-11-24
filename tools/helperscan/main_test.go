package main

import (
	"bytes"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunFormats(t *testing.T) {
	sample := filepath.Join("testdata", "sample")
	cases := map[string]string{"text": "Packages scanned: 1", "markdown": "## Helper report", "tsv": "package\tfile\tfunc"}
	for format, want := range cases {
		var stdout, stderr bytes.Buffer
		if code := run([]string{"-format", format, sample}, &stdout, &stderr); code != 0 {
			t.Fatalf("%s: exit %d, stderr %s", format, code, stderr.String())
		}
		if !strings.Contains(stdout.String(), want) {
			t.Errorf("%s: output missing %q:\n%s", format, want, stdout.String())
		}
	}
}

func TestRunSkipsUnreadablePackages(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := run([]string{filepath.Join("testdata", "broken"), filepath.Join("testdata", "sample")}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("exit %d, want 0", code)
	}
	if !strings.Contains(stderr.String(), "skipping") || !strings.Contains(stdout.String(), "Packages scanned: 1") {
		t.Errorf("stdout %q stderr %q", stdout.String(), stderr.String())
	}
}

func TestRunWriteErrorStillExitsZero(t *testing.T) {
	var stderr bytes.Buffer
	if code := run([]string{filepath.Join("testdata", "sample")}, failingWriter{}, &stderr); code != 0 || !strings.Contains(stderr.String(), "write failed") {
		t.Errorf("exit %d, stderr %q", code, stderr.String())
	}
}

func TestRunUsageErrors(t *testing.T) {
	for _, args := range [][]string{nil, {"-format", "xml", "."}, {"-unknown"}} {
		var stdout, stderr bytes.Buffer
		if code := run(args, &stdout, &stderr); code != 2 {
			t.Errorf("run(%q) = %d, want 2", args, code)
		}
	}
}
