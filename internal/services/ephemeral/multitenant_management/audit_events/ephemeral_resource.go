package auditEvents

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

const (
	EphemeralResourceName = "graph_beta_multitenant_management_audit_event"
	CreateTimeout         = 180
	UpdateTimeout         = 180
	ReadTimeout           = 180
	DeleteTimeout         = 180
)

// Ensure the implementation satisfies the expected interfaces
var (
	_ ephemeral.EphemeralResource              = &AuditEventsEphemeralResource{}
	_ ephemeral.EphemeralResourceWithConfigure = &AuditEventsEphemeralResource{}
)

// NewAuditEventsEphemeralResource is a helper function to simplify provider implementation
func NewAuditEventsEphemeralResource() ephemeral.EphemeralResource {
	return &AuditEventsEphemeralResource{
		ReadPermissions: []string{
			"ManagedTenant.Read.All",
		},
	}
}

// AuditEventsEphemeralResource is the ephemeral resource implementation
type AuditEventsEphemeralResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ProviderTypeName string
	TypeName         string
	ReadPermissions  []string
}

// Metadata returns the resource type name
func (r *AuditEventsEphemeralResource) Metadata(_ context.Context, req ephemeral.MetadataRequest, resp *ephemeral.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + EphemeralResourceName
}

// Schema defines the schema for the ephemeral resource
func (r *AuditEventsEphemeralResource) Schema(_ context.Context, _ ephemeral.SchemaRequest, resp *ephemeral.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Retrieves audit events from Microsoft 365 managed tenants as an ephemeral resource. This does not persist in state and fetches fresh data on each execution.",
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
		},
	}
}
