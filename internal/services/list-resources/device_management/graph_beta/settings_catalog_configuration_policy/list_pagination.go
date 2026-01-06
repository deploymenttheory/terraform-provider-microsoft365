package graphBetaSettingsCatalogConfigurationPolicy

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/devicemanagement"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models"
	graphcore "github.com/microsoftgraph/msgraph-sdk-go-core"
)

// listAllResourcesWithPageIterator retrieves all configuration policies using the PageIterator
// This handles pagination automatically and returns ALL policies across all pages
func (r *SettingsCatalogListResource) listAllResourcesWithPageIterator(
	ctx context.Context,
	requestConfig *devicemanagement.ConfigurationPoliciesRequestBuilderGetRequestConfiguration,
) ([]models.DeviceManagementConfigurationPolicyable, error) {
	var allPolicies []models.DeviceManagementConfigurationPolicyable

	tflog.Debug(ctx, "Fetching first page of Settings Catalog policies")

	policiesResponse, err := r.client.
		DeviceManagement().
		ConfigurationPolicies().
		Get(ctx, requestConfig)

	if err != nil {
		return nil, fmt.Errorf("failed to get policies: %w", err)
	}

	pageIterator, err := graphcore.NewPageIterator[models.DeviceManagementConfigurationPolicyable](
		policiesResponse,
		r.client.GetAdapter(),
		models.CreateDeviceManagementConfigurationPolicyCollectionResponseFromDiscriminatorValue,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create page iterator: %w", err)
	}

	pageCount := 0
	err = pageIterator.Iterate(ctx, func(item models.DeviceManagementConfigurationPolicyable) bool {
		if item != nil {
			allPolicies = append(allPolicies, item)

			if len(allPolicies)%100 == 0 {
				pageCount++
				tflog.Debug(ctx, fmt.Sprintf("PageIterator: collected %d policies (estimated page %d)", len(allPolicies), pageCount))
			}
		}
		return true
	})

	if err != nil {
		return nil, fmt.Errorf("failed to iterate pages: %w", err)
	}

	tflog.Debug(ctx, fmt.Sprintf("PageIterator complete: collected %d total policies", len(allPolicies)))

	return allPolicies, nil
}
