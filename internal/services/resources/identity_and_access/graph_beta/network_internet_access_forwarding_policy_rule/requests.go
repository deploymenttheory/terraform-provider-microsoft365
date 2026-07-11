package graphBetaNetworkInternetAccessForwardingPolicyRule

import (
	"context"
	"fmt"

	abstractions "github.com/microsoft/kiota-abstractions-go"
	s "github.com/microsoft/kiota-abstractions-go/serialization"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models/odataerrors"
)

const (
	forwardingPolicyRulesURLTemplate    = "{+baseurl}/networkAccess/forwardingPolicies/{forwardingPolicyId}/policyRules"
	forwardingPolicyRuleItemURLTemplate = forwardingPolicyRulesURLTemplate + "/{policyRuleId}"
)

var internetAccessForwardingRuleErrorMapping = abstractions.ErrorMappings{
	"XXX": odataerrors.CreateODataErrorFromDiscriminatorValue,
}

func (r *NetworkInternetAccessForwardingPolicyRuleResource) createRule(ctx context.Context, forwardingPolicyID string, requestBody s.Parsable) (*internetAccessForwardingRuleResponse, error) {
	adapter := r.client.GetAdapter()
	requestInfo, err := newInternetAccessForwardingRuleRequestInformation(ctx, adapter, abstractions.POST, forwardingPolicyID, "", requestBody)
	if err != nil {
		return nil, err
	}
	result, err := adapter.Send(ctx, requestInfo, createInternetAccessForwardingRuleResponseFromDiscriminatorValue, internetAccessForwardingRuleErrorMapping)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, fmt.Errorf("create internet access forwarding policy rule returned nil response")
	}
	rule, ok := result.(*internetAccessForwardingRuleResponse)
	if !ok {
		return nil, fmt.Errorf("create internet access forwarding policy rule returned %T, expected internetAccessForwardingRuleResponse", result)
	}
	return rule, nil
}

func (r *NetworkInternetAccessForwardingPolicyRuleResource) getRule(ctx context.Context, forwardingPolicyID, ruleID string) (*internetAccessForwardingRuleResponse, error) {
	adapter := r.client.GetAdapter()
	requestInfo, err := newInternetAccessForwardingRuleRequestInformation(ctx, adapter, abstractions.GET, forwardingPolicyID, ruleID, nil)
	if err != nil {
		return nil, err
	}
	result, err := adapter.Send(ctx, requestInfo, createInternetAccessForwardingRuleResponseFromDiscriminatorValue, internetAccessForwardingRuleErrorMapping)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, nil
	}
	rule, ok := result.(*internetAccessForwardingRuleResponse)
	if !ok {
		return nil, fmt.Errorf("get internet access forwarding policy rule returned %T, expected internetAccessForwardingRuleResponse", result)
	}
	return rule, nil
}

func (r *NetworkInternetAccessForwardingPolicyRuleResource) updateRule(ctx context.Context, forwardingPolicyID, ruleID string, requestBody s.Parsable) error {
	adapter := r.client.GetAdapter()
	requestInfo, err := newInternetAccessForwardingRuleRequestInformation(ctx, adapter, abstractions.PATCH, forwardingPolicyID, ruleID, requestBody)
	if err != nil {
		return err
	}
	return adapter.SendNoContent(ctx, requestInfo, internetAccessForwardingRuleErrorMapping)
}

func (r *NetworkInternetAccessForwardingPolicyRuleResource) deleteRule(ctx context.Context, forwardingPolicyID, ruleID string) error {
	adapter := r.client.GetAdapter()
	requestInfo, err := newInternetAccessForwardingRuleRequestInformation(ctx, adapter, abstractions.DELETE, forwardingPolicyID, ruleID, nil)
	if err != nil {
		return err
	}
	return adapter.SendNoContent(ctx, requestInfo, internetAccessForwardingRuleErrorMapping)
}

func newInternetAccessForwardingRuleRequestInformation(ctx context.Context, adapter abstractions.RequestAdapter, method abstractions.HttpMethod, forwardingPolicyID, ruleID string, requestBody s.Parsable) (*abstractions.RequestInformation, error) {
	pathParameters := map[string]string{
		"baseurl":            adapter.GetBaseUrl(),
		"forwardingPolicyId": forwardingPolicyID,
	}
	urlTemplate := forwardingPolicyRulesURLTemplate
	if ruleID != "" {
		urlTemplate = forwardingPolicyRuleItemURLTemplate
		pathParameters["policyRuleId"] = ruleID
	}

	requestInfo := abstractions.NewRequestInformationWithMethodAndUrlTemplateAndPathParameters(method, urlTemplate, pathParameters)
	requestInfo.Headers.TryAdd("Accept", "application/json")

	if requestBody != nil {
		if err := requestInfo.SetContentFromParsable(ctx, adapter, "application/json", requestBody); err != nil {
			return nil, fmt.Errorf("set internet access forwarding policy rule request content: %w", err)
		}
	}
	return requestInfo, nil
}
