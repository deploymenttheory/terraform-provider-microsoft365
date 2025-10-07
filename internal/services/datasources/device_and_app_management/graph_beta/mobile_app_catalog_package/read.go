package graphBetaMobileAppCatalogPackage

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

func (d *MobileAppCatalogPackageDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var object MobileAppCatalogPackageDataSourceModel

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

	var filteredItems []MobileAppCatalogPackageModel
	filterValue := object.FilterValue.ValueString()

	switch filterType {
	case "id":
		// For mobile app catalog packages, "id" refers to product_id
		filter := fmt.Sprintf("productId eq '%s'", filterValue)
		requestParameters := &deviceappmanagement.MobileAppCatalogPackagesRequestBuilderGetRequestConfiguration{
			QueryParameters: &deviceappmanagement.MobileAppCatalogPackagesRequestBuilderGetQueryParameters{
				Filter: &filter,
			},
		}

		allPackages, err := d.getAllMobileAppCatalogPackageWithPageIterator(ctx, requestParameters)
		if err != nil {
			errors.HandleKiotaGraphError(ctx, err, resp, "Read", d.ReadPermissions)
			return
		}

		for _, packageItem := range allPackages {
			packageModel := MapRemoteStateToDataSource(ctx, packageItem)
			filteredItems = append(filteredItems, packageModel)
		}

	case "odata":
		// Add "ConsistencyLevel: eventual" header for advanced OData queries
		headers := abstractions.NewRequestHeaders()
		headers.Add("ConsistencyLevel", "eventual")

		requestParameters := &deviceappmanagement.MobileAppCatalogPackagesRequestBuilderGetRequestConfiguration{
			Headers:         headers,
			QueryParameters: &deviceappmanagement.MobileAppCatalogPackagesRequestBuilderGetQueryParameters{},
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

		if !object.ODataCount.IsNull() && object.ODataCount.ValueBool() {
			count := true
			requestParameters.QueryParameters.Count = &count
			tflog.Debug(ctx, "Setting OData count: true")
		}

		if !object.ODataSearch.IsNull() && object.ODataSearch.ValueString() != "" {
			search := object.ODataSearch.ValueString()
			requestParameters.QueryParameters.Search = &search
			tflog.Debug(ctx, fmt.Sprintf("Setting OData search: %s", search))
		}

		if !object.ODataExpand.IsNull() && object.ODataExpand.ValueString() != "" {
			expandFields := strings.Split(object.ODataExpand.ValueString(), ",")
			for i, field := range expandFields {
				expandFields[i] = strings.TrimSpace(field)
			}
			requestParameters.QueryParameters.Expand = expandFields
			tflog.Debug(ctx, fmt.Sprintf("Setting OData expand: %v", expandFields))
		}

		tflog.Debug(ctx, "Using Microsoft Graph SDK PageIterator for mobile app catalog packages")

		allPackages, err := d.getAllMobileAppCatalogPackageWithPageIterator(ctx, requestParameters)
		if err != nil {
			tflog.Error(ctx, fmt.Sprintf("Error in OData query with pagination: %v", err))
			errors.HandleKiotaGraphError(ctx, err, resp, "Read", d.ReadPermissions)
			return
		}

		tflog.Debug(ctx, fmt.Sprintf("PageIterator returned %d results", len(allPackages)))

		for _, packageItem := range allPackages {
			packageModel := MapRemoteStateToDataSource(ctx, packageItem)
			filteredItems = append(filteredItems, packageModel)
		}

	default:
		// For "all", "product_name", and "publisher_name", get the full list and filter locally using page iterator
		tflog.Debug(ctx, "Using Microsoft Graph SDK PageIterator for mobile app catalog packages (all/product_name/publisher_name filter)")

		allPackages, err := d.getAllMobileAppCatalogPackageWithPageIterator(ctx, nil)
		if err != nil {
			errors.HandleKiotaGraphError(ctx, err, resp, "Read", d.ReadPermissions)
			return
		}

		for _, packageItem := range allPackages {
			packageModel := MapRemoteStateToDataSource(ctx, packageItem)

			switch filterType {
			case "all":
				filteredItems = append(filteredItems, packageModel)

			case "product_name":
				if packageItem.GetProductDisplayName() != nil && strings.Contains(
					strings.ToLower(*packageItem.GetProductDisplayName()),
					strings.ToLower(filterValue)) {
					filteredItems = append(filteredItems, packageModel)
				}

			case "publisher_name":
				if packageItem.GetPublisherDisplayName() != nil && strings.Contains(
					strings.ToLower(*packageItem.GetPublisherDisplayName()),
					strings.ToLower(filterValue)) {
					filteredItems = append(filteredItems, packageModel)
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

// getAllMobileAppCatalogPackageWithPageIterator gets all mobile app catalog packages using page iterator for proper pagination
func (d *MobileAppCatalogPackageDataSource) getAllMobileAppCatalogPackageWithPageIterator(ctx context.Context, requestParameters *deviceappmanagement.MobileAppCatalogPackagesRequestBuilderGetRequestConfiguration) ([]graphmodels.MobileAppCatalogPackageable, error) {
	var allPackages []graphmodels.MobileAppCatalogPackageable

	packagesResponse, err := d.client.
		DeviceAppManagement().
		MobileAppCatalogPackages().
		Get(ctx, requestParameters)

	if err != nil {
		return nil, err
	}

	pageIterator, err := graphcore.NewPageIterator[graphmodels.MobileAppCatalogPackageable](
		packagesResponse,
		d.client.GetAdapter(),
		graphmodels.CreateMobileAppCatalogPackageCollectionResponseFromDiscriminatorValue,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create page iterator: %w", err)
	}

	pageCount := 0
	err = pageIterator.Iterate(ctx, func(item graphmodels.MobileAppCatalogPackageable) bool {
		if item != nil {
			allPackages = append(allPackages, item)

			// Log every 25 items (default page size)
			if len(allPackages)%25 == 0 {
				pageCount++
				tflog.Debug(ctx, fmt.Sprintf("PageIterator: collected %d mobile app catalog packages (estimated page %d)", len(allPackages), pageCount))
			}
		}
		return true
	})

	if err != nil {
		return nil, fmt.Errorf("failed to iterate pages: %w", err)
	}

	tflog.Debug(ctx, fmt.Sprintf("PageIterator complete: collected %d total mobile app catalog packages", len(allPackages)))

	return allPackages, nil
}
