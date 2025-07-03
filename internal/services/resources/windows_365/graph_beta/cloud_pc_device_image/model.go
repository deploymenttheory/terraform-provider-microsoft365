// REF: https://learn.microsoft.com/en-us/graph/api/resources/cloudpcdeviceimage?view=graph-rest-beta
package graphBetaCloudPcDeviceImage

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type CloudPcDeviceImageResourceModel struct {
	ID                    types.String   `tfsdk:"id"`
	DisplayName           types.String   `tfsdk:"display_name"`
	ErrorCode             types.String   `tfsdk:"error_code"`
	ExpirationDate        types.String   `tfsdk:"expiration_date"`
	LastModifiedDateTime  types.String   `tfsdk:"last_modified_date_time"`
	OperatingSystem       types.String   `tfsdk:"operating_system"`
	OsBuildNumber         types.String   `tfsdk:"os_build_number"`
	OsStatus              types.String   `tfsdk:"os_status"`
	OsVersionNumber       types.String   `tfsdk:"os_version_number"`
	SourceImageResourceId types.String   `tfsdk:"source_image_resource_id"`
	Status                types.String   `tfsdk:"status"`
	Version               types.String   `tfsdk:"version"`
	Timeouts              timeouts.Value `tfsdk:"timeouts"`
}
