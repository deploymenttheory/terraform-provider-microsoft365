package graphBetaNetworkForwardingProfilePolicyLink

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (m policyLinkMatch) toModel() NetworkForwardingProfilePolicyLinkDataSourceModel {
	return NetworkForwardingProfilePolicyLinkDataSourceModel{
		ID:                    types.StringValue(m.profile.ID.ValueString() + "/" + m.link.PolicyLinkID.ValueString()),
		ForwardingProfileID:   m.profile.ID,
		ForwardingProfileName: m.profile.Name,
		TrafficForwardingType: m.profile.TrafficForwardingType,
		PolicyName:            m.link.PolicyName,
		PolicyLinkID:          m.link.PolicyLinkID,
		Priority:              m.link.Priority,
		State:                 m.link.State,
		Version:               m.link.Version,
		PolicyID:              m.link.PolicyID,
		PolicyDescription:     m.link.PolicyDescription,
		PolicyVersion:         m.link.PolicyVersion,
		PrivateAccessAppID:    m.link.PrivateAccessAppID,
	}
}
