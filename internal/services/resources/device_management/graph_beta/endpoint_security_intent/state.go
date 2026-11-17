package graphBetaEndpointSecurityIntent

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
)

// MapRemoteResourceStateToTerraform maps the remote intent, its settings and its assignments from Kiota SDK to Terraform state
func MapRemoteResourceStateToTerraform(
	ctx context.Context,
	data EndpointSecurityIntentResourceModel,
	remoteResource graphmodels.DeviceManagementIntentable,
	remoteSettings []graphmodels.DeviceManagementSettingInstanceable,
	remoteAssignments []graphmodels.DeviceManagementIntentAssignmentable,
) EndpointSecurityIntentResourceModel {
	if remoteResource == nil {
		tflog.Debug(ctx, "Remote resource is nil")
		return data
	}

	tflog.Debug(ctx, "Starting to map remote state to Terraform state", map[string]any{
		"resourceId": remoteResource.GetId(),
	})

	data.ID = convert.GraphToFrameworkString(remoteResource.GetId())
	data.TemplateId = convert.GraphToFrameworkString(remoteResource.GetTemplateId())
	data.DisplayName = convert.GraphToFrameworkString(remoteResource.GetDisplayName())
	data.Description = convert.GraphToFrameworkString(remoteResource.GetDescription())
	if data.Description.IsNull() {
		data.Description = types.StringValue("")
	}
	data.RoleScopeTagIds = convert.GraphToFrameworkStringSet(ctx, remoteResource.GetRoleScopeTagIds())
	data.IsAssigned = convert.GraphToFrameworkBool(remoteResource.GetIsAssigned())
	data.IsMigratingToConfigurationPolicy = convert.GraphToFrameworkBool(remoteResource.GetIsMigratingToConfigurationPolicy())
	data.LastModifiedDateTime = convert.GraphToFrameworkTime(remoteResource.GetLastModifiedDateTime())
	data.Settings = mapSettings(ctx, data.Settings, remoteSettings)
	data.Assignments = mapAssignments(ctx, data.Assignments, remoteAssignments)

	tflog.Debug(ctx, "Finished mapping remote state to Terraform state")

	return data
}

// mapSettings maps only the settings already tracked in state, because a template returns every one of its settings
// and the resource manages only the ones in configuration. A tracked value that differs from Graph only in JSON formatting keeps its state form.
// Tracked settings that the intent no longer has are dropped so the next plan adds them back.
func mapSettings(ctx context.Context, tracked types.Set, remoteSettings []graphmodels.DeviceManagementSettingInstanceable) types.Set {
	settingType := types.ObjectType{AttrTypes: settingAttrTypes}
	if tracked.IsNull() || tracked.IsUnknown() {
		return types.SetNull(settingType)
	}

	var trackedSettings []EndpointSecurityIntentSettingModel
	if diags := tracked.ElementsAs(ctx, &trackedSettings, false); diags.HasError() {
		tflog.Error(ctx, "Failed to read tracked settings from state", map[string]any{"errors": fmt.Sprintf("%v", diags.Errors())})
		return tracked
	}

	remoteByDefinition := make(map[string]graphmodels.DeviceManagementSettingInstanceable, len(remoteSettings))
	for _, remote := range remoteSettings {
		if remote.GetDefinitionId() != nil {
			remoteByDefinition[*remote.GetDefinitionId()] = remote
		}
	}

	elements := make([]attr.Value, 0, len(trackedSettings))
	for _, setting := range trackedSettings {
		remote, found := remoteByDefinition[setting.DefinitionId.ValueString()]
		if !found {
			tflog.Warn(ctx, "Tracked setting not returned by Graph", map[string]any{"definitionId": setting.DefinitionId.ValueString()})
			continue
		}

		remoteType := settingTypeOf(remote)
		if remoteType == "" {
			remoteType = setting.SettingType.ValueString()
		}

		remoteValue := "null"
		if remote.GetValueJson() != nil {
			remoteValue = *remote.GetValueJson()
		}
		if jsonEqual(setting.ValueJson.ValueString(), remoteValue) {
			remoteValue = setting.ValueJson.ValueString()
		}

		elements = append(elements, types.ObjectValueMust(settingAttrTypes, map[string]attr.Value{
			"definition_id": types.StringValue(setting.DefinitionId.ValueString()),
			"setting_type":  types.StringValue(remoteType),
			"value_json":    types.StringValue(remoteValue),
		}))
	}

	return types.SetValueMust(settingType, elements)
}

// mapAssignments maps the intent assignments. No assignments map to null when state has none, so an omitted block shows no drift.
func mapAssignments(ctx context.Context, current types.Set, remoteAssignments []graphmodels.DeviceManagementIntentAssignmentable) types.Set {
	assignmentType := types.ObjectType{AttrTypes: assignmentAttrTypes}

	elements := make([]attr.Value, 0, len(remoteAssignments))
	for _, assignment := range remoteAssignments {
		value, ok := mapAssignmentTarget(assignment.GetTarget())
		if !ok {
			tflog.Warn(ctx, "Skipping unsupported assignment target", map[string]any{"assignmentId": assignment.GetId()})
			continue
		}
		elements = append(elements, value)
	}

	if len(elements) == 0 && current.IsNull() {
		return types.SetNull(assignmentType)
	}

	return types.SetValueMust(assignmentType, elements)
}

// mapAssignmentTarget maps one assignment target. Targets without a filter map to the schema defaults for filter_id and filter_type.
func mapAssignmentTarget(target graphmodels.DeviceAndAppManagementAssignmentTargetable) (attr.Value, bool) {
	values := map[string]attr.Value{
		"group_id":    types.StringNull(),
		"filter_id":   types.StringValue(noFilterId),
		"filter_type": types.StringValue("none"),
	}

	switch typed := target.(type) {
	case *graphmodels.AllDevicesAssignmentTarget:
		values["type"] = types.StringValue("allDevicesAssignmentTarget")
	case *graphmodels.AllLicensedUsersAssignmentTarget:
		values["type"] = types.StringValue("allLicensedUsersAssignmentTarget")
	case *graphmodels.ExclusionGroupAssignmentTarget:
		values["type"] = types.StringValue("exclusionGroupAssignmentTarget")
		values["group_id"] = convert.GraphToFrameworkString(typed.GetGroupId())
	case *graphmodels.GroupAssignmentTarget:
		values["type"] = types.StringValue("groupAssignmentTarget")
		values["group_id"] = convert.GraphToFrameworkString(typed.GetGroupId())
	default:
		return nil, false
	}

	if filterId := target.GetDeviceAndAppManagementAssignmentFilterId(); filterId != nil && *filterId != "" {
		values["filter_id"] = types.StringValue(*filterId)
	}
	if filterType := target.GetDeviceAndAppManagementAssignmentFilterType(); filterType != nil {
		values["filter_type"] = types.StringValue(filterType.String())
	}

	return types.ObjectValueMust(assignmentAttrTypes, values), true
}
