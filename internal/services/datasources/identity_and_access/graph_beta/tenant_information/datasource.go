package graphBetaTenantInformation

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
	DataSourceName = "microsoft365_graph_beta_identity_and_access_tenant_information"
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
			"filter_type": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Type of filter to apply. Valid values are: `tenant_id`, `domain_name`.",
				Validators: []validator.String{
					stringvalidator.OneOf("tenant_id", "domain_name"),
				},
			},
			"filter_value": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Value to filter by. Should be a valid tenant ID (GUID) when filter_type is 'tenant_id', or a valid domain name when filter_type is 'domain_name'.",
			},
			"tenant_id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The unique identifier for the tenant.",
			},
			"display_name": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The display name of the tenant.",
			},
			"default_domain_name": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The default domain name of the tenant.",
			},
			"federation_brand_name": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The federation brand name of the tenant, if configured.",
			},
			"timeouts": commonschema.ResourceTimeouts(ctx),
		},
	}
}
