package graphBetaWindowsAutopatchDeployment

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodelswindowsupdates "github.com/microsoftgraph/msgraph-beta-sdk-go/models/windowsupdates"
)

func MapRemoteStateToTerraform(ctx context.Context, data *WindowsAutopatchDeploymentResourceModel, remoteResource graphmodelswindowsupdates.Deploymentable) {
	tflog.Debug(ctx, fmt.Sprintf("Mapping remote state to Terraform state for %s", ResourceName))

	if remoteResource == nil {
		tflog.Debug(ctx, "Remote resource is nil")
		return
	}

	data.ID = convert.GraphToFrameworkString(remoteResource.GetId())
	data.CreatedDateTime = convert.GraphToFrameworkTime(remoteResource.GetCreatedDateTime())
	data.LastModifiedDateTime = convert.GraphToFrameworkTime(remoteResource.GetLastModifiedDateTime())

	if content := remoteResource.GetContent(); content != nil {
		if catalogContent, ok := content.(graphmodelswindowsupdates.CatalogContentable); ok {
			if data.Content == nil {
				data.Content = &DeploymentContent{}
			}

			if catalogEntry := catalogContent.GetCatalogEntry(); catalogEntry != nil {
				data.Content.CatalogEntryId = convert.GraphToFrameworkString(catalogEntry.GetId())

				switch catalogEntry.(type) {
				case graphmodelswindowsupdates.FeatureUpdateCatalogEntryable:
					data.Content.CatalogEntryType = types.StringValue("featureUpdate")
				case graphmodelswindowsupdates.QualityUpdateCatalogEntryable:
					data.Content.CatalogEntryType = types.StringValue("qualityUpdate")
				}
			}
		}
	}

	if settings := remoteResource.GetSettings(); settings != nil {
		if data.Settings == nil {
			data.Settings = &DeploymentSettings{}
		}

		if schedule := settings.GetSchedule(); schedule != nil {
			if data.Settings.Schedule == nil {
				data.Settings.Schedule = &ScheduleSettings{}
			}

			data.Settings.Schedule.StartDateTime = convert.GraphToFrameworkTime(schedule.GetStartDateTime())

			if gradualRollout := schedule.GetGradualRollout(); gradualRollout != nil {
				if data.Settings.Schedule.GradualRollout == nil {
					data.Settings.Schedule.GradualRollout = &GradualRollout{}
				}

				if rateDriven, ok := gradualRollout.(graphmodelswindowsupdates.RateDrivenRolloutSettingsable); ok {
					data.Settings.Schedule.GradualRollout.DurationBetweenOffers = convert.GraphToFrameworkISODuration(rateDriven.GetDurationBetweenOffers())
					data.Settings.Schedule.GradualRollout.DevicesPerOffer = convert.GraphToFrameworkInt32(rateDriven.GetDevicesPerOffer())
				} else if dateDriven, ok := gradualRollout.(graphmodelswindowsupdates.DateDrivenRolloutSettingsable); ok {
					data.Settings.Schedule.GradualRollout.EndDateTime = convert.GraphToFrameworkTime(dateDriven.GetEndDateTime())
				}
			}
		}

		if monitoring := settings.GetMonitoring(); monitoring != nil {
			if data.Settings.Monitoring == nil {
				data.Settings.Monitoring = &MonitoringSettings{}
			}

			if rules := monitoring.GetMonitoringRules(); len(rules) > 0 {
				data.Settings.Monitoring.MonitoringRules = make([]MonitoringRule, 0, len(rules))

				for _, rule := range rules {
					if rule == nil {
						continue
					}

					tfRule := MonitoringRule{}

					if signal := rule.GetSignal(); signal != nil {
						tfRule.Signal = types.StringValue(signal.String())
					}

					tfRule.Threshold = convert.GraphToFrameworkInt32(rule.GetThreshold())

					if action := rule.GetAction(); action != nil {
						tfRule.Action = types.StringValue(action.String())
					}

					data.Settings.Monitoring.MonitoringRules = append(data.Settings.Monitoring.MonitoringRules, tfRule)
				}
			}
		}
	}

	if deploymentState := remoteResource.GetState(); deploymentState != nil {
		stateAttrs := map[string]attr.Value{
			"requested_value": convert.GraphToFrameworkEnum(deploymentState.GetRequestedValue()),
			"effective_value": convert.GraphToFrameworkEnum(deploymentState.GetEffectiveValue()),
		}

		stateObj, diags := types.ObjectValue(
			map[string]attr.Type{
				"requested_value": types.StringType,
				"effective_value": types.StringType,
			},
			stateAttrs,
		)

		if diags.HasError() {
			tflog.Error(ctx, "Failed to create state object", map[string]interface{}{
				"diagnostics": diags,
			})
		} else {
			data.State = stateObj
		}
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished mapping remote state to Terraform state for %s", ResourceName))
}
