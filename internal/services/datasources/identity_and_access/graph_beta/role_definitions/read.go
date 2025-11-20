package graphBetaRoleDefinitions

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
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/rolemanagement"
	graphcore "github.com/microsoftgraph/msgraph-sdk-go-core"
)

func (d *RoleDefinitionsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var object RoleDefinitionsDataSourceModel

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

	var filteredItems []RoleDefinitionModel
	filterValue := object.FilterValue.ValueString()

	switch filterType {
	case "id":
		roleDefinition, err := d.client.
			RoleManagement().
			Directory().
			RoleDefinitions().
			ByUnifiedRoleDefinitionId(filterValue).Get(ctx, nil)

		if err != nil {
			errors.HandleKiotaGraphError(ctx, err, resp, "Read", d.ReadPermissions)
			return
		}

		rdItem := MapRemoteStateToDataSource(ctx, roleDefinition)
		filteredItems = append(filteredItems, rdItem)

	case "odata":
		headers := abstractions.NewRequestHeaders()
		headers.Add("ConsistencyLevel", "eventual")

		requestParameters := d.buildODataRequestParameters(ctx, &object, headers)

		tflog.Debug(ctx, "Using Microsoft Graph SDK PageIterator for role definitions")

		allRoleDefinitions, err := d.getAllRoleDefinitionsWithPageIterator(ctx, requestParameters)
		if err != nil {
			tflog.Error(ctx, fmt.Sprintf("Error in OData query with pagination: %v", err))
			errors.HandleKiotaGraphError(ctx, err, resp, "Read", d.ReadPermissions)
			return
		}

		tflog.Debug(ctx, fmt.Sprintf("PageIterator returned %d results", len(allRoleDefinitions)))

		for _, roleDefinition := range allRoleDefinitions {
			rdItem := MapRemoteStateToDataSource(ctx, roleDefinition)
			filteredItems = append(filteredItems, rdItem)
		}

	default:
		tflog.Debug(ctx, "Using Microsoft Graph SDK PageIterator for role definitions (all/display_name filter)")

		allRoleDefinitions, err := d.getAllRoleDefinitionsWithPageIterator(ctx, nil)
		if err != nil {
			errors.HandleKiotaGraphError(ctx, err, resp, "Read", d.ReadPermissions)
			return
		}

		for _, roleDefinition := range allRoleDefinitions {
			rdItem := MapRemoteStateToDataSource(ctx, roleDefinition)

			switch filterType {
			case "all":
				filteredItems = append(filteredItems, rdItem)

			case "display_name":
				if roleDefinition.GetDisplayName() != nil && strings.Contains(
					strings.ToLower(*roleDefinition.GetDisplayName()),
					strings.ToLower(filterValue)) {
					filteredItems = append(filteredItems, rdItem)
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

func (d *RoleDefinitionsDataSource) buildODataRequestParameters(ctx context.Context, object *RoleDefinitionsDataSourceModel, headers *abstractions.RequestHeaders) *rolemanagement.DirectoryRoleDefinitionsRequestBuilderGetRequestConfiguration {
	requestParameters := &rolemanagement.DirectoryRoleDefinitionsRequestBuilderGetRequestConfiguration{
		Headers:         headers,
		QueryParameters: &rolemanagement.DirectoryRoleDefinitionsRequestBuilderGetQueryParameters{},
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

	return requestParameters
}

func (d *RoleDefinitionsDataSource) getAllRoleDefinitionsWithPageIterator(ctx context.Context, requestParameters *rolemanagement.DirectoryRoleDefinitionsRequestBuilderGetRequestConfiguration) ([]graphmodels.UnifiedRoleDefinitionable, error) {
	var allRoleDefinitions []graphmodels.UnifiedRoleDefinitionable

	roleDefinitionsResponse, err := d.client.
		RoleManagement().
		Directory().
		RoleDefinitions().
		Get(ctx, requestParameters)
	if err != nil {
		return nil, err
	}

	pageIterator, err := graphcore.NewPageIterator[graphmodels.UnifiedRoleDefinitionable](
		roleDefinitionsResponse,
		d.client.GetAdapter(),
		graphmodels.CreateUnifiedRoleDefinitionCollectionResponseFromDiscriminatorValue,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create page iterator: %w", err)
	}

	pageCount := 0
	err = pageIterator.Iterate(ctx, func(item graphmodels.UnifiedRoleDefinitionable) bool {
		if item != nil {
			allRoleDefinitions = append(allRoleDefinitions, item)

			if len(allRoleDefinitions)%25 == 0 {
				pageCount++
				tflog.Debug(ctx, fmt.Sprintf("PageIterator: collected %d role definitions (estimated page %d)", len(allRoleDefinitions), pageCount))
			}
		}
		return true
	})

	if err != nil {
		return nil, fmt.Errorf("failed to iterate pages: %w", err)
	}

	tflog.Debug(ctx, fmt.Sprintf("PageIterator complete: collected %d total role definitions", len(allRoleDefinitions)))

	return allRoleDefinitions, nil
}
