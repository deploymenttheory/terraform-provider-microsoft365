package graphBetaNetworkForwardingProfile

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type NetworkForwardingProfileDataSourceModel struct {
	ID                    types.String             `tfsdk:"id"`
	ForwardingProfileID   types.String             `tfsdk:"forwarding_profile_id"`
	Name                  types.String             `tfsdk:"name"`
	TrafficForwardingType types.String             `tfsdk:"traffic_forwarding_type"`
	ListAll               types.Bool               `tfsdk:"list_all"`
	Items                 []ForwardingProfileModel `tfsdk:"items"`
	Timeouts              timeouts.Value           `tfsdk:"timeouts"`
}

type ForwardingProfileModel struct {
	ID                    types.String                       `tfsdk:"id"`
	Name                  types.String                       `tfsdk:"name"`
	Description           types.String                       `tfsdk:"description"`
	State                 types.String                       `tfsdk:"state"`
	Version               types.String                       `tfsdk:"version"`
	LastModifiedDateTime  types.String                       `tfsdk:"last_modified_date_time"`
	TrafficForwardingType types.String                       `tfsdk:"traffic_forwarding_type"`
	Priority              types.Int32                        `tfsdk:"priority"`
	IsCustomProfile       types.Bool                         `tfsdk:"is_custom_profile"`
	ClientFallbackAction  types.String                       `tfsdk:"client_fallback_action"`
	ServicePrincipalAppID types.String                       `tfsdk:"service_principal_app_id"`
	ServicePrincipalID    types.String                       `tfsdk:"service_principal_id"`
	Policies              []ForwardingProfilePolicyLinkModel `tfsdk:"policies"`
}

type ForwardingProfilePolicyLinkModel struct {
	PolicyLinkID          types.String `tfsdk:"policy_link_id"`
	Priority              types.Int64  `tfsdk:"priority"`
	State                 types.String `tfsdk:"state"`
	Version               types.String `tfsdk:"version"`
	PolicyID              types.String `tfsdk:"policy_id"`
	PolicyName            types.String `tfsdk:"policy_name"`
	PolicyDescription     types.String `tfsdk:"policy_description"`
	PolicyVersion         types.String `tfsdk:"policy_version"`
	TrafficForwardingType types.String `tfsdk:"traffic_forwarding_type"`
	PrivateAccessAppID    types.String `tfsdk:"private_access_app_id"`
}
