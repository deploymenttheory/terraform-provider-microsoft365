// REF: https://learn.microsoft.com/en-us/graph/api/resources/m365appsinstallationoptions?view=graph-rest-1.0
package graphM365AppsInstallationOptions

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type M365AppsInstallationOptionsResourceModel struct {
	ID             types.String                       `tfsdk:"id"`
	UpdateChannel  types.String                       `tfsdk:"update_channel"`
	AppsForWindows *AppsInstallationOptionsForWindows `tfsdk:"apps_for_windows"`
	AppsForMac     *AppsInstallationOptionsForMac     `tfsdk:"apps_for_mac"`
	Timeouts       timeouts.Value                     `tfsdk:"timeouts"`
}

type AppsInstallationOptionsForWindows struct {
	IsMicrosoft365AppsEnabled types.Bool `tfsdk:"is_microsoft_365_apps_enabled"`
	IsSkypeForBusinessEnabled types.Bool `tfsdk:"is_skype_for_business_enabled"`
}

type AppsInstallationOptionsForMac struct {
	IsMicrosoft365AppsEnabled types.Bool `tfsdk:"is_microsoft_365_apps_enabled"`
	IsSkypeForBusinessEnabled types.Bool `tfsdk:"is_skype_for_business_enabled"`
}
