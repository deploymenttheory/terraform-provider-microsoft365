package graphBetaDeviceShellScript

import (
	"context"
	"regexp"

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
	ResourceName  = "graph_beta_device_and_app_management_macos_platform_script"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180
)

var (
	// Basic resource interface (CRUD operations)
	_ resource.Resource = &DeviceShellScriptResource{}

	// Allows the resource to be configured with the provider client
	_ resource.ResourceWithConfigure = &DeviceShellScriptResource{}

	// Enables import functionality
	_ resource.ResourceWithImportState = &DeviceShellScriptResource{}

	// Enables plan modification/diff suppression
	_ resource.ResourceWithModifyPlan = &DeviceShellScriptResource{}
)

func NewDeviceShellScriptResource() resource.Resource {
	return &DeviceShellScriptResource{
		ReadPermissions: []string{
			"DeviceManagementConfiguration.Read.All",
			"DeviceManagementManagedDevices.Read.All",
		},
		WritePermissions: []string{
			"DeviceManagementManagedDevices.ReadWrite.All",
			"DeviceManagementConfiguration.ReadWrite.All",
		},
		ResourcePath: "/deviceManagement/deviceShellScripts",
	}
}

type DeviceShellScriptResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ProviderTypeName string
	TypeName         string
	ReadPermissions  []string
	WritePermissions []string
	ResourcePath     string
}

// Metadata returns the resource type name.
func (r *DeviceShellScriptResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + ResourceName
}

// Configure sets the client for the resource.
func (r *DeviceShellScriptResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = common.SetGraphBetaClientForResource(ctx, req, resp, r.TypeName)
}

// ImportState imports the resource state.
func (r *DeviceShellScriptResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *DeviceShellScriptResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages an Intune macOS platform script using the 'deviceShellScripts' Graph Beta API.",
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
			"created_date_time": schema.StringAttribute{
				Description: "The date and time the device management script was created. This property is read-only.",
				Computed:    true,
			},
			"last_modified_date_time": schema.StringAttribute{
				Description: "The date and time the device management script was last modified. This property is read-only.",
				Computed:    true,
			},
			"run_as_account": schema.StringAttribute{
				Description: "Indicates the type of execution context. Possible values are: `system`, `user`.",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("system", "user"),
				},
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
			"block_execution_notifications": schema.BoolAttribute{
				Description: "Does not notify the user a script is being executed.",
				Optional:    true,
			},
			"execution_frequency": schema.StringAttribute{
				Optional:    true,
				Description: "The interval for script to run in ISO 8601 duration format (e.g., PT1H for 1 hour, P1D for 1 day). If not defined the script will run once.",
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^P(?:\d+Y)?(?:\d+M)?(?:\d+W)?(?:\d+D)?(?:T(?:\d+H)?(?:\d+M)?(?:\d+S)?)?$`),
						"must be a valid ISO 8601 duration",
					),
				},
			},
			"retry_count": schema.Int32Attribute{
				Description: "Number of times for the script to be retried if it fails.",
				Optional:    true,
			},
			"assignments": commonschema.ScriptAssignmentsSchema(),
			"timeouts":    commonschema.Timeouts(ctx),
		},
	}
}
