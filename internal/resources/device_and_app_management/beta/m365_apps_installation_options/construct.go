package graphBetaM365AppsInstallationOptions

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/constructors"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// constructResource maps the Terraform schema to the SDK model
func constructResource(ctx context.Context, data *M365AppsInstallationOptionsResourceModel) (graphmodels.M365AppsInstallationOptionsable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s resource from model", ResourceName))

	requestBody := graphmodels.NewM365AppsInstallationOptions()

	if !data.UpdateChannel.IsNull() {
		updateChannelStr := data.UpdateChannel.ValueString()
		updateChannel, err := graphmodels.ParseAppsUpdateChannelType(updateChannelStr)
		if err != nil {
			return nil, fmt.Errorf("error parsing update channel: %s", err)
		}
		if updateChannel != nil {
			requestBody.SetUpdateChannel(updateChannel.(*graphmodels.AppsUpdateChannelType))
		}
	}

	if data.AppsForWindows != nil {
		appsForWindows := graphmodels.NewAppsInstallationOptionsForWindows()
		appsForWindows.SetIsMicrosoft365AppsEnabled(data.AppsForWindows.IsMicrosoft365AppsEnabled.ValueBoolPointer())
		appsForWindows.SetIsProjectEnabled(data.AppsForWindows.IsProjectEnabled.ValueBoolPointer())
		appsForWindows.SetIsSkypeForBusinessEnabled(data.AppsForWindows.IsSkypeForBusinessEnabled.ValueBoolPointer())
		appsForWindows.SetIsVisioEnabled(data.AppsForWindows.IsVisioEnabled.ValueBoolPointer())
		requestBody.SetAppsForWindows(appsForWindows)
	}

	if data.AppsForMac != nil {
		appsForMac := graphmodels.NewAppsInstallationOptionsForMac()
		appsForMac.SetIsMicrosoft365AppsEnabled(data.AppsForMac.IsMicrosoft365AppsEnabled.ValueBoolPointer())
		appsForMac.SetIsSkypeForBusinessEnabled(data.AppsForMac.IsSkypeForBusinessEnabled.ValueBoolPointer())
		requestBody.SetAppsForMac(appsForMac)
	}

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]interface{}{
			"error": err.Error(),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing %s resource", ResourceName))

	return requestBody, nil
}
