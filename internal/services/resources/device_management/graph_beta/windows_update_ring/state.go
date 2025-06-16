package graphBetaWindowsUpdateRing

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/state"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// MapRemoteResourceStateToTerraform maps the API response to the Terraform model
func MapRemoteResourceStateToTerraform(ctx context.Context, data *WindowsUpdateRingResourceModel, apiData graphmodels.WindowsUpdateForBusinessConfigurationable) {
	if apiData == nil {
		tflog.Debug(ctx, "Remote resource is nil")
		return
	}

	tflog.Debug(ctx, "Mapping remote state to Terraform", map[string]interface{}{"resourceId": apiData.GetId()})

	data.ID = state.StringPointerValue(apiData.GetId())
	data.DisplayName = state.StringPointerValue(apiData.GetDisplayName())
	data.Description = state.StringPointerValue(apiData.GetDescription())
	data.RoleScopeTagIds = state.StringSliceToSet(ctx, apiData.GetRoleScopeTagIds())
	data.MicrosoftUpdateServiceAllowed = state.BoolPointerValue(apiData.GetMicrosoftUpdateServiceAllowed())
	data.DriversExcluded = state.BoolPointerValue(apiData.GetDriversExcluded())
	data.QualityUpdatesDeferralPeriodInDays = state.Int32PtrToTypeInt32(apiData.GetQualityUpdatesDeferralPeriodInDays())
	data.FeatureUpdatesDeferralPeriodInDays = state.Int32PtrToTypeInt32(apiData.GetFeatureUpdatesDeferralPeriodInDays())
	data.AllowWindows11Upgrade = state.BoolPointerValue(apiData.GetAllowWindows11Upgrade())
	data.QualityUpdatesPaused = state.BoolPointerValue(apiData.GetQualityUpdatesPaused())
	data.FeatureUpdatesPaused = state.BoolPointerValue(apiData.GetFeatureUpdatesPaused())
	data.SkipChecksBeforeRestart = state.BoolPointerValue(apiData.GetSkipChecksBeforeRestart())
	data.BusinessReadyUpdatesOnly = state.EnumPtrToTypeString(apiData.GetBusinessReadyUpdatesOnly())
	data.AutomaticUpdateMode = state.EnumPtrToTypeString(apiData.GetAutomaticUpdateMode())
	data.UpdateNotificationLevel = state.EnumPtrToTypeString(apiData.GetUpdateNotificationLevel())
	data.DeliveryOptimizationMode = state.EnumPtrToTypeString(apiData.GetDeliveryOptimizationMode())
	data.PrereleaseFeatures = state.EnumPtrToTypeString(apiData.GetPrereleaseFeatures())
	data.UpdateWeeks = state.EnumPtrToTypeString(apiData.GetUpdateWeeks())

	if installationSchedule := apiData.GetInstallationSchedule(); installationSchedule != nil {
		if activeHoursInstall, ok := installationSchedule.(graphmodels.WindowsUpdateActiveHoursInstallable); ok {
			if activeHoursInstall.GetActiveHoursStart() != nil {
				data.ActiveHoursStart = state.StringValue(activeHoursInstall.GetActiveHoursStart().String())
			}

			if activeHoursInstall.GetActiveHoursEnd() != nil {
				data.ActiveHoursEnd = state.StringValue(activeHoursInstall.GetActiveHoursEnd().String())
			}
		} else {
			tflog.Warn(ctx, "Installation schedule is not of type WindowsUpdateActiveHoursInstallable")
		}
	}

	data.UserPauseAccess = state.EnumPtrToTypeString(apiData.GetUserPauseAccess())
	data.UserWindowsUpdateScanAccess = state.EnumPtrToTypeString(apiData.GetUserWindowsUpdateScanAccess())
	data.FeatureUpdatesRollbackWindowInDays = state.Int32PtrToTypeInt32(apiData.GetFeatureUpdatesRollbackWindowInDays())
	data.DeadlineForFeatureUpdatesInDays = state.Int32PtrToTypeInt32(apiData.GetDeadlineForFeatureUpdatesInDays())
	data.DeadlineForQualityUpdatesInDays = state.Int32PtrToTypeInt32(apiData.GetDeadlineForQualityUpdatesInDays())
	data.DeadlineGracePeriodInDays = state.Int32PtrToTypeInt32(apiData.GetDeadlineGracePeriodInDays())
	data.PostponeRebootUntilAfterDeadline = state.BoolPointerValue(apiData.GetPostponeRebootUntilAfterDeadline())
	data.EngagedRestartDeadlineInDays = state.Int32PtrToTypeInt32(apiData.GetEngagedRestartDeadlineInDays())
	data.EngagedRestartSnoozeScheduleInDays = state.Int32PtrToTypeInt32(apiData.GetEngagedRestartSnoozeScheduleInDays())
	data.EngagedRestartTransitionScheduleInDays = state.Int32PtrToTypeInt32(apiData.GetEngagedRestartTransitionScheduleInDays())
	data.AutoRestartNotificationDismissal = state.EnumPtrToTypeString(apiData.GetAutoRestartNotificationDismissal())
	data.ScheduleRestartWarningInHours = state.Int32PtrToTypeInt32(apiData.GetScheduleRestartWarningInHours())
	data.ScheduleImminentRestartWarningInMinutes = state.Int32PtrToTypeInt32(apiData.GetScheduleImminentRestartWarningInMinutes())
	data.EngagedRestartSnoozeScheduleForFeatureUpdatesInDays = state.Int32PtrToTypeInt32(apiData.GetEngagedRestartSnoozeScheduleInDays())
	data.EngagedRestartTransitionScheduleForFeatureUpdatesInDays = state.Int32PtrToTypeInt32(apiData.GetEngagedRestartTransitionScheduleInDays())

	tflog.Debug(ctx, fmt.Sprintf("Finished stating resource %s with id %s", ResourceName, data.ID.ValueString()))

}
