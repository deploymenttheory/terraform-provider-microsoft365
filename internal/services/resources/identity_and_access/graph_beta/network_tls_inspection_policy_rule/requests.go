//nolint:wrapcheck // Preserve generated Kiota API errors for the shared Graph error handler.
package graphBetaNetworkTLSInspectionPolicyRule

import (
	"context"
	"fmt"

	abstractions "github.com/microsoft/kiota-abstractions-go"
	models "github.com/microsoftgraph/msgraph-beta-sdk-go/models/networkaccess"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models/odataerrors"
)

var tlsInspectionRuleErrorMapping = abstractions.ErrorMappings{
	"XXX": odataerrors.CreateODataErrorFromDiscriminatorValue,
}

func (r *NetworkTLSInspectionPolicyRuleResource) createTLSInspectionPolicyRule(
	ctx context.Context,
	policyID string,
	body models.TlsInspectionRuleable,
) (models.TlsInspectionRuleable, error) {
	builder := r.client.NetworkAccess().
		TlsInspectionPolicies().
		ByTlsInspectionPolicyId(policyID).
		PolicyRules()
	requestInfo, err := builder.ToPostRequestInformation(ctx, body, nil)
	if err != nil {
		return nil, err
	}
	result, err := r.client.GetAdapter().Send(
		ctx,
		requestInfo,
		models.CreateTlsInspectionRuleFromDiscriminatorValue,
		tlsInspectionRuleErrorMapping,
	)
	if err != nil {
		return nil, err
	}
	return tlsInspectionRule(result)
}

func (r *NetworkTLSInspectionPolicyRuleResource) getTLSInspectionPolicyRule(
	ctx context.Context,
	policyID, ruleID string,
) (models.TlsInspectionRuleable, error) {
	builder := r.client.NetworkAccess().
		TlsInspectionPolicies().
		ByTlsInspectionPolicyId(policyID).
		PolicyRules().
		ByPolicyRuleId(ruleID)
	requestInfo, err := builder.ToGetRequestInformation(ctx, nil)
	if err != nil {
		return nil, err
	}
	// Graph can omit @odata.type on rule responses. The generated builder then
	// falls back to PolicyRule and drops TLS-specific fields, so force its
	// generated concrete model factory while retaining generated request paths.
	result, err := r.client.GetAdapter().Send(
		ctx,
		requestInfo,
		models.CreateTlsInspectionRuleFromDiscriminatorValue,
		tlsInspectionRuleErrorMapping,
	)
	if err != nil {
		return nil, err
	}
	return tlsInspectionRule(result)
}

func (r *NetworkTLSInspectionPolicyRuleResource) updateTLSInspectionPolicyRule(
	ctx context.Context,
	policyID, ruleID string,
	body models.TlsInspectionRuleable,
) error {
	_, err := r.client.NetworkAccess().
		TlsInspectionPolicies().
		ByTlsInspectionPolicyId(policyID).
		PolicyRules().
		ByPolicyRuleId(ruleID).
		Patch(ctx, body, nil)
	return err
}

func (r *NetworkTLSInspectionPolicyRuleResource) deleteTLSInspectionPolicyRule(
	ctx context.Context,
	policyID, ruleID string,
) error {
	return r.client.NetworkAccess().
		TlsInspectionPolicies().
		ByTlsInspectionPolicyId(policyID).
		PolicyRules().
		ByPolicyRuleId(ruleID).
		Delete(ctx, nil)
}

func tlsInspectionRule(result any) (models.TlsInspectionRuleable, error) {
	if result == nil {
		return nil, errEmptyResponse
	}
	rule, ok := result.(models.TlsInspectionRuleable)
	if !ok {
		return nil, fmt.Errorf("%w: received %T", errInvalidResponse, result)
	}
	return rule, nil
}
