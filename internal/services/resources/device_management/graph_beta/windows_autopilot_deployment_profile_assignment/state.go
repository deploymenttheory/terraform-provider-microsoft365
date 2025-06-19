package graphBetaWindowsAutopilotDeploymentProfileAssignment

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// MapRemoteStateToTerraform maps a remote Windows Autopilot Deployment Profile Assignment to the Terraform resource model
func MapRemoteStateToTerraform(ctx context.Context, data WindowsAutopilotDeploymentProfileAssignmentResourceModel, assignment graphmodels.WindowsAutopilotDeploymentProfileAssignmentable) WindowsAutopilotDeploymentProfileAssignmentResourceModel {
	if assignment == nil {
		tflog.Debug(ctx, "Remote assignment is nil")
		return data
	}

	data.ID = convert.GraphToFrameworkString(assignment.GetId())
	data.Source = convert.GraphToFrameworkEnum(assignment.GetSource())

	if !data.SourceId.IsNull() && !data.SourceId.IsUnknown() {
		data.SourceId = convert.GraphToFrameworkString(assignment.GetSourceId())
	}

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

	// Map target type based on OData type
	if odataType := target.GetOdataType(); odataType != nil {
		targetModel.TargetType = types.StringValue(getTargetTypeFromOdataType(*odataType))
	}

	// Map common properties
	targetModel.DeviceAndAppManagementAssignmentFilterId = convert.GraphToFrameworkString(target.GetDeviceAndAppManagementAssignmentFilterId())
	targetModel.DeviceAndAppManagementAssignmentFilterType = convert.GraphToFrameworkEnum(target.GetDeviceAndAppManagementAssignmentFilterType())

	// Map group-specific properties
	switch typedTarget := target.(type) {
	case graphmodels.GroupAssignmentTargetable:
		targetModel.GroupId = convert.GraphToFrameworkString(typedTarget.GetGroupId())
	case graphmodels.ExclusionGroupAssignmentTargetable:
		targetModel.GroupId = convert.GraphToFrameworkString(typedTarget.GetGroupId())
	}

	return targetModel
}
