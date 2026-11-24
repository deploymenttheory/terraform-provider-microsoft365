package customrequests

import (
	"context"

	abstractions "github.com/microsoft/kiota-abstractions-go"
)

// DeleteRequestConfig contains the configuration for a custom DELETE request
type DeleteRequestConfig struct {
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

// DeleteRequestByResourceId performs a custom DELETE request using the Microsoft Graph SDK
// when the operation is not available in the generated SDK methods.
//
// e.g., DELETE https://graph.microsoft.com/beta/deviceManagement/reusablePolicySettings('191056b1-4c4a-4871-8518-162a105d011a')
//
// Parameters:
//   - ctx: The context for the request, which can be used for cancellation and timeout
//   - adapter: The request adapter for sending the request
//   - config: DeleteRequestConfig containing:
//   - APIVersion: The Graph API version to use (Beta or V1.0)
//   - Endpoint: The resource endpoint path (e.g., "/deviceManagement/configurationPolicies")
//   - ResourceID: The ID of the resource to delete
//   - ResourceIDPattern: The format for the resource ID (e.g., "('id')" or "(id)")
//   - EndpointSuffix: Optional suffix to append after the resource ID
//
// Returns:
//   - error: Returns nil if the delete was successful, otherwise an error describing what went wrong
//
// Example Usage:
//
//	config := DeleteRequestConfig{
//		APIVersion:        GraphAPIBeta,
//		Endpoint:         "deviceManagement/reusablePolicySettings",
//		ResourceID:       "4f93da7a-431b-4d48-bbce-1660aa8b0be7",
//		ResourceIDPattern: "('id')",
//	}
//
//	err := DeleteRequestByResourceId(ctx, adapter, config)
//	if err != nil {
//		log.Fatalf("Error: %v", err)
//	}
func DeleteRequestByResourceId(ctx context.Context, adapter abstractions.RequestAdapter, reqConfig DeleteRequestConfig) error {
	requestInfo := newGraphRequestInformation(adapter, abstractions.DELETE, ByIDRequestUrlTemplate(reqConfig), reqConfig.APIVersion, reqConfig.QueryParameters, reqConfig.Headers)
	requestInfo.AddRequestOptions(reqConfig.RequestOptions)
	requestInfo.Headers.TryAdd("Accept", jsonContentType)

	if err := sendForRawResponseOrNoContent(ctx, adapter, requestInfo, reqConfig.RawResponse, nil); err != nil {
		return err
	}

	return nil
}
