package graphBetaWin32LobApp

import (
	"context"
	"fmt"

	sharedmodels "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/shared_models/graph_beta"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/state"
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
		"resourceId": state.StringPtrToString(remoteResource.GetId()),
	})

	data.ID = types.StringPointerValue(remoteResource.GetId())
	data.DisplayName = types.StringPointerValue(remoteResource.GetDisplayName())
	data.Description = types.StringPointerValue(remoteResource.GetDescription())
	data.Publisher = types.StringPointerValue(remoteResource.GetPublisher())
	data.CreatedDateTime = state.TimeToString(remoteResource.GetCreatedDateTime())
	data.LastModifiedDateTime = state.TimeToString(remoteResource.GetLastModifiedDateTime())
	data.IsFeatured = types.BoolPointerValue(remoteResource.GetIsFeatured())
	data.PrivacyInformationUrl = types.StringPointerValue(remoteResource.GetPrivacyInformationUrl())
	data.InformationUrl = types.StringPointerValue(remoteResource.GetInformationUrl())
	data.Owner = types.StringPointerValue(remoteResource.GetOwner())
	data.Developer = types.StringPointerValue(remoteResource.GetDeveloper())
	data.Notes = types.StringPointerValue(remoteResource.GetNotes())
	data.UploadState = state.Int32PtrToTypeInt32(remoteResource.GetUploadState())
	data.PublishingState = state.EnumPtrToTypeString(remoteResource.GetPublishingState())
	data.IsAssigned = types.BoolPointerValue(remoteResource.GetIsAssigned())

	if largeIcon := remoteResource.GetLargeIcon(); largeIcon != nil {
		data.LargeIcon = sharedmodels.MimeContentResourceModel{
			Type:  types.StringPointerValue(largeIcon.GetTypeEscaped()),
			Value: types.StringValue(state.ByteStringToBase64(largeIcon.GetValue())),
		}
	}

	data.RoleScopeTagIds = state.StringSliceToSet(ctx, remoteResource.GetRoleScopeTagIds())
	data.DependentAppCount = state.Int32PtrToTypeInt32(remoteResource.GetDependentAppCount())
	data.SupersedingAppCount = state.Int32PtrToTypeInt32(remoteResource.GetSupersedingAppCount())
	data.SupersededAppCount = state.Int32PtrToTypeInt32(remoteResource.GetSupersededAppCount())
	data.CommittedContentVersion = types.StringPointerValue(remoteResource.GetCommittedContentVersion())
	data.FileName = types.StringPointerValue(remoteResource.GetFileName())
	data.Size = state.Int64PtrToTypeInt64(remoteResource.GetSize())
	data.InstallCommandLine = types.StringPointerValue(remoteResource.GetInstallCommandLine())
	data.UninstallCommandLine = types.StringPointerValue(remoteResource.GetUninstallCommandLine())
	data.ApplicableArchitectures = state.EnumPtrToTypeString(remoteResource.GetApplicableArchitectures())
	data.MinimumFreeDiskSpaceInMB = state.Int32PtrToTypeInt32(remoteResource.GetMinimumFreeDiskSpaceInMB())
	data.MinimumMemoryInMB = state.Int32PtrToTypeInt32(remoteResource.GetMinimumMemoryInMB())
	data.MinimumNumberOfProcessors = state.Int32PtrToTypeInt32(remoteResource.GetMinimumNumberOfProcessors())
	data.MinimumCpuSpeedInMHz = state.Int32PtrToTypeInt32(remoteResource.GetMinimumCpuSpeedInMHz())
	data.SetupFilePath = types.StringPointerValue(remoteResource.GetSetupFilePath())
	data.MinimumSupportedWindowsRelease = types.StringPointerValue(remoteResource.GetMinimumSupportedWindowsRelease())
	data.DisplayVersion = types.StringPointerValue(remoteResource.GetDisplayVersion())
	data.AllowAvailableUninstall = types.BoolPointerValue(remoteResource.GetAllowAvailableUninstall())

	// MinimumSupportedOperatingSystem
	minOS := remoteResource.GetMinimumSupportedOperatingSystem()
	if minOS != nil {
		data.MinimumSupportedOperatingSystem = WindowsMinimumOperatingSystemResourceModel{
			V8_0:     types.BoolPointerValue(minOS.GetV80()),
			V8_1:     types.BoolPointerValue(minOS.GetV81()),
			V10_0:    types.BoolPointerValue(minOS.GetV100()),
			V10_1607: types.BoolPointerValue(minOS.GetV101607()),
			V10_1703: types.BoolPointerValue(minOS.GetV101703()),
			V10_1709: types.BoolPointerValue(minOS.GetV101709()),
			V10_1803: types.BoolPointerValue(minOS.GetV101803()),
			V10_1809: types.BoolPointerValue(minOS.GetV101809()),
			V10_1903: types.BoolPointerValue(minOS.GetV101903()),
			V10_1909: types.BoolPointerValue(minOS.GetV101909()),
			V10_2004: types.BoolPointerValue(minOS.GetV102004()),
			V10_2H20: types.BoolPointerValue(minOS.GetV102H20()),
			V10_21H1: types.BoolPointerValue(minOS.GetV1021H1()),
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
					RegistryDetectionType:     state.EnumPtrToTypeString(detectionRule.GetDetectionType()),
					Check32BitOn64System:      types.BoolPointerValue(detectionRule.GetCheck32BitOn64System()),
					KeyPath:                   types.StringPointerValue(detectionRule.GetKeyPath()),
					ValueName:                 types.StringPointerValue(detectionRule.GetValueName()),
					RegistryDetectionOperator: state.EnumPtrToTypeString(detectionRule.GetOperator()),
					DetectionValue:            types.StringPointerValue(detectionRule.GetDetectionValue()),
				}
			case graphmodels.Win32LobAppProductCodeDetectionable:
				data.DetectionRules[i] = Win32LobAppRegistryDetectionRulesResourceModel{
					DetectionType:          types.StringValue("msi_information"),
					ProductCode:            types.StringPointerValue(detectionRule.GetProductCode()),
					ProductVersion:         types.StringPointerValue(detectionRule.GetProductVersion()),
					ProductVersionOperator: state.EnumPtrToTypeString(detectionRule.GetProductVersionOperator()),
				}
			case graphmodels.Win32LobAppFileSystemDetectionable:
				data.DetectionRules[i] = Win32LobAppRegistryDetectionRulesResourceModel{
					DetectionType:               types.StringValue("file_system"),
					FileSystemDetectionType:     state.EnumPtrToTypeString(detectionRule.GetDetectionType()),
					FilePath:                    types.StringPointerValue(detectionRule.GetPath()),
					FileFolderName:              types.StringPointerValue(detectionRule.GetFileOrFolderName()),
					Check32BitOn64System:        types.BoolPointerValue(detectionRule.GetCheck32BitOn64System()),
					FileSystemDetectionOperator: state.EnumPtrToTypeString(detectionRule.GetOperator()),
					DetectionValue:              types.StringPointerValue(detectionRule.GetDetectionValue()),
				}
			case graphmodels.Win32LobAppPowerShellScriptDetectionable:
				data.DetectionRules[i] = Win32LobAppRegistryDetectionRulesResourceModel{
					DetectionType:         types.StringValue("powershell_script"),
					ScriptContent:         types.StringPointerValue(detectionRule.GetScriptContent()),
					EnforceSignatureCheck: types.BoolPointerValue(detectionRule.GetEnforceSignatureCheck()),
					RunAs32Bit:            types.BoolPointerValue(detectionRule.GetRunAs32Bit()),
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
					Check32BitOn64System: types.BoolPointerValue(registryRequirement.GetCheck32BitOn64System()),
					KeyPath:              types.StringPointerValue(registryRequirement.GetKeyPath()),
					ValueName:            types.StringPointerValue(registryRequirement.GetValueName()),
					Operator:             state.EnumPtrToTypeString(registryRequirement.GetOperator()),
					DetectionValue:       types.StringPointerValue(registryRequirement.GetDetectionValue()),
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
					RuleType:             state.EnumPtrToTypeString(registryRule.GetRuleType()),
					Check32BitOn64System: types.BoolPointerValue(registryRule.GetCheck32BitOn64System()),
					KeyPath:              types.StringPointerValue(registryRule.GetKeyPath()),
					ValueName:            types.StringPointerValue(registryRule.GetValueName()),
					OperationType:        state.EnumPtrToTypeString(registryRule.GetOperationType()),
					Operator:             state.EnumPtrToTypeString(registryRule.GetOperator()),
					ComparisonValue:      types.StringPointerValue(registryRule.GetComparisonValue()),
				}
			}
		}
	}

	// Install Experience
	if installExperience := remoteResource.GetInstallExperience(); installExperience != nil {
		data.InstallExperience = Win32LobAppInstallExperienceResourceModel{
			RunAsAccount:          state.EnumPtrToTypeString(installExperience.GetRunAsAccount()),
			DeviceRestartBehavior: state.EnumPtrToTypeString(installExperience.GetDeviceRestartBehavior()),
			MaxRunTimeInMinutes:   state.Int32PtrToTypeInt32(installExperience.GetMaxRunTimeInMinutes()),
		}
	}

	// Return Codes
	if returnCodes := remoteResource.GetReturnCodes(); returnCodes != nil {
		data.ReturnCodes = make([]Win32LobAppReturnCodeResourceModel, len(returnCodes))
		for i, code := range returnCodes {
			data.ReturnCodes[i] = Win32LobAppReturnCodeResourceModel{
				ReturnCode: state.Int32PtrToTypeInt32(code.GetReturnCode()),
				Type:       state.EnumPtrToTypeString(code.GetTypeEscaped()),
			}
		}
	}

	// MSI Information
	if msiInfo := remoteResource.GetMsiInformation(); msiInfo != nil {
		data.MsiInformation = Win32LobAppMsiInformationResourceModel{
			ProductCode:    types.StringPointerValue(msiInfo.GetProductCode()),
			ProductVersion: types.StringPointerValue(msiInfo.GetProductVersion()),
			UpgradeCode:    types.StringPointerValue(msiInfo.GetUpgradeCode()),
			RequiresReboot: types.BoolPointerValue(msiInfo.GetRequiresReboot()),
			PackageType:    state.EnumPtrToTypeString(msiInfo.GetPackageType()),
		}
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished stating resource %s with id %s", ResourceName, data.ID.ValueString()))

}
