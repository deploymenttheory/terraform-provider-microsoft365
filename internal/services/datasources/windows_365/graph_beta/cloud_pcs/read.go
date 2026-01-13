package graphBetaCloudPcs

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
	"github.com/microsoftgraph/msgraph-beta-sdk-go/devicemanagement"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// Read implements the data source Read method.
// It handles three main filtering approaches:
// 1. Get by ID: Direct lookup of a single Cloud PC by its unique identifier (most efficient)
// 2. Server-side filtering: Uses OData query parameters for server-side filtering
// 3. Client-side filtering: Retrieves a list ofall Cloud PCs and filters locally for simple operations
func (d *CloudPcsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var object CloudPcsDataSourceModel
	var cloudPCs []graphmodels.CloudPCable
	var err error

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

	switch filterType {
	case "id":
		singlePC, err := d.getCloudPcById(ctx, object.FilterValue.ValueString())
		if err != nil {
			errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationRead, d.ReadPermissions)
			return
		}

		if singlePC != nil {
			cloudPCs = []graphmodels.CloudPCable{singlePC}
		}

	case "odata":
		cloudPCs, err = d.getCloudPcsWithOData(ctx, object)
		if err != nil {
			errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationRead, d.ReadPermissions)
			return
		}

	default:
		cloudPCs, err = d.listCloudPcs(ctx, filterType, object.FilterValue.ValueString())
		if err != nil {
			errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationRead, d.ReadPermissions)
			return
		}
	}

	var filteredItems []CloudPcItemModel
	for _, cloudPC := range cloudPCs {
		item := MapRemoteStateToDataSource(ctx, cloudPC)
		filteredItems = append(filteredItems, item)
	}

	object.Items = filteredItems

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("Finished Datasource Read Method: %s, found %d items", DataSourceName, len(filteredItems)))
}

// getCloudPcById retrieves a specific Cloud PC by ID
func (d *CloudPcsDataSource) getCloudPcById(ctx context.Context, id string) (graphmodels.CloudPCable, error) {
	requestParameters := &devicemanagement.VirtualEndpointCloudPCsCloudPCItemRequestBuilderGetRequestConfiguration{}

	cloudPC, err := d.client.
		DeviceManagement().
		VirtualEndpoint().
		CloudPCs().
		ByCloudPCId(id).
		Get(ctx, requestParameters)

	if err != nil {
		return nil, err
	}

	return cloudPC, nil
}

// getCloudPcsWithOData retrieves Cloud PCs using OData parameters
func (d *CloudPcsDataSource) getCloudPcsWithOData(ctx context.Context, object CloudPcsDataSourceModel) ([]graphmodels.CloudPCable, error) {
	// Handle OData filtering
	requestParameters := &devicemanagement.VirtualEndpointCloudPCsRequestBuilderGetRequestConfiguration{
		QueryParameters: &devicemanagement.VirtualEndpointCloudPCsRequestBuilderGetQueryParameters{},
	}

	if !object.ODataFilter.IsNull() && object.ODataFilter.ValueString() != "" {
		filter := object.ODataFilter.ValueString()
		requestParameters.QueryParameters.Filter = &filter
		tflog.Debug(ctx, fmt.Sprintf("Using OData $filter: %s", filter))
	}

	if !object.ODataSelect.IsNull() && object.ODataSelect.ValueString() != "" {
		selectFields := strings.Split(object.ODataSelect.ValueString(), ",")
		requestParameters.QueryParameters.Select = selectFields
		tflog.Debug(ctx, fmt.Sprintf("Using OData $select: %v", selectFields))
	}

	if !object.ODataTop.IsNull() && object.ODataTop.ValueInt64() > 0 {
		top := int32(object.ODataTop.ValueInt64())
		requestParameters.QueryParameters.Top = &top
		tflog.Debug(ctx, fmt.Sprintf("Using OData $top: %d", top))
	}

	if !object.ODataCount.IsNull() && object.ODataCount.ValueBool() {
		count := true
		requestParameters.QueryParameters.Count = &count
		tflog.Debug(ctx, "Using OData $count: true")
	}

	respList, err := d.client.
		DeviceManagement().
		VirtualEndpoint().
		CloudPCs().
		Get(ctx, requestParameters)

	if err != nil {
		return nil, err
	}

	return respList.GetValue(), nil
}

// listCloudPcs retrieves all Cloud PCs and applies client-side filtering
func (d *CloudPcsDataSource) listCloudPcs(ctx context.Context, filterType string, filterValue string) ([]graphmodels.CloudPCable, error) {
	respList, err := d.client.
		DeviceManagement().
		VirtualEndpoint().
		CloudPCs().
		Get(ctx, nil)

	if err != nil {
		return nil, err
	}

	cloudPCsList := respList.GetValue()
	tflog.Debug(ctx, fmt.Sprintf("Retrieved %d Cloud PCs", len(cloudPCsList)))

	var filteredPCs []graphmodels.CloudPCable

	// Filter the results based on filter type
	for _, cloudPC := range cloudPCsList {
		switch filterType {
		case "all":
			filteredPCs = append(filteredPCs, cloudPC)
		case "display_name":
			if cloudPC.GetDisplayName() != nil && strings.Contains(
				strings.ToLower(*cloudPC.GetDisplayName()),
				strings.ToLower(filterValue)) {
				filteredPCs = append(filteredPCs, cloudPC)
			}
		}
	}

	return filteredPCs, nil
}
