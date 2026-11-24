package customrequests

import (
	"context"
	"encoding/json"

	abstractions "github.com/microsoft/kiota-abstractions-go"
	s "github.com/microsoft/kiota-abstractions-go/serialization"
)

// GetRequestConfig contains the configuration for a custom GET request
type GetRequestConfig struct {
	// The API version to use (beta or v1.0)
	APIVersion GraphAPIVersion
	// The base endpoint (e.g., "deviceManagement/configurationPolicies")
	Endpoint string
	// The endpoint suffix appended after the ID (e.g., "/settings"). Optional.
	EndpointSuffix string
	// The resource ID syntax format (e.g., "('id')" or "(id)")
	ResourceIDPattern string
	// The ID of the resource
	ResourceID string
	// Optional query parameters to include in the request (e.g. {"$expand": "children"})
	QueryParameters map[string]string
	// Optional request headers (e.g. {"Prefer": "include-unknown-enum-members"})
	Headers map[string]string
	// Optional raw JSON response target: a pointer the response is decoded into with encoding/json. Every
	// @odata.nextLink page is collected first. Use instead of a ParsableFactory.
	RawResponse any
	// Optional per-request middleware options, e.g. a kiota RetryHandlerOptions that narrows which responses the
	// retry middleware retries for this request only.
	RequestOptions []abstractions.RequestOption
}

// ODataResponse represents the structure of an OData response
type ODataResponse struct {
	// Value is the array of JSON messages returned by the request
	Value []json.RawMessage `json:"value"`
	//  NextLink is the URL for the next page of results used by pagination
	NextLink string `json:"@odata.nextLink,omitempty"`
}

// GetRequest performs a custom GET of "<endpoint>[/<resourceID>][/<endpointSuffix>]".
//
// In Parsable mode the response is parsed with factory and returned. In raw JSON mode (config.RawResponse set and
// factory nil) every @odata.nextLink page is collected, the result is decoded into config.RawResponse and the
// returned Parsable is nil. A response with no body leaves config.RawResponse untouched.
func GetRequest(
	ctx context.Context,
	adapter abstractions.RequestAdapter,
	config GetRequestConfig,
	factory s.ParsableFactory,
	errorMappings abstractions.ErrorMappings,
) (s.Parsable, error) {
	if factory != nil && config.RawResponse != nil {
		return nil, ErrMixedResponseModes
	}

	urlTemplate := endpointUrlTemplate(config.Endpoint)
	if config.GetResourceID() != "" {
		urlTemplate += "/" + config.GetResourceID()
	}
	if config.GetEndpointSuffix() != "" {
		urlTemplate += "/" + config.GetEndpointSuffix()
	}

	requestInfo := newGraphRequestInformation(adapter, abstractions.GET, urlTemplate, config.APIVersion, config.QueryParameters, config.Headers)
	requestInfo.AddRequestOptions(config.RequestOptions)

	if config.RawResponse != nil {
		requestInfo.Headers.TryAdd("Accept", jsonContentType)
		body, err := sendGraphGetCollectingPages(ctx, adapter, requestInfo, config.Headers, config.RequestOptions, errorMappings)
		if err != nil {
			return nil, err
		}
		return nil, decodeRawResponse(body, config.RawResponse)
	}

	result, err := adapter.Send(ctx, requestInfo, factory, withODataErrorMappings(errorMappings))
	if err != nil {
		return nil, err
	}

	return result, nil
}

// GetRequestByResourceId performs a custom GET request using the Microsoft Graph SDK when the operation
// is not available in the generated SDK methods or when using raw json is easier to handle for response handling.
// This function supports both Beta and V1.0 Graph API versions and automatically handles OData pagination if present.
//
// e.g., GET https://graph.microsoft.com/beta/deviceManagement/configurationPolicies('191056b1-4c4a-4871-8518-162a105d011a')/settings
//
// The function handles:
// - Construction of the Graph API URL with proper formatting
// - Setting up the GET request with optional query parameters
// - Sending the request with proper authentication
// - Automatic pagination if the response is an OData response with a nextLink
// - Combining paginated results into a single response
// - Returning the raw JSON response
//
// Parameters:
//   - ctx: The context for the request, which can be used for cancellation and timeout
//   - adapter: The request adapter for sending the request
//   - config: GetRequestConfig containing:
//   - APIVersion: The Graph API version to use (Beta or V1.0)
//   - Endpoint: The resource endpoint path (e.g., "/deviceManagement/configurationPolicies")
//   - ResourceID: The ID of the resource to retrieve
//   - ResourceIDPattern: The format for the resource ID (e.g., "('id')" or "(id)")
//   - EndpointSuffix: Optional suffix to append after the resource ID (e.g., "/settings")
//   - QueryParameters: Optional query parameters for the request
//   - Headers: Optional request headers
//   - RawResponse: Optional pointer the (combined) response is also decoded into
//   - RequestOptions: Optional per-request middleware options, sent with every page
//
// Returns:
//   - json.RawMessage: The raw JSON response from the GET request. For paginated responses,
//     returns a combined response with all results in the "value" array
//   - error: Returns nil if the request was successful, otherwise an error describing what went wrong
//
// Example Usage:
//
//	config := GetRequestConfig{
//		APIVersion:        GraphAPIBeta,
//		Endpoint:         "/deviceManagement/configurationPolicies",
//		ResourceID:       "d557c813-b8e5-4efc-b00e-9c0bd5fd10df",
//		ResourceIDPattern: "('id')",
//		EndpointSuffix:   "/settings",
//		QueryParameters: map[string]string{
//			"$expand": "children",
//		},
//	}
//
//	response, err := GetRequestByResourceId(ctx, adapter, config)
//	if err != nil {
//		log.Fatalf("Error: %v", err)
//	}
//
//	fmt.Printf("Response: %+v\n", response)
func GetRequestByResourceId(ctx context.Context, adapter abstractions.RequestAdapter, reqConfig GetRequestConfig) (json.RawMessage, error) {
	requestInfo := newGraphRequestInformation(adapter, abstractions.GET, ByIDRequestUrlTemplate(reqConfig), reqConfig.APIVersion, reqConfig.QueryParameters, reqConfig.Headers)
	requestInfo.AddRequestOptions(reqConfig.RequestOptions)
	requestInfo.Headers.TryAdd("Accept", jsonContentType)

	body, err := sendGraphGetCollectingPages(ctx, adapter, requestInfo, reqConfig.Headers, reqConfig.RequestOptions, nil)
	if err != nil {
		return nil, err
	}

	if reqConfig.RawResponse != nil {
		if err := decodeRawResponse(body, reqConfig.RawResponse); err != nil {
			return nil, err
		}
	}

	return body, nil
}
