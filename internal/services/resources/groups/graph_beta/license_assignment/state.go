package graphBetaGroupLicenseAssignment

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// Helper function to convert UUID pointer to string value
func uuidPointerToStringValue(uuidPtr *uuid.UUID) types.String {
	if uuidPtr == nil {
		return types.StringNull()
	}
	return types.StringValue(uuidPtr.String())
}

// MapRemoteResourceStateToTerraform maps the properties of a Group to Terraform state for license assignment.
func MapRemoteResourceStateToTerraform(ctx context.Context, data *GroupLicenseAssignmentResourceModel, remoteResource graphmodels.Groupable) {
	if remoteResource == nil {
		tflog.Debug(ctx, "Remote resource is nil")
		return
	}

	tflog.Debug(ctx, "Starting to map remote state to Terraform state", map[string]any{
		"resourceId": remoteResource.GetId(),
	})

	data.GroupId = convert.GraphToFrameworkString(remoteResource.GetId())
	data.DisplayName = convert.GraphToFrameworkString(remoteResource.GetDisplayName())

	assignedLicenses := remoteResource.GetAssignedLicenses()
	managedSkuId := data.SkuId.ValueString()

	data.ID = types.StringValue(fmt.Sprintf("%s_%s", data.GroupId.ValueString(), managedSkuId))

	for _, license := range assignedLicenses {
		if license == nil {
			continue
		}

		licenseSkuId := uuidPointerToStringValue(license.GetSkuId())
		if licenseSkuId.ValueString() == managedSkuId {
			disabledPlans := license.GetDisabledPlans()
			disabledPlanValues := make([]attr.Value, 0, len(disabledPlans))
			for _, planUUID := range disabledPlans {
				disabledPlanValues = append(disabledPlanValues, types.StringValue(planUUID.String()))
			}
			data.DisabledPlans = types.SetValueMust(types.StringType, disabledPlanValues)
			break
		}
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished mapping group license assignment resource with id %s", data.ID.ValueString()))
}
