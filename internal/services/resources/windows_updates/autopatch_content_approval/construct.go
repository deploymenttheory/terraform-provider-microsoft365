package graphBetaWindowsUpdatesAutopatchContentApproval

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodelswindowsupdates "github.com/microsoftgraph/msgraph-beta-sdk-go/models/windowsupdates"
)

func constructResource(ctx context.Context, data *WindowsUpdatesAutopatchContentApprovalResourceModel) (graphmodelswindowsupdates.ComplianceChangeable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s resource from model", ResourceName))

	requestBody := graphmodelswindowsupdates.NewContentApproval()

	content := graphmodelswindowsupdates.NewCatalogContent()

	catalogEntryType := data.CatalogEntryType.ValueString()
	catalogEntryId := data.CatalogEntryId.ValueString()

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

	if data.DeploymentSettings != nil {
		deploymentSettings := graphmodelswindowsupdates.NewDeploymentSettings()

		if data.DeploymentSettings.Schedule != nil {
			schedule := graphmodelswindowsupdates.NewScheduleSettings()

			if err := convert.FrameworkToGraphTime(data.DeploymentSettings.Schedule.StartDateTime, schedule.SetStartDateTime); err != nil {
				return nil, fmt.Errorf("failed to parse start_date_time: %w", err)
			}

			if data.DeploymentSettings.Schedule.GradualRollout != nil {
				if !data.DeploymentSettings.Schedule.GradualRollout.EndDateTime.IsNull() {
					dateDrivenRollout := graphmodelswindowsupdates.NewDateDrivenRolloutSettings()

					if err := convert.FrameworkToGraphTime(data.DeploymentSettings.Schedule.GradualRollout.EndDateTime, dateDrivenRollout.SetEndDateTime); err != nil {
						return nil, fmt.Errorf("failed to parse gradual_rollout end_date_time: %w", err)
					}

					schedule.SetGradualRollout(dateDrivenRollout)
				}
			}

			deploymentSettings.SetSchedule(schedule)
		}

		requestBody.SetDeploymentSettings(deploymentSettings)
	}

	if !data.IsRevoked.IsNull() {
		isRevoked := data.IsRevoked.ValueBool()
		requestBody.SetIsRevoked(&isRevoked)
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing %s resource", ResourceName))
	return requestBody, nil
}

func constructUpdateResource(ctx context.Context, plan *WindowsUpdatesAutopatchContentApprovalResourceModel, state *WindowsUpdatesAutopatchContentApprovalResourceModel) (graphmodelswindowsupdates.ComplianceChangeable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing update request for %s resource", ResourceName))

	requestBody := graphmodelswindowsupdates.NewContentApproval()

	if !plan.IsRevoked.Equal(state.IsRevoked) {
		isRevoked := plan.IsRevoked.ValueBool()
		requestBody.SetIsRevoked(&isRevoked)
	}

	if plan.DeploymentSettings != nil {
		deploymentSettings := graphmodelswindowsupdates.NewDeploymentSettings()

		if plan.DeploymentSettings.Schedule != nil {
			schedule := graphmodelswindowsupdates.NewScheduleSettings()

			if err := convert.FrameworkToGraphTime(plan.DeploymentSettings.Schedule.StartDateTime, schedule.SetStartDateTime); err != nil {
				return nil, fmt.Errorf("failed to parse start_date_time: %w", err)
			}

			if plan.DeploymentSettings.Schedule.GradualRollout != nil {
				if !plan.DeploymentSettings.Schedule.GradualRollout.EndDateTime.IsNull() {
					dateDrivenRollout := graphmodelswindowsupdates.NewDateDrivenRolloutSettings()

					if err := convert.FrameworkToGraphTime(plan.DeploymentSettings.Schedule.GradualRollout.EndDateTime, dateDrivenRollout.SetEndDateTime); err != nil {
						return nil, fmt.Errorf("failed to parse gradual_rollout end_date_time: %w", err)
					}

					schedule.SetGradualRollout(dateDrivenRollout)
				}
			}

			deploymentSettings.SetSchedule(schedule)
		}

		requestBody.SetDeploymentSettings(deploymentSettings)
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing update request for %s resource", ResourceName))
	return requestBody, nil
}
