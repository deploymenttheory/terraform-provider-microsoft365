package graphBetaWin32LobApp

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	sharedmodels "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/shared_models/graph_beta/device_and_app_management"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// MapRemoteResourceStateToTerraform maps the remote state to the win32lobapp to Terraform state
func MapRemoteResourceStateToTerraform(ctx context.Context, data *Win32LobAppResourceModel, remoteResource graphmodels.Win32LobAppable) {
	if remoteResource == nil {
		tflog.Debug(ctx, "Remote resource is nil")
		return
	}

	tflog.Debug(ctx, "Starting to map remote state to Terraform state", map[string]any{
		"resourceId": convert.GraphToFrameworkString(remoteResource.GetId()).ValueString(),
	})

	data.ID = convert.GraphToFrameworkString(remoteResource.GetId())
	data.DisplayName = convert.GraphToFrameworkString(remoteResource.GetDisplayName())
	data.DisplayVersion = convert.GraphToFrameworkString(remoteResource.GetDisplayVersion())
	data.Description = convert.GraphToFrameworkString(remoteResource.GetDescription())
	data.Publisher = convert.GraphToFrameworkString(remoteResource.GetPublisher())
	data.CreatedDateTime = convert.GraphToFrameworkTime(remoteResource.GetCreatedDateTime())
	data.LastModifiedDateTime = convert.GraphToFrameworkTime(remoteResource.GetLastModifiedDateTime())
	data.IsFeatured = convert.GraphToFrameworkBool(remoteResource.GetIsFeatured())
	data.PrivacyInformationUrl = convert.GraphToFrameworkString(remoteResource.GetPrivacyInformationUrl())
	data.InformationUrl = convert.GraphToFrameworkString(remoteResource.GetInformationUrl())
	data.Owner = convert.GraphToFrameworkString(remoteResource.GetOwner())
	data.Developer = convert.GraphToFrameworkString(remoteResource.GetDeveloper())
	data.Notes = convert.GraphToFrameworkString(remoteResource.GetNotes())
	data.UploadState = convert.GraphToFrameworkInt32(remoteResource.GetUploadState())
	data.PublishingState = convert.GraphToFrameworkEnum(remoteResource.GetPublishingState())
	data.IsAssigned = convert.GraphToFrameworkBool(remoteResource.GetIsAssigned())
	data.InformationUrl = convert.GraphToFrameworkString(remoteResource.GetInformationUrl())
	data.PrivacyInformationUrl = convert.GraphToFrameworkString(remoteResource.GetPrivacyInformationUrl())
	data.AllowAvailableUninstall = convert.GraphToFrameworkBool(remoteResource.GetAllowAvailableUninstall())

	if data.AppIcon != nil {
		tflog.Debug(ctx, "Preserving original app_icon values from configuration")
	} else if largeIcon := remoteResource.GetLargeIcon(); largeIcon != nil {
		data.AppIcon = &sharedmodels.MobileAppIconResourceModel{
			IconFilePathSource: types.StringNull(),
			IconURLSource:      types.StringNull(),
		}
	} else {
		data.AppIcon = nil
	}

	data.RoleScopeTagIds = convert.GraphToFrameworkStringSet(ctx, remoteResource.GetRoleScopeTagIds())
	data.DependentAppCount = convert.GraphToFrameworkInt32(remoteResource.GetDependentAppCount())
	data.SupersedingAppCount = convert.GraphToFrameworkInt32(remoteResource.GetSupersedingAppCount())
	data.SupersededAppCount = convert.GraphToFrameworkInt32(remoteResource.GetSupersededAppCount())
	data.CommittedContentVersion = convert.GraphToFrameworkString(remoteResource.GetCommittedContentVersion())
	data.FileName = convert.GraphToFrameworkString(remoteResource.GetFileName())
	data.Size = convert.GraphToFrameworkInt64(remoteResource.GetSize())
	data.InstallCommandLine = convert.GraphToFrameworkString(remoteResource.GetInstallCommandLine())
	data.UninstallCommandLine = convert.GraphToFrameworkString(remoteResource.GetUninstallCommandLine())
	data.MinimumSupportedWindowsRelease = convert.GraphToFrameworkString(remoteResource.GetMinimumSupportedWindowsRelease())

	// Handle applicable architectures
	if applicableArchitectures := remoteResource.GetAllowedArchitectures(); applicableArchitectures != nil {
		data.AllowedArchitectures = convert.GraphToFrameworkBitmaskEnumAsSet(ctx, applicableArchitectures)
		tflog.Debug(ctx, "Set applicable architectures in state", map[string]any{
			"architectures": (*applicableArchitectures).String(),
		})
	}

	data.MinimumFreeDiskSpaceInMB = convert.GraphToFrameworkInt32(remoteResource.GetMinimumFreeDiskSpaceInMB())
	data.MinimumMemoryInMB = convert.GraphToFrameworkInt32(remoteResource.GetMinimumMemoryInMB())
	data.MinimumNumberOfProcessors = convert.GraphToFrameworkInt32(remoteResource.GetMinimumNumberOfProcessors())
	data.MinimumCpuSpeedInMHz = convert.GraphToFrameworkInt32(remoteResource.GetMinimumCpuSpeedInMHz())
	data.SetupFilePath = convert.GraphToFrameworkString(remoteResource.GetSetupFilePath())
	data.MinimumSupportedWindowsRelease = convert.GraphToFrameworkString(remoteResource.GetMinimumSupportedWindowsRelease())
	data.DisplayVersion = convert.GraphToFrameworkString(remoteResource.GetDisplayVersion())
	data.AllowAvailableUninstall = convert.GraphToFrameworkBool(remoteResource.GetAllowAvailableUninstall())

	// Detection Rules
	if detectionRules := remoteResource.GetDetectionRules(); detectionRules != nil {
		tflog.Debug(ctx, "Processing detection rules from API response", map[string]any{
			"count": len(detectionRules),
		})
	}

	// Requirement Rules
	if requirementRules := remoteResource.GetRequirementRules(); requirementRules != nil {
		tflog.Debug(ctx, "Processing requirement rules from API response", map[string]any{
			"count": len(requirementRules),
		})
	}

	// Rules - This is the unified field that should be populated from the API response
	if rules := remoteResource.GetRules(); rules != nil {
		data.Rules = make([]Win32LobAppRuleResourceModel, len(rules))
		for i, rule := range rules {
			// Set common rule properties
			switch ruleType := rule.(type) {
			case graphmodels.Win32LobAppRegistryRuleable:
				data.Rules[i] = Win32LobAppRuleResourceModel{
					RuleType:             convert.GraphToFrameworkEnum(ruleType.GetRuleType()),
					RuleSubType:          types.StringValue("registry"),
					Check32BitOn64System: convert.GraphToFrameworkBool(ruleType.GetCheck32BitOn64System()),
					KeyPath:              convert.GraphToFrameworkString(ruleType.GetKeyPath()),
					ValueName:            convert.GraphToFrameworkString(ruleType.GetValueName()),
					OperationType:        convert.GraphToFrameworkEnum(ruleType.GetOperationType()),
					LobAppRuleOperator:   convert.GraphToFrameworkEnum(ruleType.GetOperator()),
					ComparisonValue:      convert.GraphToFrameworkString(ruleType.GetComparisonValue()),
				}
			case graphmodels.Win32LobAppFileSystemRuleable:
				data.Rules[i] = Win32LobAppRuleResourceModel{
					RuleType:                convert.GraphToFrameworkEnum(ruleType.GetRuleType()),
					RuleSubType:             types.StringValue("file_system"),
					Check32BitOn64System:    convert.GraphToFrameworkBool(ruleType.GetCheck32BitOn64System()),
					Path:                    convert.GraphToFrameworkString(ruleType.GetPath()),
					FileOrFolderName:        convert.GraphToFrameworkString(ruleType.GetFileOrFolderName()),
					FileSystemOperationType: convert.GraphToFrameworkEnum(ruleType.GetOperationType()),
					LobAppRuleOperator:      convert.GraphToFrameworkEnum(ruleType.GetOperator()),
					ComparisonValue:         convert.GraphToFrameworkString(ruleType.GetComparisonValue()),
				}
			case graphmodels.Win32LobAppPowerShellScriptRuleable:
				data.Rules[i] = Win32LobAppRuleResourceModel{
					RuleType:                          convert.GraphToFrameworkEnum(ruleType.GetRuleType()),
					RuleSubType:                       types.StringValue("powershell_script"),
					DisplayName:                       convert.GraphToFrameworkString(ruleType.GetDisplayName()),
					EnforceSignatureCheck:             convert.GraphToFrameworkBool(ruleType.GetEnforceSignatureCheck()),
					RunAs32Bit:                        convert.GraphToFrameworkBool(ruleType.GetRunAs32Bit()),
					RunAsAccount:                      convert.GraphToFrameworkEnum(ruleType.GetRunAsAccount()),
					ScriptContent:                     helpers.DecodeBase64ToString(ctx, convert.GraphToFrameworkString(ruleType.GetScriptContent()).ValueString()),
					PowerShellScriptRuleOperationType: convert.GraphToFrameworkEnum(ruleType.GetOperationType()),
					LobAppRuleOperator:                convert.GraphToFrameworkEnum(ruleType.GetOperator()),
					ComparisonValue:                   convert.GraphToFrameworkString(ruleType.GetComparisonValue()),
				}
			default:
				tflog.Warn(ctx, "Unknown rule type", map[string]any{
					"ruleType": fmt.Sprintf("%T", rule),
				})
			}
		}
	} else {
		// If no rules are returned from the API but we have detection rules or requirement rules,
		// we need to convert those to the unified rules format
		var rulesFromDetectionAndRequirement []Win32LobAppRuleResourceModel

		// Convert detection rules to unified rules format if any exist
		if detectionRules := remoteResource.GetDetectionRules(); detectionRules != nil {
			for _, rule := range detectionRules {
				switch detectionRule := rule.(type) {
				case graphmodels.Win32LobAppRegistryDetectionable:
					rulesFromDetectionAndRequirement = append(rulesFromDetectionAndRequirement, Win32LobAppRuleResourceModel{
						RuleType:             types.StringValue("detection"),
						RuleSubType:          types.StringValue("registry"),
						Check32BitOn64System: convert.GraphToFrameworkBool(detectionRule.GetCheck32BitOn64System()),
						KeyPath:              convert.GraphToFrameworkString(detectionRule.GetKeyPath()),
						ValueName:            convert.GraphToFrameworkString(detectionRule.GetValueName()),
						OperationType:        convert.GraphToFrameworkEnum(detectionRule.GetDetectionType()),
						LobAppRuleOperator:   convert.GraphToFrameworkEnum(detectionRule.GetOperator()),
					})
				case graphmodels.Win32LobAppProductCodeDetectionable:
					rulesFromDetectionAndRequirement = append(rulesFromDetectionAndRequirement, Win32LobAppRuleResourceModel{
						RuleType:           types.StringValue("detection"),
						RuleSubType:        types.StringValue("msi_information"),
						LobAppRuleOperator: convert.GraphToFrameworkEnum(detectionRule.GetProductVersionOperator()),
						ComparisonValue:    convert.GraphToFrameworkString(detectionRule.GetProductVersion()),
					})
				case graphmodels.Win32LobAppFileSystemDetectionable:
					rulesFromDetectionAndRequirement = append(rulesFromDetectionAndRequirement, Win32LobAppRuleResourceModel{
						RuleType:                types.StringValue("detection"),
						RuleSubType:             types.StringValue("file_system"),
						Check32BitOn64System:    convert.GraphToFrameworkBool(detectionRule.GetCheck32BitOn64System()),
						Path:                    convert.GraphToFrameworkString(detectionRule.GetPath()),
						FileOrFolderName:        convert.GraphToFrameworkString(detectionRule.GetFileOrFolderName()),
						FileSystemOperationType: convert.GraphToFrameworkEnum(detectionRule.GetDetectionType()),
						LobAppRuleOperator:      convert.GraphToFrameworkEnum(detectionRule.GetOperator()),
						ComparisonValue:         convert.GraphToFrameworkString(detectionRule.GetDetectionValue()),
					})
				case graphmodels.Win32LobAppPowerShellScriptDetectionable:
					rulesFromDetectionAndRequirement = append(rulesFromDetectionAndRequirement, Win32LobAppRuleResourceModel{
						RuleType:              types.StringValue("detection"),
						RuleSubType:           types.StringValue("powershell_script"),
						EnforceSignatureCheck: convert.GraphToFrameworkBool(detectionRule.GetEnforceSignatureCheck()),
						RunAs32Bit:            convert.GraphToFrameworkBool(detectionRule.GetRunAs32Bit()),
						ScriptContent:         helpers.DecodeBase64ToString(ctx, convert.GraphToFrameworkString(detectionRule.GetScriptContent()).ValueString()),
						LobAppRuleOperator:    types.StringValue("notConfigured"),
					})
				}
			}
		}

		// Convert requirement rules to unified rules format if any exist
		if requirementRules := remoteResource.GetRequirementRules(); requirementRules != nil {
			for _, rule := range requirementRules {
				if registryRequirement, ok := rule.(graphmodels.Win32LobAppRegistryRequirementable); ok {
					rulesFromDetectionAndRequirement = append(rulesFromDetectionAndRequirement, Win32LobAppRuleResourceModel{
						RuleType:             types.StringValue("requirement"),
						RuleSubType:          types.StringValue("registry"),
						Check32BitOn64System: convert.GraphToFrameworkBool(registryRequirement.GetCheck32BitOn64System()),
						KeyPath:              convert.GraphToFrameworkString(registryRequirement.GetKeyPath()),
						ValueName:            convert.GraphToFrameworkString(registryRequirement.GetValueName()),
						OperationType:        convert.GraphToFrameworkEnum(registryRequirement.GetDetectionType()),
						LobAppRuleOperator:   convert.GraphToFrameworkEnum(registryRequirement.GetOperator()),
						ComparisonValue:      convert.GraphToFrameworkString(registryRequirement.GetDetectionValue()),
					})
				}
			}
		}

		// Set the combined rules in the state
		if len(rulesFromDetectionAndRequirement) > 0 {
			data.Rules = rulesFromDetectionAndRequirement
		}
	}

	// Install Experience
	if installExperience := remoteResource.GetInstallExperience(); installExperience != nil {
		data.InstallExperience = Win32LobAppInstallExperienceResourceModel{
			RunAsAccount:          convert.GraphToFrameworkEnum(installExperience.GetRunAsAccount()),
			DeviceRestartBehavior: convert.GraphToFrameworkEnum(installExperience.GetDeviceRestartBehavior()),
			MaxRunTimeInMinutes:   convert.GraphToFrameworkInt32(installExperience.GetMaxRunTimeInMinutes()),
		}
	}

	// Return Codes
	if returnCodes := remoteResource.GetReturnCodes(); returnCodes != nil {
		data.ReturnCodes = make([]Win32LobAppReturnCodeResourceModel, len(returnCodes))
		for i, code := range returnCodes {
			data.ReturnCodes[i] = Win32LobAppReturnCodeResourceModel{
				ReturnCode: convert.GraphToFrameworkInt32(code.GetReturnCode()),
				Type:       convert.GraphToFrameworkEnum(code.GetTypeEscaped()),
			}
		}
	}

	// MSI Information
	if msiInfo := remoteResource.GetMsiInformation(); msiInfo != nil {
		data.MsiInformation = &Win32LobAppMsiInformationResourceModel{
			ProductCode:    convert.GraphToFrameworkString(msiInfo.GetProductCode()),
			ProductVersion: convert.GraphToFrameworkString(msiInfo.GetProductVersion()),
			UpgradeCode:    convert.GraphToFrameworkString(msiInfo.GetUpgradeCode()),
			RequiresReboot: convert.GraphToFrameworkBool(msiInfo.GetRequiresReboot()),
			PackageType:    convert.GraphToFrameworkEnum(msiInfo.GetPackageType()),
			ProductName:    convert.GraphToFrameworkString(msiInfo.GetProductName()),
			Publisher:      convert.GraphToFrameworkString(msiInfo.GetPublisher()),
		}
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished stating resource %s with id %s", ResourceName, data.ID.ValueString()))

}
