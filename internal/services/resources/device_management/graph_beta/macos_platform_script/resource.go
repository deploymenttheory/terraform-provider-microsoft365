package graphBetaMacOSPlatformScript

import (
	"context"
	"regexp"

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
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

const (
	ResourceName  = "graph_beta_device_management_macos_platform_script"
	CreateTimeout = 600
	UpdateTimeout = 600
	ReadTimeout   = 180
	DeleteTimeout = 180
)

var (
	// Basic resource interface (CRUD operations)
	_ resource.Resource = &MacOSPlatformScriptResource{}

	// Allows the resource to be configured with the provider client
	_ resource.ResourceWithConfigure = &MacOSPlatformScriptResource{}

	// Enables import functionality
	_ resource.ResourceWithImportState = &MacOSPlatformScriptResource{}

	// Enables plan modification/diff suppression
	_ resource.ResourceWithModifyPlan = &MacOSPlatformScriptResource{}
)

func NewMacOSPlatformScriptResource() resource.Resource {
	return &MacOSPlatformScriptResource{
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

type MacOSPlatformScriptResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ProviderTypeName string
	TypeName         string
	ReadPermissions  []string
	WritePermissions []string
	ResourcePath     string
}

// Metadata returns the resource type name.
func (r *MacOSPlatformScriptResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	r.ProviderTypeName = req.ProviderTypeName
	r.TypeName = ResourceName
	resp.TypeName = r.FullTypeName()
}

// FullTypeName returns the full type name of the resource for logging purposes.
func (r *MacOSPlatformScriptResource) FullTypeName() string {
	return r.ProviderTypeName + "_" + r.TypeName
}

// Configure sets the client for the resource.
func (r *MacOSPlatformScriptResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, constants.PROVIDER_NAME+"_"+ResourceName)
}

// ImportState imports the resource state.
func (r *MacOSPlatformScriptResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *MacOSPlatformScriptResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages macOS shell scripts using the `/deviceManagement/deviceShellScripts` endpoint. macOS platform scripts enable automated deployment and execution of shell scripts on managed macOS devices with support for scheduled execution, retry logic, and execution context control for system maintenance and configuration tasks.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					planmodifiers.UseStateForUnknownString(),
				},
				MarkdownDescription: "The unique identifier of the macOS platform script.",
			},
			"display_name": schema.StringAttribute{
				MarkdownDescription: "Name of the macOS Platform Script.",
				Required:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Optional description for the macOS Platform Script.",
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
			"block_execution_notifications": schema.BoolAttribute{
				MarkdownDescription: "Does not notify the user a script is being executed.",
				Optional:            true,
			},
			"execution_frequency": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The interval for script to run in ISO 8601 duration format (e.g., PT1H for 1 hour, P1D for 1 day). If not defined the script will run once.",
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(constants.ISO8601DurationRegex),
						"must be a valid ISO 8601 duration. Examples: 'P1D' (1 day), 'PT1H' (1 hour), 'P1W' (1 week), 'P1Y2M3DT4H5M6S' (1 year, 2 months, 3 days, 4 hours, 5 minutes, 6 seconds)",
					),
				},
			},
			"retry_count": schema.Int32Attribute{
				MarkdownDescription: "Number of times for the script to be retried if it fails.",
				Optional:            true,
				PlanModifiers: []planmodifier.Int32{
					planmodifiers.UseStateForUnknownInt32(),
				},
			},
			"assignments": commonschemagraphbeta.PlatformScriptAssignmentsSchema(),
			"timeouts":    commonschema.Timeouts(ctx),
		},
	}
}
