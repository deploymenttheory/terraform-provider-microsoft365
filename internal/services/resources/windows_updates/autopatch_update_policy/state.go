package graphBetaWindowsUpdatesAutopatchUpdatePolicy

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodelswindowsupdates "github.com/microsoftgraph/msgraph-beta-sdk-go/models/windowsupdates"
)

func MapRemoteStateToTerraform(ctx context.Context, data *WindowsUpdatesAutopatchUpdatePolicyResourceModel, remoteResource graphmodelswindowsupdates.UpdatePolicyable) {
	tflog.Debug(ctx, fmt.Sprintf("Mapping remote state to Terraform state for %s", ResourceName))

	if remoteResource == nil {
		tflog.Debug(ctx, "Remote resource is nil")
		return
	}

	data.ID = convert.GraphToFrameworkString(remoteResource.GetId())
	data.CreatedDateTime = convert.GraphToFrameworkTime(remoteResource.GetCreatedDateTime())

	// Map audience
	if audience := remoteResource.GetAudience(); audience != nil {
		data.AudienceId = convert.GraphToFrameworkString(audience.GetId())
	}

	// Map compliance change rules - only if they exist in the response
	rules := remoteResource.GetComplianceChangeRules()
	if len(rules) > 0 {
		ruleElements := make([]attr.Value, 0, len(rules))
		
		for _, rule := range rules {
			ruleModel := ComplianceChangeRuleModel{
				CreatedDateTime:       convert.GraphToFrameworkTime(rule.GetCreatedDateTime()),
				LastEvaluatedDateTime: convert.GraphToFrameworkTime(rule.GetLastEvaluatedDateTime()),
				LastModifiedDateTime:  convert.GraphToFrameworkTime(rule.GetLastModifiedDateTime()),
			}

			// Type assert to ContentApprovalRule to access specific properties
			if contentApprovalRule, ok := rule.(graphmodelswindowsupdates.ContentApprovalRuleable); ok {
				// Map duration from additionalData to preserve raw format (avoid SDK normalization)
				if additionalData := contentApprovalRule.GetAdditionalData(); additionalData != nil {
					if rawDuration, ok := additionalData["durationBeforeDeploymentStart"]; ok {
						if durationStr, ok := rawDuration.(string); ok {
							ruleModel.DurationBeforeDeploymentStart = types.StringValue(durationStr)
						}
					}
				}

				// Map content filter
				if filter := contentApprovalRule.GetContentFilter(); filter != nil {
					filterModel := &ContentFilterModel{}
					// Check most specific types first since DriverUpdateFilter extends WindowsUpdateFilter
					switch f := filter.(type) {
					case graphmodelswindowsupdates.DriverUpdateFilterable:
						filterModel.FilterType = types.StringValue("driverUpdateFilter")
					default:
						// Check OData type for other filter types
						if odataType := f.GetOdataType(); odataType != nil {
							switch *odataType {
							case "#microsoft.graph.windowsUpdates.windowsUpdateFilter":
								filterModel.FilterType = types.StringValue("windowsUpdateFilter")
							default:
								filterModel.FilterType = types.StringValue("unknown")
							}
						} else {
							filterModel.FilterType = types.StringValue("unknown")
						}
					}
					ruleModel.ContentFilter = filterModel
				}
			}

			// Convert to object value
			ruleValue, diags := types.ObjectValueFrom(ctx, map[string]attr.Type{
				"content_filter": types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"filter_type": types.StringType,
					},
				},
				"duration_before_deployment_start": types.StringType,
				"created_date_time":                types.StringType,
				"last_evaluated_date_time":         types.StringType,
				"last_modified_date_time":          types.StringType,
			}, ruleModel)

			if !diags.HasError() {
				ruleElements = append(ruleElements, ruleValue)
			}
		}

		setValue, diags := types.SetValue(types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"content_filter": types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"filter_type": types.StringType,
					},
				},
				"duration_before_deployment_start": types.StringType,
				"created_date_time":                types.StringType,
				"last_evaluated_date_time":         types.StringType,
				"last_modified_date_time":          types.StringType,
			},
		}, ruleElements)

		if !diags.HasError() {
			data.ComplianceChangeRules = setValue
		}
	} else {
		// Set empty set if no rules
		data.ComplianceChangeRules = types.SetValueMust(types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"content_filter": types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"filter_type": types.StringType,
					},
				},
				"duration_before_deployment_start": types.StringType,
				"created_date_time":                types.StringType,
				"last_evaluated_date_time":         types.StringType,
				"last_modified_date_time":          types.StringType,
			},
		}, []attr.Value{})
	}

	// Map deployment settings - only if they exist in the response
	settings := remoteResource.GetDeploymentSettings()
	if settings != nil {
		deploymentSettingsData := DeploymentSettingsModel{}

		if schedule := settings.GetSchedule(); schedule != nil {
			scheduleSettingsData := ScheduleSettingsModel{
				StartDateTime: convert.GraphToFrameworkTime(schedule.GetStartDateTime()),
			}

			if rollout := schedule.GetGradualRollout(); rollout != nil {
				if rateDriven, ok := rollout.(graphmodelswindowsupdates.RateDrivenRolloutSettingsable); ok {
					gradualRolloutData := GradualRolloutModel{}

					// Map both duration and devices from additionalData to preserve raw format
					if additionalData := rateDriven.GetAdditionalData(); additionalData != nil {
						// Map duration between offers from additionalData (avoid SDK normalization)
						if rawDuration, ok := additionalData["durationBetweenOffers"]; ok {
							if durationStr, ok := rawDuration.(string); ok {
								gradualRolloutData.DurationBetweenOffers = types.StringValue(durationStr)
							}
						}

						// Map devices per offer from additionalData (as per Microsoft docs)
						if devicePerOffer, ok := additionalData["devicePerOffer"]; ok {
							switch v := devicePerOffer.(type) {
							case int32:
								gradualRolloutData.DevicesPerOffer = types.Int32Value(v)
							case float64:
								gradualRolloutData.DevicesPerOffer = types.Int32Value(int32(v))
							case int:
								gradualRolloutData.DevicesPerOffer = types.Int32Value(int32(v))
							case int64:
								gradualRolloutData.DevicesPerOffer = types.Int32Value(int32(v))
							}
						}
					}

					gradualRolloutObj, diags := types.ObjectValueFrom(ctx, GradualRolloutAttrTypes, gradualRolloutData)
					if diags.HasError() {
						tflog.Error(ctx, "Failed to convert GradualRollout to types.Object", map[string]any{
							"error": diags.Errors()[0].Detail(),
						})
						scheduleSettingsData.GradualRollout = types.ObjectNull(GradualRolloutAttrTypes)
					} else {
						scheduleSettingsData.GradualRollout = gradualRolloutObj
					}
				}
			} else {
				scheduleSettingsData.GradualRollout = types.ObjectNull(GradualRolloutAttrTypes)
			}

			scheduleObj, diags := types.ObjectValueFrom(ctx, ScheduleSettingsAttrTypes, scheduleSettingsData)
			if diags.HasError() {
				tflog.Error(ctx, "Failed to convert ScheduleSettings to types.Object", map[string]any{
					"error": diags.Errors()[0].Detail(),
				})
				deploymentSettingsData.Schedule = types.ObjectNull(ScheduleSettingsAttrTypes)
			} else {
				deploymentSettingsData.Schedule = scheduleObj
			}
		} else {
			deploymentSettingsData.Schedule = types.ObjectNull(ScheduleSettingsAttrTypes)
		}

		deploymentSettingsObj, diags := types.ObjectValueFrom(ctx, DeploymentSettingsAttrTypes, deploymentSettingsData)
		if diags.HasError() {
			tflog.Error(ctx, "Failed to convert DeploymentSettings to types.Object", map[string]any{
				"error": diags.Errors()[0].Detail(),
			})
			data.DeploymentSettings = types.ObjectNull(DeploymentSettingsAttrTypes)
		} else {
			data.DeploymentSettings = deploymentSettingsObj
		}
	} else {
		data.DeploymentSettings = types.ObjectNull(DeploymentSettingsAttrTypes)
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished mapping remote state to Terraform state for %s", ResourceName))
}
