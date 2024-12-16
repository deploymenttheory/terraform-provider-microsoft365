package graphbetadevicemanagementscript

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
	_ datasource.DataSource = &WindowsPlatformScriptDataSource{}

	// Allows the resource to be configured with the provider client
	_ datasource.DataSourceWithConfigure = &WindowsPlatformScriptDataSource{}
)

func NewWindowsPlatformScriptDataSource() datasource.DataSource {
	return &WindowsPlatformScriptDataSource{
		ReadPermissions: []string{
			"DeviceManagementConfiguration.Read.All",
		},
	}
}

type WindowsPlatformScriptDataSource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ProviderTypeName string
	TypeName         string
	ReadPermissions  []string
}

// Metadata returns the resource type name.
func (r *WindowsPlatformScriptDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + ResourceName
}

func (d *WindowsPlatformScriptDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Retrieves information about a windows platform script.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Unique identifier for the windows platform script.",
				Optional:            true,
				Computed:            true,
			},
			"display_name": schema.StringAttribute{
				MarkdownDescription: "Name of the windows platform script.",
				Optional:            true,
				Computed:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Description of the windows platform script.",
				Computed:            true,
			},
			"run_as_account": schema.StringAttribute{
				MarkdownDescription: "Indicates the type of execution context.",
				Computed:            true,
			},
			"enforce_signature_check": schema.BoolAttribute{
				MarkdownDescription: "Indicate whether the script signature needs be checked.",
				Computed:            true,
			},
			"file_name": schema.StringAttribute{
				MarkdownDescription: "Script file name.",
				Computed:            true,
			},
			"run_as_32_bit": schema.BoolAttribute{
				MarkdownDescription: "A value indicating whether the PowerShell script should run as 32-bit.",
				Computed:            true,
			},
			"role_scope_tag_ids": schema.ListAttribute{
				MarkdownDescription: "List of Scope Tag IDs for this PowerShellScript instance.",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"script_content": schema.StringAttribute{
				MarkdownDescription: "The script content.",
				Computed:            true,
				Sensitive:           true,
			},
			"assignments": commonschema.ScriptAssignmentsSchema(),
			"timeouts":    commonschema.Timeouts(ctx),
		},
	}
}

func (d *WindowsPlatformScriptDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.client = common.SetGraphBetaClientForDataSource(ctx, req, resp, d.TypeName)
}
