package graphBetaSettingsCatalog

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common"
	planmodifiers "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/plan_modifiers"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/schema"
	customValidator "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/validators"
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
			"DeviceManagementApps.Read.All",
			"DeviceManagementConfiguration.Read.All",
		},
		WritePermissions: []string{
			"DeviceManagementApps.ReadWrite.All",
			"DeviceManagementConfiguration.ReadWrite.All",
		},
	}
}

type SettingsCatalogResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ProviderTypeName string
	TypeName         string
	ReadPermissions  []string
	WritePermissions []string
}

// Metadata returns the resource type name.
func (r *SettingsCatalogResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_graph_beta_device_and_app_management_settings_catalog_v3"
}

// Configure sets the client for the resource.
func (r *SettingsCatalogResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = common.SetGraphBetaClientForResource(ctx, req, resp, r.TypeName)
}

// ImportState imports the resource state.
func (r *SettingsCatalogResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Function to create the full device management configuration policy schema
func (r *SettingsCatalogResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Settings Catalog profile in Microsoft Intune.",
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
				MarkdownDescription: "Policy description",
			},
			"settings": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Settings Catalog Policy settings defined as a valid JSON string, supporting up to 20 levels of nesting. Provide JSON-encoded settings structure.",
				Validators: []validator.String{
					customValidator.JSONSchemaValidator(),
				},
				PlanModifiers: []planmodifier.String{
					planmodifiers.UseStateForUnknownString(),
				},
			},
			"platforms": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					customValidator.EnumValues(
						"none", "android", "iOS", "macOS", "windows10X",
						"windows10", "linux", "unknownFutureValue",
						"androidEnterprise", "aosp",
					),
				},
				PlanModifiers: []planmodifier.String{planmodifiers.DefaultValueString("none")},

				MarkdownDescription: "Platforms for this policy",
			},
			"technologies": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Computed:    true,
				Validators: []validator.List{
					customValidator.EnumValuesList(
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
				MarkdownDescription: "Describes a list of technologies this settings catalog setting can be deployed with. Defaults to 'mdm'.",
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
			"assignments": commonschema.AssignmentsSchema(),

			"timeouts": commonschema.Timeouts(ctx),
		},
	}
}