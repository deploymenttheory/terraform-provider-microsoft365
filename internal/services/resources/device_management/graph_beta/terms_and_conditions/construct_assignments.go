package graphBetaTermsAndConditions

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"

	sharedmodels "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/shared_models/graph_beta/device_management"
)

// constructAssignment constructs and returns a TermsAndConditionsAssignment for individual assignment creation
func constructAssignment(ctx context.Context, assignment sharedmodels.DeviceManagementDeviceConfigurationAssignmentWithAllLicensedUsersInclusionGroupConfigurationManagerCollectionAssignmentModel) (graphmodels.TermsAndConditionsAssignmentable, error) {
	tflog.Debug(ctx, "Starting terms and conditions assignment construction")

	// Create TermsAndConditionsAssignment
	termsAndConditionsAssignment := graphmodels.NewTermsAndConditionsAssignment()

	if assignment.Type.IsNull() || assignment.Type.IsUnknown() {
		return nil, fmt.Errorf("assignment target type is missing or invalid")
	}

	targetType := assignment.Type.ValueString()

	target := constructTarget(ctx, targetType, assignment)
	if target == nil {
		return nil, fmt.Errorf("failed to create target for type: %s", targetType)
	}

	termsAndConditionsAssignment.SetTarget(target)

	if err := constructors.DebugLogGraphObject(ctx, "Constructed assignment", termsAndConditionsAssignment); err != nil {
		tflog.Error(ctx, "Failed to debug log assignment", map[string]any{
			"error": err.Error(),
		})
	}

	return termsAndConditionsAssignment, nil
}

// constructTarget creates the appropriate target based on the target type
func constructTarget(ctx context.Context, targetType string, assignment sharedmodels.DeviceManagementDeviceConfigurationAssignmentWithAllLicensedUsersInclusionGroupConfigurationManagerCollectionAssignmentModel) graphmodels.DeviceAndAppManagementAssignmentTargetable {
	var target graphmodels.DeviceAndAppManagementAssignmentTargetable

	switch targetType {
	case "allLicensedUsersAssignmentTarget":
		target = graphmodels.NewAllLicensedUsersAssignmentTarget()
	case "allDevicesAssignmentTarget":
		target = graphmodels.NewAllDevicesAssignmentTarget()
	case "groupAssignmentTarget":
		groupTarget := graphmodels.NewGroupAssignmentTarget()
		if !assignment.GroupId.IsNull() && !assignment.GroupId.IsUnknown() && assignment.GroupId.ValueString() != "" {
			convert.FrameworkToGraphString(assignment.GroupId, groupTarget.SetGroupId)
		} else {
			tflog.Error(ctx, "Group assignment target missing required group_id", map[string]any{
				"targetType": targetType,
			})
			return nil
		}
		target = groupTarget
	case "exclusionGroupAssignmentTarget":
		exclusionTarget := graphmodels.NewExclusionGroupAssignmentTarget()
		if !assignment.GroupId.IsNull() && !assignment.GroupId.IsUnknown() && assignment.GroupId.ValueString() != "" {
			convert.FrameworkToGraphString(assignment.GroupId, exclusionTarget.SetGroupId)
		} else {
			tflog.Error(ctx, "Exclusion group assignment target missing required group_id", map[string]any{
				"targetType": targetType,
			})
			return nil
		}
		target = exclusionTarget
	case "configurationManagerCollectionAssignmentTarget":
		configMgrTarget := graphmodels.NewConfigurationManagerCollectionAssignmentTarget()
		if !assignment.CollectionId.IsNull() && !assignment.CollectionId.IsUnknown() && assignment.CollectionId.ValueString() != "" {
			convert.FrameworkToGraphString(assignment.CollectionId, configMgrTarget.SetCollectionId)
		} else {
			tflog.Error(ctx, "Configuration manager collection assignment target missing required collection_id", map[string]any{
				"targetType": targetType,
			})
			return nil
		}
		target = configMgrTarget
	default:
		tflog.Error(ctx, "Unsupported target type", map[string]any{
			"targetType": targetType,
		})
		return nil
	}

	return target
}
