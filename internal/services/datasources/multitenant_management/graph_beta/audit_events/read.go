package graphBetaAuditEvents

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud"
	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	abstractions "github.com/microsoft/kiota-abstractions-go"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/tenantrelationships"
)

func (d *AuditEventsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var object AuditEventsDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	filterType := object.FilterType.ValueString()
	tflog.Debug(ctx, fmt.Sprintf("Starting Read method for datasource: %s with filter_type: %s", datasourceName, filterType))

	ctx, cancel := crud.HandleTimeout(ctx, object.Timeouts.Read, ReadTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	var filteredItems []AuditEventModel
	filterValue := object.FilterValue.ValueString()

	if filterType == "id" {
		auditEvent, err := d.client.
			TenantRelationships().
			ManagedTenants().
			AuditEvents().
			ByAuditEventId(filterValue).
			Get(ctx, nil)

		if err != nil {
			errors.HandleKiotaGraphError(ctx, err, resp, "Read", d.ReadPermissions)
			return
		}

		eventItem := MapRemoteStateToDataSource(ctx, auditEvent)
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

		if !object.ODataFilter.IsNull() && object.ODataFilter.ValueString() != "" {
			filter := object.ODataFilter.ValueString()
			requestParameters.QueryParameters.Filter = &filter
			tflog.Debug(ctx, fmt.Sprintf("Setting OData filter: %s", filter))
		}

		if !object.ODataTop.IsNull() {
			topValue := object.ODataTop.ValueInt32()
			requestParameters.QueryParameters.Top = &topValue
			tflog.Debug(ctx, fmt.Sprintf("Setting OData top: %d", topValue))
		}

		if !object.ODataSkip.IsNull() {
			skipValue := object.ODataSkip.ValueInt32()
			requestParameters.QueryParameters.Skip = &skipValue
			tflog.Debug(ctx, fmt.Sprintf("Setting OData skip: %d", skipValue))
		}

		if !object.ODataSelect.IsNull() && object.ODataSelect.ValueString() != "" {
			selectFields := strings.Split(object.ODataSelect.ValueString(), ",")
			requestParameters.QueryParameters.Select = selectFields
			tflog.Debug(ctx, fmt.Sprintf("Setting OData select: %v", selectFields))
		}

		if !object.ODataOrderBy.IsNull() && object.ODataOrderBy.ValueString() != "" {
			orderbyFields := strings.Split(object.ODataOrderBy.ValueString(), ",")
			requestParameters.QueryParameters.Orderby = orderbyFields
			tflog.Debug(ctx, fmt.Sprintf("Setting OData orderby: %v", orderbyFields))
		}

		respList, err := d.client.
			TenantRelationships().
			ManagedTenants().
			AuditEvents().
			Get(ctx, requestParameters)

		if err != nil {
			tflog.Error(ctx, fmt.Sprintf("Error in OData query: %v", err))
			errors.HandleKiotaGraphError(ctx, err, resp, "Read", d.ReadPermissions)
			return
		}

		tflog.Debug(ctx, fmt.Sprintf("API returned %d results", len(respList.GetValue())))

		for _, auditEvent := range respList.GetValue() {
			eventItem := MapRemoteStateToDataSource(ctx, auditEvent)
			filteredItems = append(filteredItems, eventItem)
		}
	} else {
		// For "all" and "display_name", get the full list and filter locally
		respList, err := d.client.
			TenantRelationships().
			ManagedTenants().
			AuditEvents().
			Get(ctx, nil)

		if err != nil {
			errors.HandleKiotaGraphError(ctx, err, resp, "Read", d.ReadPermissions)
			return
		}

		for _, auditEvent := range respList.GetValue() {
			eventItem := MapRemoteStateToDataSource(ctx, auditEvent)

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

	object.Items = filteredItems

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Datasource Read Method: %s, found %d items", datasourceName, len(filteredItems)))
}
