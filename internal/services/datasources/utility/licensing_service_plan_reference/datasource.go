package utilityLicensingServicePlanReference

import (
	"context"
	"regexp"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

const (
	DataSourceName = "microsoft365_utility_licensing_service_plan_reference"
	ReadTimeout    = 180
)

var (
	// Basic datasource interface (Read operations)
	_ datasource.DataSource = &licensingServicePlanReferenceDataSource{}

	// Allows the datasource to be configured with the provider client
	_ datasource.DataSourceWithConfigure = &licensingServicePlanReferenceDataSource{}
)

func NewlicensingServicePlanReferenceDataSource() datasource.DataSource {
	return &licensingServicePlanReferenceDataSource{}
}

type licensingServicePlanReferenceDataSource struct{}

func (d *licensingServicePlanReferenceDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = DataSourceName
}

// Configure implements the DataSourceWithConfigure interface.
// This utility datasource performs local lookups against embedded licensing data,
// so it doesn't require Microsoft Graph clients. However, this method is still
// required for interface compliance and maintains consistency with other datasources.
func (d *licensingServicePlanReferenceDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
}

func (d *licensingServicePlanReferenceDataSource) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Queries Microsoft 365 licensing service plan reference data. " +
			"This utility data source allows you to search for license products (SKUs) and service plans using human-readable names, " +
			"GUIDs, or string IDs. The data is sourced from Microsoft's official licensing service plan reference documentation.\n\n" +
			"**Search Modes:**\n\n" +
			"- **By Product Name**: Use `product_name` to search for license products (e.g., \"Microsoft 365 E3\")\n" +
			"- **By Product Identifier**: Use `string_id` or `guid` to look up a specific product\n" +
			"- **By Service Plan Name**: Use `service_plan_name` to find service plans (e.g., \"Exchange Online\")\n" +
			"- **By Service Plan Identifier**: Use `service_plan_id` or `service_plan_guid` for specific service plans\n\n" +
			"Only one search parameter should be specified at a time. Results include both the matching items and their relationships " +
			"(e.g., which products include a specific service plan, or which service plans are included in a product).\n\n" +
			"**Reference:** [Microsoft Licensing Service Plan Reference](https://learn.microsoft.com/en-us/entra/identity/users/licensing-service-plan-reference)",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The ID of this data source operation.",
			},
			"product_name": schema.StringAttribute{
				Optional: true,
				MarkdownDescription: "Search for products by name (case-insensitive partial match). " +
					"Returns all products whose names contain this string. " +
					"Example: `\"Microsoft 365 E3\"` or `\"Office 365\"`.",
				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(
						path.MatchRoot("string_id"),
						path.MatchRoot("guid"),
						path.MatchRoot("service_plan_id"),
						path.MatchRoot("service_plan_name"),
						path.MatchRoot("service_plan_guid"),
					),
				},
			},
			"string_id": schema.StringAttribute{
				Optional: true,
				MarkdownDescription: "Look up a product by its String ID (exact match, case-insensitive). " +
					"String IDs are used by PowerShell v1.0 and the `skuPartNumber` property in Microsoft Graph. " +
					"Example: `\"ENTERPRISEPACK\"`, `\"SPE_E3\"`.",
				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(
						path.MatchRoot("product_name"),
						path.MatchRoot("guid"),
						path.MatchRoot("service_plan_id"),
						path.MatchRoot("service_plan_name"),
						path.MatchRoot("service_plan_guid"),
					),
				},
			},
			"guid": schema.StringAttribute{
				Optional: true,
				MarkdownDescription: "Look up a product by its GUID (exact match). " +
					"GUIDs are used by the `skuId` property in Microsoft Graph. " +
					"Example: `\"6fd2c87f-b296-42f0-b197-1e91e994b900\"`.",
				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(
						path.MatchRoot("product_name"),
						path.MatchRoot("string_id"),
						path.MatchRoot("service_plan_id"),
						path.MatchRoot("service_plan_name"),
						path.MatchRoot("service_plan_guid"),
					),
					stringvalidator.RegexMatches(
						regexp.MustCompile(constants.GuidRegex),
						"must be a valid GUID in the format 00000000-0000-0000-0000-000000000000",
					),
				},
			},
			"service_plan_id": schema.StringAttribute{
				Optional: true,
				MarkdownDescription: "Search for service plans by ID (case-insensitive partial match). " +
					"Returns all service plans whose IDs contain this string. " +
					"Example: `\"EXCHANGE\"`, `\"TEAMS\"`.",
				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(
						path.MatchRoot("product_name"),
						path.MatchRoot("string_id"),
						path.MatchRoot("guid"),
						path.MatchRoot("service_plan_name"),
						path.MatchRoot("service_plan_guid"),
					),
				},
			},
			"service_plan_name": schema.StringAttribute{
				Optional: true,
				MarkdownDescription: "Search for service plans by friendly name (case-insensitive partial match). " +
					"Returns all service plans whose names contain this string. " +
					"Example: `\"Exchange Online\"`, `\"Microsoft Teams\"`.",
				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(
						path.MatchRoot("product_name"),
						path.MatchRoot("string_id"),
						path.MatchRoot("guid"),
						path.MatchRoot("service_plan_id"),
						path.MatchRoot("service_plan_guid"),
					),
				},
			},
			"service_plan_guid": schema.StringAttribute{
				Optional: true,
				MarkdownDescription: "Look up a service plan by its GUID (exact match). " +
					"Example: `\"113feb6c-3fe4-4440-bddc-54d774bf0318\"`.",
				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(
						path.MatchRoot("product_name"),
						path.MatchRoot("string_id"),
						path.MatchRoot("guid"),
						path.MatchRoot("service_plan_id"),
						path.MatchRoot("service_plan_name"),
					),
					stringvalidator.RegexMatches(
						regexp.MustCompile(constants.GuidRegex),
						"must be a valid GUID in the format 00000000-0000-0000-0000-000000000000",
					),
				},
			},
			"matching_products": schema.ListNestedAttribute{
				Computed: true,
				MarkdownDescription: "List of products matching the search criteria. " +
					"Populated when searching by product name, string_id, guid, or when searching for service plans " +
					"(returns products that include the matching service plans).",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"product_name": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The product name as displayed in management portals (e.g., \"Microsoft 365 E3\").",
						},
						"string_id": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The product String ID used by PowerShell v1.0 and skuPartNumber in Graph API.",
						},
						"guid": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The product GUID used by skuId in Graph API.",
						},
						"service_plans_included": schema.ListNestedAttribute{
							Computed:            true,
							MarkdownDescription: "Service plans included in this product.",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"id": schema.StringAttribute{
										Computed:            true,
										MarkdownDescription: "The service plan ID.",
									},
									"name": schema.StringAttribute{
										Computed:            true,
										MarkdownDescription: "The service plan friendly name.",
									},
									"guid": schema.StringAttribute{
										Computed:            true,
										MarkdownDescription: "The service plan GUID.",
									},
								},
							},
						},
					},
				},
			},
			"matching_service_plans": schema.ListNestedAttribute{
				Computed: true,
				MarkdownDescription: "List of service plans matching the search criteria. " +
					"Populated when searching by service plan name, id, or guid. " +
					"Each entry includes a list of products (SKUs) that include this service plan.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The service plan ID.",
						},
						"name": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The service plan friendly name.",
						},
						"guid": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The service plan GUID.",
						},
						"included_in_skus": schema.ListNestedAttribute{
							Computed:            true,
							MarkdownDescription: "List of products (SKUs) that include this service plan.",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"product_name": schema.StringAttribute{
										Computed:            true,
										MarkdownDescription: "The product name.",
									},
									"string_id": schema.StringAttribute{
										Computed:            true,
										MarkdownDescription: "The product String ID.",
									},
									"guid": schema.StringAttribute{
										Computed:            true,
										MarkdownDescription: "The product GUID.",
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
