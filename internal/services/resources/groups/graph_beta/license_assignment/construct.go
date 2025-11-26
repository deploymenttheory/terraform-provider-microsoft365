package graphBetaGroupLicenseAssignment

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/groups"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// constructAddLicenseRequest constructs a request to add a single license (used for Create)
func constructAddLicensesRequest(ctx context.Context, data *GroupLicenseAssignmentResourceModel) (groups.ItemAssignLicensePostRequestBodyable, error) {
	tflog.Debug(ctx, "Constructing add license request from Terraform configuration")

	requestBody := groups.NewItemAssignLicensePostRequestBody()
	assignedLicense := graphmodels.NewAssignedLicense()

	// Set SKU ID
	if !data.SkuId.IsNull() && !data.SkuId.IsUnknown() {
		skuIdStr := data.SkuId.ValueString()
		skuId, err := uuid.Parse(skuIdStr)
		if err != nil {
			return nil, fmt.Errorf("invalid SKU ID format: %s", skuIdStr)
		}
		assignedLicense.SetSkuId(&skuId)
	}

	// Set disabled plans if present
	if !data.DisabledPlans.IsNull() && !data.DisabledPlans.IsUnknown() {
		disabledPlansElements := data.DisabledPlans.Elements()
		disabledPlans := make([]uuid.UUID, 0, len(disabledPlansElements))

		for _, planVal := range disabledPlansElements {
			if strVal, ok := planVal.(types.String); ok {
				planUUID, err := uuid.Parse(strVal.ValueString())
				if err != nil {
					return nil, fmt.Errorf("invalid disabled plan ID format: %s", strVal.ValueString())
				}
				disabledPlans = append(disabledPlans, planUUID)
			}
		}

		if len(disabledPlans) > 0 {
			assignedLicense.SetDisabledPlans(disabledPlans)
		}
	}

	requestBody.SetAddLicenses([]graphmodels.AssignedLicenseable{assignedLicense})
	requestBody.SetRemoveLicenses([]uuid.UUID{})

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s (Add License)", ResourceName), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]any{
			"error": err.Error(),
		})
	}

	return requestBody, nil
}

// constructUpdateLicenseRequest constructs a request to update a single license (used for Update)
func constructUpdateLicenseRequest(ctx context.Context, data *GroupLicenseAssignmentResourceModel) (groups.ItemAssignLicensePostRequestBodyable, error) {
	tflog.Debug(ctx, "Constructing update license request from Terraform configuration")

	requestBody := groups.NewItemAssignLicensePostRequestBody()
	assignedLicense := graphmodels.NewAssignedLicense()

	// Set SKU ID
	if !data.SkuId.IsNull() && !data.SkuId.IsUnknown() {
		skuIdStr := data.SkuId.ValueString()
		skuId, err := uuid.Parse(skuIdStr)
		if err != nil {
			return nil, fmt.Errorf("invalid SKU ID format: %s", skuIdStr)
		}
		assignedLicense.SetSkuId(&skuId)
	}

	// Set disabled plans if present
	if !data.DisabledPlans.IsNull() && !data.DisabledPlans.IsUnknown() {
		disabledPlansElements := data.DisabledPlans.Elements()
		disabledPlans := make([]uuid.UUID, 0, len(disabledPlansElements))

		for _, planVal := range disabledPlansElements {
			if strVal, ok := planVal.(types.String); ok {
				planUUID, err := uuid.Parse(strVal.ValueString())
				if err != nil {
					return nil, fmt.Errorf("invalid disabled plan ID format: %s", strVal.ValueString())
				}
				disabledPlans = append(disabledPlans, planUUID)
			}
		}

		if len(disabledPlans) > 0 {
			assignedLicense.SetDisabledPlans(disabledPlans)
		}
	}

	requestBody.SetAddLicenses([]graphmodels.AssignedLicenseable{assignedLicense})
	requestBody.SetRemoveLicenses([]uuid.UUID{})

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s (Update License)", ResourceName), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]any{
			"error": err.Error(),
		})
	}

	return requestBody, nil
}

// constructRemoveLicenseRequest constructs a request to remove a single license (used for Delete)
func constructRemoveLicenseRequest(ctx context.Context, skuId string) (groups.ItemAssignLicensePostRequestBodyable, error) {
	tflog.Debug(ctx, "Constructing remove license request")

	requestBody := groups.NewItemAssignLicensePostRequestBody()
	requestBody.SetAddLicenses([]graphmodels.AssignedLicenseable{})

	licenseUUID, err := uuid.Parse(skuId)
	if err != nil {
		return nil, fmt.Errorf("invalid license ID format: %s", skuId)
	}

	requestBody.SetRemoveLicenses([]uuid.UUID{licenseUUID})

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s (Remove License)", ResourceName), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]any{
			"error": err.Error(),
		})
	}

	return requestBody, nil
}
