package graphBetaNetworkManagedTLSCertificate

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
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

const (
	DataSourceName = "microsoft365_graph_beta_identity_and_access_network_managed_tls_certificate"
	ReadTimeout    = 180
)

var (
	_ datasource.DataSource              = &NetworkManagedTLSCertificateDataSource{}
	_ datasource.DataSourceWithConfigure = &NetworkManagedTLSCertificateDataSource{}
)

func NewNetworkManagedTLSCertificateDataSource() datasource.DataSource {
	return &NetworkManagedTLSCertificateDataSource{
		ReadPermissions: []string{"NetworkAccess.Read.All"},
	}
}

type NetworkManagedTLSCertificateDataSource struct {
	client          *msgraphbetasdk.GraphServiceClient
	ReadPermissions []string
}

func (d *NetworkManagedTLSCertificateDataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = DataSourceName
}

func (d *NetworkManagedTLSCertificateDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.client = client.SetGraphBetaClientForDataSource(ctx, req, resp, DataSourceName)
}

func (d *NetworkManagedTLSCertificateDataSource) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Retrieves a Microsoft-managed certificate authority and its PEM-encoded root certificate for Microsoft Entra Global Secure Access TLS inspection using Microsoft Graph beta `/networkaccess/tls/managedCertificateAuthorityCertificates/{id}`. The collection endpoint does not return the certificate content.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The unique identifier returned by Microsoft Graph.",
			},
			"certificate_authority_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The unique identifier of the Microsoft-managed TLS certificate authority to retrieve.",
				Validators: []validator.String{
					stringvalidator.RegexMatches(regexp.MustCompile(constants.GuidRegex), "must be a valid UUID"),
				},
			},
			"name": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The internal certificate authority name.",
			},
			"common_name": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The common name of the Microsoft-managed root certificate authority.",
			},
			"organization_name": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The organization name of the Microsoft-managed root certificate authority.",
			},
			"validity_months": schema.Int32Attribute{
				Computed:            true,
				MarkdownDescription: "The configured root certificate validity period in months.",
			},
			"status": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The lifecycle status reported by Microsoft Graph, such as `unknownFutureValue` immediately after a disabled create, `disabled`, `enrolling`, or `active`.",
			},
			"created_date_time": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The date and time when the certificate authority was created.",
			},
			"validity_start_date_time": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The start of the root certificate validity period.",
			},
			"validity_end_date_time": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The end of the root certificate validity period.",
			},
			"certificate": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The PEM-encoded Microsoft-managed root CA certificate returned only by the individual certificate endpoint.",
			},
			"timeouts": commonschema.DatasourceTimeouts(ctx),
		},
	}
}
