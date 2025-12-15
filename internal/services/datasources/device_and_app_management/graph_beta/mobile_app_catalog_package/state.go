package graphBetaMobileAppCatalogPackage

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// MapRemoteStateToDataSource maps the win32CatalogApp to the data source model
func MapRemoteStateToDataSource(ctx context.Context, mobileApp graphmodels.MobileAppable) MobileAppCatalogPackageModel {
	model := MobileAppCatalogPackageModel{}

	// Map base mobile app fields
	model.ID = convert.GraphToFrameworkString(mobileApp.GetId())
	model.DisplayName = convert.GraphToFrameworkString(mobileApp.GetDisplayName())
	model.Description = convert.GraphToFrameworkString(mobileApp.GetDescription())
	model.Publisher = convert.GraphToFrameworkString(mobileApp.GetPublisher())
	model.CreatedDateTime = convert.GraphToFrameworkTime(mobileApp.GetCreatedDateTime())
	model.LastModifiedDateTime = convert.GraphToFrameworkTime(mobileApp.GetLastModifiedDateTime())
	model.IsFeatured = convert.GraphToFrameworkBool(mobileApp.GetIsFeatured())
	model.PrivacyInformationUrl = convert.GraphToFrameworkString(mobileApp.GetPrivacyInformationUrl())
	model.InformationUrl = convert.GraphToFrameworkString(mobileApp.GetInformationUrl())
	model.Owner = convert.GraphToFrameworkString(mobileApp.GetOwner())
	model.Developer = convert.GraphToFrameworkString(mobileApp.GetDeveloper())
	model.Notes = convert.GraphToFrameworkString(mobileApp.GetNotes())
	model.IsAssigned = convert.GraphToFrameworkBool(mobileApp.GetIsAssigned())
	model.RoleScopeTagIds = convert.GraphToFrameworkStringSlice(mobileApp.GetRoleScopeTagIds())
	model.DependentAppCount = convert.GraphToFrameworkInt32(mobileApp.GetDependentAppCount())
	model.SupersedingAppCount = convert.GraphToFrameworkInt32(mobileApp.GetSupersedingAppCount())
	model.SupersededAppCount = convert.GraphToFrameworkInt32(mobileApp.GetSupersededAppCount())

	// Map upload state (custom handling for enum to int32)
	if mobileApp.GetUploadState() != nil {
		model.UploadState = convert.GraphToFrameworkInt32(mobileApp.GetUploadState())
	}

	// Map publishing state enum
	model.PublishingState = convert.GraphToFrameworkEnum(mobileApp.GetPublishingState())

	// Check if it's a Win32CatalogApp to access win32-specific fields
	if win32App, ok := mobileApp.(graphmodels.Win32CatalogAppable); ok {
		model.CommittedContentVersion = convert.GraphToFrameworkString(win32App.GetCommittedContentVersion())
		model.FileName = convert.GraphToFrameworkString(win32App.GetFileName())
		model.Size = convert.GraphToFrameworkInt64(win32App.GetSize())
		model.InstallCommandLine = convert.GraphToFrameworkString(win32App.GetInstallCommandLine())
		model.UninstallCommandLine = convert.GraphToFrameworkString(win32App.GetUninstallCommandLine())
		model.ApplicableArchitectures = convert.GraphToFrameworkEnum(win32App.GetApplicableArchitectures())
		model.AllowedArchitectures = convert.GraphToFrameworkEnum(win32App.GetAllowedArchitectures())
		model.MinimumFreeDiskSpaceInMB = convert.GraphToFrameworkInt32(win32App.GetMinimumFreeDiskSpaceInMB())
		model.MinimumMemoryInMB = convert.GraphToFrameworkInt32(win32App.GetMinimumMemoryInMB())
		model.MinimumNumberOfProcessors = convert.GraphToFrameworkInt32(win32App.GetMinimumNumberOfProcessors())
		model.MinimumCpuSpeedInMHz = convert.GraphToFrameworkInt32(win32App.GetMinimumCpuSpeedInMHz())
		model.SetupFilePath = convert.GraphToFrameworkString(win32App.GetSetupFilePath())
		model.MinimumSupportedWindowsRelease = convert.GraphToFrameworkString(win32App.GetMinimumSupportedWindowsRelease())
		model.DisplayVersion = convert.GraphToFrameworkString(win32App.GetDisplayVersion())
		model.AllowAvailableUninstall = convert.GraphToFrameworkBool(win32App.GetAllowAvailableUninstall())
		model.MobileAppCatalogPackageId = convert.GraphToFrameworkString(win32App.GetMobileAppCatalogPackageId())

		// Map rules
		if win32App.GetRules() != nil {
			var rules []RuleModel
			for _, rule := range win32App.GetRules() {
				ruleModel := mapRule(ctx, rule)
				rules = append(rules, ruleModel)
			}
			model.Rules = rules
		} else {
			model.Rules = []RuleModel{}
		}

		// Map install experience
		if win32App.GetInstallExperience() != nil {
			model.InstallExperience = mapInstallExperience(ctx, win32App.GetInstallExperience())
		}

		// Map return codes
		if win32App.GetReturnCodes() != nil {
			var returnCodes []ReturnCodeModel
			for _, rc := range win32App.GetReturnCodes() {
				returnCodeModel := mapReturnCode(ctx, rc)
				returnCodes = append(returnCodes, returnCodeModel)
			}
			model.ReturnCodes = returnCodes
		} else {
			model.ReturnCodes = []ReturnCodeModel{}
		}

		// Map MSI information
		if win32App.GetMsiInformation() != nil {
			model.MsiInformation = mapMsiInformation(ctx, win32App.GetMsiInformation())
		}
	} else {
		tflog.Warn(ctx, "MobileApp is not a Win32CatalogApp, some fields will be null")
	}

	return model
}

// mapRule maps a Win32LobAppRule to RuleModel
func mapRule(ctx context.Context, rule graphmodels.Win32LobAppRuleable) RuleModel {
	model := RuleModel{
		ODataType: convert.GraphToFrameworkString(rule.GetOdataType()),
		RuleType:  convert.GraphToFrameworkEnum(rule.GetRuleType()),
	}

	// Check if it's a file system rule
	if fsRule, ok := rule.(graphmodels.Win32LobAppFileSystemRuleable); ok {
		model.Path = convert.GraphToFrameworkString(fsRule.GetPath())
		model.FileOrFolderName = convert.GraphToFrameworkString(fsRule.GetFileOrFolderName())
		model.Check32BitOn64System = convert.GraphToFrameworkBool(fsRule.GetCheck32BitOn64System())
		model.OperationType = convert.GraphToFrameworkEnum(fsRule.GetOperationType())
		model.Operator = convert.GraphToFrameworkEnum(fsRule.GetOperator())
		model.ComparisonValue = convert.GraphToFrameworkString(fsRule.GetComparisonValue())
	}

	// Check if it's a registry rule
	if regRule, ok := rule.(graphmodels.Win32LobAppRegistryRuleable); ok {
		model.Check32BitOn64System = convert.GraphToFrameworkBool(regRule.GetCheck32BitOn64System())
		model.KeyPath = convert.GraphToFrameworkString(regRule.GetKeyPath())
		model.ValueName = convert.GraphToFrameworkString(regRule.GetValueName())
		model.OperationType = convert.GraphToFrameworkEnum(regRule.GetOperationType())
		model.Operator = convert.GraphToFrameworkEnum(regRule.GetOperator())
		model.ComparisonValue = convert.GraphToFrameworkString(regRule.GetComparisonValue())
	}

	return model
}

// mapInstallExperience maps Win32LobAppInstallExperience to InstallExperienceModel
func mapInstallExperience(ctx context.Context, installExp graphmodels.Win32LobAppInstallExperienceable) *InstallExperienceModel {
	return &InstallExperienceModel{
		RunAsAccount:          convert.GraphToFrameworkEnum(installExp.GetRunAsAccount()),
		MaxRunTimeInMinutes:   convert.GraphToFrameworkInt32(installExp.GetMaxRunTimeInMinutes()),
		DeviceRestartBehavior: convert.GraphToFrameworkEnum(installExp.GetDeviceRestartBehavior()),
	}
}

// mapReturnCode maps Win32LobAppReturnCode to ReturnCodeModel
func mapReturnCode(ctx context.Context, rc graphmodels.Win32LobAppReturnCodeable) ReturnCodeModel {
	return ReturnCodeModel{
		ReturnCode: convert.GraphToFrameworkInt32(rc.GetReturnCode()),
		Type:       convert.GraphToFrameworkEnum(rc.GetTypeEscaped()),
	}
}

// mapMsiInformation maps Win32LobAppMsiInformation to MsiInformationModel
func mapMsiInformation(ctx context.Context, msiInfo graphmodels.Win32LobAppMsiInformationable) *MsiInformationModel {
	return &MsiInformationModel{
		ProductCode:    convert.GraphToFrameworkString(msiInfo.GetProductCode()),
		ProductVersion: convert.GraphToFrameworkString(msiInfo.GetProductVersion()),
		UpgradeCode:    convert.GraphToFrameworkString(msiInfo.GetUpgradeCode()),
		RequiresReboot: convert.GraphToFrameworkBool(msiInfo.GetRequiresReboot()),
		PackageType:    convert.GraphToFrameworkEnum(msiInfo.GetPackageType()),
		ProductName:    convert.GraphToFrameworkString(msiInfo.GetProductName()),
		Publisher:      convert.GraphToFrameworkString(msiInfo.GetPublisher()),
	}
}
