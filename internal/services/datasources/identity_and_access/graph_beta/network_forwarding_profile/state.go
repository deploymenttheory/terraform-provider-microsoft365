package graphBetaNetworkForwardingProfile

import (
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-framework/types"
	models "github.com/microsoftgraph/msgraph-beta-sdk-go/models/networkaccess"
)

func MapRemoteStateToDataSource(profile models.ForwardingProfileable) ForwardingProfileModel {
	if profile == nil {
		return ForwardingProfileModel{}
	}

	result := ForwardingProfileModel{
		ID:                    convert.GraphToFrameworkString(profile.GetId()),
		Name:                  convert.GraphToFrameworkString(profile.GetName()),
		Description:           convert.GraphToFrameworkString(profile.GetDescription()),
		State:                 convert.GraphToFrameworkEnum(profile.GetState()),
		Version:               convert.GraphToFrameworkString(profile.GetVersion()),
		LastModifiedDateTime:  convert.GraphToFrameworkTime(profile.GetLastModifiedDateTime()),
		TrafficForwardingType: convert.GraphToFrameworkEnum(profile.GetTrafficForwardingType()),
		Priority:              convert.GraphToFrameworkInt32(profile.GetPriority()),
		IsCustomProfile:       convert.GraphToFrameworkBool(profile.GetIsCustomProfile()),
		ClientFallbackAction:  additionalString(profile.GetAdditionalData(), "clientFallbackAction"),
		ServicePrincipalAppID: types.StringNull(),
		ServicePrincipalID:    types.StringNull(),
		Policies:              make([]ForwardingProfilePolicyLinkModel, 0),
	}

	if servicePrincipal := profile.GetServicePrincipal(); servicePrincipal != nil {
		result.ServicePrincipalAppID = convert.GraphToFrameworkString(servicePrincipal.GetAppId())
		result.ServicePrincipalID = convert.GraphToFrameworkString(servicePrincipal.GetId())
	}

	for _, link := range profile.GetPolicies() {
		if link == nil {
			continue
		}
		item := ForwardingProfilePolicyLinkModel{
			PolicyLinkID:       convert.GraphToFrameworkString(link.GetId()),
			State:              convert.GraphToFrameworkEnum(link.GetState()),
			Version:            convert.GraphToFrameworkString(link.GetVersion()),
			PrivateAccessAppID: types.StringNull(),
		}
		if forwardingLink, ok := link.(models.ForwardingPolicyLinkable); ok {
			item.Priority = convert.GraphToFrameworkInt64(forwardingLink.GetPriority())
		} else {
			item.Priority = types.Int64Null()
		}
		if policy := link.GetPolicy(); policy != nil {
			item.PolicyID = convert.GraphToFrameworkString(policy.GetId())
			item.PolicyName = convert.GraphToFrameworkString(policy.GetName())
			item.PolicyDescription = convert.GraphToFrameworkString(policy.GetDescription())
			item.PolicyVersion = convert.GraphToFrameworkString(policy.GetVersion())
			if forwardingPolicy, ok := policy.(models.ForwardingPolicyable); ok {
				item.TrafficForwardingType = convert.GraphToFrameworkEnum(forwardingPolicy.GetTrafficForwardingType())
				item.PrivateAccessAppID = convert.GraphToFrameworkString(forwardingPolicy.GetPrivateAccessAppId())
			}
		}
		result.Policies = append(result.Policies, item)
	}

	return result
}

func additionalString(values map[string]any, key string) types.String {
	if values == nil {
		return types.StringNull()
	}
	switch value := values[key].(type) {
	case string:
		return types.StringValue(value)
	case *string:
		return convert.GraphToFrameworkString(value)
	default:
		return types.StringNull()
	}
}
