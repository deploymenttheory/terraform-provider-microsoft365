package graphBetaWindowsPlatformScript

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common"
	planmodifiers "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/plan_modifiers"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
	commonschemagraphbeta "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema/graph_beta/device_management"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

const (
	ResourceName  = "graph_beta_device_management_windows_platform_script"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180
)

var (
	// Basic resource interface (CRUD operations)
	_ resource.Resource = &WindowsPlatformScriptResource{}

	// Allows the resource to be configured with the provider client
	_ resource.ResourceWithConfigure = &WindowsPlatformScriptResource{}

	// Enables import functionality
	_ resource.ResourceWithImportState = &WindowsPlatformScriptResource{}

	// Enables plan modification/diff suppression
	_ resource.ResourceWithModifyPlan = &WindowsPlatformScriptResource{}
)

func NewWindowsPlatformScriptResource() resource.Resource {
	return &WindowsPlatformScriptResource{
		ReadPermissions: []string{
			"",
		},
		WritePermissions: []string{
			"DeviceManagementConfiguration.ReadWrite.All",
			"DeviceManagementScripts.ReadWrite.All",
			"DeviceManagementManagedDevices.ReadWrite.All",
		},
		ResourcePath: "/deviceManagement/deviceManagementScripts",
	}
}

type WindowsPlatformScriptResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ProviderTypeName string
	TypeName         string
	ReadPermissions  []string
	WritePermissions []string
	ResourcePath     string
}

// Metadata returns the resource type name.
func (r *WindowsPlatformScriptResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	r.ProviderTypeName = req.ProviderTypeName
	r.TypeName = ResourceName
	resp.TypeName = r.FullTypeName()
}

// FullTypeName returns the full resource type name in the format "providername_resourcename".
func (r *WindowsPlatformScriptResource) FullTypeName() string {
	return r.ProviderTypeName + "_" + r.TypeName
}

// Configure sets the client for the resource.
func (r *WindowsPlatformScriptResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = common.SetGraphBetaClientForResource(ctx, req, resp, constants.PROVIDER_NAME+"_"+ResourceName)
}

// ImportState imports the resource state.
func (r *WindowsPlatformScriptResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *WindowsPlatformScriptResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages Windows PowerShell scripts using the `/deviceManagement/deviceManagementScripts` endpoint. Windows platform scripts enable automated deployment and execution of PowerShell scripts on managed Windows devices, supporting both system and user contexts with configurable signature checking and 32-bit execution options.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					planmodifiers.UseStateForUnknownString(),
				},
				MarkdownDescription: "The unique identifier for this Intune windows platform script",
			},
			"display_name": schema.StringAttribute{
				MarkdownDescription: "Name of the windows platform script.",
				Required:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Optional description for the windows platform script.",
				Optional:            true,
			},
			"script_content": schema.StringAttribute{
				MarkdownDescription: "The script content.",
				Required:            true,
				Sensitive:           true,
			},
			"run_as_account": schema.StringAttribute{
				MarkdownDescription: "Indicates the type of execution context. Possible values are: `system`, `user`.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("system", "user"),
				},
			},
			"enforce_signature_check": schema.BoolAttribute{
				MarkdownDescription: "Indicate whether the script signature needs be checked.",
				Optional:            true,
			},
			"file_name": schema.StringAttribute{
				MarkdownDescription: "Script file name.",
				Required:            true,
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
			"run_as_32_bit": schema.BoolAttribute{
				MarkdownDescription: "A value indicating whether the PowerShell script should run as 32-bit.",
				Optional:            true,
			},
			"assignments": commonschemagraphbeta.PlatformScriptAssignmentsSchema(),
			"timeouts":    commonschema.Timeouts(ctx),
		},
	}
}
