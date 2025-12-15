package graphBetaMobileApp

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud"
	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	abstractions "github.com/microsoft/kiota-abstractions-go"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/deviceappmanagement"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
	graphcore "github.com/microsoftgraph/msgraph-sdk-go-core"
)

func (d *MobileAppDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var object MobileAppDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	filterType := object.FilterType.ValueString()
	tflog.Debug(ctx, fmt.Sprintf("Starting Read method for datasource: %s with filter_type: %s", DataSourceName, filterType))

	ctx, cancel := crud.HandleTimeout(ctx, object.Timeouts.Read, ReadTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	var filteredItems []MobileAppModel
	filterValue := object.FilterValue.ValueString()
	appTypeFilter := object.AppTypeFilter.ValueString()

	switch filterType {
	case "id":
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
			errors.HandleKiotaGraphError(ctx, err, resp, "Read", d.ReadPermissions)
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
	case "odata":
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
			for i, field := range selectFields {
				selectFields[i] = strings.TrimSpace(field)
			}
			requestParameters.QueryParameters.Select = selectFields
			tflog.Debug(ctx, fmt.Sprintf("Setting OData select: %v", selectFields))
		}

		if !object.ODataOrderBy.IsNull() && object.ODataOrderBy.ValueString() != "" {
			orderbyFields := strings.Split(object.ODataOrderBy.ValueString(), ",")
			for i, field := range orderbyFields {
				orderbyFields[i] = strings.TrimSpace(field)
			}
			requestParameters.QueryParameters.Orderby = orderbyFields
			tflog.Debug(ctx, fmt.Sprintf("Setting OData orderby: %v", orderbyFields))
		}

		tflog.Debug(ctx, "Using Microsoft Graph SDK PageIterator for mobile apps with OData parameters")

		allMobileApps, err := d.getAllMobileAppsWithPageIterator(ctx, requestParameters)
		if err != nil {
			tflog.Error(ctx, fmt.Sprintf("Error in OData query with pagination: %v", err))
			errors.HandleKiotaGraphError(ctx, err, resp, "Read", d.ReadPermissions)
			return
		}

		tflog.Debug(ctx, fmt.Sprintf("PageIterator returned %d results", len(allMobileApps)))

		for _, mobileApp := range allMobileApps {
			if appTypeFilter != "" {
				currentAppType := getAppTypeFromMobileApp(mobileApp)
				if currentAppType != appTypeFilter {
					continue
				}
			}

			appItem := MapRemoteStateToDataSource(ctx, mobileApp)
			filteredItems = append(filteredItems, appItem)
		}
	default:
		// For "all", "display_name", and "publisher_name", get the full list and filter locally using page iterator
		tflog.Debug(ctx, "Using Microsoft Graph SDK PageIterator for mobile apps (all/display_name/publisher_name filter)")

		// Build request parameters with expand for categories
		requestParameters := &deviceappmanagement.MobileAppsRequestBuilderGetRequestConfiguration{
			QueryParameters: &deviceappmanagement.MobileAppsRequestBuilderGetQueryParameters{
				Expand: []string{"categories"},
			},
		}

		allMobileApps, err := d.getAllMobileAppsWithPageIterator(ctx, requestParameters)
		if err != nil {
			errors.HandleKiotaGraphError(ctx, err, resp, "Read", d.ReadPermissions)
			return
		}

		for _, mobileApp := range allMobileApps {
			currentAppType := getAppTypeFromMobileApp(mobileApp)

			if appTypeFilter != "" && currentAppType != appTypeFilter {
				continue
			}

			shouldInclude := false

			switch filterType {
			case "all":
				shouldInclude = true

			case "display_name":
				if mobileApp.GetDisplayName() != nil && strings.Contains(
					strings.ToLower(*mobileApp.GetDisplayName()),
					strings.ToLower(filterValue)) {
					shouldInclude = true
				}

			case "publisher_name":
				if mobileApp.GetPublisher() != nil && strings.Contains(
					strings.ToLower(*mobileApp.GetPublisher()),
					strings.ToLower(filterValue)) {
					shouldInclude = true
				}
			}

			if shouldInclude {
				appItem := MapRemoteStateToDataSource(ctx, mobileApp)
				filteredItems = append(filteredItems, appItem)
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

// getAllMobileAppsWithPageIterator gets all mobile apps using page iterator for proper pagination
func (d *MobileAppDataSource) getAllMobileAppsWithPageIterator(ctx context.Context, requestParameters *deviceappmanagement.MobileAppsRequestBuilderGetRequestConfiguration) ([]graphmodels.MobileAppable, error) {
	var allApps []graphmodels.MobileAppable

	appsResponse, err := d.client.
		DeviceAppManagement().
		MobileApps().
		Get(ctx, requestParameters)

	if err != nil {
		return nil, err
	}

	pageIterator, err := graphcore.NewPageIterator[graphmodels.MobileAppable](
		appsResponse,
		d.client.GetAdapter(),
		graphmodels.CreateMobileAppCollectionResponseFromDiscriminatorValue,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create page iterator: %w", err)
	}

	pageCount := 0
	err = pageIterator.Iterate(ctx, func(item graphmodels.MobileAppable) bool {
		if item != nil {
			allApps = append(allApps, item)

			// Log every 25 items (default page size)
			if len(allApps)%25 == 0 {
				pageCount++
				tflog.Debug(ctx, fmt.Sprintf("PageIterator: collected %d mobile apps (estimated page %d)", len(allApps), pageCount))
			}
		}
		return true
	})

	if err != nil {
		return nil, fmt.Errorf("failed to iterate pages: %w", err)
	}

	tflog.Debug(ctx, fmt.Sprintf("PageIterator complete: collected %d total mobile apps", len(allApps)))

	return allApps, nil
}
