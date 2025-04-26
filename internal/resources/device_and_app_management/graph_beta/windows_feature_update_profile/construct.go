package graphBetaWindowsFeatureUpdateProfile

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/constructors"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// Main entry point to construct the intune windows feature update profile resource for the Terraform provider.
func constructResource(ctx context.Context, data *WindowsFeatureUpdateProfileResourceModel) (graphmodels.WindowsFeatureUpdateProfileable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s resource from model", ResourceName))

	requestBody := graphmodels.NewWindowsFeatureUpdateProfile()

	constructors.SetStringProperty(data.DisplayName, requestBody.SetDisplayName)
	constructors.SetStringProperty(data.Description, requestBody.SetDescription)
	constructors.SetStringProperty(data.FeatureUpdateVersion, requestBody.SetFeatureUpdateVersion)
	constructors.SetBoolProperty(data.InstallLatestWindows10OnWindows11IneligibleDevice, requestBody.SetInstallLatestWindows10OnWindows11IneligibleDevice)
	constructors.SetBoolProperty(data.InstallFeatureUpdatesOptional, requestBody.SetInstallFeatureUpdatesOptional)

	if err := constructors.SetStringSet(ctx, data.RoleScopeTagIds, requestBody.SetRoleScopeTagIds); err != nil {
		return nil, fmt.Errorf("failed to set role scope tags: %s", err)
	}

	if data.RolloutSettings != nil {
		rolloutSettings := graphmodels.NewWindowsUpdateRolloutSettings()

		constructors.StringToTime(data.RolloutSettings.OfferStartDateTimeInUTC, rolloutSettings.SetOfferStartDateTimeInUTC)
		constructors.StringToTime(data.RolloutSettings.OfferEndDateTimeInUTC, rolloutSettings.SetOfferEndDateTimeInUTC)
		constructors.SetInt32Property(data.RolloutSettings.OfferIntervalInDays, rolloutSettings.SetOfferIntervalInDays)

		requestBody.SetRolloutSettings(rolloutSettings)
	}

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]interface{}{
			"error": err.Error(),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing %s resource", ResourceName))

	return requestBody, nil
}
