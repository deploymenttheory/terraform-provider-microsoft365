package graphcustom

import "strings"

// RequestUrlTemplateConfig interface defines the common fields needed for building URLs
type RequestUrlTemplateConfig interface {
	GetResourceIDPattern() string
	GetResourceID() string
	GetEndpoint() string
	GetEndpointSuffix() string
}

// GetRequestConfig implements RequestUrlConfig interface
func (g GetRequestConfig) GetResourceIDPattern() string { return g.ResourceIDPattern }
func (g GetRequestConfig) GetResourceID() string        { return g.ResourceID }
func (g GetRequestConfig) GetEndpoint() string          { return g.Endpoint }
func (g GetRequestConfig) GetEndpointSuffix() string    { return g.EndpointSuffix }

// DeleteRequestConfig implements RequestUrlConfig interface
func (d DeleteRequestConfig) GetResourceIDPattern() string { return d.ResourceIDPattern }
func (d DeleteRequestConfig) GetResourceID() string        { return d.ResourceID }
func (d DeleteRequestConfig) GetEndpoint() string          { return d.Endpoint }
func (d DeleteRequestConfig) GetEndpointSuffix() string    { return d.EndpointSuffix }

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
func ByIDRequestUrlTemplate(reqConfig RequestUrlTemplateConfig) string {
	idFormat := strings.ReplaceAll(reqConfig.GetResourceIDPattern(), "id", reqConfig.GetResourceID())
	endpoint := reqConfig.GetEndpoint() + idFormat
	if reqConfig.GetEndpointSuffix() != "" {
		endpoint += reqConfig.GetEndpointSuffix()
	}
	return "{+baseurl}" + endpoint
}
