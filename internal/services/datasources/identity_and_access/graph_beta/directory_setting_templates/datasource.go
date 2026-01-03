package graphBetaDirectorySettingTemplates

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

const (
	DataSourceName = "microsoft365_graph_beta_identity_and_access_directory_setting_templates"
	ReadTimeout    = 180
)

var (
	// Ensure the implementation satisfies the expected interfaces
	_ datasource.DataSource = &DirectorySettingTemplatesDataSource{}

	// Allows the resource to be configured with the provider client
	_ datasource.DataSourceWithConfigure = &DirectorySettingTemplatesDataSource{}
)

func NewDirectorySettingTemplatesDataSource() datasource.DataSource {
	return &DirectorySettingTemplatesDataSource{
		ReadPermissions: []string{
			"Directory.Read.All",
			"Directory.ReadWrite.All",
		},
		ResourcePath: "/directorySettingTemplates",
	}
}

type DirectorySettingTemplatesDataSource struct {
	client          *msgraphbetasdk.GraphServiceClient
	ReadPermissions []string
	ResourcePath    string
}

// Metadata returns the data source type name.
func (d *DirectorySettingTemplatesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = DataSourceName
}

// Configure sets the client for the data source.
func (d *DirectorySettingTemplatesDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.client = client.SetGraphBetaClientForDataSource(ctx, req, resp, DataSourceName)
}

// Schema returns the schema for the data source.
func (d *DirectorySettingTemplatesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Retrieves information about directory setting templates available in Microsoft 365. " +
			"Directory setting templates represent system-defined settings available to the tenant. " +
			"Directory settings can be created based on the available templates, and values changed from their preset defaults. " +
			"These templates are read-only and define the allowed behaviors for specific Microsoft Entra objects like Microsoft 365 groups.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The unique identifier for this data source operation.",
			},
			"filter_type": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Type of filter to apply. Valid values are: `all`, `id`, `display_name`.",
				Validators: []validator.String{
					stringvalidator.OneOf("all", "id", "display_name"),
				},
			},
			"filter_value": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Value to filter by. Not required when filter_type is 'all'.",
			},
			"directory_setting_templates": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "List of directory setting templates available to the organization.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The unique identifier for the directory setting template.",
						},
						"description": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "Description of the template.",
						},
						"display_name": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "Display name of the template.",
						},
						"values": schema.ListNestedAttribute{
							Computed:            true,
							MarkdownDescription: "Collection of setting template values that list the set of available settings, defaults, and types that make up this template.",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"name": schema.StringAttribute{
										Computed:            true,
										MarkdownDescription: "Name of the setting.",
									},
									"type": schema.StringAttribute{
										Computed:            true,
										MarkdownDescription: "Type of the setting.",
									},
									"default_value": schema.StringAttribute{
										Computed:            true,
										MarkdownDescription: "Default value for the setting.",
									},
									"description": schema.StringAttribute{
										Computed:            true,
										MarkdownDescription: "Description of the setting.",
									},
								},
							},
						},
					},
				},
			},
			"timeouts": commonschema.ResourceTimeouts(ctx),
		},
	}
}
