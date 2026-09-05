//nolint:wrapcheck // Preserve Kiota API error types for the shared Graph HTTP status and permission handler.
package graphBetaNetworkMCPPolicy

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
	mcpPoliciesURLTemplate   = "{+baseurl}/networkaccess/mcpPolicies"
	mcpPolicyItemURLTemplate = mcpPoliciesURLTemplate + "/{mcpPolicyId}"
)

var mcpPolicyErrorMapping = abstractions.ErrorMappings{
	"XXX": odataerrors.CreateODataErrorFromDiscriminatorValue,
}

func (r *NetworkMCPPolicyResource) createMCPPolicy(
	ctx context.Context,
	requestBody s.Parsable,
) (*mcpPolicyResponse, error) {
	result, err := r.sendMCPPolicy(ctx, abstractions.POST, "", requestBody)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, errEmptyResponse
	}
	return result, nil
}

func (r *NetworkMCPPolicyResource) getMCPPolicy(
	ctx context.Context,
	policyID string,
) (*mcpPolicyResponse, error) {
	return r.sendMCPPolicy(ctx, abstractions.GET, policyID, nil)
}

func (r *NetworkMCPPolicyResource) sendMCPPolicy(
	ctx context.Context,
	method abstractions.HttpMethod,
	policyID string,
	requestBody s.Parsable,
) (*mcpPolicyResponse, error) {
	adapter := r.client.GetAdapter()
	requestInfo, err := newMCPPolicyRequestInformation(
		ctx,
		adapter,
		method,
		policyID,
		requestBody,
	)
	if err != nil {
		return nil, err
	}
	result, err := adapter.SendPrimitive(ctx, requestInfo, "[]byte", mcpPolicyErrorMapping)
	if err != nil {
		return nil, err
	}
	raw, ok := result.([]byte)
	if !ok || len(raw) == 0 {
		return nil, errEmptyResponse
	}
	var policy *mcpPolicyResponse
	if err := json.Unmarshal(raw, &policy); err != nil {
		return nil, fmt.Errorf("%w: %v", errInvalidResponse, err)
	}
	if policy == nil {
		return nil, errEmptyResponse
	}
	return policy, nil
}

func (r *NetworkMCPPolicyResource) updateMCPPolicy(
	ctx context.Context,
	policyID string,
	requestBody s.Parsable,
) error {
	_, err := r.sendMCPPolicy(ctx, abstractions.PATCH, policyID, requestBody)
	return err
}

func (r *NetworkMCPPolicyResource) deleteMCPPolicy(
	ctx context.Context,
	policyID string,
) error {
	adapter := r.client.GetAdapter()
	requestInfo, err := newMCPPolicyRequestInformation(
		ctx,
		adapter,
		abstractions.DELETE,
		policyID,
		nil,
	)
	if err != nil {
		return err
	}
	return adapter.SendNoContent(ctx, requestInfo, mcpPolicyErrorMapping)
}

func newMCPPolicyRequestInformation(
	ctx context.Context,
	adapter abstractions.RequestAdapter,
	method abstractions.HttpMethod,
	policyID string,
	requestBody s.Parsable,
) (*abstractions.RequestInformation, error) {
	pathParameters := map[string]string{"baseurl": adapter.GetBaseUrl()}
	urlTemplate := mcpPoliciesURLTemplate
	if policyID != "" {
		urlTemplate = mcpPolicyItemURLTemplate
		pathParameters["mcpPolicyId"] = policyID
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
	if requestBody != nil {
		if err := requestInfo.SetContentFromParsable(
			ctx,
			adapter,
			"application/json",
			requestBody,
		); err != nil {
			return nil, fmt.Errorf("set MCP policy request content: %w", err)
		}
	}
	return requestInfo, nil
}
