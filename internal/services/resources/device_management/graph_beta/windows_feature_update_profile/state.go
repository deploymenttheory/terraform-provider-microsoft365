package graphBetaWindowsFeatureUpdateProfile

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/state"
	"github.com/hashicorp/terraform-plugin-framework/types"
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

	data.ID = types.StringPointerValue(remoteResource.GetId())
	data.DisplayName = types.StringPointerValue(remoteResource.GetDisplayName())
	data.Description = types.StringPointerValue(remoteResource.GetDescription())
	data.FeatureUpdateVersion = types.StringPointerValue(remoteResource.GetFeatureUpdateVersion())
	data.CreatedDateTime = state.TimeToString(remoteResource.GetCreatedDateTime())
	data.LastModifiedDateTime = state.TimeToString(remoteResource.GetLastModifiedDateTime())
	data.RoleScopeTagIds = state.StringSliceToSet(ctx, remoteResource.GetRoleScopeTagIds())
	data.DeployableContentDisplayName = types.StringPointerValue(remoteResource.GetDeployableContentDisplayName())
	data.EndOfSupportDate = state.TimeToString(remoteResource.GetEndOfSupportDate())
	data.InstallLatestWindows10OnWindows11IneligibleDevice = state.BoolPtrToTypeBool(remoteResource.GetInstallLatestWindows10OnWindows11IneligibleDevice())
	data.InstallFeatureUpdatesOptional = state.BoolPtrToTypeBool(remoteResource.GetInstallFeatureUpdatesOptional())

	// Handles scenarios when rollout_settings block is not included within request
	// equivilent to the rollout option 'Make update available as soon as possible'
	if rolloutSettings := remoteResource.GetRolloutSettings(); rolloutSettings != nil {
		if rolloutSettings.GetOfferStartDateTimeInUTC() != nil ||
			rolloutSettings.GetOfferEndDateTimeInUTC() != nil ||
			rolloutSettings.GetOfferIntervalInDays() != nil {

			data.RolloutSettings = &RolloutSettingsModel{
				OfferStartDateTimeInUTC: state.TimeToString(rolloutSettings.GetOfferStartDateTimeInUTC()),
				OfferEndDateTimeInUTC:   state.TimeToString(rolloutSettings.GetOfferEndDateTimeInUTC()),
				OfferIntervalInDays:     state.Int32PtrToTypeInt32(rolloutSettings.GetOfferIntervalInDays()),
			}
		}
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished stating resource %s with id %s", ResourceName, data.ID.ValueString()))
}
