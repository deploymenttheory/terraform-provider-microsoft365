package runner

import (
	"errors"
	"fmt"
)

// TestType represents the category of tests to execute.
type TestType string

const (
	TestTypeProviderCore TestType = "provider-core"
	TestTypeResources    TestType = "resources"
	TestTypeDatasources  TestType = "datasources"
)

// Valid test types for validation.
var validTestTypes = map[TestType]bool{
	TestTypeProviderCore: true,
	TestTypeResources:    true,
	TestTypeDatasources:  true,
}

// Config holds the test runner configuration.
type Config struct {
	TestType       string
	Service        string
	CoverageOutput string
	Verbose        bool
}

// Validate ensures the configuration is valid.
func (c *Config) Validate() error {
	if c.TestType == "" {
		return errors.New("test type is required")
	}

	if !validTestTypes[TestType(c.TestType)] {
		return fmt.Errorf("invalid test type: %s (must be provider-core, resources, or datasources)", c.TestType)
	}

	testType := TestType(c.TestType)
	if (testType == TestTypeResources || testType == TestTypeDatasources) && c.Service == "" {
		return fmt.Errorf("service is required for test type: %s", testType)
	}

	if c.CoverageOutput == "" {
		return errors.New("coverage output path is required")
	}

	return nil
}

// GetTestType returns the parsed TestType.
func (c *Config) GetTestType() TestType {
	return TestType(c.TestType)
}

// RequiresService returns true if this test type requires a service parameter.
func (c *Config) RequiresService() bool {
	testType := c.GetTestType()
	return testType == TestTypeResources || testType == TestTypeDatasources
}
