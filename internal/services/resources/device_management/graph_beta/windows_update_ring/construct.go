package graphBetaWindowsUpdateRing

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// constructResource creates an assign request body with assignments from the nested blocks
func constructResource(ctx context.Context, data *WindowsUpdateRingResourceModel) (graphmodels.WindowsUpdateForBusinessConfigurationable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s resource from model", ResourceName))

	requestBody := graphmodels.NewWindowsUpdateForBusinessConfiguration()

	convert.FrameworkToGraphString(data.DisplayName, requestBody.SetDisplayName)
	convert.FrameworkToGraphString(data.Description, requestBody.SetDescription)

	if err := convert.FrameworkToGraphStringSet(ctx, data.RoleScopeTagIds, requestBody.SetRoleScopeTagIds); err != nil {
		return nil, fmt.Errorf("failed to set role scope tags: %s", err)
	}

	convert.FrameworkToGraphBool(data.MicrosoftUpdateServiceAllowed, requestBody.SetMicrosoftUpdateServiceAllowed)
	convert.FrameworkToGraphBool(data.DriversExcluded, requestBody.SetDriversExcluded)
	convert.FrameworkToGraphInt32(data.QualityUpdatesDeferralPeriodInDays, requestBody.SetQualityUpdatesDeferralPeriodInDays)
	convert.FrameworkToGraphInt32(data.FeatureUpdatesDeferralPeriodInDays, requestBody.SetFeatureUpdatesDeferralPeriodInDays)
	convert.FrameworkToGraphBool(data.AllowWindows11Upgrade, requestBody.SetAllowWindows11Upgrade)
	convert.FrameworkToGraphBool(data.QualityUpdatesPaused, requestBody.SetQualityUpdatesPaused)
	convert.FrameworkToGraphBool(data.FeatureUpdatesPaused, requestBody.SetFeatureUpdatesPaused)

	if err := convert.FrameworkToGraphDateOnly(data.FeatureUpdatesPauseStartDate, requestBody.SetFeatureUpdatesPauseStartDate); err != nil {
		return nil, fmt.Errorf("error setting FeatureUpdatesPauseStartDate: %v", err)
	}

	if err := convert.FrameworkToGraphTime(data.FeatureUpdatesPauseExpiryDateTime, requestBody.SetFeatureUpdatesPauseExpiryDateTime); err != nil {
		return nil, fmt.Errorf("error setting FeatureUpdatesPauseExpiryDateTime: %v", err)
	}

	if err := convert.FrameworkToGraphTime(data.FeatureUpdatesRollbackStartDateTime, requestBody.SetFeatureUpdatesRollbackStartDateTime); err != nil {
		return nil, fmt.Errorf("error setting FeatureUpdatesRollbackStartDateTime: %v", err)
	}

	if err := convert.FrameworkToGraphDateOnly(data.QualityUpdatesPauseStartDate, requestBody.SetQualityUpdatesPauseStartDate); err != nil {
		return nil, fmt.Errorf("error setting QualityUpdatesPauseStartDate: %v", err)
	}

	if err := convert.FrameworkToGraphTime(data.QualityUpdatesPauseExpiryDateTime, requestBody.SetQualityUpdatesPauseExpiryDateTime); err != nil {
		return nil, fmt.Errorf("error setting QualityUpdatesPauseExpiryDateTime: %v", err)
	}

	if err := convert.FrameworkToGraphTime(data.QualityUpdatesRollbackStartDateTime, requestBody.SetQualityUpdatesRollbackStartDateTime); err != nil {
		return nil, fmt.Errorf("error setting QualityUpdatesRollbackStartDateTime: %v", err)
	}

	err := convert.FrameworkToGraphEnum(data.BusinessReadyUpdatesOnly, graphmodels.ParseWindowsUpdateType, requestBody.SetBusinessReadyUpdatesOnly)
	if err != nil {
		return nil, fmt.Errorf("error setting BusinessReadyUpdatesOnly: %v", err)
	}

	convert.FrameworkToGraphBool(data.SkipChecksBeforeRestart, requestBody.SetSkipChecksBeforeRestart)

	err = convert.FrameworkToGraphEnum(data.AutomaticUpdateMode, graphmodels.ParseAutomaticUpdateMode, requestBody.SetAutomaticUpdateMode)
	if err != nil {
		return nil, fmt.Errorf("error setting AutomaticUpdateMode: %v", err)
	}

	if !data.ActiveHoursStart.IsNull() && !data.ActiveHoursEnd.IsNull() &&
		data.ActiveHoursStart.ValueString() != "" && data.ActiveHoursEnd.ValueString() != "" {
		installationSchedule := graphmodels.NewWindowsUpdateActiveHoursInstall()

		if err := convert.FrameworkToGraphTimeOnly(data.ActiveHoursStart, installationSchedule.SetActiveHoursStart); err != nil {
			return nil, fmt.Errorf("error setting active hours start: %v", err)
		}

		if err := convert.FrameworkToGraphTimeOnly(data.ActiveHoursEnd, installationSchedule.SetActiveHoursEnd); err != nil {
			return nil, fmt.Errorf("error setting active hours end: %v", err)
		}
		requestBody.SetInstallationSchedule(installationSchedule)
	}

	err = convert.FrameworkToGraphEnum(data.UserPauseAccess, graphmodels.ParseEnablement, requestBody.SetUserPauseAccess)
	if err != nil {
		return nil, fmt.Errorf("error setting UserPauseAccess: %v", err)
	}

	err = convert.FrameworkToGraphEnum(data.UserWindowsUpdateScanAccess, graphmodels.ParseEnablement, requestBody.SetUserWindowsUpdateScanAccess)
	if err != nil {
		return nil, fmt.Errorf("error setting UserWindowsUpdateScanAccess: %v", err)
	}

	err = convert.FrameworkToGraphEnum(data.UpdateNotificationLevel, graphmodels.ParseWindowsUpdateNotificationDisplayOption, requestBody.SetUpdateNotificationLevel)
	if err != nil {
		return nil, fmt.Errorf("error setting UpdateNotificationLevel: %v", err)
	}

	err = convert.FrameworkToGraphEnum(data.UpdateWeeks, graphmodels.ParseWindowsUpdateForBusinessUpdateWeeks, requestBody.SetUpdateWeeks)
	if err != nil {
		return nil, fmt.Errorf("error setting UpdateWeeks: %v", err)
	}

	convert.FrameworkToGraphInt32(data.FeatureUpdatesRollbackWindowInDays, requestBody.SetFeatureUpdatesRollbackWindowInDays)

	// Handle deadline settings nested block
	if !data.DeadlineSettings.IsNull() && !data.DeadlineSettings.IsUnknown() {
		attrs := data.DeadlineSettings.Attributes()

		if deadlineForFeatureUpdatesInDays, ok := attrs["deadline_for_feature_updates_in_days"].(types.Int32); ok {
			convert.FrameworkToGraphInt32(deadlineForFeatureUpdatesInDays, requestBody.SetDeadlineForFeatureUpdatesInDays)
		}

		if deadlineForQualityUpdatesInDays, ok := attrs["deadline_for_quality_updates_in_days"].(types.Int32); ok {
			convert.FrameworkToGraphInt32(deadlineForQualityUpdatesInDays, requestBody.SetDeadlineForQualityUpdatesInDays)
		}

		if deadlineGracePeriodInDays, ok := attrs["deadline_grace_period_in_days"].(types.Int32); ok {
			convert.FrameworkToGraphInt32(deadlineGracePeriodInDays, requestBody.SetDeadlineGracePeriodInDays)
		}

		if postponeRebootUntilAfterDeadline, ok := attrs["postpone_reboot_until_after_deadline"].(types.Bool); ok {
			convert.FrameworkToGraphBool(postponeRebootUntilAfterDeadline, requestBody.SetPostponeRebootUntilAfterDeadline)
		}
	}
	convert.FrameworkToGraphInt32(data.EngagedRestartDeadlineInDays, requestBody.SetEngagedRestartDeadlineInDays)
	convert.FrameworkToGraphInt32(data.EngagedRestartSnoozeScheduleInDays, requestBody.SetEngagedRestartSnoozeScheduleInDays)
	convert.FrameworkToGraphInt32(data.EngagedRestartTransitionScheduleInDays, requestBody.SetEngagedRestartTransitionScheduleInDays)

	err = convert.FrameworkToGraphEnum(data.AutoRestartNotificationDismissal, graphmodels.ParseAutoRestartNotificationDismissalMethod, requestBody.SetAutoRestartNotificationDismissal)
	if err != nil {
		return nil, fmt.Errorf("error setting AutoRestartNotificationDismissal: %v", err)
	}

	convert.FrameworkToGraphInt32(data.ScheduleRestartWarningInHours, requestBody.SetScheduleRestartWarningInHours)
	convert.FrameworkToGraphInt32(data.ScheduleImminentRestartWarningInMinutes, requestBody.SetScheduleImminentRestartWarningInMinutes)

	err = convert.FrameworkToGraphEnum(data.DeliveryOptimizationMode, graphmodels.ParseWindowsDeliveryOptimizationMode, requestBody.SetDeliveryOptimizationMode)
	if err != nil {
		return nil, fmt.Errorf("error setting DeliveryOptimizationMode: %v", err)
	}

	err = convert.FrameworkToGraphEnum(data.PrereleaseFeatures, graphmodels.ParsePrereleaseFeatures, requestBody.SetPrereleaseFeatures)
	if err != nil {
		return nil, fmt.Errorf("error setting PrereleaseFeatures: %v", err)
	}

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]any{
			"error": err.Error(),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing %s resource", ResourceName))

	return requestBody, nil
}
