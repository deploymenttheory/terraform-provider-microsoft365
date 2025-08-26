// REF: https://learn.microsoft.com/en-us/graph/api/intune-devices-windowsmanagementapp-setasmanagedinstaller?view=graph-rest-beta
package graphBetaAppControlForBusinessManagedInstaller

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type AppControlForBusinessManagedInstallerResourceModel struct {
	ID                                          types.String   `tfsdk:"id"`
	IntuneManagementExtensionAsManagedInstaller types.String   `tfsdk:"intune_management_extension_as_managed_installer"`
	AvailableVersion                            types.String   `tfsdk:"available_version"`
	ManagedInstallerConfiguredDateTime          types.String   `tfsdk:"managed_installer_configured_date_time"`
	Timeouts                                    timeouts.Value `tfsdk:"timeouts"`
}
