package graphBetaNetworkContentPolicyRule

import (
	"context"
	"fmt"

	abstractions "github.com/microsoft/kiota-abstractions-go"
	s "github.com/microsoft/kiota-abstractions-go/serialization"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models/odataerrors"
)

const (
	contentPolicyRulesURLTemplate    = "{+baseurl}/networkaccess/filePolicies/{contentPolicyId}/policyRules"
	contentPolicyRuleItemURLTemplate = contentPolicyRulesURLTemplate + "/{policyRuleId}"
)

var contentPolicyRuleErrorMapping = abstractions.ErrorMappings{"XXX": odataerrors.CreateODataErrorFromDiscriminatorValue}

func (r *NetworkContentPolicyRuleResource) createContentPolicyRule(ctx context.Context, policyID string, body s.Parsable) (*contentPolicyRuleResponse, error) {
	return r.sendContentPolicyRule(ctx, abstractions.POST, policyID, "", body)
}

func (r *NetworkContentPolicyRuleResource) getContentPolicyRule(ctx context.Context, policyID, ruleID string) (*contentPolicyRuleResponse, error) {
	return r.sendContentPolicyRule(ctx, abstractions.GET, policyID, ruleID, nil)
}

func (r *NetworkContentPolicyRuleResource) updateContentPolicyRule(ctx context.Context, policyID, ruleID string, body s.Parsable) error {
	result, err := r.sendContentPolicyRule(ctx, abstractions.PATCH, policyID, ruleID, body)
	if err != nil {
		return err
	}
	if result == nil {
		return fmt.Errorf("update content policy rule returned nil response")
	}
	return nil
}

func (r *NetworkContentPolicyRuleResource) deleteContentPolicyRule(ctx context.Context, policyID, ruleID string) error {
	adapter := r.client.GetAdapter()
	requestInfo, err := newContentPolicyRuleRequestInformation(ctx, adapter, abstractions.DELETE, policyID, ruleID, nil)
	if err != nil {
		return err
	}
	return adapter.SendNoContent(ctx, requestInfo, contentPolicyRuleErrorMapping)
}

func (r *NetworkContentPolicyRuleResource) sendContentPolicyRule(ctx context.Context, method abstractions.HttpMethod, policyID, ruleID string, body s.Parsable) (*contentPolicyRuleResponse, error) {
	adapter := r.client.GetAdapter()
	requestInfo, err := newContentPolicyRuleRequestInformation(ctx, adapter, method, policyID, ruleID, body)
	if err != nil {
		return nil, err
	}
	result, err := adapter.Send(ctx, requestInfo, createContentPolicyRuleResponseFromDiscriminatorValue, contentPolicyRuleErrorMapping)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, nil
	}
	rule, ok := result.(*contentPolicyRuleResponse)
	if !ok {
		return nil, fmt.Errorf("content policy rule request returned %T, expected contentPolicyRuleResponse", result)
	}
	return rule, nil
}

func newContentPolicyRuleRequestInformation(ctx context.Context, adapter abstractions.RequestAdapter, method abstractions.HttpMethod, policyID, ruleID string, body s.Parsable) (*abstractions.RequestInformation, error) {
	pathParameters := map[string]string{
		"baseurl":         adapter.GetBaseUrl(),
		"contentPolicyId": policyID,
	}
	urlTemplate := contentPolicyRulesURLTemplate
	if ruleID != "" {
		urlTemplate = contentPolicyRuleItemURLTemplate
		pathParameters["policyRuleId"] = ruleID
	}
	requestInfo := abstractions.NewRequestInformationWithMethodAndUrlTemplateAndPathParameters(method, urlTemplate, pathParameters)
	requestInfo.Headers.TryAdd("Accept", "application/json")
	if body != nil {
		if err := requestInfo.SetContentFromParsable(ctx, adapter, "application/json", body); err != nil {
			return nil, fmt.Errorf("set content policy rule request content: %w", err)
		}
	}
	return requestInfo, nil
}
