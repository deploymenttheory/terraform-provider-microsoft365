package graphBetaWindowsUpdateCatalogItem

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/crud"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/errors"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// BUG NOTE: this is a bug in the Kiota middleware pipeline (its header struct isn’t safe for parallel use).
// This mutex is a workaround to serialize calls and avoid concurrent map writes.
var windowsUpdateCatalogItemsMu sync.Mutex

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

	// Fetch all catalog items (serialize to avoid concurrent map writes in Kiota headers)
	windowsUpdateCatalogItemsMu.Lock()
	respList, err := d.client.
		DeviceManagement().
		WindowsUpdateCatalogItems().
		Get(ctx, nil)
	windowsUpdateCatalogItemsMu.Unlock()

	if err != nil {
		errors.HandleGraphError(ctx, err, resp, "Read", d.ReadPermissions)
		return
	}

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

	var filteredItems []WindowsUpdateCatalogItemModel
	filterValue := object.FilterValue.ValueString()

	for _, item := range respList.GetValue() {
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

		case "release_date_time":
			itemReleaseDate := item.GetReleaseDateTime()
			if itemReleaseDate != nil && releaseDateTime != nil && itemReleaseDate.Equal(*releaseDateTime) {
				filteredItems = append(filteredItems, MapRemoteStateToDataSource(item))
			}

		case "end_of_support_date":
			itemEndOfSupportDate := item.GetEndOfSupportDate()
			if itemEndOfSupportDate != nil && endOfSupportDate != nil && itemEndOfSupportDate.Equal(*endOfSupportDate) {
				filteredItems = append(filteredItems, MapRemoteStateToDataSource(item))
			}
		}
	}

	object.Items = filteredItems

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Datasource Read Method: %s_%s, found %d items", d.ProviderTypeName, d.TypeName, len(filteredItems)))
}
