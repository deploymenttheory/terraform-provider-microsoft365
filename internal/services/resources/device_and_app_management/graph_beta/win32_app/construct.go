package graphBetaWin32App

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	helpers "github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	sharedConstructors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors/graph_beta/device_and_app_management"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	helpersCrud "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud/graph_beta/device_and_app_management"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

func constructResource(ctx context.Context, data *Win32LobAppResourceModel, installerSourcePath string) (graphmodels.Win32LobAppable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s resource from model", ResourceName))

	requestBody := graphmodels.NewWin32LobApp()

	convert.FrameworkToGraphString(data.DisplayName, requestBody.SetDisplayName)
	convert.FrameworkToGraphString(data.Description, requestBody.SetDescription)
	convert.FrameworkToGraphString(data.Publisher, requestBody.SetPublisher)
	convert.FrameworkToGraphString(data.FileName, requestBody.SetFileName)
	convert.FrameworkToGraphString(data.InstallCommandLine, requestBody.SetInstallCommandLine)
	convert.FrameworkToGraphString(data.UninstallCommandLine, requestBody.SetUninstallCommandLine)
	convert.FrameworkToGraphString(data.SetupFilePath, requestBody.SetSetupFilePath)
	convert.FrameworkToGraphString(data.CommittedContentVersion, requestBody.SetCommittedContentVersion)
	convert.FrameworkToGraphString(data.DisplayVersion, requestBody.SetDisplayVersion)
	convert.FrameworkToGraphString(data.Developer, requestBody.SetDeveloper)
	convert.FrameworkToGraphString(data.InformationUrl, requestBody.SetInformationUrl)
	convert.FrameworkToGraphString(data.PrivacyInformationUrl, requestBody.SetPrivacyInformationUrl)
	convert.FrameworkToGraphString(data.Notes, requestBody.SetNotes)
	convert.FrameworkToGraphString(data.Owner, requestBody.SetOwner)
	convert.FrameworkToGraphString(data.MinimumSupportedWindowsRelease, requestBody.SetMinimumSupportedWindowsRelease)
	convert.FrameworkToGraphBool(data.AllowAvailableUninstall, requestBody.SetAllowAvailableUninstall)
	convert.FrameworkToGraphBool(data.IsFeatured, requestBody.SetIsFeatured)

	if err := convert.FrameworkToGraphStringSet(ctx, data.RoleScopeTagIds, requestBody.SetRoleScopeTagIds); err != nil {
		return nil, fmt.Errorf("failed to set role scope tags: %s", err)
	}

	// Handle applicable architectures
	if err := convert.FrameworkToGraphBitmaskEnumFromSet(ctx, data.AllowedArchitectures,
		graphmodels.ParseWindowsArchitecture, requestBody.SetAllowedArchitectures); err != nil {
		tflog.Warn(ctx, "Failed to set applicable architectures", map[string]any{
			"error": err.Error(),
		})
	}

	// Handle app icon (either from file path or web source)
	if data.AppIcon != nil {
		largeIcon, tempFiles, err := sharedConstructors.ConstructMobileAppIcon(ctx, data.AppIcon)
		if err != nil {
			return nil, err
		}

		defer func() {
			for _, tempFile := range tempFiles {
				helpersCrud.CleanupTempFile(ctx, tempFile)
			}
		}()

		requestBody.SetLargeIcon(largeIcon)
	}

	// For creating resources, we need the installer file to extract metadata
	// Verify the installer path is provided and the file exists
	if installerSourcePath == "" {
		return nil, fmt.Errorf("installer source path is empty; a valid file path is required")
	}

	if _, err := os.Stat(installerSourcePath); err != nil {
		return nil, fmt.Errorf("installer file not found at path %s: %w", installerSourcePath, err)
	}

	filename := filepath.Base(installerSourcePath)
	tflog.Debug(ctx, fmt.Sprintf("Using filename from installer path: %s", filename))
	convert.FrameworkToGraphString(types.StringValue(filename), requestBody.SetFileName)

	if len(data.DetectionRules) > 0 {
		detectionRules := make([]graphmodels.Win32LobAppDetectionable, len(data.DetectionRules))
		for i, rule := range data.DetectionRules {
			switch rule.DetectionType.ValueString() {
			case "registry":
				registryRule := graphmodels.NewWin32LobAppRegistryDetection()
				convert.FrameworkToGraphBool(rule.Check32BitOn64System, registryRule.SetCheck32BitOn64System)
				convert.FrameworkToGraphString(rule.KeyPath, registryRule.SetKeyPath)
				convert.FrameworkToGraphString(rule.ValueName, registryRule.SetValueName)

				err := convert.FrameworkToGraphEnum(rule.RegistryDetectionOperator, graphmodels.ParseWin32LobAppDetectionOperator, registryRule.SetOperator)
				if err != nil {
					return nil, fmt.Errorf("failed to parse file system detection operator: %v", err)
				}

				err = convert.FrameworkToGraphEnum(rule.RegistryDetectionType, graphmodels.ParseWin32LobAppRegistryDetectionType, registryRule.SetDetectionType)
				if err != nil {
					return nil, fmt.Errorf("failed to parse registry detection type: %v", err)
				}

				convert.FrameworkToGraphString(rule.DetectionValue, registryRule.SetDetectionValue)

				detectionRules[i] = registryRule

			case "msi_information":
				msiRule := graphmodels.NewWin32LobAppProductCodeDetection()
				convert.FrameworkToGraphString(rule.ProductCode, msiRule.SetProductCode)
				convert.FrameworkToGraphString(rule.ProductVersion, msiRule.SetProductVersion)

				err := convert.FrameworkToGraphEnum(rule.ProductVersionOperator, graphmodels.ParseWin32LobAppDetectionOperator, msiRule.SetProductVersionOperator)
				if err != nil {
					return nil, fmt.Errorf("failed to parse MSI product version: %v", err)
				}

				detectionRules[i] = msiRule

			case "file_system":
				fileRule := graphmodels.NewWin32LobAppFileSystemDetection()
				convert.FrameworkToGraphBool(rule.Check32BitOn64System, fileRule.SetCheck32BitOn64System)
				convert.FrameworkToGraphString(rule.FilePath, fileRule.SetPath)
				convert.FrameworkToGraphString(rule.FileFolderName, fileRule.SetFileOrFolderName)
				convert.FrameworkToGraphString(rule.DetectionValue, fileRule.SetDetectionValue)

				err := convert.FrameworkToGraphEnum(rule.FileSystemDetectionType, graphmodels.ParseWin32LobAppFileSystemDetectionType, fileRule.SetDetectionType)
				if err != nil {
					return nil, fmt.Errorf("failed to parse file system detection type: %v", err)
				}

				err = convert.FrameworkToGraphEnum(rule.FileSystemDetectionOperator, graphmodels.ParseWin32LobAppDetectionOperator, fileRule.SetOperator)
				if err != nil {
					return nil, fmt.Errorf("failed to parse file system detection operator: %v", err)
				}

				detectionRules[i] = fileRule
			case "powershell_script":
				powershellRule := graphmodels.NewWin32LobAppPowerShellScriptDetection()

				// Convert script content to base64 before setting the
				// converted content as a string
				if !rule.ScriptContent.IsNull() && !rule.ScriptContent.IsUnknown() {
					scriptContent := rule.ScriptContent.ValueString()

					base64Content, err := helpers.StringToBase64(scriptContent)
					if err == nil {

						convert.FrameworkToGraphString(types.StringValue(base64Content), powershellRule.SetScriptContent)
					} else {

						convert.FrameworkToGraphString(rule.ScriptContent, powershellRule.SetScriptContent)
					}
				}

				convert.FrameworkToGraphBool(rule.EnforceSignatureCheck, powershellRule.SetEnforceSignatureCheck)
				convert.FrameworkToGraphBool(rule.RunAs32Bit, powershellRule.SetRunAs32Bit)

				detectionRules[i] = powershellRule
			}
		}
		requestBody.SetDetectionRules(detectionRules)
	}

	if len(data.RequirementRules) > 0 {
		requirementRules := make([]graphmodels.Win32LobAppRequirementable, len(data.RequirementRules))
		for i, rule := range data.RequirementRules {
			// For now, only handle registry requirements as that's what the SDK model we have supports
			// The RequirementType field indicates the type, but we'll focus on registry requirements
			registryRequirement := graphmodels.NewWin32LobAppRegistryRequirement()

			convert.FrameworkToGraphString(rule.KeyPath, registryRequirement.SetKeyPath)
			convert.FrameworkToGraphString(rule.ValueName, registryRequirement.SetValueName)
			convert.FrameworkToGraphBool(rule.Check32BitOn64System, registryRequirement.SetCheck32BitOn64System)

			err := convert.FrameworkToGraphEnum(rule.Operator, graphmodels.ParseWin32LobAppDetectionOperator, registryRequirement.SetOperator)
			if err != nil {
				return nil, fmt.Errorf("failed to parse registry requirement operator: %v", err)
			}

			err = convert.FrameworkToGraphEnum(rule.DetectionType, graphmodels.ParseWin32LobAppRegistryDetectionType, registryRequirement.SetDetectionType)
			if err != nil {
				return nil, fmt.Errorf("failed to parse registry detection type: %v", err)
			}

			convert.FrameworkToGraphString(rule.DetectionValue, registryRequirement.SetDetectionValue)

			requirementRules[i] = registryRequirement
		}
		requestBody.SetRequirementRules(requirementRules)
	}

	if len(data.Rules) > 0 {
		rules := make([]graphmodels.Win32LobAppRuleable, len(data.Rules))
		for i, rule := range data.Rules {
			// Determine which type of rule to create based on the rule_sub_type
			switch rule.RuleSubType.ValueString() {
			case "registry":
				registryRule := graphmodels.NewWin32LobAppRegistryRule()

				// Set common rule properties
				err := convert.FrameworkToGraphEnum(rule.RuleType, graphmodels.ParseWin32LobAppRuleType, registryRule.SetRuleType)
				if err != nil {
					return nil, fmt.Errorf("failed to parse rule type: %v", err)
				}

				convert.FrameworkToGraphString(rule.KeyPath, registryRule.SetKeyPath)
				convert.FrameworkToGraphString(rule.ValueName, registryRule.SetValueName)
				convert.FrameworkToGraphBool(rule.Check32BitOn64System, registryRule.SetCheck32BitOn64System)

				err = convert.FrameworkToGraphEnum(rule.LobAppRuleOperator, graphmodels.ParseWin32LobAppRuleOperator, registryRule.SetOperator)
				if err != nil {
					return nil, fmt.Errorf("failed to parse registry rule operator: %v", err)
				}

				err = convert.FrameworkToGraphEnum(rule.OperationType, graphmodels.ParseWin32LobAppRegistryRuleOperationType, registryRule.SetOperationType)
				if err != nil {
					return nil, fmt.Errorf("failed to parse registry rule operation type: %v", err)
				}

				convert.FrameworkToGraphString(rule.ComparisonValue, registryRule.SetComparisonValue)

				rules[i] = registryRule

			case "file_system":
				fileSystemRule := graphmodels.NewWin32LobAppFileSystemRule()

				// Set common rule properties
				err := convert.FrameworkToGraphEnum(rule.RuleType, graphmodels.ParseWin32LobAppRuleType, fileSystemRule.SetRuleType)
				if err != nil {
					return nil, fmt.Errorf("failed to parse rule type: %v", err)
				}

				convert.FrameworkToGraphString(rule.Path, fileSystemRule.SetPath)
				convert.FrameworkToGraphString(rule.FileOrFolderName, fileSystemRule.SetFileOrFolderName)
				convert.FrameworkToGraphBool(rule.Check32BitOn64System, fileSystemRule.SetCheck32BitOn64System)

				err = convert.FrameworkToGraphEnum(rule.LobAppRuleOperator, graphmodels.ParseWin32LobAppRuleOperator, fileSystemRule.SetOperator)
				if err != nil {
					return nil, fmt.Errorf("failed to parse file system rule operator: %v", err)
				}

				err = convert.FrameworkToGraphEnum(rule.FileSystemOperationType, graphmodels.ParseWin32LobAppFileSystemOperationType, fileSystemRule.SetOperationType)
				if err != nil {
					return nil, fmt.Errorf("failed to parse file system rule operation type: %v", err)
				}

				convert.FrameworkToGraphString(rule.ComparisonValue, fileSystemRule.SetComparisonValue)

				rules[i] = fileSystemRule

			case "powershell_script":
				powershellRule := graphmodels.NewWin32LobAppPowerShellScriptRule()

				// Set common rule properties
				err := convert.FrameworkToGraphEnum(rule.RuleType, graphmodels.ParseWin32LobAppRuleType, powershellRule.SetRuleType)
				if err != nil {
					return nil, fmt.Errorf("failed to parse rule type: %v", err)
				}

				// Only set DisplayName for requirement rules, not detection rules
				// The API doesn't allow DisplayName for detection rules
				if rule.RuleType.ValueString() == "requirement" && !rule.DisplayName.IsNull() {
					convert.FrameworkToGraphString(rule.DisplayName, powershellRule.SetDisplayName)
				}

				convert.FrameworkToGraphBool(rule.EnforceSignatureCheck, powershellRule.SetEnforceSignatureCheck)
				convert.FrameworkToGraphBool(rule.RunAs32Bit, powershellRule.SetRunAs32Bit)

				// Convert script content to base64 before setting the
				// converted content as a string
				if !rule.ScriptContent.IsNull() && !rule.ScriptContent.IsUnknown() {
					scriptContent := rule.ScriptContent.ValueString()

					base64Content, err := helpers.StringToBase64(scriptContent)
					if err == nil {

						convert.FrameworkToGraphString(types.StringValue(base64Content), powershellRule.SetScriptContent)
					} else {

						convert.FrameworkToGraphString(rule.ScriptContent, powershellRule.SetScriptContent)
					}
				}

				// Only set RunAsAccount if it's provided AND this is a requirement rule
				// The API doesn't allow RunAsAccount for detection rules
				if rule.RuleType.ValueString() == "requirement" && !rule.RunAsAccount.IsNull() && !rule.RunAsAccount.IsUnknown() {
					err = convert.FrameworkToGraphEnum(rule.RunAsAccount, graphmodels.ParseRunAsAccountType, powershellRule.SetRunAsAccount)
					if err != nil {
						return nil, fmt.Errorf("failed to parse PowerShell script rule run as account: %v", err)
					}
				}

				err = convert.FrameworkToGraphEnum(rule.LobAppRuleOperator, graphmodels.ParseWin32LobAppRuleOperator, powershellRule.SetOperator)
				if err != nil {
					return nil, fmt.Errorf("failed to parse PowerShell script rule operator: %v", err)
				}

				err = convert.FrameworkToGraphEnum(rule.PowerShellScriptRuleOperationType, graphmodels.ParseWin32LobAppPowerShellScriptRuleOperationType, powershellRule.SetOperationType)
				if err != nil {
					return nil, fmt.Errorf("failed to parse PowerShell script rule operation type: %v", err)
				}

				convert.FrameworkToGraphString(rule.ComparisonValue, powershellRule.SetComparisonValue)

				rules[i] = powershellRule

			default:
				return nil, fmt.Errorf("unsupported rule sub-type: %s", rule.RuleSubType.ValueString())
			}
		}
		requestBody.SetRules(rules)
	}

	if installExperience := data.InstallExperience; installExperience != (Win32LobAppInstallExperienceResourceModel{}) {
		installExp := graphmodels.NewWin32LobAppInstallExperience()

		err := convert.FrameworkToGraphEnum(installExperience.RunAsAccount, graphmodels.ParseRunAsAccountType, installExp.SetRunAsAccount)
		if err != nil {
			return nil, fmt.Errorf("failed to parse RunAsAccountType: %v", err)
		}

		err = convert.FrameworkToGraphEnum(installExperience.DeviceRestartBehavior, graphmodels.ParseWin32LobAppRestartBehavior, installExp.SetDeviceRestartBehavior)
		if err != nil {
			return nil, fmt.Errorf("failed to parse DeviceRestartBehavior: %v", err)
		}

		convert.FrameworkToGraphInt32(installExperience.MaxRunTimeInMinutes, installExp.SetMaxRunTimeInMinutes)

		requestBody.SetInstallExperience(installExp)
	}

	if len(data.ReturnCodes) > 0 {
		returnCodes := make([]graphmodels.Win32LobAppReturnCodeable, len(data.ReturnCodes))
		for i, code := range data.ReturnCodes {
			returnCode := graphmodels.NewWin32LobAppReturnCode()

			convert.FrameworkToGraphInt32(code.ReturnCode, returnCode.SetReturnCode)

			err := convert.FrameworkToGraphEnum(code.Type, graphmodels.ParseWin32LobAppReturnCodeType, returnCode.SetTypeEscaped)
			if err != nil {
				return nil, fmt.Errorf("failed to parse return code type: %v", err)
			}

			returnCodes[i] = returnCode
		}
		requestBody.SetReturnCodes(returnCodes)
	}

	if msiInfo := data.MsiInformation; msiInfo != (&Win32LobAppMsiInformationResourceModel{}) {
		msiInformation := graphmodels.NewWin32LobAppMsiInformation()

		convert.FrameworkToGraphString(msiInfo.ProductCode, msiInformation.SetProductCode)
		convert.FrameworkToGraphString(msiInfo.ProductVersion, msiInformation.SetProductVersion)
		convert.FrameworkToGraphString(msiInfo.UpgradeCode, msiInformation.SetUpgradeCode)
		convert.FrameworkToGraphBool(msiInfo.RequiresReboot, msiInformation.SetRequiresReboot)
		convert.FrameworkToGraphString(msiInfo.ProductName, msiInformation.SetProductName)
		convert.FrameworkToGraphString(msiInfo.Publisher, msiInformation.SetPublisher)

		err := convert.FrameworkToGraphEnum(msiInfo.PackageType, graphmodels.ParseWin32LobAppMsiPackageType, msiInformation.SetPackageType)
		if err != nil {
			return nil, fmt.Errorf("failed to parse MSI package type: %v", err)
		}
		requestBody.SetMsiInformation(msiInformation)
	}

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]any{
			"error": err.Error(),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing %s resource", ResourceName))

	return requestBody, nil
}
