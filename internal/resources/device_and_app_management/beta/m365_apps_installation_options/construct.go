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

	requestBody := models.NewAdminMicrosoft365Apps()
	requestBody.SetInstallationOptions(installationOptions)

	// Debug logging
	debugPrintRequestBody(ctx, requestBody)

	return requestBody, nil
}

func debugPrintRequestBody(ctx context.Context, requestBody models.AdminMicrosoft365Appsable) {
	requestMap := map[string]interface{}{
		"installationOptions": debugMapInstallationOptions(requestBody.GetInstallationOptions()),
	}

	requestBodyJSON, err := json.MarshalIndent(requestMap, "", "  ")
	if err != nil {
		tflog.Error(ctx, "Error marshalling request body to JSON", map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	tflog.Debug(ctx, "Constructed AdminMicrosoft365Apps resource", map[string]interface{}{
		"requestBody": string(requestBodyJSON),
	})
}

func debugMapInstallationOptions(options models.M365AppsInstallationOptionsable) map[string]interface{} {
	optionsMap := map[string]interface{}{
		"updateChannel": options.GetUpdateChannel(),
	}

	if appsForWindows := options.GetAppsForWindows(); appsForWindows != nil {
		optionsMap["appsForWindows"] = map[string]interface{}{
			"isMicrosoft365AppsEnabled": appsForWindows.GetIsMicrosoft365AppsEnabled(),
			"isProjectEnabled":          appsForWindows.GetIsProjectEnabled(),
			"isSkypeForBusinessEnabled": appsForWindows.GetIsSkypeForBusinessEnabled(),
			"isVisioEnabled":            appsForWindows.GetIsVisioEnabled(),
		}
	}

	if appsForMac := options.GetAppsForMac(); appsForMac != nil {
		optionsMap["appsForMac"] = map[string]interface{}{
			"isMicrosoft365AppsEnabled": appsForMac.GetIsMicrosoft365AppsEnabled(),
			"isSkypeForBusinessEnabled": appsForMac.GetIsSkypeForBusinessEnabled(),
		}
	}

	return optionsMap
}
