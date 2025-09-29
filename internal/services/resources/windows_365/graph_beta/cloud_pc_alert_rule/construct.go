package graphBetaCloudPcAlertRule

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	models "github.com/microsoftgraph/msgraph-beta-sdk-go/models/devicemanagement"
)

func constructResource(ctx context.Context, data *CloudPcAlertRuleResourceModel) (models.AlertRuleable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s resource from model", ResourceName))
	requestBody := models.NewAlertRule()

	convert.FrameworkToGraphEnum(data.AlertRuleTemplate, models.ParseAlertRuleTemplate, requestBody.SetAlertRuleTemplate)
	convert.FrameworkToGraphString(data.Description, requestBody.SetDescription)
	convert.FrameworkToGraphString(data.DisplayName, requestBody.SetDisplayName)
	convert.FrameworkToGraphBool(data.Enabled, requestBody.SetEnabled)
	convert.FrameworkToGraphEnum(data.Severity, models.ParseRuleSeverityType, requestBody.SetSeverity)

	// NotificationChannels (list of objects)
	if data.NotificationChannels != nil {
		channels := make([]models.NotificationChannelable, 0, len(data.NotificationChannels))
		for _, ch := range data.NotificationChannels {
			channel := models.NewNotificationChannel()
			convert.FrameworkToGraphEnum(ch.NotificationChannelType, models.ParseNotificationChannelType, channel.SetNotificationChannelType)

			// NotificationReceivers (list of objects)
			if ch.NotificationReceivers != nil {
				receivers := make([]models.NotificationReceiverable, 0, len(ch.NotificationReceivers))
				for _, r := range ch.NotificationReceivers {
					receiver := models.NewNotificationReceiver()
					convert.FrameworkToGraphString(r.ContactInformation, receiver.SetContactInformation)
					convert.FrameworkToGraphString(r.Locale, receiver.SetLocale)
					receivers = append(receivers, receiver)
				}
				channel.SetNotificationReceivers(receivers)
			}

			channels = append(channels, channel)
		}
		requestBody.SetNotificationChannels(channels)
	}

	// Threshold (single object, deprecated)
	if data.Threshold != nil {
		threshold := models.NewRuleThreshold()
		convert.FrameworkToGraphEnum(data.Threshold.Aggregation, models.ParseAggregationType, threshold.SetAggregation)
		convert.FrameworkToGraphEnum(data.Threshold.Operator, models.ParseOperatorType, threshold.SetOperator)
		convert.FrameworkToGraphInt32(data.Threshold.Target, threshold.SetTarget)
		requestBody.SetThreshold(threshold)
	}

	// Conditions (list of objects)
	if data.Conditions != nil {
		conditions := make([]models.RuleConditionable, 0, len(data.Conditions))
		for _, cond := range data.Conditions {
			c := models.NewRuleCondition()
			convert.FrameworkToGraphEnum(cond.RelationshipType, models.ParseRelationshipType, c.SetRelationshipType)
			convert.FrameworkToGraphEnum(cond.ConditionCategory, models.ParseConditionCategory, c.SetConditionCategory)
			convert.FrameworkToGraphEnum(cond.Aggregation, models.ParseAggregationType, c.SetAggregation)
			convert.FrameworkToGraphEnum(cond.Operator, models.ParseOperatorType, c.SetOperator)
			convert.FrameworkToGraphString(cond.ThresholdValue, c.SetThresholdValue)
			conditions = append(conditions, c)
		}
		requestBody.SetConditions(conditions)
	}

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]any{
			"error": err.Error(),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing %s resource", ResourceName))

	return requestBody, nil
}
