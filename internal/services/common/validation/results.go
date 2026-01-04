// Package validation provides common validation result handling for actions
package validation

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/action"
)

// Finding represents a single validation finding with auto-formatting
type Finding struct {
	Level      string   // "error" or "warning"
	IDs        []string // Device/resource IDs
	EntityType string   // e.g., "managed device", "co-managed device", "user"
	Message    string   // The validation message template
}

// Report sends this finding via SendProgress with proper formatting
func (f *Finding) Report(resp *action.InvokeResponse) {
	if len(f.IDs) == 0 {
		return
	}

	var msg string
	if f.EntityType != "" {
		msg = fmt.Sprintf("%s: %d %s(s): %s - IDs: %s",
			f.prefixForLevel(),
			len(f.IDs),
			f.EntityType,
			f.Message,
			strings.Join(f.IDs, ", "))
	} else {
		msg = fmt.Sprintf("%s: %d item(s): %s - IDs: %s",
			f.prefixForLevel(),
			len(f.IDs),
			f.Message,
			strings.Join(f.IDs, ", "))
	}

	resp.SendProgress(action.InvokeProgressEvent{
		Message: msg,
	})
}

func (f *Finding) prefixForLevel() string {
	if f.Level == "warning" {
		return "Warning"
	}
	return "Validation failed"
}

// IsError returns true if this is an error-level finding
func (f *Finding) IsError() bool {
	return f.Level == "error"
}

// Results is a collection of validation findings
type Results struct {
	findings []Finding
}

// NewResults creates a new validation results collection
func NewResults() *Results {
	return &Results{
		findings: make([]Finding, 0),
	}
}

// Add adds a finding to the results
func (r *Results) Add(level string, ids []string, entityType, message string) *Results {
	if len(ids) > 0 {
		r.findings = append(r.findings, Finding{
			Level:      level,
			IDs:        ids,
			EntityType: entityType,
			Message:    message,
		})
	}
	return r
}

// Error adds an error-level finding
func (r *Results) Error(ids []string, entityType, message string) *Results {
	return r.Add("error", ids, entityType, message)
}

// Warning adds a warning-level finding
func (r *Results) Warning(ids []string, entityType, message string) *Results {
	return r.Add("warning", ids, entityType, message)
}

// Report reports all findings via SendProgress and returns true if errors found
func (r *Results) Report(resp *action.InvokeResponse) bool {
	hasErrors := false

	for _, finding := range r.findings {
		finding.Report(resp)
		if finding.IsError() {
			hasErrors = true
		}
	}

	if hasErrors {
		resp.Diagnostics.AddError(
			"Validation Failed",
			"One or more validation checks failed",
		)
	}

	return hasErrors
}
