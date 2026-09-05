//nolint:wrapcheck // Preserve Kiota API error types for the shared Graph HTTP status and permission handler.
package graphBetaNetworkTLSInspectionPolicyRule

import (
	"context"
	"fmt"

	abstractions "github.com/microsoft/kiota-abstractions-go"
	s "github.com/microsoft/kiota-abstractions-go/serialization"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models/odataerrors"
)

const (
	tlsInspectionPolicyRulesURLTemplate    = "{+baseurl}/networkaccess/tlsInspectionPolicies/{tlsInspectionPolicyId}/policyRules"
	tlsInspectionPolicyRuleItemURLTemplate = tlsInspectionPolicyRulesURLTemplate + "/{policyRuleId}"
)

var tlsInspectionPolicyRuleErrorMapping = abstractions.ErrorMappings{
	"XXX": odataerrors.CreateODataErrorFromDiscriminatorValue,
}

func (r *NetworkTLSInspectionPolicyRuleResource) createTLSInspectionPolicyRule(
	ctx context.Context,
	policyID string,
	body s.Parsable,
) (*tlsInspectionPolicyRuleResponse, error) {
	return r.sendTLSInspectionPolicyRule(ctx, abstractions.POST, policyID, "", body)
}

func (r *NetworkTLSInspectionPolicyRuleResource) getTLSInspectionPolicyRule(
	ctx context.Context,
	policyID, ruleID string,
) (*tlsInspectionPolicyRuleResponse, error) {
	return r.sendTLSInspectionPolicyRule(ctx, abstractions.GET, policyID, ruleID, nil)
}

func (r *NetworkTLSInspectionPolicyRuleResource) updateTLSInspectionPolicyRule(
	ctx context.Context,
	policyID, ruleID string,
	body s.Parsable,
) error {
	adapter := r.client.GetAdapter()
	requestInfo, err := newTLSInspectionPolicyRuleRequestInformation(
		ctx,
		adapter,
		abstractions.PATCH,
		policyID,
		ruleID,
		body,
	)
	if err != nil {
		return err
	}
	return adapter.SendNoContent(ctx, requestInfo, tlsInspectionPolicyRuleErrorMapping)
}

func (r *NetworkTLSInspectionPolicyRuleResource) deleteTLSInspectionPolicyRule(
	ctx context.Context,
	policyID, ruleID string,
) error {
	adapter := r.client.GetAdapter()
	requestInfo, err := newTLSInspectionPolicyRuleRequestInformation(
		ctx,
		adapter,
		abstractions.DELETE,
		policyID,
		ruleID,
		nil,
	)
	if err != nil {
		return err
	}
	return adapter.SendNoContent(ctx, requestInfo, tlsInspectionPolicyRuleErrorMapping)
}

func (r *NetworkTLSInspectionPolicyRuleResource) sendTLSInspectionPolicyRule(
	ctx context.Context,
	method abstractions.HttpMethod,
	policyID, ruleID string,
	body s.Parsable,
) (*tlsInspectionPolicyRuleResponse, error) {
	adapter := r.client.GetAdapter()
	requestInfo, err := newTLSInspectionPolicyRuleRequestInformation(
		ctx,
		adapter,
		method,
		policyID,
		ruleID,
		body,
	)
	if err != nil {
		return nil, err
	}
	result, err := adapter.Send(
		ctx,
		requestInfo,
		createTLSInspectionPolicyRuleResponseFromDiscriminatorValue,
		tlsInspectionPolicyRuleErrorMapping,
	)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, errEmptyResponse
	}
	rule, ok := result.(*tlsInspectionPolicyRuleResponse)
	if !ok {
		return nil, fmt.Errorf("%w: received %T", errInvalidResponse, result)
	}
	return rule, nil
}

func newTLSInspectionPolicyRuleRequestInformation(
	ctx context.Context,
	adapter abstractions.RequestAdapter,
	method abstractions.HttpMethod,
	policyID, ruleID string,
	body s.Parsable,
) (*abstractions.RequestInformation, error) {
	pathParameters := map[string]string{
		"baseurl":               adapter.GetBaseUrl(),
		"tlsInspectionPolicyId": policyID,
	}
	urlTemplate := tlsInspectionPolicyRulesURLTemplate
	if ruleID != "" {
		urlTemplate = tlsInspectionPolicyRuleItemURLTemplate
		pathParameters["policyRuleId"] = ruleID
	}
	requestInfo := abstractions.NewRequestInformationWithMethodAndUrlTemplateAndPathParameters(
		method,
		urlTemplate,
		pathParameters,
	)
	requestInfo.Headers.TryAdd("Accept", "application/json")
	if body != nil {
		if err := requestInfo.SetContentFromParsable(
			ctx,
			adapter,
			"application/json",
			body,
		); err != nil {
			return nil, fmt.Errorf("set TLS inspection policy rule request content: %w", err)
		}
	}
	return requestInfo, nil
}
