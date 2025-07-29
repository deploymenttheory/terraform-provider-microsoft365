package graphBetaAndroidDeviceOwnerCompliancePolicy

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/devicemanagement"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// constructDeviceComplianceScheduledActionForRulesWithPatchMethod creates the request body for the scheduleActionsForRules API call
func constructDeviceComplianceScheduledActionForRulesWithPatchMethod(ctx context.Context, scheduledActionsData ScheduledActionForRuleModel) (devicemanagement.DeviceCompliancePoliciesItemScheduleActionsForRulesPostRequestBodyable, error) {
	scheduledActions := make([]graphmodels.DeviceComplianceScheduledActionForRuleable, 0, 1)
	scheduledAction := graphmodels.NewDeviceComplianceScheduledActionForRule()

	if !scheduledActionsData.RuleName.IsNull() && !scheduledActionsData.RuleName.IsUnknown() {
		ruleName := scheduledActionsData.RuleName.ValueString()
		scheduledAction.SetRuleName(&ruleName)
	}

	if !scheduledActionsData.ScheduledActionConfigurations.IsNull() && !scheduledActionsData.ScheduledActionConfigurations.IsUnknown() {
		configs, err := constructScheduledActionItem(ctx, scheduledActionsData.ScheduledActionConfigurations)
		if err != nil {
			return nil, fmt.Errorf("failed to construct scheduled action configurations: %s", err)
		}
		scheduledAction.SetScheduledActionConfigurations(configs)
	}

	scheduledActions = append(scheduledActions, scheduledAction)

	requestBody := devicemanagement.NewDeviceCompliancePoliciesItemScheduleActionsForRulesPostRequestBody()
	requestBody.SetDeviceComplianceScheduledActionForRules(scheduledActions)

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]interface{}{
			"error": err.Error(),
		})
	}

	return requestBody, nil
}

// constructScheduledActionItem converts Terraform Set to Graph SDK model using proper SDK types
func constructScheduledActionItem(ctx context.Context, configurationsData types.Set) ([]graphmodels.DeviceComplianceActionItemable, error) {
	var configModels []ScheduledActionConfigurationModel
	diags := configurationsData.ElementsAs(ctx, &configModels, false)
	if diags.HasError() {
		return nil, fmt.Errorf("failed to convert scheduled action configurations: %v", diags.Errors())
	}

	configurations := make([]graphmodels.DeviceComplianceActionItemable, 0, len(configModels))
	for _, config := range configModels {
		actionItem := graphmodels.NewDeviceComplianceActionItem()

		if err := convert.FrameworkToGraphEnum(config.ActionType,
			graphmodels.ParseDeviceComplianceActionType, actionItem.SetActionType); err != nil {
			return nil, fmt.Errorf("failed to set action type: %s", err)
		}

		// Convert GracePeriodHours to Int32 for the SDK
		convert.FrameworkToGraphInt32(config.GracePeriodHours, actionItem.SetGracePeriodHours)

		convert.FrameworkToGraphString(config.NotificationTemplateId, actionItem.SetNotificationTemplateId)

		if err := convert.FrameworkToGraphStringList(ctx, config.NotificationMessageCcList, actionItem.SetNotificationMessageCCList); err != nil {
			return nil, fmt.Errorf("failed to set notification message CC list: %s", err)
		}

		configurations = append(configurations, actionItem)
	}

	return configurations, nil
}
