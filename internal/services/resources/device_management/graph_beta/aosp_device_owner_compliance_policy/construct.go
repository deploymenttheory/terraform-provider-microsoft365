package graphBetaAospDeviceOwnerCompliancePolicy

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// constructResource maps the Terraform schema to the SDK model
func constructResource(ctx context.Context, data *DeviceCompliancePolicyResourceModel, isCreate bool) (graphmodels.DeviceCompliancePolicyable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s resource from model", ResourceName))

	aospPolicy := graphmodels.NewAospDeviceOwnerCompliancePolicy()
	if err := constructAospCompliancePolicy(ctx, data, aospPolicy); err != nil {
		return nil, fmt.Errorf("failed to construct aosp policy: %s", err)
	}
	requestBody := aospPolicy

	convert.FrameworkToGraphString(data.DisplayName, requestBody.SetDisplayName)
	convert.FrameworkToGraphString(data.Description, requestBody.SetDescription)

	if err := convert.FrameworkToGraphStringSet(ctx, data.RoleScopeTagIds, requestBody.SetRoleScopeTagIds); err != nil {
		return nil, fmt.Errorf("failed to set role scope tags: %s", err)
	}

	// Include scheduled actions during create operation as API requires at least one block action
	if !data.ScheduledActionsForRule.IsNull() && !data.ScheduledActionsForRule.IsUnknown() && isCreate {
		var scheduledActionsModels []ScheduledActionForRuleModel
		diags := data.ScheduledActionsForRule.ElementsAs(ctx, &scheduledActionsModels, false)
		if diags.HasError() {
			return nil, fmt.Errorf("failed to parse scheduled actions list: %v", diags.Errors())
		}

		if len(scheduledActionsModels) > 0 {
			firstRule := scheduledActionsModels[0]
			scheduledActions, err := constructScheduledActionsForPolicyCreation(ctx, firstRule)
			if err != nil {
				return nil, fmt.Errorf("failed to construct scheduled actions for rule: %s", err)
			}
			requestBody.SetScheduledActionsForRule(scheduledActions)
		}
	}

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]any{
			"error": err.Error(),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing %s resource", ResourceName))

	return requestBody, nil
}

// constructAospCompliancePolicy handles aosp specific settings using SDK setters
func constructAospCompliancePolicy(ctx context.Context, data *DeviceCompliancePolicyResourceModel, policy *graphmodels.AospDeviceOwnerCompliancePolicy) error {
	convert.FrameworkToGraphBool(data.PasscodeRequired, policy.SetPasswordRequired)
	convert.FrameworkToGraphInt32(data.PasscodeMinimumLength, policy.SetPasswordMinimumLength)
	convert.FrameworkToGraphInt32(data.PasscodeMinutesOfInactivityBeforeLock, policy.SetPasswordMinutesOfInactivityBeforeLock)

	if err := convert.FrameworkToGraphEnum(data.PasscodeRequiredType,
		graphmodels.ParseAndroidDeviceOwnerRequiredPasswordType, policy.SetPasswordRequiredType); err != nil {
		return fmt.Errorf("failed to set password required type: %s", err)
	}

	convert.FrameworkToGraphString(data.OsMinimumVersion, policy.SetOsMinimumVersion)
	convert.FrameworkToGraphString(data.OsMaximumVersion, policy.SetOsMaximumVersion)
	convert.FrameworkToGraphBool(data.SecurityBlockJailbrokenDevices, policy.SetSecurityBlockJailbrokenDevices)
	convert.FrameworkToGraphBool(data.StorageRequireEncryption, policy.SetStorageRequireEncryption)
	convert.FrameworkToGraphString(data.MinAndroidSecurityPatchLevel, policy.SetMinAndroidSecurityPatchLevel)

	return nil
}

// constructScheduledActionsForPolicyCreation creates scheduled actions for inclusion during policy creation
func constructScheduledActionsForPolicyCreation(ctx context.Context, scheduledActionData ScheduledActionForRuleModel) ([]graphmodels.DeviceComplianceScheduledActionForRuleable, error) {
	scheduledActions := make([]graphmodels.DeviceComplianceScheduledActionForRuleable, 0, 1)
	scheduledAction := graphmodels.NewDeviceComplianceScheduledActionForRule()

	if !scheduledActionData.RuleName.IsNull() && !scheduledActionData.RuleName.IsUnknown() {
		ruleName := scheduledActionData.RuleName.ValueString()
		scheduledAction.SetRuleName(&ruleName)
	}

	if !scheduledActionData.ScheduledActionConfigurations.IsNull() && !scheduledActionData.ScheduledActionConfigurations.IsUnknown() {
		configs, err := constructScheduledActionItem(ctx, scheduledActionData.ScheduledActionConfigurations)
		if err != nil {
			return nil, fmt.Errorf("failed to construct scheduled action configurations: %s", err)
		}
		scheduledAction.SetScheduledActionConfigurations(configs)
	}

	scheduledActions = append(scheduledActions, scheduledAction)
	return scheduledActions, nil
}
