package graphBetaWin32LobApp

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/construct"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

func constructResource(ctx context.Context, data *Win32LobAppResourceModel) (graphmodels.Win32LobAppable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s resource from model", ResourceName))

	requestBody := graphmodels.NewWin32LobApp()

	// Set string properties using the helper function
	construct.SetStringProperty(data.DisplayName, requestBody.SetDisplayName)
	construct.SetStringProperty(data.Description, requestBody.SetDescription)
	construct.SetStringProperty(data.Publisher, requestBody.SetPublisher)
	construct.SetStringProperty(data.FileName, requestBody.SetFileName)
	construct.SetStringProperty(data.InstallCommandLine, requestBody.SetInstallCommandLine)
	construct.SetStringProperty(data.UninstallCommandLine, requestBody.SetUninstallCommandLine)
	construct.SetStringProperty(data.SetupFilePath, requestBody.SetSetupFilePath)
	construct.SetStringProperty(data.CommittedContentVersion, requestBody.SetCommittedContentVersion)

	// Handle MinimumSupportedOperatingSystem
	if minOS := data.MinimumSupportedOperatingSystem; minOS != (WindowsMinimumOperatingSystemResourceModel{}) {
		minSupportedOS := graphmodels.NewWindowsMinimumOperatingSystem()

		construct.SetBoolProperty(minOS.V8_0, minSupportedOS.SetV80)
		construct.SetBoolProperty(minOS.V8_1, minSupportedOS.SetV81)
		construct.SetBoolProperty(minOS.V10_0, minSupportedOS.SetV100)
		construct.SetBoolProperty(minOS.V10_1607, minSupportedOS.SetV101607)
		construct.SetBoolProperty(minOS.V10_1703, minSupportedOS.SetV101703)
		construct.SetBoolProperty(minOS.V10_1709, minSupportedOS.SetV101709)
		construct.SetBoolProperty(minOS.V10_1803, minSupportedOS.SetV101803)
		construct.SetBoolProperty(minOS.V10_1809, minSupportedOS.SetV101809)
		construct.SetBoolProperty(minOS.V10_1903, minSupportedOS.SetV101903)
		construct.SetBoolProperty(minOS.V10_1909, minSupportedOS.SetV101909)
		construct.SetBoolProperty(minOS.V10_2004, minSupportedOS.SetV102004)
		construct.SetBoolProperty(minOS.V10_2H20, minSupportedOS.SetV102H20)
		construct.SetBoolProperty(minOS.V10_21H1, minSupportedOS.SetV1021H1)

		requestBody.SetMinimumSupportedOperatingSystem(minSupportedOS)
	}

	// Handle DetectionRules
	if len(data.DetectionRules) > 0 {
		detectionRules := make([]graphmodels.Win32LobAppDetectionable, len(data.DetectionRules))
		for i, rule := range data.DetectionRules {
			switch rule.RegistryDetectionType.ValueString() {
			case "registry":
				registryRule := graphmodels.NewWin32LobAppRegistryDetection()
				construct.SetBoolProperty(rule.Check32BitOn64System, registryRule.SetCheck32BitOn64System)
				construct.SetStringProperty(rule.KeyPath, registryRule.SetKeyPath)
				construct.SetStringProperty(rule.ValueName, registryRule.SetValueName)

				err := construct.ParseEnum(rule.RegistryDetectionOperator, graphmodels.ParseWin32LobAppDetectionOperator, registryRule.SetOperator)
				if err != nil {
					return nil, fmt.Errorf("failed to parse file system detection operator: %v", err)
				}

				err = construct.ParseEnum(rule.RegistryDetectionType, graphmodels.ParseWin32LobAppRegistryDetectionType, registryRule.SetDetectionType)
				if err != nil {
					return nil, fmt.Errorf("failed to parse registry detection type: %v", err)
				}

				construct.SetStringProperty(rule.DetectionValue, registryRule.SetDetectionValue)

				detectionRules[i] = registryRule

			case "msi_information":
				msiRule := graphmodels.NewWin32LobAppProductCodeDetection()
				construct.SetStringProperty(rule.ProductCode, msiRule.SetProductCode)
				construct.SetStringProperty(rule.ProductVersion, msiRule.SetProductVersion)

				err := construct.ParseEnum(rule.ProductVersionOperator, graphmodels.ParseWin32LobAppDetectionOperator, msiRule.SetProductVersionOperator)
				if err != nil {
					return nil, fmt.Errorf("failed to parse MSI product version: %v", err)
				}

				detectionRules[i] = msiRule

			case "file_system":
				fileRule := graphmodels.NewWin32LobAppFileSystemDetection()
				construct.SetBoolProperty(rule.Check32BitOn64System, fileRule.SetCheck32BitOn64System)
				construct.SetStringProperty(rule.FilePath, fileRule.SetPath)
				construct.SetStringProperty(rule.FileFolderName, fileRule.SetFileOrFolderName)
				construct.SetStringProperty(rule.DetectionValue, fileRule.SetDetectionValue)

				err := construct.ParseEnum(rule.FileSystemDetectionType, graphmodels.ParseWin32LobAppFileSystemDetectionType, fileRule.SetDetectionType)
				if err != nil {
					return nil, fmt.Errorf("failed to parse file system detection type: %v", err)
				}

				err = construct.ParseEnum(rule.FileSystemDetectionOperator, graphmodels.ParseWin32LobAppDetectionOperator, fileRule.SetOperator)
				if err != nil {
					return nil, fmt.Errorf("failed to parse file system detection operator: %v", err)
				}

				detectionRules[i] = fileRule
			case "powershell_script":
				powershellRule := graphmodels.NewWin32LobAppPowerShellScriptDetection()
				construct.SetStringProperty(rule.ScriptContent, powershellRule.SetScriptContent)
				construct.SetBoolProperty(rule.EnforceSignatureCheck, powershellRule.SetEnforceSignatureCheck)
				construct.SetBoolProperty(rule.RunAs32Bit, powershellRule.SetRunAs32Bit)

				detectionRules[i] = powershellRule
			}
		}
		requestBody.SetDetectionRules(detectionRules)
	}

	// Handle RequirementRules
	if len(data.RequirementRules) > 0 {
		requirementRules := make([]graphmodels.Win32LobAppRequirementable, len(data.RequirementRules))
		for i, rule := range data.RequirementRules {
			registryRequirement := graphmodels.NewWin32LobAppRegistryRequirement()

			construct.SetStringProperty(rule.KeyPath, registryRequirement.SetKeyPath)
			construct.SetStringProperty(rule.ValueName, registryRequirement.SetValueName)
			construct.SetBoolProperty(rule.Check32BitOn64System, registryRequirement.SetCheck32BitOn64System)

			err := construct.ParseEnum(rule.Operator, graphmodels.ParseWin32LobAppDetectionOperator, registryRequirement.SetOperator)
			if err != nil {
				return nil, fmt.Errorf("failed to parse registry requirement operator: %v", err)
			}

			err = construct.ParseEnum(rule.DetectionType, graphmodels.ParseWin32LobAppRegistryDetectionType, registryRequirement.SetDetectionType)
			if err != nil {
				return nil, fmt.Errorf("failed to parse registry detection type: %v", err)
			}

			construct.SetStringProperty(rule.DetectionValue, registryRequirement.SetDetectionValue)

			requirementRules[i] = registryRequirement
		}
		requestBody.SetRequirementRules(requirementRules)
	}

	// Handle Rules
	if len(data.Rules) > 0 {
		rules := make([]graphmodels.Win32LobAppRuleable, len(data.Rules))
		for i, rule := range data.Rules {
			registryRule := graphmodels.NewWin32LobAppRegistryRule()

			construct.SetStringProperty(rule.KeyPath, registryRule.SetKeyPath)
			construct.SetStringProperty(rule.ValueName, registryRule.SetValueName)
			construct.SetBoolProperty(rule.Check32BitOn64System, registryRule.SetCheck32BitOn64System)

			err := construct.ParseEnum(rule.Operator, graphmodels.ParseWin32LobAppRuleOperator, registryRule.SetOperator)
			if err != nil {
				return nil, fmt.Errorf("failed to parse registry rule operator: %v", err)
			}

			err = construct.ParseEnum(rule.OperationType, graphmodels.ParseWin32LobAppRegistryRuleOperationType, registryRule.SetOperationType)
			if err != nil {
				return nil, fmt.Errorf("failed to parse registry rule operation type: %v", err)
			}

			construct.SetStringProperty(rule.ComparisonValue, registryRule.SetComparisonValue)

			rules[i] = registryRule
		}
		requestBody.SetRules(rules)
	}

	// Handle Install Experience
	if installExperience := data.InstallExperience; installExperience != (Win32LobAppInstallExperienceResourceModel{}) {
		installExp := graphmodels.NewWin32LobAppInstallExperience()

		err := construct.ParseEnum(installExperience.RunAsAccount, graphmodels.ParseRunAsAccountType, installExp.SetRunAsAccount)
		if err != nil {
			return nil, fmt.Errorf("failed to parse RunAsAccountType: %v", err)
		}

		err = construct.ParseEnum(installExperience.DeviceRestartBehavior, graphmodels.ParseWin32LobAppRestartBehavior, installExp.SetDeviceRestartBehavior)
		if err != nil {
			return nil, fmt.Errorf("failed to parse DeviceRestartBehavior: %v", err)
		}

		construct.SetInt32Property(installExperience.MaxRunTimeInMinutes, installExp.SetMaxRunTimeInMinutes)

		requestBody.SetInstallExperience(installExp)
	}

	// Handle Return Codes
	if len(data.ReturnCodes) > 0 {
		returnCodes := make([]graphmodels.Win32LobAppReturnCodeable, len(data.ReturnCodes))
		for i, code := range data.ReturnCodes {
			returnCode := graphmodels.NewWin32LobAppReturnCode()

			construct.SetInt32Property(code.ReturnCode, returnCode.SetReturnCode)

			err := construct.ParseEnum(code.Type, graphmodels.ParseWin32LobAppReturnCodeType, returnCode.SetTypeEscaped)
			if err != nil {
				return nil, fmt.Errorf("failed to parse return code type: %v", err)
			}

			returnCodes[i] = returnCode
		}
		requestBody.SetReturnCodes(returnCodes)
	}

	// Handle MSI Information
	if msiInfo := data.MsiInformation; msiInfo != (Win32LobAppMsiInformationResourceModel{}) {
		msiInformation := graphmodels.NewWin32LobAppMsiInformation()

		construct.SetStringProperty(msiInfo.ProductCode, msiInformation.SetProductCode)
		construct.SetStringProperty(msiInfo.ProductVersion, msiInformation.SetProductVersion)
		construct.SetStringProperty(msiInfo.UpgradeCode, msiInformation.SetUpgradeCode)
		construct.SetBoolProperty(msiInfo.RequiresReboot, msiInformation.SetRequiresReboot)

		err := construct.ParseEnum[*graphmodels.Win32LobAppMsiPackageType](msiInfo.PackageType, graphmodels.ParseWin32LobAppMsiPackageType, msiInformation.SetPackageType)
		if err != nil {
			return nil, fmt.Errorf("failed to parse MSI package type: %v", err)
		}
		requestBody.SetMsiInformation(msiInformation)
	}

	if err := construct.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]interface{}{
			"error": err.Error(),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing %s resource", ResourceName))

	return requestBody, nil
}
