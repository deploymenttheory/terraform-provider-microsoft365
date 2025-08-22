package graphBetaWindowsDeviceCompliancePolicies

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

// validateRequest validates the request data, specifically checking if notification_message_cc_list group IDs exist
func validateRequest(ctx context.Context, client *msgraphbetasdk.GraphServiceClient, data *DeviceCompliancePolicyResourceModel) diag.Diagnostics {
	var diags diag.Diagnostics

	if data.ScheduledActionsForRule.IsNull() || data.ScheduledActionsForRule.IsUnknown() {
		return diags
	}

	var scheduledActionsModels []ScheduledActionForRuleModel
	convertDiags := data.ScheduledActionsForRule.ElementsAs(ctx, &scheduledActionsModels, false)
	if convertDiags.HasError() {
		diags.Append(convertDiags...)
		return diags
	}

	for _, scheduledAction := range scheduledActionsModels {
		if scheduledAction.ScheduledActionConfigurations.IsNull() || scheduledAction.ScheduledActionConfigurations.IsUnknown() {
			continue
		}

		var configModels []ScheduledActionConfigurationModel
		convertDiags := scheduledAction.ScheduledActionConfigurations.ElementsAs(ctx, &configModels, false)
		if convertDiags.HasError() {
			diags.Append(convertDiags...)
			return diags
		}

		for _, config := range configModels {
			if config.NotificationMessageCcList.IsNull() || config.NotificationMessageCcList.IsUnknown() {
				continue
			}

			var groupIds []string
			convertDiags := config.NotificationMessageCcList.ElementsAs(ctx, &groupIds, false)
			if convertDiags.HasError() {
				diags.Append(convertDiags...)
				return diags
			}

			// Validate each group ID
			for _, groupId := range groupIds {
				if groupId == "" {
					continue // Skip empty group IDs
				}

				tflog.Debug(ctx, "Validating Microsoft 365 group ID", map[string]interface{}{
					"groupId": groupId,
				})

				group, err := client.
					Groups().
					ByGroupId(groupId).
					Get(ctx, nil)

				if err != nil {
					diags.AddError(
						"Invalid Group ID in notification_message_cc_list",
						fmt.Sprintf("The group ID '%s' in notification_message_cc_list is invalid or is not accessible. "+
							"Please verify the group ID is correct and your provider authentication method has the necessary permissions to read Microsoft 365 groups. "+
							"Error: %s", groupId, err.Error()),
					)
				} else {
					// Check if the group is mail-enabled
					if group.GetMailEnabled() == nil || !*group.GetMailEnabled() {
						diags.AddError(
							"Invalid Group Type in notification_message_cc_list",
							fmt.Sprintf("The group ID '%s' in notification_message_cc_list is not mail-enabled. "+
								"Only mail-enabled groups can receive notification messages. "+
								"Please use a Microsoft 365 group or distribution group that has mail enabled.", groupId),
						)
					} else {
						tflog.Debug(ctx, "Successfully validated mail-enabled Microsoft 365 group", map[string]interface{}{
							"groupId":     groupId,
							"mailEnabled": *group.GetMailEnabled(),
							"displayName": group.GetDisplayName(),
						})
					}
				}
			}
		}
	}

	return diags
}
