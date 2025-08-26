package graphBetaCloudPcGalleryImage

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud"
	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/devicemanagement"
)

func (d *CloudPcGalleryImageDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var object CloudPcGalleryImageDataSourceModel

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

	var filteredItems []CloudPcGalleryImageItemModel
	filterValue := object.FilterValue.ValueString()

	if filterType == "id" {
		requestParameters := &devicemanagement.VirtualEndpointGalleryImagesCloudPcGalleryImageItemRequestBuilderGetRequestConfiguration{}

		galleryImage, err := d.client.
			DeviceManagement().
			VirtualEndpoint().
			GalleryImages().
			ByCloudPcGalleryImageId(filterValue).
			Get(ctx, requestParameters)

		if err != nil {
			errors.HandleKiotaGraphError(ctx, err, resp, "Read", d.ReadPermissions)
			return
		}

		item := MapRemoteStateToDataSource(ctx, galleryImage)
		filteredItems = append(filteredItems, item)

	} else {
		respList, err := d.client.
			DeviceManagement().
			VirtualEndpoint().
			GalleryImages().
			Get(ctx, nil)

		if err != nil {
			errors.HandleKiotaGraphError(ctx, err, resp, "Read", d.ReadPermissions)
			return
		}

		for _, galleryImage := range respList.GetValue() {
			item := MapRemoteStateToDataSource(ctx, galleryImage)
			switch filterType {
			case "all":
				filteredItems = append(filteredItems, item)
			case "display_name":
				if galleryImage.GetDisplayName() != nil && strings.Contains(
					strings.ToLower(*galleryImage.GetDisplayName()),
					strings.ToLower(filterValue)) {
					filteredItems = append(filteredItems, item)
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
