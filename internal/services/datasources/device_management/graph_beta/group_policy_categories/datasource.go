package graphBetaGroupPolicyCategories

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

const (
	DataSourceName = "microsoft365_graph_beta_device_management_group_policy_category"
	ReadTimeout    = 300
)

var (
	_ datasource.DataSource              = &GroupPolicyCategoryDataSource{}
	_ datasource.DataSourceWithConfigure = &GroupPolicyCategoryDataSource{}
)

func NewGroupPolicyCategoryDataSource() datasource.DataSource {
	return &GroupPolicyCategoryDataSource{
		ReadPermissions: []string{
			"DeviceManagementConfiguration.Read.All",
		},
	}
}

type GroupPolicyCategoryDataSource struct {
	client          *msgraphbetasdk.GraphServiceClient
	ReadPermissions []string
}

func (d *GroupPolicyCategoryDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = DataSourceName
}

func (d *GroupPolicyCategoryDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.client = client.SetGraphBetaClientForDataSource(ctx, req, resp, DataSourceName)
}

func (d *GroupPolicyCategoryDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Retrieves Group Policy categories and settings from Microsoft Intune using the `/deviceManagement/groupPolicyCategories` endpoint. This data source is used to query Group Policy definitions with their categories, presentations, and configuration details for ADMX-backed policies.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The unique identifier for this data source",
			},
			"setting_name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The display name of the Group Policy setting to search for (case-insensitive)",
			},
			"category": schema.SingleNestedAttribute{
				Computed:            true,
				MarkdownDescription: "The Group Policy category information from the first API call",
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Computed:            true,
						MarkdownDescription: "The unique identifier of the category",
					},
					"display_name": schema.StringAttribute{
						Computed:            true,
						MarkdownDescription: "The display name of the category",
					},
					"is_root": schema.BoolAttribute{
						Computed:            true,
						MarkdownDescription: "Indicates if the category is a root category",
					},
					"ingestion_source": schema.StringAttribute{
						Computed:            true,
						MarkdownDescription: "The source of the category (e.g., builtIn, custom)",
					},
					"parent": schema.SingleNestedAttribute{
						Computed:            true,
						MarkdownDescription: "The parent category if this is not a root category",
						Attributes: map[string]schema.Attribute{
							"id": schema.StringAttribute{
								Computed:            true,
								MarkdownDescription: "The unique identifier of the parent category",
							},
							"display_name": schema.StringAttribute{
								Computed:            true,
								MarkdownDescription: "The display name of the parent category",
							},
							"is_root": schema.BoolAttribute{
								Computed:            true,
								MarkdownDescription: "Indicates if the parent category is a root category",
							},
						},
					},
				},
			},
			"definition": schema.SingleNestedAttribute{
				Computed:            true,
				MarkdownDescription: "The detailed Group Policy definition information from the second API call",
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Computed:            true,
						MarkdownDescription: "The unique identifier of the definition",
					},
					"display_name": schema.StringAttribute{
						Computed:            true,
						MarkdownDescription: "The display name of the definition",
					},
					"category_path": schema.StringAttribute{
						Computed:            true,
						MarkdownDescription: "The category path of the definition",
					},
					"class_type": schema.StringAttribute{
						Computed:            true,
						MarkdownDescription: "The class type of the definition (e.g., machine, user)",
					},
					"policy_type": schema.StringAttribute{
						Computed:            true,
						MarkdownDescription: "The policy type of the definition",
					},
					"version": schema.StringAttribute{
						Computed:            true,
						MarkdownDescription: "The version of the definition",
					},
					"has_related_definitions": schema.BoolAttribute{
						Computed:            true,
						MarkdownDescription: "Indicates if the definition has related definitions",
					},
					"explain_text": schema.StringAttribute{
						Computed:            true,
						MarkdownDescription: "The explanation text for the definition",
					},
					"supported_on": schema.StringAttribute{
						Computed:            true,
						MarkdownDescription: "The supported platforms for the definition",
					},
					"group_policy_category_id": schema.StringAttribute{
						Computed:            true,
						MarkdownDescription: "The ID of the group policy category this definition belongs to",
					},
					"min_device_csp_version": schema.StringAttribute{
						Computed:            true,
						MarkdownDescription: "The minimum device CSP version required",
					},
					"min_user_csp_version": schema.StringAttribute{
						Computed:            true,
						MarkdownDescription: "The minimum user CSP version required",
					},
					"last_modified_date_time": schema.StringAttribute{
						Computed:            true,
						MarkdownDescription: "The date and time the definition was last modified",
					},
				},
			},
			"presentations": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "The list of presentations associated with the group policy definition from the third API call",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The ID of the presentation",
						},
						"odata_type": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The OData type of the presentation (e.g., #microsoft.graph.groupPolicyPresentationDropdownList)",
						},
						"label": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The localized text label for the presentation",
						},
						"required": schema.BoolAttribute{
							Computed:            true,
							MarkdownDescription: "Whether a value is required for the parameter box (if applicable)",
						},
						"last_modified_date_time": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The date and time the entity was last modified",
						},
						// For dropdown lists
						"default_item": schema.SingleNestedAttribute{
							Computed:            true,
							MarkdownDescription: "The default item for dropdown list presentations",
							Attributes: map[string]schema.Attribute{
								"display_name": schema.StringAttribute{
									Computed:            true,
									MarkdownDescription: "The display name of the default item",
								},
								"value": schema.StringAttribute{
									Computed:            true,
									MarkdownDescription: "The value of the default item",
								},
							},
						},
						"items": schema.ListNestedAttribute{
							Computed:            true,
							MarkdownDescription: "The list of items for dropdown list presentations",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"display_name": schema.StringAttribute{
										Computed:            true,
										MarkdownDescription: "The display name of the item",
									},
									"value": schema.StringAttribute{
										Computed:            true,
										MarkdownDescription: "The value of the item",
									},
								},
							},
						},
						// For text boxes
						"default_value": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The default value for text box presentations",
						},
						"max_length": schema.Int64Attribute{
							Computed:            true,
							MarkdownDescription: "The maximum length for text box presentations",
						},
						// For checkboxes
						"default_checked": schema.BoolAttribute{
							Computed:            true,
							MarkdownDescription: "Whether the checkbox is checked by default",
						},
						// For decimal text boxes
						"default_decimal_value": schema.Int64Attribute{
							Computed:            true,
							MarkdownDescription: "The default decimal value for decimal text box presentations",
						},
						"min_value": schema.Int64Attribute{
							Computed:            true,
							MarkdownDescription: "The minimum value for decimal text box presentations",
						},
						"max_value": schema.Int64Attribute{
							Computed:            true,
							MarkdownDescription: "The maximum value for decimal text box presentations",
						},
						"spin": schema.BoolAttribute{
							Computed:            true,
							MarkdownDescription: "Whether spin controls are enabled for decimal text box presentations",
						},
						"spin_step": schema.Int64Attribute{
							Computed:            true,
							MarkdownDescription: "The spin step for decimal text box presentations",
						},
						// For list boxes
						"explicit_value": schema.BoolAttribute{
							Computed:            true,
							MarkdownDescription: "Whether the user must specify the registry subkey value and name for list box presentations",
						},
						"value_prefix": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The value prefix for list box presentations",
						},
					},
				},
			},
			"timeouts": commonschema.ResourceTimeouts(ctx),
		},
	}
}
