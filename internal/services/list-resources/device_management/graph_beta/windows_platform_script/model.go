package graphBetaWindowsPlatformScript

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// WindowsPlatformScriptListConfigModel represents the configuration for listing Windows platform scripts
type WindowsPlatformScriptListConfigModel struct {
	DisplayNameFilter  types.String `tfsdk:"display_name_filter"`
	FileNameFilter     types.String `tfsdk:"file_name_filter"`
	RunAsAccountFilter types.String `tfsdk:"run_as_account_filter"`
	IsAssignedFilter   types.Bool   `tfsdk:"is_assigned_filter"`
	ODataFilter        types.String `tfsdk:"odata_filter"`
}
