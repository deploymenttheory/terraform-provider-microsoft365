package customrequests

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	abstractions "github.com/microsoft/kiota-abstractions-go"
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

	requestInfo := abstractions.NewRequestInformation()
	requestInfo.Method = abstractions.GET
	requestInfo.UrlTemplate = ByIDRequestUrlTemplate(reqConfig)
	requestInfo.PathParameters = map[string]string{
		"baseurl": fmt.Sprintf("https://graph.microsoft.com/%s", reqConfig.APIVersion),
	}
	requestInfo.Headers.Add("Accept", "application/json")

	if reqConfig.QueryParameters != nil {
		for key, value := range reqConfig.QueryParameters {
			requestInfo.QueryParameters[key] = value
		}
	}

	body, err := makeRequest(ctx, adapter, requestInfo)
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

	tflog.Debug(ctx, "Pagination detected, retrieving additional pages", map[string]interface{}{
		"itemsRetrieved": len(allResults),
	})

	for nextLink != "" {
		requestInfo = abstractions.NewRequestInformation()
		requestInfo.Method = abstractions.GET
		requestInfo.UrlTemplate = nextLink
		requestInfo.Headers.Add("Accept", "application/json")

		body, err = makeRequest(ctx, adapter, requestInfo)
		if err != nil {
			return nil, err
		}

		var pageResponse ODataResponse
		if err := json.Unmarshal(body, &pageResponse); err != nil {
			return nil, fmt.Errorf("error parsing paginated response: %w", err)
		}

		allResults = append(allResults, pageResponse.Value...)
		nextLink = pageResponse.NextLink

		tflog.Debug(ctx, "Retrieved additional page", map[string]interface{}{
			"itemsRetrieved": len(allResults),
			"hasNextPage":    nextLink != "",
		})
	}

	combinedResponse := map[string]interface{}{
		"value": allResults,
	}

	return json.Marshal(combinedResponse)
}

// makeRequest executes an HTTP request using the provided Kiota request adapter and request information.
// This helper function handles the conversion of Kiota's RequestInformation into a native HTTP request,
// executes the request, and returns the raw response body.
//
// Parameters:
//   - ctx: The context for the request, which can be used for cancellation and timeout
//   - adapter: The Kiota request adapter that converts RequestInformation to a native request
//   - requestInfo: The Kiota RequestInformation containing the request configuration
//
// Returns:
//   - []byte: The raw response body from the HTTP request
//   - error: Returns nil if the request was successful, otherwise an error describing what went wrong
//
// The function performs the following steps:
// 1. Converts the Kiota RequestInformation to a native HTTP request
// 2. Executes the HTTP request using a standard http.Client
// 3. Reads and returns the complete response body
func makeRequest(ctx context.Context, adapter abstractions.RequestAdapter, requestInfo *abstractions.RequestInformation) ([]byte, error) {
	nativeReq, err := adapter.ConvertToNativeRequest(ctx, requestInfo)
	if err != nil {
		return nil, fmt.Errorf("error converting to native HTTP request: %w", err)
	}

	httpReq := nativeReq.(*http.Request)
	client := &http.Client{}

	tflog.Debug(ctx, "Making request", map[string]interface{}{
		"url": httpReq.URL.String(),
	})

	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
