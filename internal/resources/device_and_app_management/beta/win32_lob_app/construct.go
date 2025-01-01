package graphBetaWin32LobApp

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/constructors"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

func constructResource(ctx context.Context, data *Win32LobAppResourceModel) (graphmodels.Win32LobAppable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s resource from model", ResourceName))

	requestBody := graphmodels.NewWin32LobApp()

	// Set string properties using the helper function
	constructors.SetStringProperty(data.DisplayName, requestBody.SetDisplayName)
	constructors.SetStringProperty(data.Description, requestBody.SetDescription)
	constructors.SetStringProperty(data.Publisher, requestBody.SetPublisher)
	constructors.SetStringProperty(data.FileName, requestBody.SetFileName)
	constructors.SetStringProperty(data.InstallCommandLine, requestBody.SetInstallCommandLine)
	constructors.SetStringProperty(data.UninstallCommandLine, requestBody.SetUninstallCommandLine)
	constructors.SetStringProperty(data.SetupFilePath, requestBody.SetSetupFilePath)
	constructors.SetStringProperty(data.CommittedContentVersion, requestBody.SetCommittedContentVersion)

	// Handle MinimumSupportedOperatingSystem
	if minOS := data.MinimumSupportedOperatingSystem; minOS != (WindowsMinimumOperatingSystemResourceModel{}) {
		minSupportedOS := graphmodels.NewWindowsMinimumOperatingSystem()

		constructors.SetBoolProperty(minOS.V8_0, minSupportedOS.SetV80)
		constructors.SetBoolProperty(minOS.V8_1, minSupportedOS.SetV81)
		constructors.SetBoolProperty(minOS.V10_0, minSupportedOS.SetV100)
		constructors.SetBoolProperty(minOS.V10_1607, minSupportedOS.SetV101607)
		constructors.SetBoolProperty(minOS.V10_1703, minSupportedOS.SetV101703)
		constructors.SetBoolProperty(minOS.V10_1709, minSupportedOS.SetV101709)
		constructors.SetBoolProperty(minOS.V10_1803, minSupportedOS.SetV101803)
		constructors.SetBoolProperty(minOS.V10_1809, minSupportedOS.SetV101809)
		constructors.SetBoolProperty(minOS.V10_1903, minSupportedOS.SetV101903)
		constructors.SetBoolProperty(minOS.V10_1909, minSupportedOS.SetV101909)
		constructors.SetBoolProperty(minOS.V10_2004, minSupportedOS.SetV102004)
		constructors.SetBoolProperty(minOS.V10_2H20, minSupportedOS.SetV102H20)
		constructors.SetBoolProperty(minOS.V10_21H1, minSupportedOS.SetV1021H1)

		requestBody.SetMinimumSupportedOperatingSystem(minSupportedOS)
	}

	// Handle DetectionRules
	if len(data.DetectionRules) > 0 {
		detectionRules := make([]graphmodels.Win32LobAppDetectionable, len(data.DetectionRules))
		for i, rule := range data.DetectionRules {
			switch rule.RegistryDetectionType.ValueString() {
			case "registry":
				registryRule := graphmodels.NewWin32LobAppRegistryDetection()
				constructors.SetBoolProperty(rule.Check32BitOn64System, registryRule.SetCheck32BitOn64System)
				constructors.SetStringProperty(rule.KeyPath, registryRule.SetKeyPath)
				constructors.SetStringProperty(rule.ValueName, registryRule.SetValueName)

				err := constructors.SetEnumProperty(rule.RegistryDetectionOperator, graphmodels.ParseWin32LobAppDetectionOperator, registryRule.SetOperator)
				if err != nil {
					return nil, fmt.Errorf("failed to parse file system detection operator: %v", err)
				}

				err = constructors.SetEnumProperty(rule.RegistryDetectionType, graphmodels.ParseWin32LobAppRegistryDetectionType, registryRule.SetDetectionType)
				if err != nil {
					return nil, fmt.Errorf("failed to parse registry detection type: %v", err)
				}

				constructors.SetStringProperty(rule.DetectionValue, registryRule.SetDetectionValue)

				detectionRules[i] = registryRule

			case "msi_information":
				msiRule := graphmodels.NewWin32LobAppProductCodeDetection()
				constructors.SetStringProperty(rule.ProductCode, msiRule.SetProductCode)
				constructors.SetStringProperty(rule.ProductVersion, msiRule.SetProductVersion)

				err := constructors.SetEnumProperty(rule.ProductVersionOperator, graphmodels.ParseWin32LobAppDetectionOperator, msiRule.SetProductVersionOperator)
				if err != nil {
					return nil, fmt.Errorf("failed to parse MSI product version: %v", err)
				}

				detectionRules[i] = msiRule

			case "file_system":
				fileRule := graphmodels.NewWin32LobAppFileSystemDetection()
				constructors.SetBoolProperty(rule.Check32BitOn64System, fileRule.SetCheck32BitOn64System)
				constructors.SetStringProperty(rule.FilePath, fileRule.SetPath)
				constructors.SetStringProperty(rule.FileFolderName, fileRule.SetFileOrFolderName)
				constructors.SetStringProperty(rule.DetectionValue, fileRule.SetDetectionValue)

				err := constructors.SetEnumProperty(rule.FileSystemDetectionType, graphmodels.ParseWin32LobAppFileSystemDetectionType, fileRule.SetDetectionType)
				if err != nil {
					return nil, fmt.Errorf("failed to parse file system detection type: %v", err)
				}

				err = constructors.SetEnumProperty(rule.FileSystemDetectionOperator, graphmodels.ParseWin32LobAppDetectionOperator, fileRule.SetOperator)
				if err != nil {
					return nil, fmt.Errorf("failed to parse file system detection operator: %v", err)
				}

				detectionRules[i] = fileRule
			case "powershell_script":
				powershellRule := graphmodels.NewWin32LobAppPowerShellScriptDetection()
				constructors.SetStringProperty(rule.ScriptContent, powershellRule.SetScriptContent)
				constructors.SetBoolProperty(rule.EnforceSignatureCheck, powershellRule.SetEnforceSignatureCheck)
				constructors.SetBoolProperty(rule.RunAs32Bit, powershellRule.SetRunAs32Bit)

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

			constructors.SetStringProperty(rule.KeyPath, registryRequirement.SetKeyPath)
			constructors.SetStringProperty(rule.ValueName, registryRequirement.SetValueName)
			constructors.SetBoolProperty(rule.Check32BitOn64System, registryRequirement.SetCheck32BitOn64System)

			err := constructors.SetEnumProperty(rule.Operator, graphmodels.ParseWin32LobAppDetectionOperator, registryRequirement.SetOperator)
			if err != nil {
				return nil, fmt.Errorf("failed to parse registry requirement operator: %v", err)
			}

			err = constructors.SetEnumProperty(rule.DetectionType, graphmodels.ParseWin32LobAppRegistryDetectionType, registryRequirement.SetDetectionType)
			if err != nil {
				return nil, fmt.Errorf("failed to parse registry detection type: %v", err)
			}

			constructors.SetStringProperty(rule.DetectionValue, registryRequirement.SetDetectionValue)

			requirementRules[i] = registryRequirement
		}
		requestBody.SetRequirementRules(requirementRules)
	}

	// Handle Rules
	if len(data.Rules) > 0 {
		rules := make([]graphmodels.Win32LobAppRuleable, len(data.Rules))
		for i, rule := range data.Rules {
			registryRule := graphmodels.NewWin32LobAppRegistryRule()

			constructors.SetStringProperty(rule.KeyPath, registryRule.SetKeyPath)
			constructors.SetStringProperty(rule.ValueName, registryRule.SetValueName)
			constructors.SetBoolProperty(rule.Check32BitOn64System, registryRule.SetCheck32BitOn64System)

			err := constructors.SetEnumProperty(rule.Operator, graphmodels.ParseWin32LobAppRuleOperator, registryRule.SetOperator)
			if err != nil {
				return nil, fmt.Errorf("failed to parse registry rule operator: %v", err)
			}

			err = constructors.SetEnumProperty(rule.OperationType, graphmodels.ParseWin32LobAppRegistryRuleOperationType, registryRule.SetOperationType)
			if err != nil {
				return nil, fmt.Errorf("failed to parse registry rule operation type: %v", err)
			}

			constructors.SetStringProperty(rule.ComparisonValue, registryRule.SetComparisonValue)

			rules[i] = registryRule
		}
		requestBody.SetRules(rules)
	}

	// Handle Install Experience
	if installExperience := data.InstallExperience; installExperience != (Win32LobAppInstallExperienceResourceModel{}) {
		installExp := graphmodels.NewWin32LobAppInstallExperience()

		err := constructors.SetEnumProperty(installExperience.RunAsAccount, graphmodels.ParseRunAsAccountType, installExp.SetRunAsAccount)
		if err != nil {
			return nil, fmt.Errorf("failed to parse RunAsAccountType: %v", err)
		}

		err = constructors.SetEnumProperty(installExperience.DeviceRestartBehavior, graphmodels.ParseWin32LobAppRestartBehavior, installExp.SetDeviceRestartBehavior)
		if err != nil {
			return nil, fmt.Errorf("failed to parse DeviceRestartBehavior: %v", err)
		}

		constructors.SetInt32Property(installExperience.MaxRunTimeInMinutes, installExp.SetMaxRunTimeInMinutes)

		requestBody.SetInstallExperience(installExp)
	}

	// Handle Return Codes
	if len(data.ReturnCodes) > 0 {
		returnCodes := make([]graphmodels.Win32LobAppReturnCodeable, len(data.ReturnCodes))
		for i, code := range data.ReturnCodes {
			returnCode := graphmodels.NewWin32LobAppReturnCode()

			constructors.SetInt32Property(code.ReturnCode, returnCode.SetReturnCode)

			err := constructors.SetEnumProperty(code.Type, graphmodels.ParseWin32LobAppReturnCodeType, returnCode.SetTypeEscaped)
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

		constructors.SetStringProperty(msiInfo.ProductCode, msiInformation.SetProductCode)
		constructors.SetStringProperty(msiInfo.ProductVersion, msiInformation.SetProductVersion)
		constructors.SetStringProperty(msiInfo.UpgradeCode, msiInformation.SetUpgradeCode)
		constructors.SetBoolProperty(msiInfo.RequiresReboot, msiInformation.SetRequiresReboot)

		err := constructors.SetEnumProperty[*graphmodels.Win32LobAppMsiPackageType](msiInfo.PackageType, graphmodels.ParseWin32LobAppMsiPackageType, msiInformation.SetPackageType)
		if err != nil {
			return nil, fmt.Errorf("failed to parse MSI package type: %v", err)
		}
		requestBody.SetMsiInformation(msiInformation)
	}

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]interface{}{
			"error": err.Error(),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing %s resource", ResourceName))

	return requestBody, nil
}
