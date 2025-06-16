package graphBetaLinuxPlatformScript

import (
	"context"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

func MapRemoteSettingsStateToTerraform(ctx context.Context, data *LinuxPlatformScriptResourceModel, remoteSettings []graphmodels.DeviceManagementConfigurationSettingable) {
	if remoteSettings == nil {
		tflog.Debug(ctx, "Remote settings are nil")
		return
	}

	tflog.Debug(ctx, "Finished mapping settings state to Terraform state")
}
