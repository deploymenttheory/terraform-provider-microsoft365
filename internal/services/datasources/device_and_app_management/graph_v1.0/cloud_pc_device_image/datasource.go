package graphCloudPcDeviceImage

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	msgraphsdk "github.com/microsoftgraph/msgraph-sdk-go"
)

const (
	DataSourceName = "graph_device_and_app_management_cloud_pc_device_image"
	ReadTimeout    = 180
)

var (
	// Basic resource interface (CRUD operations)
	_ datasource.DataSource = &CloudPcDeviceImageDataSource{}

	// Allows the resource to be configured with the provider client
	_ datasource.DataSourceWithConfigure = &CloudPcDeviceImageDataSource{}
)

func NewCloudPcDeviceImageDataSource() datasource.DataSource {
	return &CloudPcDeviceImageDataSource{
		ReadPermissions: []string{
			"CloudPC.Read.All",
			"CloudPC.ReadWrite.All",
		},
	}
}

type CloudPcDeviceImageDataSource struct {
	client           *msgraphsdk.GraphServiceClient
	ProviderTypeName string
	TypeName         string
	ReadPermissions  []string
}

// Metadata sets the data source name
func (d *CloudPcDeviceImageDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + DataSourceName
}

// Configure sets the client for the data source
func (d *CloudPcDeviceImageDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.client = client.SetGraphStableClientForDataSource(ctx, req, resp, d.TypeName)
}

// Schema defines the schema for the data source
func (d *CloudPcDeviceImageDataSource) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Required:    true,
				Description: "The unique identifier (ID) of the image resource on the Cloud PC.",
			},
			"display_name": schema.StringAttribute{
				Computed:    true,
				Description: "The display name of the associated device image.",
			},
			"error_code": schema.StringAttribute{
				Computed:    true,
				Description: "The error code of the status of the image that indicates why the upload failed, if applicable.",
			},
			"expiration_date": schema.StringAttribute{
				Computed:    true,
				Description: "The date when the image became unavailable.",
			},
			"last_modified_date_time": schema.StringAttribute{
				Computed:    true,
				Description: "The date and time when the image was last modified.",
			},
			"operating_system": schema.StringAttribute{
				Computed:    true,
				Description: "The operating system (OS) of the image.",
			},
			"os_build_number": schema.StringAttribute{
				Computed:    true,
				Description: "The OS build version of the image.",
			},
			"os_status": schema.StringAttribute{
				Computed:    true,
				Description: "The OS status of this image.",
			},
			"source_image_resource_id": schema.StringAttribute{
				Computed:    true,
				Description: "The unique identifier (ID) of the source image resource on Azure.",
			},
			"status": schema.StringAttribute{
				Computed:    true,
				Description: "The status of the image on the Cloud PC.",
			},
			"version": schema.StringAttribute{
				Computed:    true,
				Description: "The image version.",
			},
		},
	}
}
