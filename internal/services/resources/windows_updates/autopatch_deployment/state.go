package graphBetaWindowsUpdatesAutopatchDeployment

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/sentinels"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodelswindowsupdates "github.com/microsoftgraph/msgraph-beta-sdk-go/models/windowsupdates"
)

func MapRemoteStateToTerraform(ctx context.Context, data *WindowsUpdatesAutopatchDeploymentResourceModel, remoteResource graphmodelswindowsupdates.Deploymentable) {
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

			monitoringRulesSet, err := mapMonitoringRulesToState(monitoring.GetMonitoringRules())
			if err != nil {
				tflog.Error(ctx, "Failed to map monitoring rules to state", map[string]any{
					"error": err.Error(),
				})
				data.Settings.Monitoring.MonitoringRules = types.SetNull(monitoringRuleAttrType())
			} else {
				data.Settings.Monitoring.MonitoringRules = monitoringRulesSet
			}
		}
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished mapping remote state to Terraform state for %s", ResourceName))
}

// monitoringRuleAttrType returns the object type for MonitoringRule.
func monitoringRuleAttrType() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"signal":    types.StringType,
			"threshold": types.Int32Type,
			"action":    types.StringType,
		},
	}
}

// mapMonitoringRulesToState maps the SDK monitoring rules collection to a Terraform types.Set.
func mapMonitoringRulesToState(rules []graphmodelswindowsupdates.MonitoringRuleable) (types.Set, error) {
	ruleType := monitoringRuleAttrType()

	if len(rules) == 0 {
		return types.SetNull(ruleType), nil
	}

	ruleValues := make([]attr.Value, 0, len(rules))

	for _, rule := range rules {
		if rule == nil {
			continue
		}

		var signal types.String
		if s := rule.GetSignal(); s != nil {
			signal = types.StringValue(s.String())
		} else {
			signal = types.StringNull()
		}

		threshold := convert.GraphToFrameworkInt32(rule.GetThreshold())

		var action types.String
		if a := rule.GetAction(); a != nil {
			action = types.StringValue(a.String())
		} else {
			action = types.StringNull()
		}

		ruleAttrs := map[string]attr.Value{
			"signal":    signal,
			"threshold": threshold,
			"action":    action,
		}

		ruleValue, diags := types.ObjectValue(ruleType.(types.ObjectType).AttrTypes, ruleAttrs)
		if diags.HasError() {
			return types.SetNull(ruleType), sentinels.ErrCreateMonitoringRuleObject
		}
		ruleValues = append(ruleValues, ruleValue)
	}

	set, diags := types.SetValue(ruleType, ruleValues)
	if diags.HasError() {
		return types.SetNull(ruleType), sentinels.ErrCreateMonitoringRulesSet
	}

	return set, nil
}
