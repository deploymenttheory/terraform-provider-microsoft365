package graphBetaCloudPcSourceDeviceImage

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
	DataSourceName = "microsoft365_graph_beta_windows_365_cloud_pc_source_device_image"
	ReadTimeout    = 180
)

var (
	_ datasource.DataSource              = &CloudPcSourceDeviceImageDataSource{}
	_ datasource.DataSourceWithConfigure = &CloudPcSourceDeviceImageDataSource{}
)

func NewCloudPcSourceDeviceImageDataSource() datasource.DataSource {
	return &CloudPcSourceDeviceImageDataSource{
		ReadPermissions: []string{
			"CloudPC.Read.All",
		},
	}
}

type CloudPcSourceDeviceImageDataSource struct {
	client          *msgraphbetasdk.GraphServiceClient
	ReadPermissions []string
}

func (d *CloudPcSourceDeviceImageDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = DataSourceName
}

func (d *CloudPcSourceDeviceImageDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.client = client.SetGraphBetaClientForDataSource(ctx, req, resp, DataSourceName)
}

func (d *CloudPcSourceDeviceImageDataSource) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Retrieves Cloud PC source device images from Microsoft Intune using the `/deviceManagement/virtualEndpoint/deviceImages/getSourceImages` endpoint. This data source is used to discover available Azure images for uploading and provisioning Cloud PCs with filtering options.",
		Attributes: map[string]schema.Attribute{
			"filter_type": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Type of filter to apply. Valid values are: `all`, `id`, `display_name`. Use 'all' to retrieve all images, 'id' to retrieve a specific image by its unique identifier, or 'display_name' to filter by the image's display name.",
				Validators: []validator.String{
					stringvalidator.OneOf("all", "id", "display_name"),
				},
			},
			"filter_value": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Value to filter by. Not required when filter_type is 'all'. For 'id', provide the image ID. For 'display_name', provide a substring to match against image display names.",
			},
			"items": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "The list of Cloud PC source device images that match the filter criteria.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The unique identifier for the source device image.",
						},
						"resource_id": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The resource ID for the source device image.",
						},
						"display_name": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The display name of the source device image.",
						},
						"subscription_id": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The subscription ID associated with the image.",
						},
						"subscription_display_name": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The display name of the subscription.",
						},
					},
				},
			},
			"timeouts": commonschema.ResourceTimeouts(ctx),
		},
	}
}
