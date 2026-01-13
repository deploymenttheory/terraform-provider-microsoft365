package graphBetaServicePrincipal

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
	abstractions "github.com/microsoft/kiota-abstractions-go"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/serviceprincipals"
	graphcore "github.com/microsoftgraph/msgraph-sdk-go-core"
)

func (d *ServicePrincipalDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var object ServicePrincipalDataSourceModel

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

	var filteredItems []ServicePrincipalModel
	filterValue := object.FilterValue.ValueString()

	switch filterType {
	case "id":
		servicePrincipal, err := d.client.
			ServicePrincipals().
			ByServicePrincipalId(filterValue).
			Get(ctx, nil)

		if err != nil {
			errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationRead, d.ReadPermissions)
			return
		}

		spItem := MapRemoteStateToDataSource(ctx, servicePrincipal)
		filteredItems = append(filteredItems, spItem)

	case "app_id":
		filter := fmt.Sprintf("appId eq '%s'", filterValue)
		requestParameters := &serviceprincipals.ServicePrincipalsRequestBuilderGetRequestConfiguration{
			QueryParameters: &serviceprincipals.ServicePrincipalsRequestBuilderGetQueryParameters{
				Filter: &filter,
			},
		}

		// Use PageIterator for consistency, even though app_id filter typically returns few results
		tflog.Debug(ctx, "Using Microsoft Graph SDK PageIterator for service principals (app_id filter)")

		allServicePrincipals, err := d.getAllServicePrincipalsWithPageIterator(ctx, requestParameters)
		if err != nil {
			errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationRead, d.ReadPermissions)
			return
		}

		for _, servicePrincipal := range allServicePrincipals {
			spItem := MapRemoteStateToDataSource(ctx, servicePrincipal)
			filteredItems = append(filteredItems, spItem)
		}

	case "odata":
		// Add "ConsistencyLevel: eventual" header for advanced OData queries
		headers := abstractions.NewRequestHeaders()
		headers.Add("ConsistencyLevel", "eventual")

		requestParameters := &serviceprincipals.ServicePrincipalsRequestBuilderGetRequestConfiguration{
			Headers:         headers,
			QueryParameters: &serviceprincipals.ServicePrincipalsRequestBuilderGetQueryParameters{},
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

		tflog.Debug(ctx, "Using Microsoft Graph SDK PageIterator for service principals")

		allServicePrincipals, err := d.getAllServicePrincipalsWithPageIterator(ctx, requestParameters)
		if err != nil {
			tflog.Error(ctx, fmt.Sprintf("Error in OData query with pagination: %v", err))
			errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationRead, d.ReadPermissions)
			return
		}

		tflog.Debug(ctx, fmt.Sprintf("PageIterator returned %d results", len(allServicePrincipals)))

		for _, servicePrincipal := range allServicePrincipals {
			spItem := MapRemoteStateToDataSource(ctx, servicePrincipal)
			filteredItems = append(filteredItems, spItem)
		}

	default:
		// For "all" and "display_name", get the full list and filter locally using page iterator
		tflog.Debug(ctx, "Using Microsoft Graph SDK PageIterator for service principals (all/display_name filter)")

		allServicePrincipals, err := d.getAllServicePrincipalsWithPageIterator(ctx, nil)
		if err != nil {
			errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationRead, d.ReadPermissions)
			return
		}

		for _, servicePrincipal := range allServicePrincipals {
			spItem := MapRemoteStateToDataSource(ctx, servicePrincipal)

			switch filterType {
			case "all":
				filteredItems = append(filteredItems, spItem)

			case "display_name":
				if servicePrincipal.GetDisplayName() != nil && strings.Contains(
					strings.ToLower(*servicePrincipal.GetDisplayName()),
					strings.ToLower(filterValue)) {
					filteredItems = append(filteredItems, spItem)
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

// getAllServicePrincipalsWithPageIterator gets all service principals using page iterator for proper pagination
func (d *ServicePrincipalDataSource) getAllServicePrincipalsWithPageIterator(ctx context.Context, requestParameters *serviceprincipals.ServicePrincipalsRequestBuilderGetRequestConfiguration) ([]graphmodels.ServicePrincipalable, error) {
	var allServicePrincipals []graphmodels.ServicePrincipalable

	servicePrincipalsResponse, err := d.client.
		ServicePrincipals().
		Get(ctx, requestParameters)

	if err != nil {
		return nil, err
	}

	pageIterator, err := graphcore.NewPageIterator[graphmodels.ServicePrincipalable](
		servicePrincipalsResponse,
		d.client.GetAdapter(),
		graphmodels.CreateServicePrincipalCollectionResponseFromDiscriminatorValue,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create page iterator: %w", err)
	}

	pageCount := 0
	err = pageIterator.Iterate(ctx, func(item graphmodels.ServicePrincipalable) bool {
		if item != nil {
			allServicePrincipals = append(allServicePrincipals, item)

			// Log every 25 items (default page size)
			if len(allServicePrincipals)%25 == 0 {
				pageCount++
				tflog.Debug(ctx, fmt.Sprintf("PageIterator: collected %d service principals (estimated page %d)", len(allServicePrincipals), pageCount))
			}
		}
		return true
	})

	if err != nil {
		return nil, fmt.Errorf("failed to iterate pages: %w", err)
	}

	tflog.Debug(ctx, fmt.Sprintf("PageIterator complete: collected %d total service principals", len(allServicePrincipals)))

	return allServicePrincipals, nil
}
