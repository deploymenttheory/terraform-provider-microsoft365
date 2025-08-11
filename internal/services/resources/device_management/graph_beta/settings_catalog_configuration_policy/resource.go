package graphBetaSettingsCatalogConfigurationPolicy

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	planmodifiers "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/plan_modifiers"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
	commonschemagraphbeta "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema/graph_beta/device_management"
	customValidator "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/validators"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

const (
	ResourceName  = "graph_beta_device_management_settings_catalog_configuration_policy"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180
)

var (
	// Basic resource interface (CRUD operations)
	_ resource.Resource = &SettingsCatalogResource{}

	// Allows the resource to be configured with the provider client
	_ resource.ResourceWithConfigure = &SettingsCatalogResource{}

	// Enables import functionality
	_ resource.ResourceWithImportState = &SettingsCatalogResource{}

	// Enables plan modification/diff suppression
	_ resource.ResourceWithModifyPlan = &SettingsCatalogResource{}
)

func NewSettingsCatalogResource() resource.Resource {
	return &SettingsCatalogResource{
		ReadPermissions: []string{
			"DeviceManagementConfiguration.Read.All",
		},
		WritePermissions: []string{
			"DeviceManagementConfiguration.ReadWrite.All",
		},
		ResourcePath: "/deviceManagement/configurationPolicies",
	}
}

type SettingsCatalogResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ProviderTypeName string
	TypeName         string
	ReadPermissions  []string
	WritePermissions []string
	ResourcePath     string
}

// Metadata returns the resource type name.
func (r *SettingsCatalogResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	r.ProviderTypeName = req.ProviderTypeName
	r.TypeName = ResourceName
	resp.TypeName = r.FullTypeName()
}

// FullTypeName returns the full type name of the resource for logging purposes.
func (r *SettingsCatalogResource) FullTypeName() string {
	return r.ProviderTypeName + "_" + r.TypeName
}

// Configure sets the client for the resource.
func (r *SettingsCatalogResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, constants.PROVIDER_NAME+"_"+ResourceName)
}

// ImportState imports the resource state.
func (r *SettingsCatalogResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Function to create the full device management configuration policy schema
func (r *SettingsCatalogResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages Intune Settings Catalog policies using the `/deviceManagement/configurationPolicies` endpoint. " +
			"Settings Catalog provides a unified configuration experience for Windows, macOS, iOS/iPadOS, and Android devices through a modern, " +
			"simplified interface that replaces traditional device configuration profiles and legacy Intune configuration templates. You can simplify the hcl creation process by using the " +
			"`Export-IntuneSettingsCatalogConfigurationToHCL.ps1` [https://github.com/deploymenttheory/terraform-provider-microsoft365/blob/main/scripts/powershell/Export-IntuneSettingsCatalogConfigurationToHCL.ps1] " +
			"script to export the settings catalog and settings catalog templates. You can export by a singular resource ID or by exporting all policies. " +
			"This will build the hcl representation of the settings catalog configuration programmatically.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					planmodifiers.UseStateForUnknownString(),
				},
				MarkdownDescription: "The unique identifier for this policy",
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Policy name",
			},
			"description": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				PlanModifiers:       []planmodifier.String{planmodifiers.DefaultValueString("")},
				MarkdownDescription: "Optional description for the settings catalog policy.",
			},
			"template_reference": schema.SingleNestedAttribute{
				Required:            true,
				MarkdownDescription: "Policy template reference information",
				Attributes: map[string]schema.Attribute{
					"template_id": schema.StringAttribute{
						Optional:            true,
						Computed:            true,
						Default:             stringdefault.StaticString(""),
						MarkdownDescription: "Template ID to reference. This is a UUID that identifies a specific template in Microsoft Intune.",
					},
					"template_family": schema.StringAttribute{
						Computed:            true,
						MarkdownDescription: "Describes the TemplateFamily for the Template entity. This is a read-only property.",
						Validators: []validator.String{
							stringvalidator.OneOf(
								"none",
								"endpointSecurityAntivirus",
								"endpointSecurityDiskEncryption",
								"endpointSecurityFirewall",
								"endpointSecurityEndpointDetectionAndResponse",
								"endpointSecurityAttackSurfaceReduction",
								"endpointSecurityAccountProtection",
								"endpointSecurityApplicationControl",
								"endpointSecurityEndpointPrivilegeManagement",
								"enrollmentConfiguration",
								"appQuietTime",
								"baseline",
								"unknownFutureValue",
								"deviceConfigurationScripts",
								"deviceConfigurationPolicies",
								"windowsOsRecoveryPolicies",
								"companyPortal",
							),
						},
					},
					"template_display_name": schema.StringAttribute{
						Computed:            true,
						MarkdownDescription: "Template Display Name of the referenced template. This property is read-only.",
					},
					"template_display_version": schema.StringAttribute{
						Computed:            true,
						MarkdownDescription: "Template Display Version of the referenced Template. This property is read-only.",
					},
				},
			},
			"configuration_policy": schema.SingleNestedAttribute{
				Optional:            true,
				MarkdownDescription: "Settings Catalog (configuration policy) settings",
				Attributes:          DeviceConfigV2Attributes(),
			},
			"platforms": schema.StringAttribute{
				Optional: true,
				Computed: true,
				MarkdownDescription: "Platform type for this settings catalog policy." +
					"Can be one of: `none`, `android`, `iOS`, `macOS`, `windows10X`, `windows10`, `linux`," +
					"`unknownFutureValue`, `androidEnterprise`, or `aosp`. Defaults to `none`.",
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
				Optional:    true,
				Computed:    true,
				MarkdownDescription: "Describes a list of technologies this settings catalog setting can be deployed with. Valid values are:" +
					" `none`, `mdm`, `windows10XManagement`, `configManager`, `intuneManagementExtension`, `thirdParty`, `documentGateway`, `appleRemoteManagement`," +
					" `microsoftSense`, `exchangeOnline`, `mobileApplicationManagement`, `linuxMdm`, `enrollment`, `endpointPrivilegeManagement`, `unknownFutureValue`, " +
					"`windowsOsRecovery`, and `android`. Defaults to `mdm`.",
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
				PlanModifiers: []planmodifier.List{
					planmodifiers.DefaultListValue([]attr.Value{types.StringValue("mdm")}),
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
				MarkdownDescription: "Creation date and time of the settings catalog policy",
			},
			"last_modified_date_time": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Last modification date and time of the settings catalog policy",
			},
			"settings_count": schema.Int32Attribute{
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
			"assignments": commonschemagraphbeta.DeviceConfigurationWithAllGroupAssignmentsAndFilterSchema(),
			"timeouts":    commonschema.Timeouts(ctx),
		},
	}
}
