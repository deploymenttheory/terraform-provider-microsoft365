package graphBetaWindowsUpdateCatalog

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud"
	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Read handles the Read operation for Windows Update Catalog data source.
func (d *WindowsUpdateCatalogEnteriesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var object WindowsUpdateCatalogEnteriesDataSourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting Read method for: %s", DataSourceName))

	resp.Diagnostics.Append(req.Config.Get(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	filterType := object.FilterType.ValueString()
	tflog.Debug(ctx, fmt.Sprintf("Reading %s with filter_type: %s", DataSourceName, filterType))

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

	// Fetch all catalog entries from the Windows Updates service
	respList, err := d.client.
		Admin().
		Windows().
		Updates().
		Catalog().
		Entries().
		Get(ctx, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationRead, d.ReadPermissions)
		return
	}

	var filteredEntries []WindowsUpdateCatalogEntry
	filterValue := object.FilterValue.ValueString()

	for _, entry := range respList.GetValue() {
		switch filterType {
		case "all":
			filteredEntries = append(filteredEntries, MapRemoteStateToDataSource(entry))

		case "id":
			if entry.GetId() != nil && *entry.GetId() == filterValue {
				filteredEntries = append(filteredEntries, MapRemoteStateToDataSource(entry))
			}

		case "display_name":
			if entry.GetDisplayName() != nil && strings.Contains(
				strings.ToLower(*entry.GetDisplayName()),
				strings.ToLower(filterValue)) {
				filteredEntries = append(filteredEntries, MapRemoteStateToDataSource(entry))
			}

		case "catalog_entry_type":
			mappedEntry := MapRemoteStateToDataSource(entry)
			if mappedEntry.CatalogEntryType.ValueString() == filterValue {
				filteredEntries = append(filteredEntries, mappedEntry)
			}
		}
	}

	object.Entries = filteredEntries

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Datasource Read Method: %s, found %d entries", DataSourceName, len(filteredEntries)))
}
