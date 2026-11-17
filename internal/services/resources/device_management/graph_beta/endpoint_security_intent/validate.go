package graphBetaEndpointSecurityIntent

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"

	sharedmodels "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/shared_models/graph_beta/device_management"
)

// ValidateConfig rejects settings and assignments that Graph would reject or silently change after apply.
func (r *EndpointSecurityIntentResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var data EndpointSecurityIntentResourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	validateSettings(ctx, &data, resp)
	validateAssignments(ctx, &data, resp)
}

// validateSettings checks that each definition ID appears once and that each value is valid JSON.
func validateSettings(ctx context.Context, data *EndpointSecurityIntentResourceModel, resp *resource.ValidateConfigResponse) {
	if data.Settings.IsNull() || data.Settings.IsUnknown() {
		return
	}

	var settings []EndpointSecurityIntentSettingModel
	resp.Diagnostics.Append(data.Settings.ElementsAs(ctx, &settings, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	seen := make(map[string]bool, len(settings))
	for _, setting := range settings {
		if !setting.DefinitionId.IsUnknown() && !setting.DefinitionId.IsNull() {
			definitionId := setting.DefinitionId.ValueString()
			if seen[definitionId] {
				detail := fmt.Sprintf("Setting %q is configured more than once. Each definition_id may appear only once.", definitionId)
				resp.Diagnostics.AddAttributeError(path.Root("settings"), "Duplicate setting", detail)
			}
			seen[definitionId] = true
		}

		if !setting.ValueJson.IsUnknown() && !setting.ValueJson.IsNull() && !json.Valid([]byte(setting.ValueJson.ValueString())) {
			detail := fmt.Sprintf("value_json for setting %q is not valid JSON. Use jsonencode() to build it.", setting.DefinitionId.ValueString())
			resp.Diagnostics.AddAttributeError(path.Root("settings"), "Invalid setting value", detail)
		}
	}
}

// validateAssignments checks that group targets have a group ID, that other targets do not, and that filters are only used where Graph supports them.
func validateAssignments(ctx context.Context, data *EndpointSecurityIntentResourceModel, resp *resource.ValidateConfigResponse) {
	if data.Assignments.IsNull() || data.Assignments.IsUnknown() {
		return
	}

	var assignments []sharedmodels.DeviceManagementDeviceConfigurationAssignmentWithGroupFilterModel
	resp.Diagnostics.Append(data.Assignments.ElementsAs(ctx, &assignments, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	for _, assignment := range assignments {
		if assignment.Type.IsUnknown() || assignment.GroupId.IsUnknown() {
			continue
		}

		targetType := assignment.Type.ValueString()
		groupTarget := targetType == "groupAssignmentTarget" || targetType == "exclusionGroupAssignmentTarget"

		if groupTarget && assignment.GroupId.IsNull() {
			resp.Diagnostics.AddAttributeError(path.Root("assignments"), "Missing group ID", fmt.Sprintf("Assignments of type %q require group_id.", targetType))
		}
		if !groupTarget && !assignment.GroupId.IsNull() {
			detail := fmt.Sprintf("group_id only applies to group assignment targets, not %q.", targetType)
			resp.Diagnostics.AddAttributeError(path.Root("assignments"), "Unexpected group ID", detail)
		}

		filterType := assignment.FilterType.ValueString()
		if targetType == "exclusionGroupAssignmentTarget" && (filterType == "include" || filterType == "exclude") {
			resp.Diagnostics.AddAttributeError(path.Root("assignments"), "Unsupported assignment filter", "Exclusion group assignments do not support assignment filters.")
		}
	}
}
