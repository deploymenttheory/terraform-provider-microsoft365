package graphBetaCloudPcDeviceImages

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
	DataSourceName = "graph_beta_windows_365_cloud_pc_device_image"
	ReadTimeout    = 180
)

var (
	_ datasource.DataSource              = &CloudPcDeviceImageDataSource{}
	_ datasource.DataSourceWithConfigure = &CloudPcDeviceImageDataSource{}
)

func NewCloudPcDeviceImageDataSource() datasource.DataSource {
	return &CloudPcDeviceImageDataSource{
		ReadPermissions: []string{
			"CloudPC.Read.All",
		},
	}
}

type CloudPcDeviceImageDataSource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ProviderTypeName string
	TypeName         string
	ReadPermissions  []string
}

func (d *CloudPcDeviceImageDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + DataSourceName
}

func (d *CloudPcDeviceImageDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.client = client.SetGraphBetaClientForDataSource(ctx, req, resp, d.TypeName)
}

func (d *CloudPcDeviceImageDataSource) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Retrieves Cloud PC Device Images from Microsoft Intune. Using the endpoint '/deviceManagement/virtualEndpoint/deviceImages'. Supports filtering by all, id, or display_name.",
		Attributes: map[string]schema.Attribute{
			"filter_type": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Type of filter to apply. Valid values are: `all`, `id`, `display_name`. Use 'all' to retrieve all Device Images, 'id' to retrieve a specific Device Image by its unique identifier, or 'display_name' to filter by name.",
				Validators: []validator.String{
					stringvalidator.OneOf("all", "id", "display_name"),
				},
			},
			"filter_value": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Value to filter by. Not required when filter_type is 'all'. For 'id', provide the Device Image ID. For 'display_name', provide the name to match.",
			},
			"items": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "The list of Cloud PC Device Images that match the filter criteria.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The unique identifier for the Device Image.",
						},
						"display_name": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The display name of the Device Image.",
						},
						"expiration_date": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The expiration date of the Device Image.",
						},
						"os_build_number": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The OS build number of the Device Image (e.g., '21H2').",
						},
						"os_status": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The OS support status (e.g., 'supported').",
						},
						"operating_system": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The operating system of the Device Image (e.g., 'Windows 10 Enterprise').",
						},
						"version": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The version of the Device Image.",
						},
						"source_image_resource_id": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The Azure resource ID of the source image.",
						},
						"last_modified_date_time": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The date and time when the Device Image was last modified.",
						},
						"status": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The current status of the Device Image (e.g., 'ready').",
						},
						"status_details": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "Additional details about the status of the Device Image.",
						},
						"error_code": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "Error code if there was an issue with the Device Image.",
						},
						"os_version_number": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The OS version number of the Device Image (e.g., '10.0.22631.3593').",
						},
					},
				},
			},
			"timeouts": commonschema.Timeouts(ctx),
		},
	}
}
