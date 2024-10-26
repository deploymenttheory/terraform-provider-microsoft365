package client

import (
	"context"
	"fmt"

	abstractions "github.com/microsoft/kiota-abstractions-go"
	s "github.com/microsoft/kiota-abstractions-go/serialization"
)

// GraphAPIVersion represents Microsoft Graph API version
type GraphAPIVersion string

const (
	GraphAPIBeta GraphAPIVersion = "beta"
	GraphAPIV1   GraphAPIVersion = "v1.0"
)

// CustomPutRequestConfig contains the configuration for a custom PUT request
type CustomPutRequestConfig struct {
	// The API version to use (beta or v1.0)
	APIVersion GraphAPIVersion
	// The base endpoint (e.g., "deviceManagement/configurationPolicies")
	Endpoint string
	// The ID of the resource
	ResourceID string
	// The request body
	RequestBody s.Parsable
}

// SendCustomPutRequestByResourceId performs a custom PUT request using the Microsoft Graph SDK when the operation
// is not available in the generated SDK methods. This function supports both Beta and V1.0 Graph API versions
// and expects a 204 No Content response from the server on success.
//
// e.g., PUT https://graph.microsoft.com/beta/deviceManagement/configurationPolicies('d557c813-b8e5-4efc-b00e-9c0bd5fd10df')
//
// The function handles:
// - Construction of the Graph API URL with proper formatting
// - Setting up the PUT request with the provided request body
// - Sending the request with proper authentication
// - Error handling and return
//
// Parameters:
//   - ctx: The context for the request, which can be used for cancellation and timeout
//   - clients: The GraphClients instance containing both Beta and V1.0 client configurations
//   - config: CustomPutRequestConfig containing:
//   - APIVersion: The Graph API version to use (Beta or V1.0)
//   - Endpoint: The resource endpoint path (e.g., "deviceManagement/configurationPolicies")
//   - ResourceID: The ID of the resource to update
//   - RequestBody: The body of the PUT request implementing serialization.Parsable
//
// Returns:
//   - error: Returns nil if the request was successful (204 No Content received),
//     otherwise returns an error describing what went wrong
//
// Example Usage:
//
//	config := CustomPutRequestConfig{
//	    APIVersion: GraphAPIBeta,
//	    Endpoint:   "deviceManagement/configurationPolicies",
//	    ResourceID: "d557c813-b8e5-4efc-b00e-9c0bd5fd10df",
//	    RequestBody: myRequestBody,
//	}
//	err := SendCustomPutRequestByResourceId(ctx, clients, config)
func SendCustomPutRequestByResourceId(ctx context.Context, adapter abstractions.RequestAdapter, config CustomPutRequestConfig) error {
	requestInfo := abstractions.NewRequestInformation()
	requestInfo.Method = abstractions.PUT
	requestInfo.UrlTemplate = "{+baseurl}/" + config.Endpoint + "('{id}')"
	requestInfo.PathParameters = map[string]string{
		"baseurl": fmt.Sprintf("https://graph.microsoft.com/%s", config.APIVersion),
		"id":      config.ResourceID,
	}

	err := requestInfo.SetContentFromParsable(ctx, adapter, "application/json", config.RequestBody)
	if err != nil {
		return fmt.Errorf("error setting content: %v", err)
	}

	err = adapter.SendNoContent(ctx, requestInfo, nil)
	if err != nil {
		return fmt.Errorf("error sending request: %v", err)
	}

	return nil
}
