package client

import (
	"context"
	"fmt"

	abstractions "github.com/microsoft/kiota-abstractions-go"
	s "github.com/microsoft/kiota-abstractions-go/serialization"
)

// CustomPostRequestConfig contains the configuration for a custom POST request
type CustomPostRequestConfig struct {
	// The API version to use (beta or v1.0)
	APIVersion GraphAPIVersion
	// The base endpoint (e.g., "deviceManagement/configurationPolicies")
	Endpoint string
	// The request body
	RequestBody s.Parsable
	// Optional query parameters for the request
	QueryParameters map[string]string
}

// SendCustomPostRequest performs a custom POST request using the Microsoft Graph SDK when the operation
// is not available in the generated SDK methods. This function supports both Beta and V1.0 Graph API versions
// and returns the parsed response model.
//
// Parameters:
//   - ctx: The context for the request, which can be used for cancellation and timeout
//   - adapter: The RequestAdapter interface for making HTTP requests
//   - config: CustomPostRequestConfig containing:
//   - APIVersion: The Graph API version to use (Beta or V1.0)
//   - Endpoint: The resource endpoint path
//   - RequestBody: The body of the POST request implementing serialization.Parsable
//   - QueryParameters: Optional map of query parameters
//   - factory: The factory function to create the response model
//   - errorMappings: Optional error mappings for custom error handling
//
// Returns:
//   - s.Parsable: The parsed response model
//   - error: Any error that occurred during the request
func SendCustomPostRequest(
	ctx context.Context,
	adapter abstractions.RequestAdapter,
	config CustomPostRequestConfig,
	factory s.ParsableFactory,
	errorMappings abstractions.ErrorMappings,
) (s.Parsable, error) {
	requestInfo := abstractions.NewRequestInformation()
	requestInfo.Method = abstractions.POST
	requestInfo.UrlTemplate = "{+baseurl}/" + config.Endpoint

	requestInfo.PathParameters = map[string]string{
		"baseurl": fmt.Sprintf("https://graph.microsoft.com/%s", config.APIVersion),
	}

	// Add query parameters if provided
	if config.QueryParameters != nil {
		for key, value := range config.QueryParameters {
			requestInfo.QueryParameters[key] = value
		}
	}

	err := requestInfo.SetContentFromParsable(ctx, adapter, "application/json", config.RequestBody)
	if err != nil {
		return nil, fmt.Errorf("error setting content: %v", err)
	}

	// Send the request using the adapter's Send method
	result, err := adapter.Send(ctx, requestInfo, factory, errorMappings)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %v", err)
	}

	return result, nil
}

// SendCustomPostRequestNoContent performs a custom POST request that doesn't expect a response body.
// This is useful for operations that return 204 No Content.
//
// Parameters are the same as SendCustomPostRequest except it doesn't take a responseModel parameter
// and uses the SendNoContent method of the adapter.
//
// Returns:
//   - error: Returns nil if the request was successful (204 No Content received),
//     otherwise returns an error describing what went wrong
func SendCustomPostRequestNoContent(ctx context.Context, adapter abstractions.RequestAdapter, config CustomPostRequestConfig) error {
	requestInfo := abstractions.NewRequestInformation()
	requestInfo.Method = abstractions.POST
	requestInfo.UrlTemplate = "{+baseurl}/" + config.Endpoint

	requestInfo.PathParameters = map[string]string{
		"baseurl": fmt.Sprintf("https://graph.microsoft.com/%s", config.APIVersion),
	}

	if config.QueryParameters != nil {
		for key, value := range config.QueryParameters {
			requestInfo.QueryParameters[key] = value
		}
	}

	err := requestInfo.SetContentFromParsable(ctx, adapter, "application/json", config.RequestBody)
	if err != nil {
		return fmt.Errorf("error setting content: %v", err)
	}

	// Use SendNoContent for requests that don't return a response body
	err = adapter.SendNoContent(ctx, requestInfo, nil)
	if err != nil {
		return fmt.Errorf("error sending post request: %v", err)
	}

	return nil
}