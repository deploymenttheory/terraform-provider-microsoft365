package graphBetaIdentityAndAccessRoleDefinition

import (
	"context"
	"sort"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// roleDefinitionConsistencyPredicate returns a consistency predicate for ReadWithRetry
// that verifies the role definition write has propagated across Microsoft Entra replicas
// before accepting the read as authoritative.
//
// Role definitions carry no modified timestamp, so the predicate compares the written values:
// display name, description, enabled flag, version and the allowed resource actions.
// A stale replica returns the pre-write role (or 404 on create, which ReadWithRetry retries).
//
// See: https://devblogs.microsoft.com/identity/designing-for-eventual-consistency-for-microsoft-entra/
func roleDefinitionConsistencyPredicate(expected *RoleDefinitionResourceModel) func(ctx context.Context, state tfsdk.State) bool {
	return func(ctx context.Context, state tfsdk.State) bool {
		var actual RoleDefinitionResourceModel
		if diags := state.Get(ctx, &actual); diags.HasError() {
			return false
		}

		if actual.ID.IsNull() || actual.ID.IsUnknown() || actual.ID.ValueString() == "" {
			return false
		}

		if !stringMatches(expected.DisplayName, actual.DisplayName) || !stringMatches(expected.Description, actual.Description) {
			return false
		}

		if !stringMatches(expected.Version, actual.Version) {
			return false
		}

		if !expected.IsEnabled.IsNull() && !expected.IsEnabled.IsUnknown() && expected.IsEnabled.ValueBool() != actual.IsEnabled.ValueBool() {
			return false
		}

		return permissionsKey(expected.RolePermissions) == permissionsKey(actual.RolePermissions)
	}
}

// stringMatches compares a planned value with the read value, treating null and empty as equal
func stringMatches(expected, actual types.String) bool {
	if expected.IsUnknown() {
		return true
	}
	return expected.ValueString() == actual.ValueString()
}

// permissionsKey builds an order-independent, case-insensitive key for a set of role permissions
func permissionsKey(permissions []RoleDefinitionPermission) string {
	entries := make([]string, 0, len(permissions))
	for _, permission := range permissions {
		actions := make([]string, 0, len(permission.AllowedResourceActions.Elements()))
		for _, element := range permission.AllowedResourceActions.Elements() {
			if action, ok := element.(types.String); ok {
				actions = append(actions, strings.ToLower(action.ValueString()))
			}
		}
		sort.Strings(actions)
		entries = append(entries, strings.Join(actions, ","))
	}
	sort.Strings(entries)
	return strings.Join(entries, "|")
}
