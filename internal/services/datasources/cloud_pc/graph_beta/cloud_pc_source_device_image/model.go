package cloudPcSourceDeviceImage

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// CloudPcSourceDeviceImageDataSourceModel represents the Terraform data source model for source device images
// See: https://learn.microsoft.com/en-us/graph/api/cloudpcdeviceimage-getsourceimages?view=graph-rest-beta
type CloudPcSourceDeviceImageDataSourceModel struct {
	FilterType  types.String                        `tfsdk:"filter_type"`
	FilterValue types.String                        `tfsdk:"filter_value"`
	Items       []CloudPcSourceDeviceImageItemModel `tfsdk:"items"`
	Timeouts    timeouts.Value                      `tfsdk:"timeouts"`
}

type CloudPcSourceDeviceImageItemModel struct {
	ID                      types.String `tfsdk:"id"`
	ResourceId              types.String `tfsdk:"resource_id"`
	DisplayName             types.String `tfsdk:"display_name"`
	SubscriptionId          types.String `tfsdk:"subscription_id"`
	SubscriptionDisplayName types.String `tfsdk:"subscription_display_name"`
}
