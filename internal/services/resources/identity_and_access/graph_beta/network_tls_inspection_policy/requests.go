//nolint:wrapcheck // Preserve generated Kiota API errors for the shared Graph error handler.
package graphBetaNetworkTLSInspectionPolicy

import (
	"context"

	models "github.com/microsoftgraph/msgraph-beta-sdk-go/models/networkaccess"
)

func (r *NetworkTLSInspectionPolicyResource) createTLSInspectionPolicy(
	ctx context.Context,
	body models.TlsInspectionPolicyable,
) (models.TlsInspectionPolicyable, error) {
	result, err := r.client.NetworkAccess().TlsInspectionPolicies().Post(ctx, body, nil)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, errEmptyResponse
	}
	return result, nil
}

func (r *NetworkTLSInspectionPolicyResource) getTLSInspectionPolicy(
	ctx context.Context,
	policyID string,
) (models.TlsInspectionPolicyable, error) {
	result, err := r.client.NetworkAccess().TlsInspectionPolicies().
		ByTlsInspectionPolicyId(policyID).
		Get(ctx, nil)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, errEmptyResponse
	}
	return result, nil
}

func (r *NetworkTLSInspectionPolicyResource) updateTLSInspectionPolicy(
	ctx context.Context,
	policyID string,
	body models.TlsInspectionPolicyable,
) error {
	_, err := r.client.NetworkAccess().TlsInspectionPolicies().
		ByTlsInspectionPolicyId(policyID).
		Patch(ctx, body, nil)
	return err
}

func (r *NetworkTLSInspectionPolicyResource) deleteTLSInspectionPolicy(
	ctx context.Context,
	policyID string,
) error {
	return r.client.NetworkAccess().TlsInspectionPolicies().
		ByTlsInspectionPolicyId(policyID).
		Delete(ctx, nil)
}
