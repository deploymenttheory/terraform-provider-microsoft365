package graphBetaNetworkThreatIntelligencePolicy

import (
	"context"

	abstractions "github.com/microsoft/kiota-abstractions-go"
	kiotahttp "github.com/microsoft/kiota-http-go"
	models "github.com/microsoftgraph/msgraph-beta-sdk-go/models/networkaccess"
	builders "github.com/microsoftgraph/msgraph-beta-sdk-go/networkaccess"
)

//nolint:wrapcheck // Shared Graph diagnostics require the original SDK error type.
func (r *NetworkThreatIntelligencePolicyResource) createThreatIntelligencePolicy(
	ctx context.Context,
	body models.ThreatIntelligencePolicyable,
) (models.ThreatIntelligencePolicyable, error) {
	return r.client.NetworkAccess().
		ThreatIntelligencePolicies().
		Post(ctx, body, &builders.ThreatIntelligencePoliciesRequestBuilderPostRequestConfiguration{
			// Kiota 1.5.6 can retain a gzip header while replaying the original JSON body.
			// Use its standard compression option so existing retries replay valid requests.
			Options: []abstractions.RequestOption{kiotahttp.NewCompressionOptionsReference(false)},
		})
}

//nolint:wrapcheck // Shared Graph diagnostics require the original SDK error type.
func (r *NetworkThreatIntelligencePolicyResource) getThreatIntelligencePolicy(
	ctx context.Context,
	id string,
) (models.ThreatIntelligencePolicyable, error) {
	return r.client.NetworkAccess().
		ThreatIntelligencePolicies().
		ByThreatIntelligencePolicyId(id).
		Get(ctx, nil)
}

//nolint:wrapcheck // Shared Graph diagnostics require the original SDK error type.
func (r *NetworkThreatIntelligencePolicyResource) updateThreatIntelligencePolicy(
	ctx context.Context,
	id string,
	body models.ThreatIntelligencePolicyable,
) error {
	_, err := r.client.NetworkAccess().
		ThreatIntelligencePolicies().
		ByThreatIntelligencePolicyId(id).
		Patch(ctx, body, &builders.ThreatIntelligencePoliciesThreatIntelligencePolicyItemRequestBuilderPatchRequestConfiguration{
			// Kiota 1.5.6 can retain a gzip header while replaying the original JSON body.
			// Use its standard compression option so existing retries replay valid requests.
			Options: []abstractions.RequestOption{kiotahttp.NewCompressionOptionsReference(false)},
		})
	return err
}

//nolint:wrapcheck // Shared Graph diagnostics require the original SDK error type.
func (r *NetworkThreatIntelligencePolicyResource) deleteThreatIntelligencePolicy(
	ctx context.Context,
	id string,
) error {
	return r.client.NetworkAccess().
		ThreatIntelligencePolicies().
		ByThreatIntelligencePolicyId(id).
		Delete(ctx, nil)
}
