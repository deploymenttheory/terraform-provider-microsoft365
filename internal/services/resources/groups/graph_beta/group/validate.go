package graphBetaGroup

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// validateVisibility ensures visibility constraints are met according to Microsoft Graph API rules
func validateVisibility(ctx context.Context, data *GroupResourceModel, diagnostics *diag.Diagnostics) {
	// Skip validation if visibility is null or unknown
	if data.Visibility.IsNull() || data.Visibility.IsUnknown() {
		return
	}

	visibility := data.Visibility.ValueString()

	// Validate HiddenMembership constraint
	if visibility == "HiddenMembership" {
		// HiddenMembership can only be set for Microsoft 365 groups (Unified groups)
		if !data.GroupTypes.IsNull() && !data.GroupTypes.IsUnknown() {
			groupTypes := make([]string, 0, len(data.GroupTypes.Elements()))
			diag := data.GroupTypes.ElementsAs(ctx, &groupTypes, false)
			if diag.HasError() {
				diagnostics.Append(diag...)
				return
			}

			// Check if this is a Microsoft 365 group (contains "Unified")
			isUnifiedGroup := false
			for _, gt := range groupTypes {
				if gt == "Unified" {
					isUnifiedGroup = true
					break
				}
			}

			if !isUnifiedGroup {
				diagnostics.AddAttributeError(
					path.Root("visibility"),
					"Invalid Visibility Value",
					"The visibility value 'HiddenMembership' can only be set for Microsoft 365 groups. "+
						"Ensure 'group_types' includes 'Unified' when using 'HiddenMembership' visibility.",
				)
				return
			}
		}
	}

	// Validate role assignable groups must be Private
	if !data.IsAssignableToRole.IsNull() && !data.IsAssignableToRole.IsUnknown() {
		if data.IsAssignableToRole.ValueBool() && visibility != "Private" {
			diagnostics.AddAttributeError(
				path.Root("visibility"),
				"Invalid Visibility for Role-Assignable Group",
				"Groups assignable to roles must have visibility set to 'Private'. "+
					fmt.Sprintf("Current visibility is '%s'.", visibility),
			)
		}
	}
}

// validateRoleAssignability ensures role assignability constraints are met
func validateRoleAssignability(ctx context.Context, data *GroupResourceModel, diagnostics *diag.Diagnostics) {
	// Skip validation if is_assignable_to_role is null or unknown
	if data.IsAssignableToRole.IsNull() || data.IsAssignableToRole.IsUnknown() {
		return
	}

	// Only validate if is_assignable_to_role is true
	if !data.IsAssignableToRole.ValueBool() {
		return
	}

	// Rule 1: securityEnabled must be true
	if !data.SecurityEnabled.IsNull() && !data.SecurityEnabled.IsUnknown() {
		if !data.SecurityEnabled.ValueBool() {
			diagnostics.AddAttributeError(
				path.Root("is_assignable_to_role"),
				"Invalid Role-Assignable Group Configuration",
				"When 'is_assignable_to_role' is true, 'security_enabled' must also be true.",
			)
		}
	}

	// Rule 2: visibility must be Private
	if !data.Visibility.IsNull() && !data.Visibility.IsUnknown() {
		if data.Visibility.ValueString() != "Private" {
			diagnostics.AddAttributeError(
				path.Root("is_assignable_to_role"),
				"Invalid Role-Assignable Group Configuration",
				fmt.Sprintf("When 'is_assignable_to_role' is true, 'visibility' must be 'Private'. Current visibility is '%s'.",
					data.Visibility.ValueString()),
			)
		}
	}

	// Rule 3: group cannot be dynamic (groupTypes cannot contain DynamicMembership)
	if !data.GroupTypes.IsNull() && !data.GroupTypes.IsUnknown() {
		groupTypes := make([]string, 0, len(data.GroupTypes.Elements()))
		diag := data.GroupTypes.ElementsAs(ctx, &groupTypes, false)
		if diag.HasError() {
			diagnostics.Append(diag...)
			return
		}

		for _, gt := range groupTypes {
			if gt == "DynamicMembership" {
				diagnostics.AddAttributeError(
					path.Root("is_assignable_to_role"),
					"Invalid Role-Assignable Group Configuration",
					"When 'is_assignable_to_role' is true, the group cannot be a dynamic group. "+
						"Remove 'DynamicMembership' from 'group_types'.",
				)
				break
			}
		}
	}
}

// ValidateGroupConfiguration is a helper function that can be called during plan validation
func ValidateGroupConfiguration(ctx context.Context, data *GroupResourceModel) diag.Diagnostics {
	var diagnostics diag.Diagnostics

	tflog.Debug(ctx, "Validating group configuration")

	// Validate visibility constraints
	validateVisibility(ctx, data, &diagnostics)

	// Validate role assignability constraints
	validateRoleAssignability(ctx, data, &diagnostics)

	tflog.Debug(ctx, "Finished validating group configuration")

	return diagnostics
}
