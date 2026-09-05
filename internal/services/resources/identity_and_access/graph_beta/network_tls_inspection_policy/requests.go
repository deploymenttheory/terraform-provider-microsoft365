//nolint:wrapcheck // Preserve Kiota API error types for the shared Graph HTTP status and permission handler.
package graphBetaNetworkTLSInspectionPolicy

import (
	"context"
	"fmt"

	abstractions "github.com/microsoft/kiota-abstractions-go"
	s "github.com/microsoft/kiota-abstractions-go/serialization"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models/odataerrors"
)

const (
	tlsInspectionPoliciesURLTemplate   = "{+baseurl}/networkaccess/tlsInspectionPolicies"
	tlsInspectionPolicyItemURLTemplate = tlsInspectionPoliciesURLTemplate + "/{tlsInspectionPolicyId}"
)

var tlsInspectionPolicyErrorMapping = abstractions.ErrorMappings{
	"XXX": odataerrors.CreateODataErrorFromDiscriminatorValue,
}

func (r *NetworkTLSInspectionPolicyResource) createTLSInspectionPolicy(
	ctx context.Context,
	requestBody s.Parsable,
) (*tlsInspectionPolicyResponse, error) {
	result, err := r.sendTLSInspectionPolicy(ctx, abstractions.POST, "", requestBody)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, errEmptyResponse
	}
	return result, nil
}

func (r *NetworkTLSInspectionPolicyResource) getTLSInspectionPolicy(
	ctx context.Context,
	policyID string,
) (*tlsInspectionPolicyResponse, error) {
	return r.sendTLSInspectionPolicy(ctx, abstractions.GET, policyID, nil)
}

func (r *NetworkTLSInspectionPolicyResource) sendTLSInspectionPolicy(
	ctx context.Context,
	method abstractions.HttpMethod,
	policyID string,
	requestBody s.Parsable,
) (*tlsInspectionPolicyResponse, error) {
	adapter := r.client.GetAdapter()
	requestInfo, err := newTLSInspectionPolicyRequestInformation(
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
		createTLSInspectionPolicyResponseFromDiscriminatorValue,
		tlsInspectionPolicyErrorMapping,
	)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, errEmptyResponse
	}
	policy, ok := result.(*tlsInspectionPolicyResponse)
	if !ok {
		return nil, fmt.Errorf("%w: received %T", errInvalidResponse, result)
	}
	return policy, nil
}

func (r *NetworkTLSInspectionPolicyResource) updateTLSInspectionPolicy(
	ctx context.Context,
	policyID string,
	requestBody s.Parsable,
) error {
	adapter := r.client.GetAdapter()
	requestInfo, err := newTLSInspectionPolicyRequestInformation(
		ctx,
		adapter,
		abstractions.PATCH,
		policyID,
		requestBody,
	)
	if err != nil {
		return err
	}
	return adapter.SendNoContent(ctx, requestInfo, tlsInspectionPolicyErrorMapping)
}

func (r *NetworkTLSInspectionPolicyResource) deleteTLSInspectionPolicy(
	ctx context.Context,
	policyID string,
) error {
	adapter := r.client.GetAdapter()
	requestInfo, err := newTLSInspectionPolicyRequestInformation(
		ctx,
		adapter,
		abstractions.DELETE,
		policyID,
		nil,
	)
	if err != nil {
		return err
	}
	return adapter.SendNoContent(ctx, requestInfo, tlsInspectionPolicyErrorMapping)
}

func newTLSInspectionPolicyRequestInformation(
	ctx context.Context,
	adapter abstractions.RequestAdapter,
	method abstractions.HttpMethod,
	policyID string,
	requestBody s.Parsable,
) (*abstractions.RequestInformation, error) {
	pathParameters := map[string]string{"baseurl": adapter.GetBaseUrl()}
	urlTemplate := tlsInspectionPoliciesURLTemplate
	if policyID != "" {
		urlTemplate = tlsInspectionPolicyItemURLTemplate
		pathParameters["tlsInspectionPolicyId"] = policyID
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
			return nil, fmt.Errorf("set TLS inspection policy request content: %w", err)
		}
	}
	return requestInfo, nil
}
