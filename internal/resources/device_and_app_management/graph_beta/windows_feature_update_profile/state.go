package graphBetaWindowsFeatureUpdateProfile

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/state"
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

	if rolloutSettings := remoteResource.GetRolloutSettings(); rolloutSettings != nil {
		data.RolloutSettings = &RolloutSettingsModel{
			OfferStartDateTimeInUTC: state.TimeToString(rolloutSettings.GetOfferStartDateTimeInUTC()),
			OfferEndDateTimeInUTC:   state.TimeToString(rolloutSettings.GetOfferEndDateTimeInUTC()),
			OfferIntervalInDays:     state.Int32PtrToTypeInt32(rolloutSettings.GetOfferIntervalInDays()),
		}
	} else {
		data.RolloutSettings = nil
	}

	tflog.Debug(ctx, "Finished mapping remote state to Terraform", map[string]interface{}{"resourceId": data.ID.ValueString()})
}
