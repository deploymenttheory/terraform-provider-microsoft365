package graphBetaLinuxDeviceCompliancePolicy

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/devicemanagement"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// constructScheduledActions creates the request body for the scheduleActionsForRules API call
func constructScheduledActions(ctx context.Context, scheduledActionsData ScheduledActionForRuleModel) (devicemanagement.CompliancePoliciesItemSetScheduledActionsPostRequestBodyable, error) {

	requestBody := devicemanagement.NewCompliancePoliciesItemSetScheduledActionsPostRequestBody()

	managementAction := graphmodels.NewDeviceManagementComplianceScheduledActionForRule()

	// Set rule name
	if !scheduledActionsData.RuleName.IsNull() && !scheduledActionsData.RuleName.IsUnknown() {
		ruleName := scheduledActionsData.RuleName.ValueString()
		managementAction.SetRuleName(&ruleName)
	} else {
		// Default rule name as per SDK documentation
		defaultRuleName := "PasswordRequired"
		managementAction.SetRuleName(&defaultRuleName)
	}

	if !scheduledActionsData.ScheduledActionConfigurations.IsNull() && !scheduledActionsData.ScheduledActionConfigurations.IsUnknown() {
		var configModels []ScheduledActionConfigurationModel
		diags := scheduledActionsData.ScheduledActionConfigurations.ElementsAs(ctx, &configModels, false)
		if diags.HasError() {
			return nil, fmt.Errorf("failed to convert scheduled action configurations: %v", diags.Errors())
		}

		configurations := make([]graphmodels.DeviceManagementComplianceActionItemable, 0, len(configModels))
		for _, config := range configModels {
			actionItem := graphmodels.NewDeviceManagementComplianceActionItem()

			if !config.ActionType.IsNull() && !config.ActionType.IsUnknown() {
				if err := convert.FrameworkToGraphEnum(config.ActionType,
					graphmodels.ParseDeviceManagementComplianceActionType, actionItem.SetActionType); err != nil {
					return nil, fmt.Errorf("failed to set action type: %s", err)
				}
			} else {
				blockAction := graphmodels.BLOCK_DEVICEMANAGEMENTCOMPLIANCEACTIONTYPE
				actionItem.SetActionType(&blockAction)
			}

			convert.FrameworkToGraphInt32(config.GracePeriodHours, actionItem.SetGracePeriodHours)
			convert.FrameworkToGraphString(config.NotificationTemplateId, actionItem.SetNotificationTemplateId)

			if err := convert.FrameworkToGraphStringList(ctx, config.NotificationMessageCcList, actionItem.SetNotificationMessageCCList); err != nil {
				return nil, fmt.Errorf("failed to set notification message CC list: %s", err)
			}

			configurations = append(configurations, actionItem)
		}

		managementAction.SetScheduledActionConfigurations(configurations)
	}

	managementScheduledActions := []graphmodels.DeviceManagementComplianceScheduledActionForRuleable{managementAction}
	requestBody.SetScheduledActions(managementScheduledActions)

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]interface{}{
			"error": err.Error(),
		})
	}

	return requestBody, nil
}
