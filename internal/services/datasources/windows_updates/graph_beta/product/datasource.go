package graphBetaWindowsUpdateProduct

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

const (
	DataSourceName = "microsoft365_graph_beta_windows_updates_product"
	ReadTimeout    = 180
)

var (
	_ datasource.DataSource              = &WindowsUpdateProductDataSource{}
	_ datasource.DataSourceWithConfigure = &WindowsUpdateProductDataSource{}
)

func NewWindowsUpdateProductDataSource() datasource.DataSource {
	return &WindowsUpdateProductDataSource{
		ReadPermissions: []string{
			"WindowsUpdates.ReadWrite.All",
		},
	}
}

type WindowsUpdateProductDataSource struct {
	client          *msgraphbetasdk.GraphServiceClient
	ReadPermissions []string
}

func (d *WindowsUpdateProductDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = DataSourceName
}

func (d *WindowsUpdateProductDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.client = client.SetGraphBetaClientForDataSource(ctx, req, resp, DataSourceName)
}

func (d *WindowsUpdateProductDataSource) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Retrieves Windows Update product information from Microsoft Graph. This data source can search by catalog ID or KB number using the `/admin/windows/updates/products/FindByCatalogId` or `/admin/windows/updates/products/FindByKbNumber` endpoints.",
		Attributes: map[string]schema.Attribute{
			"search_type": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Type of search to perform. Valid values are: `catalog_id`, `kb_number`.",
				Validators: []validator.String{
					stringvalidator.OneOf("catalog_id", "kb_number"),
				},
			},
			"search_value": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Value to search by. For catalog_id, provide the catalog identifier. For kb_number, provide the KB article number (e.g., '5029332').",
			},
			"products": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "The list of Windows Update products that match the search criteria.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The unique identifier for the product.",
						},
						"name": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The name of the product, e.g., 'Windows 11, version 22H2'.",
						},
						"group_name": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The name of the product group, e.g., 'Windows 11'.",
						},
						"friendly_names": schema.ListAttribute{
							ElementType:         types.StringType,
							Computed:            true,
							MarkdownDescription: "The friendly names of the product, e.g., 'Version 22H2 (OS build 22621)'.",
						},
						"revisions": schema.ListNestedAttribute{
							Computed:            true,
							MarkdownDescription: "Product revisions associated with the search criteria.",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"id": schema.StringAttribute{
										Computed:            true,
										MarkdownDescription: "The unique identifier for the product revision, e.g., '10.0.22621.2215'.",
									},
									"display_name": schema.StringAttribute{
										Computed:            true,
										MarkdownDescription: "The display name of the product revision.",
									},
									"release_date_time": schema.StringAttribute{
										Computed:            true,
										MarkdownDescription: "The release date and time in RFC3339 format.",
									},
									"version": schema.StringAttribute{
										Computed:            true,
										MarkdownDescription: "The version of the product revision, e.g., '22H2'.",
									},
									"os_build": schema.SingleNestedAttribute{
										Computed:            true,
										MarkdownDescription: "The OS build information.",
										Attributes: map[string]schema.Attribute{
											"major_version": schema.Int32Attribute{
												Computed:            true,
												MarkdownDescription: "The major version number.",
											},
											"minor_version": schema.Int32Attribute{
												Computed:            true,
												MarkdownDescription: "The minor version number.",
											},
											"build_number": schema.Int32Attribute{
												Computed:            true,
												MarkdownDescription: "The build number.",
											},
											"update_build_revision": schema.Int64Attribute{
												Computed:            true,
												MarkdownDescription: "The update build revision number.",
											},
										},
									},
									"catalog_entry": schema.SingleNestedAttribute{
										Computed:            true,
										MarkdownDescription: "The catalog entry associated with this revision.",
										Attributes: map[string]schema.Attribute{
											"id": schema.StringAttribute{
												Computed:            true,
												MarkdownDescription: "The catalog entry identifier.",
											},
											"display_name": schema.StringAttribute{
												Computed:            true,
												MarkdownDescription: "The display name of the catalog entry.",
											},
											"release_date_time": schema.StringAttribute{
												Computed:            true,
												MarkdownDescription: "The release date and time in RFC3339 format.",
											},
											"deployable_until_date_time": schema.StringAttribute{
												Computed:            true,
												MarkdownDescription: "The date and time until which the update can be deployed, in RFC3339 format.",
											},
											"catalog_name": schema.StringAttribute{
												Computed:            true,
												MarkdownDescription: "The catalog name.",
											},
											"is_expeditable": schema.BoolAttribute{
												Computed:            true,
												MarkdownDescription: "Indicates whether the update can be expedited.",
											},
											"quality_update_classification": schema.StringAttribute{
												Computed:            true,
												MarkdownDescription: "The classification of the quality update.",
											},
											"quality_update_cadence": schema.StringAttribute{
												Computed:            true,
												MarkdownDescription: "The release cadence of the quality update.",
											},
											"short_name": schema.StringAttribute{
												Computed:            true,
												MarkdownDescription: "The short name of the update.",
											},
										},
									},
									"knowledge_base_article": schema.SingleNestedAttribute{
										Computed:            true,
										MarkdownDescription: "The knowledge base article associated with this revision.",
										Attributes: map[string]schema.Attribute{
											"id": schema.StringAttribute{
												Computed:            true,
												MarkdownDescription: "The KB article ID, e.g., 'KB5029351'.",
											},
											"url": schema.StringAttribute{
												Computed:            true,
												MarkdownDescription: "The URL to the KB article.",
											},
										},
									},
								},
							},
						},
						"known_issues": schema.ListNestedAttribute{
							Computed:            true,
							MarkdownDescription: "Known issues related to the product.",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"id": schema.StringAttribute{
										Computed:            true,
										MarkdownDescription: "The unique identifier for the known issue.",
									},
									"title": schema.StringAttribute{
										Computed:            true,
										MarkdownDescription: "The title of the known issue.",
									},
									"description": schema.StringAttribute{
										Computed:            true,
										MarkdownDescription: "The description of the known issue.",
									},
									"status": schema.StringAttribute{
										Computed:            true,
										MarkdownDescription: "The status of the known issue.",
									},
									"web_view_url": schema.StringAttribute{
										Computed:            true,
										MarkdownDescription: "The URL to view the known issue in the admin portal.",
									},
									"start_date_time": schema.StringAttribute{
										Computed:            true,
										MarkdownDescription: "The start date and time of the known issue in RFC3339 format.",
									},
									"resolved_date_time": schema.StringAttribute{
										Computed:            true,
										MarkdownDescription: "The resolved date and time of the known issue in RFC3339 format.",
									},
									"last_updated_date_time": schema.StringAttribute{
										Computed:            true,
										MarkdownDescription: "The last updated date and time of the known issue in RFC3339 format.",
									},
									"originating_knowledge_base_article": schema.SingleNestedAttribute{
										Computed:            true,
										MarkdownDescription: "The KB article that originated the known issue.",
										Attributes: map[string]schema.Attribute{
											"id": schema.StringAttribute{
												Computed:            true,
												MarkdownDescription: "The KB article ID.",
											},
											"url": schema.StringAttribute{
												Computed:            true,
												MarkdownDescription: "The URL to the KB article.",
											},
										},
									},
									"resolving_knowledge_base_article": schema.SingleNestedAttribute{
										Computed:            true,
										MarkdownDescription: "The KB article that resolved the known issue.",
										Attributes: map[string]schema.Attribute{
											"id": schema.StringAttribute{
												Computed:            true,
												MarkdownDescription: "The KB article ID.",
											},
											"url": schema.StringAttribute{
												Computed:            true,
												MarkdownDescription: "The URL to the KB article.",
											},
										},
									},
									"safeguard_hold_ids": schema.ListAttribute{
										ElementType:         types.StringType,
										Computed:            true,
										MarkdownDescription: "List of safeguard hold IDs associated with the known issue.",
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
