package graphBetaWindowsDeviceCompliancePolicy

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// constructResource maps the Terraform schema to the SDK model
func constructResource(ctx context.Context, data *DeviceCompliancePolicyResourceModel, isCreate bool) (graphmodels.DeviceCompliancePolicyable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s resource from model", ResourceName))

	windowsPolicy := graphmodels.NewWindows10CompliancePolicy()
	if err := constructWindows10CompliancePolicy(ctx, data, windowsPolicy); err != nil {
		return nil, fmt.Errorf("failed to construct Windows 10 policy: %s", err)
	}
	requestBody := windowsPolicy

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

// constructWindows10CompliancePolicy handles Windows 10 specific settings using SDK setters
func constructWindows10CompliancePolicy(ctx context.Context, data *DeviceCompliancePolicyResourceModel, policy *graphmodels.Windows10CompliancePolicy) error {
	// all of these fields are now deprecated.

	//convert.FrameworkToGraphBool(data.RequireHealthyDeviceReport, policy.SetRequireHealthyDeviceReport)
	//convert.FrameworkToGraphBool(data.EarlyLaunchAntiMalwareDriverEnabled, policy.SetEarlyLaunchAntiMalwareDriverEnabled)
	// convert.FrameworkToGraphBool(data.MemoryIntegrityEnabled, policy.SetMemoryIntegrityEnabled)
	// convert.FrameworkToGraphBool(data.KernelDmaProtectionEnabled, policy.SetKernelDmaProtectionEnabled)
	// convert.FrameworkToGraphBool(data.VirtualizationBasedSecurityEnabled, policy.SetVirtualizationBasedSecurityEnabled)
	// convert.FrameworkToGraphBool(data.FirmwareProtectionEnabled, policy.SetFirmwareProtectionEnabled)

	// Handle device_health object
	if !data.DeviceHealth.IsNull() && !data.DeviceHealth.IsUnknown() {
		var deviceHealth DeviceHealthModel
		diags := data.DeviceHealth.As(ctx, &deviceHealth, basetypes.ObjectAsOptions{})
		if !diags.HasError() {
			convert.FrameworkToGraphBool(deviceHealth.BitLockerEnabled, policy.SetBitLockerEnabled)
			convert.FrameworkToGraphBool(deviceHealth.SecureBootEnabled, policy.SetSecureBootEnabled)
			convert.FrameworkToGraphBool(deviceHealth.CodeIntegrityEnabled, policy.SetCodeIntegrityEnabled)
		}
	}

	// Handle system_security object
	if !data.SystemSecurity.IsNull() && !data.SystemSecurity.IsUnknown() {
		var systemSecurity SystemSecurityModel
		diags := data.SystemSecurity.As(ctx, &systemSecurity, basetypes.ObjectAsOptions{})
		if !diags.HasError() {
			convert.FrameworkToGraphBool(systemSecurity.StorageRequireEncryption, policy.SetStorageRequireEncryption)
			convert.FrameworkToGraphBool(systemSecurity.ActiveFirewallRequired, policy.SetActiveFirewallRequired)
			convert.FrameworkToGraphBool(systemSecurity.DefenderEnabled, policy.SetDefenderEnabled)
			convert.FrameworkToGraphString(systemSecurity.DefenderVersion, policy.SetDefenderVersion)
			convert.FrameworkToGraphBool(systemSecurity.SignatureOutOfDate, policy.SetSignatureOutOfDate)
			convert.FrameworkToGraphBool(systemSecurity.RtpEnabled, policy.SetRtpEnabled)
			convert.FrameworkToGraphBool(systemSecurity.AntivirusRequired, policy.SetAntivirusRequired)
			convert.FrameworkToGraphBool(systemSecurity.AntiSpywareRequired, policy.SetAntiSpywareRequired)
			convert.FrameworkToGraphBool(systemSecurity.ConfigurationManagerComplianceRequired, policy.SetConfigurationManagerComplianceRequired)
			convert.FrameworkToGraphBool(systemSecurity.TpmRequired, policy.SetTpmRequired)
			convert.FrameworkToGraphBool(systemSecurity.PasswordRequired, policy.SetPasswordRequired)
			convert.FrameworkToGraphBool(systemSecurity.PasswordBlockSimple, policy.SetPasswordBlockSimple)
			convert.FrameworkToGraphBool(systemSecurity.PasswordRequiredToUnlockFromIdle, policy.SetPasswordRequiredToUnlockFromIdle)
			convert.FrameworkToGraphInt32(systemSecurity.PasswordMinimumCharacterSetCount, policy.SetPasswordMinimumCharacterSetCount)

			if err := convert.FrameworkToGraphEnum(systemSecurity.PasswordRequiredType,
				graphmodels.ParseRequiredPasswordType, policy.SetPasswordRequiredType); err != nil {
				return fmt.Errorf("failed to set password required type: %s", err)
			}

		}
	}

	// Handle microsoft_defender_for_endpoint object
	if !data.MicrosoftDefenderForEndpoint.IsNull() && !data.MicrosoftDefenderForEndpoint.IsUnknown() {
		var microsoftDefender MicrosoftDefenderForEndpointModel
		diags := data.MicrosoftDefenderForEndpoint.As(ctx, &microsoftDefender, basetypes.ObjectAsOptions{})
		if !diags.HasError() {
			convert.FrameworkToGraphBool(microsoftDefender.DeviceThreatProtectionEnabled, policy.SetDeviceThreatProtectionEnabled)

			if err := convert.FrameworkToGraphEnum(microsoftDefender.DeviceThreatProtectionRequiredSecurityLevel,
				graphmodels.ParseDeviceThreatProtectionLevel, policy.SetDeviceThreatProtectionRequiredSecurityLevel); err != nil {
				return fmt.Errorf("failed to set device threat protection required security level: %s", err)
			}
		}
	}

	// Handle device_properties object
	if !data.DeviceProperties.IsNull() && !data.DeviceProperties.IsUnknown() {
		var deviceProperties DevicePropertiesModel
		diags := data.DeviceProperties.As(ctx, &deviceProperties, basetypes.ObjectAsOptions{})
		if !diags.HasError() {
			convert.FrameworkToGraphString(deviceProperties.OsMinimumVersion, policy.SetOsMinimumVersion)
			convert.FrameworkToGraphString(deviceProperties.OsMaximumVersion, policy.SetOsMaximumVersion)
			convert.FrameworkToGraphString(deviceProperties.MobileOsMinimumVersion, policy.SetMobileOsMinimumVersion)
			convert.FrameworkToGraphString(deviceProperties.MobileOsMaximumVersion, policy.SetMobileOsMaximumVersion)

			if !deviceProperties.ValidOperatingSystemBuildRanges.IsNull() && !deviceProperties.ValidOperatingSystemBuildRanges.IsUnknown() {
				buildRanges, err := constructValidOperatingSystemBuildRanges(ctx, deviceProperties.ValidOperatingSystemBuildRanges)
				if err != nil {
					return fmt.Errorf("failed to construct valid operating system build ranges from device_properties: %s", err)
				}
				policy.SetValidOperatingSystemBuildRanges(buildRanges)
			}
		}
	}

	if !data.WslDistributions.IsNull() && !data.WslDistributions.IsUnknown() {
		wslDistributions, err := constructWslDistributions(ctx, data.WslDistributions)
		if err != nil {
			return fmt.Errorf("failed to construct WSL distributions: %s", err)
		}
		policy.SetWslDistributions(wslDistributions)
	}

	// Handle device compliance policy script
	if !data.DeviceCompliancePolicyScript.IsNull() && !data.DeviceCompliancePolicyScript.IsUnknown() {
		script, err := constructDeviceCompliancePolicyScript(ctx, data.DeviceCompliancePolicyScript)
		if err != nil {
			return fmt.Errorf("failed to construct device compliance policy script: %s", err)
		}
		policy.SetDeviceCompliancePolicyScript(script)

		// Set custom compliance required flag if we have a script
		// Note: There's no direct SetCustomComplianceRequired method in the SDK
		// The API will infer this from the presence of a script
		if !data.CustomComplianceRequired.IsNull() && !data.CustomComplianceRequired.IsUnknown() && data.CustomComplianceRequired.ValueBool() {
			// The SDK doesn't have a direct setter for this field
			// It's inferred from the presence of a device compliance policy script
			tflog.Debug(ctx, "Custom compliance is required, script is set")
		}
	}

	return nil
}

// constructDeviceCompliancePolicyScript converts Terraform Object to Graph SDK model
func constructDeviceCompliancePolicyScript(ctx context.Context, scriptData types.Object) (graphmodels.DeviceCompliancePolicyScriptable, error) {
	attrs := scriptData.Attributes()

	script := graphmodels.NewDeviceCompliancePolicyScript()

	if scriptIdAttr, ok := attrs["device_compliance_script_id"].(types.String); ok && !scriptIdAttr.IsNull() {
		scriptId := scriptIdAttr.ValueString()
		script.SetDeviceComplianceScriptId(&scriptId)
	}

	if rulesContentAttr, ok := attrs["rules_content"].(types.String); ok && !rulesContentAttr.IsNull() {
		rulesContentStr := rulesContentAttr.ValueString()
		// The rules content comes as a JSON string from the user which needs to be encoded to base64 for the API request
		encodedBytes := []byte(rulesContentStr)
		script.SetRulesContent(encodedBytes)
	}

	return script, nil
}

// constructWslDistributions converts Terraform Set to Graph SDK model using proper SDK types
func constructWslDistributions(ctx context.Context, wslDistributionsData types.Set) ([]graphmodels.WslDistributionConfigurationable, error) {
	var wslDistributionModels []WslDistributionModel
	diags := wslDistributionsData.ElementsAs(ctx, &wslDistributionModels, false)
	if diags.HasError() {
		return nil, fmt.Errorf("failed to convert WSL distributions: %v", diags.Errors())
	}

	wslDistributions := make([]graphmodels.WslDistributionConfigurationable, 0, len(wslDistributionModels))
	for _, wslDist := range wslDistributionModels {
		distribution := graphmodels.NewWslDistributionConfiguration()

		convert.FrameworkToGraphString(wslDist.Distribution, distribution.SetDistribution)
		convert.FrameworkToGraphString(wslDist.MinimumOSVersion, distribution.SetMinimumOSVersion)
		convert.FrameworkToGraphString(wslDist.MaximumOSVersion, distribution.SetMaximumOSVersion)

		wslDistributions = append(wslDistributions, distribution)
	}

	return wslDistributions, nil
}

// constructValidOperatingSystemBuildRanges converts Terraform Set to Graph SDK model using proper SDK types
func constructValidOperatingSystemBuildRanges(ctx context.Context, buildRangesData types.Set) ([]graphmodels.OperatingSystemVersionRangeable, error) {
	var buildRangeModels []ValidOperatingSystemBuildRangeModel
	diags := buildRangesData.ElementsAs(ctx, &buildRangeModels, false)
	if diags.HasError() {
		return nil, fmt.Errorf("failed to convert valid operating system build ranges: %v", diags.Errors())
	}

	buildRanges := make([]graphmodels.OperatingSystemVersionRangeable, 0, len(buildRangeModels))
	for _, buildRange := range buildRangeModels {
		versionRange := graphmodels.NewOperatingSystemVersionRange()

		convert.FrameworkToGraphString(buildRange.LowOSVersion, versionRange.SetLowestVersion)
		convert.FrameworkToGraphString(buildRange.HighOSVersion, versionRange.SetHighestVersion)
		convert.FrameworkToGraphString(buildRange.Description, versionRange.SetDescription)

		buildRanges = append(buildRanges, versionRange)
	}

	return buildRanges, nil
}

// constructScheduledActionsForPolicyCreation creates scheduled actions for inclusion during policy creation
func constructScheduledActionsForPolicyCreation(ctx context.Context, scheduledActionData ScheduledActionForRuleModel) ([]graphmodels.DeviceComplianceScheduledActionForRuleable, error) {
	scheduledActions := make([]graphmodels.DeviceComplianceScheduledActionForRuleable, 0, 1)
	scheduledAction := graphmodels.NewDeviceComplianceScheduledActionForRule()

	// Always set rule name to "PasswordRequired" - API requirement but not user configurable
	// value is never returned by the API
	ruleName := "PasswordRequired"
	scheduledAction.SetRuleName(&ruleName)

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
