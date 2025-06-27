package cloudPcs

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/devicemanagement"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

func (d *CloudPcsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var object CloudPcsDataSourceModel

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

	var filteredItems []CloudPcItemModel
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
			errors.HandleGraphError(ctx, err, resp, "Read", d.ReadPermissions)
			return
		}

		item := MapRemoteStateToDataSource(ctx, cloudPC)
		filteredItems = append(filteredItems, item)

	} else {
		// For other filter types, we'll fetch all Cloud PCs and filter client-side
		// except for product_type which can be filtered server-side using OData
		var respList interface{}
		var err error

		if filterType == "product_type" {
			// Use OData filter for product_type
			requestParameters := &devicemanagement.VirtualEndpointCloudPCsRequestBuilderGetRequestConfiguration{}
			requestParameters.QueryParameters = &devicemanagement.VirtualEndpointCloudPCsRequestBuilderGetQueryParameters{
				Filter: &[]string{fmt.Sprintf("productType eq '%s'", filterValue)}[0],
			}

			respList, err = d.client.
				DeviceManagement().
				VirtualEndpoint().
				CloudPCs().
				Get(ctx, requestParameters)
		} else {
			// For other filters, get all items
			respList, err = d.client.
				DeviceManagement().
				VirtualEndpoint().
				CloudPCs().
				Get(ctx, nil)
		}

		if err != nil {
			errors.HandleGraphError(ctx, err, resp, "Read", d.ReadPermissions)
			return
		}

		// Get the value property which contains the Cloud PCs
		cloudPCsList := respList.(graphmodels.CloudPCCollectionResponseable).GetValue()

		for _, cloudPC := range cloudPCsList {
			item := MapRemoteStateToDataSource(ctx, cloudPC)

			switch filterType {
			case "all":
				filteredItems = append(filteredItems, item)
			case "display_name":
				if cloudPC.GetDisplayName() != nil && strings.Contains(
					strings.ToLower(*cloudPC.GetDisplayName()),
					strings.ToLower(filterValue)) {
					filteredItems = append(filteredItems, item)
				}
			case "user_principal_name":
				if cloudPC.GetUserPrincipalName() != nil && strings.Contains(
					strings.ToLower(*cloudPC.GetUserPrincipalName()),
					strings.ToLower(filterValue)) {
					filteredItems = append(filteredItems, item)
				}
			case "status":
				// For status, compare the string representation
				statusString := convert.GraphToFrameworkEnum(cloudPC.GetStatus())
				if statusString.ValueString() == filterValue {
					filteredItems = append(filteredItems, item)
				}
			case "product_type":
				// For product_type, we already filtered server-side
				filteredItems = append(filteredItems, item)
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
