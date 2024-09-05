package graphBetaM365AppsInstallationOptions

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/state"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

func MapRemoteStateToTerraform(ctx context.Context, data *M365AppsInstallationOptionsResourceModel, remoteResource models.AdminMicrosoft365Appsable) {
	if remoteResource == nil {
		tflog.Debug(ctx, "Remote resource is nil")
		return
	}

	tflog.Debug(ctx, "Starting to map remote state to Terraform state for M365AppsInstallationOptions")

	installationOptions := remoteResource.GetInstallationOptions()
	if installationOptions == nil {
		tflog.Debug(ctx, "Installation options are nil")
		return
	}

	// Map UpdateChannel
	data.UpdateChannel = state.EnumPtrToTypeString(installationOptions.GetUpdateChannel())

	// Map AppsForWindows
	if remoteWindows := installationOptions.GetAppsForWindows(); remoteWindows != nil {
		data.AppsForWindows = &AppsInstallationOptionsForWindows{
			IsMicrosoft365AppsEnabled: types.BoolPointerValue(remoteWindows.GetIsMicrosoft365AppsEnabled()),
			IsProjectEnabled:          types.BoolPointerValue(remoteWindows.GetIsProjectEnabled()),
			IsSkypeForBusinessEnabled: types.BoolPointerValue(remoteWindows.GetIsSkypeForBusinessEnabled()),
			IsVisioEnabled:            types.BoolPointerValue(remoteWindows.GetIsVisioEnabled()),
		}
	}

	// Map AppsForMac
	if remoteMac := installationOptions.GetAppsForMac(); remoteMac != nil {
		data.AppsForMac = &AppsInstallationOptionsForMac{
			IsMicrosoft365AppsEnabled: types.BoolPointerValue(remoteMac.GetIsMicrosoft365AppsEnabled()),
			IsSkypeForBusinessEnabled: types.BoolPointerValue(remoteMac.GetIsSkypeForBusinessEnabled()),
		}
	}

	tflog.Debug(ctx, "Finished mapping remote state to Terraform state for M365AppsInstallationOptions")
}
