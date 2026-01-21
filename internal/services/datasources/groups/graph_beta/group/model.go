// REF: https://learn.microsoft.com/en-us/graph/api/group-get?view=graph-rest-beta
package graphBetaGroup

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/datasource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type GroupDataSourceModel struct {
	ID                            types.String   `tfsdk:"id"`
	DisplayName                   types.String   `tfsdk:"display_name"`
	ObjectId                      types.String   `tfsdk:"object_id"`
	MailNickname                  types.String   `tfsdk:"mail_nickname"`
	ODataQuery                    types.String   `tfsdk:"odata_query"`
	MailEnabled                   types.Bool     `tfsdk:"mail_enabled"`
	SecurityEnabled               types.Bool     `tfsdk:"security_enabled"`
	Description                   types.String   `tfsdk:"description"`
	Classification                types.String   `tfsdk:"classification"`
	GroupTypes                    types.Set      `tfsdk:"group_types"`
	Visibility                    types.String   `tfsdk:"visibility"`
	AssignableToRole              types.Bool     `tfsdk:"assignable_to_role"`
	DynamicMembershipEnabled      types.Bool     `tfsdk:"dynamic_membership_enabled"`
	MembershipRule                types.String   `tfsdk:"membership_rule"`
	MembershipRuleProcessingState types.String   `tfsdk:"membership_rule_processing_state"`
	CreatedDateTime               types.String   `tfsdk:"created_date_time"`
	Mail                          types.String   `tfsdk:"mail"`
	ProxyAddresses                types.Set      `tfsdk:"proxy_addresses"`
	AssignedLicenses              types.List     `tfsdk:"assigned_licenses"`
	HasMembersWithLicenseErrors   types.Bool     `tfsdk:"has_members_with_license_errors"`
	HideFromAddressLists          types.Bool     `tfsdk:"hide_from_address_lists"`
	HideFromOutlookClients        types.Bool     `tfsdk:"hide_from_outlook_clients"`
	OnPremisesSyncEnabled         types.Bool     `tfsdk:"onpremises_sync_enabled"`
	OnPremisesLastSyncDateTime    types.String   `tfsdk:"onpremises_last_sync_date_time"`
	OnPremisesSamAccountName      types.String   `tfsdk:"onpremises_sam_account_name"`
	OnPremisesDomainName          types.String   `tfsdk:"onpremises_domain_name"`
	OnPremisesNetBiosName         types.String   `tfsdk:"onpremises_netbios_name"`
	OnPremisesSecurityIdentifier  types.String   `tfsdk:"onpremises_security_identifier"`
	Members                       types.Set      `tfsdk:"members"`
	Owners                        types.Set      `tfsdk:"owners"`
	Timeouts                      timeouts.Value `tfsdk:"timeouts"`
}

type AssignedLicenseModel struct {
	SkuId         types.String `tfsdk:"sku_id"`
	DisabledPlans types.Set    `tfsdk:"disabled_plans"`
}
