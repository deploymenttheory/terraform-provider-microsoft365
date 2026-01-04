package report

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/action"
)

// When creates a validation reporter that reports via SendProgress
// Usage: report.When(resp).Error("message").Warning("message").FailIfErrors()
func When(resp *action.InvokeResponse) *Reporter {
	return &Reporter{
		resp:      resp,
		hasErrors: false,
	}
}

// Reporter reports validation findings via SendProgress events
type Reporter struct {
	resp      *action.InvokeResponse
	hasErrors bool
}

// Error reports a validation error via SendProgress and marks validation as failed
// Returns self for chaining
// Usage: .Error("Validation failed: %d device(s) do not exist: %s", len(ids), strings.Join(ids, ", "))
func (r *Reporter) Error(format string, args ...interface{}) *Reporter {
	r.resp.SendProgress(action.InvokeProgressEvent{
		Message: fmt.Sprintf(format, args...),
	})
	r.hasErrors = true
	return r
}

// Warning reports a validation warning via SendProgress (doesn't fail validation)
// Returns self for chaining
// Usage: .Warning("Warning: %d device(s) may already be enrolled: %s", len(ids), strings.Join(ids, ", "))
func (r *Reporter) Warning(format string, args ...interface{}) *Reporter {
	r.resp.SendProgress(action.InvokeProgressEvent{
		Message: fmt.Sprintf(format, args...),
	})
	return r
}

// FailIfErrors adds diagnostic error if validation failed
// Returns true if errors were found
func (r *Reporter) FailIfErrors() bool {
	if r.hasErrors {
		r.resp.Diagnostics.AddError(
			"Validation Failed",
			"One or more validation checks failed",
		)
		return true
	}
	return false
}
