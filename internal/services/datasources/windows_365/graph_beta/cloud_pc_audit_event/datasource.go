package graphBetaCloudPcAuditEvent

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
	DataSourceName = "microsoft365_graph_beta_windows_365_cloud_pc_audit_event"
	ReadTimeout    = 180
)

var (
	_ datasource.DataSource              = &CloudPcAuditEventDataSource{}
	_ datasource.DataSourceWithConfigure = &CloudPcAuditEventDataSource{}
)

func NewCloudPcAuditEventDataSource() datasource.DataSource {
	return &CloudPcAuditEventDataSource{
		ReadPermissions: []string{
			"CloudPC.Read.All",
		},
	}
}

type CloudPcAuditEventDataSource struct {
	client          *msgraphbetasdk.GraphServiceClient
	ReadPermissions []string
}

func (d *CloudPcAuditEventDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = DataSourceName
}

func (d *CloudPcAuditEventDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.client = client.SetGraphBetaClientForDataSource(ctx, req, resp, DataSourceName)
}

func (d *CloudPcAuditEventDataSource) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Retrieves Cloud PC audit events from Microsoft Intune. Using the endpoint '/deviceManagement/virtualEndpoint/auditEvents'. Supports filtering by all, id, or display_name.",
		Attributes: map[string]schema.Attribute{
			"filter_type": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Type of filter to apply. Valid values are: `all`, `id`, `display_name`. Use 'all' to retrieve all audit events, 'id' to retrieve a specific event by its unique identifier, or 'display_name' to filter by the event's display name.",
				Validators: []validator.String{
					stringvalidator.OneOf("all", "id", "display_name"),
				},
			},
			"filter_value": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Value to filter by. Not required when filter_type is 'all'. For 'id', provide the audit event ID. For 'display_name', provide a substring to match against event display names.",
			},
			"items": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "The list of Cloud PC audit events that match the filter criteria.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The unique identifier for the audit event.",
						},
						"display_name": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The display name of the audit event, describing the action or object affected.",
						},
						"component_name": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The name of the component or controller that generated the event.",
						},
						"activity": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "A friendly name for the activity performed. May be null if not applicable.",
						},
						"activity_date_time": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The date and time (in UTC) when the activity was performed, in ISO 8601 format.",
						},
						"activity_type": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The type of activity that was performed, such as 'Delete CloudPcOnPremisesConnection'.",
						},
						"activity_operation_type": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The HTTP operation type of the activity (e.g., 'Create', 'Delete', 'Patch').",
						},
						"activity_result": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The result of the activity (e.g., 'Success', 'Failure', 'ClientError', 'Timeout').",
						},
						"correlation_id": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The client request ID used to correlate activity within the system.",
						},
						"category": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The audit category, such as 'Cloud PC'.",
						},
						"actor": schema.SingleNestedAttribute{
							Computed:            true,
							MarkdownDescription: "Details about the Microsoft Entra user or application associated with the audit event.",
							Attributes: map[string]schema.Attribute{
								"type": schema.StringAttribute{
									Computed: true, MarkdownDescription: "The type of actor (e.g., 'application', 'itPro', 'partner')."},
								"user_permissions":         schema.ListAttribute{ElementType: types.StringType, Computed: true, MarkdownDescription: "List of user and application permissions at the time of the event."},
								"application_id":           schema.StringAttribute{Computed: true, MarkdownDescription: "The Microsoft Entra application ID of the actor."},
								"application_display_name": schema.StringAttribute{Computed: true, MarkdownDescription: "The display name of the application."},
								"user_principal_name":      schema.StringAttribute{Computed: true, MarkdownDescription: "The user principal name (UPN) of the actor, if applicable."},
								"service_principal_name":   schema.StringAttribute{Computed: true, MarkdownDescription: "The service principal name (SPN) of the actor, if applicable."},
								"ip_address":               schema.StringAttribute{Computed: true, MarkdownDescription: "The IP address from which the activity was performed, if available."},
								"user_id":                  schema.StringAttribute{Computed: true, MarkdownDescription: "The Microsoft Entra user ID of the actor."},
								"user_role_scope_tags": schema.ListNestedAttribute{
									Computed:            true,
									MarkdownDescription: "List of role scope tags associated with the user at the time of the event.",
									NestedObject: schema.NestedAttributeObject{
										Attributes: map[string]schema.Attribute{
											"display_name":      schema.StringAttribute{Computed: true, MarkdownDescription: "The display name of the role scope tag."},
											"role_scope_tag_id": schema.StringAttribute{Computed: true, MarkdownDescription: "The unique identifier of the role scope tag."},
										},
									},
								},
								"remote_tenant_id": schema.StringAttribute{Computed: true, MarkdownDescription: "The delegated partner tenant ID, if applicable."},
								"remote_user_id":   schema.StringAttribute{Computed: true, MarkdownDescription: "The delegated partner user ID, if applicable."},
							},
						},
						"resources": schema.ListNestedAttribute{
							Computed:            true,
							MarkdownDescription: "List of resources affected by the audit event.",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"display_name": schema.StringAttribute{Computed: true, MarkdownDescription: "The display name of the resource entity."},
									"modified_properties": schema.ListNestedAttribute{
										Computed:            true,
										MarkdownDescription: "A list of properties that were modified as part of the event.",
										NestedObject: schema.NestedAttributeObject{
											Attributes: map[string]schema.Attribute{
												"display_name": schema.StringAttribute{Computed: true, MarkdownDescription: "The display name of the modified property."},
												"new_value":    schema.StringAttribute{Computed: true, MarkdownDescription: "The new value of the property after the change."},
												"old_value":    schema.StringAttribute{Computed: true, MarkdownDescription: "The old value of the property before the change."},
											},
										},
									},
									"resource_type": schema.StringAttribute{Computed: true, MarkdownDescription: "The type of the resource affected (e.g., 'CloudPcOnPremisesConnection')."},
									"resource_id":   schema.StringAttribute{Computed: true, MarkdownDescription: "The unique identifier of the resource affected."},
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
