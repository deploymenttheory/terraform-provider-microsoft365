package graphBetaWindowsUpdateRing

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-framework/types"
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

	data.ID = convert.GraphToFrameworkString(apiData.GetId())
	data.DisplayName = convert.GraphToFrameworkString(apiData.GetDisplayName())
	data.Description = convert.GraphToFrameworkString(apiData.GetDescription())
	data.RoleScopeTagIds = convert.GraphToFrameworkStringSet(ctx, apiData.GetRoleScopeTagIds())
	data.MicrosoftUpdateServiceAllowed = convert.GraphToFrameworkBool(apiData.GetMicrosoftUpdateServiceAllowed())
	data.DriversExcluded = convert.GraphToFrameworkBool(apiData.GetDriversExcluded())
	data.QualityUpdatesDeferralPeriodInDays = convert.GraphToFrameworkInt32(apiData.GetQualityUpdatesDeferralPeriodInDays())
	data.FeatureUpdatesDeferralPeriodInDays = convert.GraphToFrameworkInt32(apiData.GetFeatureUpdatesDeferralPeriodInDays())
	data.AllowWindows11Upgrade = convert.GraphToFrameworkBool(apiData.GetAllowWindows11Upgrade())
	data.QualityUpdatesPaused = convert.GraphToFrameworkBool(apiData.GetQualityUpdatesPaused())
	data.FeatureUpdatesPaused = convert.GraphToFrameworkBool(apiData.GetFeatureUpdatesPaused())
	data.SkipChecksBeforeRestart = convert.GraphToFrameworkBool(apiData.GetSkipChecksBeforeRestart())
	data.BusinessReadyUpdatesOnly = convert.GraphToFrameworkEnum(apiData.GetBusinessReadyUpdatesOnly())
	data.AutomaticUpdateMode = convert.GraphToFrameworkEnum(apiData.GetAutomaticUpdateMode())
	data.UpdateNotificationLevel = convert.GraphToFrameworkEnum(apiData.GetUpdateNotificationLevel())
	data.DeliveryOptimizationMode = convert.GraphToFrameworkEnum(apiData.GetDeliveryOptimizationMode())
	data.PrereleaseFeatures = convert.GraphToFrameworkEnum(apiData.GetPrereleaseFeatures())
	data.UpdateWeeks = convert.GraphToFrameworkEnum(apiData.GetUpdateWeeks())

	if installationSchedule := apiData.GetInstallationSchedule(); installationSchedule != nil {
		if activeHoursInstall, ok := installationSchedule.(graphmodels.WindowsUpdateActiveHoursInstallable); ok {
			if activeHoursInstall.GetActiveHoursStart() != nil {
				data.ActiveHoursStart = types.StringValue(activeHoursInstall.GetActiveHoursStart().String())
			}

			if activeHoursInstall.GetActiveHoursEnd() != nil {
				data.ActiveHoursEnd = types.StringValue(activeHoursInstall.GetActiveHoursEnd().String())
			}
		} else {
			tflog.Warn(ctx, "Installation schedule is not of type WindowsUpdateActiveHoursInstallable")
		}
	}

	data.UserPauseAccess = convert.GraphToFrameworkEnum(apiData.GetUserPauseAccess())
	data.UserWindowsUpdateScanAccess = convert.GraphToFrameworkEnum(apiData.GetUserWindowsUpdateScanAccess())
	data.FeatureUpdatesRollbackWindowInDays = convert.GraphToFrameworkInt32(apiData.GetFeatureUpdatesRollbackWindowInDays())
	data.DeadlineForFeatureUpdatesInDays = convert.GraphToFrameworkInt32(apiData.GetDeadlineForFeatureUpdatesInDays())
	data.DeadlineForQualityUpdatesInDays = convert.GraphToFrameworkInt32(apiData.GetDeadlineForQualityUpdatesInDays())
	data.DeadlineGracePeriodInDays = convert.GraphToFrameworkInt32(apiData.GetDeadlineGracePeriodInDays())
	data.PostponeRebootUntilAfterDeadline = convert.GraphToFrameworkBool(apiData.GetPostponeRebootUntilAfterDeadline())
	data.EngagedRestartDeadlineInDays = convert.GraphToFrameworkInt32(apiData.GetEngagedRestartDeadlineInDays())
	data.EngagedRestartSnoozeScheduleInDays = convert.GraphToFrameworkInt32(apiData.GetEngagedRestartSnoozeScheduleInDays())
	data.EngagedRestartTransitionScheduleInDays = convert.GraphToFrameworkInt32(apiData.GetEngagedRestartTransitionScheduleInDays())
	data.AutoRestartNotificationDismissal = convert.GraphToFrameworkEnum(apiData.GetAutoRestartNotificationDismissal())
	data.ScheduleRestartWarningInHours = convert.GraphToFrameworkInt32(apiData.GetScheduleRestartWarningInHours())
	data.ScheduleImminentRestartWarningInMinutes = convert.GraphToFrameworkInt32(apiData.GetScheduleImminentRestartWarningInMinutes())
	data.EngagedRestartSnoozeScheduleForFeatureUpdatesInDays = convert.GraphToFrameworkInt32(apiData.GetEngagedRestartSnoozeScheduleInDays())
	data.EngagedRestartTransitionScheduleForFeatureUpdatesInDays = convert.GraphToFrameworkInt32(apiData.GetEngagedRestartTransitionScheduleInDays())

	tflog.Debug(ctx, fmt.Sprintf("Finished stating resource %s with id %s", ResourceName, data.ID.ValueString()))

}
