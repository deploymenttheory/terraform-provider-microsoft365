package graphBetaWin32LobApp

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/construct"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/utilities"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

func constructResource(ctx context.Context, data *Win32LobAppResourceModel) (graphmodels.Win32LobAppable, error) {
	tflog.Debug(ctx, "Constructing Win32LobApp resource")
	construct.DebugPrintStruct(ctx, "Constructed Win32LobApp Resource from model", data)

	win32LobApp := graphmodels.NewWin32LobApp()

	if !data.DisplayName.IsNull() && !data.DisplayName.IsUnknown() {
		displayName := data.DisplayName.ValueString()
		win32LobApp.SetDisplayName(&displayName)
	}

	if !data.Description.IsNull() && !data.Description.IsUnknown() {
		description := data.Description.ValueString()
		win32LobApp.SetDescription(&description)
	}

	if !data.Publisher.IsNull() && !data.Publisher.IsUnknown() {
		publisher := data.Publisher.ValueString()
		win32LobApp.SetPublisher(&publisher)
	}

	if !data.FileName.IsNull() && !data.FileName.IsUnknown() {
		fileName := data.FileName.ValueString()
		win32LobApp.SetFileName(&fileName)
	}

	if !data.InstallCommandLine.IsNull() && !data.InstallCommandLine.IsUnknown() {
		installCommandLine := data.InstallCommandLine.ValueString()
		win32LobApp.SetInstallCommandLine(&installCommandLine)
	}

	if !data.UninstallCommandLine.IsNull() && !data.UninstallCommandLine.IsUnknown() {
		uninstallCommandLine := data.UninstallCommandLine.ValueString()
		win32LobApp.SetUninstallCommandLine(&uninstallCommandLine)
	}

	if !data.SetupFilePath.IsNull() && !data.SetupFilePath.IsUnknown() {
		setupFilePath := data.SetupFilePath.ValueString()
		win32LobApp.SetSetupFilePath(&setupFilePath)
	}

	if !data.CommittedContentVersion.IsNull() && !data.CommittedContentVersion.IsUnknown() {
		contentVersion := data.CommittedContentVersion.ValueString()
		win32LobApp.SetCommittedContentVersion(&contentVersion)
	}

	// Handle MinimumSupportedOperatingSystem
	if minOS := data.MinimumSupportedOperatingSystem; minOS != (WindowsMinimumOperatingSystemResourceModel{}) {
		minSupportedOS := graphmodels.NewWindowsMinimumOperatingSystem()

		if !minOS.V8_0.IsNull() && !minOS.V8_0.IsUnknown() {
			minSupportedOS.SetV80(utilities.BoolPtr(minOS.V8_0.ValueBool()))
		}
		if !minOS.V8_1.IsNull() && !minOS.V8_1.IsUnknown() {
			minSupportedOS.SetV81(utilities.BoolPtr(minOS.V8_1.ValueBool()))
		}
		if !minOS.V10_0.IsNull() && !minOS.V10_0.IsUnknown() {
			minSupportedOS.SetV100(utilities.BoolPtr(minOS.V10_0.ValueBool()))
		}
		if !minOS.V10_1607.IsNull() && !minOS.V10_1607.IsUnknown() {
			minSupportedOS.SetV101607(utilities.BoolPtr(minOS.V10_1607.ValueBool()))
		}
		if !minOS.V10_1703.IsNull() && !minOS.V10_1703.IsUnknown() {
			minSupportedOS.SetV101703(utilities.BoolPtr(minOS.V10_1703.ValueBool()))
		}
		if !minOS.V10_1709.IsNull() && !minOS.V10_1709.IsUnknown() {
			minSupportedOS.SetV101709(utilities.BoolPtr(minOS.V10_1709.ValueBool()))
		}
		if !minOS.V10_1803.IsNull() && !minOS.V10_1803.IsUnknown() {
			minSupportedOS.SetV101803(utilities.BoolPtr(minOS.V10_1803.ValueBool()))
		}
		if !minOS.V10_1809.IsNull() && !minOS.V10_1809.IsUnknown() {
			minSupportedOS.SetV101809(utilities.BoolPtr(minOS.V10_1809.ValueBool()))
		}
		if !minOS.V10_1903.IsNull() && !minOS.V10_1903.IsUnknown() {
			minSupportedOS.SetV101903(utilities.BoolPtr(minOS.V10_1903.ValueBool()))
		}
		if !minOS.V10_1909.IsNull() && !minOS.V10_1909.IsUnknown() {
			minSupportedOS.SetV101909(utilities.BoolPtr(minOS.V10_1909.ValueBool()))
		}
		if !minOS.V10_2004.IsNull() && !minOS.V10_2004.IsUnknown() {
			minSupportedOS.SetV102004(utilities.BoolPtr(minOS.V10_2004.ValueBool()))
		}
		if !minOS.V10_2H20.IsNull() && !minOS.V10_2H20.IsUnknown() {
			minSupportedOS.SetV102H20(utilities.BoolPtr(minOS.V10_2H20.ValueBool()))
		}
		if !minOS.V10_21H1.IsNull() && !minOS.V10_21H1.IsUnknown() {
			minSupportedOS.SetV1021H1(utilities.BoolPtr(minOS.V10_21H1.ValueBool()))
		}

		win32LobApp.SetMinimumSupportedOperatingSystem(minSupportedOS)
	}

	// Handle DetectionRules
	if len(data.DetectionRules) > 0 {
		detectionRules := make([]graphmodels.Win32LobAppDetectionable, len(data.DetectionRules))
		for i, rule := range data.DetectionRules {
			registryRule := graphmodels.NewWin32LobAppRegistryDetection()
			registryRule.SetKeyPath(rule.KeyPath.ValueString())
			registryRule.SetValueName(rule.ValueName.ValueString())
			registryRule.SetCheck32BitOn64System(rule.Check32BitOn64System.ValueBool())
			registryRule.SetOperator(graphmodels.ParseWin32LobAppDetectionOperator(rule.Operator.ValueString()))
			registryRule.SetDetectionValue(rule.DetectionValue.ValueString())
			detectionRules[i] = registryRule
		}
		win32LobApp.SetDetectionRules(detectionRules)
	}

	// Handle RequirementRules
	if len(data.RequirementRules) > 0 {
		requirementRules := make([]graphmodels.Win32LobAppRequirementable, len(data.RequirementRules))
		for i, rule := range data.RequirementRules {
			registryRequirement := graphmodels.NewWin32LobAppRegistryRequirement()
			registryRequirement.SetKeyPath(rule.KeyPath.ValueString())
			registryRequirement.SetValueName(rule.ValueName.ValueString())
			registryRequirement.SetCheck32BitOn64System(rule.Check32BitOn64System.ValueBool())
			registryRequirement.SetOperator(graphmodels.ParseWin32LobAppDetectionOperator(rule.Operator.ValueString()))
			registryRequirement.SetDetectionValue(rule.DetectionValue.ValueString())
			requirementRules[i] = registryRequirement
		}
		win32LobApp.SetRequirementRules(requirementRules)
	}

	// Handle Rules
	if len(data.Rules) > 0 {
		rules := make([]graphmodels.Win32LobAppRuleable, len(data.Rules))
		for i, rule := range data.Rules {
			registryRule := graphmodels.NewWin32LobAppRegistryRule()
			registryRule.SetKeyPath(rule.KeyPath.ValueString())
			registryRule.SetValueName(rule.ValueName.ValueString())
			registryRule.SetCheck32BitOn64System(rule.Check32BitOn64System.ValueBool())
			registryRule.SetOperator(graphmodels.ParseWin32LobAppDetectionOperator(rule.Operator.ValueString()))
			registryRule.SetComparisonValue(rule.ComparisonValue.ValueString())
			registryRule.SetOperationType(graphmodels.ParseWin32LobAppRuleOperationType(rule.OperationType.ValueString()))
			rules[i] = registryRule
		}
		win32LobApp.SetRules(rules)
	}

	// Handle Install Experience
	if installExperience := data.InstallExperience; installExperience != (Win32LobAppInstallExperienceResourceModel{}) {
		installExp := graphmodels.NewWin32LobAppInstallExperience()
		installExp.SetRunAsAccount(graphmodels.ParseRunAsAccountType(installExperience.RunAsAccount.ValueString()))
		installExp.SetDeviceRestartBehavior(graphmodels.ParseWin32LobAppRestartBehavior(installExperience.DeviceRestartBehavior.ValueString()))
		installExp.SetMaxRunTimeInMinutes(installExperience.MaxRunTimeInMinutes.ValueInt64())
		win32LobApp.SetInstallExperience(installExp)
	}

	// Handle Return Codes
	if len(data.ReturnCodes) > 0 {
		returnCodes := make([]graphmodels.Win32LobAppReturnCodeable, len(data.ReturnCodes))
		for i, code := range data.ReturnCodes {
			returnCode := graphmodels.NewWin32LobAppReturnCode()
			returnCode.SetReturnCode(code.ReturnCode.ValueInt64())
			returnCode.SetTypeEscaped(graphmodels.ParseWin32LobAppReturnCodeType(code.Type.ValueString()))
			returnCodes[i] = returnCode
		}
		win32LobApp.SetReturnCodes(returnCodes)
	}

	// Handle MSI Information
	if msiInfo := data.MsiInformation; msiInfo != (Win32LobAppMsiInformationResourceModel{}) {
		msiInformation := graphmodels.NewWin32LobAppMsiInformation()
		msiInformation.SetProductCode(msiInfo.ProductCode.ValueString())
		msiInformation.SetProductVersion(msiInfo.ProductVersion.ValueString())
		msiInformation.SetUpgradeCode(msiInfo.UpgradeCode.ValueString())
		msiInformation.SetRequiresReboot(msiInfo.RequiresReboot.ValueBool())
		msiInformation.SetPackageType(graphmodels.ParseWin32LobAppMsiPackageType(msiInfo.PackageType.ValueString()))
		win32LobApp.SetMsiInformation(msiInformation)
	}

	tflog.Debug(ctx, "Finished constructing Win32LobApp resource")
	return win32LobApp, nil
}
