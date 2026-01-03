// REF: https://learn.microsoft.com/en-us/graph/api/resources/managedtenants-aggregatedpolicycompliance?view=graph-rest-beta
// REF: https://learn.microsoft.com/en-us/graph/api/managedtenants-managedtenant-list-aggregatedpolicycompliances?view=graph-rest-beta&tabs=http

package graphBetaAggregatedPolicyCompliances

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// AggregatedPolicyCompliancesDataSourceModel defines the data source model
type AggregatedPolicyCompliancesDataSourceModel struct {
	FilterType   types.String                      `tfsdk:"filter_type"`   // Required field to specify how to filter
	FilterValue  types.String                      `tfsdk:"filter_value"`  // Value to filter by (not used for "all" or "odata")
	ODataFilter  types.String                      `tfsdk:"odata_filter"`  // OData filter parameter
	ODataTop     types.Int32                       `tfsdk:"odata_top"`     // OData top parameter for limiting results
	ODataSkip    types.Int32                       `tfsdk:"odata_skip"`    // OData skip parameter for pagination
	ODataSelect  types.String                      `tfsdk:"odata_select"`  // OData select parameter for field selection
	ODataOrderBy types.String                      `tfsdk:"odata_orderby"` // OData orderby parameter for sorting
	Items        []AggregatedPolicyComplianceModel `tfsdk:"items"`         // List of aggregated policy compliances that match the filters
	Timeouts     timeouts.Value                    `tfsdk:"timeouts"`
}

// AggregatedPolicyComplianceModel represents a single aggregated policy compliance
type AggregatedPolicyComplianceModel struct {
	ID                          types.String `tfsdk:"id"`
	CompliancePolicyId          types.String `tfsdk:"compliance_policy_id"`
	CompliancePolicyName        types.String `tfsdk:"compliance_policy_name"`
	CompliancePolicyPlatform    types.String `tfsdk:"compliance_policy_platform"`
	CompliancePolicyType        types.String `tfsdk:"compliance_policy_type"`
	LastRefreshedDateTime       types.String `tfsdk:"last_refreshed_date_time"`
	NumberOfCompliantDevices    types.Int64  `tfsdk:"number_of_compliant_devices"`
	NumberOfErrorDevices        types.Int64  `tfsdk:"number_of_error_devices"`
	NumberOfNonCompliantDevices types.Int64  `tfsdk:"number_of_non_compliant_devices"`
	PolicyModifiedDateTime      types.String `tfsdk:"policy_modified_date_time"`
	TenantDisplayName           types.String `tfsdk:"tenant_display_name"`
	TenantId                    types.String `tfsdk:"tenant_id"`
}
