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

// constructResource maps the Terraform configuration to a group license assignment request
func constructResource(ctx context.Context, data *GroupLicenseAssignmentResourceModel) (groups.ItemAssignLicensePostRequestBodyable, error) {
	tflog.Debug(ctx, "Constructing group license assignment request from Terraform configuration")

	requestBody := groups.NewItemAssignLicensePostRequestBody()

	addLicenses := make([]graphmodels.AssignedLicenseable, 0)
	for _, license := range data.AddLicenses {
		assignedLicense := graphmodels.NewAssignedLicense()

		if !license.SkuId.IsNull() && !license.SkuId.IsUnknown() {
			skuIdStr := license.SkuId.ValueString()
			skuId, err := uuid.Parse(skuIdStr)
			if err != nil {
				return nil, fmt.Errorf("invalid SKU ID format: %s", skuIdStr)
			}
			assignedLicense.SetSkuId(&skuId)
		}

		if !license.DisabledPlans.IsNull() && !license.DisabledPlans.IsUnknown() {
			disabledPlansElements := license.DisabledPlans.Elements()
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

		addLicenses = append(addLicenses, assignedLicense)
	}
	requestBody.SetAddLicenses(addLicenses)

	if !data.RemoveLicenses.IsNull() && !data.RemoveLicenses.IsUnknown() {
		removeLicensesElements := data.RemoveLicenses.Elements()
		removeLicenses := make([]uuid.UUID, 0, len(removeLicensesElements))

		for _, licenseVal := range removeLicensesElements {
			if strVal, ok := licenseVal.(types.String); ok {
				licenseUUID, err := uuid.Parse(strVal.ValueString())
				if err != nil {
					return nil, fmt.Errorf("invalid remove license ID format: %s", strVal.ValueString())
				}
				removeLicenses = append(removeLicenses, licenseUUID)
			}
		}

		requestBody.SetRemoveLicenses(removeLicenses)
	} else {
		requestBody.SetRemoveLicenses([]uuid.UUID{})
	}

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]any{
			"error": err.Error(),
		})
	}

	return requestBody, nil
}
