package graphBetaNetworkWebFilteringPolicyRule

import (
	"context"
	"fmt"

	abstractions "github.com/microsoft/kiota-abstractions-go"
	s "github.com/microsoft/kiota-abstractions-go/serialization"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models/odataerrors"
)

const (
	// The webFilteringPolicy policyRules endpoint is used by the Entra portal but
	// is not represented by generated Microsoft Graph beta SDK request builders.
	// Keep these templates aligned with observed portal traffic rather than the
	// documented generic filteringPolicy/policyRules endpoint:
	// https://learn.microsoft.com/graph/api/networkaccess-filteringpolicy-post-policyrules
	webFilteringPolicyRulesURLTemplate    = "{+baseurl}/networkaccess/webFilteringPolicies/{webFilteringPolicyId}/policyRules"
	webFilteringPolicyRuleItemURLTemplate = webFilteringPolicyRulesURLTemplate + "/{policyRuleId}"
)

var webFilteringPolicyRuleErrorMapping = abstractions.ErrorMappings{
	"XXX": odataerrors.CreateODataErrorFromDiscriminatorValue,
}

func (r *NetworkWebFilteringPolicyRuleResource) createWebFilteringPolicyRule(ctx context.Context, policyID string, requestBody s.Parsable) (*webFilteringPolicyRuleResponse, error) {
	adapter := r.client.GetAdapter()
	requestInfo, err := newWebFilteringPolicyRuleRequestInformation(ctx, adapter, abstractions.POST, policyID, "", requestBody)
	if err != nil {
		return nil, err
	}

	result, err := adapter.Send(ctx, requestInfo, createWebFilteringPolicyRuleResponseFromDiscriminatorValue, webFilteringPolicyRuleErrorMapping)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, fmt.Errorf("create web filtering policy rule returned nil response")
	}

	rule, ok := result.(*webFilteringPolicyRuleResponse)
	if !ok {
		return nil, fmt.Errorf("create web filtering policy rule returned %T, expected webFilteringPolicyRuleResponse", result)
	}

	return rule, nil
}

func (r *NetworkWebFilteringPolicyRuleResource) getWebFilteringPolicyRule(ctx context.Context, policyID, ruleID string) (*webFilteringPolicyRuleResponse, error) {
	adapter := r.client.GetAdapter()
	requestInfo, err := newWebFilteringPolicyRuleRequestInformation(ctx, adapter, abstractions.GET, policyID, ruleID, nil)
	if err != nil {
		return nil, err
	}

	result, err := adapter.Send(ctx, requestInfo, createWebFilteringPolicyRuleResponseFromDiscriminatorValue, webFilteringPolicyRuleErrorMapping)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, nil
	}

	rule, ok := result.(*webFilteringPolicyRuleResponse)
	if !ok {
		return nil, fmt.Errorf("get web filtering policy rule returned %T, expected webFilteringPolicyRuleResponse", result)
	}

	return rule, nil
}

func (r *NetworkWebFilteringPolicyRuleResource) updateWebFilteringPolicyRule(ctx context.Context, policyID, ruleID string, requestBody s.Parsable) error {
	adapter := r.client.GetAdapter()
	requestInfo, err := newWebFilteringPolicyRuleRequestInformation(ctx, adapter, abstractions.PATCH, policyID, ruleID, requestBody)
	if err != nil {
		return err
	}

	return adapter.SendNoContent(ctx, requestInfo, webFilteringPolicyRuleErrorMapping)
}

func (r *NetworkWebFilteringPolicyRuleResource) deleteWebFilteringPolicyRule(ctx context.Context, policyID, ruleID string) error {
	adapter := r.client.GetAdapter()
	requestInfo, err := newWebFilteringPolicyRuleRequestInformation(ctx, adapter, abstractions.DELETE, policyID, ruleID, nil)
	if err != nil {
		return err
	}

	return adapter.SendNoContent(ctx, requestInfo, webFilteringPolicyRuleErrorMapping)
}

func newWebFilteringPolicyRuleRequestInformation(ctx context.Context, adapter abstractions.RequestAdapter, method abstractions.HttpMethod, policyID, ruleID string, requestBody s.Parsable) (*abstractions.RequestInformation, error) {
	pathParameters := map[string]string{
		"baseurl":              adapter.GetBaseUrl(),
		"webFilteringPolicyId": policyID,
	}

	urlTemplate := webFilteringPolicyRulesURLTemplate
	if ruleID != "" {
		urlTemplate = webFilteringPolicyRuleItemURLTemplate
		pathParameters["policyRuleId"] = ruleID
	}

	requestInfo := abstractions.NewRequestInformationWithMethodAndUrlTemplateAndPathParameters(method, urlTemplate, pathParameters)
	requestInfo.Headers.TryAdd("Accept", "application/json")

	if requestBody != nil {
		if err := requestInfo.SetContentFromParsable(ctx, adapter, "application/json", requestBody); err != nil {
			return nil, fmt.Errorf("set web filtering policy rule request content: %w", err)
		}
	}

	return requestInfo, nil
}
