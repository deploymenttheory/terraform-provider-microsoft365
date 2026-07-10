package graphBetaNetworkForwardingProfilePolicyLink

import (
	"context"
	"fmt"

	abstractions "github.com/microsoft/kiota-abstractions-go"
	s "github.com/microsoft/kiota-abstractions-go/serialization"
	models "github.com/microsoftgraph/msgraph-beta-sdk-go/models/networkaccess"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models/odataerrors"
)

const (
	forwardingProfilePolicyLinkURLTemplate = "{+baseurl}/networkAccess/forwardingProfiles/{forwardingProfileId}/policies/{policyLinkId}"
)

var forwardingPolicyLinkErrorMapping = abstractions.ErrorMappings{
	"XXX": odataerrors.CreateODataErrorFromDiscriminatorValue,
}

func (r *NetworkForwardingProfilePolicyLinkResource) getPolicyLink(ctx context.Context, forwardingProfileID, policyLinkID string) (models.PolicyLinkable, error) {
	result, err := r.client.NetworkAccess().ForwardingProfiles().ByForwardingProfileId(forwardingProfileID).Policies().ByPolicyLinkId(policyLinkID).Get(ctx, nil)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r *NetworkForwardingProfilePolicyLinkResource) patchPolicyLinkState(ctx context.Context, forwardingProfileID, policyLinkID string, requestBody s.Parsable) error {
	adapter := r.client.GetAdapter()
	requestInfo, err := newForwardingPolicyLinkRequestInformation(ctx, adapter, abstractions.PATCH, forwardingProfileID, policyLinkID, requestBody)
	if err != nil {
		return err
	}
	return adapter.SendNoContent(ctx, requestInfo, forwardingPolicyLinkErrorMapping)
}

func newForwardingPolicyLinkRequestInformation(ctx context.Context, adapter abstractions.RequestAdapter, method abstractions.HttpMethod, forwardingProfileID, policyLinkID string, requestBody s.Parsable) (*abstractions.RequestInformation, error) {
	pathParameters := map[string]string{
		"baseurl":             adapter.GetBaseUrl(),
		"forwardingProfileId": forwardingProfileID,
		"policyLinkId":        policyLinkID,
	}

	requestInfo := abstractions.NewRequestInformationWithMethodAndUrlTemplateAndPathParameters(method, forwardingProfilePolicyLinkURLTemplate, pathParameters)
	requestInfo.Headers.TryAdd("Accept", "application/json")

	if requestBody != nil {
		if err := requestInfo.SetContentFromParsable(ctx, adapter, "application/json", requestBody); err != nil {
			return nil, fmt.Errorf("set forwarding policy link request content: %w", err)
		}
	}
	return requestInfo, nil
}
