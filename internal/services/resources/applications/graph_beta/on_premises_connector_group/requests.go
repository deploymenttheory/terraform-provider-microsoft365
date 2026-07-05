package graphBetaApplicationsOnPremisesConnectorGroup

import (
	"context"
	"fmt"

	abstractions "github.com/microsoft/kiota-abstractions-go"
	s "github.com/microsoft/kiota-abstractions-go/serialization"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models/odataerrors"
)

const (
	graphBetaBaseURL = "https://graph.microsoft.com/beta"

	// Microsoft Learn documents connectorGroup CRUD at:
	// https://learn.microsoft.com/en-us/graph/api/resources/connectorgroup?view=graph-rest-beta
	// https://learn.microsoft.com/en-us/graph/api/connectorgroup-post?view=graph-rest-beta
	//
	// The current msgraph-beta-sdk-go package exposes generated request builders
	// for this path. A dedicated Kiota request helper is still used because the
	// generated ConnectorGroup model parses region through an enum generated from
	// Microsoft Graph beta OData CSDL metadata, while direct API verification on
	// 2026-07-05 returned "region": "japan", which is not in that enum. This
	// same verification also showed create returning HTTP 200 instead of the
	// Learn-documented 201, read-only create body values such as isDefault being
	// ignored, and PATCH returning 204 with no response body. This helper keeps
	// the normal Kiota adapter pipeline but pairs it with the custom
	// request/response types in this package so actual API values are preserved.
	connectorGroupsURLTemplate = "{+baseurl}/onPremisesPublishingProfiles/applicationProxy/connectorGroups"
	connectorGroupURLTemplate  = connectorGroupsURLTemplate + "/{connectorGroupId}"
)

var connectorGroupErrorMapping = abstractions.ErrorMappings{
	"XXX": odataerrors.CreateODataErrorFromDiscriminatorValue,
}

func (r *OnPremisesConnectorGroupResource) createConnectorGroup(ctx context.Context, requestBody s.Parsable) (*connectorGroupResponse, error) {
	adapter := r.client.GetAdapter()
	requestInfo, err := newConnectorGroupRequestInformation(ctx, adapter, abstractions.POST, "", requestBody)
	if err != nil {
		return nil, err
	}

	result, err := adapter.Send(
		ctx,
		requestInfo,
		createConnectorGroupResponseFromDiscriminatorValue,
		connectorGroupErrorMapping,
	)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, fmt.Errorf("create connector group returned nil response")
	}

	connectorGroup, ok := result.(*connectorGroupResponse)
	if !ok {
		return nil, fmt.Errorf("create connector group returned %T, expected connectorGroupResponse", result)
	}

	return connectorGroup, nil
}

func (r *OnPremisesConnectorGroupResource) getConnectorGroup(ctx context.Context, connectorGroupID string) (*connectorGroupResponse, error) {
	adapter := r.client.GetAdapter()
	requestInfo, err := newConnectorGroupRequestInformation(ctx, adapter, abstractions.GET, connectorGroupID, nil)
	if err != nil {
		return nil, err
	}

	result, err := adapter.Send(
		ctx,
		requestInfo,
		createConnectorGroupResponseFromDiscriminatorValue,
		connectorGroupErrorMapping,
	)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, nil
	}

	connectorGroup, ok := result.(*connectorGroupResponse)
	if !ok {
		return nil, fmt.Errorf("get connector group returned %T, expected connectorGroupResponse", result)
	}

	return connectorGroup, nil
}

func (r *OnPremisesConnectorGroupResource) updateConnectorGroup(ctx context.Context, connectorGroupID string, requestBody s.Parsable) (*connectorGroupResponse, error) {
	adapter := r.client.GetAdapter()
	requestInfo, err := newConnectorGroupRequestInformation(ctx, adapter, abstractions.PATCH, connectorGroupID, requestBody)
	if err != nil {
		return nil, err
	}

	result, err := adapter.Send(
		ctx,
		requestInfo,
		createConnectorGroupResponseFromDiscriminatorValue,
		connectorGroupErrorMapping,
	)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, nil
	}

	connectorGroup, ok := result.(*connectorGroupResponse)
	if !ok {
		return nil, fmt.Errorf("update connector group returned %T, expected connectorGroupResponse", result)
	}

	return connectorGroup, nil
}

func (r *OnPremisesConnectorGroupResource) deleteConnectorGroup(ctx context.Context, connectorGroupID string) error {
	adapter := r.client.GetAdapter()
	requestInfo, err := newConnectorGroupRequestInformation(ctx, adapter, abstractions.DELETE, connectorGroupID, nil)
	if err != nil {
		return err
	}

	return adapter.SendNoContent(ctx, requestInfo, connectorGroupErrorMapping)
}

func newConnectorGroupRequestInformation(ctx context.Context, adapter abstractions.RequestAdapter, method abstractions.HttpMethod, connectorGroupID string, requestBody s.Parsable) (*abstractions.RequestInformation, error) {
	pathParameters := map[string]string{
		"baseurl": graphBetaBaseURL,
	}

	urlTemplate := connectorGroupsURLTemplate
	if connectorGroupID != "" {
		urlTemplate = connectorGroupURLTemplate
		pathParameters["connectorGroupId"] = connectorGroupID
	}

	requestInfo := abstractions.NewRequestInformationWithMethodAndUrlTemplateAndPathParameters(method, urlTemplate, pathParameters)
	requestInfo.Headers.TryAdd("Accept", "application/json")

	if requestBody != nil {
		if err := requestInfo.SetContentFromParsable(ctx, adapter, "application/json", requestBody); err != nil {
			return nil, fmt.Errorf("set connector group request content: %w", err)
		}
	}

	return requestInfo, nil
}
