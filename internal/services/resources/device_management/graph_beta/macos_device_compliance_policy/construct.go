package graphBetaMacosDeviceCompliancePolicy

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

	macosPolicy := graphmodels.NewMacOSCompliancePolicy()
	if err := constructMacOSCompliancePolicy(ctx, data, macosPolicy); err != nil {
		return nil, fmt.Errorf("failed to construct macOS policy: %s", err)
	}
	requestBody := macosPolicy

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
		tflog.Error(ctx, "Failed to debug log object", map[string]interface{}{
			"error": err.Error(),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing %s resource", ResourceName))

	return requestBody, nil
}

// constructMacOSCompliancePolicy handles macOS specific settings using SDK setters
func constructMacOSCompliancePolicy(ctx context.Context, data *DeviceCompliancePolicyResourceModel, policy *graphmodels.MacOSCompliancePolicy) error {
	convert.FrameworkToGraphBool(data.PasswordRequired, policy.SetPasswordRequired)
	convert.FrameworkToGraphBool(data.PasswordBlockSimple, policy.SetPasswordBlockSimple)
	convert.FrameworkToGraphInt32(data.PasswordMinutesOfInactivityBeforeLock, policy.SetPasswordMinutesOfInactivityBeforeLock)
	convert.FrameworkToGraphInt32(data.PasswordExpirationDays, policy.SetPasswordExpirationDays)
	convert.FrameworkToGraphInt32(data.PasswordMinimumLength, policy.SetPasswordMinimumLength)
	convert.FrameworkToGraphInt32(data.PasswordMinimumCharacterSetCount, policy.SetPasswordMinimumCharacterSetCount)
	convert.FrameworkToGraphInt32(data.PasswordPreviousPasswordBlockCount, policy.SetPasswordPreviousPasswordBlockCount)

	if err := convert.FrameworkToGraphEnum(data.PasswordRequiredType,
		graphmodels.ParseRequiredPasswordType, policy.SetPasswordRequiredType); err != nil {
		return fmt.Errorf("failed to set password required type: %s", err)
	}

	convert.FrameworkToGraphString(data.OsMinimumVersion, policy.SetOsMinimumVersion)
	convert.FrameworkToGraphString(data.OsMaximumVersion, policy.SetOsMaximumVersion)
	convert.FrameworkToGraphString(data.OsMinimumBuildVersion, policy.SetOsMinimumBuildVersion)
	convert.FrameworkToGraphString(data.OsMaximumBuildVersion, policy.SetOsMaximumBuildVersion)
	convert.FrameworkToGraphBool(data.SystemIntegrityProtectionEnabled, policy.SetSystemIntegrityProtectionEnabled)
	convert.FrameworkToGraphBool(data.DeviceThreatProtectionEnabled, policy.SetDeviceThreatProtectionEnabled)
	convert.FrameworkToGraphBool(data.StorageRequireEncryption, policy.SetStorageRequireEncryption)

	if err := convert.FrameworkToGraphEnum(data.DeviceThreatProtectionRequiredSecurityLevel,
		graphmodels.ParseDeviceThreatProtectionLevel, policy.SetDeviceThreatProtectionRequiredSecurityLevel); err != nil {
		return fmt.Errorf("failed to set device threat protection required security level: %s", err)
	}

	if err := convert.FrameworkToGraphEnum(data.AdvancedThreatProtectionRequiredSecurityLevel,
		graphmodels.ParseDeviceThreatProtectionLevel, policy.SetAdvancedThreatProtectionRequiredSecurityLevel); err != nil {
		return fmt.Errorf("failed to set advanced threat protection required security level: %s", err)
	}

	if err := convert.FrameworkToGraphEnum(data.GatekeeperAllowedAppSource,
		graphmodels.ParseMacOSGatekeeperAppSources, policy.SetGatekeeperAllowedAppSource); err != nil {
		return fmt.Errorf("failed to set gatekeeper allowed app source: %s", err)
	}

	convert.FrameworkToGraphBool(data.FirewallEnabled, policy.SetFirewallEnabled)
	convert.FrameworkToGraphBool(data.FirewallBlockAllIncoming, policy.SetFirewallBlockAllIncoming)
	convert.FrameworkToGraphBool(data.FirewallEnableStealthMode, policy.SetFirewallEnableStealthMode)

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
