package runner

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
)

// Package paths for provider core tests.
var providerCorePackages = []string{
	"./internal/client/...",
	"./internal/helpers/...",
	"./internal/provider/...",
	"./internal/utilities/...",
}

// Runner executes Go tests with coverage tracking.
type Runner struct {
	config   *Config
	credProv CredentialProvider
	stdout   io.Writer
	stderr   io.Writer
}

// New creates a new test runner with the given configuration and credential provider.
func New(cfg *Config, credProv CredentialProvider) *Runner {
	return &Runner{
		config:   cfg,
		credProv: credProv,
		stdout:   os.Stdout,
		stderr:   os.Stderr,
	}
}

// WithOutput sets custom output writers for stdout and stderr.
func (r *Runner) WithOutput(stdout, stderr io.Writer) *Runner {
	r.stdout = stdout
	r.stderr = stderr
	return r
}

// Run executes the test suite based on configuration.
func (r *Runner) Run() error {
	creds, err := r.credProv.GetCredentials()
	if err != nil {
		return fmt.Errorf("failed to retrieve credentials: %w", err)
	}

	if !creds.IsValid() {
		fmt.Fprintf(r.stdout, "⚠️  No credentials configured for %s, skipping tests\n", r.config.Service)
		return createEmptyCoverageFile(r.config.CoverageOutput)
	}

	fmt.Fprintf(r.stdout, "✅ Credentials found for %s\n", r.config.Service)

	switch r.config.GetTestType() {
	case TestTypeProviderCore:
		return r.runProviderCoreTests()
	case TestTypeResources:
		return r.runServiceTests("resources", r.config.Service)
	case TestTypeDatasources:
		return r.runServiceTests("datasources", r.config.Service)
	default:
		return fmt.Errorf("unknown test type: %s", r.config.TestType)
	}
}

// runProviderCoreTests executes tests for provider core packages.
func (r *Runner) runProviderCoreTests() error {
	fmt.Fprintln(r.stdout, "Running provider core tests...")

	args := []string{"test", "-v", "-race"}
	args = append(args, buildCoverageArgs(r.config.CoverageOutput)...)
	args = append(args, providerCorePackages...)

	return r.executeGoTest(args)
}

// runServiceTests executes tests for a specific service area.
func (r *Runner) runServiceTests(category, service string) error {
	fmt.Fprintf(r.stdout, "Running tests for %s/%s...\n", category, service)

	testDir := filepath.Join(".", "internal", "services", category, service)

	exists, err := directoryExists(testDir)
	if err != nil {
		return fmt.Errorf("failed to check test directory: %w", err)
	}

	if !exists {
		fmt.Fprintf(r.stdout, "Directory not found: %s, skipping\n", testDir)
		return createEmptyCoverageFile(r.config.CoverageOutput)
	}

	hasTests, err := hasTestFiles(testDir)
	if err != nil {
		return fmt.Errorf("failed to check for test files: %w", err)
	}

	if !hasTests {
		fmt.Fprintf(r.stdout, "No test files found in %s, creating empty coverage file\n", testDir)
		return createEmptyCoverageFile(r.config.CoverageOutput)
	}

	testCount, _ := countTestFiles(testDir)
	fmt.Fprintf(r.stdout, "Found %d test files\n", testCount)

	args := []string{"test", "-v", "-race"}
	args = append(args, buildCoverageArgs(r.config.CoverageOutput)...)
	args = append(args, "./"+testDir+"/...")

	return r.executeGoTest(args)
}

// executeGoTest runs the go test command with the given arguments.
func (r *Runner) executeGoTest(args []string) error {
	cmd := exec.Command("go", args...)
	cmd.Stdout = r.stdout
	cmd.Stderr = r.stderr
	cmd.Env = os.Environ()

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("test execution failed: %w", err)
	}

	return nil
}

// buildCoverageArgs constructs coverage-related command-line arguments.
func buildCoverageArgs(outputPath string) []string {
	return []string{
		"-coverprofile=" + outputPath,
		"-covermode=atomic",
	}
}

// createEmptyCoverageFile creates a minimal valid Go coverage file.
func createEmptyCoverageFile(path string) error {
	return os.WriteFile(path, []byte("mode: atomic\n"), 0644)
}

// directoryExists checks if a directory exists.
func directoryExists(path string) (bool, error) {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return info.IsDir(), nil
}

// hasTestFiles checks if a directory contains any Go test files.
func hasTestFiles(dir string) (bool, error) {
	files, err := filepath.Glob(filepath.Join(dir, "*_test.go"))
	if err != nil {
		return false, err
	}
	return len(files) > 0, nil
}

// countTestFiles returns the number of test files in a directory.
func countTestFiles(dir string) (int, error) {
	files, err := filepath.Glob(filepath.Join(dir, "*_test.go"))
	if err != nil {
		return 0, err
	}
	return len(files), nil
}
