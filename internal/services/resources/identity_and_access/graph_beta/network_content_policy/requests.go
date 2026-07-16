package graphBetaNetworkContentPolicy

import (
	"context"
	"fmt"

	abstractions "github.com/microsoft/kiota-abstractions-go"
	s "github.com/microsoft/kiota-abstractions-go/serialization"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models/odataerrors"
)

const (
	contentPoliciesURLTemplate   = "{+baseurl}/networkaccess/filePolicies"
	contentPolicyItemURLTemplate = contentPoliciesURLTemplate + "/{contentPolicyId}"
)

var contentPolicyErrorMapping = abstractions.ErrorMappings{"XXX": odataerrors.CreateODataErrorFromDiscriminatorValue}

func (r *NetworkContentPolicyResource) createContentPolicy(ctx context.Context, requestBody s.Parsable) (*contentPolicyResponse, error) {
	result, err := r.sendContentPolicy(ctx, abstractions.POST, "", requestBody)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, fmt.Errorf("create content policy returned nil response")
	}
	return result, nil
}

func (r *NetworkContentPolicyResource) getContentPolicy(ctx context.Context, policyID string) (*contentPolicyResponse, error) {
	return r.sendContentPolicy(ctx, abstractions.GET, policyID, nil)
}

func (r *NetworkContentPolicyResource) sendContentPolicy(ctx context.Context, method abstractions.HttpMethod, policyID string, requestBody s.Parsable) (*contentPolicyResponse, error) {
	adapter := r.client.GetAdapter()
	requestInfo, err := newContentPolicyRequestInformation(ctx, adapter, method, policyID, requestBody)
	if err != nil {
		return nil, err
	}
	result, err := adapter.Send(ctx, requestInfo, createContentPolicyResponseFromDiscriminatorValue, contentPolicyErrorMapping)
	if err != nil || result == nil {
		return nil, err
	}
	policy, ok := result.(*contentPolicyResponse)
	if !ok {
		return nil, fmt.Errorf("content policy request returned %T, expected contentPolicyResponse", result)
	}
	return policy, nil
}

func (r *NetworkContentPolicyResource) updateContentPolicy(ctx context.Context, policyID string, requestBody s.Parsable) error {
	adapter := r.client.GetAdapter()
	requestInfo, err := newContentPolicyRequestInformation(ctx, adapter, abstractions.PATCH, policyID, requestBody)
	if err != nil {
		return err
	}
	result, err := adapter.Send(ctx, requestInfo, createContentPolicyResponseFromDiscriminatorValue, contentPolicyErrorMapping)
	if err != nil {
		return err
	}
	if result == nil {
		return fmt.Errorf("update content policy returned nil response")
	}
	return nil
}

func (r *NetworkContentPolicyResource) deleteContentPolicy(ctx context.Context, policyID string) error {
	adapter := r.client.GetAdapter()
	requestInfo, err := newContentPolicyRequestInformation(ctx, adapter, abstractions.DELETE, policyID, nil)
	if err != nil {
		return err
	}
	return adapter.SendNoContent(ctx, requestInfo, contentPolicyErrorMapping)
}

func newContentPolicyRequestInformation(ctx context.Context, adapter abstractions.RequestAdapter, method abstractions.HttpMethod, policyID string, requestBody s.Parsable) (*abstractions.RequestInformation, error) {
	pathParameters := map[string]string{"baseurl": adapter.GetBaseUrl()}
	urlTemplate := contentPoliciesURLTemplate
	if policyID != "" {
		urlTemplate = contentPolicyItemURLTemplate
		pathParameters["contentPolicyId"] = policyID
	}
	requestInfo := abstractions.NewRequestInformationWithMethodAndUrlTemplateAndPathParameters(method, urlTemplate, pathParameters)
	requestInfo.Headers.TryAdd("Accept", "application/json")
	if requestBody != nil {
		if err := requestInfo.SetContentFromParsable(ctx, adapter, "application/json", requestBody); err != nil {
			return nil, fmt.Errorf("set content policy request content: %w", err)
		}
	}
	return requestInfo, nil
}
