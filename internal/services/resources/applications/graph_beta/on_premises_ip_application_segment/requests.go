package graphBetaApplicationsOnPremisesIpApplicationSegment

import (
	"context"
	"fmt"

	abstractions "github.com/microsoft/kiota-abstractions-go"
	s "github.com/microsoft/kiota-abstractions-go/serialization"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models/odataerrors"
)

const (
	graphBetaBaseURL = "https://graph.microsoft.com/beta"

	// Microsoft Learn documents the application-scoped endpoint here:
	// https://learn.microsoft.com/en-us/graph/api/onpremisespublishingprofile-post-applicationsegments?view=graph-rest-beta
	//
	// POST /applications/{applicationObjectId}/onPremisesPublishing/segmentsConfiguration/microsoft.graph.ipSegmentConfiguration/applicationSegments
	//
	// The current msgraph-beta-sdk-go package does not expose a generated
	// request builder for this application-scoped segmentsConfiguration path.
	// A dedicated request helper is therefore used instead of the generic
	// resource path helpers or generated application builders. It still builds a
	// Kiota RequestInformation and sends it through the provider's Graph adapter,
	// so authentication, retry, middleware, serialization, and OData error
	// handling remain on the normal Kiota pipeline.
	ipApplicationSegmentsURLTemplate = "{+baseurl}/applications/{applicationObjectId}/onPremisesPublishing/segmentsConfiguration/microsoft.graph.ipSegmentConfiguration/applicationSegments"
	ipApplicationSegmentURLTemplate  = ipApplicationSegmentsURLTemplate + "/{ipApplicationSegmentId}"
)

var ipApplicationSegmentErrorMapping = abstractions.ErrorMappings{
	"XXX": odataerrors.CreateODataErrorFromDiscriminatorValue,
}

func (r *OnPremisesIpApplicationSegmentResource) createIpApplicationSegment(ctx context.Context, applicationObjectID string, requestBody s.Parsable) (*ipApplicationSegmentResponse, error) {
	adapter := r.client.GetAdapter()
	requestInfo, err := newIpApplicationSegmentRequestInformation(ctx, adapter, abstractions.POST, applicationObjectID, "", requestBody)
	if err != nil {
		return nil, err
	}

	result, err := adapter.Send(
		ctx,
		requestInfo,
		createIpApplicationSegmentResponseFromDiscriminatorValue,
		ipApplicationSegmentErrorMapping,
	)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, fmt.Errorf("create ip application segment returned nil response")
	}

	segment, ok := result.(*ipApplicationSegmentResponse)
	if !ok {
		return nil, fmt.Errorf("create ip application segment returned %T, expected ipApplicationSegmentResponse", result)
	}

	return segment, nil
}

func (r *OnPremisesIpApplicationSegmentResource) getIpApplicationSegment(ctx context.Context, applicationObjectID, segmentID string) (*ipApplicationSegmentResponse, error) {
	adapter := r.client.GetAdapter()
	requestInfo, err := newIpApplicationSegmentRequestInformation(ctx, adapter, abstractions.GET, applicationObjectID, segmentID, nil)
	if err != nil {
		return nil, err
	}

	result, err := adapter.Send(
		ctx,
		requestInfo,
		createIpApplicationSegmentResponseFromDiscriminatorValue,
		ipApplicationSegmentErrorMapping,
	)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, nil
	}

	segment, ok := result.(*ipApplicationSegmentResponse)
	if !ok {
		return nil, fmt.Errorf("get ip application segment returned %T, expected ipApplicationSegmentResponse", result)
	}

	return segment, nil
}

func (r *OnPremisesIpApplicationSegmentResource) updateIpApplicationSegment(ctx context.Context, applicationObjectID, segmentID string, requestBody s.Parsable) error {
	adapter := r.client.GetAdapter()
	requestInfo, err := newIpApplicationSegmentRequestInformation(ctx, adapter, abstractions.PATCH, applicationObjectID, segmentID, requestBody)
	if err != nil {
		return err
	}

	return adapter.SendNoContent(ctx, requestInfo, ipApplicationSegmentErrorMapping)
}

func (r *OnPremisesIpApplicationSegmentResource) deleteIpApplicationSegment(ctx context.Context, applicationObjectID, segmentID string) error {
	adapter := r.client.GetAdapter()
	requestInfo, err := newIpApplicationSegmentRequestInformation(ctx, adapter, abstractions.DELETE, applicationObjectID, segmentID, nil)
	if err != nil {
		return err
	}

	return adapter.SendNoContent(ctx, requestInfo, ipApplicationSegmentErrorMapping)
}

func newIpApplicationSegmentRequestInformation(ctx context.Context, adapter abstractions.RequestAdapter, method abstractions.HttpMethod, applicationObjectID, segmentID string, requestBody s.Parsable) (*abstractions.RequestInformation, error) {
	pathParameters := map[string]string{
		"baseurl":             graphBetaBaseURL,
		"applicationObjectId": applicationObjectID,
	}

	urlTemplate := ipApplicationSegmentsURLTemplate
	if segmentID != "" {
		urlTemplate = ipApplicationSegmentURLTemplate
		pathParameters["ipApplicationSegmentId"] = segmentID
	}

	requestInfo := abstractions.NewRequestInformationWithMethodAndUrlTemplateAndPathParameters(method, urlTemplate, pathParameters)
	requestInfo.Headers.TryAdd("Accept", "application/json")

	if requestBody != nil {
		if err := requestInfo.SetContentFromParsable(ctx, adapter, "application/json", requestBody); err != nil {
			return nil, fmt.Errorf("set ip application segment request content: %w", err)
		}
	}

	return requestInfo, nil
}
