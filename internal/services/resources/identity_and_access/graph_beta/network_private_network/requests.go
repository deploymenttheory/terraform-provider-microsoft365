package graphBetaNetworkPrivateNetwork

import (
	"context"
	"fmt"

	abstractions "github.com/microsoft/kiota-abstractions-go"
	s "github.com/microsoft/kiota-abstractions-go/serialization"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models/odataerrors"
)

const (
	privateNetworksURLTemplate    = "{+baseurl}/networkaccess/privateNetworks"
	privateNetworkItemURLTemplate = privateNetworksURLTemplate + "/{privateNetworkId}"
)

var privateNetworkErrorMapping = abstractions.ErrorMappings{
	"XXX": odataerrors.CreateODataErrorFromDiscriminatorValue,
}

func (r *NetworkPrivateNetworkResource) createPrivateNetwork(ctx context.Context, requestBody s.Parsable) (*privateNetworkResponse, error) {
	adapter := r.client.GetAdapter()
	requestInfo, err := newPrivateNetworkRequestInformation(ctx, adapter, abstractions.POST, "", requestBody)
	if err != nil {
		return nil, err
	}

	result, err := adapter.Send(ctx, requestInfo, createPrivateNetworkResponseFromDiscriminatorValue, privateNetworkErrorMapping)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, fmt.Errorf("create private network returned nil response")
	}

	privateNetwork, ok := result.(*privateNetworkResponse)
	if !ok {
		return nil, fmt.Errorf("create private network returned %T, expected privateNetworkResponse", result)
	}

	return privateNetwork, nil
}

func (r *NetworkPrivateNetworkResource) getPrivateNetwork(ctx context.Context, privateNetworkID string) (*privateNetworkResponse, error) {
	adapter := r.client.GetAdapter()
	requestInfo, err := newPrivateNetworkRequestInformation(ctx, adapter, abstractions.GET, privateNetworkID, nil)
	if err != nil {
		return nil, err
	}

	result, err := adapter.Send(ctx, requestInfo, createPrivateNetworkResponseFromDiscriminatorValue, privateNetworkErrorMapping)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, nil
	}

	privateNetwork, ok := result.(*privateNetworkResponse)
	if !ok {
		return nil, fmt.Errorf("get private network returned %T, expected privateNetworkResponse", result)
	}

	return privateNetwork, nil
}

func (r *NetworkPrivateNetworkResource) updatePrivateNetwork(ctx context.Context, privateNetworkID string, requestBody s.Parsable) (*privateNetworkResponse, error) {
	adapter := r.client.GetAdapter()
	requestInfo, err := newPrivateNetworkRequestInformation(ctx, adapter, abstractions.PATCH, privateNetworkID, requestBody)
	if err != nil {
		return nil, err
	}

	result, err := adapter.Send(ctx, requestInfo, createPrivateNetworkResponseFromDiscriminatorValue, privateNetworkErrorMapping)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, nil
	}

	privateNetwork, ok := result.(*privateNetworkResponse)
	if !ok {
		return nil, fmt.Errorf("update private network returned %T, expected privateNetworkResponse", result)
	}

	return privateNetwork, nil
}

func (r *NetworkPrivateNetworkResource) deletePrivateNetwork(ctx context.Context, privateNetworkID string) error {
	adapter := r.client.GetAdapter()
	requestInfo, err := newPrivateNetworkRequestInformation(ctx, adapter, abstractions.DELETE, privateNetworkID, nil)
	if err != nil {
		return err
	}

	return adapter.SendNoContent(ctx, requestInfo, privateNetworkErrorMapping)
}

func newPrivateNetworkRequestInformation(ctx context.Context, adapter abstractions.RequestAdapter, method abstractions.HttpMethod, privateNetworkID string, requestBody s.Parsable) (*abstractions.RequestInformation, error) {
	pathParameters := map[string]string{
		"baseurl": adapter.GetBaseUrl(),
	}

	urlTemplate := privateNetworksURLTemplate
	if privateNetworkID != "" {
		urlTemplate = privateNetworkItemURLTemplate
		pathParameters["privateNetworkId"] = privateNetworkID
	}

	requestInfo := abstractions.NewRequestInformationWithMethodAndUrlTemplateAndPathParameters(method, urlTemplate, pathParameters)
	requestInfo.Headers.TryAdd("Accept", "application/json")

	if requestBody != nil {
		if err := requestInfo.SetContentFromParsable(ctx, adapter, "application/json", requestBody); err != nil {
			return nil, fmt.Errorf("set private network request content: %w", err)
		}
	}

	return requestInfo, nil
}
