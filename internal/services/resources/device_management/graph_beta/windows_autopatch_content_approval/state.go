package graphBetaWindowsAutopatchContentApproval

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodelswindowsupdates "github.com/microsoftgraph/msgraph-beta-sdk-go/models/windowsupdates"
)

func MapRemoteStateToTerraform(ctx context.Context, data *WindowsAutopatchContentApprovalResourceModel, remoteResource graphmodelswindowsupdates.ComplianceChangeable) {
	tflog.Debug(ctx, fmt.Sprintf("Mapping remote state to Terraform state for %s", ResourceName))

	if remoteResource == nil {
		tflog.Debug(ctx, "Remote resource is nil")
		return
	}

	contentApproval, ok := remoteResource.(graphmodelswindowsupdates.ContentApprovalable)
	if !ok {
		tflog.Warn(ctx, "Remote resource is not a ContentApproval type")
		return
	}

	data.ID = convert.GraphToFrameworkString(contentApproval.GetId())
	data.CreatedDateTime = convert.GraphToFrameworkTime(contentApproval.GetCreatedDateTime())
	data.IsRevoked = convert.GraphToFrameworkBool(contentApproval.GetIsRevoked())
	data.RevokedDateTime = convert.GraphToFrameworkTime(contentApproval.GetRevokedDateTime())

	if content := contentApproval.GetContent(); content != nil {
		if catalogContent, ok := content.(graphmodelswindowsupdates.CatalogContentable); ok {
			if catalogEntry := catalogContent.GetCatalogEntry(); catalogEntry != nil {
				data.CatalogEntryId = convert.GraphToFrameworkString(catalogEntry.GetId())

				switch catalogEntry.(type) {
				case graphmodelswindowsupdates.FeatureUpdateCatalogEntryable:
					data.CatalogEntryType = types.StringValue("featureUpdate")
				case graphmodelswindowsupdates.QualityUpdateCatalogEntryable:
					data.CatalogEntryType = types.StringValue("qualityUpdate")
				}
			}
		}
	}

	if deploymentSettings := contentApproval.GetDeploymentSettings(); deploymentSettings != nil {
		data.DeploymentSettings = &DeploymentSettings{}

		if schedule := deploymentSettings.GetSchedule(); schedule != nil {
			data.DeploymentSettings.Schedule = &Schedule{}
			data.DeploymentSettings.Schedule.StartDateTime = convert.GraphToFrameworkTime(schedule.GetStartDateTime())

			if gradualRollout := schedule.GetGradualRollout(); gradualRollout != nil {
				if dateDrivenRollout, ok := gradualRollout.(graphmodelswindowsupdates.DateDrivenRolloutSettingsable); ok {
					data.DeploymentSettings.Schedule.GradualRollout = &GradualRollout{}
					data.DeploymentSettings.Schedule.GradualRollout.EndDateTime = convert.GraphToFrameworkTime(dateDrivenRollout.GetEndDateTime())
				}
			}
		}
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished mapping remote state to Terraform state for %s", ResourceName))
}
