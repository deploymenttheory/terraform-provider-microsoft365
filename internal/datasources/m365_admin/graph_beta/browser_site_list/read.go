package graphBetaBrowserSiteList

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/crud"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/errors"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Read handles the Read operation for Browser Site Lists data source.
func (d *BrowserSiteListDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var object BrowserSiteListDataSourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting Read method for: %s_%s", d.ProviderTypeName, d.TypeName))

	resp.Diagnostics.Append(req.Config.Get(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	filterType := object.FilterType.ValueString()
	tflog.Debug(ctx, fmt.Sprintf("Reading %s_%s with filter_type: %s", d.ProviderTypeName, d.TypeName, filterType))

	if filterType != "all" && (object.FilterValue.IsNull() || object.FilterValue.ValueString() == "") {
		resp.Diagnostics.AddError(
			"Missing Required Parameter",
			fmt.Sprintf("filter_value must be provided when filter_type is '%s'", filterType),
		)
		return
	}

	ctx, cancel := crud.HandleTimeout(ctx, object.Timeouts.Read, ReadTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	// Serialize access to the Graph API to prevent concurrent map writes in Kiota

	respList, err := d.client.
		Admin().
		Edge().
		InternetExplorerMode().
		SiteLists().
		Get(ctx, nil)

	if err != nil {
		errors.HandleGraphError(ctx, err, resp, "Read", d.ReadPermissions)
		return
	}

	if respList == nil || respList.GetValue() == nil {
		tflog.Debug(ctx, "Received a nil or empty response from the Microsoft Graph API")
		object.Items = []BrowserSiteListResourceModel{}
		resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Received %d site lists from the API", len(respList.GetValue())))

	var filteredItems []BrowserSiteListResourceModel
	filterValue := object.FilterValue.ValueString()

	for _, item := range respList.GetValue() {
		if item == nil {
			continue
		}

		if item.GetId() != nil && item.GetDisplayName() != nil {
			tflog.Debug(ctx, fmt.Sprintf("Processing site list: ID=%s, DisplayName=%s",
				*item.GetId(),
				*item.GetDisplayName()))
		}

		switch filterType {
		case "all":
			filteredItems = append(filteredItems, MapRemoteStateToDataSource(item))

		case "id":
			if item.GetId() != nil && *item.GetId() == filterValue {
				filteredItems = append(filteredItems, MapRemoteStateToDataSource(item))
			}

		case "display_name":
			if item.GetDisplayName() != nil && strings.Contains(
				strings.ToLower(*item.GetDisplayName()),
				strings.ToLower(filterValue)) {
				filteredItems = append(filteredItems, MapRemoteStateToDataSource(item))
			}
		}
	}

	// Handle pagination if needed - check if more results are available
	if respList.GetOdataNextLink() != nil && *respList.GetOdataNextLink() != "" {
		tflog.Warn(ctx, "Pagination detected but not implemented - only returning first page of results")
	}

	object.Items = filteredItems

	tflog.Debug(ctx, fmt.Sprintf("Found %d site lists after filtering", len(filteredItems)))

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Datasource Read Method: %s_%s, found %d items", d.ProviderTypeName, d.TypeName, len(filteredItems)))
}
