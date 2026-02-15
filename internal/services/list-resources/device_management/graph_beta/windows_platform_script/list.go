package graphBetaWindowsPlatformScript

import (
	"context"
	"fmt"

	sharedmodels "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/shared_models/graph_beta"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/devicemanagement"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// List retrieves and returns a stream of Windows platform scripts.
func (r *WindowsPlatformScriptListResource) List(ctx context.Context, req list.ListRequest, stream *list.ListResultsStream) {
	var data WindowsPlatformScriptListConfigModel

	diags := req.Config.Get(ctx, &data)
	if diags.HasError() {
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}

	requestConfig := &devicemanagement.DeviceManagementScriptsRequestBuilderGetRequestConfiguration{
		QueryParameters: &devicemanagement.DeviceManagementScriptsRequestBuilderGetQueryParameters{},
	}

	var filter string

	// If custom OData filter is provided, use it directly
	if !data.ODataFilter.IsNull() && !data.ODataFilter.IsUnknown() {
		filter = data.ODataFilter.ValueString()
		tflog.Debug(ctx, fmt.Sprintf("Using custom OData filter: %s", filter))
	} else {
		// Build filter from individual parameters using dedicated builder functions
		// Note: isAssigned is NOT included here because the API field is unreliable
		// We'll filter by assignments locally after fetching scripts
		filter = combineFilters(
			buildDisplayNameFilter(&data),
			buildFileNameFilter(&data),
			buildRunAsAccountFilter(&data),
		)

		if filter != "" {
			tflog.Debug(ctx, fmt.Sprintf("Built filter query: %s", filter))
		}
	}

	// Apply filter to request configuration
	if filter != "" {
		requestConfig.QueryParameters.Filter = &filter
	}

	allScripts, err := r.listAllResourcesWithPageIterator(ctx, requestConfig)

	if err != nil {
		result := req.NewListResult(ctx)
		result.Diagnostics.AddError(
			"Error listing Windows platform scripts",
			fmt.Sprintf("Could not list Windows platform scripts: %s", err.Error()),
		)
		stream.Results = func(push func(list.ListResult) bool) {
			push(result)
		}
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Retrieved %d total Windows platform scripts", len(allScripts)))

	// Apply local assignment filter if specified
	// This is done locally because the API's isAssigned field is unreliable
	if !data.IsAssignedFilter.IsNull() && !data.IsAssignedFilter.IsUnknown() {
		filterByAssigned := data.IsAssignedFilter.ValueBool()
		tflog.Debug(ctx, fmt.Sprintf("Applying local assignment filter: isAssigned=%t", filterByAssigned))

		filteredScripts := []models.DeviceManagementScriptable{}
		for _, script := range allScripts {
			if script == nil || script.GetId() == nil {
				continue
			}

			hasAssignments, err := r.listResourceAssignments(ctx, *script.GetId())
			if err != nil {
				tflog.Warn(ctx, fmt.Sprintf("Failed to check assignments for script %s: %v", *script.GetId(), err))
				continue
			}

			if hasAssignments == filterByAssigned {
				filteredScripts = append(filteredScripts, script)
			}
		}

		tflog.Debug(ctx, fmt.Sprintf("After assignment filter: %d scripts match", len(filteredScripts)))
		allScripts = filteredScripts
	}

	stream.Results = streamResults(ctx, req, allScripts)
}

// streamResults processes a list of Graph API resources and streams them as list results.
// This helper extracts the ID and display name from each resource and sets the identity.
func streamResults(ctx context.Context, req list.ListRequest, items []models.DeviceManagementScriptable) func(push func(list.ListResult) bool) {
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
