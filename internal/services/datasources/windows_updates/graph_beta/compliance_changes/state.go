package graphBetaWindowsUpdatesComplianceChanges

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodelswindowsupdates "github.com/microsoftgraph/msgraph-beta-sdk-go/models/windowsupdates"
)

func MapRemoteStateToDataSource(ctx context.Context, remoteChange graphmodelswindowsupdates.ComplianceChangeable) ComplianceChange {
	tflog.Debug(ctx, fmt.Sprintf("Mapping remote state to data source for %s", DataSourceName))

	change := ComplianceChange{
		ID:              convert.GraphToFrameworkString(remoteChange.GetId()),
		CreatedDateTime: convert.GraphToFrameworkTime(remoteChange.GetCreatedDateTime()),
		IsRevoked:       convert.GraphToFrameworkBool(remoteChange.GetIsRevoked()),
		RevokedDateTime: convert.GraphToFrameworkTime(remoteChange.GetRevokedDateTime()),
	}

	if contentApproval, ok := remoteChange.(graphmodelswindowsupdates.ContentApprovalable); ok {
		if content := contentApproval.GetContent(); content != nil {
			if catalogContent, ok := content.(graphmodelswindowsupdates.CatalogContentable); ok {
				change.Content = &ComplianceContent{}

				if catalogEntry := catalogContent.GetCatalogEntry(); catalogEntry != nil {
					change.Content.CatalogEntryId = convert.GraphToFrameworkString(catalogEntry.GetId())

					odataType := catalogEntry.GetOdataType()
					if odataType != nil {
						switch *odataType {
						case "#microsoft.graph.windowsUpdates.featureUpdateCatalogEntry":
							change.Content.CatalogEntryType = types.StringValue("featureUpdate")
						case "#microsoft.graph.windowsUpdates.qualityUpdateCatalogEntry":
							change.Content.CatalogEntryType = types.StringValue("qualityUpdate")
						case "#microsoft.graph.windowsUpdates.driverUpdateCatalogEntry":
							change.Content.CatalogEntryType = types.StringValue("driverUpdate")
						}
					}
				}
			}
		}

		if deploymentSettings := contentApproval.GetDeploymentSettings(); deploymentSettings != nil {
			change.DeploymentSettings = &DeploymentSettings{}

			if schedule := deploymentSettings.GetSchedule(); schedule != nil {
				change.DeploymentSettings.Schedule = &ScheduleSettings{
					StartDateTime: convert.GraphToFrameworkTime(schedule.GetStartDateTime()),
				}

				if gradualRollout := schedule.GetGradualRollout(); gradualRollout != nil {
					if rateDriven, ok := gradualRollout.(graphmodelswindowsupdates.RateDrivenRolloutSettingsable); ok {
						change.DeploymentSettings.Schedule.GradualRollout = &GradualRollout{
							DevicesPerOffer: convert.GraphToFrameworkInt32(rateDriven.GetDevicesPerOffer()),
						}

						if additionalData := rateDriven.GetAdditionalData(); additionalData != nil {
							if duration, ok := additionalData["durationBetweenOffers"].(string); ok {
								change.DeploymentSettings.Schedule.GradualRollout.DurationBetweenOffers = convert.GraphToFrameworkString(&duration)
							}
						}
					}
				}
			}
		}
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished mapping remote state to data source for %s", DataSourceName))
	return change
}
