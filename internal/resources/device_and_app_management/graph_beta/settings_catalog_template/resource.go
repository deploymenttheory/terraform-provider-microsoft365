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
		Description: "Manages a Settings Catalog policy template in Microsoft Intune for Windows, macOS, iOS/iPadOS and Android.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					planmodifiers.UseStateForUnknownString(),
				},
				MarkdownDescription: "The unique identifier for this policy template",
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
				MarkdownDescription: "Defines which device management template type with settings catalog setting that will be deployed. " +
					"Options available are `macOS_disk_encryption`, `macOS_firewall`, `macOS_endpoint_detection_and_response`, `macOS_anti_virus`, " +
					"`windows_account_protection`, `windows_anti_virus`, `windows_app_control_for_business`, `windows_attack_surface_reduction`, " +
					"`windows_disk_encryption`, `windows_firewall`, `windows_firewall_rules`, `windows_hyper-v_firewall_rules`, " +
					"`windows_firewall_config_manager`, `windows_firewall_profile_config_manager`, `windows_firewall_rules_config_manager`, " +
					"`windows_endpoint_detection_and_response`, `windows_config_manager_anti_virus`, `windows_config_manager_attack_surface_reduction`, " +
					"`windows_config_manager_endpoint_detection_and_response`, `linux_endpoint_detection_and_response`, `linux_anti_virus`.",
				Validators: []validator.String{
					stringvalidator.OneOf(
						"linux_endpoint_detection_and_response",
						"linux_anti_virus_microsoft_defender_antivirus",
						"linux_anti_virus_microsoft_defender_antivirus_exclusions",
						"macOS_disk_encryption",
						"macOS_firewall",
						"macOS_endpoint_detection_and_response",
						"macOS_anti_virus_microsoft_defender_antivirus",
						"macOS_anti_virus_microsoft_defender_antivirus_exclusions",
						"windows_account_protection",
						"windows_anti_virus_defender_update_controls",
						"windows_anti_virus_microsoft_defender_antivirus",
						"windows_anti_virus_microsoft_defender_antivirus_exclusions",
						"windows_anti_virus_security_experience",
						"windows_app_control_for_business",
						"windows_attack_surface_reduction",
						"windows_disk_encryption",
						"windows_firewall",
						"windows_firewall_rules",
						"windows_hyper-v_firewall_rules",
						"windows_firewall_config_manager",
						"windows_firewall_profile_config_manager",
						"windows_firewall_rules_config_manager",
						"windows_endpoint_detection_and_response",
						"windows_config_manager_anti_virus_microsoft_defender_antivirus",
						"windows_config_manager_anti_virus_windows_security_experience",
						"windows_config_manager_attack_surface_reduction",
						"windows_config_manager_endpoint_detection_and_response",
					),
				},
			},
			"settings": schema.StringAttribute{
				Required: true,
				MarkdownDescription: "Settings Catalog Policy template settings defined as a valid JSON string. Provide JSON-encoded settings structure. " +
					"This can either be extracted from an existing policy using the Intune gui `export JSON` functionality, via a script such as" +
					" [this PowerShell script](https://github.com/deploymenttheory/terraform-provider-microsoft365/blob/main/scripts/ExportSettingsCatalogConfigurationById.ps1) " +
					"or created from scratch. The JSON structure should match the graph schema of the settings catalog. Please look at the " +
					"terraform documentation for the settings catalog for examples and how to correctly format the HCL.\n\n" +
					"A correctly formatted field in the HCL should begin and end like this:\n" +
					"```hcl\n" +
					"settings = jsonencode({\n" +
					"  \"settings\": [\n" +
					"    {\n" +
					"        \"id\": \"0\",\n" +
					"        \"settingInstance\": {\n" +
					"            }\n" +
					"        }\n" +
					"    },\n" +
					"```\n\n" +
					"Note: When setting secret values (identified by `@odata.type: \"#microsoft.graph.deviceManagementConfigurationSecretSettingValue\"`), " +
					"ensure the `valueState` is set to `\"notEncrypted\"`. The value `\"encryptedValueToken\"` is reserved for server responses and " +
					"should not be used when creating or updating settings.",
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
					"Can be one of: none, android, iOS, macOS, windows10X, windows10, linux," +
					"unknownFutureValue, androidEnterprise, or aosp. Defaults to none.",
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
				ElementType:         types.StringType,
				Computed:            true,
				MarkdownDescription: "Describes a list of technologies this settings catalog setting can be deployed with. Valid values are: none, mdm, windows10XManagement, configManager, intuneManagementExtension, thirdParty, documentGateway, appleRemoteManagement, microsoftSense, exchangeOnline, mobileApplicationManagement, linuxMdm, enrollment, endpointPrivilegeManagement, unknownFutureValue, windowsOsRecovery, and android. Defaults to ['mdm'].",
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
				MarkdownDescription: "List of scope tag IDs for this Windows Settings Catalog profile.",
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
				MarkdownDescription: "Creation date and time of the settings catalog policy",
			},
			"last_modified_date_time": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Last modification date and time of the settings catalog policy",
			},
			"settings_count": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "Number of settings catalog settings with the policy. This will change over time as the resource is updated.",
			},
			"is_assigned": schema.BoolAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.Bool{
					planmodifiers.UseStateForUnknownBool(),
				},
				MarkdownDescription: "Indicates if the policy is assigned to any scope",
			},
			"assignments": commonschemagraphbeta.ConfigurationPolicyAssignmentsSchema(),
			"timeouts":    commonschema.Timeouts(ctx),
		},
	}
}
