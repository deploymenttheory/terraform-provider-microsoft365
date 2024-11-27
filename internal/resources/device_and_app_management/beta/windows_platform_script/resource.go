package graphBetaWindowsPlatformScript

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/schema"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

const (
	ResourceName = "graph_beta_device_and_app_management_windows_platform_script"
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
			"DeviceManagementConfiguration.Read.All",
		},
		WritePermissions: []string{
			"DeviceManagementConfiguration.ReadWrite.All",
		},
	}
}

type WindowsPlatformScriptResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ProviderTypeName string
	TypeName         string
	ReadPermissions  []string
	WritePermissions []string
}

// Metadata returns the resource type name.
func (r *WindowsPlatformScriptResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + ResourceName
}

// Configure sets the client for the resource.
func (r *WindowsPlatformScriptResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = common.SetGraphBetaClientForResource(ctx, req, resp, r.TypeName)
}

// ImportState imports the resource state.
func (r *WindowsPlatformScriptResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *WindowsPlatformScriptResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages an Intune windows platform script using the 'deviceManagementScripts' Graph Beta API.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Unique Identifier for the device management script.",
				Computed:    true,
			},
			"display_name": schema.StringAttribute{
				Description: "Name of the device management script.",
				Required:    true,
			},
			"description": schema.StringAttribute{
				Description: "Optional description for the device management script.",
				Optional:    true,
			},
			"script_content": schema.StringAttribute{
				Description: "The script content.",
				Required:    true,
				Sensitive:   true,
			},
			"run_as_account": schema.StringAttribute{
				Description: "Indicates the type of execution context. Possible values are: `system`, `user`.",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("system", "user"),
				},
			},
			"enforce_signature_check": schema.BoolAttribute{
				Description: "Indicate whether the script signature needs be checked.",
				Optional:    true,
			},
			"file_name": schema.StringAttribute{
				Description: "Script file name.",
				Required:    true,
			},
			"role_scope_tag_ids": schema.ListAttribute{
				Description: "List of Scope Tag IDs for this PowerShellScript instance.",
				Optional:    true,
				ElementType: types.StringType,
			},
			"run_as_32_bit": schema.BoolAttribute{
				Description: "A value indicating whether the PowerShell script should run as 32-bit.",
				Optional:    true,
			},
			"assignments": commonschema.ScriptAssignmentsSchema(),
			"timeouts":    commonschema.Timeouts(ctx),
		},
	}
}
