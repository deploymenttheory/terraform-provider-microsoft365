package graphBetaGroupLicenseAssignment

import (
	"context"
	"fmt"
	"strings"

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
// It returns true when the managed SKU is present in the group's assignedLicenses, and false when
// the group exists but the license assignment does not — callers use this to distinguish
// "group found" from "license assignment found".
func MapRemoteResourceStateToTerraform(ctx context.Context, data *GroupLicenseAssignmentResourceModel, remoteResource graphmodels.Groupable) bool {
	if remoteResource == nil {
		tflog.Debug(ctx, "Remote resource is nil")
		return false
	}

	tflog.Debug(ctx, "Starting to map remote state to Terraform state", map[string]any{
		"resourceId": remoteResource.GetId(),
	})

	data.GroupId = convert.GraphToFrameworkString(remoteResource.GetId())
	data.DisplayName = convert.GraphToFrameworkString(remoteResource.GetDisplayName())

	assignedLicenses := remoteResource.GetAssignedLicenses()
	managedSkuId := data.SkuId.ValueString()

	data.ID = types.StringValue(fmt.Sprintf("%s_%s", data.GroupId.ValueString(), managedSkuId))

	// Default to empty set so that if the SKU is not present in assignedLicenses (e.g. during
	// propagation delay) we don't leave stale disabled-plan data in state.
	data.DisabledPlans = types.SetValueMust(types.StringType, []attr.Value{})

	skuFound := false
	for _, license := range assignedLicenses {
		if license == nil {
			continue
		}

		licenseSkuId := uuidPointerToStringValue(license.GetSkuId())
		tflog.Debug(ctx, fmt.Sprintf("Checking license SKU: %s (looking for: %s)", licenseSkuId.ValueString(), managedSkuId))

		// The API returns SKU ids in canonical lowercase form while the configured
		// sku_id may use any casing, so compare case-insensitively.
		if strings.EqualFold(licenseSkuId.ValueString(), managedSkuId) {
			skuFound = true
			disabledPlans := license.GetDisabledPlans()
			tflog.Debug(ctx, fmt.Sprintf("Found matching license, disabled plans count: %d", len(disabledPlans)))

			disabledPlanValues := make([]attr.Value, 0, len(disabledPlans))
			for _, planUUID := range disabledPlans {
				planStr := planUUID.String()
				tflog.Debug(ctx, fmt.Sprintf("Mapping disabled plan: %s", planStr))
				disabledPlanValues = append(disabledPlanValues, types.StringValue(planStr))
			}
			data.DisabledPlans = types.SetValueMust(types.StringType, disabledPlanValues)
			tflog.Debug(ctx, fmt.Sprintf("Set disabled_plans in state with %d items", len(disabledPlanValues)))
			break
		}
	}

	tflog.Debug(ctx, fmt.Sprintf("Final disabled_plans in state: %v", data.DisabledPlans))

	tflog.Debug(ctx, fmt.Sprintf("Finished mapping group license assignment resource with id %s (sku found: %t)", data.ID.ValueString(), skuFound))

	return skuFound
}
