package graphBetaWindowsBackupAndRestore

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// MapRemoteStateToTerraform maps the remote GraphServiceClient object to a Terraform state.
func MapRemoteStateToTerraform(ctx context.Context, data *WindowsBackupAndRestoreResourceModel, remoteResource graphmodels.DeviceEnrollmentConfigurationable) {
	if remoteResource == nil {
		tflog.Debug(ctx, "Remote resource is nil")
		return
	}

	tflog.Debug(ctx, "Starting to map remote state to Terraform state", map[string]any{
		"resourceId": convert.GraphToFrameworkString(remoteResource.GetId()).ValueString(),
	})

	data.ID = convert.GraphToFrameworkString(remoteResource.GetId())
	data.DisplayName = convert.GraphToFrameworkString(remoteResource.GetDisplayName())
	data.Description = convert.GraphToFrameworkString(remoteResource.GetDescription())
	data.Priority = convert.GraphToFrameworkInt32(remoteResource.GetPriority())
	data.CreatedDateTime = convert.GraphToFrameworkTime(remoteResource.GetCreatedDateTime())
	data.LastModifiedDateTime = convert.GraphToFrameworkTime(remoteResource.GetLastModifiedDateTime())
	data.Version = convert.GraphToFrameworkInt32(remoteResource.GetVersion())
	data.RoleScopeTagIds = convert.GraphToFrameworkStringSet(ctx, remoteResource.GetRoleScopeTagIds())
	data.DeviceEnrollmentConfigurationType = convert.GraphToFrameworkEnum(remoteResource.GetDeviceEnrollmentConfigurationType())

	// This resource only handles Windows Restore device enrollment configurations
	if windowsRestoreConfig, ok := remoteResource.(*graphmodels.WindowsRestoreDeviceEnrollmentConfiguration); ok {
		mapWindowsRestoreConfigurationToState(ctx, data, windowsRestoreConfig)
	} else {
		tflog.Error(ctx, "Remote resource is not a Windows Restore device enrollment configuration")
		return
	}

	assignments := remoteResource.GetAssignments()
	tflog.Debug(ctx, "Retrieved assignments from remote resource", map[string]any{
		"assignmentCount": len(assignments),
		"resourceId":      data.ID.ValueString(),
	})

	if len(assignments) == 0 {
		tflog.Debug(ctx, "No assignments found, setting assignments to null", map[string]any{
			"resourceId": data.ID.ValueString(),
		})
		data.Assignments = types.SetNull(WindowsBackupAndRestoreAssignmentType())
	} else {
		tflog.Debug(ctx, "Starting assignment mapping process", map[string]any{
			"resourceId":      data.ID.ValueString(),
			"assignmentCount": len(assignments),
		})
		MapAssignmentsToTerraform(ctx, data, assignments)
		tflog.Debug(ctx, "Completed assignment mapping process", map[string]any{
			"resourceId": data.ID.ValueString(),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished mapping resource %s with id %s", ResourceName, data.ID.ValueString()))
}

// mapWindowsRestoreConfigurationToState maps Windows Restore specific settings using SDK getters.
func mapWindowsRestoreConfigurationToState(ctx context.Context, data *WindowsBackupAndRestoreResourceModel, config *graphmodels.WindowsRestoreDeviceEnrollmentConfiguration) {
	data.State = convert.GraphToFrameworkEnum(config.GetState())
}

// WindowsBackupAndRestoreAssignmentType returns the type for assignments
func WindowsBackupAndRestoreAssignmentType() types.ObjectType {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"target": types.ObjectType{
				AttrTypes: map[string]attr.Type{
					"device_and_app_management_assignment_filter": types.ObjectType{
						AttrTypes: map[string]attr.Type{
							"filter_id":   types.StringType,
							"filter_type": types.StringType,
						},
					},
					"device_and_app_management_assignment_filter_id": types.StringType,
					"group_id": types.StringType,
					"intent":   types.StringType,
					"target":   types.StringType,
				},
			},
		},
	}
}
