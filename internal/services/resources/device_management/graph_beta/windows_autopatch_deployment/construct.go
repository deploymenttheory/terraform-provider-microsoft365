package graphBetaWindowsAutopatchDeployment

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodelswindowsupdates "github.com/microsoftgraph/msgraph-beta-sdk-go/models/windowsupdates"
)

func constructResource(ctx context.Context, data *WindowsAutopatchDeploymentResourceModel) (graphmodelswindowsupdates.Deploymentable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s resource from model", ResourceName))

	requestBody := graphmodelswindowsupdates.NewDeployment()

	if data.Content != nil {
		content := graphmodelswindowsupdates.NewCatalogContent()

		catalogEntryType := data.Content.CatalogEntryType.ValueString()
		catalogEntryId := data.Content.CatalogEntryId.ValueString()

		switch catalogEntryType {
		case "featureUpdate":
			catalogEntry := graphmodelswindowsupdates.NewFeatureUpdateCatalogEntry()
			catalogEntry.SetId(&catalogEntryId)
			content.SetCatalogEntry(catalogEntry)
		case "qualityUpdate":
			catalogEntry := graphmodelswindowsupdates.NewQualityUpdateCatalogEntry()
			catalogEntry.SetId(&catalogEntryId)
			content.SetCatalogEntry(catalogEntry)
		default:
			return nil, fmt.Errorf("invalid catalog_entry_type: %s", catalogEntryType)
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

	audience := graphmodelswindowsupdates.NewDeploymentAudience()
	requestBody.SetAudience(audience)

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing %s resource", ResourceName))
	return requestBody, nil
}

func constructUpdateResource(ctx context.Context, plan *WindowsAutopatchDeploymentResourceModel, state *WindowsAutopatchDeploymentResourceModel) (graphmodelswindowsupdates.Deploymentable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing update request for %s resource", ResourceName))

	requestBody := graphmodelswindowsupdates.NewDeployment()

	if plan.Settings != nil {
		settings, err := constructDeploymentSettings(ctx, plan.Settings)
		if err != nil {
			return nil, fmt.Errorf("failed to construct deployment settings: %w", err)
		}
		requestBody.SetSettings(settings)
	}

	if !plan.State.IsNull() && !plan.State.IsUnknown() {
		stateAttrs := plan.State.Attributes()
		if requestedValue, ok := stateAttrs["requested_value"].(types.String); ok && !requestedValue.IsNull() {
			deploymentState := graphmodelswindowsupdates.NewDeploymentState()

			requestedValueStr := requestedValue.ValueString()
			var requestedValueEnum graphmodelswindowsupdates.RequestedDeploymentStateValue

			switch requestedValueStr {
			case "none":
				requestedValueEnum = graphmodelswindowsupdates.NONE_REQUESTEDDEPLOYMENTSTATEVALUE
			case "paused":
				requestedValueEnum = graphmodelswindowsupdates.PAUSED_REQUESTEDDEPLOYMENTSTATEVALUE
			case "archived":
				requestedValueEnum = graphmodelswindowsupdates.ARCHIVED_REQUESTEDDEPLOYMENTSTATEVALUE
			default:
				return nil, fmt.Errorf("invalid requested_value: %s", requestedValueStr)
			}

			deploymentState.SetRequestedValue(&requestedValueEnum)
			requestBody.SetState(deploymentState)
		}
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

				if err := convert.FrameworkToGraphISODuration(data.Schedule.GradualRollout.DurationBetweenOffers, rateDrivenRollout.SetDurationBetweenOffers); err != nil {
					return nil, fmt.Errorf("failed to parse duration_between_offers: %w", err)
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

		if len(data.Monitoring.MonitoringRules) > 0 {
			rules := make([]graphmodelswindowsupdates.MonitoringRuleable, 0, len(data.Monitoring.MonitoringRules))

			for _, ruleData := range data.Monitoring.MonitoringRules {
				rule := graphmodelswindowsupdates.NewMonitoringRule()

				signalStr := ruleData.Signal.ValueString()
				var signal graphmodelswindowsupdates.MonitoringSignal
				switch signalStr {
				case "rollback":
					signal = graphmodelswindowsupdates.ROLLBACK_MONITORINGSIGNAL
				default:
					return nil, fmt.Errorf("invalid monitoring signal: %s", signalStr)
				}
				rule.SetSignal(&signal)

				convert.FrameworkToGraphInt32(ruleData.Threshold, rule.SetThreshold)

				actionStr := ruleData.Action.ValueString()
				var action graphmodelswindowsupdates.MonitoringAction
				switch actionStr {
				case "pauseDeployment":
					action = graphmodelswindowsupdates.PAUSEDEPLOYMENT_MONITORINGACTION
				case "alertError":
					action = graphmodelswindowsupdates.ALERTERROR_MONITORINGACTION
				default:
					return nil, fmt.Errorf("invalid monitoring action: %s", actionStr)
				}
				rule.SetAction(&action)

				rules = append(rules, rule)
			}

			monitoring.SetMonitoringRules(rules)
		}

		settings.SetMonitoring(monitoring)
	}

	return settings, nil
}
