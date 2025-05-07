package graphBetaLinuxPlatformScript

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/schema"
	commonschemagraphbeta "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/schema/graph_beta/device_and_app_management"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

const (
	DataSourceName = "graph_beta_device_and_app_management_linux_platform_script"
	ReadTimeout    = 180
)

var (
	// Basic data source interface
	_ datasource.DataSource = &LinuxPlatformScriptDataSource{}

	// Allows the data source to be configured with the provider client
	_ datasource.DataSourceWithConfigure = &LinuxPlatformScriptDataSource{}
)

func NewLinuxPlatformScriptDataSource() datasource.DataSource {
	return &LinuxPlatformScriptDataSource{
		ReadPermissions: []string{
			"DeviceManagementConfiguration.Read.All",
		},
	}
}

type LinuxPlatformScriptDataSource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ProviderTypeName string
	TypeName         string
	ReadPermissions  []string
}

// Metadata returns the data source type name.
func (d *LinuxPlatformScriptDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + DataSourceName
}

// Configure sets the client for the data source.
func (d *LinuxPlatformScriptDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.client = common.SetGraphBetaClientForDataSource(ctx, req, resp, d.TypeName)
}

// Schema defines the data source schema.
func (d *LinuxPlatformScriptDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Gets information about an Intune Linux platform script using the 'configurationPolicies' Graph Beta API.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The unique identifier of the linux platform script.",
			},
			"name": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Name of the linux device management script.",
			},
			"description": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Description for the linux device management script.",
			},
			"script_content": schema.StringAttribute{
				MarkdownDescription: "The linux script content. This will be base64 encoded as part of the request.",
				Computed:            true,
			},
			"role_scope_tag_ids": schema.ListAttribute{
				ElementType:         types.StringType,
				Computed:            true,
				MarkdownDescription: "List of scope tag IDs for this linux device management script.",
			},
			"platforms": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Platform type for this linux platform script. Will always be 'linux'.",
			},
			"technologies": schema.ListAttribute{
				ElementType:         types.StringType,
				Computed:            true,
				MarkdownDescription: "Describes the technologies this settings catalog setting can be deployed with. Usually contains 'linuxMdm'.",
			},
			"execution_context": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Execution context for the linux platform script. Can be one of: user or root.",
			},
			"execution_frequency": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Execution frequency for the Linux platform script. Can be one of: `15minutes`, `30minutes`, `1hour`, `2hour`, `3hour`, `6hour`, `12hour`, `1day`, or `1week`.",
			},
			"execution_retries": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Number of times the Linux platform script should be retried on failure. Can be one of: `1`, `2`, or `3`.",
			},
			"assignments": commonschemagraphbeta.ConfigurationPolicyAssignmentsSchema(),
			"timeouts":    commonschema.Timeouts(ctx),
		},
	}
}
