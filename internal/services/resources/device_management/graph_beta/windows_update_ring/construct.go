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

	err := convert.FrameworkToGraphEnum(data.BusinessReadyUpdatesOnly, graphmodels.ParseWindowsUpdateType, requestBody.SetBusinessReadyUpdatesOnly)
	if err != nil {
		return nil, fmt.Errorf("error setting BusinessReadyUpdatesOnly: %v", err)
	}

	convert.FrameworkToGraphBool(data.SkipChecksBeforeRestart, requestBody.SetSkipChecksBeforeRestart)

	err = convert.FrameworkToGraphEnum(data.AutomaticUpdateMode, graphmodels.ParseAutomaticUpdateMode, requestBody.SetAutomaticUpdateMode)
	if err != nil {
		return nil, fmt.Errorf("error setting AutomaticUpdateMode: %v", err)
	}

	if !data.ActiveHoursStart.IsNull() && !data.ActiveHoursEnd.IsNull() {
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
	convert.FrameworkToGraphInt32(data.DeadlineForFeatureUpdatesInDays, requestBody.SetDeadlineForFeatureUpdatesInDays)
	convert.FrameworkToGraphInt32(data.DeadlineForQualityUpdatesInDays, requestBody.SetDeadlineForQualityUpdatesInDays)
	convert.FrameworkToGraphInt32(data.DeadlineGracePeriodInDays, requestBody.SetDeadlineGracePeriodInDays)
	convert.FrameworkToGraphBool(data.PostponeRebootUntilAfterDeadline, requestBody.SetPostponeRebootUntilAfterDeadline)
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

	if len(data.AdditionalProperties) > 0 {
		requestBody.SetAdditionalData(data.AdditionalProperties)
	}

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]interface{}{
			"error": err.Error(),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing %s resource", ResourceName))

	return requestBody, nil
}
