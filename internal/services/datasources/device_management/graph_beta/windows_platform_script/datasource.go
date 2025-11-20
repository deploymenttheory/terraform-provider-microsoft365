package graphBetaDeviceManagementScript

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

const (
	DataSourceName = "microsoft365_graph_beta_device_management_windows_platform_script"
	ReadTimeout  = 180
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
	
	
	ReadPermissions  []string
}

// Metadata returns the resource type name.
func (r *WindowsPlatformScriptDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = DataSourceName
}

// Configure sets the client for the data source
func (d *WindowsPlatformScriptDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.client = client.SetGraphBetaClientForDataSource(ctx, req, resp, DataSourceName)
}

// Schema defines the schema for the data source
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
			"role_scope_tag_ids": schema.SetAttribute{
				MarkdownDescription: "List of Scope Tag IDs for this PowerShellScript instance.",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"timeouts": commonschema.Timeouts(ctx),
		},
	}
}
