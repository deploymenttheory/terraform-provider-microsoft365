package customrequests

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/hashicorp/terraform-plugin-log/tflog"
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
	// Optional query parameters to include in the request
	QueryParameters map[string]string
}

// ODataResponse represents the structure of an OData response
type ODataResponse struct {
	// Value is the array of JSON messages returned by the request
	Value []json.RawMessage `json:"value"`
	//  NextLink is the URL for the next page of results used by pagination
	NextLink string `json:"@odata.nextLink,omitempty"`
}

func GetRequest(
	ctx context.Context,
	adapter abstractions.RequestAdapter,
	config GetRequestConfig,
	factory s.ParsableFactory,
	errorMappings abstractions.ErrorMappings,
) (s.Parsable, error) {
	urlTemplate := endpointUrlTemplate(config.Endpoint)
	if config.GetResourceID() != "" {
		urlTemplate += "/" + config.GetResourceID()
	}
	if config.GetEndpointSuffix() != "" {
		urlTemplate += "/" + config.GetEndpointSuffix()
	}

	requestInfo := newGraphRequestInformation(adapter, abstractions.GET, urlTemplate, config.APIVersion)

	if config.QueryParameters != nil {
		for key, value := range config.QueryParameters {
			requestInfo.QueryParameters[key] = value
		}
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

	requestInfo := newGraphRequestInformation(adapter, abstractions.GET, ByIDRequestUrlTemplate(reqConfig), reqConfig.APIVersion)
	requestInfo.Headers.Add("Accept", "application/json")

	if reqConfig.QueryParameters != nil {
		for key, value := range reqConfig.QueryParameters {
			requestInfo.QueryParameters[key] = value
		}
	}

	body, err := sendGraphRequestForBytes(ctx, adapter, requestInfo)
	if err != nil {
		return nil, err
	}

	var firstResponse ODataResponse
	if err := json.Unmarshal(body, &firstResponse); err != nil {
		return body, nil
	}

	// If no NextLink or no Value array, this isn't a paginated response
	if firstResponse.NextLink == "" || firstResponse.Value == nil {
		return body, nil
	}

	var allResults []json.RawMessage
	allResults = append(allResults, firstResponse.Value...)
	nextLink := firstResponse.NextLink

	tflog.Debug(ctx, "Pagination detected, retrieving additional pages", map[string]any{
		"itemsRetrieved": len(allResults),
	})

	for nextLink != "" {
		nextURL, err := url.Parse(nextLink)
		if err != nil {
			return nil, fmt.Errorf("error parsing @odata.nextLink %q: %w", nextLink, err)
		}

		requestInfo = abstractions.NewRequestInformation()
		requestInfo.Method = abstractions.GET
		requestInfo.SetUri(*nextURL)
		requestInfo.Headers.Add("Accept", "application/json")

		body, err = sendGraphRequestForBytes(ctx, adapter, requestInfo)
		if err != nil {
			return nil, err
		}

		var pageResponse ODataResponse
		if err := json.Unmarshal(body, &pageResponse); err != nil {
			return nil, fmt.Errorf("error parsing paginated response: %w", err)
		}

		allResults = append(allResults, pageResponse.Value...)
		nextLink = pageResponse.NextLink

		tflog.Debug(ctx, "Retrieved additional page", map[string]any{
			"itemsRetrieved": len(allResults),
			"hasNextPage":    nextLink != "",
		})
	}

	combinedResponse := map[string]any{
		"value": allResults,
	}

	return json.Marshal(combinedResponse)
}

// sendGraphRequestForBytes sends requestInfo through the adapter's pipeline (authentication, retry, compression,
// logging and status checks) and returns the raw response body. Any 4xx or 5xx response is returned as an ODataError.
// A response with no body (e.g. 204 No Content) returns nil.
func sendGraphRequestForBytes(ctx context.Context, adapter abstractions.RequestAdapter, requestInfo *abstractions.RequestInformation) ([]byte, error) {
	tflog.Debug(ctx, "Making custom msgraph GET request", map[string]any{"url_template": requestInfo.UrlTemplate})

	result, err := adapter.SendPrimitive(ctx, requestInfo, "[]byte", withODataErrorMappings(nil))
	if err != nil {
		return nil, err
	}

	body, _ := result.([]byte)
	return body, nil
}
