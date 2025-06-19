package graphBetaWindowsFeatureUpdateProfile

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// MapRemoteResourceStateToTerraform maps the Graph API model into the Terraform state model
func MapRemoteResourceStateToTerraform(ctx context.Context, data *WindowsFeatureUpdateProfileResourceModel, remoteResource graphmodels.WindowsFeatureUpdateProfileable) {
	if remoteResource == nil {
		tflog.Debug(ctx, "Remote resource is nil")
		return
	}

	tflog.Debug(ctx, "Mapping remote state to Terraform", map[string]interface{}{"resourceId": remoteResource.GetId()})

	data.ID = convert.GraphToFrameworkString(remoteResource.GetId())
	data.DisplayName = convert.GraphToFrameworkString(remoteResource.GetDisplayName())
	data.Description = convert.GraphToFrameworkString(remoteResource.GetDescription())
	data.FeatureUpdateVersion = convert.GraphToFrameworkString(remoteResource.GetFeatureUpdateVersion())
	data.CreatedDateTime = convert.GraphToFrameworkTime(remoteResource.GetCreatedDateTime())
	data.LastModifiedDateTime = convert.GraphToFrameworkTime(remoteResource.GetLastModifiedDateTime())
	data.RoleScopeTagIds = convert.GraphToFrameworkStringSet(ctx, remoteResource.GetRoleScopeTagIds())
	data.DeployableContentDisplayName = convert.GraphToFrameworkString(remoteResource.GetDeployableContentDisplayName())
	data.EndOfSupportDate = convert.GraphToFrameworkTime(remoteResource.GetEndOfSupportDate())
	data.InstallLatestWindows10OnWindows11IneligibleDevice = convert.GraphToFrameworkBool(remoteResource.GetInstallLatestWindows10OnWindows11IneligibleDevice())
	data.InstallFeatureUpdatesOptional = convert.GraphToFrameworkBool(remoteResource.GetInstallFeatureUpdatesOptional())

	// Handles scenarios when rollout_settings block is not included within request
	// equivilent to the rollout option 'Make update available as soon as possible'
	if rolloutSettings := remoteResource.GetRolloutSettings(); rolloutSettings != nil {
		if rolloutSettings.GetOfferStartDateTimeInUTC() != nil ||
			rolloutSettings.GetOfferEndDateTimeInUTC() != nil ||
			rolloutSettings.GetOfferIntervalInDays() != nil {

			data.RolloutSettings = &RolloutSettingsModel{
				OfferStartDateTimeInUTC: convert.GraphToFrameworkTime(rolloutSettings.GetOfferStartDateTimeInUTC()),
				OfferEndDateTimeInUTC:   convert.GraphToFrameworkTime(rolloutSettings.GetOfferEndDateTimeInUTC()),
				OfferIntervalInDays:     convert.GraphToFrameworkInt32(rolloutSettings.GetOfferIntervalInDays()),
			}
		}
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished stating resource %s with id %s", ResourceName, data.ID.ValueString()))
}
