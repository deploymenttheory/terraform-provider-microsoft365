package graphBetaSubscribedSkus

import (
	"context"
	"regexp"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/datasources/common"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/datasources/common/schema"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

const (
	DataSourceName = "graph_beta_subscribed_skus"
	ReadTimeout    = 30
)

var (
	// Ensure the implementation satisfies the expected interfaces
	_ datasource.DataSource              = &SubscribedSkusDataSource{}
	_ datasource.DataSourceWithConfigure = &SubscribedSkusDataSource{}

	// Compiled regex for UUID validation
	uuidRegex = regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)
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
	client           *msgraphbetasdk.GraphServiceClient
	ProviderTypeName string
	TypeName         string
	ReadPermissions  []string
	ResourcePath     string
}

// Metadata returns the data source type name.
func (d *SubscribedSkusDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + DataSourceName
}

// Configure sets the client for the data source.
func (d *SubscribedSkusDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.client = common.SetGraphBetaClientForDataSource(ctx, req, resp, d.TypeName)
}

// Schema returns the schema for the data source.
func (d *SubscribedSkusDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Retrieves information about Microsoft 365 license SKUs (Stock Keeping Units) that an organization has subscribed to. " +
			"This data source provides details about available licenses, their consumption, and service plans.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The unique identifier for this data source operation.",
			},
			"sku_id": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Filter results by a specific SKU ID (GUID). When specified, only the matching SKU will be returned.",
				Validators: []validator.String{
					stringvalidator.RegexMatches(uuidRegex, "Must be a valid UUID format (xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx)"),
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
						"consumed_units": schema.Int64Attribute{
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
								"enabled": schema.Int64Attribute{
									Computed:            true,
									MarkdownDescription: "The number of units that are enabled.",
								},
								"locked_out": schema.Int64Attribute{
									Computed:            true,
									MarkdownDescription: "The number of units that are locked out.",
								},
								"suspended": schema.Int64Attribute{
									Computed:            true,
									MarkdownDescription: "The number of units that are suspended.",
								},
								"warning": schema.Int64Attribute{
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
			"timeouts": commonschema.Timeouts(ctx),
		},
	}
}

// Helper functions for object type definitions
func getSubscribedSkuObjectType() map[string]attr.Type {
	return map[string]attr.Type{
		"id":                types.StringType,
		"account_id":        types.StringType,
		"account_name":      types.StringType,
		"applies_to":        types.StringType,
		"capability_status": types.StringType,
		"consumed_units":    types.Int64Type,
		"sku_id":            types.StringType,
		"sku_part_number":   types.StringType,
		"prepaid_units":     types.ObjectType{AttrTypes: getPrepaidUnitsObjectType()},
		"service_plans":     types.ListType{ElemType: types.ObjectType{AttrTypes: getServicePlanObjectType()}},
		"subscription_ids":  types.ListType{ElemType: types.StringType},
	}
}

func getPrepaidUnitsObjectType() map[string]attr.Type {
	return map[string]attr.Type{
		"enabled":    types.Int64Type,
		"locked_out": types.Int64Type,
		"suspended":  types.Int64Type,
		"warning":    types.Int64Type,
	}
}

func getServicePlanObjectType() map[string]attr.Type {
	return map[string]attr.Type{
		"service_plan_id":     types.StringType,
		"service_plan_name":   types.StringType,
		"provisioning_status": types.StringType,
		"applies_to":          types.StringType,
	}
}
