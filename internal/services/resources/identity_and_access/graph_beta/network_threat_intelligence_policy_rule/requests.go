package graphBetaNetworkThreatIntelligencePolicyRule

import (
	"context"

	abstractions "github.com/microsoft/kiota-abstractions-go"
	kiotahttp "github.com/microsoft/kiota-http-go"
	models "github.com/microsoftgraph/msgraph-beta-sdk-go/models/networkaccess"
	builders "github.com/microsoftgraph/msgraph-beta-sdk-go/networkaccess"
)

//nolint:wrapcheck // Shared Graph diagnostics require the original SDK error type.
func (r *NetworkThreatIntelligencePolicyRuleResource) createThreatIntelligencePolicyRule(
	ctx context.Context,
	parent string,
	body models.PolicyRuleable,
) (models.PolicyRuleable, error) {
	return r.client.NetworkAccess().
		ThreatIntelligencePolicies().
		ByThreatIntelligencePolicyId(parent).
		PolicyRules().
		Post(ctx, body, &builders.ThreatIntelligencePoliciesItemPolicyRulesRequestBuilderPostRequestConfiguration{
			// Kiota 1.5.6 can retain a gzip header while replaying the original JSON body.
			// Use its standard compression option so existing retries replay valid requests.
			Options: []abstractions.RequestOption{kiotahttp.NewCompressionOptionsReference(false)},
		})
}

//nolint:wrapcheck // Shared Graph diagnostics require the original SDK error type.
func (r *NetworkThreatIntelligencePolicyRuleResource) getThreatIntelligencePolicyRule(
	ctx context.Context,
	parent, id string,
) (models.PolicyRuleable, error) {
	return r.client.NetworkAccess().
		ThreatIntelligencePolicies().
		ByThreatIntelligencePolicyId(parent).
		PolicyRules().
		ByPolicyRuleId(id).
		Get(ctx, nil)
}

//nolint:wrapcheck // Shared Graph diagnostics require the original SDK error type.
func (r *NetworkThreatIntelligencePolicyRuleResource) updateThreatIntelligencePolicyRule(
	ctx context.Context,
	parent, id string,
	body models.PolicyRuleable,
) error {
	_, err := r.client.NetworkAccess().
		ThreatIntelligencePolicies().
		ByThreatIntelligencePolicyId(parent).
		PolicyRules().
		ByPolicyRuleId(id).
		Patch(ctx, body, &builders.ThreatIntelligencePoliciesItemPolicyRulesPolicyRuleItemRequestBuilderPatchRequestConfiguration{
			// Kiota 1.5.6 can retain a gzip header while replaying the original JSON body.
			// Use its standard compression option so existing retries replay valid requests.
			Options: []abstractions.RequestOption{kiotahttp.NewCompressionOptionsReference(false)},
		})
	return err
}

//nolint:wrapcheck // Shared Graph diagnostics require the original SDK error type.
func (r *NetworkThreatIntelligencePolicyRuleResource) deleteThreatIntelligencePolicyRule(
	ctx context.Context,
	parent, id string,
) error {
	return r.client.NetworkAccess().
		ThreatIntelligencePolicies().
		ByThreatIntelligencePolicyId(parent).
		PolicyRules().
		ByPolicyRuleId(id).
		Delete(ctx, nil)
}
