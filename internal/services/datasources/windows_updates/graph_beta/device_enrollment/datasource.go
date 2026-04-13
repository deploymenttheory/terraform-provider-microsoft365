package graphBetaWindowsUpdatesDeviceEnrollment

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

const (
	DataSourceName = "microsoft365_graph_beta_windows_updates_device_enrollment"
	ReadTimeout    = 180
)

var (
	_ datasource.DataSource              = &DeviceEnrollmentDataSource{}
	_ datasource.DataSourceWithConfigure = &DeviceEnrollmentDataSource{}
)

func NewDeviceEnrollmentDataSource() datasource.DataSource {
	return &DeviceEnrollmentDataSource{
		ReadPermissions: []string{
			"WindowsUpdates.Read.All",
		},
	}
}

type DeviceEnrollmentDataSource struct {
	client          *msgraphbetasdk.GraphServiceClient
	ReadPermissions []string
}

func (d *DeviceEnrollmentDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = DataSourceName
}

func (d *DeviceEnrollmentDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.client = client.SetGraphBetaClientForDataSource(ctx, req, resp, DataSourceName)
}

func (d *DeviceEnrollmentDataSource) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Retrieves Windows Autopatch enrollment status for Azure AD devices using the `/admin/windows/updates/updatableAssets` endpoint. " +
			"This data source supports multiple lookup methods: by Entra device ID, by device name, list all enrolled devices, or use custom OData queries for advanced filtering.",
		Attributes: map[string]schema.Attribute{
			"entra_device_id": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The Entra ID (Azure AD) device object ID to query enrollment status for. One of `entra_device_id`, `device_name`, or `list_all` must be specified.",
				Validators: []validator.String{
					stringvalidator.AtLeastOneOf(
						path.MatchRoot("device_name"),
						path.MatchRoot("list_all"),
					),
					stringvalidator.ConflictsWith(
						path.MatchRoot("device_name"),
						path.MatchRoot("list_all"),
					),
				},
			},
			"device_name": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The device name to search for. The data source will resolve the name to an Entra device ID and then fetch enrollment status. One of `entra_device_id`, `device_name`, or `list_all` must be specified.",
				Validators: []validator.String{
					stringvalidator.AtLeastOneOf(
						path.MatchRoot("entra_device_id"),
						path.MatchRoot("list_all"),
					),
					stringvalidator.ConflictsWith(
						path.MatchRoot("entra_device_id"),
						path.MatchRoot("list_all"),
					),
				},
			},
			"list_all": schema.BoolAttribute{
				Optional: true,
				MarkdownDescription: "Set to `true` to list all enrolled devices. " +
					"Cannot be combined with `entra_device_id` or `device_name`. " +
					"When using this option, the data source returns all devices in the `devices` attribute. " +
					"Use `update_category` or `odata_filter` to narrow results.",
				Validators: []validator.Bool{
					// AtLeastOneOf is handled by string validators above
				},
			},
			"update_category": schema.StringAttribute{
				Optional: true,
				MarkdownDescription: "Optional filter to only return devices enrolled in a specific update category. " +
					"Valid values: `feature`, `quality`, `driver`. " +
					"Can be used with any lookup method to filter results.",
				Validators: []validator.String{
					stringvalidator.OneOf("feature", "quality", "driver"),
				},
			},
			"odata_filter": schema.StringAttribute{
				Optional: true,
				MarkdownDescription: "Custom OData filter query for advanced filtering when using `list_all`. " +
					"Example: `id eq '12345678-1234-1234-1234-123456789012'`. " +
					"Only applicable when `list_all` is `true`.",
			},
			"devices": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "List of enrolled devices with their enrollment status. When querying by `entra_device_id` or `device_name`, this will contain a single device. When using `list_all`, this may contain multiple devices.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The Entra ID (Azure AD) device object ID.",
						},
						"enrollments": schema.ListNestedAttribute{
							Computed:            true,
							MarkdownDescription: "List of update management enrollments for this device.",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"update_category": schema.StringAttribute{
										Computed:            true,
										MarkdownDescription: "The update category the device is enrolled in (feature, quality, driver).",
									},
								},
							},
						},
						"errors": schema.ListNestedAttribute{
							Computed:            true,
							MarkdownDescription: "List of errors associated with this device's enrollment.",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"error_code": schema.StringAttribute{
										Computed:            true,
										MarkdownDescription: "The error code.",
									},
									"error_message": schema.StringAttribute{
										Computed:            true,
										MarkdownDescription: "The error message.",
									},
								},
							},
						},
					},
				},
			},
			"timeouts": commonschema.DatasourceTimeouts(ctx),
		},
	}
}
