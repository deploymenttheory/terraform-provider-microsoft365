// REF: https://learn.microsoft.com/en-us/graph/api/resources/managedtenants-auditevent?view=graph-rest-beta
// REF: https://learn.microsoft.com/en-us/graph/api/managedtenants-managedtenant-list-auditevents?view=graph-rest-beta&tabs=http

package auditEvents

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// AuditEventsEphemeralModel defines the ephemeral resource model for audit events
type AuditEventsEphemeralModel struct {
	FilterType   types.String      `tfsdk:"filter_type"`   // Required field to specify how to filter
	FilterValue  types.String      `tfsdk:"filter_value"`  // Value to filter by (not used for "all" or "odata")
	ODataFilter  types.String      `tfsdk:"odata_filter"`  // OData filter parameter
	ODataTop     types.Int32       `tfsdk:"odata_top"`     // OData top parameter for limiting results
	ODataSkip    types.Int32       `tfsdk:"odata_skip"`    // OData skip parameter for pagination
	ODataSelect  types.String      `tfsdk:"odata_select"`  // OData select parameter for field selection
	ODataOrderBy types.String      `tfsdk:"odata_orderby"` // OData orderby parameter for sorting
	Items        []AuditEventModel `tfsdk:"items"`         // List of audit events that match the filters
}

// AuditEventModel represents a single audit event
type AuditEventModel struct {
	ID                types.String `tfsdk:"id"`
	Activity          types.String `tfsdk:"activity"`
	ActivityDateTime  types.String `tfsdk:"activity_date_time"`
	ActivityId        types.String `tfsdk:"activity_id"`
	Category          types.String `tfsdk:"category"`
	HttpVerb          types.String `tfsdk:"http_verb"`
	InitiatedByAppId  types.String `tfsdk:"initiated_by_app_id"`
	InitiatedByUpn    types.String `tfsdk:"initiated_by_upn"`
	InitiatedByUserId types.String `tfsdk:"initiated_by_user_id"`
	IpAddress         types.String `tfsdk:"ip_address"`
	RequestId         types.String `tfsdk:"request_id"`
	RequestUrl        types.String `tfsdk:"request_url"`
	TenantIds         types.String `tfsdk:"tenant_ids"`
	TenantNames       types.String `tfsdk:"tenant_names"`
}
