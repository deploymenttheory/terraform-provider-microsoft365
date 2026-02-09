package graphBetaConditionalAccessPolicy

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/identity"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models"
	graphcore "github.com/microsoftgraph/msgraph-sdk-go-core"
)

// listAllResourcesWithPageIterator retrieves all Conditional Access policies using the PageIterator
// This handles pagination automatically and returns ALL policies across all pages
func (r *ConditionalAccessPolicyListResource) listAllResourcesWithPageIterator(
	ctx context.Context,
	requestConfig *identity.ConditionalAccessPoliciesRequestBuilderGetRequestConfiguration,
) ([]models.ConditionalAccessPolicyable, error) {
	var allPolicies []models.ConditionalAccessPolicyable

	tflog.Debug(ctx, "Fetching first page of Conditional Access policies")

	policiesResponse, err := r.client.
		Identity().
		ConditionalAccess().
		Policies().
		Get(ctx, requestConfig)

	if err != nil {
		return nil, fmt.Errorf("failed to get policies: %w", err)
	}

	pageIterator, err := graphcore.NewPageIterator[models.ConditionalAccessPolicyable](
		policiesResponse,
		r.client.GetAdapter(),
		models.CreateConditionalAccessPolicyCollectionResponseFromDiscriminatorValue,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create page iterator: %w", err)
	}

	pageCount := 0
	err = pageIterator.Iterate(ctx, func(item models.ConditionalAccessPolicyable) bool {
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
