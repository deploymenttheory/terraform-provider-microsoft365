package graphCloudPcDeviceImage

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/schema"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	msgraphsdk "github.com/microsoftgraph/msgraph-sdk-go"
)

var _ datasource.DataSource = &CloudPcDeviceImageDataSource{}
var _ datasource.DataSourceWithConfigure = &CloudPcDeviceImageDataSource{}

func NewCloudPcDeviceImageDataSource() datasource.DataSource {
	return &CloudPcDeviceImageDataSource{}
}

type CloudPcDeviceImageDataSource struct {
	client           *msgraphsdk.GraphServiceClient
	ProviderTypeName string
	TypeName         string
}

type CloudPcDeviceImageDataSourceModel struct {
	ID                    types.String   `tfsdk:"id"`
	DisplayName           types.String   `tfsdk:"display_name"`
	ErrorCode             types.String   `tfsdk:"error_code"`
	ExpirationDate        types.String   `tfsdk:"expiration_date"`
	LastModifiedDateTime  types.String   `tfsdk:"last_modified_date_time"`
	OperatingSystem       types.String   `tfsdk:"operating_system"`
	OSBuildNumber         types.String   `tfsdk:"os_build_number"`
	OSStatus              types.String   `tfsdk:"os_status"`
	SourceImageResourceID types.String   `tfsdk:"source_image_resource_id"`
	Status                types.String   `tfsdk:"status"`
	Version               types.String   `tfsdk:"version"`
	Timeouts              timeouts.Value `tfsdk:"timeouts"`
}

func (d *CloudPcDeviceImageDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_graph_cloud_pc_device_image"
}

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
			"timeouts": commonschema.Timeouts(ctx),
		},
	}
}

func (d *CloudPcDeviceImageDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.client = common.SetGraphStableClientForDataSource(ctx, req, resp, "CloudPcDeviceImageDataSource")
}
