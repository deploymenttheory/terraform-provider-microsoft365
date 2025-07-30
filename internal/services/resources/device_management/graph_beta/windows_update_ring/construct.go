package graphBetaWindowsUpdateRing

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
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
	if data.DeadlineSettings != nil {
		convert.FrameworkToGraphInt32(data.DeadlineSettings.DeadlineForFeatureUpdatesInDays, requestBody.SetDeadlineForFeatureUpdatesInDays)
		convert.FrameworkToGraphInt32(data.DeadlineSettings.DeadlineForQualityUpdatesInDays, requestBody.SetDeadlineForQualityUpdatesInDays)
		convert.FrameworkToGraphInt32(data.DeadlineSettings.DeadlineGracePeriodInDays, requestBody.SetDeadlineGracePeriodInDays)
		convert.FrameworkToGraphBool(data.DeadlineSettings.PostponeRebootUntilAfterDeadline, requestBody.SetPostponeRebootUntilAfterDeadline)
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
		tflog.Error(ctx, "Failed to debug log object", map[string]interface{}{
			"error": err.Error(),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing %s resource", ResourceName))

	return requestBody, nil
}

// constructFeatureUpdateRollBack creates a request body for feature update rollback settings
func constructFeatureUpdateRollBack(ctx context.Context, data *WindowsUpdateRingResourceModel) (graphmodels.WindowsUpdateForBusinessConfigurationable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s feature update rollback settings", ResourceName))

	requestBody := graphmodels.NewWindowsUpdateForBusinessConfiguration()

	// Only set the rollback setting if uninstall settings exist
	if data.UninstallSettings != nil {
		convert.FrameworkToGraphBool(data.UninstallSettings.FeatureUpdatesWillBeRolledBack, requestBody.SetFeatureUpdatesWillBeRolledBack)
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing %s feature update rollback settings", ResourceName))
	return requestBody, nil
}

// constructQualityUpdateRollBack creates a request body for quality update rollback settings
func constructQualityUpdateRollBack(ctx context.Context, data *WindowsUpdateRingResourceModel) (graphmodels.WindowsUpdateForBusinessConfigurationable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s quality update rollback settings", ResourceName))

	requestBody := graphmodels.NewWindowsUpdateForBusinessConfiguration()

	// Only set the rollback setting if uninstall settings exist
	if data.UninstallSettings != nil {
		convert.FrameworkToGraphBool(data.UninstallSettings.QualityUpdatesWillBeRolledBack, requestBody.SetQualityUpdatesWillBeRolledBack)
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing %s quality update rollback settings", ResourceName))
	return requestBody, nil
}

// constructFeatureUpdatesPause creates a request body for feature updates pause/resume operations
func constructFeatureUpdatesPause(ctx context.Context, pause bool) (graphmodels.WindowsUpdateForBusinessConfigurationable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s feature updates pause operation (pause: %t)", ResourceName, pause))

	requestBody := graphmodels.NewWindowsUpdateForBusinessConfiguration()
	requestBody.SetFeatureUpdatesPaused(&pause)

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing %s feature updates pause operation", ResourceName))
	return requestBody, nil
}

// constructQualityUpdatesPause creates a request body for quality updates pause/resume operations
func constructQualityUpdatesPause(ctx context.Context, pause bool) (graphmodels.WindowsUpdateForBusinessConfigurationable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s quality updates pause operation (pause: %t)", ResourceName, pause))

	requestBody := graphmodels.NewWindowsUpdateForBusinessConfiguration()
	requestBody.SetQualityUpdatesPaused(&pause)

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing %s quality updates pause operation", ResourceName))
	return requestBody, nil
}

// constructFeatureUpdatesUninstall creates a request body for feature updates uninstall operation
func constructFeatureUpdatesUninstall(ctx context.Context, uninstall bool) (graphmodels.WindowsUpdateForBusinessConfigurationable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s feature updates uninstall operation (uninstall: %t)", ResourceName, uninstall))

	requestBody := graphmodels.NewWindowsUpdateForBusinessConfiguration()
	requestBody.SetFeatureUpdatesWillBeRolledBack(&uninstall)

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing %s feature updates uninstall operation", ResourceName))
	return requestBody, nil
}

// constructQualityUpdatesUninstall creates a request body for quality updates uninstall operation
func constructQualityUpdatesUninstall(ctx context.Context, uninstall bool) (graphmodels.WindowsUpdateForBusinessConfigurationable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s quality updates uninstall operation (uninstall: %t)", ResourceName, uninstall))

	requestBody := graphmodels.NewWindowsUpdateForBusinessConfiguration()
	requestBody.SetQualityUpdatesWillBeRolledBack(&uninstall)

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing %s quality updates uninstall operation", ResourceName))
	return requestBody, nil
}
