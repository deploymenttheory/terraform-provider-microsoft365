package graphBetaWindowsUpdatesAutopatchUpdatePolicy

import (
	"context"
	"fmt"
	"regexp"
	"strconv"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodelswindowsupdates "github.com/microsoftgraph/msgraph-beta-sdk-go/models/windowsupdates"
)

// denormalizeISODuration converts SDK-normalized durations back to day-based format
// to match user configuration (P1W -> P7D, P2W -> P14D, etc.)
func denormalizeISODuration(duration string) string {
	// Pattern for week-based durations (e.g., P1W, P2W)
	weekPattern := regexp.MustCompile(`^P(\d+)W$`)
	if matches := weekPattern.FindStringSubmatch(duration); len(matches) == 2 {
		weeks, err := strconv.Atoi(matches[1])
		if err != nil {
			return duration
		}

		// Convert weeks to days
		days := weeks * 7
		return "P" + strconv.Itoa(days) + "D"
	}

	// For other formats, return as-is
	return duration
}

func MapRemoteStateToTerraform(ctx context.Context, data *WindowsUpdatesAutopatchUpdatePolicyResourceModel, remoteResource graphmodelswindowsupdates.UpdatePolicyable) {
	tflog.Debug(ctx, fmt.Sprintf("Mapping remote state to Terraform state for %s", ResourceName))

	if remoteResource == nil {
		tflog.Debug(ctx, "Remote resource is nil")
		return
	}

	data.ID = convert.GraphToFrameworkString(remoteResource.GetId())
	data.CreatedDateTime = convert.GraphToFrameworkTime(remoteResource.GetCreatedDateTime())

	// Note: compliance_changes is write-only and not returned by the API
	// The value from config/state is preserved automatically

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

			tflog.Debug(ctx, "Mapping compliance rule", map[string]any{
				"created_date_time":        ruleModel.CreatedDateTime.ValueString(),
				"last_evaluated_date_time": ruleModel.LastEvaluatedDateTime.ValueString(),
				"last_modified_date_time":  ruleModel.LastModifiedDateTime.ValueString(),
			})

			// Type assert to ContentApprovalRule to access specific properties
			if contentApprovalRule, ok := rule.(graphmodelswindowsupdates.ContentApprovalRuleable); ok {
				// Map durationBeforeDeploymentStart using SDK getter
				// Note: SDK normalizes durations (e.g., P7D becomes P1W)
				// We denormalize back to days to match user config
				if duration := contentApprovalRule.GetDurationBeforeDeploymentStart(); duration != nil {
					durationStr := denormalizeISODuration(duration.String())
					ruleModel.DurationBeforeDeploymentStart = types.StringValue(durationStr)
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

					// SDK normalizes durations (P7D -> P1W), so we denormalize back to match hcl config which
					// the request body expects.
					if duration := rateDriven.GetDurationBetweenOffers(); duration != nil {
						durationStr := denormalizeISODuration(duration.String())
						gradualRolloutData.DurationBetweenOffers = types.StringValue(durationStr)
					}

					if devices := rateDriven.GetDevicesPerOffer(); devices != nil {
						gradualRolloutData.DevicesPerOffer = types.Int32Value(*devices)
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
