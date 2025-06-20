package customrequests

import (
	"context"
	"fmt"

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
	requestInfo := abstractions.NewRequestInformation()
	requestInfo.Method = abstractions.DELETE
	requestInfo.UrlTemplate = ByIDRequestUrlTemplate(reqConfig)
	requestInfo.PathParameters = map[string]string{
		"baseurl": fmt.Sprintf("https://graph.microsoft.com/%s", reqConfig.APIVersion),
	}
	requestInfo.Headers.Add("Accept", "application/json")

	_, err := makeRequest(ctx, adapter, requestInfo)
	if err != nil {
		return err
	}

	return nil
}
