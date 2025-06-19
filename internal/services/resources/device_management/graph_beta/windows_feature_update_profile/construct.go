package graphBetaWindowsFeatureUpdateProfile

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// Main entry point to construct the intune windows feature update profile resource for the Terraform provider.
//
// If forUpdate is true, only PATCH-allowed fields are populated.
func constructResource(ctx context.Context, data *WindowsFeatureUpdateProfileResourceModel, forUpdate bool) (graphmodels.WindowsFeatureUpdateProfileable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s resource from model (forUpdate=%v)", ResourceName, forUpdate))

	requestBody := graphmodels.NewWindowsFeatureUpdateProfile()

	convert.FrameworkToGraphString(data.DisplayName, requestBody.SetDisplayName)
	convert.FrameworkToGraphString(data.Description, requestBody.SetDescription)
	convert.FrameworkToGraphString(data.FeatureUpdateVersion, requestBody.SetFeatureUpdateVersion)
	convert.FrameworkToGraphBool(data.InstallFeatureUpdatesOptional, requestBody.SetInstallFeatureUpdatesOptional)

	if err := convert.FrameworkToGraphStringSet(ctx, data.RoleScopeTagIds, requestBody.SetRoleScopeTagIds); err != nil {
		return nil, fmt.Errorf("failed to set role scope tags: %s", err)
	}

	if data.RolloutSettings != nil {
		rolloutSettings := graphmodels.NewWindowsUpdateRolloutSettings()

		convert.FrameworkToGraphTime(data.RolloutSettings.OfferStartDateTimeInUTC, rolloutSettings.SetOfferStartDateTimeInUTC)
		convert.FrameworkToGraphTime(data.RolloutSettings.OfferEndDateTimeInUTC, rolloutSettings.SetOfferEndDateTimeInUTC)
		convert.FrameworkToGraphInt32(data.RolloutSettings.OfferIntervalInDays, rolloutSettings.SetOfferIntervalInDays)

		requestBody.SetRolloutSettings(rolloutSettings)
	}

	// Immutable field once created. Excluded from update req construction.
	if !forUpdate {
		convert.FrameworkToGraphBool(data.InstallLatestWindows10OnWindows11IneligibleDevice, requestBody.SetInstallLatestWindows10OnWindows11IneligibleDevice)
	}

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]interface{}{
			"error": err.Error(),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing %s resource", ResourceName))

	return requestBody, nil
}
