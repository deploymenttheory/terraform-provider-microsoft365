package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	abstractions "github.com/microsoft/kiota-abstractions-go"
)

// CustomGetRequestConfig contains the configuration for a custom GET request
type CustomGetRequestConfig struct {
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

// SendCustomGetRequestByResourceId performs a custom GET request using the Microsoft Graph SDK when the operation
// is not available in the generated SDK methods or when using raw json is easier to handle for response handling during stating operations.
// This function supports both Beta and V1.0 Graph API versions.
//
// e.g., GET https://graph.microsoft.com/beta/deviceManagement/configurationPolicies('191056b1-4c4a-4871-8518-162a105d011a')/settings
//
// The function handles:
// - Construction of the Graph API URL with proper formatting
// - Setting up the GET request with optional query parameters
// - Sending the request with proper authentication
// - Returning the raw JSON response
//
// Parameters:
//   - ctx: The context for the request, which can be used for cancellation and timeout
//   - adapter: The request adapter for sending the request
//   - config: CustomGetRequestConfig containing:
//   - APIVersion: The Graph API version to use (Beta or V1.0)
//   - Endpoint: The resource endpoint path (e.g., "deviceManagement/configurationPolicies")
//   - ResourceID: The ID of the resource to retrieve
//   - QueryParameters: Optional query parameters for the request
//
// Returns:
//   - json.RawMessage: The raw JSON response from the GET request
//   - error: Returns nil if the request was successful, otherwise an error describing what went wrong
//
// Example Usage:
//
//		config := CustomGetRequestConfig{
//		APIVersion: GraphAPIBeta,
//		Endpoint:   "deviceManagement/configurationPolicies('{id}')/settings",
//		ResourceID: "d557c813-b8e5-4efc-b00e-9c0bd5fd10df",
//		QueryParameters: map[string]string{
//				"$expand": "children",
//		},
//	}
//
// response, err := SendCustomGetRequestByResourceId(ctx, adapter, config, factory, errorMappings)
//
//	if err != nil {
//		log.Fatalf("Error: %v", err)
//	}
//
// fmt.Printf("Response: %+v\n", response)
func SendCustomGetRequestByResourceId(ctx context.Context, adapter abstractions.RequestAdapter, reqConfig CustomGetRequestConfig) (json.RawMessage, error) {
	requestInfo := abstractions.NewRequestInformation()
	requestInfo.Method = abstractions.GET

	// Build endpoint with ID syntax
	idFormat := strings.ReplaceAll(reqConfig.ResourceIDPattern, "id", reqConfig.ResourceID)
	endpoint := reqConfig.Endpoint + idFormat
	if reqConfig.EndpointSuffix != "" {
		endpoint += reqConfig.EndpointSuffix
	}

	requestInfo.UrlTemplate = "{+baseurl}/" + endpoint
	requestInfo.PathParameters = map[string]string{
		"baseurl": fmt.Sprintf("https://graph.microsoft.com/%s", reqConfig.APIVersion),
	}
	requestInfo.Headers.Add("Accept", "application/json")

	if reqConfig.QueryParameters != nil {
		for key, value := range reqConfig.QueryParameters {
			requestInfo.QueryParametersAny[key] = value
		}
	}

	nativeReq, err := adapter.ConvertToNativeRequest(ctx, requestInfo)
	if err != nil {
		return nil, fmt.Errorf("error converting to native request: %w", err)
	}

	httpReq := nativeReq.(*http.Request)
	client := &http.Client{}

	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("error executing request: %w", err)
	}
	defer resp.Body.Close()

	tflog.Debug(ctx, "Request URL", map[string]interface{}{
		"url": httpReq.URL.String(),
	})

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %w", err)
	}

	return body, nil
}
