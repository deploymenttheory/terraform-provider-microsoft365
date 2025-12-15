package graphBetaApplicationCategory

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud"
	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/deviceappmanagement"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
	graphcore "github.com/microsoftgraph/msgraph-sdk-go-core"
)

// getAllApplicationCategoriesWithPageIterator retrieves all application categories using PageIterator
func (d *ApplicationCategoryDataSource) getAllApplicationCategoriesWithPageIterator(ctx context.Context, requestConfig *deviceappmanagement.MobileAppCategoriesRequestBuilderGetRequestConfiguration) ([]graphmodels.MobileAppCategoryable, error) {
	var allCategories []graphmodels.MobileAppCategoryable

	categoriesResponse, err := d.client.DeviceAppManagement().MobileAppCategories().Get(ctx, requestConfig)
	if err != nil {
		return nil, err
	}

	pageIterator, err := graphcore.NewPageIterator[graphmodels.MobileAppCategoryable](
		categoriesResponse,
		d.client.GetAdapter(),
		graphmodels.CreateMobileAppCategoryCollectionResponseFromDiscriminatorValue,
	)
	if err != nil {
		return nil, fmt.Errorf("error creating page iterator: %v", err)
	}

	err = pageIterator.Iterate(ctx, func(category graphmodels.MobileAppCategoryable) bool {
		allCategories = append(allCategories, category)
		return true
	})
	if err != nil {
		return nil, fmt.Errorf("error iterating pages: %v", err)
	}

	tflog.Debug(ctx, fmt.Sprintf("PageIterator complete: collected %d total application categories", len(allCategories)))

	return allCategories, nil
}

// Read handles the Read operation for Application Categories data source.
func (d *ApplicationCategoryDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var object ApplicationCategoryDataSourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting Read method for: %s", DataSourceName))

	resp.Diagnostics.Append(req.Config.Get(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	filterType := object.FilterType.ValueString()
	tflog.Debug(ctx, fmt.Sprintf("Reading %s with filter_type: %s", DataSourceName, filterType))

	if (filterType == "id" || filterType == "display_name") && (object.FilterValue.IsNull() || object.FilterValue.ValueString() == "") {
		resp.Diagnostics.AddError(
			"Missing Required Parameter",
			fmt.Sprintf("filter_value must be provided when filter_type is '%s'", filterType),
		)
		return
	}

	if filterType == "odata" && (object.ODataFilter.IsNull() || object.ODataFilter.ValueString() == "") {
		resp.Diagnostics.AddError(
			"Missing Required Parameter",
			"odata_filter must be provided when filter_type is 'odata'",
		)
		return
	}

	ctx, cancel := crud.HandleTimeout(ctx, object.Timeouts.Read, ReadTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	var filteredItems []ApplicationCategoryModel
	filterValue := object.FilterValue.ValueString()

	switch filterType {
	case "id":
		tflog.Debug(ctx, fmt.Sprintf("Fetching application category by ID: %s", filterValue))

		respItem, err := d.client.
			DeviceAppManagement().
			MobileAppCategories().
			ByMobileAppCategoryId(filterValue).
			Get(ctx, nil)

		if err != nil {
			errors.HandleKiotaGraphError(ctx, err, resp, "Read", d.ReadPermissions)
			return
		}

		filteredItems = append(filteredItems, MapRemoteStateToDataSource(respItem))

	case "all":
		tflog.Debug(ctx, "Fetching all application categories")

		allCategories, err := d.getAllApplicationCategoriesWithPageIterator(ctx, nil)
		if err != nil {
			errors.HandleKiotaGraphError(ctx, err, resp, "Read", d.ReadPermissions)
			return
		}

		for _, category := range allCategories {
			filteredItems = append(filteredItems, MapRemoteStateToDataSource(category))
		}

	case "display_name":
		tflog.Debug(ctx, fmt.Sprintf("Fetching application categories by display name: %s", filterValue))

		allCategories, err := d.getAllApplicationCategoriesWithPageIterator(ctx, nil)
		if err != nil {
			errors.HandleKiotaGraphError(ctx, err, resp, "Read", d.ReadPermissions)
			return
		}

		lowerFilterValue := strings.ToLower(filterValue)
		for _, category := range allCategories {
			if category.GetDisplayName() != nil {
				lowerDisplayName := strings.ToLower(*category.GetDisplayName())
				if strings.Contains(lowerDisplayName, lowerFilterValue) {
					filteredItems = append(filteredItems, MapRemoteStateToDataSource(category))
				}
			}
		}

	case "odata":
		odataFilter := object.ODataFilter.ValueString()
		tflog.Debug(ctx, fmt.Sprintf("Fetching application categories with OData filter: %s", odataFilter))

		requestConfig := &deviceappmanagement.MobileAppCategoriesRequestBuilderGetRequestConfiguration{
			QueryParameters: &deviceappmanagement.MobileAppCategoriesRequestBuilderGetQueryParameters{
				Filter: &odataFilter,
			},
		}

		if !object.ODataTop.IsNull() {
			top := object.ODataTop.ValueInt32()
			requestConfig.QueryParameters.Top = &top
		}

		allCategories, err := d.getAllApplicationCategoriesWithPageIterator(ctx, requestConfig)
		if err != nil {
			errors.HandleKiotaGraphError(ctx, err, resp, "Read", d.ReadPermissions)
			return
		}

		for _, category := range allCategories {
			filteredItems = append(filteredItems, MapRemoteStateToDataSource(category))
		}

	default:
		resp.Diagnostics.AddError(
			"Invalid Filter Type",
			fmt.Sprintf("Unsupported filter_type: %s. Valid values are: all, id, display_name, odata", filterType),
		)
		return
	}

	object.Items = filteredItems

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Datasource Read Method: %s, found %d items", DataSourceName, len(filteredItems)))
}
