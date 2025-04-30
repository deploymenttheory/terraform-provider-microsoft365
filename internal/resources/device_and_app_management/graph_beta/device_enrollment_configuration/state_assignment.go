package graphBetaDeviceEnrollmentConfiguration

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// MapRemoteAssignmentsToTerraform maps the Graph API assignments into the Terraform consolidated model
func MapRemoteAssignmentsToTerraform(ctx context.Context, data *DeviceEnrollmentConfigurationResourceModel, assignments []graphmodels.EnrollmentConfigurationAssignmentable) {
	if assignments == nil || len(assignments) == 0 {
		tflog.Debug(ctx, "No assignments found, clearing assignments field")
		data.Assignments = nil
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Mapping %d assignments from Graph API", len(assignments)))

	includeGroupIDs := []string{}
	excludeGroupIDs := []string{}

	for _, assignment := range assignments {
		if assignment == nil || assignment.GetTarget() == nil {
			continue
		}

		target := assignment.GetTarget()

		if groupTarget, ok := target.(graphmodels.GroupAssignmentTargetable); ok && groupTarget != nil && groupTarget.GetGroupId() != nil {
			groupID := *groupTarget.GetGroupId()

			if odataType := groupTarget.GetOdataType(); odataType != nil {
				if strings.Contains(*odataType, "exclusion") || strings.Contains(*odataType, "ExclusionGroupAssignmentTarget") {
					excludeGroupIDs = append(excludeGroupIDs, groupID)
					tflog.Debug(ctx, fmt.Sprintf("Mapped exclusion group ID: %s", groupID))
				} else {
					includeGroupIDs = append(includeGroupIDs, groupID)
					tflog.Debug(ctx, fmt.Sprintf("Mapped inclusion group ID: %s", groupID))
				}
			} else {
				tflog.Error(ctx, fmt.Sprintf("Inclusion group ID returned an unsupported odataType: %s", groupID))
			}
		}
	}

	var assignmentBlocks []AssignmentResourceModel

	if len(includeGroupIDs) > 0 {
		includeSet, diags := types.SetValueFrom(ctx, types.StringType, includeGroupIDs)
		if diags.HasError() {
			tflog.Error(ctx, "Failed to create Set from inclusion group IDs", map[string]interface{}{
				"error": diags.Errors(),
			})
		} else {
			assignmentBlocks = append(assignmentBlocks, AssignmentResourceModel{
				Target:   types.StringValue("include"),
				GroupIds: includeSet,
			})
		}
	}

	if len(excludeGroupIDs) > 0 {
		excludeSet, diags := types.SetValueFrom(ctx, types.StringType, excludeGroupIDs)
		if diags.HasError() {
			tflog.Error(ctx, "Failed to create Set from exclusion group IDs", map[string]interface{}{
				"error": diags.Errors(),
			})
		} else {
			assignmentBlocks = append(assignmentBlocks, AssignmentResourceModel{
				Target:   types.StringValue("exclude"),
				GroupIds: excludeSet,
			})
		}
	}

	data.Assignments = assignmentBlocks

	tflog.Debug(ctx, "Finished mapping assignments into Terraform state", map[string]interface{}{
		"assignmentBlockCount": len(assignmentBlocks),
	})
}
