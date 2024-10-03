package graphBetaWin32LobApp

import (
	"context"
	"fmt"

	sharedmodels "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/shared_models/graph_beta"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/state"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

func MapRemoteStateToTerraform(ctx context.Context, data *Win32LobAppResourceModel, remoteResource graphmodels.Win32LobAppable) {
	if remoteResource == nil {
		tflog.Debug(ctx, "Remote resource is nil")
		return
	}

	tflog.Debug(ctx, "Starting to map remote state to Terraform state", map[string]interface{}{
		"resourceId": state.StringPtrToString(remoteResource.GetId()),
	})

	data.ID = types.StringValue(state.StringPtrToString(remoteResource.GetId()))
	data.DisplayName = types.StringValue(state.StringPtrToString(remoteResource.GetDisplayName()))
	data.Description = types.StringValue(state.StringPtrToString(remoteResource.GetDescription()))
	data.Publisher = types.StringValue(state.StringPtrToString(remoteResource.GetPublisher()))
	data.CreatedDateTime = state.TimeToString(remoteResource.GetCreatedDateTime())
	data.LastModifiedDateTime = state.TimeToString(remoteResource.GetLastModifiedDateTime())
	data.IsFeatured = state.BoolPtrToTypeBool(remoteResource.GetIsFeatured())
	data.PrivacyInformationUrl = types.StringValue(state.StringPtrToString(remoteResource.GetPrivacyInformationUrl()))
	data.InformationUrl = types.StringValue(state.StringPtrToString(remoteResource.GetInformationUrl()))
	data.Owner = types.StringValue(state.StringPtrToString(remoteResource.GetOwner()))
	data.Developer = types.StringValue(state.StringPtrToString(remoteResource.GetDeveloper()))
	data.Notes = types.StringValue(state.StringPtrToString(remoteResource.GetNotes()))
	data.UploadState = state.Int32PtrToTypeInt64(remoteResource.GetUploadState())
	data.PublishingState = state.EnumPtrToTypeString(remoteResource.GetPublishingState())
	data.IsAssigned = state.BoolPtrToTypeBool(remoteResource.GetIsAssigned())

	if largeIcon := remoteResource.GetLargeIcon(); largeIcon != nil {
		data.LargeIcon = sharedmodels.MimeContentResourceModel{
			Type:  types.StringValue(state.StringPtrToString(largeIcon.GetTypeEscaped())),
			Value: types.StringValue(state.ByteToString(largeIcon.GetValue())),
		}
	}

	data.RoleScopeTagIds = state.SliceToTypeStringSlice(remoteResource.GetRoleScopeTagIds())
	data.DependentAppCount = state.Int32PtrToTypeInt64(remoteResource.GetDependentAppCount())
	data.SupersedingAppCount = state.Int32PtrToTypeInt64(remoteResource.GetSupersedingAppCount())
	data.SupersededAppCount = state.Int32PtrToTypeInt64(remoteResource.GetSupersededAppCount())
	data.CommittedContentVersion = types.StringValue(state.StringPtrToString(remoteResource.GetCommittedContentVersion()))
	data.FileName = types.StringValue(state.StringPtrToString(remoteResource.GetFileName()))
	data.Size = state.Int64PtrToTypeInt64(remoteResource.GetSize())
	data.InstallCommandLine = types.StringValue(state.StringPtrToString(remoteResource.GetInstallCommandLine()))
	data.UninstallCommandLine = types.StringValue(state.StringPtrToString(remoteResource.GetUninstallCommandLine()))
	data.ApplicableArchitectures = state.EnumPtrToTypeString(remoteResource.GetApplicableArchitectures())
	data.MinimumFreeDiskSpaceInMB = state.Int32PtrToTypeInt64(remoteResource.GetMinimumFreeDiskSpaceInMB())
	data.MinimumMemoryInMB = state.Int32PtrToTypeInt64(remoteResource.GetMinimumMemoryInMB())
	data.MinimumNumberOfProcessors = state.Int32PtrToTypeInt64(remoteResource.GetMinimumNumberOfProcessors())
	data.MinimumCpuSpeedInMHz = state.Int32PtrToTypeInt64(remoteResource.GetMinimumCpuSpeedInMHz())
	data.SetupFilePath = types.StringValue(state.StringPtrToString(remoteResource.GetSetupFilePath()))
	data.MinimumSupportedWindowsRelease = types.StringValue(state.StringPtrToString(remoteResource.GetMinimumSupportedWindowsRelease()))
	data.DisplayVersion = types.StringValue(state.StringPtrToString(remoteResource.GetDisplayVersion()))
	data.AllowAvailableUninstall = state.BoolPtrToTypeBool(remoteResource.GetAllowAvailableUninstall())

	// Handle MinimumSupportedOperatingSystem
	minOS := remoteResource.GetMinimumSupportedOperatingSystem()
	if minOS != nil {
		data.MinimumSupportedOperatingSystem = WindowsMinimumOperatingSystemResourceModel{
			V8_0:     state.BoolPtrToTypeBool(minOS.GetV80()),
			V8_1:     state.BoolPtrToTypeBool(minOS.GetV81()),
			V10_0:    state.BoolPtrToTypeBool(minOS.GetV100()),
			V10_1607: state.BoolPtrToTypeBool(minOS.GetV101607()),
			V10_1703: state.BoolPtrToTypeBool(minOS.GetV101703()),
			V10_1709: state.BoolPtrToTypeBool(minOS.GetV101709()),
			V10_1803: state.BoolPtrToTypeBool(minOS.GetV101803()),
			V10_1809: state.BoolPtrToTypeBool(minOS.GetV101809()),
			V10_1903: state.BoolPtrToTypeBool(minOS.GetV101903()),
			V10_1909: state.BoolPtrToTypeBool(minOS.GetV101909()),
			V10_2004: state.BoolPtrToTypeBool(minOS.GetV102004()),
			V10_2H20: state.BoolPtrToTypeBool(minOS.GetV102H20()),
			V10_21H1: state.BoolPtrToTypeBool(minOS.GetV1021H1()),
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
					Check32BitOn64System:      state.BoolPtrToTypeBool(detectionRule.GetCheck32BitOn64System()),
					KeyPath:                   types.StringValue(state.StringPtrToString(detectionRule.GetKeyPath())),
					ValueName:                 types.StringValue(state.StringPtrToString(detectionRule.GetValueName())),
					RegistryDetectionOperator: state.EnumPtrToTypeString(detectionRule.GetOperator()),
					DetectionValue:            types.StringValue(state.StringPtrToString(detectionRule.GetDetectionValue())),
				}
			case graphmodels.Win32LobAppProductCodeDetectionable:
				data.DetectionRules[i] = Win32LobAppRegistryDetectionRulesResourceModel{
					DetectionType:          types.StringValue("msi_information"),
					ProductCode:            types.StringValue(state.StringPtrToString(detectionRule.GetProductCode())),
					ProductVersion:         types.StringValue(state.StringPtrToString(detectionRule.GetProductVersion())),
					ProductVersionOperator: state.EnumPtrToTypeString(detectionRule.GetProductVersionOperator()),
				}
			case graphmodels.Win32LobAppFileSystemDetectionable:
				data.DetectionRules[i] = Win32LobAppRegistryDetectionRulesResourceModel{
					DetectionType:               types.StringValue("file_system"),
					FileSystemDetectionType:     state.EnumPtrToTypeString(detectionRule.GetDetectionType()),
					FilePath:                    types.StringValue(state.StringPtrToString(detectionRule.GetPath())),
					FileFolderName:              types.StringValue(state.StringPtrToString(detectionRule.GetFileOrFolderName())),
					Check32BitOn64System:        state.BoolPtrToTypeBool(detectionRule.GetCheck32BitOn64System()),
					FileSystemDetectionOperator: state.EnumPtrToTypeString(detectionRule.GetOperator()),
					DetectionValue:              types.StringValue(state.StringPtrToString(detectionRule.GetDetectionValue())),
				}
			case graphmodels.Win32LobAppPowerShellScriptDetectionable:
				data.DetectionRules[i] = Win32LobAppRegistryDetectionRulesResourceModel{
					DetectionType:         types.StringValue("powershell_script"),
					ScriptContent:         types.StringValue(state.StringPtrToString(detectionRule.GetScriptContent())),
					EnforceSignatureCheck: state.BoolPtrToTypeBool(detectionRule.GetEnforceSignatureCheck()),
					RunAs32Bit:            state.BoolPtrToTypeBool(detectionRule.GetRunAs32Bit()),
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
					Check32BitOn64System: state.BoolPtrToTypeBool(registryRequirement.GetCheck32BitOn64System()),
					KeyPath:              types.StringValue(state.StringPtrToString(registryRequirement.GetKeyPath())),
					ValueName:            types.StringValue(state.StringPtrToString(registryRequirement.GetValueName())),
					Operator:             state.EnumPtrToTypeString(registryRequirement.GetOperator()),
					DetectionValue:       types.StringValue(state.StringPtrToString(registryRequirement.GetDetectionValue())),
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
					Check32BitOn64System: state.BoolPtrToTypeBool(registryRule.GetCheck32BitOn64System()),
					KeyPath:              types.StringValue(state.StringPtrToString(registryRule.GetKeyPath())),
					ValueName:            types.StringValue(state.StringPtrToString(registryRule.GetValueName())),
					OperationType:        state.EnumPtrToTypeString(registryRule.GetOperationType()),
					Operator:             state.EnumPtrToTypeString(registryRule.GetOperator()),
					ComparisonValue:      types.StringValue(state.StringPtrToString(registryRule.GetComparisonValue())),
				}
			}
		}
	}

	// Install Experience
	if installExperience := remoteResource.GetInstallExperience(); installExperience != nil {
		data.InstallExperience = Win32LobAppInstallExperienceResourceModel{
			RunAsAccount:          state.EnumPtrToTypeString(installExperience.GetRunAsAccount()),
			DeviceRestartBehavior: state.EnumPtrToTypeString(installExperience.GetDeviceRestartBehavior()),
			MaxRunTimeInMinutes:   state.Int32PtrToTypeInt64(installExperience.GetMaxRunTimeInMinutes()),
		}
	}

	// Return Codes
	if returnCodes := remoteResource.GetReturnCodes(); returnCodes != nil {
		data.ReturnCodes = make([]Win32LobAppReturnCodeResourceModel, len(returnCodes))
		for i, code := range returnCodes {
			data.ReturnCodes[i] = Win32LobAppReturnCodeResourceModel{
				ReturnCode: state.Int32PtrToTypeInt64(code.GetReturnCode()),
				Type:       state.EnumPtrToTypeString(code.GetTypeEscaped()),
			}
		}
	}

	// MSI Information
	if msiInfo := remoteResource.GetMsiInformation(); msiInfo != nil {
		data.MsiInformation = Win32LobAppMsiInformationResourceModel{
			ProductCode:    types.StringValue(state.StringPtrToString(msiInfo.GetProductCode())),
			ProductVersion: types.StringValue(state.StringPtrToString(msiInfo.GetProductVersion())),
			UpgradeCode:    types.StringValue(state.StringPtrToString(msiInfo.GetUpgradeCode())),
			RequiresReboot: state.BoolPtrToTypeBool(msiInfo.GetRequiresReboot()),
			PackageType:    state.EnumPtrToTypeString(msiInfo.GetPackageType()),
		}
	}

	tflog.Debug(ctx, "Finished mapping remote state to Terraform state", map[string]interface{}{
		"resourceId": data.ID.ValueString(),
	})
}
