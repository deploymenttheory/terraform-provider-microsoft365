package graphBetaWindowsUpdatesDeviceEnrollment

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/datasource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type DeviceEnrollmentDataSourceModel struct {
	EntraDeviceId  types.String     `tfsdk:"entra_device_id"`
	DeviceName     types.String     `tfsdk:"device_name"`
	ListAll        types.Bool       `tfsdk:"list_all"`
	UpdateCategory types.String     `tfsdk:"update_category"`
	ODataFilter    types.String     `tfsdk:"odata_filter"`
	Devices        []EnrolledDevice `tfsdk:"devices"`
	Timeouts       timeouts.Value   `tfsdk:"timeouts"`
}

type EnrolledDevice struct {
	ID          types.String                 `tfsdk:"id"`
	Enrollments []UpdateManagementEnrollment `tfsdk:"enrollments"`
	Errors      []UpdatableAssetError        `tfsdk:"errors"`
}

type UpdateManagementEnrollment struct {
	UpdateCategory types.String `tfsdk:"update_category"`
}

type UpdatableAssetError struct {
	ErrorCode    types.String `tfsdk:"error_code"`
	ErrorMessage types.String `tfsdk:"error_message"`
}
