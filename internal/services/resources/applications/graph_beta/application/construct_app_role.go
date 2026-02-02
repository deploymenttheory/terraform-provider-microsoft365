package graphBetaApplication

import (
	"context"
	"fmt"
	"slices"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// ConstructAppRoleIsEnabledToFalse builds app roles with specified roles disabled
// This is used as the first step when removing app roles - they must be disabled before deletion
func ConstructAppRoleIsEnabledToFalse(ctx context.Context, currentRoles types.Set, roleIdsToDisable []string) ([]graphmodels.AppRoleable, error) {
	var stateRoles []ApplicationAppRole
	diags := currentRoles.ElementsAs(ctx, &stateRoles, false)
	if diags.HasError() {
		return nil, fmt.Errorf("failed to extract current app_roles: %s", diags.Errors()[0].Detail())
	}

	tflog.Debug(ctx, fmt.Sprintf("Disabling %d app roles out of %d total", len(roleIdsToDisable), len(stateRoles)))

	result := make([]graphmodels.AppRoleable, 0, len(stateRoles))
	for _, stateRole := range stateRoles {
		appRole := graphmodels.NewAppRole()

		// Set all fields from current state
		if err := convert.FrameworkToGraphUUID(stateRole.Id, appRole.SetId); err != nil {
			return nil, fmt.Errorf("failed to set app role id: %w", err)
		}
		convert.FrameworkToGraphString(stateRole.Description, appRole.SetDescription)
		convert.FrameworkToGraphString(stateRole.DisplayName, appRole.SetDisplayName)
		convert.FrameworkToGraphString(stateRole.Value, appRole.SetValue)

		var memberTypes []string
		stateRole.AllowedMemberTypes.ElementsAs(ctx, &memberTypes, false)
		appRole.SetAllowedMemberTypes(memberTypes)

		// Disable if this role is in the disable list, otherwise keep current state
		if slices.Contains(roleIdsToDisable, stateRole.Id.ValueString()) {
			disabled := false
			appRole.SetIsEnabled(&disabled)
		} else {
			convert.FrameworkToGraphBool(stateRole.IsEnabled, appRole.SetIsEnabled)
		}

		result = append(result, appRole)
	}

	return result, nil
}

// ConstructAppRolesForUpdate builds the complete app roles collection for PATCH updates
// This handles both additions and removals in a single operation
func ConstructAppRolesForUpdate(ctx context.Context, data types.Set) ([]graphmodels.AppRoleable, error) {
	if data.IsNull() || data.IsUnknown() {
		return nil, nil
	}

	var appRoles []ApplicationAppRole
	diags := data.ElementsAs(ctx, &appRoles, false)
	if diags.HasError() {
		return nil, fmt.Errorf("failed to extract app_roles: %s", diags.Errors()[0].Detail())
	}

	tflog.Debug(ctx, fmt.Sprintf("Constructing %d app roles for update", len(appRoles)))

	result := make([]graphmodels.AppRoleable, 0, len(appRoles))
	for _, role := range appRoles {
		appRole := graphmodels.NewAppRole()

		// Required fields
		if err := convert.FrameworkToGraphUUID(role.Id, appRole.SetId); err != nil {
			return nil, fmt.Errorf("failed to set app role id: %w", err)
		}

		convert.FrameworkToGraphString(role.Description, appRole.SetDescription)
		convert.FrameworkToGraphString(role.DisplayName, appRole.SetDisplayName)
		convert.FrameworkToGraphBool(role.IsEnabled, appRole.SetIsEnabled)

		// Optional value field
		convert.FrameworkToGraphString(role.Value, appRole.SetValue)

		// AllowedMemberTypes - convert Set to []string
		if !role.AllowedMemberTypes.IsNull() && !role.AllowedMemberTypes.IsUnknown() {
			var memberTypes []string
			diags := role.AllowedMemberTypes.ElementsAs(ctx, &memberTypes, false)
			if diags.HasError() {
				return nil, fmt.Errorf("failed to extract allowed_member_types: %s", diags.Errors()[0].Detail())
			}
			appRole.SetAllowedMemberTypes(memberTypes)
		}

		result = append(result, appRole)
	}

	return result, nil
}
