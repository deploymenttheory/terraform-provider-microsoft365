package graphBetaGroupLicenseAssignment

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/groups"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// constructGroupLicenseAssignmentRequest maps the Terraform configuration to a group license assignment request
func constructGroupLicenseAssignmentRequest(ctx context.Context, data *GroupLicenseAssignmentResourceModel) (groups.ItemAssignLicensePostRequestBodyable, error) {
	tflog.Debug(ctx, "Constructing group license assignment request from Terraform configuration")

	requestBody := groups.NewItemAssignLicensePostRequestBody()

	// Process add_licenses
	addLicenses := make([]graphmodels.AssignedLicenseable, 0)
	for _, license := range data.AddLicenses {
		assignedLicense := graphmodels.NewAssignedLicense()

		// Set SKU ID - convert string to UUID
		if !license.SkuId.IsNull() && !license.SkuId.IsUnknown() {
			skuIdStr := license.SkuId.ValueString()
			skuId, err := uuid.Parse(skuIdStr)
			if err != nil {
				return nil, fmt.Errorf("invalid SKU ID format: %s", skuIdStr)
			}
			assignedLicense.SetSkuId(&skuId)
		}

		// Set disabled plans if provided - convert strings to UUIDs
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

	// Process remove_licenses - convert strings to UUIDs
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
		// Set empty array if no licenses to remove
		requestBody.SetRemoveLicenses([]uuid.UUID{})
	}

	tflog.Debug(ctx, fmt.Sprintf("Constructed group license assignment request with %d licenses to add and %d licenses to remove",
		len(addLicenses), len(requestBody.GetRemoveLicenses())))

	return requestBody, nil
}
