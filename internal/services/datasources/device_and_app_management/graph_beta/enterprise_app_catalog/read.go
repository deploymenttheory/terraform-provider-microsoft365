package enterpriseappcatalog

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	abstractions "github.com/microsoft/kiota-abstractions-go"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/deviceappmanagement"
)

// Read fetches mobile app catalog packages and sets them in the data source state
func (d *EnterpriseAppCatalogDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var object MobileAppDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	filterType := object.FilterType.ValueString()
	tflog.Debug(ctx, fmt.Sprintf("Starting Read method for datasource: %s with filter_type: %s", datasourceName, filterType))

	// Set up timeout context
	ctx, cancel := crud.HandleTimeout(ctx, object.Timeouts.Read, ReadTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	// Add "ConsistencyLevel: eventual" header for advanced queries
	headers := abstractions.NewRequestHeaders()
	headers.Add("ConsistencyLevel", "eventual")

	var filteredItems []MobileAppCatalogPackageModel
	filterValue := object.FilterValue.ValueString()

	// Handle different filter types
	switch filterType {
	case "id":
		// Get a specific package by product ID
		requestConfiguration := &deviceappmanagement.MobileAppCatalogPackagesRequestBuilderGetRequestConfiguration{
			Headers: headers,
		}
		queryParams := &deviceappmanagement.MobileAppCatalogPackagesRequestBuilderGetQueryParameters{}

		filter := fmt.Sprintf("productId eq '%s'", filterValue)
		queryParams.Filter = &filter
		queryParams.Orderby = []string{"versionDisplayName asc"}

		requestConfiguration.QueryParameters = queryParams

		tflog.Debug(ctx, fmt.Sprintf("Filtering by product ID: %s", filterValue))

		response, err := d.client.
			DeviceAppManagement().
			MobileAppCatalogPackages().
			Get(ctx, requestConfiguration)

		if err != nil {
			errors.HandleGraphError(ctx, err, resp, "Read", d.ReadPermissions)
			return
		}

		if response != nil && response.GetValue() != nil {
			packages := response.GetValue()
			tflog.Debug(ctx, fmt.Sprintf("Retrieved %d mobile app catalog packages for product ID %s", len(packages), filterValue))

			terraformPackages, diags := MapRemoteStateToDataSource(ctx, packages)
			resp.Diagnostics.Append(diags...)
			if resp.Diagnostics.HasError() {
				return
			}
			filteredItems = terraformPackages
		}

	case "odata":
		// Use OData query parameters
		requestConfiguration := &deviceappmanagement.MobileAppCatalogPackagesRequestBuilderGetRequestConfiguration{
			Headers: headers,
		}
		queryParams := &deviceappmanagement.MobileAppCatalogPackagesRequestBuilderGetQueryParameters{}

		if !object.ODataFilter.IsNull() && object.ODataFilter.ValueString() != "" {
			filter := object.ODataFilter.ValueString()
			queryParams.Filter = &filter
			tflog.Debug(ctx, fmt.Sprintf("Setting OData filter: %s", filter))
		}

		if !object.ODataTop.IsNull() {
			topValue := object.ODataTop.ValueInt32()
			queryParams.Top = &topValue
			tflog.Debug(ctx, fmt.Sprintf("Setting OData top: %d", topValue))
		}

		if !object.ODataSkip.IsNull() {
			skipValue := object.ODataSkip.ValueInt32()
			queryParams.Skip = &skipValue
			tflog.Debug(ctx, fmt.Sprintf("Setting OData skip: %d", skipValue))
		}

		if !object.ODataSelect.IsNull() && object.ODataSelect.ValueString() != "" {
			selectFields := strings.Split(object.ODataSelect.ValueString(), ",")
			queryParams.Select = selectFields
			tflog.Debug(ctx, fmt.Sprintf("Setting OData select: %v", selectFields))
		}

		if !object.ODataOrderBy.IsNull() && object.ODataOrderBy.ValueString() != "" {
			orderbyFields := strings.Split(object.ODataOrderBy.ValueString(), ",")
			queryParams.Orderby = orderbyFields
			tflog.Debug(ctx, fmt.Sprintf("Setting OData orderby: %v", orderbyFields))
		}

		requestConfiguration.QueryParameters = queryParams

		response, err := d.client.
			DeviceAppManagement().
			MobileAppCatalogPackages().
			Get(ctx, requestConfiguration)

		if err != nil {
			errors.HandleGraphError(ctx, err, resp, "Read", d.ReadPermissions)
			return
		}

		if response != nil && response.GetValue() != nil {
			packages := response.GetValue()
			tflog.Debug(ctx, fmt.Sprintf("Retrieved %d mobile app catalog packages with OData query", len(packages)))

			terraformPackages, diags := MapRemoteStateToDataSource(ctx, packages)
			resp.Diagnostics.Append(diags...)
			if resp.Diagnostics.HasError() {
				return
			}
			filteredItems = terraformPackages
		}

	case "all":
		// Get all packages with groupby to get unique products
		requestConfiguration := &deviceappmanagement.MobileAppCatalogPackagesRequestBuilderGetRequestConfiguration{
			Headers: headers,
		}
		queryParams := &deviceappmanagement.MobileAppCatalogPackagesRequestBuilderGetQueryParameters{}

		// Use filter to group by product attributes
		filter := "isNull(versionDisplayName) eq false"
		queryParams.Filter = &filter
		queryParams.Orderby = []string{"productDisplayName asc"}

		if !object.ODataTop.IsNull() {
			topValue := object.ODataTop.ValueInt32()
			queryParams.Top = &topValue
		} else {
			topValue := int32(50) // Default to 50 if not specified
			queryParams.Top = &topValue
		}

		requestConfiguration.QueryParameters = queryParams

		tflog.Debug(ctx, "Getting all mobile app catalog packages (grouped by product)")

		response, err := d.client.
			DeviceAppManagement().
			MobileAppCatalogPackages().
			Get(ctx, requestConfiguration)

		if err != nil {
			errors.HandleGraphError(ctx, err, resp, "Read", d.ReadPermissions)
			return
		}

		if response != nil && response.GetValue() != nil {
			packages := response.GetValue()
			tflog.Debug(ctx, fmt.Sprintf("Retrieved %d mobile app catalog packages", len(packages)))

			terraformPackages, diags := MapRemoteStateToDataSource(ctx, packages)
			resp.Diagnostics.Append(diags...)
			if resp.Diagnostics.HasError() {
				return
			}
			filteredItems = terraformPackages
		}

	case "display_name":
		// Get all packages and filter by display name
		requestConfiguration := &deviceappmanagement.MobileAppCatalogPackagesRequestBuilderGetRequestConfiguration{
			Headers: headers,
		}
		queryParams := &deviceappmanagement.MobileAppCatalogPackagesRequestBuilderGetQueryParameters{}

		// Use filter to search by display name
		filter := fmt.Sprintf("contains(productDisplayName, '%s')", filterValue)
		queryParams.Filter = &filter
		queryParams.Orderby = []string{"productDisplayName asc"}

		if !object.ODataTop.IsNull() {
			topValue := object.ODataTop.ValueInt32()
			queryParams.Top = &topValue
		} else {
			topValue := int32(100) // Use a higher limit for filtering locally
			queryParams.Top = &topValue
		}

		requestConfiguration.QueryParameters = queryParams

		tflog.Debug(ctx, fmt.Sprintf("Searching for packages with display name containing: %s", filterValue))

		response, err := d.client.
			DeviceAppManagement().
			MobileAppCatalogPackages().
			Get(ctx, requestConfiguration)

		if err != nil {
			errors.HandleGraphError(ctx, err, resp, "Read", d.ReadPermissions)
			return
		}

		if response != nil && response.GetValue() != nil {
			packages := response.GetValue()
			tflog.Debug(ctx, fmt.Sprintf("Retrieved %d mobile app catalog packages, filtering by display name", len(packages)))

			// First convert all packages to Terraform models
			allPackages, diags := MapRemoteStateToDataSource(ctx, packages)
			resp.Diagnostics.Append(diags...)
			if resp.Diagnostics.HasError() {
				return
			}

			// Then filter by display name
			for _, pkg := range allPackages {
				if !pkg.ProductDisplayName.IsNull() && strings.Contains(
					strings.ToLower(pkg.ProductDisplayName.ValueString()),
					strings.ToLower(filterValue)) {
					filteredItems = append(filteredItems, pkg)
				}
			}

			tflog.Debug(ctx, fmt.Sprintf("Found %d packages matching display name filter", len(filteredItems)))
		}

	default:
		resp.Diagnostics.AddError(
			"Invalid Filter Type",
			fmt.Sprintf("Filter type '%s' is not supported. Supported types are: id, odata, all, display_name", filterType),
		)
		return
	}

	// If we got no results, return an empty list
	if filteredItems == nil {
		filteredItems = []MobileAppCatalogPackageModel{}
	}

	// Apply app type filter if provided
	if !object.AppTypeFilter.IsNull() && object.AppTypeFilter.ValueString() != "" {
		appTypeFilter := object.AppTypeFilter.ValueString()
		var appTypeFilteredItems []MobileAppCatalogPackageModel

		tflog.Debug(ctx, fmt.Sprintf("Filtering by app type: %s", appTypeFilter))

		for _, item := range filteredItems {
			// This would need to be implemented based on how you determine app types
			// For now, we'll just pass through all items since we don't have app type info
			appTypeFilteredItems = append(appTypeFilteredItems, item)
		}

		filteredItems = appTypeFilteredItems
	}

	// Fetch detailed app configuration if requested
	if !object.IncludeAppConfig.IsNull() && object.IncludeAppConfig.ValueBool() {
		tflog.Debug(ctx, fmt.Sprintf("Fetching detailed app configuration for %d items", len(filteredItems)))

		// Process each item to fetch its app configuration
		for i := range filteredItems {
			if filteredItems[i].Id.IsNull() || filteredItems[i].Id.ValueString() == "" {
				continue
			}

			packageId := filteredItems[i].Id.ValueString()
			appConfigData, diags := fetchAppConfigurationData(ctx, d.client.GetAdapter(), packageId)
			resp.Diagnostics.Append(diags...)
			if resp.Diagnostics.HasError() || appConfigData == nil {
				continue
			}

			// Map the response data to the model
			appConfig, diags := MapAppConfigToModel(ctx, appConfigData)
			resp.Diagnostics.Append(diags...)
			if resp.Diagnostics.HasError() {
				continue
			}

			filteredItems[i].AppConfig = appConfig
			tflog.Debug(ctx, fmt.Sprintf("Added app configuration for package ID %s", packageId))
		}
	}

	object.Items = filteredItems

	// Set the state
	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Datasource Read Method: %s, found %d items", datasourceName, len(filteredItems)))
}

// fetchAppConfigurationData retrieves detailed app configuration data for a package ID
func fetchAppConfigurationData(ctx context.Context, client abstractions.RequestAdapter, packageId string) (map[string]interface{}, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Construct the URL for the convertFromMobileAppCatalogPackage endpoint
	url := fmt.Sprintf("/deviceAppManagement/mobileApps/convertFromMobileAppCatalogPackage(mobileAppCatalogPackageId='%s')", packageId)

	// Create request information
	requestInfo := abstractions.NewRequestInformation()
	requestInfo.Method = http.MethodGet
	requestInfo.SetUri(url)
	requestInfo.Headers.Add("Accept", "application/json")

	tflog.Debug(ctx, "Fetching app configuration", map[string]interface{}{
		"packageId": packageId,
		"url":       url,
	})

	// Execute the request
	response, err := client.Send(ctx, requestInfo, nil, nil)
	if err != nil {
		diags.AddError(
			"Error Fetching App Configuration",
			fmt.Sprintf("Failed to fetch app configuration for package ID %s: %v", packageId, err),
		)
		return nil, diags
	}
	defer response.Close()

	// Check for successful response
	if response.StatusCode != http.StatusOK {
		diags.AddError(
			"Error Fetching App Configuration",
			fmt.Sprintf("Received non-OK status code %d when fetching app configuration for package ID %s", response.StatusCode, packageId),
		)
		return nil, diags
	}

	// Parse the response body
	var responseData map[string]interface{}
	err = json.NewDecoder(response).Decode(&responseData)
	if err != nil {
		diags.AddError(
			"Error Parsing App Configuration",
			fmt.Sprintf("Failed to parse app configuration response for package ID %s: %v", packageId, err),
		)
		return nil, diags
	}

	return responseData, diags
}
