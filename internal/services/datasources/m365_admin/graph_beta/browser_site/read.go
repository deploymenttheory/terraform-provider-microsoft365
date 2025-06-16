package graphBetaBrowserSite

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Read handles the Read operation for Browser Site data source.
func (d *BrowserSiteDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var object BrowserSiteDataSourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting Read method for: %s_%s", d.ProviderTypeName, d.TypeName))

	resp.Diagnostics.Append(req.Config.Get(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	filterType := object.FilterType.ValueString()
	tflog.Debug(ctx, fmt.Sprintf("Reading %s_%s with filter_type: %s", d.ProviderTypeName, d.TypeName, filterType))

	// Check if we have a valid browser site list ID
	if object.BrowserSiteListAssignmentID.IsNull() || object.BrowserSiteListAssignmentID.ValueString() == "" {
		resp.Diagnostics.AddError(
			"Missing Required Parameter",
			"browser_site_list_assignment_id must be provided to query browser sites.",
		)
		return
	}

	// For non-"all" filter types, ensure we have a filter value
	if filterType != "all" && (object.FilterValue.IsNull() || object.FilterValue.ValueString() == "") {
		resp.Diagnostics.AddError(
			"Missing Required Parameter",
			fmt.Sprintf("filter_value must be provided when filter_type is '%s'", filterType),
		)
		return
	}

	ctx, cancel := crud.HandleTimeout(ctx, object.Timeouts.Read, ReadDataSourceTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	browserSiteListId := object.BrowserSiteListAssignmentID.ValueString()

	// Serialize access to the Graph API to prevent concurrent map writes in Kiota

	respList, err := d.client.
		Admin().
		Edge().
		InternetExplorerMode().
		SiteLists().
		ByBrowserSiteListId(browserSiteListId).
		Sites().
		Get(ctx, nil)

	if err != nil {
		errors.HandleGraphError(ctx, err, resp, "Read", d.ReadPermissions)
		return
	}

	if respList == nil || respList.GetValue() == nil {
		tflog.Debug(ctx, "Received a nil or empty response from the Microsoft Graph API")
		object.Items = []BrowserSiteResourceModel{}
		resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Received %d browser sites from the API", len(respList.GetValue())))

	var filteredItems []BrowserSiteResourceModel
	filterValue := object.FilterValue.ValueString()

	for _, item := range respList.GetValue() {
		if item == nil {
			continue
		}

		if item.GetId() != nil && item.GetWebUrl() != nil {
			tflog.Debug(ctx, fmt.Sprintf("Processing browser site: ID=%s, WebUrl=%s",
				*item.GetId(),
				*item.GetWebUrl()))
		}

		var matchesFilter bool
		switch filterType {
		case "all":
			matchesFilter = true

		case "id":
			matchesFilter = item.GetId() != nil && *item.GetId() == filterValue

		case "web_url":
			matchesFilter = item.GetWebUrl() != nil && strings.Contains(
				strings.ToLower(*item.GetWebUrl()),
				strings.ToLower(filterValue))
		}

		if matchesFilter {
			// Create a new item and map the remote state to it
			var siteItem BrowserSiteResourceModel
			MapRemoteStateToDataSource(ctx, &siteItem, item, object.BrowserSiteListAssignmentID)
			filteredItems = append(filteredItems, siteItem)
		}
	}

	// Handle pagination if needed - check if more results are available
	if respList.GetOdataNextLink() != nil && *respList.GetOdataNextLink() != "" {
		tflog.Warn(ctx, "Pagination detected but not implemented - only returning first page of results")
	}

	object.Items = filteredItems

	tflog.Debug(ctx, fmt.Sprintf("Found %d browser sites after filtering", len(filteredItems)))

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Datasource Read Method: %s_%s, found %d items", d.ProviderTypeName, d.TypeName, len(filteredItems)))
}
