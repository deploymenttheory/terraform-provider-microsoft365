package graphBetaWin32LobApp

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/state"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	models "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

func MapRemoteStateToTerraform(ctx context.Context, data *Win32LobAppResourceModel, remoteResource models.Win32LobAppable) {
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
		data.LargeIcon = MimeContent{
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
		data.MinimumSupportedOperatingSystem = WindowsMinimumOperatingSystemModel{
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

	// Handle RequirementRules
	if requirementRules := remoteResource.GetRequirementRules(); len(requirementRules) > 0 {
		data.RequirementRules = make([]RequirementRule, len(requirementRules))
		for i, rule := range requirementRules {
			data.RequirementRules[i] = RequirementRule{
				DetectionValue: types.StringValue(state.StringPtrToString(rule.GetDetectionValue())),
				Operator:       state.EnumPtrToTypeString(rule.GetOperator()),
			}

			// Type assertion to handle specific requirement types
			switch r := rule.(type) {
			case *models.Win32LobAppFileSystemRequirement:
				data.RequirementRules[i].RequirementType = types.StringValue("file")
				data.RequirementRules[i].Path = types.StringValue(state.StringPtrToString(r.GetPath()))
				data.RequirementRules[i].FileOrFolderName = types.StringValue(state.StringPtrToString(r.GetFileOrFolderName()))
				data.RequirementRules[i].Check32BitOn64System = state.BoolPtrToTypeBool(r.GetCheck32BitOn64System())
			case *models.Win32LobAppRegistryRequirement:
				data.RequirementRules[i].RequirementType = types.StringValue("registry")
				data.RequirementRules[i].KeyPath = types.StringValue(state.StringPtrToString(r.GetKeyPath()))
				data.RequirementRules[i].ValueName = types.StringValue(state.StringPtrToString(r.GetValueName()))
				data.RequirementRules[i].Check32BitOn64System = state.BoolPtrToTypeBool(r.GetCheck32BitOn64System())
			case *models.Win32LobAppPowerShellScriptRequirement:
				data.RequirementRules[i].RequirementType = types.StringValue("script")
				data.RequirementRules[i].ScriptContent = types.StringValue(state.StringPtrToString(r.GetScriptContent()))
				data.RequirementRules[i].EnforceSignatureCheck = state.BoolPtrToTypeBool(r.GetEnforceSignatureCheck())
				data.RequirementRules[i].RunAs32Bit = state.BoolPtrToTypeBool(r.GetRunAs32Bit())
			}
		}
	}

	// Handle RequirementRules
	requirementRules := remoteResource.GetRequirementRules()
	if len(requirementRules) > 0 {
		data.RequirementRules = make([]RequirementRule, len(requirementRules))
		for i, rule := range requirementRules {
			data.RequirementRules[i] = RequirementRule{
				RequirementType:      state.EnumPtrToTypeString(rule.GetRequirementType()),
				Path:                 types.StringValue(state.StringPtrToString(rule.GetPath())),
				FileOrFolderName:     types.StringValue(state.StringPtrToString(rule.GetFileOrFolderName())),
				Check32BitOn64System: state.BoolPtrToTypeBool(rule.GetCheck32BitOn64System()),
				Operator:             state.EnumPtrToTypeString(rule.GetOperator()),
				DetectionValue:       types.StringValue(state.StringPtrToString(rule.GetDetectionValue())),
			}
		}
	}

	// Handle Rules
	rules := remoteResource.GetRules()
	if len(rules) > 0 {
		data.Rules = make([]Rule, len(rules))
		for i, rule := range rules {
			data.Rules[i] = Rule{
				RuleType:             types.StringValue(state.StringPtrToString(rule.GetRuleType())),
				Check32BitOn64System: state.BoolPtrToTypeBool(rule.GetCheck32BitOn64System()),
				KeyPath:              types.StringValue(state.StringPtrToString(rule.GetKeyPath())),
				ValueName:            types.StringValue(state.StringPtrToString(rule.GetValueName())),
				OperationType:        types.StringValue(state.StringPtrToString(rule.GetOperationType())),
				Operator:             types.StringValue(state.StringPtrToString(rule.GetOperator())),
				ComparisonValue:      types.StringValue(state.StringPtrToString(rule.GetComparisonValue())),
			}
		}
	}

	// Handle InstallExperience
	if installExperience := remoteResource.GetInstallExperience(); installExperience != nil {
		data.InstallExperience = InstallExperience{
			RunAsAccount:          state.EnumPtrToTypeString(installExperience.GetRunAsAccount()),
			DeviceRestartBehavior: state.EnumPtrToTypeString(installExperience.GetDeviceRestartBehavior()),
		}
	}

	// Handle ReturnCodes
	returnCodes := remoteResource.GetReturnCodes()
	if len(returnCodes) > 0 {
		data.ReturnCodes = make([]ReturnCode, len(returnCodes))
		for i, code := range returnCodes {
			data.ReturnCodes[i] = ReturnCode{
				ReturnCode: state.Int32PtrToTypeInt64(code.GetReturnCode()),
				Type:       state.EnumPtrToTypeString(code.GetType()),
			}
		}
	}

	// Handle MsiInformation
	if msiInfo := remoteResource.GetMsiInformation(); msiInfo != nil {
		data.MsiInformation = MsiInformation{
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
