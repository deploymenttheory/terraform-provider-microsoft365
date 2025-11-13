package handler

import (
	"errors"
	"strings"
	"time"
)

// Config holds the failure handler configuration.
type Config struct {
	Owner      string
	Repo       string
	RunID      string
	FailedJobs string
}

// Validate ensures the configuration is complete.
func (c *Config) Validate() error {
	if c.Owner == "" {
		return errors.New("repository owner is required")
	}
	if c.Repo == "" {
		return errors.New("repository name is required")
	}
	if c.RunID == "" {
		return errors.New("workflow run ID is required")
	}
	if c.FailedJobs == "" {
		return errors.New("failed jobs list is required")
	}
	return nil
}

// ParseFailedJobs returns the failed jobs as a slice.
func (c *Config) ParseFailedJobs() []string {
	jobs := strings.Split(c.FailedJobs, ",")
	result := make([]string, 0, len(jobs))
	for _, job := range jobs {
		if trimmed := strings.TrimSpace(job); trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

// GetBranchName generates a unique branch name for the failure report.
func (c *Config) GetBranchName() string {
	timestamp := time.Now().Format("20060102-150405")
	return "nightly-test-failures-" + timestamp
}

// GetDate returns the current date in YYYY-MM-DD format.
func (c *Config) GetDate() string {
	return time.Now().Format("2006-01-02")
}

// GetWorkflowURL returns the full URL to the workflow run.
func (c *Config) GetWorkflowURL() string {
	return buildURL("https://github.com", c.Owner, c.Repo, "actions/runs", c.RunID)
}

// GetCodecovURL returns the full URL to the Codecov dashboard.
func (c *Config) GetCodecovURL() string {
	return buildURL("https://codecov.io/gh", c.Owner, c.Repo)
}

// buildURL constructs a URL from path segments.
func buildURL(base string, segments ...string) string {
	parts := append([]string{base}, segments...)
	return strings.Join(parts, "/")
}
