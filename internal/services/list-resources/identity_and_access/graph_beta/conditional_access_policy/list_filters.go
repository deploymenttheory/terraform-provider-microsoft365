package graphBetaConditionalAccessPolicy

import (
	"fmt"
	"strings"
)

// buildDisplayNameFilter constructs an OData filter for partial display name matching
func buildDisplayNameFilter(data *ConditionalAccessPolicyListConfigModel) string {
	if data.DisplayNameFilter.IsNull() || data.DisplayNameFilter.IsUnknown() {
		return ""
	}
	return fmt.Sprintf("contains(displayName,'%s')", data.DisplayNameFilter.ValueString())
}

// buildStateFilter constructs an OData filter for state
func buildStateFilter(data *ConditionalAccessPolicyListConfigModel) string {
	if data.StateFilter.IsNull() || data.StateFilter.IsUnknown() {
		return ""
	}
	return fmt.Sprintf("state eq '%s'", data.StateFilter.ValueString())
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
