package graphBetaWindowsUpdatesAutopatchDeployment

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodelswindowsupdates "github.com/microsoftgraph/msgraph-beta-sdk-go/models/windowsupdates"
)

func constructResource(ctx context.Context, data *WindowsUpdatesAutopatchDeploymentResourceModel) (graphmodelswindowsupdates.Deploymentable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s resource from model", ResourceName))

	requestBody := graphmodelswindowsupdates.NewDeployment()

	if data.Content != nil {
		content := graphmodelswindowsupdates.NewCatalogContent()

		catalogEntryId := data.Content.CatalogEntryId.ValueString()

		switch data.Content.CatalogEntryType.ValueString() {
		case "featureUpdate":
			catalogEntry := graphmodelswindowsupdates.NewFeatureUpdateCatalogEntry()
			catalogEntry.SetId(&catalogEntryId)
			content.SetCatalogEntry(catalogEntry)
		case "qualityUpdate":
			catalogEntry := graphmodelswindowsupdates.NewQualityUpdateCatalogEntry()
			catalogEntry.SetId(&catalogEntryId)
			content.SetCatalogEntry(catalogEntry)
		default:
			return nil, fmt.Errorf("invalid catalog_entry_type: %s", data.Content.CatalogEntryType.ValueString())
		}

		requestBody.SetContent(content)
	}

	if data.Settings != nil {
		settings, err := constructDeploymentSettings(ctx, data.Settings)
		if err != nil {
			return nil, fmt.Errorf("failed to construct deployment settings: %w", err)
		}
		requestBody.SetSettings(settings)
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing %s resource", ResourceName))
	return requestBody, nil
}

func constructUpdateResource(ctx context.Context, plan *WindowsUpdatesAutopatchDeploymentResourceModel) (graphmodelswindowsupdates.Deploymentable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing update request for %s resource", ResourceName))

	requestBody := graphmodelswindowsupdates.NewDeployment()

	if plan.Settings != nil {
		settings, err := constructDeploymentSettings(ctx, plan.Settings)
		if err != nil {
			return nil, fmt.Errorf("failed to construct deployment settings: %w", err)
		}
		requestBody.SetSettings(settings)
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing update request for %s resource", ResourceName))
	return requestBody, nil
}

func constructDeploymentSettings(ctx context.Context, data *DeploymentSettings) (graphmodelswindowsupdates.DeploymentSettingsable, error) {
	settings := graphmodelswindowsupdates.NewDeploymentSettings()

	if data.Schedule != nil {
		schedule := graphmodelswindowsupdates.NewScheduleSettings()

		if err := convert.FrameworkToGraphTime(data.Schedule.StartDateTime, schedule.SetStartDateTime); err != nil {
			return nil, fmt.Errorf("failed to parse start_date_time: %w", err)
		}

		if data.Schedule.GradualRollout != nil {
			if !data.Schedule.GradualRollout.DurationBetweenOffers.IsNull() && !data.Schedule.GradualRollout.DevicesPerOffer.IsNull() {
				rateDrivenRollout := graphmodelswindowsupdates.NewRateDrivenRolloutSettings()

				// Use AdditionalData to send the raw duration string rather than
				// SetDurationBetweenOffers, which goes through ISODuration.String()
				// and normalizes day-based durations to weeks (e.g. P7D → P1W).
				// The Graph API rejects week-based durations for this field.
				rawDuration := data.Schedule.GradualRollout.DurationBetweenOffers.ValueString()
				if rawDuration != "" {
					rateDrivenRollout.GetAdditionalData()["durationBetweenOffers"] = rawDuration
				}

				convert.FrameworkToGraphInt32(data.Schedule.GradualRollout.DevicesPerOffer, rateDrivenRollout.SetDevicesPerOffer)

				schedule.SetGradualRollout(rateDrivenRollout)
			} else if !data.Schedule.GradualRollout.EndDateTime.IsNull() {
				dateDrivenRollout := graphmodelswindowsupdates.NewDateDrivenRolloutSettings()

				if err := convert.FrameworkToGraphTime(data.Schedule.GradualRollout.EndDateTime, dateDrivenRollout.SetEndDateTime); err != nil {
					return nil, fmt.Errorf("failed to parse gradual_rollout end_date_time: %w", err)
				}

				schedule.SetGradualRollout(dateDrivenRollout)
			}
		}

		settings.SetSchedule(schedule)
	}

	if data.Monitoring != nil {
		monitoring := graphmodelswindowsupdates.NewMonitoringSettings()

		if !data.Monitoring.MonitoringRules.IsNull() && !data.Monitoring.MonitoringRules.IsUnknown() && len(data.Monitoring.MonitoringRules.Elements()) > 0 {
			var ruleModels []MonitoringRule
			diags := data.Monitoring.MonitoringRules.ElementsAs(ctx, &ruleModels, false)
			if diags.HasError() {
				return nil, fmt.Errorf("failed to extract monitoring_rules: %s", diags.Errors())
			}

			rules := make([]graphmodelswindowsupdates.MonitoringRuleable, 0, len(ruleModels))

			for _, ruleData := range ruleModels {
				rule := graphmodelswindowsupdates.NewMonitoringRule()

				if err := convert.FrameworkToGraphEnum(
					ruleData.Signal,
					graphmodelswindowsupdates.ParseMonitoringSignal,
					rule.SetSignal,
				); err != nil {
					return nil, fmt.Errorf("failed to parse monitoring signal: %w", err)
				}

				convert.FrameworkToGraphInt32(ruleData.Threshold, rule.SetThreshold)

				if err := convert.FrameworkToGraphEnum(
					ruleData.Action,
					graphmodelswindowsupdates.ParseMonitoringAction,
					rule.SetAction,
				); err != nil {
					return nil, fmt.Errorf("failed to parse monitoring action: %w", err)
				}

				rules = append(rules, rule)
			}

			monitoring.SetMonitoringRules(rules)
		}

		settings.SetMonitoring(monitoring)
	}

	return settings, nil
}
