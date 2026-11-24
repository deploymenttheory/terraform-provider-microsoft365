package customrequests

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	abstractions "github.com/microsoft/kiota-abstractions-go"
	s "github.com/microsoft/kiota-abstractions-go/serialization"
)

// PatchRequestConfig contains the configuration for a custom PATCH request
type PatchRequestConfig struct {
	// The API version to use (beta or v1.0)
	APIVersion GraphAPIVersion
	// The base endpoint (e.g., "/deviceManagement/deviceCategories")
	Endpoint string
	// The ID of the resource
	ResourceID string
	// Pattern for the resource ID (e.g., "/{id}" or "('id')")
	ResourceIDPattern string
	// Optional suffix to append after the ID
	EndpointSuffix string
	// The request body
	RequestBody s.Parsable
}

type PatchResponse struct {
	StatusCode int
	Error      error
}

// PatchRequestByResourceId performs a custom PATCH request using the Microsoft Graph SDK when the operation
// is not available in the generated SDK methods. This function supports both Beta and V1.0 Graph API versions
// and expects a 204 No Content response from the server on success.
//
// e.g., PATCH https://graph.microsoft.com/beta/deviceManagement/deviceCategories('d557c813-b8e5-4efc-b00e-9c0bd5fd10df')
//
// The function handles:
// - Construction of the Graph API URL with proper formatting
// - Setting up the PATCH request with the provided request body
// - Sending the request with proper authentication
// - Error handling and return
//
// Parameters:
//   - ctx: The context for the request, which can be used for cancellation and timeout
//   - adapter: The RequestAdapter instance containing client configurations
//   - config: PatchRequestConfig containing:
//   - APIVersion: The Graph API version to use (Beta or V1.0)
//   - Endpoint: The resource endpoint path (e.g., "deviceManagement/deviceCategories")
//   - ResourceID: The ID of the resource to update
//   - ResourceIDPattern: The pattern for formatting the resource ID (e.g., "/{id}" or "('id')")
//   - EndpointSuffix: Optional suffix to append after the ID
//   - RequestBody: The body of the PATCH request implementing serialization.Parsable
//
// Returns:
//   - error: Returns nil if the request was successful (204 No Content received),
//     otherwise returns an error describing what went wrong
//
// Example Usage:
//
//	config := PatchRequestConfig{
//	    APIVersion: GraphAPIBeta,
//	    Endpoint:   "/deviceManagement/deviceCategories",
//	    ResourceIDPattern: "/{id}",
//	    ResourceID: "d557c813-b8e5-4efc-b00e-9c0bd5fd10df",
//	    RequestBody: myRequestBody,
//	}
//	err := PatchRequestByResourceId(ctx, adapter, config)
func PatchRequestByResourceId(ctx context.Context, adapter abstractions.RequestAdapter, config PatchRequestConfig) error {
	requestInfo := newGraphRequestInformation(adapter, abstractions.PATCH, ByIDRequestUrlTemplate(config), config.APIVersion)

	tflog.Debug(ctx, "Making custom msgraph PATCH request", map[string]any{"url": requestInfo.UrlTemplate})

	err := requestInfo.SetContentFromParsable(ctx, adapter, "application/json", config.RequestBody)
	if err != nil {
		return fmt.Errorf("error setting content: %v", err)
	}

	if err := adapter.SendNoContent(ctx, requestInfo, withODataErrorMappings(nil)); err != nil {
		return fmt.Errorf("custom msgraph PATCH request failed: %w", err)
	}

	return nil
}
