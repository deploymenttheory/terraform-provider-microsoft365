package graphBetaRoleDefinitionAssignment

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ModifyPlan handles plan modification for the RoleAssignment resource.
func (r *RoleAssignmentResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	// Only proceed if we have a plan (not during destroy)
	if req.Plan.Raw.IsNull() {
		return
	}

	var plan RoleAssignmentResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Validate scope configuration
	if len(plan.ScopeConfig) == 0 {
		resp.Diagnostics.AddError(
			"Missing Scope Configuration",
			"Exactly one scope_configuration block is required.",
		)
		return
	}

	if len(plan.ScopeConfig) > 1 {
		resp.Diagnostics.AddError(
			"Multiple Scope Configurations",
			"Only one scope_configuration block is allowed.",
		)
		return
	}

	scopeConfig := plan.ScopeConfig[0]
	scopeType := scopeConfig.Type.ValueString()

	// Validate resource_scopes based on type
	switch scopeType {
	case "ResourceScopes":
		if scopeConfig.ResourceScopes.IsNull() || len(scopeConfig.ResourceScopes.Elements()) == 0 {
			resp.Diagnostics.AddError(
				"Missing Resource Scopes",
				"resource_scopes is required when type is 'ResourceScopes' and must contain at least one scope ID.",
			)
			return
		}
	case "AllLicensedUsers", "AllDevices":
		if !scopeConfig.ResourceScopes.IsNull() && len(scopeConfig.ResourceScopes.Elements()) > 0 {
			resp.Diagnostics.AddError(
				"Invalid Resource Scopes",
				fmt.Sprintf("resource_scopes must be empty when type is '%s'.", scopeType),
			)
			return
		}
		// Set to null if not already to ensure consistency
		if !scopeConfig.ResourceScopes.IsNull() {
			plan.ScopeConfig[0].ResourceScopes = types.SetNull(types.StringType)
			resp.Diagnostics.Append(resp.Plan.Set(ctx, &plan)...)
		}
	default:
		resp.Diagnostics.AddError(
			"Invalid Scope Type",
			fmt.Sprintf("Invalid scope type '%s'. Valid values are: ResourceScopes, AllLicensedUsers, AllDevices.", scopeType),
		)
		return
	}
}
