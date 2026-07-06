package graphBetaNetworkWebContentFilteringPolicy

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
	webContentFilteringPoliciesURLTemplate   = "{+baseurl}/networkaccess/webFilteringPolicies"
	webContentFilteringPolicyItemURLTemplate = webContentFilteringPoliciesURLTemplate + "/{webContentFilteringPolicyId}"
)

var webContentFilteringPolicyErrorMapping = abstractions.ErrorMappings{
	"XXX": odataerrors.CreateODataErrorFromDiscriminatorValue,
}

func (r *NetworkWebContentFilteringPolicyResource) createWebContentFilteringPolicy(ctx context.Context, requestBody s.Parsable) (*webContentFilteringPolicyResponse, error) {
	adapter := r.client.GetAdapter()
	requestInfo, err := newWebContentFilteringPolicyRequestInformation(ctx, adapter, abstractions.POST, "", requestBody)
	if err != nil {
		return nil, err
	}

	result, err := adapter.Send(ctx, requestInfo, createWebContentFilteringPolicyResponseFromDiscriminatorValue, webContentFilteringPolicyErrorMapping)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, fmt.Errorf("create web content filtering policy returned nil response")
	}

	policy, ok := result.(*webContentFilteringPolicyResponse)
	if !ok {
		return nil, fmt.Errorf("create web content filtering policy returned %T, expected webContentFilteringPolicyResponse", result)
	}

	return policy, nil
}

func (r *NetworkWebContentFilteringPolicyResource) getWebContentFilteringPolicy(ctx context.Context, policyID string) (*webContentFilteringPolicyResponse, error) {
	adapter := r.client.GetAdapter()
	requestInfo, err := newWebContentFilteringPolicyRequestInformation(ctx, adapter, abstractions.GET, policyID, nil)
	if err != nil {
		return nil, err
	}

	result, err := adapter.Send(ctx, requestInfo, createWebContentFilteringPolicyResponseFromDiscriminatorValue, webContentFilteringPolicyErrorMapping)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, nil
	}

	policy, ok := result.(*webContentFilteringPolicyResponse)
	if !ok {
		return nil, fmt.Errorf("get web content filtering policy returned %T, expected webContentFilteringPolicyResponse", result)
	}

	return policy, nil
}

func (r *NetworkWebContentFilteringPolicyResource) updateWebContentFilteringPolicy(ctx context.Context, policyID string, requestBody s.Parsable) error {
	adapter := r.client.GetAdapter()
	requestInfo, err := newWebContentFilteringPolicyRequestInformation(ctx, adapter, abstractions.PATCH, policyID, requestBody)
	if err != nil {
		return err
	}

	return adapter.SendNoContent(ctx, requestInfo, webContentFilteringPolicyErrorMapping)
}

func (r *NetworkWebContentFilteringPolicyResource) deleteWebContentFilteringPolicy(ctx context.Context, policyID string) error {
	adapter := r.client.GetAdapter()
	requestInfo, err := newWebContentFilteringPolicyRequestInformation(ctx, adapter, abstractions.DELETE, policyID, nil)
	if err != nil {
		return err
	}

	return adapter.SendNoContent(ctx, requestInfo, webContentFilteringPolicyErrorMapping)
}

func newWebContentFilteringPolicyRequestInformation(ctx context.Context, adapter abstractions.RequestAdapter, method abstractions.HttpMethod, policyID string, requestBody s.Parsable) (*abstractions.RequestInformation, error) {
	pathParameters := map[string]string{
		"baseurl": adapter.GetBaseUrl(),
	}

	urlTemplate := webContentFilteringPoliciesURLTemplate
	if policyID != "" {
		urlTemplate = webContentFilteringPolicyItemURLTemplate
		pathParameters["webContentFilteringPolicyId"] = policyID
	}

	requestInfo := abstractions.NewRequestInformationWithMethodAndUrlTemplateAndPathParameters(method, urlTemplate, pathParameters)
	requestInfo.Headers.TryAdd("Accept", "application/json")

	if requestBody != nil {
		if err := requestInfo.SetContentFromParsable(ctx, adapter, "application/json", requestBody); err != nil {
			return nil, fmt.Errorf("set web content filtering policy request content: %w", err)
		}
	}

	return requestInfo, nil
}
