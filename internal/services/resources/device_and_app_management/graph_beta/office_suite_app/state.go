package graphBetaOfficeSuiteApp

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	sharedstater "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/state/graph_beta/device_and_app_management"
	sharedmodels "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/shared_models/graph_beta/device_and_app_management"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// MapRemoteResourceStateToTerraform maps the remote OfficeSuiteApp resource to the Terraform state
func MapRemoteResourceStateToTerraform(ctx context.Context, data *OfficeSuiteAppResourceModel, remoteResource graphmodels.OfficeSuiteAppable) {
	if remoteResource == nil {
		tflog.Debug(ctx, "Remote resource is nil")
		return
	}

	tflog.Debug(ctx, "Starting to map remote state to Terraform state", map[string]interface{}{
		"resourceId": convert.GraphToFrameworkString(remoteResource.GetId()).ValueString(),
	})

	// Base mobile app fields
	data.ID = convert.GraphToFrameworkString(remoteResource.GetId())
	data.DisplayName = convert.GraphToFrameworkString(remoteResource.GetDisplayName())
	data.Description = convert.GraphToFrameworkString(remoteResource.GetDescription())
	data.Publisher = convert.GraphToFrameworkString(remoteResource.GetPublisher())
	data.IsFeatured = convert.GraphToFrameworkBool(remoteResource.GetIsFeatured())
	data.PrivacyInformationUrl = convert.GraphToFrameworkString(remoteResource.GetPrivacyInformationUrl())
	data.InformationUrl = convert.GraphToFrameworkString(remoteResource.GetInformationUrl())
	data.Owner = convert.GraphToFrameworkString(remoteResource.GetOwner())
	data.Developer = convert.GraphToFrameworkString(remoteResource.GetDeveloper())
	data.Notes = convert.GraphToFrameworkString(remoteResource.GetNotes())
	data.CreatedDateTime = convert.GraphToFrameworkTime(remoteResource.GetCreatedDateTime())
	data.LastModifiedDateTime = convert.GraphToFrameworkTime(remoteResource.GetLastModifiedDateTime())
	data.UploadState = convert.GraphToFrameworkInt32(remoteResource.GetUploadState())
	data.PublishingState = convert.GraphToFrameworkEnum(remoteResource.GetPublishingState())
	data.IsAssigned = convert.GraphToFrameworkBool(remoteResource.GetIsAssigned())
	data.DependentAppCount = convert.GraphToFrameworkInt32(remoteResource.GetDependentAppCount())
	data.SupersedingAppCount = convert.GraphToFrameworkInt32(remoteResource.GetSupersedingAppCount())
	data.SupersededAppCount = convert.GraphToFrameworkInt32(remoteResource.GetSupersededAppCount())

	// Handle app icon - preserve original configuration values
	if data.AppIcon != nil {
		tflog.Debug(ctx, "Preserving original app_icon values from configuration")
	} else if largeIcon := remoteResource.GetLargeIcon(); largeIcon != nil {
		data.AppIcon = &sharedmodels.MobileAppIconResourceModel{
			IconFilePathSource: types.StringNull(),
			IconURLSource:      types.StringNull(),
		}
	} else {
		data.AppIcon = nil
	}

	data.RoleScopeTagIds = convert.GraphToFrameworkStringSet(ctx, remoteResource.GetRoleScopeTagIds())
	data.Categories = sharedstater.MapMobileAppCategoriesStateToTerraform(ctx, remoteResource.GetCategories())

	// Handle Office Suite App configuration blocks (mutually exclusive)
	// Check if user originally configured using Configuration Designer or XML Configuration
	if data.ConfigurationDesigner != nil {
		// Configuration Designer block was used - populate it
		if data.ConfigurationDesigner == nil {
			data.ConfigurationDesigner = &OfficeSuiteAppConfigurationDesignerModel{}
		}

		data.ConfigurationDesigner.AutoAcceptEula = convert.GraphToFrameworkBool(remoteResource.GetAutoAcceptEula())

		// Handle excluded apps
		if excludedApps := remoteResource.GetExcludedApps(); excludedApps != nil {
			data.ConfigurationDesigner.ExcludedApps = &OfficeSuiteAppExcludedAppsModel{
				Access:             convert.GraphToFrameworkBool(excludedApps.GetAccess()),
				Bing:               convert.GraphToFrameworkBool(excludedApps.GetBing()),
				Excel:              convert.GraphToFrameworkBool(excludedApps.GetExcel()),
				Groove:             convert.GraphToFrameworkBool(excludedApps.GetGroove()),
				InfoPath:           convert.GraphToFrameworkBool(excludedApps.GetInfoPath()),
				Lync:               convert.GraphToFrameworkBool(excludedApps.GetLync()),
				OneDrive:           convert.GraphToFrameworkBool(excludedApps.GetOneDrive()),
				OneNote:            convert.GraphToFrameworkBool(excludedApps.GetOneNote()),
				Outlook:            convert.GraphToFrameworkBool(excludedApps.GetOutlook()),
				PowerPoint:         convert.GraphToFrameworkBool(excludedApps.GetPowerPoint()),
				Publisher:          convert.GraphToFrameworkBool(excludedApps.GetPublisher()),
				SharePointDesigner: convert.GraphToFrameworkBool(excludedApps.GetSharePointDesigner()),
				Teams:              convert.GraphToFrameworkBool(excludedApps.GetTeams()),
				Visio:              convert.GraphToFrameworkBool(excludedApps.GetVisio()),
				Word:               convert.GraphToFrameworkBool(excludedApps.GetWord()),
			}
		} else {
			data.ConfigurationDesigner.ExcludedApps = nil
		}

		data.ConfigurationDesigner.LocalesToInstall = convert.GraphToFrameworkStringSet(ctx, remoteResource.GetLocalesToInstall())
		data.ConfigurationDesigner.OfficePlatformArchitecture = convert.GraphToFrameworkEnum(remoteResource.GetOfficePlatformArchitecture())
		data.ConfigurationDesigner.OfficeSuiteAppDefaultFileFormat = convert.GraphToFrameworkEnum(remoteResource.GetOfficeSuiteAppDefaultFileFormat())

		// Handle product IDs - convert from OfficeProductId enum to strings
		if productIds := remoteResource.GetProductIds(); productIds != nil {
			productIdStrings := make([]string, len(productIds))
			for i, productId := range productIds {
				productIdStrings[i] = productId.String()
			}
			data.ConfigurationDesigner.ProductIds = convert.GraphToFrameworkStringSet(ctx, productIdStrings)
		} else {
			data.ConfigurationDesigner.ProductIds = types.SetNull(types.StringType)
		}

		data.ConfigurationDesigner.ShouldUninstallOlderVersionsOfOffice = convert.GraphToFrameworkBool(remoteResource.GetShouldUninstallOlderVersionsOfOffice())
		data.ConfigurationDesigner.TargetVersion = convert.GraphToFrameworkString(remoteResource.GetTargetVersion())
		data.ConfigurationDesigner.UpdateChannel = convert.GraphToFrameworkEnum(remoteResource.GetUpdateChannel())
		data.ConfigurationDesigner.UpdateVersion = convert.GraphToFrameworkString(remoteResource.GetUpdateVersion())
		data.ConfigurationDesigner.UseSharedComputerActivation = convert.GraphToFrameworkBool(remoteResource.GetUseSharedComputerActivation())

		// Ensure XML configuration is null when using Configuration Designer
		data.XMLConfiguration = nil

	} else if data.XMLConfiguration != nil {
		// XML Configuration block was used - populate it
		if data.XMLConfiguration == nil {
			data.XMLConfiguration = &OfficeSuiteAppXMLConfigurationModel{}
		}

		data.XMLConfiguration.OfficeConfigurationXml = convert.GraphToFrameworkBytes(remoteResource.GetOfficeConfigurationXml())

		// Ensure Configuration Designer is null when using XML Configuration
		data.ConfigurationDesigner = nil

	} else {
		// Neither block was configured originally - determine which one to use based on remote resource
		if remoteResource.GetOfficeConfigurationXml() != nil && len(remoteResource.GetOfficeConfigurationXml()) > 0 {
			// Has XML configuration, use XML block
			data.XMLConfiguration = &OfficeSuiteAppXMLConfigurationModel{
				OfficeConfigurationXml: convert.GraphToFrameworkBytes(remoteResource.GetOfficeConfigurationXml()),
			}
			data.ConfigurationDesigner = nil
		} else {
			// No XML configuration or has individual settings, use Configuration Designer block
			data.ConfigurationDesigner = &OfficeSuiteAppConfigurationDesignerModel{
				AutoAcceptEula:                          convert.GraphToFrameworkBool(remoteResource.GetAutoAcceptEula()),
				LocalesToInstall:                       convert.GraphToFrameworkStringSet(ctx, remoteResource.GetLocalesToInstall()),
				OfficePlatformArchitecture:             convert.GraphToFrameworkEnum(remoteResource.GetOfficePlatformArchitecture()),
				OfficeSuiteAppDefaultFileFormat:        convert.GraphToFrameworkEnum(remoteResource.GetOfficeSuiteAppDefaultFileFormat()),
				ShouldUninstallOlderVersionsOfOffice:   convert.GraphToFrameworkBool(remoteResource.GetShouldUninstallOlderVersionsOfOffice()),
				TargetVersion:                          convert.GraphToFrameworkString(remoteResource.GetTargetVersion()),
				UpdateChannel:                          convert.GraphToFrameworkEnum(remoteResource.GetUpdateChannel()),
				UpdateVersion:                          convert.GraphToFrameworkString(remoteResource.GetUpdateVersion()),
				UseSharedComputerActivation:            convert.GraphToFrameworkBool(remoteResource.GetUseSharedComputerActivation()),
			}

			// Handle excluded apps
			if excludedApps := remoteResource.GetExcludedApps(); excludedApps != nil {
				data.ConfigurationDesigner.ExcludedApps = &OfficeSuiteAppExcludedAppsModel{
					Access:             convert.GraphToFrameworkBool(excludedApps.GetAccess()),
					Bing:               convert.GraphToFrameworkBool(excludedApps.GetBing()),
					Excel:              convert.GraphToFrameworkBool(excludedApps.GetExcel()),
					Groove:             convert.GraphToFrameworkBool(excludedApps.GetGroove()),
					InfoPath:           convert.GraphToFrameworkBool(excludedApps.GetInfoPath()),
					Lync:               convert.GraphToFrameworkBool(excludedApps.GetLync()),
					OneDrive:           convert.GraphToFrameworkBool(excludedApps.GetOneDrive()),
					OneNote:            convert.GraphToFrameworkBool(excludedApps.GetOneNote()),
					Outlook:            convert.GraphToFrameworkBool(excludedApps.GetOutlook()),
					PowerPoint:         convert.GraphToFrameworkBool(excludedApps.GetPowerPoint()),
					Publisher:          convert.GraphToFrameworkBool(excludedApps.GetPublisher()),
					SharePointDesigner: convert.GraphToFrameworkBool(excludedApps.GetSharePointDesigner()),
					Teams:              convert.GraphToFrameworkBool(excludedApps.GetTeams()),
					Visio:              convert.GraphToFrameworkBool(excludedApps.GetVisio()),
					Word:               convert.GraphToFrameworkBool(excludedApps.GetWord()),
				}
			}

			// Handle product IDs - convert from OfficeProductId enum to strings
			if productIds := remoteResource.GetProductIds(); productIds != nil {
				productIdStrings := make([]string, len(productIds))
				for i, productId := range productIds {
					productIdStrings[i] = productId.String()
				}
				data.ConfigurationDesigner.ProductIds = convert.GraphToFrameworkStringSet(ctx, productIdStrings)
			} else {
				data.ConfigurationDesigner.ProductIds = types.SetNull(types.StringType)
			}

			data.XMLConfiguration = nil
		}
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished stating resource %s with id %s", ResourceName, data.ID.ValueString()))
}