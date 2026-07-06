package graphBetaNetworkWebFilteringPolicy

import (
	"context"
	"fmt"

	abstractions "github.com/microsoft/kiota-abstractions-go"
	s "github.com/microsoft/kiota-abstractions-go/serialization"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models/odataerrors"
)

const (
	// This endpoint is used by the Microsoft Entra Global Secure Access web
	// content filtering blade. It is not currently present in Microsoft Graph beta
	// metadata or the generated Go SDK, so this resource uses a custom Kiota
	// RequestInformation path instead of a generated request builder. Do not
	// confuse it with the documented /networkAccess/filteringPolicies surface:
	// https://learn.microsoft.com/graph/api/resources/networkaccess-filteringpolicy
	webFilteringPoliciesURLTemplate   = "{+baseurl}/networkaccess/webFilteringPolicies"
	webFilteringPolicyItemURLTemplate = webFilteringPoliciesURLTemplate + "/{webFilteringPolicyId}"
)

var webFilteringPolicyErrorMapping = abstractions.ErrorMappings{
	"XXX": odataerrors.CreateODataErrorFromDiscriminatorValue,
}

func (r *NetworkWebFilteringPolicyResource) createWebFilteringPolicy(ctx context.Context, requestBody s.Parsable) (*webFilteringPolicyResponse, error) {
	adapter := r.client.GetAdapter()
	requestInfo, err := newWebFilteringPolicyRequestInformation(ctx, adapter, abstractions.POST, "", requestBody)
	if err != nil {
		return nil, err
	}

	result, err := adapter.Send(ctx, requestInfo, createWebFilteringPolicyResponseFromDiscriminatorValue, webFilteringPolicyErrorMapping)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, fmt.Errorf("create web filtering policy returned nil response")
	}

	policy, ok := result.(*webFilteringPolicyResponse)
	if !ok {
		return nil, fmt.Errorf("create web filtering policy returned %T, expected webFilteringPolicyResponse", result)
	}

	return policy, nil
}

func (r *NetworkWebFilteringPolicyResource) getWebFilteringPolicy(ctx context.Context, policyID string) (*webFilteringPolicyResponse, error) {
	adapter := r.client.GetAdapter()
	requestInfo, err := newWebFilteringPolicyRequestInformation(ctx, adapter, abstractions.GET, policyID, nil)
	if err != nil {
		return nil, err
	}

	result, err := adapter.Send(ctx, requestInfo, createWebFilteringPolicyResponseFromDiscriminatorValue, webFilteringPolicyErrorMapping)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, nil
	}

	policy, ok := result.(*webFilteringPolicyResponse)
	if !ok {
		return nil, fmt.Errorf("get web filtering policy returned %T, expected webFilteringPolicyResponse", result)
	}

	return policy, nil
}

func (r *NetworkWebFilteringPolicyResource) updateWebFilteringPolicy(ctx context.Context, policyID string, requestBody s.Parsable) error {
	adapter := r.client.GetAdapter()
	requestInfo, err := newWebFilteringPolicyRequestInformation(ctx, adapter, abstractions.PATCH, policyID, requestBody)
	if err != nil {
		return err
	}

	return adapter.SendNoContent(ctx, requestInfo, webFilteringPolicyErrorMapping)
}

func (r *NetworkWebFilteringPolicyResource) deleteWebFilteringPolicy(ctx context.Context, policyID string) error {
	adapter := r.client.GetAdapter()
	requestInfo, err := newWebFilteringPolicyRequestInformation(ctx, adapter, abstractions.DELETE, policyID, nil)
	if err != nil {
		return err
	}

	return adapter.SendNoContent(ctx, requestInfo, webFilteringPolicyErrorMapping)
}

func newWebFilteringPolicyRequestInformation(ctx context.Context, adapter abstractions.RequestAdapter, method abstractions.HttpMethod, policyID string, requestBody s.Parsable) (*abstractions.RequestInformation, error) {
	pathParameters := map[string]string{
		"baseurl": adapter.GetBaseUrl(),
	}

	urlTemplate := webFilteringPoliciesURLTemplate
	if policyID != "" {
		urlTemplate = webFilteringPolicyItemURLTemplate
		pathParameters["webFilteringPolicyId"] = policyID
	}

	requestInfo := abstractions.NewRequestInformationWithMethodAndUrlTemplateAndPathParameters(method, urlTemplate, pathParameters)
	requestInfo.Headers.TryAdd("Accept", "application/json")

	if requestBody != nil {
		if err := requestInfo.SetContentFromParsable(ctx, adapter, "application/json", requestBody); err != nil {
			return nil, fmt.Errorf("set web filtering policy request content: %w", err)
		}
	}

	return requestInfo, nil
}
