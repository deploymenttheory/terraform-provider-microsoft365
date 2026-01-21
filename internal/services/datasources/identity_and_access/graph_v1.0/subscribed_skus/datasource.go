package graphSubscribedSkus

import (
	"context"
	"regexp"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	msgraphsdk "github.com/microsoftgraph/msgraph-sdk-go"
)

const (
	DataSourceName = "microsoft365_graph_beta_identity_and_access_subscribed_skus"
	ReadTimeout    = 180
)

var (
	// Ensure the implementation satisfies the expected interfaces
	_ datasource.DataSource = &SubscribedSkusDataSource{}

	// Allows the resource to be configured with the provider client
	_ datasource.DataSourceWithConfigure = &SubscribedSkusDataSource{}
)

func NewSubscribedSkusDataSource() datasource.DataSource {
	return &SubscribedSkusDataSource{
		ReadPermissions: []string{
			"LicenseAssignment.Read.All",
			"Directory.Read.All",
			"Organization.Read.All",
		},
		ResourcePath: "/subscribedSkus",
	}
}

type SubscribedSkusDataSource struct {
	client          *msgraphsdk.GraphServiceClient
	ReadPermissions []string
	ResourcePath    string
}

// Metadata returns the data source type name.
func (d *SubscribedSkusDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = DataSourceName
}

// Configure sets the client for the data source.
func (d *SubscribedSkusDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.client = client.SetGraphStableClientForDataSource(ctx, req, resp, DataSourceName)
}

// Schema returns the schema for the data source.
func (d *SubscribedSkusDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Retrieves Microsoft 365 license SKUs from Microsoft Entra ID using the `/subscribedSkus` endpoint. This data source is used to query subscribed license SKUs with consumption details and service plans.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The unique identifier for this data source operation.",
			},
			"sku_id": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Filter results by a specific SKU ID (GUID). When specified, only the matching SKU will be returned.",
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(constants.GuidRegex),
						"must be a valid GUID in the format 00000000-0000-0000-0000-000000000000",
					),
				},
			},
			"sku_part_number": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Filter results by SKU part number (e.g., 'ENTERPRISEPREMIUM', 'AAD_PREMIUM'). When specified, only matching SKUs will be returned.",
			},
			"applies_to": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Filter results by target class. Possible values: 'User', 'Company'. When specified, only SKUs that apply to the specified target will be returned.",
				Validators: []validator.String{
					stringvalidator.OneOf("User", "Company"),
				},
			},
			"subscribed_skus": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "List of subscribed SKUs available to the organization.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The unique identifier for the subscribed SKU object.",
						},
						"account_id": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The unique ID of the account this SKU belongs to.",
						},
						"account_name": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The name of the account this SKU belongs to.",
						},
						"applies_to": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The target class for this SKU. Only SKUs with target class 'User' are assignable. Possible values: 'User', 'Company'.",
						},
						"capability_status": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The status of the SKU. Possible values: 'Enabled', 'Warning', 'Suspended', 'Deleted', 'LockedOut'.",
						},
						"consumed_units": schema.Int32Attribute{
							Computed:            true,
							MarkdownDescription: "The number of licenses that have been assigned.",
						},
						"sku_id": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The unique identifier (GUID) for the service SKU.",
						},
						"sku_part_number": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The SKU part number; for example: 'AAD_PREMIUM' or 'RMSBASIC'.",
						},
						"prepaid_units": schema.SingleNestedAttribute{
							Computed:            true,
							MarkdownDescription: "Information about the number and status of prepaid licenses.",
							Attributes: map[string]schema.Attribute{
								"enabled": schema.Int32Attribute{
									Computed:            true,
									MarkdownDescription: "The number of units that are enabled.",
								},
								"locked_out": schema.Int32Attribute{
									Computed:            true,
									MarkdownDescription: "The number of units that are locked out.",
								},
								"suspended": schema.Int32Attribute{
									Computed:            true,
									MarkdownDescription: "The number of units that are suspended.",
								},
								"warning": schema.Int32Attribute{
									Computed:            true,
									MarkdownDescription: "The number of units that are in warning state.",
								},
							},
						},
						"service_plans": schema.ListNestedAttribute{
							Computed:            true,
							MarkdownDescription: "Information about the service plans that are available with the SKU.",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"service_plan_id": schema.StringAttribute{
										Computed:            true,
										MarkdownDescription: "The unique identifier of the service plan.",
									},
									"service_plan_name": schema.StringAttribute{
										Computed:            true,
										MarkdownDescription: "The name of the service plan.",
									},
									"provisioning_status": schema.StringAttribute{
										Computed:            true,
										MarkdownDescription: "The provisioning status of the service plan.",
									},
									"applies_to": schema.StringAttribute{
										Computed:            true,
										MarkdownDescription: "The object the service plan can be assigned to.",
									},
								},
							},
						},
						"subscription_ids": schema.ListAttribute{
							ElementType:         types.StringType,
							Computed:            true,
							MarkdownDescription: "A list of all subscription IDs associated with this SKU.",
						},
					},
				},
			},
			"timeouts": commonschema.ResourceTimeouts(ctx),
		},
	}
}
