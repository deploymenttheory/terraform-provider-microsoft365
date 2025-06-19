package graphBetaAppleUserInitiatedEnrollmentProfileAssignment

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// MapRemoteStateToTerraform maps a remote Apple User Initiated Enrollment Profile Assignment to the Terraform resource model
func MapRemoteStateToTerraform(ctx context.Context, data AppleUserInitiatedEnrollmentProfileAssignmentResourceModel, assignment graphmodels.AppleEnrollmentProfileAssignmentable) AppleUserInitiatedEnrollmentProfileAssignmentResourceModel {
	if assignment == nil {
		tflog.Debug(ctx, "Remote assignment is nil")
		return data
	}

	data.ID = convert.GraphToFrameworkString(assignment.GetId())

	if target := assignment.GetTarget(); target != nil {
		data.Target = mapRemoteTargetToTerraform(target)
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished mapping resource %s with id %s", ResourceName, data.ID.ValueString()))

	return data
}

// mapRemoteTargetToTerraform maps a remote assignment target to a Terraform assignment target
func mapRemoteTargetToTerraform(target graphmodels.DeviceAndAppManagementAssignmentTargetable) AssignmentTargetResourceModel {
	targetModel := AssignmentTargetResourceModel{}

	if target == nil {
		return targetModel
	}

	if odataType := target.GetOdataType(); odataType != nil {
		targetModel.TargetType = types.StringValue(getTargetTypeFromOdataType(*odataType))
	}

	targetModel.DeviceAndAppManagementAssignmentFilterId = convert.GraphToFrameworkString(target.GetDeviceAndAppManagementAssignmentFilterId())
	targetModel.DeviceAndAppManagementAssignmentFilterType = convert.GraphToFrameworkEnum(target.GetDeviceAndAppManagementAssignmentFilterType())

	// Map target-specific properties
	switch typedTarget := target.(type) {
	case graphmodels.GroupAssignmentTargetable:
		groupId := typedTarget.GetGroupId()
		if groupId != nil {
			targetModel.GroupId = convert.GraphToFrameworkString(groupId)
			targetModel.EntraObjectId = convert.GraphToFrameworkString(groupId)
		}
	case graphmodels.ExclusionGroupAssignmentTargetable:
		targetModel.GroupId = convert.GraphToFrameworkString(typedTarget.GetGroupId())
	}

	return targetModel
}

// getTargetTypeFromOdataType converts OData type to target type string
func getTargetTypeFromOdataType(odataType string) string {
	switch odataType {
	case "#microsoft.graph.allLicensedUsersAssignmentTarget":
		return "allUsers"
	case "#microsoft.graph.groupAssignmentTarget":
		return "group"
	case "#microsoft.graph.exclusionGroupAssignmentTarget":
		return "exclusionGroup"
	default:
		return ""
	}
}
