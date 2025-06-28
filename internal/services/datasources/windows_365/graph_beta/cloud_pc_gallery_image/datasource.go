package graphBetaCloudPcGalleryImage

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
	DataSourceName = "graph_beta_windows_365_cloud_pc_gallery_image"
	ReadTimeout    = 180
)

var (
	_ datasource.DataSource              = &CloudPcGalleryImageDataSource{}
	_ datasource.DataSourceWithConfigure = &CloudPcGalleryImageDataSource{}
)

func NewCloudPcGalleryImageDataSource() datasource.DataSource {
	return &CloudPcGalleryImageDataSource{
		ReadPermissions: []string{
			"CloudPC.Read.All",
		},
	}
}

type CloudPcGalleryImageDataSource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ProviderTypeName string
	TypeName         string
	ReadPermissions  []string
}

func (d *CloudPcGalleryImageDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + DataSourceName
}

func (d *CloudPcGalleryImageDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.client = client.SetGraphBetaClientForDataSource(ctx, req, resp, d.TypeName)
}

func (d *CloudPcGalleryImageDataSource) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Retrieves Cloud PC Gallery Images from Microsoft Intune. Using the endpoint '/deviceManagement/virtualEndpoint/galleryImages'. Supports filtering by all, id, or display_name.",
		Attributes: map[string]schema.Attribute{
			"filter_type": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Type of filter to apply. Valid values are: `all`, `id`, `display_name`. Use 'all' to retrieve all gallery images, 'id' to retrieve a specific image by its unique identifier, or 'display_name' to filter by the image's display name.",
				Validators: []validator.String{
					stringvalidator.OneOf("all", "id", "display_name"),
				},
			},
			"filter_value": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Value to filter by. Not required when filter_type is 'all'. For 'id', provide the gallery image ID. For 'display_name', provide a substring to match against image display names.",
			},
			"items": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "The list of Cloud PC gallery images that match the filter criteria.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The unique identifier (ID) of the gallery image resource on Cloud PC. The ID format is {publisherName_offerName_skuName}.",
						},
						"display_name": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The display name of this gallery image. For example, Windows 11 Enterprise + Microsoft 365 Apps 22H2.",
						},
						"start_date": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The date when the Cloud PC image is available for provisioning new Cloud PCs.",
						},
						"end_date": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The date when the status of image becomes supportedWithWarning.",
						},
						"expiration_date": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The date when the image is no longer available. Users are unable to provision new Cloud PCs if the current time is later than this date.",
						},
						"os_version_number": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The operating system version of this gallery image. For example, 10.0.22000.296.",
						},
						"publisher_name": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The publisher name of this gallery image that is passed to ARM to retrieve the image resource.",
						},
						"offer_name": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The offer name of this gallery image that is passed to ARM to retrieve the image resource.",
						},
						"sku_name": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The SKU name of this image that is passed to ARM to retrieve the image resource.",
						},
						"size_in_gb": schema.Int64Attribute{
							Computed:            true,
							MarkdownDescription: "Indicates the size of this image in gigabytes. For example, 64.",
						},
						"status": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The status of the gallery image on the Cloud PC. Possible values are: supported, supportedWithWarning, notSupported, unknownFutureValue.",
						},
					},
				},
			},
			"timeouts": commonschema.Timeouts(ctx),
		},
	}
}
