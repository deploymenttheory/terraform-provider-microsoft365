package graphDeviceConfigurationAssignment

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoftgraph/msgraph-sdk-go/models"
)

func MapRemoteStateToTerraform(ctx context.Context, data *DeviceConfigurationAssignmentResourceModel, remoteResource models.DeviceConfigurationAssignmentable) {
	if remoteResource == nil {
		tflog.Debug(ctx, "Remote resource is nil")
		return
	}

	tflog.Debug(ctx, "Starting to map remote state to Terraform state", map[string]interface{}{
		"resourceId": convert.GraphToFrameworkString(remoteResource.GetId()).ValueString(),
	})

	// Map basic properties using convert helpers
	data.ID = convert.GraphToFrameworkString(remoteResource.GetId())

	// Map target information
	target := remoteResource.GetTarget()
	if target != nil {
		mapTargetToTerraform(ctx, data, target)
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished mapping remote state for resource %s with id %s", ResourceName, data.ID.ValueString()))
}

// mapTargetToTerraform maps the target object to Terraform state - NO FILTERS
func mapTargetToTerraform(ctx context.Context, data *DeviceConfigurationAssignmentResourceModel, target models.DeviceAndAppManagementAssignmentTargetable) {
	// Clean type switch for target-specific properties only
	switch v := target.(type) {
	case *models.GroupAssignmentTarget:
		data.TargetType = types.StringValue("groupAssignment")
		data.GroupId = convert.GraphToFrameworkString(v.GetGroupId())

	case *models.AllDevicesAssignmentTarget:
		data.TargetType = types.StringValue("allDevices")
		data.GroupId = types.StringValue("")

	case *models.AllLicensedUsersAssignmentTarget:
		data.TargetType = types.StringValue("allLicensedUsers")
		data.GroupId = types.StringValue("")

	case *models.ExclusionGroupAssignmentTarget:
		data.TargetType = types.StringValue("exclusionGroupAssignment")
		data.GroupId = convert.GraphToFrameworkString(v.GetGroupId())

	case *models.ConfigurationManagerCollectionAssignmentTarget:
		data.TargetType = types.StringValue("configurationManagerCollection")
		data.GroupId = convert.GraphToFrameworkString(v.GetCollectionId())

	default:
		tflog.Warn(ctx, "Unknown target type", map[string]interface{}{
			"targetType": fmt.Sprintf("%T", target),
		})
		data.TargetType = types.StringNull()
		data.GroupId = types.StringNull()
	}
}
