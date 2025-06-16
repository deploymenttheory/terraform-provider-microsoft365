package graphBetaWindowsUpdateRing

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// constructResource creates an assign request body with assignments from the nested blocks
func constructResource(ctx context.Context, data *WindowsUpdateRingResourceModel) (graphmodels.WindowsUpdateForBusinessConfigurationable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s resource from model", ResourceName))

	requestBody := graphmodels.NewWindowsUpdateForBusinessConfiguration()

	constructors.SetStringProperty(data.DisplayName, requestBody.SetDisplayName)
	constructors.SetStringProperty(data.Description, requestBody.SetDescription)

	if err := constructors.SetStringSet(ctx, data.RoleScopeTagIds, requestBody.SetRoleScopeTagIds); err != nil {
		return nil, fmt.Errorf("failed to set role scope tags: %s", err)
	}

	constructors.SetBoolProperty(data.MicrosoftUpdateServiceAllowed, requestBody.SetMicrosoftUpdateServiceAllowed)
	constructors.SetBoolProperty(data.DriversExcluded, requestBody.SetDriversExcluded)
	constructors.SetInt32Property(data.QualityUpdatesDeferralPeriodInDays, requestBody.SetQualityUpdatesDeferralPeriodInDays)
	constructors.SetInt32Property(data.FeatureUpdatesDeferralPeriodInDays, requestBody.SetFeatureUpdatesDeferralPeriodInDays)
	constructors.SetBoolProperty(data.AllowWindows11Upgrade, requestBody.SetAllowWindows11Upgrade)
	constructors.SetBoolProperty(data.QualityUpdatesPaused, requestBody.SetQualityUpdatesPaused)
	constructors.SetBoolProperty(data.FeatureUpdatesPaused, requestBody.SetFeatureUpdatesPaused)

	err := constructors.SetEnumProperty(data.BusinessReadyUpdatesOnly, graphmodels.ParseWindowsUpdateType, requestBody.SetBusinessReadyUpdatesOnly)
	if err != nil {
		return nil, fmt.Errorf("error setting BusinessReadyUpdatesOnly: %v", err)
	}

	constructors.SetBoolProperty(data.SkipChecksBeforeRestart, requestBody.SetSkipChecksBeforeRestart)

	err = constructors.SetEnumProperty(data.AutomaticUpdateMode, graphmodels.ParseAutomaticUpdateMode, requestBody.SetAutomaticUpdateMode)
	if err != nil {
		return nil, fmt.Errorf("error setting AutomaticUpdateMode: %v", err)
	}

	if !data.ActiveHoursStart.IsNull() && !data.ActiveHoursEnd.IsNull() {
		installationSchedule := graphmodels.NewWindowsUpdateActiveHoursInstall()

		if err := constructors.StringToTimeOnly(data.ActiveHoursStart, installationSchedule.SetActiveHoursStart); err != nil {
			return nil, fmt.Errorf("error setting active hours start: %v", err)
		}

		if err := constructors.StringToTimeOnly(data.ActiveHoursEnd, installationSchedule.SetActiveHoursEnd); err != nil {
			return nil, fmt.Errorf("error setting active hours end: %v", err)
		}
		requestBody.SetInstallationSchedule(installationSchedule)
	}

	err = constructors.SetEnumProperty(data.UserPauseAccess, graphmodels.ParseEnablement, requestBody.SetUserPauseAccess)
	if err != nil {
		return nil, fmt.Errorf("error setting UserPauseAccess: %v", err)
	}

	err = constructors.SetEnumProperty(data.UserWindowsUpdateScanAccess, graphmodels.ParseEnablement, requestBody.SetUserWindowsUpdateScanAccess)
	if err != nil {
		return nil, fmt.Errorf("error setting UserWindowsUpdateScanAccess: %v", err)
	}

	err = constructors.SetEnumProperty(data.UpdateNotificationLevel, graphmodels.ParseWindowsUpdateNotificationDisplayOption, requestBody.SetUpdateNotificationLevel)
	if err != nil {
		return nil, fmt.Errorf("error setting UpdateNotificationLevel: %v", err)
	}

	err = constructors.SetEnumProperty(data.UpdateWeeks, graphmodels.ParseWindowsUpdateForBusinessUpdateWeeks, requestBody.SetUpdateWeeks)
	if err != nil {
		return nil, fmt.Errorf("error setting UpdateWeeks: %v", err)
	}

	constructors.SetInt32Property(data.FeatureUpdatesRollbackWindowInDays, requestBody.SetFeatureUpdatesRollbackWindowInDays)
	constructors.SetInt32Property(data.DeadlineForFeatureUpdatesInDays, requestBody.SetDeadlineForFeatureUpdatesInDays)
	constructors.SetInt32Property(data.DeadlineForQualityUpdatesInDays, requestBody.SetDeadlineForQualityUpdatesInDays)
	constructors.SetInt32Property(data.DeadlineGracePeriodInDays, requestBody.SetDeadlineGracePeriodInDays)
	constructors.SetBoolProperty(data.PostponeRebootUntilAfterDeadline, requestBody.SetPostponeRebootUntilAfterDeadline)
	constructors.SetInt32Property(data.EngagedRestartDeadlineInDays, requestBody.SetEngagedRestartDeadlineInDays)
	constructors.SetInt32Property(data.EngagedRestartSnoozeScheduleInDays, requestBody.SetEngagedRestartSnoozeScheduleInDays)
	constructors.SetInt32Property(data.EngagedRestartTransitionScheduleInDays, requestBody.SetEngagedRestartTransitionScheduleInDays)

	err = constructors.SetEnumProperty(data.AutoRestartNotificationDismissal, graphmodels.ParseAutoRestartNotificationDismissalMethod, requestBody.SetAutoRestartNotificationDismissal)
	if err != nil {
		return nil, fmt.Errorf("error setting AutoRestartNotificationDismissal: %v", err)
	}

	constructors.SetInt32Property(data.ScheduleRestartWarningInHours, requestBody.SetScheduleRestartWarningInHours)
	constructors.SetInt32Property(data.ScheduleImminentRestartWarningInMinutes, requestBody.SetScheduleImminentRestartWarningInMinutes)

	err = constructors.SetEnumProperty(data.DeliveryOptimizationMode, graphmodels.ParseWindowsDeliveryOptimizationMode, requestBody.SetDeliveryOptimizationMode)
	if err != nil {
		return nil, fmt.Errorf("error setting DeliveryOptimizationMode: %v", err)
	}

	err = constructors.SetEnumProperty(data.PrereleaseFeatures, graphmodels.ParsePrereleaseFeatures, requestBody.SetPrereleaseFeatures)
	if err != nil {
		return nil, fmt.Errorf("error setting PrereleaseFeatures: %v", err)
	}

	if data.AdditionalProperties != nil && len(data.AdditionalProperties) > 0 {
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
