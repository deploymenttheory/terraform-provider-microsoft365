package graphBetaEndpointSecurityIntent

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/devicemanagement"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	sharedmodels "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/shared_models/graph_beta/device_management"
)

// constructResource builds the createInstance request body for a new intent.
// The configured settings are sent as the settings delta, so Graph applies them on top of the template defaults in the same call.
// Unset collections are omitted rather than sent as null.
func constructResource(ctx context.Context, data *EndpointSecurityIntentResourceModel) (devicemanagement.TemplatesItemCreateInstancePostRequestBodyable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s create request from model", ResourceName))

	requestBody := devicemanagement.NewTemplatesItemCreateInstancePostRequestBody()

	convert.FrameworkToGraphString(data.DisplayName, requestBody.SetDisplayName)
	convert.FrameworkToGraphString(data.Description, requestBody.SetDescription)

	if err := convert.FrameworkToGraphStringSetOmitUnset(ctx, data.RoleScopeTagIds, requestBody.SetRoleScopeTagIds); err != nil {
		return nil, fmt.Errorf("failed to set role scope tag IDs: %w", err)
	}

	settings, err := constructSettings(ctx, data.Settings)
	if err != nil {
		return nil, err
	}
	if len(settings) > 0 {
		requestBody.SetSettingsDelta(settings)
	}

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]any{
			"error": err.Error(),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing %s create request", ResourceName))

	return requestBody, nil
}

// constructUpdateResource builds the PATCH request body for the intent's own properties.
// Settings and assignments are not part of the PATCH body; they are written with updateSettings and assign.
func constructUpdateResource(ctx context.Context, data *EndpointSecurityIntentResourceModel) (graphmodels.DeviceManagementIntentable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s update request from model", ResourceName))

	requestBody := graphmodels.NewDeviceManagementIntent()

	convert.FrameworkToGraphString(data.DisplayName, requestBody.SetDisplayName)
	convert.FrameworkToGraphString(data.Description, requestBody.SetDescription)

	if err := convert.FrameworkToGraphStringSetOmitUnset(ctx, data.RoleScopeTagIds, requestBody.SetRoleScopeTagIds); err != nil {
		return nil, fmt.Errorf("failed to set role scope tag IDs: %w", err)
	}

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]any{
			"error": err.Error(),
		})
	}

	return requestBody, nil
}

// constructUpdateSettings builds the updateSettings request body from the configured settings.
func constructUpdateSettings(ctx context.Context, data *EndpointSecurityIntentResourceModel) (devicemanagement.IntentsItemUpdateSettingsPostRequestBodyable, error) {
	settings, err := constructSettings(ctx, data.Settings)
	if err != nil {
		return nil, err
	}

	requestBody := devicemanagement.NewIntentsItemUpdateSettingsPostRequestBody()
	requestBody.SetSettings(settings)

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final settings JSON to be sent to Graph API for resource %s", ResourceName), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]any{
			"error": err.Error(),
		})
	}

	return requestBody, nil
}

// constructSettings converts the settings set into typed Graph setting instances carrying only the definition ID and value JSON.
func constructSettings(ctx context.Context, settingsSet types.Set) ([]graphmodels.DeviceManagementSettingInstanceable, error) {
	settings := make([]graphmodels.DeviceManagementSettingInstanceable, 0)
	if settingsSet.IsNull() || settingsSet.IsUnknown() {
		return settings, nil
	}

	var values []EndpointSecurityIntentSettingModel
	if diags := settingsSet.ElementsAs(ctx, &values, false); diags.HasError() {
		return nil, fmt.Errorf("%w: %v", errInvalidSettings, diags.Errors())
	}

	for _, value := range values {
		instance, err := newSettingInstance(value.SettingType.ValueString())
		if err != nil {
			return nil, fmt.Errorf("setting %s: %w", value.DefinitionId.ValueString(), err)
		}
		convert.FrameworkToGraphString(value.DefinitionId, instance.SetDefinitionId)
		convert.FrameworkToGraphString(value.ValueJson, instance.SetValueJson)
		settings = append(settings, instance)
	}

	return settings, nil
}

// constructAssignments builds the assign request body. An empty list removes every assignment, so it is always sent.
func constructAssignments(ctx context.Context, data *EndpointSecurityIntentResourceModel) (devicemanagement.IntentsItemAssignPostRequestBodyable, error) {
	assignments := make([]graphmodels.DeviceManagementIntentAssignmentable, 0)

	if !data.Assignments.IsNull() && !data.Assignments.IsUnknown() {
		var values []sharedmodels.DeviceManagementDeviceConfigurationAssignmentWithGroupFilterModel
		if diags := data.Assignments.ElementsAs(ctx, &values, false); diags.HasError() {
			return nil, fmt.Errorf("%w: %v", errInvalidAssignments, diags.Errors())
		}

		for _, value := range values {
			target, err := constructAssignmentTarget(value)
			if err != nil {
				return nil, err
			}
			assignment := graphmodels.NewDeviceManagementIntentAssignment()
			assignment.SetTarget(target)
			assignments = append(assignments, assignment)
		}
	}

	requestBody := devicemanagement.NewIntentsItemAssignPostRequestBody()
	requestBody.SetAssignments(assignments)

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final assignment JSON to be sent to Graph API for resource %s", ResourceName), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]any{
			"error": err.Error(),
		})
	}

	return requestBody, nil
}

// constructAssignmentTarget converts one assignment into its Graph target, including the assignment filter when one is set.
func constructAssignmentTarget(
	value sharedmodels.DeviceManagementDeviceConfigurationAssignmentWithGroupFilterModel,
) (graphmodels.DeviceAndAppManagementAssignmentTargetable, error) {
	var target graphmodels.DeviceAndAppManagementAssignmentTargetable

	switch value.Type.ValueString() {
	case "allDevicesAssignmentTarget":
		target = graphmodels.NewAllDevicesAssignmentTarget()
	case "allLicensedUsersAssignmentTarget":
		target = graphmodels.NewAllLicensedUsersAssignmentTarget()
	case "groupAssignmentTarget":
		group := graphmodels.NewGroupAssignmentTarget()
		convert.FrameworkToGraphString(value.GroupId, group.SetGroupId)
		target = group
	case "exclusionGroupAssignmentTarget":
		group := graphmodels.NewExclusionGroupAssignmentTarget()
		convert.FrameworkToGraphString(value.GroupId, group.SetGroupId)
		target = group
	default:
		return nil, fmt.Errorf("%w: unsupported assignment target type %q", errInvalidAssignments, value.Type.ValueString())
	}

	filterType := value.FilterType.ValueString()
	if filterType == "include" || filterType == "exclude" {
		convert.FrameworkToGraphString(value.FilterId, target.SetDeviceAndAppManagementAssignmentFilterId)
		setFilterType := target.SetDeviceAndAppManagementAssignmentFilterType
		if err := convert.FrameworkToGraphEnum(value.FilterType, graphmodels.ParseDeviceAndAppManagementAssignmentFilterType, setFilterType); err != nil {
			return nil, fmt.Errorf("failed to set assignment filter type: %w", err)
		}
	}

	return target, nil
}
