package graphBetaOfficeSuiteApp

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	sharedConstructors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors/graph_beta/device_and_app_management"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	helpers "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud/graph_beta/device_and_app_management"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// constructResource constructs an OfficeSuiteApp resource using the provided model data.
func constructResource(ctx context.Context, data *OfficeSuiteAppResourceModel, isUpdate bool) (graphmodels.OfficeSuiteAppable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s resource", ResourceName))

	requestBody := graphmodels.NewOfficeSuiteApp()

	convert.FrameworkToGraphString(data.DisplayName, requestBody.SetDisplayName)
	convert.FrameworkToGraphString(data.Description, requestBody.SetDescription)
	// Use hardcoded values from Microsoft Graph API documentation
	publisher := "Microsoft"
	developer := "Microsoft"
	owner := "Microsoft"
	requestBody.SetPublisher(&publisher)
	requestBody.SetDeveloper(&developer)
	requestBody.SetOwner(&owner)
	convert.FrameworkToGraphString(data.Notes, requestBody.SetNotes)
	convert.FrameworkToGraphBool(data.IsFeatured, requestBody.SetIsFeatured)
	// Use hardcoded privacy information URL from Microsoft Graph API documentation
	privacyUrl := "https://privacy.microsoft.com/privacystatement"
	requestBody.SetPrivacyInformationUrl(&privacyUrl)
	convert.FrameworkToGraphString(data.InformationUrl, requestBody.SetInformationUrl)

	// Handle role scope tag IDs
	if err := convert.FrameworkToGraphStringSet(ctx, data.RoleScopeTagIds, requestBody.SetRoleScopeTagIds); err != nil {
		return nil, fmt.Errorf("failed to set role scope tag IDs: %v", err)
	}

	// Handle app icon (either from file path or web source)
	if data.AppIcon != nil {
		largeIcon, tempFiles, err := sharedConstructors.ConstructMobileAppIcon(ctx, data.AppIcon)
		if err != nil {
			return nil, err
		}

		defer func() {
			for _, tempFile := range tempFiles {
				helpers.CleanupTempFile(ctx, tempFile)
			}
		}()

		requestBody.SetLargeIcon(largeIcon)
	}

	// Handle Office Suite App configuration blocks (mutually exclusive)
	if data.ConfigurationDesigner != nil {
		// Configuration Designer block - individual settings
		convert.FrameworkToGraphBool(data.ConfigurationDesigner.AutoAcceptEula, requestBody.SetAutoAcceptEula)

		// Handle excluded apps
		if data.ConfigurationDesigner.ExcludedApps != nil {
			excludedApps := graphmodels.NewExcludedApps()

			convert.FrameworkToGraphBool(data.ConfigurationDesigner.ExcludedApps.Access, excludedApps.SetAccess)
			convert.FrameworkToGraphBool(data.ConfigurationDesigner.ExcludedApps.Bing, excludedApps.SetBing)
			convert.FrameworkToGraphBool(data.ConfigurationDesigner.ExcludedApps.Excel, excludedApps.SetExcel)
			convert.FrameworkToGraphBool(data.ConfigurationDesigner.ExcludedApps.Groove, excludedApps.SetGroove)
			convert.FrameworkToGraphBool(data.ConfigurationDesigner.ExcludedApps.InfoPath, excludedApps.SetInfoPath)
			convert.FrameworkToGraphBool(data.ConfigurationDesigner.ExcludedApps.Lync, excludedApps.SetLync)
			convert.FrameworkToGraphBool(data.ConfigurationDesigner.ExcludedApps.OneDrive, excludedApps.SetOneDrive)
			convert.FrameworkToGraphBool(data.ConfigurationDesigner.ExcludedApps.OneNote, excludedApps.SetOneNote)
			convert.FrameworkToGraphBool(data.ConfigurationDesigner.ExcludedApps.Outlook, excludedApps.SetOutlook)
			convert.FrameworkToGraphBool(data.ConfigurationDesigner.ExcludedApps.PowerPoint, excludedApps.SetPowerPoint)
			convert.FrameworkToGraphBool(data.ConfigurationDesigner.ExcludedApps.Publisher, excludedApps.SetPublisher)
			convert.FrameworkToGraphBool(data.ConfigurationDesigner.ExcludedApps.SharePointDesigner, excludedApps.SetSharePointDesigner)
			convert.FrameworkToGraphBool(data.ConfigurationDesigner.ExcludedApps.Teams, excludedApps.SetTeams)
			convert.FrameworkToGraphBool(data.ConfigurationDesigner.ExcludedApps.Visio, excludedApps.SetVisio)
			convert.FrameworkToGraphBool(data.ConfigurationDesigner.ExcludedApps.Word, excludedApps.SetWord)

			requestBody.SetExcludedApps(excludedApps)
		}

		// Handle locales to install
		if err := convert.FrameworkToGraphStringSet(ctx, data.ConfigurationDesigner.LocalesToInstall, requestBody.SetLocalesToInstall); err != nil {
			return nil, fmt.Errorf("failed to set locales to install: %v", err)
		}

		// OfficePlatformArchitecture - use enum helper function
		if err := convert.FrameworkToGraphEnum(data.ConfigurationDesigner.OfficePlatformArchitecture,
			graphmodels.ParseWindowsArchitecture, requestBody.SetOfficePlatformArchitecture); err != nil {
			return nil, fmt.Errorf("failed to set office platform architecture: %v", err)
		}

		// OfficeSuiteAppDefaultFileFormat - use enum helper function
		if err := convert.FrameworkToGraphEnum(data.ConfigurationDesigner.OfficeSuiteAppDefaultFileFormat,
			graphmodels.ParseOfficeSuiteDefaultFileFormatType, requestBody.SetOfficeSuiteAppDefaultFileFormat); err != nil {
			return nil, fmt.Errorf("failed to set office suite app default file format: %v", err)
		}

		// Handle product IDs - convert strings to OfficeProductId enums
		if err := convert.FrameworkToGraphObjectsFromStringSet(ctx, data.ConfigurationDesigner.ProductIds,
			func(ctx context.Context, stringValues []string) []graphmodels.OfficeProductId {
				var productIds []graphmodels.OfficeProductId
				for _, productIdStr := range stringValues {
					switch productIdStr {
					case "o365ProPlusRetail":
						productIds = append(productIds, graphmodels.O365PROPLUSRETAIL_OFFICEPRODUCTID)
					case "projectProRetail":
						productIds = append(productIds, graphmodels.PROJECTPRORETAIL_OFFICEPRODUCTID)
					case "visioProRetail":
						productIds = append(productIds, graphmodels.VISIOPRORETAIL_OFFICEPRODUCTID)
					case "o365BusinessRetail":
						productIds = append(productIds, graphmodels.O365BUSINESSRETAIL_OFFICEPRODUCTID)
					}
				}
				return productIds
			}, requestBody.SetProductIds); err != nil {
			return nil, fmt.Errorf("failed to set product IDs: %v", err)
		}

		convert.FrameworkToGraphBool(data.ConfigurationDesigner.ShouldUninstallOlderVersionsOfOffice, requestBody.SetShouldUninstallOlderVersionsOfOffice)
		convert.FrameworkToGraphString(data.ConfigurationDesigner.TargetVersion, requestBody.SetTargetVersion)

		// UpdateChannel - use enum helper function
		if err := convert.FrameworkToGraphEnum(data.ConfigurationDesigner.UpdateChannel,
			graphmodels.ParseOfficeUpdateChannel, requestBody.SetUpdateChannel); err != nil {
			return nil, fmt.Errorf("failed to set update channel: %v", err)
		}

		convert.FrameworkToGraphString(data.ConfigurationDesigner.UpdateVersion, requestBody.SetUpdateVersion)
		convert.FrameworkToGraphBool(data.ConfigurationDesigner.UseSharedComputerActivation, requestBody.SetUseSharedComputerActivation)

	} else if data.XMLConfiguration != nil {
		// XML Configuration block - use XML configuration
		convert.FrameworkToGraphBytes(data.XMLConfiguration.OfficeConfigurationXml, requestBody.SetOfficeConfigurationXml)
	}

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]interface{}{
			"error": err.Error(),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing %s resource", ResourceName))

	return requestBody, nil
}
