package graphBetaUsersUser

import (
	"context"
	"fmt"

	sharedmodels "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/shared_models/graph_beta"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/users"
)

// List retrieves and returns a stream of users.
func (r *UserListResource) List(ctx context.Context, req list.ListRequest, stream *list.ListResultsStream) {
	var data UserListConfigModel

	diags := req.Config.Get(ctx, &data)
	if diags.HasError() {
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}

	requestConfig := &users.UsersRequestBuilderGetRequestConfiguration{
		QueryParameters: &users.UsersRequestBuilderGetQueryParameters{},
	}

	var filter string

	// If custom OData filter is provided, use it directly
	if !data.ODataFilter.IsNull() && !data.ODataFilter.IsUnknown() {
		filter = data.ODataFilter.ValueString()
		tflog.Debug(ctx, fmt.Sprintf("Using custom OData filter: %s", filter))
	} else {
		// Build filter from individual parameters using dedicated builder functions
		filter = combineFilters(
			buildDisplayNameFilter(&data),
			buildUserPrincipalNameFilter(&data),
			buildAccountEnabledFilter(&data),
			buildUserTypeFilter(&data),
		)

		if filter != "" {
			tflog.Debug(ctx, fmt.Sprintf("Built filter query: %s", filter))
		}
	}

	if filter != "" {
		requestConfig.QueryParameters.Filter = &filter
	}

	allUsers, err := r.listAllResourcesWithPageIterator(ctx, requestConfig)

	if err != nil {
		result := req.NewListResult(ctx)
		result.Diagnostics.AddError(
			"Error listing users",
			fmt.Sprintf("Could not list users: %s", err.Error()),
		)
		stream.Results = func(push func(list.ListResult) bool) {
			push(result)
		}
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Retrieved %d total users", len(allUsers)))

	stream.Results = streamResults(ctx, req, allUsers)
}

// streamResults processes a list of Graph API resources and streams them as list results.
// This helper extracts the ID and display name from each resource and sets the identity.
func streamResults(ctx context.Context, req list.ListRequest, items []models.Userable) func(push func(list.ListResult) bool) {
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

			var identity sharedmodels.ResourceIdentity
			identity.ID = *item.GetId()

			result.Diagnostics.Append(result.Identity.Set(ctx, identity)...)

			if !push(result) {
				return
			}
		}
	}
}
