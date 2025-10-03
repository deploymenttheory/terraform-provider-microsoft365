package auditEvents

import (
	"context"
	"fmt"
	"strings"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	abstractions "github.com/microsoft/kiota-abstractions-go"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/tenantrelationships"
)

// Open is called when the ephemeral resource is created
func (r *AuditEventsEphemeralResource) Open(ctx context.Context, req ephemeral.OpenRequest, resp *ephemeral.OpenResponse) {
	tflog.Debug(ctx, "Starting Open method for audit events ephemeral resource")

	var data AuditEventsEphemeralModel

	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	filterType := data.FilterType.ValueString()
	tflog.Debug(ctx, fmt.Sprintf("Starting ephemeral audit events with filter_type: %s", filterType))

	var filteredItems []AuditEventModel
	filterValue := data.FilterValue.ValueString()

	switch filterType {
	case "id":
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

	case "odata":
		// Add "ConsistencyLevel: eventual" header for advanced OData queries
		headers := abstractions.NewRequestHeaders()
		headers.Add("ConsistencyLevel", "eventual")

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

	default:
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

	var data AuditEventsEphemeralModel

	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

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
