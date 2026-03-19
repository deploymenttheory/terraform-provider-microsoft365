package graphBetaWindowsUpdatesAutopatchUpdatePolicy

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodelswindowsupdates "github.com/microsoftgraph/msgraph-beta-sdk-go/models/windowsupdates"
)

func constructResource(ctx context.Context, data *WindowsUpdatesAutopatchUpdatePolicyResourceModel, isUpdate bool) (graphmodelswindowsupdates.UpdatePolicyable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s resource from model", ResourceName))
	tflog.Debug(ctx, "Input data", map[string]any{
		"audience_id":         data.AudienceId.ValueString(),
		"compliance_changes":  data.ComplianceChanges.ValueBool(),
		"has_change_rules":    !data.ComplianceChangeRules.IsNull(),
		"has_deploy_settings": !data.DeploymentSettings.IsNull(),
	})

	requestBody := graphmodelswindowsupdates.NewUpdatePolicy()
	
	// Set @odata.type for updatePolicy
	odataType := "#microsoft.graph.windowsUpdates.updatePolicy"
	requestBody.SetOdataType(&odataType)

	// For UPDATE operations, don't send audience, complianceChanges, or complianceChangeRules
	// API testing shows that UPDATE only accepts deploymentSettings, despite docs saying otherwise
	if !isUpdate {
		if !data.AudienceId.IsNull() && !data.AudienceId.IsUnknown() {
			audience := graphmodelswindowsupdates.NewDeploymentAudience()
			convert.FrameworkToGraphString(data.AudienceId, audience.SetId)
			requestBody.SetAudience(audience)
			tflog.Debug(ctx, "Set audience", map[string]any{"id": data.AudienceId.ValueString()})
		}

		// Set complianceChanges (required for CREATE per API docs)
		if !data.ComplianceChanges.IsNull() && data.ComplianceChanges.ValueBool() {
			complianceChange := graphmodelswindowsupdates.NewContentApproval()
			complianceChanges := []graphmodelswindowsupdates.ComplianceChangeable{
				complianceChange,
			}
			requestBody.SetComplianceChanges(complianceChanges)
		}
	}

	// complianceChangeRules can only be set during CREATE, not UPDATE
	// API returns 400 "Schema validation failed for resource 'deploymentPolicy'" if included in UPDATE
	if !isUpdate && !data.ComplianceChangeRules.IsNull() && !data.ComplianceChangeRules.IsUnknown() {
		var rules []ComplianceChangeRuleModel
		data.ComplianceChangeRules.ElementsAs(ctx, &rules, false)

		complianceRules := make([]graphmodelswindowsupdates.ComplianceChangeRuleable, 0, len(rules))
		for _, rule := range rules {
			complianceRule := graphmodelswindowsupdates.NewContentApprovalRule()

			if rule.ContentFilter != nil && !rule.ContentFilter.FilterType.IsNull() {
				filterType := rule.ContentFilter.FilterType.ValueString()
				switch filterType {
				case "driverUpdateFilter":
					filter := graphmodelswindowsupdates.NewDriverUpdateFilter()
					complianceRule.SetContentFilter(filter)
				case "windowsUpdateFilter":
					filter := graphmodelswindowsupdates.NewWindowsUpdateFilter()
					complianceRule.SetContentFilter(filter)
				}
			}

			// Use AdditionalData to send the raw duration string rather than
			// SetDurationBeforeDeploymentStart, which goes through ISODuration.String()
			// and normalizes day-based durations to weeks (e.g. P7D → P1W).
			// The Graph API rejects week-based durations for this field.
			if !rule.DurationBeforeDeploymentStart.IsNull() {
				rawDuration := rule.DurationBeforeDeploymentStart.ValueString()
				if rawDuration != "" {
					additionalData := complianceRule.GetAdditionalData()
					if additionalData == nil {
						additionalData = make(map[string]any)
					}
					additionalData["durationBeforeDeploymentStart"] = rawDuration
					complianceRule.SetAdditionalData(additionalData)
				}
			}

			complianceRules = append(complianceRules, complianceRule)
		}
		requestBody.SetComplianceChangeRules(complianceRules)
	}

	// Set deployment settings
	if !data.DeploymentSettings.IsNull() && !data.DeploymentSettings.IsUnknown() {
		var deploymentSettingsData DeploymentSettingsModel
		diags := data.DeploymentSettings.As(ctx, &deploymentSettingsData, basetypes.ObjectAsOptions{})
		if diags.HasError() {
			return nil, fmt.Errorf("failed to extract deployment_settings data: %s", diags.Errors()[0].Detail())
		}

		settings := graphmodelswindowsupdates.NewDeploymentSettings()

		if !deploymentSettingsData.Schedule.IsNull() && !deploymentSettingsData.Schedule.IsUnknown() {
			var scheduleData ScheduleSettingsModel
			diags := deploymentSettingsData.Schedule.As(ctx, &scheduleData, basetypes.ObjectAsOptions{})
			if diags.HasError() {
				return nil, fmt.Errorf("failed to extract schedule data: %s", diags.Errors()[0].Detail())
			}

			schedule := graphmodelswindowsupdates.NewScheduleSettings()

			// Only set startDateTime if explicitly provided (not required for creation)
			if !scheduleData.StartDateTime.IsNull() && !scheduleData.StartDateTime.IsUnknown() {
				tflog.Debug(ctx, "Setting schedule start_date_time", map[string]any{
					"value": scheduleData.StartDateTime.ValueString(),
				})
				if err := convert.FrameworkToGraphTime(scheduleData.StartDateTime, schedule.SetStartDateTime); err != nil {
					return nil, fmt.Errorf("invalid start_date_time: %w", err)
				}
			}

			if !scheduleData.GradualRollout.IsNull() && !scheduleData.GradualRollout.IsUnknown() {
				var rolloutData GradualRolloutModel
				diags := scheduleData.GradualRollout.As(ctx, &rolloutData, basetypes.ObjectAsOptions{})
				if diags.HasError() {
					return nil, fmt.Errorf("failed to extract gradual_rollout data: %s", diags.Errors()[0].Detail())
				}

				rollout := graphmodelswindowsupdates.NewRateDrivenRolloutSettings()

				// Use AdditionalData to send raw values (avoid SDK normalization)
				additionalData := rollout.GetAdditionalData()
				if additionalData == nil {
					additionalData = make(map[string]any)
				}

				// Set duration between offers as raw string (avoid P7D → P1W normalization)
				if !rolloutData.DurationBetweenOffers.IsNull() {
					rawDuration := rolloutData.DurationBetweenOffers.ValueString()
					if rawDuration != "" {
						additionalData["durationBetweenOffers"] = rawDuration
					}
				}

				// Set devicesPerOffer via additionalData (note: plural "devices", not singular as shown in some MS docs)
				if !rolloutData.DevicesPerOffer.IsNull() {
					additionalData["devicesPerOffer"] = rolloutData.DevicesPerOffer.ValueInt32()
				}

				rollout.SetAdditionalData(additionalData)

				schedule.SetGradualRollout(rollout)
			}

			settings.SetSchedule(schedule)
		}

		requestBody.SetDeploymentSettings(settings)
	}

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]any{
			"error": err.Error(),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing deployment settings for %s resource", ResourceName))
	return requestBody, nil
}
