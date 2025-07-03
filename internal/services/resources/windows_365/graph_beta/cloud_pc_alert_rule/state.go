package graphBetaCloudPcAlertRule

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	models "github.com/microsoftgraph/msgraph-beta-sdk-go/models/devicemanagement"
)

// MapRemoteStateToTerraform maps the remote AlertRule to the Terraform resource model
func MapRemoteStateToTerraform(ctx context.Context, data *CloudPcAlertRuleResourceModel, remoteResource models.AlertRuleable) {
	if remoteResource == nil {
		tflog.Debug(ctx, "Remote resource is nil")
		return
	}

	tflog.Debug(ctx, "Starting to map remote state to Terraform state", map[string]interface{}{
		"resourceId": convert.GraphToFrameworkString(remoteResource.GetId()).ValueString(),
	})

	data.ID = convert.GraphToFrameworkString(remoteResource.GetId())
	data.AlertRuleTemplate = convert.GraphToFrameworkEnum(remoteResource.GetAlertRuleTemplate())
	data.Description = convert.GraphToFrameworkString(remoteResource.GetDescription())
	data.DisplayName = convert.GraphToFrameworkString(remoteResource.GetDisplayName())
	data.Enabled = convert.GraphToFrameworkBool(remoteResource.GetEnabled())
	data.IsSystemRule = convert.GraphToFrameworkBool(remoteResource.GetIsSystemRule())
	data.Severity = convert.GraphToFrameworkEnum(remoteResource.GetSeverity())

	// NotificationChannels (list of objects)
	if remoteResource.GetNotificationChannels() != nil {
		channels := remoteResource.GetNotificationChannels()
		result := make([]NotificationChannelModel, 0, len(channels))
		for _, ch := range channels {
			if ch == nil {
				continue
			}
			var channel NotificationChannelModel

			channel.NotificationChannelType = convert.GraphToFrameworkEnum(ch.GetNotificationChannelType())
			if ch.GetNotificationReceivers() != nil {
				receivers := ch.GetNotificationReceivers()
				receiverModels := make([]NotificationReceiverModel, 0, len(receivers))
				for _, r := range receivers {
					if r == nil {
						continue
					}
					var receiver NotificationReceiverModel

					receiver.ContactInformation = convert.GraphToFrameworkString(r.GetContactInformation())
					receiver.Locale = convert.GraphToFrameworkString(r.GetLocale())
					receiverModels = append(receiverModels, receiver)
				}
				channel.NotificationReceivers = receiverModels
			}
			result = append(result, channel)
		}
		data.NotificationChannels = result
	}

	// Threshold (single object, deprecated)
	if remoteResource.GetThreshold() != nil {
		th := remoteResource.GetThreshold()
		var threshold RuleThresholdModel

		threshold.Aggregation = convert.GraphToFrameworkEnum(th.GetAggregation())
		threshold.Operator = convert.GraphToFrameworkEnum(th.GetOperator())
		threshold.Target = convert.GraphToFrameworkInt32(th.GetTarget())
		data.Threshold = &threshold
	}

	// Conditions (list of objects)
	if remoteResource.GetConditions() != nil {
		conditions := remoteResource.GetConditions()
		result := make([]RuleConditionModel, 0, len(conditions))
		for _, cond := range conditions {
			if cond == nil {
				continue
			}
			var c RuleConditionModel

			c.RelationshipType = convert.GraphToFrameworkEnum(cond.GetRelationshipType())
			c.ConditionCategory = convert.GraphToFrameworkEnum(cond.GetConditionCategory())
			c.Aggregation = convert.GraphToFrameworkEnum(cond.GetAggregation())
			c.Operator = convert.GraphToFrameworkEnum(cond.GetOperator())
			c.ThresholdValue = convert.GraphToFrameworkString(cond.GetThresholdValue())
			result = append(result, c)
		}
		data.Conditions = result
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished mapping remote state for resource %s with id %s", ResourceName, data.ID.ValueString()))
}
