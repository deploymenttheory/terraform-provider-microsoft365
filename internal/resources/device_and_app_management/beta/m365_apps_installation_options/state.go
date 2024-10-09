package graphBetaM365AppsInstallationOptions

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/state"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

func MapRemoteStateToTerraform(ctx context.Context, data *M365AppsInstallationOptionsResourceModel, remoteResource graphmodels.M365AppsInstallationOptionsable) {
	if remoteResource == nil {
		tflog.Debug(ctx, "Remote resource is nil")
		return
	}

	tflog.Debug(ctx, "Starting to map remote state to Terraform state for M365AppsInstallationOptions")

	// Check if UpdateChannel exists and map it
	updateChannel := remoteResource.GetUpdateChannel()
	if updateChannel != nil {
		data.UpdateChannel = state.EnumPtrToTypeString(updateChannel)
	}

	// Map AppsForWindows if it exists
	if remoteWindows := remoteResource.GetAppsForWindows(); remoteWindows != nil {
		data.AppsForWindows = &AppsInstallationOptionsForWindows{
			IsMicrosoft365AppsEnabled: state.BoolPtrToTypeBool(remoteWindows.GetIsMicrosoft365AppsEnabled()),
			IsProjectEnabled:          state.BoolPtrToTypeBool(remoteWindows.GetIsProjectEnabled()),
			IsSkypeForBusinessEnabled: state.BoolPtrToTypeBool(remoteWindows.GetIsSkypeForBusinessEnabled()),
			IsVisioEnabled:            state.BoolPtrToTypeBool(remoteWindows.GetIsVisioEnabled()),
		}
	} else {
		data.AppsForWindows = nil
	}

	// Map AppsForMac if it exists
	if remoteMac := remoteResource.GetAppsForMac(); remoteMac != nil {
		data.AppsForMac = &AppsInstallationOptionsForMac{
			IsMicrosoft365AppsEnabled: state.BoolPtrToTypeBool(remoteMac.GetIsMicrosoft365AppsEnabled()),
			IsSkypeForBusinessEnabled: state.BoolPtrToTypeBool(remoteMac.GetIsSkypeForBusinessEnabled()),
		}
	} else {
		data.AppsForMac = nil
	}

	tflog.Debug(ctx, "Finished mapping remote state to Terraform state for M365AppsInstallationOptions")
}
