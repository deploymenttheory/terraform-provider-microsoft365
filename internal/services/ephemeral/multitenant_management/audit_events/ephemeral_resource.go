package auditEvents

import (
	"context"
	"fmt"
	"strings"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	abstractions "github.com/microsoft/kiota-abstractions-go"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/tenantrelationships"
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
	resp.TypeName = req.ProviderTypeName + "_multitenant_management_audit_events"
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

// Open is called when the ephemeral resource is created
func (r *AuditEventsEphemeralResource) Open(ctx context.Context, req ephemeral.OpenRequest, resp *ephemeral.OpenResponse) {
	tflog.Debug(ctx, "Starting Open method for audit events ephemeral resource")

	// Create a new model to hold the configuration
	var data AuditEventsEphemeralModel

	// Get the configuration from the request
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	filterType := data.FilterType.ValueString()
	tflog.Debug(ctx, fmt.Sprintf("Starting ephemeral audit events with filter_type: %s", filterType))

	var filteredItems []AuditEventModel
	filterValue := data.FilterValue.ValueString()

	if filterType == "id" {
		auditEvent, err := r.client.
			TenantRelationships().
			ManagedTenants().
			AuditEvents().
			ByAuditEventId(filterValue).
			Get(ctx, nil)

		if err != nil {
			errors.HandleKiotaGraphError(ctx, err, resp, "Open", r.ReadPermissions)
			return
		}

		eventItem := MapRemoteStateToEphemeral(ctx, auditEvent)
		filteredItems = append(filteredItems, eventItem)
	} else if filterType == "odata" {
		// Add "ConsistencyLevel: eventual" header for advanced OData queries
		headers := abstractions.NewRequestHeaders()
		headers.Add("ConsistencyLevel", "eventual")

		// Initialize request parameters for audit events
		requestParameters := &tenantrelationships.ManagedTenantsAuditEventsRequestBuilderGetRequestConfiguration{
			Headers:         headers,
			QueryParameters: &tenantrelationships.ManagedTenantsAuditEventsRequestBuilderGetQueryParameters{},
		}

		if !data.ODataFilter.IsNull() && data.ODataFilter.ValueString() != "" {
			filter := data.ODataFilter.ValueString()
			requestParameters.QueryParameters.Filter = &filter
			tflog.Debug(ctx, fmt.Sprintf("Setting OData filter: %s", filter))
		}

		if !data.ODataTop.IsNull() {
			topValue := data.ODataTop.ValueInt32()
			requestParameters.QueryParameters.Top = &topValue
			tflog.Debug(ctx, fmt.Sprintf("Setting OData top: %d", topValue))
		}

		if !data.ODataSkip.IsNull() {
			skipValue := data.ODataSkip.ValueInt32()
			requestParameters.QueryParameters.Skip = &skipValue
			tflog.Debug(ctx, fmt.Sprintf("Setting OData skip: %d", skipValue))
		}

		if !data.ODataSelect.IsNull() && data.ODataSelect.ValueString() != "" {
			selectFields := strings.Split(data.ODataSelect.ValueString(), ",")
			requestParameters.QueryParameters.Select = selectFields
			tflog.Debug(ctx, fmt.Sprintf("Setting OData select: %v", selectFields))
		}

		if !data.ODataOrderBy.IsNull() && data.ODataOrderBy.ValueString() != "" {
			orderbyFields := strings.Split(data.ODataOrderBy.ValueString(), ",")
			requestParameters.QueryParameters.Orderby = orderbyFields
			tflog.Debug(ctx, fmt.Sprintf("Setting OData orderby: %v", orderbyFields))
		}

		respList, err := r.client.
			TenantRelationships().
			ManagedTenants().
			AuditEvents().
			Get(ctx, requestParameters)

		if err != nil {
			tflog.Error(ctx, fmt.Sprintf("Error in OData query: %v", err))
			errors.HandleKiotaGraphError(ctx, err, resp, "Open", r.ReadPermissions)
			return
		}

		tflog.Debug(ctx, fmt.Sprintf("API returned %d results", len(respList.GetValue())))

		for _, auditEvent := range respList.GetValue() {
			eventItem := MapRemoteStateToEphemeral(ctx, auditEvent)
			filteredItems = append(filteredItems, eventItem)
		}
	} else {
		// For "all" and "display_name", get the full list and filter locally
		respList, err := r.client.
			TenantRelationships().
			ManagedTenants().
			AuditEvents().
			Get(ctx, nil)

		if err != nil {
			errors.HandleKiotaGraphError(ctx, err, resp, "Open", r.ReadPermissions)
			return
		}

		for _, auditEvent := range respList.GetValue() {
			eventItem := MapRemoteStateToEphemeral(ctx, auditEvent)

			switch filterType {
			case "all":
				filteredItems = append(filteredItems, eventItem)

			case "display_name":
				// For audit events, we'll search in activity field for display_name filter
				if auditEvent.GetActivity() != nil && strings.Contains(
					strings.ToLower(*auditEvent.GetActivity()),
					strings.ToLower(filterValue)) {
					filteredItems = append(filteredItems, eventItem)
				}
			}
		}
	}

	data.Items = filteredItems

	// Set the result
	diags = resp.Result.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)

	tflog.Debug(ctx, fmt.Sprintf("Completed ephemeral audit events Open method, found %d items", len(filteredItems)))
}

// Configure is called to pass the provider configured client to the resource
func (r *AuditEventsEphemeralResource) Configure(ctx context.Context, req ephemeral.ConfigureRequest, resp *ephemeral.ConfigureResponse) {
	r.client = client.SetGraphBetaClientForEphemeralResource(ctx, req, resp, r.TypeName)
}

// ValidateConfig validates the configuration
func (r *AuditEventsEphemeralResource) ValidateConfig(ctx context.Context, req ephemeral.ValidateConfigRequest, resp *ephemeral.ValidateConfigResponse) {
	tflog.Debug(ctx, "Validating audit events ephemeral resource configuration")

	// Create a new model to hold the configuration
	var data AuditEventsEphemeralModel

	// Get the configuration from the request
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Validate filter_type and filter_value combination
	filterType := data.FilterType.ValueString()
	filterValue := data.FilterValue.ValueString()

	if filterType == "id" && filterValue == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("filter_value"),
			"Missing Filter Value",
			"filter_value is required when filter_type is 'id'.",
		)
	}

	if filterType == "display_name" && filterValue == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("filter_value"),
			"Missing Filter Value",
			"filter_value is required when filter_type is 'display_name'.",
		)
	}

	if filterType == "odata" && data.ODataFilter.IsNull() && filterValue == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("odata_filter"),
			"Missing OData Filter",
			"odata_filter is required when filter_type is 'odata'.",
		)
	}
}

// Close is called when the ephemeral resource is no longer needed
func (r *AuditEventsEphemeralResource) Close(ctx context.Context, req ephemeral.CloseRequest, resp *ephemeral.CloseResponse) {
	tflog.Debug(ctx, "Closing audit events ephemeral resource")

	// For audit events, we don't maintain any persistent connections, tokens, or resources
	// that need explicit cleanup. The Graph client is managed by the provider.
	// This is a no-op cleanup since audit events are read-only and stateless.

	tflog.Info(ctx, "Audit events ephemeral resource session closed - no cleanup required")
}

// Renew is called when the ephemeral resource needs to be renewed
func (r *AuditEventsEphemeralResource) Renew(ctx context.Context, req ephemeral.RenewRequest, resp *ephemeral.RenewResponse) {
	tflog.Debug(ctx, "Renewing audit events ephemeral resource")

	// According to the Terraform ephemeral resource framework documentation:
	// "Renew cannot return new result data for the ephemeral resource instance, so this logic
	// is only appropriate for remote objects like HashiCorp Vault leases, which can be renewed
	// without changing their data."
	//
	// Since audit events are read-only data queries (not leases or tokens that need renewal),
	// this is a no-op. If fresh audit event data is needed, a new ephemeral resource instance
	// should be created instead.

	tflog.Info(ctx, "Audit events do not require renewal - data is read-only")
}
