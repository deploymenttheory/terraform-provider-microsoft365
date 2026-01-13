package graphBetaMobileAppRelationship

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
	"github.com/microsoftgraph/msgraph-beta-sdk-go/deviceappmanagement"
)

func (d *MobileAppRelationshipDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var object MobileAppRelationshipDataSourceModel

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

	var filteredItems []MobileAppRelationshipModel
	filterValue := object.FilterValue.ValueString()

	if filterType == "id" {
		// Get a specific mobile app relationship by ID
		mobileAppRelationship, err := d.client.
			DeviceAppManagement().
			MobileAppRelationships().
			ByMobileAppRelationshipId(filterValue).
			Get(ctx, nil)

		if err != nil {
			errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationRead, d.ReadPermissions)
			return
		}

		appItem := MapRemoteStateToDataSource(ctx, mobileAppRelationship)
		filteredItems = append(filteredItems, appItem)
	} else if filterType == "odata" {
		// Add "ConsistencyLevel: eventual" header for advanced OData queries
		headers := abstractions.NewRequestHeaders()
		headers.Add("ConsistencyLevel", "eventual")

		// Initialize request parameters
		requestParameters := &deviceappmanagement.MobileAppRelationshipsRequestBuilderGetRequestConfiguration{
			Headers: headers,
		}

		if !object.ODataFilter.IsNull() && object.ODataFilter.ValueString() != "" {
			filter := object.ODataFilter.ValueString()
			requestParameters.QueryParameters = &deviceappmanagement.MobileAppRelationshipsRequestBuilderGetQueryParameters{
				Filter: &filter,
			}
			tflog.Debug(ctx, fmt.Sprintf("Setting OData filter: %s", filter))
		}

		if !object.ODataTop.IsNull() {
			topValue := object.ODataTop.ValueInt32()
			if requestParameters.QueryParameters == nil {
				requestParameters.QueryParameters = &deviceappmanagement.MobileAppRelationshipsRequestBuilderGetQueryParameters{}
			}
			requestParameters.QueryParameters.Top = &topValue
			tflog.Debug(ctx, fmt.Sprintf("Setting OData top: %d", topValue))
		}

		if !object.ODataSkip.IsNull() {
			skipValue := object.ODataSkip.ValueInt32()
			if requestParameters.QueryParameters == nil {
				requestParameters.QueryParameters = &deviceappmanagement.MobileAppRelationshipsRequestBuilderGetQueryParameters{}
			}
			requestParameters.QueryParameters.Skip = &skipValue
			tflog.Debug(ctx, fmt.Sprintf("Setting OData skip: %d", skipValue))
		}

		if !object.ODataSelect.IsNull() && object.ODataSelect.ValueString() != "" {
			selectFields := strings.Split(object.ODataSelect.ValueString(), ",")
			if requestParameters.QueryParameters == nil {
				requestParameters.QueryParameters = &deviceappmanagement.MobileAppRelationshipsRequestBuilderGetQueryParameters{}
			}
			requestParameters.QueryParameters.Select = selectFields
			tflog.Debug(ctx, fmt.Sprintf("Setting OData select: %v", selectFields))
		}

		if !object.ODataOrderBy.IsNull() && object.ODataOrderBy.ValueString() != "" {
			orderbyFields := strings.Split(object.ODataOrderBy.ValueString(), ",")
			if requestParameters.QueryParameters == nil {
				requestParameters.QueryParameters = &deviceappmanagement.MobileAppRelationshipsRequestBuilderGetQueryParameters{}
			}
			requestParameters.QueryParameters.Orderby = orderbyFields
			tflog.Debug(ctx, fmt.Sprintf("Setting OData orderby: %v", orderbyFields))
		}

		respList, err := d.client.
			DeviceAppManagement().
			MobileAppRelationships().
			Get(ctx, requestParameters)

		if err != nil {
			tflog.Error(ctx, fmt.Sprintf("Error in OData query: %v", err))
			errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationRead, d.ReadPermissions)
			return
		}

		tflog.Debug(ctx, fmt.Sprintf("API returned %d results", len(respList.GetValue())))

		for _, relationship := range respList.GetValue() {
			appItem := MapRemoteStateToDataSource(ctx, relationship)
			filteredItems = append(filteredItems, appItem)
		}
	} else if filterType == "source_id" || filterType == "target_id" {
		// Get all relationships and filter by source_id or target_id
		respList, err := d.client.
			DeviceAppManagement().
			MobileAppRelationships().
			Get(ctx, nil)

		if err != nil {
			errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationRead, d.ReadPermissions)
			return
		}

		for _, relationship := range respList.GetValue() {
			if filterType == "source_id" && relationship.GetSourceId() != nil && *relationship.GetSourceId() == filterValue {
				appItem := MapRemoteStateToDataSource(ctx, relationship)
				filteredItems = append(filteredItems, appItem)
			} else if filterType == "target_id" && relationship.GetTargetId() != nil && *relationship.GetTargetId() == filterValue {
				appItem := MapRemoteStateToDataSource(ctx, relationship)
				filteredItems = append(filteredItems, appItem)
			}
		}
	} else if filterType == "all" {
		// Get all relationships
		respList, err := d.client.
			DeviceAppManagement().
			MobileAppRelationships().
			Get(ctx, nil)

		if err != nil {
			errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationRead, d.ReadPermissions)
			return
		}

		for _, relationship := range respList.GetValue() {
			appItem := MapRemoteStateToDataSource(ctx, relationship)
			filteredItems = append(filteredItems, appItem)
		}
	}

	object.Items = filteredItems

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Datasource Read Method: %s, found %d items", DataSourceName, len(filteredItems)))
}
