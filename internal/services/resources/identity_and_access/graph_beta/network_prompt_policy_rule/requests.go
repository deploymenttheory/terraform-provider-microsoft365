//nolint:wrapcheck // Preserve Kiota API error types for the shared Graph HTTP status and permission handler.
package graphBetaNetworkPromptPolicyRule

import (
	"context"
	"fmt"

	abstractions "github.com/microsoft/kiota-abstractions-go"
	s "github.com/microsoft/kiota-abstractions-go/serialization"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models/odataerrors"
)

const (
	promptPolicyRulesURLTemplate    = "{+baseurl}/networkaccess/promptPolicies/{promptPolicyId}/policyRules"
	promptPolicyRuleItemURLTemplate = promptPolicyRulesURLTemplate + "/{policyRuleId}"
)

var promptPolicyRuleErrorMapping = abstractions.ErrorMappings{
	"XXX": odataerrors.CreateODataErrorFromDiscriminatorValue,
}

func (r *NetworkPromptPolicyRuleResource) createPromptPolicyRule(
	ctx context.Context,
	policyID string,
	body s.Parsable,
) (*promptPolicyRuleResponse, error) {
	return r.sendPromptPolicyRule(ctx, abstractions.POST, policyID, "", body)
}

func (r *NetworkPromptPolicyRuleResource) getPromptPolicyRule(
	ctx context.Context,
	policyID, ruleID string,
) (*promptPolicyRuleResponse, error) {
	return r.sendPromptPolicyRule(ctx, abstractions.GET, policyID, ruleID, nil)
}

func (r *NetworkPromptPolicyRuleResource) updatePromptPolicyRule(
	ctx context.Context,
	policyID, ruleID string,
	body s.Parsable,
) error {
	adapter := r.client.GetAdapter()
	requestInfo, err := newPromptPolicyRuleRequestInformation(
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
	_, err = adapter.Send(
		ctx,
		requestInfo,
		createPromptPolicyRuleResponseFromDiscriminatorValue,
		promptPolicyRuleErrorMapping,
	)
	return err
}

func (r *NetworkPromptPolicyRuleResource) deletePromptPolicyRule(
	ctx context.Context,
	policyID, ruleID string,
) error {
	adapter := r.client.GetAdapter()
	requestInfo, err := newPromptPolicyRuleRequestInformation(
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
	return adapter.SendNoContent(ctx, requestInfo, promptPolicyRuleErrorMapping)
}

func (r *NetworkPromptPolicyRuleResource) sendPromptPolicyRule(
	ctx context.Context,
	method abstractions.HttpMethod,
	policyID, ruleID string,
	body s.Parsable,
) (*promptPolicyRuleResponse, error) {
	adapter := r.client.GetAdapter()
	requestInfo, err := newPromptPolicyRuleRequestInformation(
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
		createPromptPolicyRuleResponseFromDiscriminatorValue,
		promptPolicyRuleErrorMapping,
	)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, errEmptyResponse
	}
	rule, ok := result.(*promptPolicyRuleResponse)
	if !ok {
		return nil, fmt.Errorf("%w: received %T", errInvalidResponse, result)
	}
	return rule, nil
}

func newPromptPolicyRuleRequestInformation(
	ctx context.Context,
	adapter abstractions.RequestAdapter,
	method abstractions.HttpMethod,
	policyID, ruleID string,
	body s.Parsable,
) (*abstractions.RequestInformation, error) {
	pathParameters := map[string]string{
		"baseurl":        adapter.GetBaseUrl(),
		"promptPolicyId": policyID,
	}
	urlTemplate := promptPolicyRulesURLTemplate
	if ruleID != "" {
		urlTemplate = promptPolicyRuleItemURLTemplate
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
			return nil, fmt.Errorf("set prompt policy rule request content: %w", err)
		}
	}
	return requestInfo, nil
}
