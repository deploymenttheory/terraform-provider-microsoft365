package graphBetaUserLicenseAssignment

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/state"
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

// MapRemoteResourceStateToTerraform maps the properties of a User to Terraform state for license assignment.
func MapRemoteResourceStateToTerraform(ctx context.Context, data *UserLicenseAssignmentResourceModel, remoteResource graphmodels.Userable) {
	if remoteResource == nil {
		tflog.Debug(ctx, "Remote resource is nil")
		return
	}

	tflog.Debug(ctx, "Starting to map remote state to Terraform state", map[string]interface{}{
		"resourceId": remoteResource.GetId(),
	})

	// Set basic user information
	data.ID = state.StringPointerValue(remoteResource.GetId())
	data.UserId = state.StringPointerValue(remoteResource.GetId())
	data.UserPrincipalName = state.StringPointerValue(remoteResource.GetUserPrincipalName())

	tflog.Debug(ctx, fmt.Sprintf("Finished mapping user license assignment resource with id %s", data.ID.ValueString()))
}

// MapLicenseDetailsToTerraform maps the license details response to Terraform state.
func MapLicenseDetailsToTerraform(ctx context.Context, data *UserLicenseAssignmentResourceModel, licenseDetailsResponse graphmodels.LicenseDetailsCollectionResponseable) {
	if licenseDetailsResponse == nil {
		tflog.Debug(ctx, "License details response is nil")
		return
	}

	licenseDetails := licenseDetailsResponse.GetValue()
	if len(licenseDetails) == 0 {
		tflog.Debug(ctx, "No license details found")
		data.AssignedLicenses = types.ListNull(types.ObjectType{
			AttrTypes: getLicenseDetailsObjectType(),
		})
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Mapping %d license details", len(licenseDetails)))

	// Map license details to Terraform state
	assignedLicenses := make([]attr.Value, 0, len(licenseDetails))

	for _, license := range licenseDetails {
		if license == nil {
			continue
		}

		licenseObj := map[string]attr.Value{
			"sku_id":          uuidPointerToStringValue(license.GetSkuId()),
			"sku_part_number": state.StringPointerValue(license.GetSkuPartNumber()),
			"service_plans":   mapServicePlansToTerraform(ctx, license.GetServicePlans()),
		}

		objValue, diag := types.ObjectValue(getLicenseDetailsObjectType(), licenseObj)
		if diag.HasError() {
			tflog.Error(ctx, fmt.Sprintf("Error creating license object: %v", diag.Errors()))
			continue
		}

		assignedLicenses = append(assignedLicenses, objValue)
	}

	assignedLicensesList, diag := types.ListValue(
		types.ObjectType{AttrTypes: getLicenseDetailsObjectType()},
		assignedLicenses,
	)
	if diag.HasError() {
		tflog.Error(ctx, fmt.Sprintf("Error creating assigned licenses list: %v", diag.Errors()))
		data.AssignedLicenses = types.ListNull(types.ObjectType{
			AttrTypes: getLicenseDetailsObjectType(),
		})
		return
	}

	data.AssignedLicenses = assignedLicensesList
	tflog.Debug(ctx, fmt.Sprintf("Successfully mapped %d assigned licenses", len(assignedLicenses)))
}

// mapServicePlansToTerraform maps service plans to Terraform state.
func mapServicePlansToTerraform(ctx context.Context, servicePlans []graphmodels.ServicePlanInfoable) types.List {
	if len(servicePlans) == 0 {
		return types.ListNull(types.ObjectType{
			AttrTypes: getServicePlanObjectType(),
		})
	}

	servicePlanValues := make([]attr.Value, 0, len(servicePlans))

	for _, plan := range servicePlans {
		if plan == nil {
			continue
		}

		planObj := map[string]attr.Value{
			"service_plan_id":     uuidPointerToStringValue(plan.GetServicePlanId()),
			"service_plan_name":   state.StringPointerValue(plan.GetServicePlanName()),
			"provisioning_status": state.StringPointerValue(plan.GetProvisioningStatus()),
			"applies_to":          state.StringPointerValue(plan.GetAppliesTo()),
		}

		objValue, diag := types.ObjectValue(getServicePlanObjectType(), planObj)
		if diag.HasError() {
			tflog.Error(ctx, fmt.Sprintf("Error creating service plan object: %v", diag.Errors()))
			continue
		}

		servicePlanValues = append(servicePlanValues, objValue)
	}

	servicePlansList, diag := types.ListValue(
		types.ObjectType{AttrTypes: getServicePlanObjectType()},
		servicePlanValues,
	)
	if diag.HasError() {
		tflog.Error(ctx, fmt.Sprintf("Error creating service plans list: %v", diag.Errors()))
		return types.ListNull(types.ObjectType{
			AttrTypes: getServicePlanObjectType(),
		})
	}

	return servicePlansList
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
