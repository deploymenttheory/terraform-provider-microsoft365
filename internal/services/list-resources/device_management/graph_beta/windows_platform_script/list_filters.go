package graphBetaWindowsPlatformScript

import (
	"fmt"
	"strings"
)

// buildDisplayNameFilter constructs an OData filter for partial display name matching
func buildDisplayNameFilter(data *WindowsPlatformScriptListConfigModel) string {
	if data.DisplayNameFilter.IsNull() || data.DisplayNameFilter.IsUnknown() {
		return ""
	}
	return fmt.Sprintf("contains(displayName,'%s')", data.DisplayNameFilter.ValueString())
}

// buildFileNameFilter constructs an OData filter for partial file name matching
func buildFileNameFilter(data *WindowsPlatformScriptListConfigModel) string {
	if data.FileNameFilter.IsNull() || data.FileNameFilter.IsUnknown() {
		return ""
	}
	return fmt.Sprintf("contains(fileName,'%s')", data.FileNameFilter.ValueString())
}

// buildRunAsAccountFilter constructs an OData filter for run as account
func buildRunAsAccountFilter(data *WindowsPlatformScriptListConfigModel) string {
	if data.RunAsAccountFilter.IsNull() || data.RunAsAccountFilter.IsUnknown() {
		return ""
	}
	return fmt.Sprintf("runAsAccount eq '%s'", data.RunAsAccountFilter.ValueString())
}

// combineFilters joins non-empty filter parts with AND logic
func combineFilters(filters ...string) string {
	var parts []string
	for _, filter := range filters {
		if filter != "" {
			parts = append(parts, filter)
		}
	}
	if len(parts) == 0 {
		return ""
	}
	return strings.Join(parts, " and ")
}
