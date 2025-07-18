// MapRemoteResourceStateToTerraform states the base properties of a SettingsCatalogProfileResourceModel to a Terraform state
package graphBetaMacOSPlatformScript

import (
	"context"
	"fmt"
	"sort"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	sharedmodels "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/shared_models/graph_beta/device_management"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// MapRemoteResourceStateToTerraform maps the base properties of a MacOSPlatformScriptResourceModel to a Terraform state.
func MapRemoteResourceStateToTerraform(ctx context.Context, data *MacOSPlatformScriptResourceModel, remoteResource graphmodels.DeviceShellScriptable) {
	if remoteResource == nil {
		tflog.Debug(ctx, "Remote resource is nil")
		return
	}

	tflog.Debug(ctx, "Starting to map remote state to Terraform state", map[string]interface{}{
		"resourceName": remoteResource.GetDisplayName(),
		"resourceId":   remoteResource.GetId(),
	})

	data.ID = convert.GraphToFrameworkString(remoteResource.GetId())
	data.DisplayName = convert.GraphToFrameworkString(remoteResource.GetDisplayName())
	data.Description = convert.GraphToFrameworkString(remoteResource.GetDescription())
	data.RunAsAccount = convert.GraphToFrameworkEnum(remoteResource.GetRunAsAccount())
	data.FileName = convert.GraphToFrameworkString(remoteResource.GetFileName())
	data.RoleScopeTagIds = convert.GraphToFrameworkStringSet(ctx, remoteResource.GetRoleScopeTagIds())
	data.BlockExecutionNotifications = convert.GraphToFrameworkBool(remoteResource.GetBlockExecutionNotifications())
	data.ExecutionFrequency = convert.GraphToFrameworkISODuration(remoteResource.GetExecutionFrequency())
	data.ScriptContent = convert.GraphToFrameworkBytes(remoteResource.GetScriptContent())
	data.RetryCount = convert.GraphToFrameworkInt32(remoteResource.GetRetryCount())

	assignments := remoteResource.GetAssignments()

	// If there are no assignments, set data.Assignments to nil
	if len(assignments) == 0 {
		tflog.Debug(ctx, "No assignments found, setting assignments to nil")
		data.Assignments = nil
	} else {
		MapAssignmentsToTerraform(ctx, data, assignments)
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished stating resource %s with id %s", ResourceName, data.ID.ValueString()))
}

// MapAssignmentsToTerraform processes script assignments directly from the slice returned by GetAssignments
// There appears to be no other way to do this, as the assignments are not returned by any other api call
// despite all of the docs saying there is.
func MapAssignmentsToTerraform(ctx context.Context, data *MacOSPlatformScriptResourceModel, assignments []graphmodels.DeviceManagementScriptAssignmentable) {
	if len(assignments) == 0 {
		tflog.Debug(ctx, "No assignments to process, setting assignments to nil")
		data.Assignments = nil
		return
	}

	tflog.Debug(ctx, "Processing assignments from resource response")

	processAssignments(ctx, data, assignments)
}

// processAssignments handles the direct processing of assignment slices
// This contains the core logic from MapRemoteAssignmentStateToTerraform but works with the slice type
func processAssignments(ctx context.Context, data *MacOSPlatformScriptResourceModel, assignments []graphmodels.DeviceManagementScriptAssignmentable) {
	tflog.Debug(ctx, "Starting to map assignments directly to Terraform state")

	// If no assignments are provided, set data.Assignments to nil and return
	if len(assignments) == 0 {
		tflog.Debug(ctx, "No assignments to process in processAssignments, setting assignments to nil")
		data.Assignments = nil
		return
	}

	scriptAssignments := &sharedmodels.DeviceManagementScriptAssignmentResourceModel{
		AllDevices: types.BoolValue(false),
		AllUsers:   types.BoolValue(false),
	}

	var allDeviceAssignments []graphmodels.DeviceManagementScriptAssignmentable
	var allUserAssignments []graphmodels.DeviceManagementScriptAssignmentable
	var includeGroupAssignments []graphmodels.DeviceManagementScriptAssignmentable
	var excludeGroupAssignments []graphmodels.DeviceManagementScriptAssignmentable

	for _, assignment := range assignments {
		if target := assignment.GetTarget(); target != nil {
			if odataType := target.GetOdataType(); odataType != nil {
				switch *odataType {
				case "#microsoft.graph.allDevicesAssignmentTarget":
					allDeviceAssignments = append(allDeviceAssignments, assignment)
				case "#microsoft.graph.allLicensedUsersAssignmentTarget":
					allUserAssignments = append(allUserAssignments, assignment)
				case "#microsoft.graph.groupAssignmentTarget":
					includeGroupAssignments = append(includeGroupAssignments, assignment)
				case "#microsoft.graph.exclusionGroupAssignmentTarget":
					excludeGroupAssignments = append(excludeGroupAssignments, assignment)
				}
			}
		}
	}

	if len(allDeviceAssignments) > 0 {
		scriptAssignments.AllDevices = types.BoolValue(true)
	}

	if len(allUserAssignments) > 0 {
		scriptAssignments.AllUsers = types.BoolValue(true)
	}

	if len(includeGroupAssignments) > 0 {
		includeGroupIds := make([]types.String, 0)
		for _, assignment := range includeGroupAssignments {
			if target, ok := assignment.GetTarget().(graphmodels.GroupAssignmentTargetable); ok {
				if groupId := target.GetGroupId(); groupId != nil {
					includeGroupIds = append(includeGroupIds, types.StringValue(*groupId))
				}
			}
		}

		// Sort include group IDs alphanumerically
		sort.Slice(includeGroupIds, func(i, j int) bool {
			return includeGroupIds[i].ValueString() < includeGroupIds[j].ValueString()
		})

		scriptAssignments.IncludeGroupIds = includeGroupIds
	}

	if len(excludeGroupAssignments) > 0 {
		excludeGroupIds := make([]types.String, 0)
		for _, assignment := range excludeGroupAssignments {
			if target, ok := assignment.GetTarget().(graphmodels.GroupAssignmentTargetable); ok {
				if groupId := target.GetGroupId(); groupId != nil {
					excludeGroupIds = append(excludeGroupIds, types.StringValue(*groupId))
				}
			}
		}

		// Sort exclude group IDs alphanumerically
		sort.Slice(excludeGroupIds, func(i, j int) bool {
			return excludeGroupIds[i].ValueString() < excludeGroupIds[j].ValueString()
		})

		scriptAssignments.ExcludeGroupIds = excludeGroupIds
	}

	// Only set assignments if there are actual assignments
	if len(allDeviceAssignments) > 0 || len(allUserAssignments) > 0 ||
		len(includeGroupAssignments) > 0 || len(excludeGroupAssignments) > 0 {
		data.Assignments = scriptAssignments
	} else {
		data.Assignments = nil
	}

	tflog.Debug(ctx, "Finished mapping assignments directly to Terraform state")
}
