package graphBetaConditionalAccessPolicy

import (
	"context"
	"fmt"

	sharedmodels "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/shared_models/graph_beta"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/identity"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// List retrieves and returns a stream of Conditional Access policies.
func (r *ConditionalAccessPolicyListResource) List(ctx context.Context, req list.ListRequest, stream *list.ListResultsStream) {
	var data ConditionalAccessPolicyListConfigModel

	diags := req.Config.Get(ctx, &data)
	if diags.HasError() {
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}

	// Build query parameters
	requestConfig := &identity.ConditionalAccessPoliciesRequestBuilderGetRequestConfiguration{
		QueryParameters: &identity.ConditionalAccessPoliciesRequestBuilderGetQueryParameters{},
	}

	// Build OData filter query
	var filter string

	// If custom OData filter is provided, use it directly
	if !data.ODataFilter.IsNull() && !data.ODataFilter.IsUnknown() {
		filter = data.ODataFilter.ValueString()
		tflog.Debug(ctx, fmt.Sprintf("Using custom OData filter: %s", filter))
	} else {
		// Build filter from individual parameters using dedicated builder functions
		filter = combineFilters(
			buildDisplayNameFilter(&data),
			buildStateFilter(&data),
		)

		if filter != "" {
			tflog.Debug(ctx, fmt.Sprintf("Built filter query: %s", filter))
		}
	}

	// Apply filter to request configuration
	if filter != "" {
		requestConfig.QueryParameters.Filter = &filter
	}

	// Get all policies using PageIterator (handles pagination automatically)
	allPolicies, err := r.listAllResourcesWithPageIterator(ctx, requestConfig)

	if err != nil {
		result := req.NewListResult(ctx)
		result.Diagnostics.AddError(
			"Error listing Conditional Access policies",
			fmt.Sprintf("Could not list Conditional Access policies: %s", err.Error()),
		)
		stream.Results = func(push func(list.ListResult) bool) {
			push(result)
		}
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Retrieved %d total Conditional Access policies", len(allPolicies)))

	stream.Results = streamResults(ctx, req, allPolicies)
}

// streamResults processes a list of Graph API resources and streams them as list results.
// This helper extracts the ID and display name from each resource and sets the identity.
func streamResults(ctx context.Context, req list.ListRequest, items []models.ConditionalAccessPolicyable) func(push func(list.ListResult) bool) {
	return func(push func(list.ListResult) bool) {
		for _, item := range items {

			if item == nil {
				continue
			}

			if item.GetId() == nil || *item.GetId() == "" {
				continue
			}

			result := req.NewListResult(ctx)

			if item.GetDisplayName() != nil {
				result.DisplayName = *item.GetDisplayName()
			}

			// Set the identity using shared struct
			var identity sharedmodels.ResourceIdentity
			identity.ID = *item.GetId()

			result.Diagnostics.Append(result.Identity.Set(ctx, identity)...)

			if !push(result) {
				return
			}
		}
	}
}
