package graphBetaM365AppsInstallationOptions

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/constructors"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

func constructResource(ctx context.Context, data *M365AppsInstallationOptionsResourceModel) (graphmodels.M365AppsInstallationOptionsable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s resource from model", ResourceName))

	requestBody := graphmodels.NewM365AppsInstallationOptions()

	if err := constructors.SetEnumProperty(data.UpdateChannel,
		graphmodels.ParseAppsUpdateChannelType,
		requestBody.SetUpdateChannel); err != nil {
		return nil, fmt.Errorf("failed to set update channel: %v", err)
	}

	// Handle Windows apps
	if data.AppsForWindows != nil {
		appsForWindows := graphmodels.NewAppsInstallationOptionsForWindows()
		constructors.SetBoolProperty(data.AppsForWindows.IsMicrosoft365AppsEnabled, appsForWindows.SetIsMicrosoft365AppsEnabled)
		constructors.SetBoolProperty(data.AppsForWindows.IsProjectEnabled, appsForWindows.SetIsProjectEnabled)
		constructors.SetBoolProperty(data.AppsForWindows.IsSkypeForBusinessEnabled, appsForWindows.SetIsSkypeForBusinessEnabled)
		constructors.SetBoolProperty(data.AppsForWindows.IsVisioEnabled, appsForWindows.SetIsVisioEnabled)
		requestBody.SetAppsForWindows(appsForWindows)
	}

	// Handle Mac apps
	if data.AppsForMac != nil {
		appsForMac := graphmodels.NewAppsInstallationOptionsForMac()
		constructors.SetBoolProperty(data.AppsForMac.IsMicrosoft365AppsEnabled, appsForMac.SetIsMicrosoft365AppsEnabled)
		constructors.SetBoolProperty(data.AppsForMac.IsSkypeForBusinessEnabled, appsForMac.SetIsSkypeForBusinessEnabled)
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