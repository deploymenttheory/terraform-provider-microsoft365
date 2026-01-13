package graphBetaCloudPC

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
	"github.com/microsoftgraph/msgraph-beta-sdk-go/devicemanagement"
)

func (d *CloudPCDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var object CloudPCDataSourceModel

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

	var filteredItems []CloudPCItemDataSourceModel
	filterValue := object.FilterValue.ValueString()

	if filterType == "id" {
		requestParameters := &devicemanagement.VirtualEndpointCloudPCsCloudPCItemRequestBuilderGetRequestConfiguration{}

		cloudPC, err := d.client.
			DeviceManagement().
			VirtualEndpoint().
			CloudPCs().
			ByCloudPCId(filterValue).
			Get(ctx, requestParameters)

		if err != nil {
			errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationRead, d.ReadPermissions)
			return
		}

		cloudPCItem := MapRemoteStateToDataSource(ctx, cloudPC)
		filteredItems = append(filteredItems, cloudPCItem)

	} else if filterType == "odata" {
		// Add "ConsistencyLevel: eventual" header for advanced OData queries
		headers := abstractions.NewRequestHeaders()
		headers.Add("ConsistencyLevel", "eventual")

		// Initialize request parameters
		requestParameters := &devicemanagement.VirtualEndpointCloudPCsRequestBuilderGetRequestConfiguration{
			Headers:         headers,
			QueryParameters: &devicemanagement.VirtualEndpointCloudPCsRequestBuilderGetQueryParameters{},
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
			DeviceManagement().
			VirtualEndpoint().
			CloudPCs().
			Get(ctx, requestParameters)

		if err != nil {
			tflog.Error(ctx, fmt.Sprintf("Error in OData query: %v", err))
			errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationRead, d.ReadPermissions)
			return
		}

		tflog.Debug(ctx, fmt.Sprintf("API returned %d results", len(respList.GetValue())))

		for _, cloudPC := range respList.GetValue() {
			cloudPCItem := MapRemoteStateToDataSource(ctx, cloudPC)
			filteredItems = append(filteredItems, cloudPCItem)
		}

	} else {
		// For "all" and "display_name", get the full list and filter locally
		respList, err := d.client.
			DeviceManagement().
			VirtualEndpoint().
			CloudPCs().
			Get(ctx, nil)

		if err != nil {
			errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationRead, d.ReadPermissions)
			return
		}

		for _, cloudPC := range respList.GetValue() {
			cloudPCItem := MapRemoteStateToDataSource(ctx, cloudPC)

			switch filterType {
			case "all":
				filteredItems = append(filteredItems, cloudPCItem)

			case "display_name":
				if cloudPC.GetDisplayName() != nil && strings.Contains(
					strings.ToLower(*cloudPC.GetDisplayName()),
					strings.ToLower(filterValue)) {
					filteredItems = append(filteredItems, cloudPCItem)
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
