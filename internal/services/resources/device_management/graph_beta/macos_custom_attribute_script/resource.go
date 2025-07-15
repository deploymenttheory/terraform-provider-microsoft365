package graphBetaMacOSCustomAttributeScript

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	planmodifiers "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/plan_modifiers"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
	commonschemagraphbeta "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema/graph_beta/device_management"
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
	ResourceName  = "graph_beta_device_management_macos_custom_attribute_script"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180
)

var (
	// Basic resource interface (CRUD operations)
	_ resource.Resource = &DeviceCustomAttributeShellScriptResource{}

	// Allows the resource to be configured with the provider client
	_ resource.ResourceWithConfigure = &DeviceCustomAttributeShellScriptResource{}

	// Enables import functionality
	_ resource.ResourceWithImportState = &DeviceCustomAttributeShellScriptResource{}

	// Enables plan modification/diff suppression
	_ resource.ResourceWithModifyPlan = &DeviceCustomAttributeShellScriptResource{}
)

func NewDeviceCustomAttributeShellScriptResource() resource.Resource {
	return &DeviceCustomAttributeShellScriptResource{
		ReadPermissions: []string{
			"DeviceManagementConfiguration.Read.All",
			"DeviceManagementManagedDevices.Read.All",
		},
		WritePermissions: []string{
			"DeviceManagementManagedDevices.ReadWrite.All",
			"DeviceManagementConfiguration.ReadWrite.All",
		},
		ResourcePath: "/deviceManagement/deviceCustomAttributeShellScripts",
	}
}

type DeviceCustomAttributeShellScriptResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ProviderTypeName string
	TypeName         string
	ReadPermissions  []string
	WritePermissions []string
	ResourcePath     string
}

func (r *DeviceCustomAttributeShellScriptResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	r.ProviderTypeName = req.ProviderTypeName
	r.TypeName = ResourceName
	resp.TypeName = r.FullTypeName()
}

func (r *DeviceCustomAttributeShellScriptResource) FullTypeName() string {
	return r.ProviderTypeName + "_" + r.TypeName
}

func (r *DeviceCustomAttributeShellScriptResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, constants.PROVIDER_NAME+"_"+ResourceName)
}

func (r *DeviceCustomAttributeShellScriptResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *DeviceCustomAttributeShellScriptResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a device custom attribute shell script using the `/deviceManagement/deviceCustomAttributeShellScripts` endpoint.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					planmodifiers.UseStateForUnknownString(),
				},
				MarkdownDescription: "The unique identifier for the custom attribute shell script.",
			},
			"display_name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Name of the device management script.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"description": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Optional description for the device management script.",
			},
			"custom_attribute_name": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The name of the custom attribute.",
			},
			"custom_attribute_type": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The expected type of the custom attribute's value. Possible values: integer, string, dateTime.",
				Validators: []validator.String{
					stringvalidator.OneOf("integer", "string", "dateTime"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"script_content": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The script content.",
			},
			"run_as_account": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Indicates the type of execution context. Possible values: system, user.",
				Validators: []validator.String{
					stringvalidator.OneOf("system", "user"),
				},
			},
			"file_name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Script file name.",
			},
			"role_scope_tag_ids": schema.SetAttribute{
				ElementType:         types.StringType,
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Set of scope tag IDs for this Settings Catalog template profile.",
				PlanModifiers: []planmodifier.Set{
					planmodifiers.DefaultSetValue(
						[]attr.Value{types.StringValue("0")},
					),
				},
			},
			"created_date_time": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The date and time the script was created. Read-only.",
			},
			"last_modified_date_time": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The date and time the script was last modified. Read-only.",
			},
			"assignments": commonschemagraphbeta.PlatformScriptAssignmentsSchema(),
			"timeouts":    commonschema.Timeouts(ctx),
		},
	}
}
