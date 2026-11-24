package customrequests

import (
	"context"

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
	// Optional raw JSON request body: a json.RawMessage or []byte sent as it is, or any value marshalled with
	// encoding/json. Use instead of RequestBody; setting both is an error.
	RawBody any
	// Optional query parameters to include in the request (e.g. {"$select": "id"})
	QueryParameters map[string]string
	// Optional request headers (e.g. {"Prefer": "include-unknown-enum-members"})
	Headers map[string]string
	// Optional raw JSON response target: a pointer the response body is decoded into with encoding/json. A response
	// with no body (e.g. 204 No Content) leaves it untouched.
	RawResponse any
	// Optional per-request middleware options, e.g. a kiota RetryHandlerOptions that narrows which responses the
	// retry middleware retries for this request only.
	RequestOptions []abstractions.RequestOption
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
	requestInfo := newGraphRequestInformation(adapter, abstractions.PATCH, ByIDRequestUrlTemplate(config), config.APIVersion, config.QueryParameters, config.Headers)
	requestInfo.AddRequestOptions(config.RequestOptions)

	tflog.Debug(ctx, "Making custom msgraph PATCH request", map[string]any{"url": requestInfo.UrlTemplate})

	if err := setRequestContent(ctx, adapter, requestInfo, config.RequestBody, config.RawBody); err != nil {
		return err
	}

	if err := sendForRawResponseOrNoContent(ctx, adapter, requestInfo, config.RawResponse, nil); err != nil {
		return err
	}

	return nil
}
