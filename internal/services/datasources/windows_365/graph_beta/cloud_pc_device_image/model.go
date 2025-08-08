// REF: https://learn.microsoft.com/en-us/graph/api/virtualendpoint-list-deviceimages?view=graph-rest-beta

package graphBetaCloudPcDeviceImages

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// CloudPcDeviceImageDataSourceModel represents the Terraform data source model for Cloud PC Device Images
type CloudPcDeviceImageDataSourceModel struct {
	FilterType  types.String             `tfsdk:"filter_type"`
	FilterValue types.String             `tfsdk:"filter_value"`
	Items       []CloudPcDeviceImageItem `tfsdk:"items"`
	Timeouts    timeouts.Value           `tfsdk:"timeouts"`
}

// CloudPcDeviceImageItem represents an individual Cloud PC Device Image
type CloudPcDeviceImageItem struct {
	ID                    types.String `tfsdk:"id"`
	DisplayName           types.String `tfsdk:"display_name"`
	ExpirationDate        types.String `tfsdk:"expiration_date"`
	OSBuildNumber         types.String `tfsdk:"os_build_number"`
	OSStatus              types.String `tfsdk:"os_status"`
	OperatingSystem       types.String `tfsdk:"operating_system"`
	Version               types.String `tfsdk:"version"`
	SourceImageResourceID types.String `tfsdk:"source_image_resource_id"`
	LastModifiedDateTime  types.String `tfsdk:"last_modified_date_time"`
	Status                types.String `tfsdk:"status"`
	StatusDetails         types.String `tfsdk:"status_details"`
	ErrorCode             types.String `tfsdk:"error_code"`
	OSVersionNumber       types.String `tfsdk:"os_version_number"`
}
