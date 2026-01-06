package graphBetaSettingsCatalogConfigurationPolicy

import (
	"context"
	"fmt"

	sharedmodels "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/shared_models/graph_beta"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/devicemanagement"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// List retrieves and returns a stream of Settings Catalog configuration policies.
func (r *SettingsCatalogListResource) List(ctx context.Context, req list.ListRequest, stream *list.ListResultsStream) {
	var data SettingsCatalogListConfigModel

	diags := req.Config.Get(ctx, &data)
	if diags.HasError() {
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}

	// Build query parameters
	requestConfig := &devicemanagement.ConfigurationPoliciesRequestBuilderGetRequestConfiguration{
		QueryParameters: &devicemanagement.ConfigurationPoliciesRequestBuilderGetQueryParameters{},
	}

	// Build OData filter query
	var filter string

	// If custom OData filter is provided, use it directly
	if !data.ODataFilter.IsNull() && !data.ODataFilter.IsUnknown() {
		filter = data.ODataFilter.ValueString()
		tflog.Debug(ctx, fmt.Sprintf("Using custom OData filter: %s", filter))
	} else {
		// Build filter from individual parameters using dedicated builder functions
		// Note: isAssigned is NOT included here because the API field is unreliable
		// We'll filter by assignments locally after fetching policies
		filter = combineFilters(
			buildNameFilter(&data),
			buildPlatformFilter(ctx, &data),
			buildTemplateFamilyFilter(&data),
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
			"Error listing Settings Catalog policies",
			fmt.Sprintf("Could not list Settings Catalog policies: %s", err.Error()),
		)
		stream.Results = func(push func(list.ListResult) bool) {
			push(result)
		}
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Retrieved %d total Settings Catalog policies", len(allPolicies)))

	// Apply local assignment filter if specified
	// This is done locally because the API's isAssigned field is unreliable
	if !data.IsAssignedFilter.IsNull() && !data.IsAssignedFilter.IsUnknown() {
		filterByAssigned := data.IsAssignedFilter.ValueBool()
		tflog.Debug(ctx, fmt.Sprintf("Applying local assignment filter: isAssigned=%t", filterByAssigned))

		filteredPolicies := []models.DeviceManagementConfigurationPolicyable{}
		for _, policy := range allPolicies {
			if policy == nil || policy.GetId() == nil {
				continue
			}

			hasAssignments, err := r.listResourceAssignments(ctx, *policy.GetId())
			if err != nil {
				tflog.Warn(ctx, fmt.Sprintf("Failed to check assignments for policy %s: %v", *policy.GetId(), err))
				continue
			}

			// Include policy if it matches the filter
			if hasAssignments == filterByAssigned {
				filteredPolicies = append(filteredPolicies, policy)
			}
		}

		tflog.Debug(ctx, fmt.Sprintf("After assignment filter: %d policies match", len(filteredPolicies)))
		allPolicies = filteredPolicies
	}

	stream.Results = streamResults(ctx, req, allPolicies)
}

// streamResults processes a list of Graph API resources and streams them as list results.
// This helper extracts the ID and display name from each resource and sets the identity.
func streamResults(ctx context.Context, req list.ListRequest, items []models.DeviceManagementConfigurationPolicyable) func(push func(list.ListResult) bool) {
	return func(push func(list.ListResult) bool) {
		for _, item := range items {

			if item == nil {
				continue
			}

			if item.GetId() == nil || *item.GetId() == "" {
				continue
			}

			result := req.NewListResult(ctx)

			if item.GetName() != nil {
				result.DisplayName = *item.GetName()
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
