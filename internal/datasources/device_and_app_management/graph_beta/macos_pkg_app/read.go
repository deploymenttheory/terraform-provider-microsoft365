package graphBetaMacOSPKGApp

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/crud"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/errors"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/deviceappmanagement"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// Read handles the Read operation for macOS PKG App data source.
func (d *MacOSPKGAppDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var object MacOSPKGAppDataSourceModel

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

	var filteredItems []MacOSPKGAppModel
	filterValue := object.FilterValue.ValueString()

	// For ID filter, we can make a direct API call
	if filterType == "id" {
		requestParameters := &deviceappmanagement.MobileAppsMobileAppItemRequestBuilderGetRequestConfiguration{
			QueryParameters: &deviceappmanagement.MobileAppsMobileAppItemRequestBuilderGetQueryParameters{
				Expand: []string{"categories"},
			},
		}

		respBaseResource, err := d.client.
			DeviceAppManagement().
			MobileApps().
			ByMobileAppId(filterValue).
			Get(ctx, requestParameters)

		if err != nil {
			errors.HandleGraphError(ctx, err, resp, "Read", d.ReadPermissions)
			return
		}

		macOSPkgApp, ok := respBaseResource.(graphmodels.MacOSPkgAppable)
		if !ok {
			resp.Diagnostics.AddError(
				"Resource type mismatch",
				fmt.Sprintf("Expected resource of type MacOSPkgAppable but got %T", respBaseResource),
			)
			return
		}

		filteredItems = append(filteredItems, MapRemoteStateToDataSource(ctx, macOSPkgApp))
	} else {
		// For all other filters, we need to get all macOS PKG apps and filter locally

		respList, err := d.client.
			DeviceAppManagement().
			MobileApps().
			Get(ctx, nil)

		if err != nil {
			errors.HandleGraphError(ctx, err, resp, "Read", d.ReadPermissions)
			return
		}

		for _, app := range respList.GetValue() {
			// Skip if not a macOS PKG app
			macOSPkgApp, ok := app.(graphmodels.MacOSPkgAppable)
			if !ok {
				continue
			}

			switch filterType {
			case "all":
				filteredItems = append(filteredItems, MapRemoteStateToDataSource(ctx, macOSPkgApp))

			case "display_name":
				if macOSPkgApp.GetDisplayName() != nil && strings.Contains(
					strings.ToLower(*macOSPkgApp.GetDisplayName()),
					strings.ToLower(filterValue)) {
					filteredItems = append(filteredItems, MapRemoteStateToDataSource(ctx, macOSPkgApp))
				}
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
