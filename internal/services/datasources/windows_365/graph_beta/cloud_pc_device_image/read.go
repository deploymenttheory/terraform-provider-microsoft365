package graphBetaCloudPcDeviceImages

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
// It supports three filtering approaches:
// 1. Get by ID: Direct lookup of a single Device Image by its unique identifier
// 2. Get all: Retrieves all Device Images
// 3. Client-side filtering by display name: Retrieves all Device Images and filters locally
func (d *CloudPcDeviceImageDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var object CloudPcDeviceImageDataSourceModel
	var deviceImages []graphmodels.CloudPcDeviceImageable
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
		// Validate filter value is provided
		if object.FilterValue.IsNull() || object.FilterValue.ValueString() == "" {
			resp.Diagnostics.AddError(
				"Missing Required Parameter",
				"When filter_type is 'id', filter_value must be provided with a valid Device Image ID.",
			)
			return
		}

		singleImage, err := d.getDeviceImageById(ctx, object.FilterValue.ValueString())
		if err != nil {
			errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationRead, d.ReadPermissions)
			return
		}

		if singleImage != nil {
			deviceImages = []graphmodels.CloudPcDeviceImageable{singleImage}
		}

	case "display_name":
		// Validate filter value is provided
		if object.FilterValue.IsNull() || object.FilterValue.ValueString() == "" {
			resp.Diagnostics.AddError(
				"Missing Required Parameter",
				"When filter_type is 'display_name', filter_value must be provided with a name to match.",
			)
			return
		}

		deviceImages, err = d.listDeviceImages(ctx, filterType, object.FilterValue.ValueString())
		if err != nil {
			errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationRead, d.ReadPermissions)
			return
		}

	case "all":
		deviceImages, err = d.listDeviceImages(ctx, filterType, "")
		if err != nil {
			errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationRead, d.ReadPermissions)
			return
		}

	default:
		resp.Diagnostics.AddError(
			"Invalid Filter Type",
			fmt.Sprintf("Filter type '%s' is not supported. Supported filter types are: 'all', 'id', 'display_name'.", filterType),
		)
		return
	}

	var filteredItems []CloudPcDeviceImageItem
	for _, deviceImage := range deviceImages {
		item := StateDatasource(ctx, deviceImage)
		filteredItems = append(filteredItems, item)
	}

	object.Items = filteredItems

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("Finished Datasource Read Method: %s, found %d items", DataSourceName, len(filteredItems)))
}

// getDeviceImageById retrieves a specific Device Image by ID
func (d *CloudPcDeviceImageDataSource) getDeviceImageById(ctx context.Context, id string) (graphmodels.CloudPcDeviceImageable, error) {
	requestParameters := &devicemanagement.VirtualEndpointDeviceImagesCloudPcDeviceImageItemRequestBuilderGetRequestConfiguration{}

	deviceImage, err := d.client.
		DeviceManagement().
		VirtualEndpoint().
		DeviceImages().
		ByCloudPcDeviceImageId(id).
		Get(ctx, requestParameters)

	if err != nil {
		return nil, err
	}

	return deviceImage, nil
}

// listDeviceImages retrieves all Device Images and applies client-side filtering if needed
func (d *CloudPcDeviceImageDataSource) listDeviceImages(ctx context.Context, filterType string, filterValue string) ([]graphmodels.CloudPcDeviceImageable, error) {
	respList, err := d.client.
		DeviceManagement().
		VirtualEndpoint().
		DeviceImages().
		Get(ctx, nil)

	if err != nil {
		return nil, err
	}

	deviceImagesList := respList.GetValue()
	tflog.Debug(ctx, fmt.Sprintf("Retrieved %d Device Images", len(deviceImagesList)))

	var filteredImages []graphmodels.CloudPcDeviceImageable

	// Filter the results based on filter type
	for _, deviceImage := range deviceImagesList {
		switch filterType {
		case "all":
			filteredImages = append(filteredImages, deviceImage)
		case "display_name":
			if deviceImage.GetDisplayName() != nil && strings.Contains(
				strings.ToLower(*deviceImage.GetDisplayName()),
				strings.ToLower(filterValue)) {
				filteredImages = append(filteredImages, deviceImage)
			}
		}
	}

	return filteredImages, nil
}
