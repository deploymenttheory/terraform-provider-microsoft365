// read.go
package graphBetaRoleScopeTag

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

// Read handles the Read operation for Role Scope Tag data source.
func (d *RoleScopeTagDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var object RoleScopeTagDataSourceModel

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

	var filteredItems []RoleScopeTagModel
	filterValue := object.FilterValue.ValueString()

	// For ID filter, we can make a direct API call
	if filterType == "id" {

		respItem, err := d.client.
			DeviceManagement().
			RoleScopeTags().
			ByRoleScopeTagId(filterValue).
			Get(ctx, nil)

		if err != nil {
			errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationRead, d.ReadPermissions)
			return
		}

		filteredItems = append(filteredItems, MapRemoteStateToDataSource(respItem))
	} else {
		// For all other filters, we need to get all role scope tags and filter locally

		respList, err := d.client.
			DeviceManagement().
			RoleScopeTags().
			Get(ctx, nil)

		if err != nil {
			errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationRead, d.ReadPermissions)
			return
		}

		for _, item := range respList.GetValue() {
			switch filterType {
			case "all":
				filteredItems = append(filteredItems, MapRemoteStateToDataSource(item))

			case "display_name":
				if item.GetDisplayName() != nil && strings.Contains(
					strings.ToLower(*item.GetDisplayName()),
					strings.ToLower(filterValue)) {
					filteredItems = append(filteredItems, MapRemoteStateToDataSource(item))
				}
			}
		}
	}

	object.Items = filteredItems

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Datasource Read Method: %s, found %d items", DataSourceName, len(filteredItems)))
}
