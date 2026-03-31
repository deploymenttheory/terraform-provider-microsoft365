package graphBetaWindowsUpdatesApplicableContent

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
	DataSourceName = "microsoft365_graph_beta_windows_updates_applicable_content"
	ReadTimeout    = 180
)

var (
	_ datasource.DataSource              = &ApplicableContentDataSource{}
	_ datasource.DataSourceWithConfigure = &ApplicableContentDataSource{}
)

func NewApplicableContentDataSource() datasource.DataSource {
	return &ApplicableContentDataSource{
		ReadPermissions: []string{
			"WindowsUpdates.ReadWrite.All",
		},
	}
}

type ApplicableContentDataSource struct {
	client          *msgraphbetasdk.GraphServiceClient
	ReadPermissions []string
}

func (d *ApplicableContentDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = DataSourceName
}

func (d *ApplicableContentDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.client = client.SetGraphBetaClientForDataSource(ctx, req, resp, DataSourceName)
}

func (d *ApplicableContentDataSource) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Retrieves applicable content (driver and firmware updates) for a deployment audience using the `/admin/windows/updates/deploymentAudiences/{audienceId}/applicableContent` endpoint. " +
			"This data source shows which updates are applicable to devices in a deployment audience, along with which devices match each update. " +
			"Supports filtering by catalog entry type, driver class, manufacturer, and specific devices.",
		Attributes: map[string]schema.Attribute{
			"audience_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The ID of the deployment audience to query for applicable content. This is required as applicable content is scoped to a specific audience.",
			},
			"catalog_entry_type": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Optional filter to only return content of a specific type. Valid values: `driver`, `quality`, `feature`. When specified, only catalog entries of this type will be returned.",
			},
			"driver_class": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Optional filter to only return driver updates of a specific class. Examples: `Display`, `Network`, `Storage`, `Audio`, `Bluetooth`. Only applicable when filtering for driver updates.",
			},
			"manufacturer": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Optional filter to only return updates from a specific manufacturer. Examples: `Intel`, `NVIDIA`, `AMD`, `Microsoft`, `Realtek`.",
			},
			"device_id": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Optional Azure AD device ID to filter results. When specified, only shows applicable content that matches this specific device.",
			},
			"include_no_matches": schema.BoolAttribute{
				Optional:            true,
				MarkdownDescription: "Whether to include content with no matched devices. Defaults to `true`. Set to `false` to only return content that has at least one matched device.",
			},
			"odata_filter": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Optional custom OData filter expression for advanced filtering. This is applied client-side after retrieving the applicable content. Example: `catalogEntry/displayName contains 'NVIDIA'`.",
			},
			"applicable_content": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "List of applicable content entries (drivers/firmware) for the audience.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"catalog_entry_id": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The ID of the catalog entry for this applicable content.",
						},
						"catalog_entry": schema.SingleNestedAttribute{
							Computed:            true,
							MarkdownDescription: "Details about the driver update catalog entry.",
							Attributes: map[string]schema.Attribute{
								"id": schema.StringAttribute{
									Computed:            true,
									MarkdownDescription: "The unique identifier for the driver update catalog entry.",
								},
								"display_name": schema.StringAttribute{
									Computed:            true,
									MarkdownDescription: "The display name of the driver update.",
								},
								"release_date_time": schema.StringAttribute{
									Computed:            true,
									MarkdownDescription: "The release date and time in RFC3339 format.",
								},
								"deployable_until_date_time": schema.StringAttribute{
									Computed:            true,
									MarkdownDescription: "The date and time until which the driver can be deployed, in RFC3339 format.",
								},
								"description": schema.StringAttribute{
									Computed:            true,
									MarkdownDescription: "Description of the driver update.",
								},
								"driver_class": schema.StringAttribute{
									Computed:            true,
									MarkdownDescription: "The class of the driver, e.g., 'Display', 'Network'.",
								},
								"provider": schema.StringAttribute{
									Computed:            true,
									MarkdownDescription: "The provider of the driver update.",
								},
								"manufacturer": schema.StringAttribute{
									Computed:            true,
									MarkdownDescription: "The manufacturer of the driver.",
								},
								"version": schema.StringAttribute{
									Computed:            true,
									MarkdownDescription: "The version of the driver.",
								},
								"version_date_time": schema.StringAttribute{
									Computed:            true,
									MarkdownDescription: "The version date and time in RFC3339 format.",
								},
							},
						},
						"matched_devices": schema.ListNestedAttribute{
							Computed:            true,
							MarkdownDescription: "List of devices that match this driver update.",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"device_id": schema.StringAttribute{
										Computed:            true,
										MarkdownDescription: "The Azure AD device ID.",
									},
									"recommended_by": schema.ListAttribute{
										ElementType:         types.StringType,
										Computed:            true,
										MarkdownDescription: "List of entities recommending this driver, e.g., ['Microsoft', 'Contoso'].",
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
