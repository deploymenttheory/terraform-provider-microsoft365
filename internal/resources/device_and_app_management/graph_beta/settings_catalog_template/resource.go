package graphBetaDeviceManagementTemplate

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common"
	planmodifiers "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/plan_modifiers"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/schema"
	commonschemagraphbeta "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/schema/graph_beta/device_and_app_management"
	customValidator "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/validators"
	sharedValidators "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/validators/graph_beta/device_and_app_management"
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
	ResourceName  = "graph_beta_device_and_app_management_settings_catalog_template"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180
)

var (
	// Basic resource interface (CRUD operations)
	_ resource.Resource = &DeviceManagementTemplateResource{}

	// Allows the resource to be configured with the provider client
	_ resource.ResourceWithConfigure = &DeviceManagementTemplateResource{}

	// Enables import functionality
	_ resource.ResourceWithImportState = &DeviceManagementTemplateResource{}

	// Enables plan modification/diff suppression
	_ resource.ResourceWithModifyPlan = &DeviceManagementTemplateResource{}
)

func NewDeviceManagementTemplateResource() resource.Resource {
	return &DeviceManagementTemplateResource{
		ReadPermissions: []string{
			"DeviceManagementConfiguration.Read.All",
		},
		WritePermissions: []string{
			"DeviceManagementConfiguration.ReadWrite.All",
		},
		ResourcePath: "/deviceManagement/configurationPolicies",
	}
}

type DeviceManagementTemplateResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ProviderTypeName string
	TypeName         string
	ReadPermissions  []string
	WritePermissions []string
	ResourcePath     string
}

// Metadata returns the resource type name.
func (r *DeviceManagementTemplateResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + ResourceName
}

// Configure sets the client for the resource.
func (r *DeviceManagementTemplateResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = common.SetGraphBetaClientForResource(ctx, req, resp, r.TypeName)
}

// ImportState imports the resource state.
func (r *DeviceManagementTemplateResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Function to create the full device management configuration policy schema
func (r *DeviceManagementTemplateResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
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
					"`linux_anti_virus_microsoft_defender_antivirus`: Customers using Microsoft Defender for Endpoint on Linux can configure and deploy Antivirus settings to Linux devices.\n\n" +
					"`linux_anti_virus_microsoft_defender_antivirus_exclusions`: This template allows you to manage settings for Microsoft Defender Antivirus that define Antivirus exclusions for paths, extensions and processes. Antivirus exclusion are also managed by Microsoft Defender Antivirus policy, which includes identical settings for exclusions. Settings from both templates (Antivirus and Antivirus exclusions) are subject to policy merge, and create a super set of exclusions for applicable devices and users.\n\n" +
					"`linux_endpoint_detection_and_response`: Endpoint detection and response settings for Linux devices.\n\n" +
					"`macOS settings catalog templates`\n\n" +
					"`macOS_anti_virus_microsoft_defender_antivirus`: Microsoft Defender Antivirus is the next-generation protection component of Microsoft Defender for Endpoint on Mac. Next-generation protection brings together machine learning, big-data analysis, in-depth threat resistance research, and cloud infrastructure to protect devices in your enterprise organization.\n\n" +
					"`macOS_anti_virus_microsoft_defender_antivirus_exclusions`: This template allows you to manage settings for Microsoft Defender Antivirus that define Antivirus exclusions for paths, extensions and processes. Antivirus exclusion are also managed by Microsoft Defender Antivirus policy, which includes identical settings for exclusions. Settings from both templates (Antivirus and Antivirus exclusions) are subject to policy merge, and create a super set of exclusions for applicable devices and users.\n\n" +
					"`macOS_disk_encryption`: Disk encryption settings for macOS devices.\n\n" +
					"`macOS_endpoint_detection_and_response`: Endpoint detection and response settings for macOS devices.\n\n" +
					//"`macOS_firewall`: Firewall configuration for macOS devices.\n\n" + TODO: uses another api endpoint entirely
					"`Windows settings catalog templates`\n\n" +
					"`windows_account_protection`: Account protection policies help protect user credentials by using technology such as Windows Hello for Business and Credential Guard.\n\n" +
					"`windows_anti_virus_defender_update_controls`: Configure the gradual release rollout of Defender Updates to targeted device groups. Use a ringed approach to test, validate, and rollout updates to devices through release channels. Updates available are platform, engine, security intelligence updates. These policy types have pause, resume, manual rollback commands similar to Windows Update ring policies.\n\n" +
					"`windows_anti_virus_microsoft_defender_antivirus`: Windows Defender Antivirus is the next-generation protection component of Microsoft Defender for Endpoint. Next-generation protection brings together machine learning, big-data analysis, in-depth threat resistance research, and cloud infrastructure to protect devices in your enterprise organization.\n\n" +
					"`windows_anti_virus_microsoft_defender_antivirus_exclusions`: This template allows you to manage settings for Microsoft Defender Antivirus that define Antivirus exclusions for paths, extensions and processes. Antivirus exclusion are also managed by Microsoft Defender Antivirus policy, which includes identical settings for exclusions. Settings from both templates (Antivirus and Antivirus exclusions) are subject to policy merge, and create a super set of exclusions for applicable devices and users.\n\n" +
					"`windows_anti_virus_security_experience`: The Windows Security app is used by a number of Windows security features to provide notifications about the health and security of the machine. These include notifications about firewalls, antivirus products, Windows Defender SmartScreen, and others.\n\n" +
					"`windows_app_control_for_business`: Application control settings for Windows devices.\n\n" +
					"`windows_attack_surface_reduction`: Attack surface reduction rules for Windows devices.\n\n" +
					"`windows_disk_encryption`: Disk encryption settings for Windows devices.\n\n" +
					"`windows_endpoint_detection_and_response`: Endpoint detection and response settings for Windows devices.\n\n" +
					"`windows_firewall`: Firewall settings for Windows devices.\n\n" +
					"`windows_firewall_config_manager`: Firewall configuration manager for Windows devices.\n\n" +
					"`windows_firewall_profile_config_manager`: Profile-specific firewall configuration for Windows devices.\n\n" +
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
					"`windows_(config_mgr)_firewall_rules`: Rules-based firewall configuration for Windows devices managed via Microsoft Configuration Manager.\n",
				Validators: []validator.String{
					stringvalidator.OneOf(
						"linux_anti_virus_microsoft_defender_antivirus",
						"linux_anti_virus_microsoft_defender_antivirus_exclusions",
						"linux_endpoint_detection_and_response",
						"macOS_anti_virus_microsoft_defender_antivirus",
						"macOS_anti_virus_microsoft_defender_antivirus_exclusions",
						"macOS_disk_encryption",
						"macOS_endpoint_detection_and_response",
						//"macOS_firewall",
						"windows_account_protection",
						"windows_anti_virus_defender_update_controls",
						"windows_anti_virus_microsoft_defender_antivirus",
						"windows_anti_virus_microsoft_defender_antivirus_exclusions",
						"windows_anti_virus_security_experience",
						"windows_app_control_for_business",
						"windows_attack_surface_reduction",
						"windows_disk_encryption",
						"windows_endpoint_detection_and_response",
						"windows_firewall",
						"windows_firewall_config_manager",
						"windows_firewall_profile_config_manager",
						"windows_firewall_rules",
						"windows_firewall_rules_config_manager",
						"windows_hyper-v_firewall_rules",
						"windows_local_admin_password_solution_(windows_LAPS)",
						"windows_local_user_group_membership",
						"windows_(config_mgr)_anti_virus_microsoft_defender_antivirus",
						"windows_(config_mgr)_anti_virus_windows_security_experience",
						"windows_(config_mgr)_attack_surface_reduction",
						"windows_(config_mgr)_endpoint_detection_and_response",
						"windows_(config_mgr)_firewall",
						"windows_(config_mgr)_firewall_profile",
						"windows_(config_mgr)_firewall_rules",
					),
				},
			},
			"settings": schema.StringAttribute{
				Required: true,
				MarkdownDescription: "Settings Catalog Policy template settings defined as a JSON string. Please provide a valid JSON-encoded settings structure. " +
					"This can either be extracted from an existing policy using the Intune gui `export JSON` functionality if supported, via a script such as this powershell script." +
					" [ExportSettingsCatalogTemplateConfigurationById](https://github.com/deploymenttheory/terraform-provider-microsoft365/blob/main/scripts/ExportSettingsCatalogTemplateConfigurationById.ps1) " +
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
					sharedValidators.SettingsCatalogValidator(),
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
			"role_scope_tag_ids": schema.ListAttribute{
				ElementType:         types.StringType,
				Optional:            true,
				MarkdownDescription: "List of scope tag IDs for this Settings Catalog template profile.",
				PlanModifiers: []planmodifier.List{
					planmodifiers.DefaultListValue(
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
			"settings_count": schema.Int64Attribute{
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
			"assignments": commonschemagraphbeta.ConfigurationPolicyAssignmentsSchema(),
			"timeouts":    commonschema.Timeouts(ctx),
		},
	}
}
