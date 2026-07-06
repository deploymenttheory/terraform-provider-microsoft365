package graphBetaNetworkFilteringProfilePolicyLink

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	models "github.com/microsoftgraph/msgraph-beta-sdk-go/models/networkaccess"
)

func MapRemoteStateToTerraform(ctx context.Context, data *NetworkFilteringProfilePolicyLinkResourceModel, filteringProfileID string, link models.PolicyLinkable) {
	if link == nil {
		tflog.Debug(ctx, "Remote policy link is nil")
		return
	}

	policyLinkID := convert.GraphToFrameworkString(link.GetId())
	data.FilteringProfileID = types.StringValue(filteringProfileID)
	data.PolicyLinkID = policyLinkID
	data.ID = compositeID(filteringProfileID, policyLinkID.ValueString())
	data.State = convert.GraphToFrameworkEnum(link.GetState())
	data.Version = convert.GraphToFrameworkString(link.GetVersion())

	if policy := link.GetPolicy(); policy != nil {
		data.PolicyID = convert.GraphToFrameworkString(policy.GetId())
		if policyType := policy.GetOdataType(); policyType != nil {
			data.PolicyODataType = types.StringValue(*policyType)
			data.PolicyType = policyTypeFromPolicyODataType(*policyType)
		}
	}

	if odataType := link.GetOdataType(); odataType != nil {
		data.PolicyLinkODataType = types.StringValue(*odataType)
		if data.PolicyType.IsNull() || data.PolicyType.IsUnknown() {
			data.PolicyType = policyTypeFromLinkODataType(*odataType)
		}
	}

	if filteringLink, ok := link.(models.FilteringPolicyLinkable); ok {
		data.Priority = convert.GraphToFrameworkInt64(filteringLink.GetPriority())
		data.LoggingState = convert.GraphToFrameworkEnum(filteringLink.GetLoggingState())
		data.CreatedDateTime = convert.GraphToFrameworkTime(filteringLink.GetCreatedDateTime())
		data.LastModifiedDateTime = convert.GraphToFrameworkTime(filteringLink.GetLastModifiedDateTime())
	} else {
		data.Priority = types.Int64Null()
		data.LoggingState = types.StringNull()
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished stating %s with id %s", ResourceName, data.ID.ValueString()))
}

func policyTypeFromPolicyODataType(odataType string) types.String {
	switch odataType {
	case filteringPolicyODataType:
		return types.StringValue(policyTypeFiltering)
	case webFilteringPolicyODataType:
		return types.StringValue(policyTypeWebFiltering)
	case cloudFirewallPolicyODataType:
		return types.StringValue(policyTypeCloudFirewall)
	case threatIntelligencePolicyODataType:
		return types.StringValue(policyTypeThreatIntelligence)
	case tlsInspectionPolicyODataType:
		return types.StringValue(policyTypeTlsInspection)
	default:
		return types.StringValue(policyTypeCustom)
	}
}

func policyTypeFromLinkODataType(odataType string) types.String {
	switch odataType {
	case filteringPolicyLinkODataType:
		return types.StringValue(policyTypeFiltering)
	case webFilteringPolicyLinkODataType:
		return types.StringValue(policyTypeWebFiltering)
	case cloudFirewallPolicyLinkODataType:
		return types.StringValue(policyTypeCloudFirewall)
	case threatIntelligencePolicyLinkODataType:
		return types.StringValue(policyTypeThreatIntelligence)
	case tlsInspectionPolicyLinkODataType:
		return types.StringValue(policyTypeTlsInspection)
	default:
		return types.StringValue(policyTypeCustom)
	}
}

func findPolicyLink(profile models.FilteringProfileable, policyLinkID, policyID string) models.PolicyLinkable {
	if profile == nil {
		return nil
	}

	for _, link := range profile.GetPolicies() {
		if link == nil {
			continue
		}
		if policyLinkID != "" {
			if id := link.GetId(); id != nil && *id == policyLinkID {
				return link
			}
			continue
		}
		if policyID != "" {
			if policy := link.GetPolicy(); policy != nil {
				if id := policy.GetId(); id != nil && *id == policyID {
					return link
				}
			}
		}
	}

	return nil
}
