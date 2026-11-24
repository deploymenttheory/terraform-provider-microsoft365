package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// fixturePackage copies testdata/<name>/*.go.fixture into a temporary <name> directory as .go files and returns it.
// The fixtures are stored as .go.fixture so the pipeline's changed-package selection never treats them as packages to lint.
func fixturePackage(t *testing.T, name string) string {
	t.Helper()
	sources, err := filepath.Glob(filepath.Join("testdata", name, "*.go.fixture"))
	if err != nil || len(sources) == 0 {
		t.Fatalf("no fixtures for %s: %v", name, err)
	}
	dir := filepath.Join(t.TempDir(), name)
	if err := os.Mkdir(dir, 0o750); err != nil {
		t.Fatal(err)
	}
	for _, source := range sources {
		content, err := os.ReadFile(source)
		if err != nil {
			t.Fatal(err)
		}
		target := filepath.Join(dir, strings.TrimSuffix(filepath.Base(source), ".fixture"))
		if err := os.WriteFile(target, content, 0o600); err != nil {
			t.Fatal(err)
		}
	}
	return dir
}
