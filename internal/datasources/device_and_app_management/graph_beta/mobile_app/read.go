package graphBetaMobileApp

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/crud"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/errors"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	abstractions "github.com/microsoft/kiota-abstractions-go"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/deviceappmanagement"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

func (d *MobileAppDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var object MobileAppDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	filterType := object.FilterType.ValueString()
	tflog.Debug(ctx, fmt.Sprintf("Starting Read method for datasource: %s with filter_type: %s", datasourceName, filterType))

	ctx, cancel := crud.HandleTimeout(ctx, object.Timeouts.Read, ReadTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	var filteredItems []MobileAppModel
	filterValue := object.FilterValue.ValueString()
	appTypeFilter := object.AppTypeFilter.ValueString()

	if filterType == "id" {
		requestParameters := &deviceappmanagement.MobileAppsMobileAppItemRequestBuilderGetRequestConfiguration{
			QueryParameters: &deviceappmanagement.MobileAppsMobileAppItemRequestBuilderGetQueryParameters{
				Expand: []string{"categories"},
			},
		}

		mobileApp, err := d.client.
			DeviceAppManagement().
			MobileApps().
			ByMobileAppId(filterValue).
			Get(ctx, requestParameters)

		if err != nil {
			errors.HandleGraphError(ctx, err, resp, "Read", d.ReadPermissions)
			return
		}

		if appTypeFilter != "" {
			appType := getAppTypeFromMobileApp(mobileApp)
			if appType != appTypeFilter {
				tflog.Debug(ctx, fmt.Sprintf("App with ID %s is of type %s, which doesn't match the requested filter %s",
					filterValue, appType, appTypeFilter))

				object.Items = []MobileAppModel{}
				resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
				return
			}
		}
		appItem := MapRemoteStateToDataSource(ctx, mobileApp)
		filteredItems = append(filteredItems, appItem)
	} else if filterType == "odata" {
		// Add "ConsistencyLevel: eventual" header for advanced OData queries
		headers := abstractions.NewRequestHeaders()
		headers.Add("ConsistencyLevel", "eventual")

		// Initialize request parameters with expand
		requestParameters := &deviceappmanagement.MobileAppsRequestBuilderGetRequestConfiguration{
			Headers: headers,
			QueryParameters: &deviceappmanagement.MobileAppsRequestBuilderGetQueryParameters{
				Expand: []string{"categories"},
			},
		}

		if !object.ODataFilter.IsNull() && object.ODataFilter.ValueString() != "" {
			filter := object.ODataFilter.ValueString()
			requestParameters.QueryParameters.Filter = &filter
			tflog.Debug(ctx, fmt.Sprintf("Setting OData filter: %s", filter))
		}

		if !object.ODataTop.IsNull() {
			topValue := object.ODataTop.ValueInt32()
			requestParameters.QueryParameters.Top = &topValue
			tflog.Debug(ctx, fmt.Sprintf("Setting OData top: %d", topValue))
		}

		if !object.ODataSkip.IsNull() {
			skipValue := object.ODataSkip.ValueInt32()
			requestParameters.QueryParameters.Skip = &skipValue
			tflog.Debug(ctx, fmt.Sprintf("Setting OData skip: %d", skipValue))
		}

		if !object.ODataSelect.IsNull() && object.ODataSelect.ValueString() != "" {
			selectFields := strings.Split(object.ODataSelect.ValueString(), ",")
			requestParameters.QueryParameters.Select = selectFields
			tflog.Debug(ctx, fmt.Sprintf("Setting OData select: %v", selectFields))
		}

		if !object.ODataOrderBy.IsNull() && object.ODataOrderBy.ValueString() != "" {
			orderbyFields := strings.Split(object.ODataOrderBy.ValueString(), ",")
			requestParameters.QueryParameters.Orderby = orderbyFields
			tflog.Debug(ctx, fmt.Sprintf("Setting OData orderby: %v", orderbyFields))
		}

		respList, err := d.client.
			DeviceAppManagement().
			MobileApps().
			Get(ctx, requestParameters)

		if err != nil {
			tflog.Error(ctx, fmt.Sprintf("Error in OData query: %v", err))
			errors.HandleGraphError(ctx, err, resp, "Read", d.ReadPermissions)
			return
		}

		tflog.Debug(ctx, fmt.Sprintf("API returned %d results", len(respList.GetValue())))

		for _, mobileApp := range respList.GetValue() {
			appItem := MapRemoteStateToDataSource(ctx, mobileApp)
			filteredItems = append(filteredItems, appItem)
		}
	} else {
		// For "all" and "display_name", get the full list and filter locally
		respList, err := d.client.
			DeviceAppManagement().
			MobileApps().
			Get(ctx, nil)

		if err != nil {
			errors.HandleGraphError(ctx, err, resp, "Read", d.ReadPermissions)
			return
		}

		for _, app := range respList.GetValue() {
			mobileApp, ok := app.(graphmodels.MobileAppable)
			if !ok {
				continue
			}

			currentAppType := getAppTypeFromMobileApp(mobileApp)

			if appTypeFilter != "" && currentAppType != appTypeFilter {
				continue
			}

			appItem := MapRemoteStateToDataSource(ctx, mobileApp)

			switch filterType {
			case "all":
				filteredItems = append(filteredItems, appItem)

			case "display_name":
				if mobileApp.GetDisplayName() != nil && strings.Contains(
					strings.ToLower(*mobileApp.GetDisplayName()),
					strings.ToLower(filterValue)) {
					filteredItems = append(filteredItems, appItem)
				}
			}
		}
	}

	object.Items = filteredItems

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Datasource Read Method: %s, found %d items", datasourceName, len(filteredItems)))
}
