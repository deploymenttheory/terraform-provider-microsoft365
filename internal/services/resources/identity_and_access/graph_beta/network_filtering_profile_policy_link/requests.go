package graphBetaNetworkFilteringProfilePolicyLink

import (
	"context"
	"fmt"

	abstractions "github.com/microsoft/kiota-abstractions-go"
	s "github.com/microsoft/kiota-abstractions-go/serialization"
	models "github.com/microsoftgraph/msgraph-beta-sdk-go/models/networkaccess"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models/odataerrors"
)

const (
	// Microsoft Learn documents two different create shapes on the same page:
	// https://learn.microsoft.com/en-us/graph/api/networkaccess-filteringpolicylink-post?view=graph-rest-beta
	// The Entra admin center XHR uses POST /networkAccess/filteringProfiles/{id}/policies with an
	// explicit policyLink body that references an existing policy. Keep the URL and request bodies
	// explicit so V2 webFilteringPolicyLink can be sent even though the current beta SDK discriminator
	// does not include that link type.
	filteringProfilePoliciesURLTemplate   = "{+baseurl}/networkAccess/filteringProfiles/{filteringProfileId}/policies"
	filteringProfilePolicyLinkURLTemplate = filteringProfilePoliciesURLTemplate + "/{policyLinkId}"
)

var policyLinkErrorMapping = abstractions.ErrorMappings{
	"XXX": odataerrors.CreateODataErrorFromDiscriminatorValue,
}

func (r *NetworkFilteringProfilePolicyLinkResource) createPolicyLink(ctx context.Context, filteringProfileID string, requestBody s.Parsable) (models.PolicyLinkable, error) {
	adapter := r.client.GetAdapter()
	requestInfo, err := newPolicyLinkRequestInformation(ctx, adapter, abstractions.POST, filteringProfileID, "", requestBody)
	if err != nil {
		return nil, err
	}

	result, err := adapter.Send(ctx, requestInfo, models.CreatePolicyLinkFromDiscriminatorValue, policyLinkErrorMapping)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, fmt.Errorf("create policy link returned nil response")
	}

	link, ok := result.(models.PolicyLinkable)
	if !ok {
		return nil, fmt.Errorf("create policy link returned %T, expected PolicyLinkable", result)
	}

	return link, nil
}

func (r *NetworkFilteringProfilePolicyLinkResource) updatePolicyLink(ctx context.Context, filteringProfileID, policyLinkID string, requestBody s.Parsable) error {
	adapter := r.client.GetAdapter()
	requestInfo, err := newPolicyLinkRequestInformation(ctx, adapter, abstractions.PATCH, filteringProfileID, policyLinkID, requestBody)
	if err != nil {
		return err
	}

	return adapter.SendNoContent(ctx, requestInfo, policyLinkErrorMapping)
}

func newPolicyLinkRequestInformation(ctx context.Context, adapter abstractions.RequestAdapter, method abstractions.HttpMethod, filteringProfileID, policyLinkID string, requestBody s.Parsable) (*abstractions.RequestInformation, error) {
	pathParameters := map[string]string{
		"baseurl":            adapter.GetBaseUrl(),
		"filteringProfileId": filteringProfileID,
	}

	urlTemplate := filteringProfilePoliciesURLTemplate
	if policyLinkID != "" {
		urlTemplate = filteringProfilePolicyLinkURLTemplate
		pathParameters["policyLinkId"] = policyLinkID
	}

	requestInfo := abstractions.NewRequestInformationWithMethodAndUrlTemplateAndPathParameters(method, urlTemplate, pathParameters)
	requestInfo.Headers.TryAdd("Accept", "application/json")

	if requestBody != nil {
		if err := requestInfo.SetContentFromParsable(ctx, adapter, "application/json", requestBody); err != nil {
			return nil, fmt.Errorf("set policy link request content: %w", err)
		}
	}

	return requestInfo, nil
}
