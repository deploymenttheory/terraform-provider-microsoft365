package graphBetaNetworkFilteringProfile

import (
	"context"
	"fmt"

	abstractions "github.com/microsoft/kiota-abstractions-go"
	s "github.com/microsoft/kiota-abstractions-go/serialization"
	models "github.com/microsoftgraph/msgraph-beta-sdk-go/models/networkaccess"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models/odataerrors"
)

const (
	filteringProfilesURLTemplate    = "{+baseurl}/networkAccess/filteringProfiles"
	filteringProfileItemURLTemplate = filteringProfilesURLTemplate + "/{filteringProfileId}"
)

var filteringProfileErrorMapping = abstractions.ErrorMappings{
	"XXX": odataerrors.CreateODataErrorFromDiscriminatorValue,
}

func (r *NetworkFilteringProfileResource) createFilteringProfile(ctx context.Context, requestBody s.Parsable) (models.FilteringProfileable, error) {
	adapter := r.client.GetAdapter()
	requestInfo, err := newFilteringProfileRequestInformation(ctx, adapter, abstractions.POST, "", requestBody)
	if err != nil {
		return nil, err
	}

	result, err := adapter.Send(
		ctx,
		requestInfo,
		models.CreateFilteringProfileFromDiscriminatorValue,
		filteringProfileErrorMapping,
	)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, fmt.Errorf("create filtering profile returned nil response")
	}

	profile, ok := result.(models.FilteringProfileable)
	if !ok {
		return nil, fmt.Errorf("create filtering profile returned %T, expected FilteringProfileable", result)
	}

	return profile, nil
}

func (r *NetworkFilteringProfileResource) updateFilteringProfile(ctx context.Context, filteringProfileID string, requestBody s.Parsable) error {
	adapter := r.client.GetAdapter()
	requestInfo, err := newFilteringProfileRequestInformation(ctx, adapter, abstractions.PATCH, filteringProfileID, requestBody)
	if err != nil {
		return err
	}

	return adapter.SendNoContent(ctx, requestInfo, filteringProfileErrorMapping)
}

func newFilteringProfileRequestInformation(ctx context.Context, adapter abstractions.RequestAdapter, method abstractions.HttpMethod, filteringProfileID string, requestBody s.Parsable) (*abstractions.RequestInformation, error) {
	pathParameters := map[string]string{
		"baseurl": adapter.GetBaseUrl(),
	}

	urlTemplate := filteringProfilesURLTemplate
	if filteringProfileID != "" {
		urlTemplate = filteringProfileItemURLTemplate
		pathParameters["filteringProfileId"] = filteringProfileID
	}

	requestInfo := abstractions.NewRequestInformationWithMethodAndUrlTemplateAndPathParameters(method, urlTemplate, pathParameters)
	requestInfo.Headers.TryAdd("Accept", "application/json")

	if requestBody != nil {
		if err := requestInfo.SetContentFromParsable(ctx, adapter, "application/json", requestBody); err != nil {
			return nil, fmt.Errorf("set filtering profile request content: %w", err)
		}
	}

	return requestInfo, nil
}
