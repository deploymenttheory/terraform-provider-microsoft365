package graphBetaMacOSPlatformScript

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/schema"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

const (
	ResourceName = "graph_beta_device_and_app_management_macos_platform_script"
)

var (
	// Basic resource interface (CRUD operations)
	_ datasource.DataSource = &MacOSPlatformScriptDataSource{}

	// Allows the resource to be configured with the provider client
	_ datasource.DataSourceWithConfigure = &MacOSPlatformScriptDataSource{}
)

func NewMacOSPlatformScriptDataSource() datasource.DataSource {
	return &MacOSPlatformScriptDataSource{
		ReadPermissions: []string{
			"DeviceManagementConfiguration.Read.All",
			"DeviceManagementManagedDevices.Read.All",
		},
	}
}

type MacOSPlatformScriptDataSource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ProviderTypeName string
	TypeName         string
	ReadPermissions  []string
}

// Metadata returns the resource type name.
func (r *MacOSPlatformScriptDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + ResourceName
}

func (d *MacOSPlatformScriptDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages an Intune macOS platform script using the 'MacOSPlatformScripts' Graph Beta API.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				Optional:            true,
				MarkdownDescription: "Unique Identifier for the macOS Platform Script.",
			},
			"display_name": schema.StringAttribute{
				Computed:            true,
				Optional:            true,
				MarkdownDescription: "Name of the macOS Platform Script.",
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Optional description for the macOS Platform Script.",
				Optional:            true,
			},
			"script_content": schema.StringAttribute{
				MarkdownDescription: "The script content.",
				Computed:            true,
				Sensitive:           true,
			},
			"created_date_time": schema.StringAttribute{
				MarkdownDescription: "The date and time the macOS Platform Script was created. This property is read-only.",
				Computed:            true,
			},
			"last_modified_date_time": schema.StringAttribute{
				MarkdownDescription: "The date and time the macOS Platform Script was last modified. This property is read-only.",
				Computed:            true,
			},
			"run_as_account": schema.StringAttribute{
				MarkdownDescription: "Indicates the type of execution context. Possible values are: `system`, `user`.",
				Computed:            true,
			},
			"file_name": schema.StringAttribute{
				MarkdownDescription: "Script file name.",
				Computed:            true,
			},
			"role_scope_tag_ids": schema.ListAttribute{
				MarkdownDescription: "List of Scope Tag IDs for this PowerShellScript instance.",
				Optional:            true,
				ElementType:         types.StringType,
			},
			"block_execution_notifications": schema.BoolAttribute{
				MarkdownDescription: "Does not notify the user a script is being executed.",
				Optional:            true,
			},
			"execution_frequency": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The interval for script to run in ISO 8601 duration format (e.g., PT1H for 1 hour, P1D for 1 day). If not defined the script will run once.",
			},
			"retry_count": schema.Int32Attribute{
				MarkdownDescription: "Number of times for the script to be retried if it fails.",
				Optional:            true,
			},
			"assignments": commonschema.ScriptAssignmentsSchema(),
			"timeouts":    commonschema.Timeouts(ctx),
		},
	}
}

func (d *MacOSPlatformScriptDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.client = common.SetGraphBetaClientForDataSource(ctx, req, resp, d.TypeName)
}
