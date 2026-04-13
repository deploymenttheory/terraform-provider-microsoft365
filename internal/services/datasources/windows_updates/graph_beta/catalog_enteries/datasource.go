package graphBetaWindowsUpdateCatalog

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
	DataSourceName = "microsoft365_graph_beta_windows_updates_catalog_enteries"
	ReadTimeout    = 180
)

var (
	_ datasource.DataSource              = &WindowsUpdateCatalogEnteriesDataSource{}
	_ datasource.DataSourceWithConfigure = &WindowsUpdateCatalogEnteriesDataSource{}
)

// NewWindowsUpdateCatalogEnteriesDataSource creates a new data source for Windows Update Catalog
func NewWindowsUpdateCatalogEnteriesDataSource() datasource.DataSource {
	return &WindowsUpdateCatalogEnteriesDataSource{
		ReadPermissions: []string{
			"WindowsUpdates.Read.All",
		},
	}
}

// WindowsUpdateCatalogEnteriesDataSource defines the data source implementation
type WindowsUpdateCatalogEnteriesDataSource struct {
	client          *msgraphbetasdk.GraphServiceClient
	ReadPermissions []string
}

// Metadata returns the data source type name
func (d *WindowsUpdateCatalogEnteriesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = DataSourceName
}

// Configure configures the data source with the provider client
func (d *WindowsUpdateCatalogEnteriesDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.client = client.SetGraphBetaClientForDataSource(ctx, req, resp, DataSourceName)
}

// Schema defines the schema for the data source
func (d *WindowsUpdateCatalogEnteriesDataSource) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Retrieves Windows Update catalog entries from Microsoft Graph using the `/admin/windows/updates/catalog/entries` endpoint. This data source returns feature update and quality update catalog entries with deployment information.",
		Attributes: map[string]schema.Attribute{
			"filter_type": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Type of filter to apply. Valid values are: `all`, `id`, `display_name`, `catalog_entry_type`.",
				Validators: []validator.String{
					stringvalidator.OneOf("all", "id", "display_name", "catalog_entry_type"),
				},
			},
			"filter_value": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Value to filter by. Not required when filter_type is 'all'. For catalog_entry_type, use 'featureUpdate' or 'qualityUpdate'.",
			},
			"entries": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "The list of Windows Update Catalog Entries that match the filter criteria.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The unique identifier for the catalog entry.",
						},
						"display_name": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The display name of the catalog entry.",
						},
						"release_date_time": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The release date and time of the catalog entry in RFC3339 format.",
						},
						"deployable_until_date_time": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The date and time until which the update can be deployed, in RFC3339 format. Null if no expiration.",
						},
						"catalog_entry_type": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The type of catalog entry. Values: 'featureUpdate' or 'qualityUpdate'.",
						},
						"version": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The version of the feature update (feature updates only).",
						},
						"catalog_name": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The catalog name of the quality update (quality updates only).",
						},
						"short_name": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The short name of the quality update (quality updates only).",
						},
						"is_expeditable": schema.BoolAttribute{
							Computed:            true,
							MarkdownDescription: "Indicates whether the quality update can be expedited (quality updates only).",
						},
						"quality_update_classification": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The classification of the quality update, e.g., 'security' (quality updates only).",
						},
						"quality_update_cadence": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The release cadence of the quality update, e.g., 'monthly' (quality updates only).",
						},
						"cve_severity_information": schema.SingleNestedAttribute{
							Computed:            true,
							MarkdownDescription: "CVE severity information for the quality update (quality updates only).",
							Attributes: map[string]schema.Attribute{
								"max_severity": schema.StringAttribute{
									Computed:            true,
									MarkdownDescription: "The maximum severity level of CVEs, e.g., 'critical'.",
								},
								"max_base_score": schema.Float64Attribute{
									Computed:            true,
									MarkdownDescription: "The maximum CVSS base score.",
								},
								"exploited_cves": schema.ListNestedAttribute{
									Computed:            true,
									MarkdownDescription: "List of exploited CVEs.",
									NestedObject: schema.NestedAttributeObject{
										Attributes: map[string]schema.Attribute{
											"number": schema.StringAttribute{
												Computed:            true,
												MarkdownDescription: "The CVE number, e.g., 'CVE-2023-32046'.",
											},
											"url": schema.StringAttribute{
												Computed:            true,
												MarkdownDescription: "The URL to the CVE details.",
											},
										},
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
