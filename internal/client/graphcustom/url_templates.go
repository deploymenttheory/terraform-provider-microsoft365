package graphcustom

import "strings"

// ByIDRequestUrlTemplate constructs a URL template for a single resource request using the provided configuration.
// The function combines the endpoint path with a resource ID and optional suffix to create a complete URL template.
// For example, if the config contains:
//   - Endpoint: "/deviceManagement/configurationPolicies"
//   - ResourceIDPattern: "('id')"
//   - ResourceID: "12345"
//   - EndpointSuffix: "/settings"
//
// The resulting template would be: "{+baseurl}/deviceManagement/configurationPolicies('12345')/settings"
//
// Parameters:
//   - reqConfig: GetRequestConfig containing the endpoint path, resource ID pattern, actual ID, and optional suffix
//
// Returns:
//   - string: The constructed URL template ready for use with the Kiota request adapter
func ByIDRequestUrlTemplate(reqConfig GetRequestConfig) string {
	idFormat := strings.ReplaceAll(reqConfig.ResourceIDPattern, "id", reqConfig.ResourceID)
	endpoint := reqConfig.Endpoint + idFormat
	if reqConfig.EndpointSuffix != "" {
		endpoint += reqConfig.EndpointSuffix
	}
	return "{+baseurl}" + endpoint
}
