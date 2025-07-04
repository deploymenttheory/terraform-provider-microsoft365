package graphBetaMacosPlatformScriptAssignment

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// MapRemoteStateToTerraform maps a remote assignment to the Terraform resource model
func MapRemoteStateToTerraform(ctx context.Context, data MacosPlatformScriptAssignmentResourceModel, assignment graphmodels.DeviceManagementScriptAssignmentable) MacosPlatformScriptAssignmentResourceModel {
	if assignment == nil {
		tflog.Debug(ctx, "Remote assignment is nil")
		return data
	}

	data.ID = convert.GraphToFrameworkString(assignment.GetId())

	if target := assignment.GetTarget(); target != nil {
		data.Target = mapRemoteTargetToTerraform(target)
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished stating resource %s with id %s", ResourceName, data.ID.ValueString()))

	return data
}

// mapRemoteTargetToTerraform maps a remote assignment target to a Terraform assignment target
func mapRemoteTargetToTerraform(remoteTarget graphmodels.DeviceAndAppManagementAssignmentTargetable) AssignmentTargetResourceModel {
	target := AssignmentTargetResourceModel{
		DeviceAndAppManagementAssignmentFilterId:   convert.GraphToFrameworkString(remoteTarget.GetDeviceAndAppManagementAssignmentFilterId()),
		DeviceAndAppManagementAssignmentFilterType: convert.GraphToFrameworkEnum(remoteTarget.GetDeviceAndAppManagementAssignmentFilterType()),
	}

	switch v := remoteTarget.(type) {
	case *graphmodels.GroupAssignmentTarget:
		target.TargetType = types.StringValue("groupAssignment")
		target.GroupId = convert.GraphToFrameworkString(v.GetGroupId())
	case *graphmodels.AllDevicesAssignmentTarget:
		target.TargetType = types.StringValue("allDevices")
	case *graphmodels.AllLicensedUsersAssignmentTarget:
		target.TargetType = types.StringValue("allLicensedUsers")
	}

	return target
}
