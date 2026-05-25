package graphBetaSettingsCatalogInventoryPolicy

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	planmodifiers "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/plan_modifiers"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
	commonschemagraphbeta "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema/graph_beta/device_management"
	configPolicy "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_management/graph_beta/settings_catalog_configuration_policy"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/identityschema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

const (
	ResourceName  = "microsoft365_graph_beta_device_management_settings_catalog_inventory_policy"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180
)

var (
	_ resource.Resource                = &InventoryPolicyResource{}
	_ resource.ResourceWithConfigure   = &InventoryPolicyResource{}
	_ resource.ResourceWithImportState = &InventoryPolicyResource{}
	_ resource.ResourceWithModifyPlan  = &InventoryPolicyResource{}
	_ resource.ResourceWithIdentity    = &InventoryPolicyResource{}
)

func NewInventoryPolicyResource() resource.Resource {
	return &InventoryPolicyResource{
		ReadPermissions: []string{
			"DeviceManagementConfiguration.Read.All",
		},
		WritePermissions: []string{
			"DeviceManagementConfiguration.ReadWrite.All",
		},
		ResourcePath: "/deviceManagement/inventoryPolicies",
	}
}

type InventoryPolicyResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ReadPermissions  []string
	WritePermissions []string
	ResourcePath     string
}

func (r *InventoryPolicyResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = ResourceName
}

func (r *InventoryPolicyResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, ResourceName)
}

func (r *InventoryPolicyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *InventoryPolicyResource) IdentitySchema(ctx context.Context, req resource.IdentitySchemaRequest, resp *resource.IdentitySchemaResponse) {
	resp.IdentitySchema = identityschema.Schema{
		Attributes: map[string]identityschema.Attribute{
			"id": identityschema.StringAttribute{
				RequiredForImport: true,
			},
		},
	}
}

func (r *InventoryPolicyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages Intune Settings Catalog Inventory policies using the `/deviceManagement/inventoryPolicies` endpoint. This resource enables collection of application inventory data from managed Windows devices through the extensibility technology platform.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					planmodifiers.UseStateForUnknownString(),
				},
				MarkdownDescription: "The unique identifier for this inventory policy.",
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Policy name.",
			},
			"description": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Policy description. Maximum 1500 characters.",
				Validators: []validator.String{
					stringvalidator.LengthAtMost(1500),
				},
			},
			"configuration_policy": schema.SingleNestedAttribute{
				Optional:            true,
				MarkdownDescription: "Settings Catalog inventory policy settings.",
				Attributes:          configPolicy.DeviceConfigV2Attributes(),
			},
			"platforms": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Platform type for this inventory policy.",
				Validators: []validator.String{
					stringvalidator.OneOf(
						"none", "android", "iOS", "macOS", "windows10X",
						"windows10", "linux", "unknownFutureValue",
						"androidEnterprise", "aosp",
					),
				},
				PlanModifiers: []planmodifier.String{planmodifiers.DefaultValueString("none")},
			},
			"technologies": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("extensibility"),
				MarkdownDescription: "Technology for deployment. Valid values: `extensibility`.",
				Validators: []validator.String{
					stringvalidator.OneOf("extensibility"),
				},
			},
			"role_scope_tag_ids": schema.SetAttribute{
				ElementType:         types.StringType,
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Set of scope tag IDs for this entity instance.",
				PlanModifiers: []planmodifier.Set{
					planmodifiers.DefaultSetValue(
						[]attr.Value{types.StringValue("0")},
					),
				},
			},
			"created_date_time": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Creation date and time of the inventory policy.",
			},
			"last_modified_date_time": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Last modification date and time of the inventory policy.",
			},
			"settings_count": schema.Int32Attribute{
				Computed:            true,
				MarkdownDescription: "Number of settings within the policy.",
			},
			"assignments": commonschemagraphbeta.DeviceConfigurationWithAllGroupAssignmentsAndFilterSchema(),
			"timeouts":    commonschema.ResourceTimeouts(ctx),
		},
	}
}
