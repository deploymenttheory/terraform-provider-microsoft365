package graphBetaWin32LobApp

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	sharedmodels "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/shared_models/graph_beta"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// MapRemoteStateToTerraform maps the remote state to the win32lobapp to Terraform state
func MapRemoteStateToTerraform(ctx context.Context, data *Win32LobAppResourceModel, remoteResource graphmodels.Win32LobAppable) {
	if remoteResource == nil {
		tflog.Debug(ctx, "Remote resource is nil")
		return
	}

	tflog.Debug(ctx, "Starting to map remote state to Terraform state", map[string]interface{}{
		"resourceId": convert.GraphToFrameworkString(remoteResource.GetId()).ValueString(),
	})

	data.ID = convert.GraphToFrameworkString(remoteResource.GetId())
	data.DisplayName = convert.GraphToFrameworkString(remoteResource.GetDisplayName())
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

	if largeIcon := remoteResource.GetLargeIcon(); largeIcon != nil {
		data.LargeIcon = sharedmodels.MimeContentResourceModel{
			Type:  convert.GraphToFrameworkString(largeIcon.GetTypeEscaped()),
			Value: convert.GraphToFrameworkBytes(largeIcon.GetValue()),
		}
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
	data.ApplicableArchitectures = convert.GraphToFrameworkEnum(remoteResource.GetApplicableArchitectures())
	data.MinimumFreeDiskSpaceInMB = convert.GraphToFrameworkInt32(remoteResource.GetMinimumFreeDiskSpaceInMB())
	data.MinimumMemoryInMB = convert.GraphToFrameworkInt32(remoteResource.GetMinimumMemoryInMB())
	data.MinimumNumberOfProcessors = convert.GraphToFrameworkInt32(remoteResource.GetMinimumNumberOfProcessors())
	data.MinimumCpuSpeedInMHz = convert.GraphToFrameworkInt32(remoteResource.GetMinimumCpuSpeedInMHz())
	data.SetupFilePath = convert.GraphToFrameworkString(remoteResource.GetSetupFilePath())
	data.MinimumSupportedWindowsRelease = convert.GraphToFrameworkString(remoteResource.GetMinimumSupportedWindowsRelease())
	data.DisplayVersion = convert.GraphToFrameworkString(remoteResource.GetDisplayVersion())
	data.AllowAvailableUninstall = convert.GraphToFrameworkBool(remoteResource.GetAllowAvailableUninstall())

	// MinimumSupportedOperatingSystem
	minOS := remoteResource.GetMinimumSupportedOperatingSystem()
	if minOS != nil {
		data.MinimumSupportedOperatingSystem = WindowsMinimumOperatingSystemResourceModel{
			V8_0:     convert.GraphToFrameworkBool(minOS.GetV80()),
			V8_1:     convert.GraphToFrameworkBool(minOS.GetV81()),
			V10_0:    convert.GraphToFrameworkBool(minOS.GetV100()),
			V10_1607: convert.GraphToFrameworkBool(minOS.GetV101607()),
			V10_1703: convert.GraphToFrameworkBool(minOS.GetV101703()),
			V10_1709: convert.GraphToFrameworkBool(minOS.GetV101709()),
			V10_1803: convert.GraphToFrameworkBool(minOS.GetV101803()),
			V10_1809: convert.GraphToFrameworkBool(minOS.GetV101809()),
			V10_1903: convert.GraphToFrameworkBool(minOS.GetV101903()),
			V10_1909: convert.GraphToFrameworkBool(minOS.GetV101909()),
			V10_2004: convert.GraphToFrameworkBool(minOS.GetV102004()),
			V10_2H20: convert.GraphToFrameworkBool(minOS.GetV102H20()),
			V10_21H1: convert.GraphToFrameworkBool(minOS.GetV1021H1()),
		}
	}

	// Detection Rules
	if detectionRules := remoteResource.GetDetectionRules(); detectionRules != nil {
		data.DetectionRules = make([]Win32LobAppRegistryDetectionRulesResourceModel, len(detectionRules))
		for i, rule := range detectionRules {
			switch detectionRule := rule.(type) {
			case graphmodels.Win32LobAppRegistryDetectionable:
				data.DetectionRules[i] = Win32LobAppRegistryDetectionRulesResourceModel{
					DetectionType:             types.StringValue("registry"),
					RegistryDetectionType:     convert.GraphToFrameworkEnum(detectionRule.GetDetectionType()),
					Check32BitOn64System:      convert.GraphToFrameworkBool(detectionRule.GetCheck32BitOn64System()),
					KeyPath:                   convert.GraphToFrameworkString(detectionRule.GetKeyPath()),
					ValueName:                 convert.GraphToFrameworkString(detectionRule.GetValueName()),
					RegistryDetectionOperator: convert.GraphToFrameworkEnum(detectionRule.GetOperator()),
					DetectionValue:            convert.GraphToFrameworkString(detectionRule.GetDetectionValue()),
				}
			case graphmodels.Win32LobAppProductCodeDetectionable:
				data.DetectionRules[i] = Win32LobAppRegistryDetectionRulesResourceModel{
					DetectionType:          types.StringValue("msi_information"),
					ProductCode:            convert.GraphToFrameworkString(detectionRule.GetProductCode()),
					ProductVersion:         convert.GraphToFrameworkString(detectionRule.GetProductVersion()),
					ProductVersionOperator: convert.GraphToFrameworkEnum(detectionRule.GetProductVersionOperator()),
				}
			case graphmodels.Win32LobAppFileSystemDetectionable:
				data.DetectionRules[i] = Win32LobAppRegistryDetectionRulesResourceModel{
					DetectionType:               types.StringValue("file_system"),
					FileSystemDetectionType:     convert.GraphToFrameworkEnum(detectionRule.GetDetectionType()),
					FilePath:                    convert.GraphToFrameworkString(detectionRule.GetPath()),
					FileFolderName:              convert.GraphToFrameworkString(detectionRule.GetFileOrFolderName()),
					Check32BitOn64System:        convert.GraphToFrameworkBool(detectionRule.GetCheck32BitOn64System()),
					FileSystemDetectionOperator: convert.GraphToFrameworkEnum(detectionRule.GetOperator()),
					DetectionValue:              convert.GraphToFrameworkString(detectionRule.GetDetectionValue()),
				}
			case graphmodels.Win32LobAppPowerShellScriptDetectionable:
				data.DetectionRules[i] = Win32LobAppRegistryDetectionRulesResourceModel{
					DetectionType:         types.StringValue("powershell_script"),
					ScriptContent:         convert.GraphToFrameworkString(detectionRule.GetScriptContent()),
					EnforceSignatureCheck: convert.GraphToFrameworkBool(detectionRule.GetEnforceSignatureCheck()),
					RunAs32Bit:            convert.GraphToFrameworkBool(detectionRule.GetRunAs32Bit()),
				}
			default:
				tflog.Warn(ctx, "Unknown detection rule type", map[string]interface{}{
					"ruleType": fmt.Sprintf("%T", rule),
				})
			}
		}
	}

	// Requirement Rules
	if requirementRules := remoteResource.GetRequirementRules(); requirementRules != nil {
		data.RequirementRules = make([]Win32LobAppRegistryRequirementResourceModel, len(requirementRules))
		for i, rule := range requirementRules {
			if registryRequirement, ok := rule.(graphmodels.Win32LobAppRegistryRequirementable); ok {
				data.RequirementRules[i] = Win32LobAppRegistryRequirementResourceModel{
					Check32BitOn64System: convert.GraphToFrameworkBool(registryRequirement.GetCheck32BitOn64System()),
					KeyPath:              convert.GraphToFrameworkString(registryRequirement.GetKeyPath()),
					ValueName:            convert.GraphToFrameworkString(registryRequirement.GetValueName()),
					Operator:             convert.GraphToFrameworkEnum(registryRequirement.GetOperator()),
					DetectionValue:       convert.GraphToFrameworkString(registryRequirement.GetDetectionValue()),
					DetectionType:        types.StringValue("registry"),
				}
			}
		}
	}

	// Rules
	if rules := remoteResource.GetRules(); rules != nil {
		data.Rules = make([]Win32LobAppRegistryRuleResourceModel, len(rules))
		for i, rule := range rules {
			if registryRule, ok := rule.(graphmodels.Win32LobAppRegistryRuleable); ok {
				data.Rules[i] = Win32LobAppRegistryRuleResourceModel{
					RuleType:             convert.GraphToFrameworkEnum(registryRule.GetRuleType()),
					Check32BitOn64System: convert.GraphToFrameworkBool(registryRule.GetCheck32BitOn64System()),
					KeyPath:              convert.GraphToFrameworkString(registryRule.GetKeyPath()),
					ValueName:            convert.GraphToFrameworkString(registryRule.GetValueName()),
					OperationType:        convert.GraphToFrameworkEnum(registryRule.GetOperationType()),
					Operator:             convert.GraphToFrameworkEnum(registryRule.GetOperator()),
					ComparisonValue:      convert.GraphToFrameworkString(registryRule.GetComparisonValue()),
				}
			}
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
		data.MsiInformation = Win32LobAppMsiInformationResourceModel{
			ProductCode:    convert.GraphToFrameworkString(msiInfo.GetProductCode()),
			ProductVersion: convert.GraphToFrameworkString(msiInfo.GetProductVersion()),
			UpgradeCode:    convert.GraphToFrameworkString(msiInfo.GetUpgradeCode()),
			RequiresReboot: convert.GraphToFrameworkBool(msiInfo.GetRequiresReboot()),
			PackageType:    convert.GraphToFrameworkEnum(msiInfo.GetPackageType()),
		}
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished stating resource %s with id %s", ResourceName, data.ID.ValueString()))

}
