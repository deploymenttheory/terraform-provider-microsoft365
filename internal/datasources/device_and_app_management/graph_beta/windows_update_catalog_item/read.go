package graphBetaWindowsUpdateCatalogItem

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

// Read handles the Read operation for Windows Update Catalog Items data source.
func (d *WindowsUpdateCatalogItemDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var object WindowsUpdateCatalogItemDataSourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting Read method for: %s_%s", d.ProviderTypeName, d.TypeName))

	resp.Diagnostics.Append(req.Config.Get(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	filterType := object.FilterType.ValueString()
	tflog.Debug(ctx, fmt.Sprintf("Reading %s_%s with filter_type: %s", d.ProviderTypeName, d.TypeName, filterType))

	// Validate filter_value is provided when filter_type is not "all"
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

	// Fetch all catalog items
	respList, err := d.client.
		DeviceManagement().
		WindowsUpdateCatalogItems().
		Get(ctx, nil)

	if err != nil {
		errors.HandleGraphError(ctx, err, resp, "Read", d.ReadPermissions)
		return
	}

	// Parse date filter if necessary
	var releaseDateTime, endOfSupportDate *time.Time

	if filterType == "release_date_time" || filterType == "end_of_support_date" {
		parsedTime, err := time.Parse(time.RFC3339, object.FilterValue.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Invalid Date Format",
				fmt.Sprintf("Could not parse date value as RFC3339 format: %s", err),
			)
			return
		}

		if filterType == "release_date_time" {
			releaseDateTime = &parsedTime
		} else {
			endOfSupportDate = &parsedTime
		}
	}

	// Filter the results based on the specified filter_type and filter_value
	var filteredItems []WindowsUpdateCatalogItemModel
	filterValue := object.FilterValue.ValueString()

	for _, item := range respList.GetValue() {
		switch filterType {
		case "all":
			// No filtering, include all items
			filteredItems = append(filteredItems, MapRemoteStateToDataSource(item))

		case "id":
			// Filter by ID (exact match)
			if item.GetId() != nil && *item.GetId() == filterValue {
				filteredItems = append(filteredItems, MapRemoteStateToDataSource(item))
			}

		case "display_name":
			// Filter by display name (case-insensitive substring match)
			if item.GetDisplayName() != nil && strings.Contains(
				strings.ToLower(*item.GetDisplayName()),
				strings.ToLower(filterValue)) {
				filteredItems = append(filteredItems, MapRemoteStateToDataSource(item))
			}

		case "release_date_time":
			// Filter by release date time (exact match)
			itemReleaseDate := item.GetReleaseDateTime()
			if itemReleaseDate != nil && releaseDateTime != nil && itemReleaseDate.Equal(*releaseDateTime) {
				filteredItems = append(filteredItems, MapRemoteStateToDataSource(item))
			}

		case "end_of_support_date":
			// Filter by end of support date (exact match)
			itemEndOfSupportDate := item.GetEndOfSupportDate()
			if itemEndOfSupportDate != nil && endOfSupportDate != nil && itemEndOfSupportDate.Equal(*endOfSupportDate) {
				filteredItems = append(filteredItems, MapRemoteStateToDataSource(item))
			}
		}
	}

	// Update the model with the filtered items
	object.Items = filteredItems

	// Set the data in the response
	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Datasource Read Method: %s_%s, found %d items", d.ProviderTypeName, d.TypeName, len(filteredItems)))
}
