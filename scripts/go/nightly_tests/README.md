# Nightly Test Scripts

Idiomatic Go tools for managing nightly acceptance tests for the Microsoft 365 Terraform Provider.

## Overview

These tools follow Go best practices with proper package structure, interfaces, constants, and reusable functions. Each tool is designed to be testable, maintainable, and follow Go idioms.

## Architecture

```
scripts/go/nightly_tests/
├── testrunner/
│   ├── main.go                    # CLI entry point
│   ├── go.mod                     # Module definition
│   └── runner/
│       ├── types.go               # Config and TestType definitions
│       ├── credentials.go         # CredentialProvider interface
│       └── runner.go              # Test execution logic
├── credmapper/
│   ├── main.go                    # CLI entry point
│   ├── go.mod                     # Module definition
│   └── mapper/
│       ├── credentials.go         # Service credentials constants
│       └── mapper.go              # Mapper and Exporter interface
├── coveragemerger/
│   ├── main.go                    # CLI entry point
│   ├── go.mod                     # Module definition
│   └── merger/
│       └── merger.go              # Coverage merge logic
└── failurehandler/
    ├── main.go                    # CLI entry point
    ├── go.mod                     # Module definition
    └── handler/
        ├── types.go               # Config and validation
        ├── handler.go             # Failure handling logic
        └── templates.go           # Report template builders
```

## Design Principles

### 1. Constants and Variables

All fixed values are defined as package-level constants:

```go
const (
    EnvClientID     = "M365_CLIENT_ID"
    EnvClientSecret = "M365_CLIENT_SECRET"
    TestTypeProviderCore = "provider-core"
)
```

### 2. Interfaces for Abstraction

Interfaces enable testability and flexibility:

```go
type CredentialProvider interface {
    GetCredentials() (*Credentials, error)
}

type Exporter interface {
    Export(key, value string) error
    HasCapability() bool
}
```

### 3. Separation of Concerns

- `main.go`: CLI parsing and setup
- Package layer: Business logic
- Types: Configuration and validation

### 4. Reusable Functions

Generic, composable functions:

```go
func directoryExists(path string) (bool, error)
func hasTestFiles(dir string) (bool, error)
func buildCoverageArgs(outputPath string) []string
```

### 5. Proper Error Handling

Errors are wrapped with context:

```go
if err != nil {
    return fmt.Errorf("failed to create output file: %w", err)
}
```

### 6. Testable Design

Logic separated from I/O with dependency injection:

```go
func New(cfg *Config, credProv CredentialProvider) *Runner
func (r *Runner) WithOutput(stdout, stderr io.Writer) *Runner
```

## Tools

### testrunner

Executes tests with proper credential validation and coverage tracking.

**Key Types:**

```go
type Config struct {
    TestType       string
    Service        string
    CoverageOutput string
    Verbose        bool
}

type Runner struct {
    config   *Config
    credProv CredentialProvider
    stdout   io.Writer
    stderr   io.Writer
}
```

**Usage:**

```bash
./testrunner -type=provider-core -coverage=coverage.txt
./testrunner -type=resources -service=applications -coverage=coverage.txt
```

### credmapper

Maps service names to credentials using explicit constants.

**Key Constants:**

```go
const (
    ServiceApplications      = "applications"
    ServiceBackupStorage     = "backup_storage"
    EnvApplicationsClientID  = "M365_CLIENT_ID_APPLICATIONS"
    // ... 11 services mapped
)
```

**Key Types:**

```go
type ServiceCredential struct {
    ClientIDVar     string
    ClientSecretVar string
}

type Mapper struct {
    exporter Exporter
}
```

**Usage:**

```bash
./credmapper -service=applications
```

### coveragemerger

Merges Go coverage files with proper stream processing.

**Key Constants:**

```go
const (
    CoverageFileExt       = ".txt"
    CoverageMode          = "mode: atomic"
    DefaultFilePermission = 0644
)
```

**Key Types:**

```go
type Merger struct {
    inputDir   string
    outputPath string
    stdout     io.Writer
}
```

**Usage:**

```bash
./coveragemerger -input=coverage-artifacts -output=coverage-merged.txt
```

### failurehandler

Creates PRs and issues using the GitHub SDK with proper configuration.

**Key Constants:**

```go
const (
    LabelBug            = "bug"
    LabelTesting        = "testing"
    LabelAutomated      = "automated"
    LabelNightlyFailure = "nightly-failure"
)
```

**Key Types:**

```go
type Config struct {
    Owner      string
    Repo       string
    RunID      string
    FailedJobs string
}

type Handler struct {
    client *github.Client
    config *Config
}
```

**Usage:**

```bash
GITHUB_TOKEN=xxx ./failurehandler \
  -owner=deploymenttheory \
  -repo=terraform-provider-microsoft365 \
  -run-id=123456789 \
  -failed-jobs="Provider-Core,Resources"
```

## Testing

Each package is designed for testability:

```go
// Mock credential provider for testing
type MockCredentialProvider struct {
    Creds *Credentials
    Err   error
}

func (m *MockCredentialProvider) GetCredentials() (*Credentials, error) {
    return m.Creds, m.Err
}

// Test with mock
func TestRunner_Run(t *testing.T) {
    mockCreds := &MockCredentialProvider{
        Creds: &Credentials{ClientID: "test", ClientSecret: "secret"},
    }
    runner := New(cfg, mockCreds)
    // ... test
}
```

## Benefits of Idiomatic Go

1. **Type Safety**: Constants prevent typos and invalid values
2. **Testability**: Interfaces enable mocking and unit testing
3. **Maintainability**: Clear separation of concerns
4. **Reusability**: Generic functions can be used across tools
5. **Documentation**: Types and constants are self-documenting
6. **Error Context**: Wrapped errors provide clear failure paths
7. **Performance**: Efficient use of buffers and streams

## Dependencies

- **Standard Library**: Most tools use only stdlib
- **failurehandler**:
  - `github.com/google/go-github/v66` - GitHub API client
  - `golang.org/x/oauth2` - OAuth2 authentication

## Building

Each tool is an independent module:

```bash
cd testrunner && go build
cd credmapper && go build
cd coveragemerger && go build
cd failurehandler && go build
```

## Contributing

When adding new functionality:

1. Define constants for fixed values
2. Create interfaces for pluggable behavior
3. Separate configuration from logic
4. Make functions small and reusable
5. Add proper error wrapping
6. Consider testability in design
7. Follow Go naming conventions (exported vs unexported)
