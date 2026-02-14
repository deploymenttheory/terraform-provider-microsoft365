package graphBetaUsersUser

import (
	"fmt"
	"strings"
)

// buildDisplayNameFilter constructs an OData filter for partial display name matching
func buildDisplayNameFilter(data *UserListConfigModel) string {
	if data.DisplayNameFilter.IsNull() || data.DisplayNameFilter.IsUnknown() {
		return ""
	}
	return fmt.Sprintf("startsWith(displayName,'%s')", data.DisplayNameFilter.ValueString())
}

// buildUserPrincipalNameFilter constructs an OData filter for user principal name matching
func buildUserPrincipalNameFilter(data *UserListConfigModel) string {
	if data.UserPrincipalNameFilter.IsNull() || data.UserPrincipalNameFilter.IsUnknown() {
		return ""
	}
	return fmt.Sprintf("startsWith(userPrincipalName,'%s')", data.UserPrincipalNameFilter.ValueString())
}

// buildAccountEnabledFilter constructs an OData filter for account enabled status
func buildAccountEnabledFilter(data *UserListConfigModel) string {
	if data.AccountEnabledFilter.IsNull() || data.AccountEnabledFilter.IsUnknown() {
		return ""
	}
	return fmt.Sprintf("accountEnabled eq %t", data.AccountEnabledFilter.ValueBool())
}

// buildUserTypeFilter constructs an OData filter for user type
func buildUserTypeFilter(data *UserListConfigModel) string {
	if data.UserTypeFilter.IsNull() || data.UserTypeFilter.IsUnknown() {
		return ""
	}
	return fmt.Sprintf("userType eq '%s'", data.UserTypeFilter.ValueString())
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
