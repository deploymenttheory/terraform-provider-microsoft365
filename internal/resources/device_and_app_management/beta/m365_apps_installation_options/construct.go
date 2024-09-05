package graphBetaM365AppsInstallationOptions

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// constructResource maps the Terraform schema to the SDK model
func constructResource(ctx context.Context, data *M365AppsInstallationOptionsResourceModel) (models.AdminMicrosoft365Appsable, error) {
	installationOptions := models.NewM365AppsInstallationOptions()

	if !data.UpdateChannel.IsNull() {
		updateChannelStr := data.UpdateChannel.ValueString()
		updateChannel, err := models.ParseAppsUpdateChannelType(updateChannelStr)
		if err != nil {
			return nil, fmt.Errorf("error parsing update channel: %s", err)
		}
		if updateChannel != nil {
			installationOptions.SetUpdateChannel(updateChannel.(*models.AppsUpdateChannelType))
		}
	}

	if data.AppsForWindows != nil {
		appsForWindows := models.NewAppsInstallationOptionsForWindows()
		appsForWindows.SetIsMicrosoft365AppsEnabled(data.AppsForWindows.IsMicrosoft365AppsEnabled.ValueBoolPointer())
		appsForWindows.SetIsProjectEnabled(data.AppsForWindows.IsProjectEnabled.ValueBoolPointer())
		appsForWindows.SetIsSkypeForBusinessEnabled(data.AppsForWindows.IsSkypeForBusinessEnabled.ValueBoolPointer())
		appsForWindows.SetIsVisioEnabled(data.AppsForWindows.IsVisioEnabled.ValueBoolPointer())
		installationOptions.SetAppsForWindows(appsForWindows)
	}

	if data.AppsForMac != nil {
		appsForMac := models.NewAppsInstallationOptionsForMac()
		appsForMac.SetIsMicrosoft365AppsEnabled(data.AppsForMac.IsMicrosoft365AppsEnabled.ValueBoolPointer())
		appsForMac.SetIsSkypeForBusinessEnabled(data.AppsForMac.IsSkypeForBusinessEnabled.ValueBoolPointer())
		installationOptions.SetAppsForMac(appsForMac)
	}

	// Create an AdminMicrosoft365Apps object and set the installation options
	requestBody := models.NewAdminMicrosoft365Apps()
	requestBody.SetInstallationOptions(installationOptions)

	requestBodyJSON, err := json.MarshalIndent(requestBody, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("error marshalling request body to JSON: %s", err)
	}

	tflog.Debug(ctx, "Constructed AdminMicrosoft365Apps resource:\n"+string(requestBodyJSON))

	return requestBody, nil
}
