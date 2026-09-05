//nolint:wrapcheck // Preserve Kiota API error types for the shared Graph HTTP status and permission handler.
package graphBetaNetworkPromptPolicy

import (
	"context"
	"fmt"

	abstractions "github.com/microsoft/kiota-abstractions-go"
	s "github.com/microsoft/kiota-abstractions-go/serialization"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models/odataerrors"
)

const (
	promptPoliciesURLTemplate   = "{+baseurl}/networkaccess/promptPolicies"
	promptPolicyItemURLTemplate = promptPoliciesURLTemplate + "/{promptPolicyId}"
)

var promptPolicyErrorMapping = abstractions.ErrorMappings{
	"XXX": odataerrors.CreateODataErrorFromDiscriminatorValue,
}

func (r *NetworkPromptPolicyResource) createPromptPolicy(
	ctx context.Context,
	requestBody s.Parsable,
) (*promptPolicyResponse, error) {
	result, err := r.sendPromptPolicy(ctx, abstractions.POST, "", requestBody)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, errEmptyResponse
	}
	return result, nil
}

func (r *NetworkPromptPolicyResource) getPromptPolicy(
	ctx context.Context,
	policyID string,
) (*promptPolicyResponse, error) {
	return r.sendPromptPolicy(ctx, abstractions.GET, policyID, nil)
}

func (r *NetworkPromptPolicyResource) sendPromptPolicy(
	ctx context.Context,
	method abstractions.HttpMethod,
	policyID string,
	requestBody s.Parsable,
) (*promptPolicyResponse, error) {
	adapter := r.client.GetAdapter()
	requestInfo, err := newPromptPolicyRequestInformation(
		ctx,
		adapter,
		method,
		policyID,
		requestBody,
	)
	if err != nil {
		return nil, err
	}
	result, err := adapter.Send(
		ctx,
		requestInfo,
		createPromptPolicyResponseFromDiscriminatorValue,
		promptPolicyErrorMapping,
	)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, errEmptyResponse
	}
	policy, ok := result.(*promptPolicyResponse)
	if !ok {
		return nil, fmt.Errorf("%w: received %T", errInvalidResponse, result)
	}
	return policy, nil
}

func (r *NetworkPromptPolicyResource) updatePromptPolicy(
	ctx context.Context,
	policyID string,
	requestBody s.Parsable,
) error {
	adapter := r.client.GetAdapter()
	requestInfo, err := newPromptPolicyRequestInformation(
		ctx,
		adapter,
		abstractions.PATCH,
		policyID,
		requestBody,
	)
	if err != nil {
		return err
	}
	_, err = adapter.Send(
		ctx,
		requestInfo,
		createPromptPolicyResponseFromDiscriminatorValue,
		promptPolicyErrorMapping,
	)
	return err
}

func (r *NetworkPromptPolicyResource) deletePromptPolicy(
	ctx context.Context,
	policyID string,
) error {
	adapter := r.client.GetAdapter()
	requestInfo, err := newPromptPolicyRequestInformation(
		ctx,
		adapter,
		abstractions.DELETE,
		policyID,
		nil,
	)
	if err != nil {
		return err
	}
	return adapter.SendNoContent(ctx, requestInfo, promptPolicyErrorMapping)
}

func newPromptPolicyRequestInformation(
	ctx context.Context,
	adapter abstractions.RequestAdapter,
	method abstractions.HttpMethod,
	policyID string,
	requestBody s.Parsable,
) (*abstractions.RequestInformation, error) {
	pathParameters := map[string]string{"baseurl": adapter.GetBaseUrl()}
	urlTemplate := promptPoliciesURLTemplate
	if policyID != "" {
		urlTemplate = promptPolicyItemURLTemplate
		pathParameters["promptPolicyId"] = policyID
	}
	requestInfo := abstractions.NewRequestInformationWithMethodAndUrlTemplateAndPathParameters(
		method,
		urlTemplate,
		pathParameters,
	)
	requestInfo.Headers.TryAdd("Accept", "application/json")
	if requestBody != nil {
		if err := requestInfo.SetContentFromParsable(
			ctx,
			adapter,
			"application/json",
			requestBody,
		); err != nil {
			return nil, fmt.Errorf("set prompt policy request content: %w", err)
		}
	}
	return requestInfo, nil
}
