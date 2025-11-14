package graphBetaWindowsQualityUpdateExpeditePolicy

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// Main entry point to construct the intune windows quality update profile resource for the Terraform provider.
func constructResource(ctx context.Context, data *WindowsQualityUpdateExpeditePolicyResourceModel) (graphmodels.WindowsQualityUpdateProfileable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s resource from model", ResourceName))

	requestBody := graphmodels.NewWindowsQualityUpdateProfile()

	convert.FrameworkToGraphString(data.DisplayName, requestBody.SetDisplayName)
	convert.FrameworkToGraphString(data.Description, requestBody.SetDescription)
	convert.FrameworkToGraphString(data.ReleaseDateDisplayName, requestBody.SetReleaseDateDisplayName)
	convert.FrameworkToGraphString(data.DeployableContentDisplayName, requestBody.SetDeployableContentDisplayName)

	if err := convert.FrameworkToGraphStringSet(ctx, data.RoleScopeTagIds, requestBody.SetRoleScopeTagIds); err != nil {
		return nil, fmt.Errorf("failed to set role scope tags: %s", err)
	}

	if !data.ExpeditedUpdateSettings.IsNull() && !data.ExpeditedUpdateSettings.IsUnknown() {
		expeditedSettings := graphmodels.NewExpeditedWindowsQualityUpdateSettings()

		// Extract attributes from the object
		attrs := data.ExpeditedUpdateSettings.Attributes()

		if qualityUpdateRelease, ok := attrs["quality_update_release"].(types.String); ok {
			convert.FrameworkToGraphString(qualityUpdateRelease, expeditedSettings.SetQualityUpdateRelease)
		}

		if daysUntilForcedReboot, ok := attrs["days_until_forced_reboot"].(types.Int32); ok {
			convert.FrameworkToGraphInt32(daysUntilForcedReboot, expeditedSettings.SetDaysUntilForcedReboot)
		}

		requestBody.SetExpeditedUpdateSettings(expeditedSettings)
	}

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]any{
			"error": err.Error(),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing %s resource", ResourceName))

	return requestBody, nil
}
