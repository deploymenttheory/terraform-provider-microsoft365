package customrequests

import (
	"context"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	abstractions "github.com/microsoft/kiota-abstractions-go"
	s "github.com/microsoft/kiota-abstractions-go/serialization"
)

// GraphAPIVersion represents Microsoft Graph API version
type GraphAPIVersion string

const (
	GraphAPIBeta GraphAPIVersion = "beta"
	GraphAPIV1   GraphAPIVersion = "v1.0"
)

// PutRequestConfig contains the configuration for a custom PUT request
type PutRequestConfig struct {
	// The API version to use (beta or v1.0)
	APIVersion GraphAPIVersion
	// The base endpoint (e.g., "/deviceManagement/configurationPolicies")
	Endpoint string
	// The ID of the resource
	ResourceID string
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
}

type PutResponse struct {
	StatusCode int
	Error      error
}

// PutRequestByResourceId performs a custom PUT request using the Microsoft Graph SDK when the operation
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
//	err := PutRequestByResourceId(ctx, clients, config)
func PutRequestByResourceId(ctx context.Context, adapter abstractions.RequestAdapter, config PutRequestConfig) error {
	requestInfo := newGraphRequestInformation(adapter, abstractions.PUT, endpointUrlTemplate(config.Endpoint)+"('{id}')", config.APIVersion, config.QueryParameters, config.Headers)
	requestInfo.PathParameters["id"] = config.ResourceID

	tflog.Debug(ctx, "Making custom msgraph PUT request", map[string]any{"url_template": requestInfo.UrlTemplate, "id": config.ResourceID})

	if err := setRequestContent(ctx, adapter, requestInfo, config.RequestBody, config.RawBody); err != nil {
		return err
	}

	if err := sendForRawResponseOrNoContent(ctx, adapter, requestInfo, config.RawResponse, nil); err != nil {
		return err
	}

	return nil
}
