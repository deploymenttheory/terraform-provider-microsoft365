//nolint:wrapcheck // Preserve Kiota API error types for the shared Graph HTTP status and permission handler.
package graphBetaNetworkMCPPolicyRule

import (
	"context"
	"encoding/json"
	"fmt"
	kiotahttp "github.com/microsoft/kiota-http-go"
	"net/http"
	"time"

	abstractions "github.com/microsoft/kiota-abstractions-go"
	s "github.com/microsoft/kiota-abstractions-go/serialization"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models/odataerrors"
)

const (
	mcpPolicyRulesURLTemplate    = "{+baseurl}/networkaccess/mcpPolicies/{mcpPolicyId}/policyRules"
	mcpPolicyRuleItemURLTemplate = mcpPolicyRulesURLTemplate + "/{policyRuleId}"
)

var mcpPolicyRuleErrorMapping = abstractions.ErrorMappings{
	"XXX": odataerrors.CreateODataErrorFromDiscriminatorValue,
}

func (r *NetworkMCPPolicyRuleResource) createMCPPolicyRule(
	ctx context.Context,
	policyID string,
	body s.Parsable,
) (*mcpPolicyRuleResponse, error) {
	return r.sendMCPPolicyRule(ctx, abstractions.POST, policyID, "", body)
}

func (r *NetworkMCPPolicyRuleResource) getMCPPolicyRule(
	ctx context.Context,
	policyID, ruleID string,
) (*mcpPolicyRuleResponse, error) {
	return r.sendMCPPolicyRule(ctx, abstractions.GET, policyID, ruleID, nil)
}

func (r *NetworkMCPPolicyRuleResource) updateMCPPolicyRule(
	ctx context.Context,
	policyID, ruleID string,
	body s.Parsable,
) error {
	_, err := r.sendMCPPolicyRule(ctx, abstractions.PATCH, policyID, ruleID, body)
	return err
}

func (r *NetworkMCPPolicyRuleResource) deleteMCPPolicyRule(
	ctx context.Context,
	policyID, ruleID string,
) error {
	adapter := r.client.GetAdapter()
	requestInfo, err := newMCPPolicyRuleRequestInformation(
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
	return adapter.SendNoContent(ctx, requestInfo, mcpPolicyRuleErrorMapping)
}

func (r *NetworkMCPPolicyRuleResource) sendMCPPolicyRule(
	ctx context.Context,
	method abstractions.HttpMethod,
	policyID, ruleID string,
	body s.Parsable,
) (*mcpPolicyRuleResponse, error) {
	adapter := r.client.GetAdapter()
	requestInfo, err := newMCPPolicyRuleRequestInformation(
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
	// Keep Kiota authentication, middleware and Graph error mapping, but decode the
	// unmodeled MCP JSON ourselves instead of dropping derived condition fields.
	result, err := adapter.SendPrimitive(ctx, requestInfo, "[]byte", mcpPolicyRuleErrorMapping)
	if err != nil {
		return nil, err
	}
	raw, ok := result.([]byte)
	if !ok || len(raw) == 0 {
		return nil, errEmptyResponse
	}
	var rule *mcpPolicyRuleResponse
	if err := json.Unmarshal(raw, &rule); err != nil {
		return nil, fmt.Errorf("%w: %v", errInvalidResponse, err)
	}
	if rule == nil {
		return nil, errEmptyResponse
	}
	return rule, nil
}

func newMCPPolicyRuleRequestInformation(
	ctx context.Context,
	adapter abstractions.RequestAdapter,
	method abstractions.HttpMethod,
	policyID, ruleID string,
	body s.Parsable,
) (*abstractions.RequestInformation, error) {
	pathParameters := map[string]string{
		"baseurl":     adapter.GetBaseUrl(),
		"mcpPolicyId": policyID,
	}
	urlTemplate := mcpPolicyRulesURLTemplate
	if ruleID != "" {
		urlTemplate = mcpPolicyRuleItemURLTemplate
		pathParameters["policyRuleId"] = ruleID
	}
	requestInfo := abstractions.NewRequestInformationWithMethodAndUrlTemplateAndPathParameters(
		method,
		urlTemplate,
		pathParameters,
	)
	if method == abstractions.POST {
		// Creation has no idempotency key. Do not replay it after an ambiguous service response.
		requestInfo.AddRequestOptions([]abstractions.RequestOption{&kiotahttp.RetryHandlerOptions{ShouldRetry: func(time.Duration, int, *http.Request, *http.Response) bool { return false }}})
	}
	requestInfo.Headers.TryAdd("Accept", "application/json")
	if body != nil {
		if err := requestInfo.SetContentFromParsable(
			ctx,
			adapter,
			"application/json",
			body,
		); err != nil {
			return nil, fmt.Errorf("set MCP policy rule request content: %w", err)
		}
	}
	return requestInfo, nil
}
