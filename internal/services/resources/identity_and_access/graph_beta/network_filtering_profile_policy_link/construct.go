package graphBetaNetworkFilteringProfilePolicyLink

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const (
	// Microsoft Learn documents policyLink as an abstract type with these Graph beta
	// derived types:
	// https://learn.microsoft.com/en-us/graph/api/resources/networkaccess-policylink?view=graph-rest-beta
	//
	// webFilteringPolicyLink is intentionally included even though it is not in the
	// current Microsoft Graph beta Go SDK policyLink discriminator. Entra admin center
	// XHR was observed linking V2 web filtering policies by POSTing this explicit
	// @odata.type to /networkAccess/filteringProfiles/{filteringProfileId}/policies.
	//
	// prompt, content, and netskope_dlp are portal-first cases:
	// the Entra admin center JavaScript bundle defines promptPolicyLink,
	// filePolicyLink, contentPolicyLink, and securityProviderPolicyLink, but these
	// types were not present in Microsoft Graph beta $metadata or the generated Go
	// SDK when this resource was added. Observed Entra admin center XHR for the UI's
	// content policy option sends filePolicyLink/filePolicy.
	//
	// Netskope DLP policy selection was observed querying
	// /networkAccess/securityProviderPolicies with a Netskope schema filter, so
	// Terraform exposes that UI concept as netskope_dlp while keeping the
	// generic securityProviderPolicy Graph type internal.
	//
	// Keep these values as strings instead of SDK constructors so portal-first policy
	// link types can be sent before the generated SDK catches up.
	filteringPolicyLinkODataType          = "#microsoft.graph.networkaccess.filteringPolicyLink"
	filteringPolicyODataType              = "#microsoft.graph.networkaccess.filteringPolicy"
	webFilteringPolicyLinkODataType       = "#microsoft.graph.networkaccess.webFilteringPolicyLink"
	webFilteringPolicyODataType           = "#microsoft.graph.networkaccess.webFilteringPolicy"
	promptPolicyLinkODataType             = "#microsoft.graph.networkaccess.promptPolicyLink"
	promptPolicyODataType                 = "#microsoft.graph.networkaccess.promptPolicy"
	filePolicyLinkODataType               = "#microsoft.graph.networkaccess.filePolicyLink"
	filePolicyODataType                   = "#microsoft.graph.networkaccess.filePolicy"
	securityProviderPolicyLinkODataType   = "#microsoft.graph.networkaccess.securityProviderPolicyLink"
	securityProviderPolicyODataType       = "#microsoft.graph.networkaccess.securityProviderPolicy"
	cloudFirewallPolicyLinkODataType      = "#microsoft.graph.networkaccess.cloudFirewallPolicyLink"
	cloudFirewallPolicyODataType          = "#microsoft.graph.networkaccess.cloudFirewallPolicy"
	threatIntelligencePolicyLinkODataType = "#microsoft.graph.networkaccess.threatIntelligencePolicyLink"
	threatIntelligencePolicyODataType     = "#microsoft.graph.networkaccess.threatIntelligencePolicy"
	tlsInspectionPolicyLinkODataType      = "#microsoft.graph.networkaccess.tlsInspectionPolicyLink"
	tlsInspectionPolicyODataType          = "#microsoft.graph.networkaccess.tlsInspectionPolicy"
)

var errUnsupportedPolicyType = errors.New("unsupported policy_type")

// policyODataTypes holds the OData types of a policy link and of the policy it links.
type policyODataTypes struct {
	link   string
	policy string
}

var policyTypeToODataTypes = map[string]policyODataTypes{
	policyTypeFiltering: {
		link:   filteringPolicyLinkODataType,
		policy: filteringPolicyODataType,
	},
	policyTypeWebFiltering: {
		link:   webFilteringPolicyLinkODataType,
		policy: webFilteringPolicyODataType,
	},
	policyTypePrompt: {
		link:   promptPolicyLinkODataType,
		policy: promptPolicyODataType,
	},
	policyTypeContent: {
		link:   filePolicyLinkODataType,
		policy: filePolicyODataType,
	},
	policyTypeNetskopeDlp: {
		link:   securityProviderPolicyLinkODataType,
		policy: securityProviderPolicyODataType,
	},
	policyTypeCloudFirewall: {
		link:   cloudFirewallPolicyLinkODataType,
		policy: cloudFirewallPolicyODataType,
	},
	policyTypeThreatIntelligence: {
		link:   threatIntelligencePolicyLinkODataType,
		policy: threatIntelligencePolicyODataType,
	},
	policyTypeTlsInspection: {
		link:   tlsInspectionPolicyLinkODataType,
		policy: tlsInspectionPolicyODataType,
	},
}

// constructResource builds the policy link body. The create body references the linked policy and, for filtering
// links, carries priority and loggingState; the update body carries only the link type and state.
//
// It also fills the computed request fields on data so state matches what was sent: filtering links default
// priority to 100 and logging_state to enabled, and every other link type has neither.
func constructResource(ctx context.Context, data *NetworkFilteringProfilePolicyLinkResourceModel, includePolicy bool) (policyLinkRequest, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s resource from model", ResourceName))

	odataTypes, ok := policyTypeToODataTypes[data.PolicyType.ValueString()]
	if !ok {
		return policyLinkRequest{}, fmt.Errorf("%w %q", errUnsupportedPolicyType, data.PolicyType.ValueString())
	}

	if data.PolicyType.ValueString() == policyTypeFiltering {
		if data.Priority.IsNull() || data.Priority.IsUnknown() {
			data.Priority = types.Int64Value(100)
		}
		if data.LoggingState.IsNull() || data.LoggingState.IsUnknown() {
			data.LoggingState = types.StringValue("enabled")
		}
	} else {
		data.Priority = types.Int64Null()
		data.LoggingState = types.StringNull()
	}

	requestBody := policyLinkRequest{ODataType: odataTypes.link, State: data.State.ValueString()}
	if !includePolicy {
		return requestBody, nil
	}

	if data.PolicyID.ValueString() != "" {
		requestBody.Policy = &policyReference{ID: data.PolicyID.ValueString(), ODataType: odataTypes.policy}
	}
	if data.PolicyType.ValueString() == policyTypeFiltering {
		requestBody.Priority = data.Priority.ValueInt64Pointer()
		requestBody.LoggingState = data.LoggingState.ValueString()
	}

	return requestBody, nil
}
