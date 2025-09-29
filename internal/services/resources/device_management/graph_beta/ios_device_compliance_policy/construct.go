package graphBetaIosDeviceCompliancePolicy

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

	iosPolicy := graphmodels.NewIosCompliancePolicy()
	if err := constructIosCompliancePolicy(ctx, data, iosPolicy); err != nil {
		return nil, fmt.Errorf("failed to construct ios policy: %s", err)
	}
	requestBody := iosPolicy

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

// constructIosCompliancePolicy handles ios specific settings using SDK setters
func constructIosCompliancePolicy(ctx context.Context, data *DeviceCompliancePolicyResourceModel, policy *graphmodels.IosCompliancePolicy) error {
	convert.FrameworkToGraphBool(data.PasscodeRequired, policy.SetPasscodeRequired)
	convert.FrameworkToGraphBool(data.PasscodeBlockSimple, policy.SetPasscodeBlockSimple)
	convert.FrameworkToGraphInt32(data.PasscodeMinutesOfInactivityBeforeLock, policy.SetPasscodeMinutesOfInactivityBeforeLock)
	convert.FrameworkToGraphInt32(data.PasscodeMinutesOfInactivityBeforeScreenTimeout, policy.SetPasscodeMinutesOfInactivityBeforeScreenTimeout)
	convert.FrameworkToGraphInt32(data.PasscodeExpirationDays, policy.SetPasscodeExpirationDays)
	convert.FrameworkToGraphInt32(data.PasscodeMinimumLength, policy.SetPasscodeMinimumLength)
	convert.FrameworkToGraphInt32(data.PasscodeMinimumCharacterSetCount, policy.SetPasscodeMinimumCharacterSetCount)
	convert.FrameworkToGraphInt32(data.PasscodePreviousPasscodeBlockCount, policy.SetPasscodePreviousPasscodeBlockCount)

	if err := convert.FrameworkToGraphEnum(data.PasscodeRequiredType,
		graphmodels.ParseRequiredPasswordType, policy.SetPasscodeRequiredType); err != nil {
		return fmt.Errorf("failed to set passcode required type: %s", err)
	}

	convert.FrameworkToGraphString(data.OsMinimumVersion, policy.SetOsMinimumVersion)
	convert.FrameworkToGraphString(data.OsMaximumVersion, policy.SetOsMaximumVersion)
	convert.FrameworkToGraphString(data.OsMinimumBuildVersion, policy.SetOsMinimumBuildVersion)
	convert.FrameworkToGraphString(data.OsMaximumBuildVersion, policy.SetOsMaximumBuildVersion)
	convert.FrameworkToGraphBool(data.DeviceThreatProtectionEnabled, policy.SetDeviceThreatProtectionEnabled)
	convert.FrameworkToGraphBool(data.ManagedEmailProfileRequired, policy.SetManagedEmailProfileRequired)
	convert.FrameworkToGraphBool(data.SecurityBlockJailbrokenDevices, policy.SetSecurityBlockJailbrokenDevices)

	if err := convert.FrameworkToGraphEnum(data.DeviceThreatProtectionRequiredSecurityLevel,
		graphmodels.ParseDeviceThreatProtectionLevel, policy.SetDeviceThreatProtectionRequiredSecurityLevel); err != nil {
		return fmt.Errorf("failed to set device threat protection required security level: %s", err)
	}

	if err := convert.FrameworkToGraphEnum(data.AdvancedThreatProtectionRequiredSecurityLevel,
		graphmodels.ParseDeviceThreatProtectionLevel, policy.SetAdvancedThreatProtectionRequiredSecurityLevel); err != nil {
		return fmt.Errorf("failed to set advanced threat protection required security level: %s", err)
	}

	if !data.RestrictedApps.IsNull() && !data.RestrictedApps.IsUnknown() {
		var restrictedAppsModels []RestrictedAppModel
		diags := data.RestrictedApps.ElementsAs(ctx, &restrictedAppsModels, false)
		if diags.HasError() {
			return fmt.Errorf("failed to parse restricted apps list: %v", diags.Errors())
		}

		if len(restrictedAppsModels) > 0 {
			restrictedApps, err := constructRestrictedApps(ctx, restrictedAppsModels)
			if err != nil {
				return fmt.Errorf("failed to construct restricted apps: %s", err)
			}
			policy.SetRestrictedApps(restrictedApps)
		}
	}

	return nil
}

// constructRestrictedApps constructs the restricted apps for the policy
func constructRestrictedApps(ctx context.Context, restrictedAppsModels []RestrictedAppModel) ([]graphmodels.AppListItemable, error) {
	restrictedApps := make([]graphmodels.AppListItemable, 0, len(restrictedAppsModels))
	for _, restrictedApp := range restrictedAppsModels {
		restrictedAppModel := graphmodels.NewAppListItem()

		convert.FrameworkToGraphString(restrictedApp.Name, restrictedAppModel.SetName)
		convert.FrameworkToGraphString(restrictedApp.AppId, restrictedAppModel.SetAppId)
		convert.FrameworkToGraphString(restrictedApp.AppStoreUrl, restrictedAppModel.SetAppStoreUrl)
		convert.FrameworkToGraphString(restrictedApp.Publisher, restrictedAppModel.SetPublisher)
		restrictedApps = append(restrictedApps, restrictedAppModel)
	}
	return restrictedApps, nil
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
