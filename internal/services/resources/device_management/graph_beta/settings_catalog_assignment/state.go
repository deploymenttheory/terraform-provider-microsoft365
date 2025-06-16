package graphBetaDeviceManagementConfigurationPolicyAssignment

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/state"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// MapRemoteStateToTerraform maps a remote assignment to the Terraform resource model
func MapRemoteStateToTerraform(ctx context.Context, data DeviceManagementConfigurationPolicyAssignmentResourceModel, assignment graphmodels.DeviceManagementConfigurationPolicyAssignmentable) DeviceManagementConfigurationPolicyAssignmentResourceModel {
	if assignment == nil {
		tflog.Debug(ctx, "Remote assignment is nil")
		return data
	}

	data.ID = state.StringPointerValue(assignment.GetId())
	data.Source = state.EnumPtrToTypeString(assignment.GetSource())
	data.SourceId = state.StringPointerValue(assignment.GetSourceId())

	if target := assignment.GetTarget(); target != nil {
		data.Target = mapRemoteTargetToTerraform(target)
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished stating resource %s with id %s", ResourceName, data.ID.ValueString()))

	return data
}

// mapRemoteTargetToTerraform maps a remote assignment target to a Terraform assignment target
func mapRemoteTargetToTerraform(remoteTarget graphmodels.DeviceAndAppManagementAssignmentTargetable) AssignmentTargetResourceModel {
	target := AssignmentTargetResourceModel{
		DeviceAndAppManagementAssignmentFilterId:   types.StringPointerValue(remoteTarget.GetDeviceAndAppManagementAssignmentFilterId()),
		DeviceAndAppManagementAssignmentFilterType: state.EnumPtrToTypeString(remoteTarget.GetDeviceAndAppManagementAssignmentFilterType()),
	}

	switch v := remoteTarget.(type) {
	case *graphmodels.GroupAssignmentTarget:
		target.TargetType = types.StringValue("groupAssignment")
		target.GroupId = types.StringPointerValue(v.GetGroupId())
	case *graphmodels.ExclusionGroupAssignmentTarget:
		target.TargetType = types.StringValue("exclusionGroupAssignment")
		target.GroupId = types.StringPointerValue(v.GetGroupId())
	case *graphmodels.ConfigurationManagerCollectionAssignmentTarget:
		target.TargetType = types.StringValue("configurationManagerCollection")
		target.CollectionId = types.StringPointerValue(v.GetCollectionId())
	case *graphmodels.AllDevicesAssignmentTarget:
		target.TargetType = types.StringValue("allDevices")
	case *graphmodels.AllLicensedUsersAssignmentTarget:
		target.TargetType = types.StringValue("allLicensedUsers")
	}

	return target
}
