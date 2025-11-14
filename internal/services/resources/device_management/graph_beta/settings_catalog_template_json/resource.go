package graphBetaDeviceManagementTemplateJson

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	planmodifiers "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/plan_modifiers"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
	commonschemagraphbeta "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema/graph_beta/device_management"
	customValidator "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/validate/attribute"
	sharedValidators "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/validate/graph_beta/device_management"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

const (
	ResourceName  = "graph_beta_device_management_settings_catalog_template_json"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180
)

var (
	// Basic resource interface (CRUD operations)
	_ resource.Resource = &DeviceManagementTemplateJsonResource{}

	// Allows the resource to be configured with the provider client
	_ resource.ResourceWithConfigure = &DeviceManagementTemplateJsonResource{}

	// Enables import functionality
	_ resource.ResourceWithImportState = &DeviceManagementTemplateJsonResource{}

	// Enables plan modification/diff suppression
	_ resource.ResourceWithModifyPlan = &DeviceManagementTemplateJsonResource{}
)

func NewDeviceManagementTemplateJsonResource() resource.Resource {
	return &DeviceManagementTemplateJsonResource{
		ReadPermissions: []string{
			"DeviceManagementConfiguration.Read.All",
		},
		WritePermissions: []string{
			"DeviceManagementConfiguration.ReadWrite.All",
		},
		ResourcePath: "/deviceManagement/configurationPolicies",
	}
}

type DeviceManagementTemplateJsonResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ProviderTypeName string
	TypeName         string
	ReadPermissions  []string
	WritePermissions []string
	ResourcePath     string
}

// Metadata returns the resource type name.
func (r *DeviceManagementTemplateJsonResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	r.ProviderTypeName = req.ProviderTypeName
	r.TypeName = ResourceName
	resp.TypeName = r.FullTypeName()
}

// FullTypeName returns the full type name of the resource for logging purposes.
func (r *DeviceManagementTemplateJsonResource) FullTypeName() string {
	return r.ProviderTypeName + "_" + r.TypeName
}

// Configure sets the client for the resource.
func (r *DeviceManagementTemplateJsonResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, constants.PROVIDER_NAME+"_"+ResourceName)
}

// ImportState imports the resource state.
func (r *DeviceManagementTemplateJsonResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Function to create the full device management configuration policy schema
func (r *DeviceManagementTemplateJsonResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Settings Catalog policy template in Microsoft Intune for `Windows`, `macOS`, `Linux`, `iOS/iPadOS` and `Android`.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					planmodifiers.UseStateForUnknownString(),
				},
				MarkdownDescription: "The unique identifier for this settings catalog policy template",
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Settings Catalog Policy template name",
			},
			"description": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				PlanModifiers:       []planmodifier.String{planmodifiers.DefaultValueString("")},
				MarkdownDescription: "Settings Catalog Policy template description",
			},
			"settings_catalog_template_type": schema.StringAttribute{
				Required: true,
				MarkdownDescription: "Defines the intune settings catalog template type to be deployed using the settings catalog.\n\n" +
					"This value will automatically set the correct `platform` , `templateID` , `creationSource` and `technologies` values for the settings catalog policy." +
					"This value must correctly correlate to the settings defined in the `settings` field." +
					"The available options include templates for various platforms and configurations, such as macOS, Windows, and Linux. Options available are:\n\n" +
					"`Linux settings catalog templates`\n\n" +
					"`linux_anti_virus_microsoft_defender_antivirus`: This template allows you to configure Microsoft Defender for Endpoint and deploy Antivirus settings to Linux devices.\n\n" +
					"`linux_anti_virus_microsoft_defender_antivirus_exclusions`: This template allows you to manage settings for Microsoft Defender Antivirus that define Antivirus exclusions for paths, extensions and processes. Antivirus exclusion are also managed by Microsoft Defender Antivirus policy, which includes identical settings for exclusions. Settings from both templates (Antivirus and Antivirus exclusions) are subject to policy merge, and create a super set of exclusions for applicable devices and users.\n\n" +
					"`linux_endpoint_detection_and_response`: Endpoint detection and response settings for Linux devices.\n\n" +
					"`macOS settings catalog templates`\n\n" +
					"`macOS_anti_virus_microsoft_defender_antivirus`: Microsoft Defender Antivirus is the next-generation protection component of Microsoft Defender for Endpoint on Mac. Next-generation protection brings together machine learning, big-data analysis, in-depth threat resistance research, and cloud infrastructure to protect devices in your enterprise organization.\n\n" +
					"`macOS_anti_virus_microsoft_defender_antivirus_exclusions`: This template allows you to manage settings for Microsoft Defender Antivirus that define Antivirus exclusions for paths, extensions and processes. Antivirus exclusion are also managed by Microsoft Defender Antivirus policy, which includes identical settings for exclusions. Settings from both templates (Antivirus and Antivirus exclusions) are subject to policy merge, and create a super set of exclusions for applicable devices and users.\n\n" +
					"`macOS_endpoint_detection_and_response`: Endpoint detection and response settings for macOS devices.\n\n" +
					"`Security Baselines`\n\n" +
					"`security_baseline_for_windows_10_and_later_version_24H2`: The Security Baseline for Windows 10 and later represents the recommendations for configuring Windows for security conscious customers using the Microsoft full security stack. This baseline includes relevant MDM settings consistent with the security recommendations outlined in the group policy Windows security baseline. Use this baseline to tailor and adjust Microsoft-recommended policy settings within an MDM environment.\n\n" +
					"`security_baseline_for_microsoft_defender_for_endpoint_version_24H1`: he Microsoft Defender for Endpoint Security baseline for Windows 10 and newer represents the security best practices for the Microsoft security stack on devices managed by Intune (MDM). Use the baseline to tailor and adjust Microsoft-recommended policy settings.\n\n" +
					"`security_baseline_for_microsoft_edge_version_128`: The Security Baseline for Microsoft Edge represents the recommendations for configuring Microsoft Edge for security conscious customers using the Microsoft full security stack. This baseline aligns with the security recommendations for Edge security baseline for group policy. Use this baseline to configure and customize Microsoft-recommended policy settings.\n\n" +
					"`security_baseline_for_windows_365`: Windows 365 Security Baselines are a set of policy templates that you can deploy with Microsoft Intune to configure and enforce security settings for Windows 10, Windows 11, Microsoft Edge, and Microsoft Defender for Endpoint on your Cloud PCs. They are based on security best practices and real-world implementations, and they include versioning features to help you update your policies to the latest release. You can also customize the baselines to meet your specific business needs.\n\n" +
					"`security_baseline_for_microsoft_365_apps`: The Microsoft 365 Apps for enterprise security baseline provides a starting point for IT admins to evaluate and balance the security benefits with productivity needs of their users. This baseline aligns with the security recommendations for Microsoft 365 Apps for enterprise group policy security baseline. Use this baseline to configure and customize Microsoft-recommended policy settings.\n\n" +
					"`Windows settings catalog templates`\n\n" +
					"`windows_account_protection`: Account protection policies help protect user credentials by using technology such as Windows Hello for Business and Credential Guard.\n\n" +
					"`windows_anti_virus_defender_update_controls`: This template allows you to configure the gradual release rollout of Defender Updates to targeted device groups. Use a ringed approach to test, validate, and rollout updates to devices through release channels. Updates available are platform, engine, security intelligence updates. These policy types have pause, resume, manual rollback commands similar to Windows Update ring policies.\n\n" +
					"`windows_anti_virus_microsoft_defender_antivirus`: Windows Defender Antivirus is the next-generation protection component of Microsoft Defender for Endpoint. Next-generation protection brings together machine learning, big-data analysis, in-depth threat resistance research, and cloud infrastructure to protect devices in your enterprise organization.\n\n" +
					"`windows_anti_virus_microsoft_defender_antivirus_exclusions`: This template allows you to manage settings for Microsoft Defender Antivirus that define Antivirus exclusions for paths, extensions and processes. Antivirus exclusion are also managed by Microsoft Defender Antivirus policy, which includes identical settings for exclusions. Settings from both templates (Antivirus and Antivirus exclusions) are subject to policy merge, and create a super set of exclusions for applicable devices and users.\n\n" +
					"`windows_anti_virus_security_experience`: This template allows you to configure the Windows Security app is used by a number of Windows security features to provide notifications about the health and security of the machine. These include notifications about firewalls, antivirus products, Windows Defender SmartScreen, and others.\n\n" +
					"`windows_imported_administrative_templates`: This template allows you to configure imported custom and third-party/partner ADMX and ADML templates into the Intune admin center. Once imported, you can create a device configuration policy using these settings, and then assign the policy to your managed devices..\n\n" +
					"`windows_app_control_for_business`: Application control settings for Windows devices.\n\n" +
					"`windows_attack_surface_reduction_app_and_browser_isolation`:This template allows you to configure the Microsoft Defender Application Guard (Application Guard) to help prevent old and newly emerging attacks to help keep employees productive. Using MSFT's unique hardware isolation approach, their goal is to destroy the playbook that attackers use by making current attack methods obsolete.\n\n" +
					"`windows_attack_surface_reduction_attack_surface_reduction_rules`: This template allows you to configure the Attack surface reduction rules target behaviors that malware and malicious apps typically use to infect computers, including: Executable files and scripts used in Office apps or web mail that attempt to download or run files Obfuscated or otherwise suspicious scripts Behaviors that apps don't usually initiate during normal day-to-day work\n\n" +
					"`windows_attack_surface_reduction_app_device_control`:This template allows you to configure the securing removable media, and Microsoft Defender for Endpoint provides multiple monitoring and control features to help prevent threats in unauthorized peripherals from compromising your devices.\n\n" +
					"`windows_attack_surface_reduction_exploit_protection`: This template allows you to configure the protection against malware that uses exploits to infect devices and spread. Exploit protection consists of a number of mitigations that can be applied to either the operating system or individual apps.\n\n" +
					"`windows_disk_encryption_bitlocker`: This template allows you to configure the BitLocker Drive Encryption data protection features that integrates with the operating system and addresses the threats of data theft or exposure from lost, stolen, or inappropriately decommissioned computers.\n\n" +
					"`windows_disk_encryption_personal_data`: This template allows you to configure the Personal Data Encryption feature that encrypts select folders and its contents on deployed devices. Personal Data Encryption utilizes Windows Hello for Business to link data encryption keys with user credentials. This feature can minimize the number of credentials the user has to remember to gain access to content. Users will only be able to access their protected content once they've signed into Windows using Windows Hello for Business.\n\n" +
					"`windows_endpoint_detection_and_response`: Endpoint detection and response settings for Windows devices.\n\n" +
					"`windows_firewall_rules`: Firewall rules for Windows devices.\n\n" +
					"`windows_firewall_rules_config_manager`: Rules-based firewall configuration for Windows devices.\n\n" +
					"`windows_hyper-v_firewall_rules`: Hyper-V firewall rules for Windows devices.\n\n" +
					"`windows_local_admin_password_solution_(windows_LAPS)`: Windows Local Administrator Password Solution(Windows LAPS) is a Windows feature that automatically manages and backs up the password of a local administrator account on your Azure Active Directory - joined or Windows Server Active Directory - joined devices.\n\n" +
					"`windows_local_user_group_membership`: Local user group membership policies help to add, remove, or replace members of local groups on Windows devices..\n\n" +
					"`Windows Configuration Manager settings catalog templates`\n\n" +
					"`windows_(config_mgr)_anti_virus_microsoft_defender_antivirus`: Microsoft Defender Antivirus settings for Windows devices managed via Microsoft Configuration Manager.\n\n" +
					"`windows_(config_mgr)_anti_virus_windows_security_experience`: Security experience settings for Windows devices managed via Microsoft Configuration Manager.\n\n" +
					"`windows_(config_mgr)_attack_surface_reduction`: Attack surface reduction settings for Windows devices managed via Microsoft Configuration Manager.\n\n" +
					"`windows_(config_mgr)_endpoint_detection_and_response`: Endpoint detection and response settings for Windows devices managed via Microsoft Configuration Manager.\n\n" +
					"`windows_(config_mgr)_firewall`: Firewall settings for Windows devices managed via Microsoft Configuration Manager.\n\n" +
					"`windows_(config_mgr)_firewall_profile`: Profile-specific firewall configuration for Windows devices managed via Microsoft Configuration Manager.\n\n" +
					"`windows_(config_mgr)_firewall_rules`: Rules-based firewall configuration for Windows devices managed via Microsoft Configuration Manager.\n\n" +
					"`windows_(config_mgr)_attack_surface_reduction_app_and_browser_isolation`: This template allows you to configure the Microsoft Defender Application Guard (Application Guard) settings for devices managed via Microsoft Configuration Manager to help prevent old and newly emerging attacks through hardware-based isolation.\n\n" +
					"`windows_(config_mgr)_attack_surface_reduction_attack_surface_reduction_rules`: This template allows you to configure Attack Surface Reduction rules for devices managed via Microsoft Configuration Manager. These rules target behaviors commonly used by malware and malicious apps, including suspicious scripts and unusual app behaviors.\n\n" +
					"`windows_(config_mgr)_attack_surface_reduction_web_protection`: This template allows you to configure web protection settings for devices managed via Microsoft Configuration Manager, helping to protect your organization from web-based threats and malicious content.\n\n" +
					"`windows_(config_mgr)_attack_surface_reduction_exploit_protection`: This template allows you to configure exploit protection settings for devices managed via Microsoft Configuration Manager. These settings help protect against malware that uses exploits to infect devices and spread through your network.\n\n",
				// TODO - these template types currently use a legacy api endpoint. will implement in the future
				//"macOS_disk_encryption_filevault", uses templateManagement api endpoint
				//"macOS_firewall", uses templateManagement api endpoint
				//"windows_attack_surface_reduction_web_protection_(microsoft_edge_legacy)", uses templateManagement api endpoint
				//"windows_attack_surface_reduction_application_control", uses templateManagement api endpoint
				Validators: []validator.String{
					stringvalidator.OneOf(
						"linux_anti_virus_microsoft_defender_antivirus",
						"linux_anti_virus_microsoft_defender_antivirus_exclusions",
						"linux_endpoint_detection_and_response",
						"macOS_anti_virus_microsoft_defender_antivirus",
						"macOS_anti_virus_microsoft_defender_antivirus_exclusions",
						"macOS_endpoint_detection_and_response",
						"security_baseline_for_windows_10_and_later_version_24H2",
						"security_baseline_for_microsoft_defender_for_endpoint_version_24H1",
						"security_baseline_for_microsoft_edge_version_128",
						"security_baseline_for_windows_365",
						"security_baseline_for_microsoft_365_apps",
						"windows_account_protection",
						"windows_anti_virus_defender_update_controls",
						"windows_anti_virus_microsoft_defender_antivirus",
						"windows_anti_virus_microsoft_defender_antivirus_exclusions",
						"windows_anti_virus_security_experience",
						"windows_app_control_for_business",
						"windows_attack_surface_reduction_app_and_browser_isolation",
						"windows_attack_surface_reduction_attack_surface_reduction_rules",
						"windows_attack_surface_reduction_app_device_control",
						"windows_attack_surface_reduction_exploit_protection",
						"windows_disk_encryption_bitlocker",
						"windows_disk_encryption_personal_data",
						"windows_endpoint_detection_and_response",
						"windows_firewall",
						"windows_firewall_rules",
						"windows_hyper-v_firewall_rules",
						"windows_local_admin_password_solution_(windows_LAPS)",
						"windows_local_user_group_membership",
						"windows_(config_mgr)_attack_surface_reduction_app_and_browser_isolation",
						"windows_(config_mgr)_attack_surface_reduction_attack_surface_reduction_rules",
						"windows_(config_mgr)_attack_surface_reduction_web_protection",
						"windows_(config_mgr)_attack_surface_reduction_exploit_protection",
						"windows_(config_mgr)_anti_virus_microsoft_defender_antivirus",
						"windows_(config_mgr)_anti_virus_windows_security_experience",
						"windows_(config_mgr)_endpoint_detection_and_response",
						"windows_(config_mgr)_firewall",
						"windows_(config_mgr)_firewall_profile",
						"windows_(config_mgr)_firewall_rules",
						// TODO - these template types currently use a legacy api endpoint. will implement in the future
						//"macOS_disk_encryption_filevault", uses templateManagement api endpoint
						//"macOS_firewall", uses templateManagement api endpoint
						//"windows_attack_surface_reduction_web_protection_(microsoft_edge_legacy)", uses templateManagement api endpoint
						//"windows_attack_surface_reduction_application_control", uses templateManagement api endpoint
					),
				},
			},
			"settings": schema.StringAttribute{
				Required: true,
				MarkdownDescription: "Settings Catalog Policy template settings defined as a JSON string. Please provide a valid JSON-encoded settings structure. " +
					"This can either be extracted from an existing policy using the Intune gui `export JSON` functionality if supported, via a script such as this powershell script." +
					" [ExportSettingsCatalogTemplateConfigurationById](https://github.com/deploymenttheory/terraform-provider-microsoft365/blob/main/scripts/powershell/Export-IntuneSettingsCatalogTemplateConfigurationById.ps1) " +
					"or created from scratch. The JSON structure should match the graph schema of the settings catalog. Please look at the " +
					"terraform documentation for the settings catalog template for examples and how to correctly format the HCL.\n\n" +
					"A correctly formatted field in the HCL should begin and end like this:\n" +
					"```hcl\n" +
					"settings = jsonencode({\n" +
					"  \"settings\": [\n" +
					"    {\n" +
					"      \"id\": \"0\",\n" +
					"      \"settingInstance\": {\n" +
					"      }\n" +
					"    }\n" +
					"  ]\n" +
					"})\n" +
					"```\n\n" +
					"**Note:** Settings must always be provided as an array within the settings field, even when configuring a single setting." +
					"This is required because the Microsoft Graph SDK for Go always returns settings in an array format\n\n" +
					"**Note:** When configuring secret values (identified by @odata.type: \"#microsoft.graph.deviceManagementConfigurationSecretSettingValue\") " +
					"ensure the valueState is set to \"notEncrypted\". The value \"encryptedValueToken\" is reserved for server" +
					"responses and should not be used when creating or updating settings.\n\n" +
					"```hcl\n" +
					"settings = jsonencode({\n" +
					"  \"settings\": [\n" +
					"    {\n" +
					"      \"id\": \"0\",\n" +
					"      \"settingInstance\": {\n" +
					"        \"@odata.type\": \"#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance\",\n" +
					"        \"settingDefinitionId\": \"com.apple.loginwindow_autologinpassword\",\n" +
					"        \"settingInstanceTemplateReference\": null,\n" +
					"        \"simpleSettingValue\": {\n" +
					"          \"@odata.type\": \"#microsoft.graph.deviceManagementConfigurationSecretSettingValue\",\n" +
					"          \"valueState\": \"notEncrypted\",\n" +
					"          \"value\": \"your_secret_value\",\n" +
					"          \"settingValueTemplateReference\": null\n" +
					"        }\n" +
					"      }\n" +
					"    }\n" +
					"  ]\n" +
					"})\n" +
					"```\n\n",
				Validators: []validator.String{
					customValidator.JSONSchemaValidator(),
					sharedValidators.SettingsCatalogJSONValidator(),
				},
				PlanModifiers: []planmodifier.String{
					planmodifiers.NormalizeJSONPlanModifier{},
				},
			},
			"platforms": schema.StringAttribute{
				Computed: true,
				MarkdownDescription: "Platform type for this settings catalog policy." +
					"Can be one of: `none`, `android`, `iOS`, `macOS`, `windows10X`, `windows10`, `linux`," +
					"`unknownFutureValue`, `androidEnterprise`, or `aosp`. This is automatically set based on the `settings_catalog_template_type` field.",
				Validators: []validator.String{
					stringvalidator.OneOf(
						"none", "android", "iOS", "macOS", "windows10X",
						"windows10", "linux", "unknownFutureValue",
						"androidEnterprise", "aosp",
					),
				},
				PlanModifiers: []planmodifier.String{planmodifiers.DefaultValueString("none")},
			},
			"technologies": schema.ListAttribute{
				ElementType: types.StringType,
				Computed:    true,
				MarkdownDescription: "Describes a list of technologies this settings catalog setting can be deployed with. Valid values are: `none`," +
					"`mdm`, `windows10XManagement`, `configManager`, `intuneManagementExtension`, `thirdParty`, `documentGateway`, `appleRemoteManagement`, `microsoftSense`," +
					"`exchangeOnline`, `mobileApplicationManagement`, `linuxMdm`, `enrollment`, `endpointPrivilegeManagement`, `unknownFutureValue`, `windowsOsRecovery`, " +
					"and `android`. This is automatically set based on the `settings_catalog_template_type` field.",
				Validators: []validator.List{
					customValidator.StringListAllowedValues(
						"none", "mdm", "windows10XManagement", "configManager",
						"intuneManagementExtension", "thirdParty", "documentGateway",
						"appleRemoteManagement", "microsoftSense", "exchangeOnline",
						"mobileApplicationManagement", "linuxMdm", "enrollment",
						"endpointPrivilegeManagement", "unknownFutureValue",
						"windowsOsRecovery", "android",
					),
				},
			},
			"role_scope_tag_ids": schema.SetAttribute{
				ElementType:         types.StringType,
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Set of scope tag IDs for this Entity instance.",
				PlanModifiers: []planmodifier.Set{
					planmodifiers.DefaultSetValue(
						[]attr.Value{types.StringValue("0")},
					),
				},
			},
			"created_date_time": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				MarkdownDescription: "Creation date and time of the settings catalog policy template",
			},
			"last_modified_date_time": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Last modification date and time of the settings catalog policy template",
			},
			"settings_count": schema.Int32Attribute{
				Computed:            true,
				MarkdownDescription: "Number of settings catalog settings with the policy template. This will change over time as the resource is updated.",
			},
			"is_assigned": schema.BoolAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.Bool{
					planmodifiers.UseStateForUnknownBool(),
				},
				MarkdownDescription: "Indicates if the policy template is assigned to any user or device scope",
			},
			"assignments": commonschemagraphbeta.DeviceConfigurationWithAllGroupAssignmentsAndFilterSchema(),
			"timeouts":    commonschema.Timeouts(ctx),
		},
	}
}
