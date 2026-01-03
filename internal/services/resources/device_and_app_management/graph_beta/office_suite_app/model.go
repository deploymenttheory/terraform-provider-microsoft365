// REF: https://learn.microsoft.com/en-us/graph/api/resources/intune-shared-mobileapp?view=graph-rest-beta
// REF: https://learn.microsoft.com/en-us/graph/api/resources/intune-apps-officesuiteapp?view=graph-rest-beta

package graphBetaOfficeSuiteApp

import (
	sharedmodels "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/shared_models/graph_beta/device_and_app_management"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// OfficeSuiteAppResourceModel represents the Terraform resource model for an Office Suite App
type OfficeSuiteAppResourceModel struct {
	ID                    types.String                             `tfsdk:"id"`
	DisplayName           types.String                             `tfsdk:"display_name"`
	Description           types.String                             `tfsdk:"description"`
	Publisher             types.String                             `tfsdk:"publisher"`
	Categories            types.Set                                `tfsdk:"categories"`
	AppIcon               *sharedmodels.MobileAppIconResourceModel `tfsdk:"app_icon"`
	CreatedDateTime       types.String                             `tfsdk:"created_date_time"`
	LastModifiedDateTime  types.String                             `tfsdk:"last_modified_date_time"`
	IsFeatured            types.Bool                               `tfsdk:"is_featured"`
	PrivacyInformationUrl types.String                             `tfsdk:"privacy_information_url"`
	InformationUrl        types.String                             `tfsdk:"information_url"`
	Owner                 types.String                             `tfsdk:"owner"`
	Developer             types.String                             `tfsdk:"developer"`
	Notes                 types.String                             `tfsdk:"notes"`
	UploadState           types.Int32                              `tfsdk:"upload_state"`
	PublishingState       types.String                             `tfsdk:"publishing_state"`
	IsAssigned            types.Bool                               `tfsdk:"is_assigned"`
	RoleScopeTagIds       types.Set                                `tfsdk:"role_scope_tag_ids"`
	DependentAppCount     types.Int32                              `tfsdk:"dependent_app_count"`
	SupersedingAppCount   types.Int32                              `tfsdk:"superseding_app_count"`
	SupersededAppCount    types.Int32                              `tfsdk:"superseded_app_count"`

	// Office Suite App configuration blocks (mutually exclusive)
	ConfigurationDesigner *OfficeSuiteAppConfigurationDesignerModel `tfsdk:"configuration_designer"`
	XMLConfiguration      *OfficeSuiteAppXMLConfigurationModel      `tfsdk:"xml_configuration"`

	Timeouts timeouts.Value `tfsdk:"timeouts"`
}

// OfficeSuiteAppExcludedAppsModel represents the excluded apps configuration
type OfficeSuiteAppExcludedAppsModel struct {
	Access             types.Bool `tfsdk:"access"`
	Bing               types.Bool `tfsdk:"bing"`
	Excel              types.Bool `tfsdk:"excel"`
	Groove             types.Bool `tfsdk:"groove"`
	InfoPath           types.Bool `tfsdk:"info_path"`
	Lync               types.Bool `tfsdk:"lync"`
	OneDrive           types.Bool `tfsdk:"one_drive"`
	OneNote            types.Bool `tfsdk:"one_note"`
	Outlook            types.Bool `tfsdk:"outlook"`
	PowerPoint         types.Bool `tfsdk:"power_point"`
	Publisher          types.Bool `tfsdk:"publisher"`
	SharePointDesigner types.Bool `tfsdk:"share_point_designer"`
	Teams              types.Bool `tfsdk:"teams"`
	Visio              types.Bool `tfsdk:"visio"`
	Word               types.Bool `tfsdk:"word"`
}

// OfficeSuiteAppConfigurationDesignerModel represents the configuration designer block
type OfficeSuiteAppConfigurationDesignerModel struct {
	AutoAcceptEula                       types.Bool                       `tfsdk:"auto_accept_eula"`
	ExcludedApps                         *OfficeSuiteAppExcludedAppsModel `tfsdk:"excluded_apps"`
	LocalesToInstall                     types.Set                        `tfsdk:"locales_to_install"`
	OfficePlatformArchitecture           types.String                     `tfsdk:"office_platform_architecture"`
	OfficeSuiteAppDefaultFileFormat      types.String                     `tfsdk:"office_suite_app_default_file_format"`
	ProductIds                           types.Set                        `tfsdk:"product_ids"`
	ShouldUninstallOlderVersionsOfOffice types.Bool                       `tfsdk:"should_uninstall_older_versions_of_office"`
	TargetVersion                        types.String                     `tfsdk:"target_version"`
	UpdateChannel                        types.String                     `tfsdk:"update_channel"`
	UpdateVersion                        types.String                     `tfsdk:"update_version"`
	UseSharedComputerActivation          types.Bool                       `tfsdk:"use_shared_computer_activation"`
}

// OfficeSuiteAppXMLConfigurationModel represents the XML configuration block
type OfficeSuiteAppXMLConfigurationModel struct {
	OfficeConfigurationXml types.String `tfsdk:"office_configuration_xml"`
}
