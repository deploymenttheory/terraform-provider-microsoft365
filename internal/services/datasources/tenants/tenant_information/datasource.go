package graphBetaTenantInformation

import (
	"context"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
)

const (
	DataSourceName = "microsoft365_graph_beta_tenants_tenant_information"
	ReadTimeout    = 180
)

var (
	_ datasource.DataSource              = &TenantInformationDataSource{}
	_ datasource.DataSourceWithConfigure = &TenantInformationDataSource{}
)

func NewTenantInformationDataSource() datasource.DataSource {
	return &TenantInformationDataSource{
		ReadPermissions: []string{
			"CrossTenantInformation.ReadBasic.All",
		},
	}
}

type TenantInformationDataSource struct {
	client          *msgraphbetasdk.GraphServiceClient
	ReadPermissions []string
}

func (d *TenantInformationDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = DataSourceName
}

func (d *TenantInformationDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.client = client.SetGraphBetaClientForDataSource(ctx, req, resp, DataSourceName)
}

func (d *TenantInformationDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Retrieves tenant information from Microsoft Entra ID using the `/tenantRelationships/findTenantInformationByTenantId` or `/tenantRelationships/findTenantInformationByDomainName` endpoint. This data source is used to query tenant details for cross-tenant access configuration and validation.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The unique identifier for this data source operation.",
			},

			"tenant_id": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Tenant UUID to look up. Specify exactly one of tenant_id or domain_name.",
				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(path.MatchRoot("domain_name")),
					stringvalidator.RegexMatches(regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`), "must be a tenant UUID"),
				},
			},
			"domain_name": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Domain to look up. Specify exactly one of tenant_id or domain_name.",
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 253),
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^[A-Za-z0-9](?:[A-Za-z0-9-]{0,61}[A-Za-z0-9])?(?:\.[A-Za-z0-9](?:[A-Za-z0-9-]{0,61}[A-Za-z0-9])?)+$`),
						"must be a domain name",
					),
				},
			},
			"items": schema.ListNestedAttribute{
				Computed: true, MarkdownDescription: "Tenant information returned by the lookup. A successful lookup returns one item.",
				NestedObject: schema.NestedAttributeObject{Attributes: map[string]schema.Attribute{
					"tenant_id":             schema.StringAttribute{Computed: true, MarkdownDescription: "Unique tenant identifier."},
					"display_name":          schema.StringAttribute{Computed: true, MarkdownDescription: "Tenant display name."},
					"default_domain_name":   schema.StringAttribute{Computed: true, MarkdownDescription: "Tenant primary domain."},
					"federation_brand_name": schema.StringAttribute{Computed: true, MarkdownDescription: "Tenant sign-in brand, when configured."},
				}},
			},
			"timeouts": commonschema.ResourceTimeouts(ctx),
		},
	}
}
