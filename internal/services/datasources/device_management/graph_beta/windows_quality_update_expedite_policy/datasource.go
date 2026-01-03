package graphBetaWindowsQualityUpdateExpeditePolicy

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
	DataSourceName = "microsoft365_graph_beta_device_management_windows_quality_update_expedite_policy"
	ReadTimeout    = 180
)

var (
	_ datasource.DataSource              = &WindowsQualityUpdateExpeditePolicyDataSource{}
	_ datasource.DataSourceWithConfigure = &WindowsQualityUpdateExpeditePolicyDataSource{}
)

// NewWindowsQualityUpdateExpeditePolicyDataSource creates a new data source for Windows Quality Update Expedite Policies
func NewWindowsQualityUpdateExpeditePolicyDataSource() datasource.DataSource {
	return &WindowsQualityUpdateExpeditePolicyDataSource{
		ReadPermissions: []string{
			"DeviceManagementConfiguration.Read.All",
		},
	}
}

// WindowsQualityUpdateExpeditePolicyDataSource defines the data source implementation
type WindowsQualityUpdateExpeditePolicyDataSource struct {
	client *msgraphbetasdk.GraphServiceClient

	ReadPermissions []string
}

// Metadata returns the data source type name
func (d *WindowsQualityUpdateExpeditePolicyDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = DataSourceName
}

// Configure configures the data source with the provider client
func (d *WindowsQualityUpdateExpeditePolicyDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.client = client.SetGraphBetaClientForDataSource(ctx, req, resp, DataSourceName)
}

// Schema defines the schema for the data source
func (d *WindowsQualityUpdateExpeditePolicyDataSource) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Retrieves Windows Quality Update Expedite Policies from Microsoft Intune with explicit filtering options. " +
			"These policies control the expedited deployment of quality updates to Windows devices.",
		Attributes: map[string]schema.Attribute{
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
			"items": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "The list of Windows Quality Update Expedite Policies that match the filter criteria.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The ID of the Windows Quality Update Expedite Policy.",
						},
						"display_name": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The display name of the Windows Quality Update Expedite Policy.",
						},
						"description": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The description of the Windows Quality Update Expedite Policy.",
						},
					},
				},
			},
			"timeouts": commonschema.ResourceTimeouts(ctx),
		},
	}
}
