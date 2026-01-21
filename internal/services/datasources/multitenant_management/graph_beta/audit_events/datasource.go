package graphBetaAuditEvents

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
	DataSourceName = "microsoft365_graph_beta_multitenant_management_audit_events"
	ReadTimeout    = 180
)

var (
	// Basic datasource interface (CRUD operations)
	_ datasource.DataSource = &AuditEventsDataSource{}

	// Allows the datasource to be configured with the provider client
	_ datasource.DataSourceWithConfigure = &AuditEventsDataSource{}
)

func NewAuditEventsDataSource() datasource.DataSource {
	return &AuditEventsDataSource{
		ReadPermissions: []string{
			"ManagedTenant.Read.All",
		},
	}
}

type AuditEventsDataSource struct {
	client          *msgraphbetasdk.GraphServiceClient
	ReadPermissions []string
}

// Metadata returns the datasource type name.
func (r *AuditEventsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = DataSourceName
}

// Configure sets the client for the data source
func (d *AuditEventsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.client = client.SetGraphBetaClientForDataSource(ctx, req, resp, DataSourceName)
}

// Schema defines the schema for the data source
func (d *AuditEventsDataSource) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Retrieves audit events from Microsoft 365 managed tenants using the `/tenantRelationships/managedTenants/auditEvents` endpoint. This data source is used to query administrative activities and changes across managed tenants for security monitoring and compliance reporting.",
		Attributes: map[string]schema.Attribute{
			"filter_type": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Type of filter to apply. Valid values are: `all`, `id`, `display_name`, `odata`.",
				Validators: []validator.String{
					stringvalidator.OneOf("all", "id", "display_name", "odata"),
				},
			},
			"filter_value": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Value to filter by. Not required when filter_type is 'all'.",
			},
			"odata_filter": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "OData $filter parameter for filtering results. Only used when filter_type is 'odata'.",
			},
			"odata_top": schema.Int32Attribute{
				Optional:            true,
				MarkdownDescription: "OData $top parameter to limit the number of results. Only used when filter_type is 'odata'.",
			},
			"odata_skip": schema.Int32Attribute{
				Optional:            true,
				MarkdownDescription: "OData $skip parameter for pagination. Only used when filter_type is 'odata'.",
			},
			"odata_select": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "OData $select parameter to specify which fields to include. Only used when filter_type is 'odata'.",
			},
			"odata_orderby": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "OData $orderby parameter to sort results. Only used when filter_type is 'odata'.",
			},
			"items": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "The list of audit events that match the filter criteria.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "Unique identifier for the audit event.",
						},
						"activity": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "A string that describes the activity that was performed.",
						},
						"activity_date_time": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The date and time when the activity was performed.",
						},
						"activity_id": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The identifier of the activity request that made the audit event.",
						},
						"category": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "A category that represents a logical grouping of activities.",
						},
						"http_verb": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The HTTP verb that was used when making the API request.",
						},
						"initiated_by_app_id": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The identifier for the app that was used to make the request.",
						},
						"initiated_by_upn": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The UPN of the user who initiated the activity.",
						},
						"initiated_by_user_id": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The identifier for the user who initiated the activity.",
						},
						"ip_address": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The IP address of where the activity was initiated. This may be an IPv4 or IPv6 address.",
						},
						"request_id": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The identifier of the request that made the audit event.",
						},
						"request_url": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The URL of the request that made the audit event.",
						},
						"tenant_ids": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The Azure Active Directory tenant identifier for the managed tenant that was affected by this event.",
						},
						"tenant_names": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The tenant name that was affected by this event.",
						},
					},
				},
			},
			"timeouts": commonschema.ResourceTimeouts(ctx),
		},
	}
}
