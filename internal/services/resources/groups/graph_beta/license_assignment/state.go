package graphBetaGroupLicenseAssignment

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/state"
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

	tflog.Debug(ctx, "Starting to map remote state to Terraform state", map[string]interface{}{
		"resourceId": remoteResource.GetId(),
	})

	// Set basic group information
	data.ID = state.StringPointerValue(remoteResource.GetId())
	data.GroupId = state.StringPointerValue(remoteResource.GetId())
	data.DisplayName = state.StringPointerValue(remoteResource.GetDisplayName())

	// Map assigned licenses
	assignedLicenses := remoteResource.GetAssignedLicenses()
	if len(assignedLicenses) > 0 {
		MapAssignedLicensesToTerraform(ctx, data, assignedLicenses)
	} else {
		data.AssignedLicenses = types.ListNull(types.ObjectType{
			AttrTypes: getLicenseDetailsObjectType(),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished mapping group license assignment resource with id %s", data.ID.ValueString()))
}

// MapAssignedLicensesToTerraform maps the assigned licenses to Terraform state.
func MapAssignedLicensesToTerraform(ctx context.Context, data *GroupLicenseAssignmentResourceModel, assignedLicenses []graphmodels.AssignedLicenseable) {
	if len(assignedLicenses) == 0 {
		tflog.Debug(ctx, "No assigned licenses found")
		data.AssignedLicenses = types.ListNull(types.ObjectType{
			AttrTypes: getLicenseDetailsObjectType(),
		})
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Mapping %d assigned licenses", len(assignedLicenses)))

	// Map assigned licenses to Terraform state
	assignedLicensesList := make([]attr.Value, 0, len(assignedLicenses))

	for _, license := range assignedLicenses {
		if license == nil {
			continue
		}

		// Create a simplified representation for assigned licenses
		// Note: Unlike user license details, group assigned licenses don't include service plans
		licenseObj := map[string]attr.Value{
			"sku_id":          uuidPointerToStringValue(license.GetSkuId()),
			"sku_part_number": types.StringNull(), // Not available in AssignedLicense
			"service_plans":   types.ListNull(types.ObjectType{AttrTypes: getServicePlanObjectType()}),
		}

		objValue, diag := types.ObjectValue(getLicenseDetailsObjectType(), licenseObj)
		if diag.HasError() {
			tflog.Error(ctx, fmt.Sprintf("Error creating license object: %v", diag.Errors()))
			continue
		}

		assignedLicensesList = append(assignedLicensesList, objValue)
	}

	assignedLicensesListValue, diag := types.ListValue(
		types.ObjectType{AttrTypes: getLicenseDetailsObjectType()},
		assignedLicensesList,
	)
	if diag.HasError() {
		tflog.Error(ctx, fmt.Sprintf("Error creating assigned licenses list: %v", diag.Errors()))
		data.AssignedLicenses = types.ListNull(types.ObjectType{
			AttrTypes: getLicenseDetailsObjectType(),
		})
		return
	}

	data.AssignedLicenses = assignedLicensesListValue
	tflog.Debug(ctx, fmt.Sprintf("Successfully mapped %d assigned licenses", len(assignedLicensesList)))
}

// getLicenseDetailsObjectType returns the object type for license details.
func getLicenseDetailsObjectType() map[string]attr.Type {
	return map[string]attr.Type{
		"sku_id":          types.StringType,
		"sku_part_number": types.StringType,
		"service_plans": types.ListType{
			ElemType: types.ObjectType{
				AttrTypes: getServicePlanObjectType(),
			},
		},
	}
}

// getServicePlanObjectType returns the object type for service plans.
func getServicePlanObjectType() map[string]attr.Type {
	return map[string]attr.Type{
		"service_plan_id":     types.StringType,
		"service_plan_name":   types.StringType,
		"provisioning_status": types.StringType,
		"applies_to":          types.StringType,
	}
}
