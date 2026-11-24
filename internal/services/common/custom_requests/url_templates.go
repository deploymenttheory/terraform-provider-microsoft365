package customrequests

import (
	"strings"

	abstractions "github.com/microsoft/kiota-abstractions-go"
	"github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// defaultGraphHost is used only when the adapter has no base URL configured.
const defaultGraphHost = "https://graph.microsoft.com"

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

// Implement RequestUrlTemplateConfig interface for PatchRequestConfig
func (p PatchRequestConfig) GetResourceIDPattern() string { return p.ResourceIDPattern }
func (p PatchRequestConfig) GetResourceID() string        { return p.ResourceID }
func (p PatchRequestConfig) GetEndpoint() string          { return p.Endpoint }
func (p PatchRequestConfig) GetEndpointSuffix() string    { return p.EndpointSuffix }

// graphBaseURL returns the adapter's base URL with its version segment set to apiVersion.
// The adapter's host is kept, so national clouds (e.g. graph.microsoft.us) are honoured.
// e.g. adapter base "https://graph.microsoft.us/v1.0" with GraphAPIBeta returns "https://graph.microsoft.us/beta".
// When the adapter has no base URL, the public cloud host is used. When apiVersion is empty, the adapter's version is kept.
func graphBaseURL(adapter abstractions.RequestAdapter, apiVersion GraphAPIVersion) string {
	base := strings.TrimSuffix(adapter.GetBaseUrl(), "/")
	if base == "" {
		base = defaultGraphHost + "/" + string(GraphAPIV1)
	}

	if apiVersion == "" {
		return base
	}

	for _, version := range []GraphAPIVersion{GraphAPIBeta, GraphAPIV1} {
		if strings.HasSuffix(base, "/"+string(version)) {
			base = strings.TrimSuffix(base, "/"+string(version))
			break
		}
	}

	return base + "/" + string(apiVersion)
}

// newGraphRequestInformation builds a RequestInformation whose {+baseurl} resolves to the adapter's host and the given API version.
// The resolved base URL is written into the template literally because the Kiota adapter overwrites the "baseurl" path
// parameter with its own base URL (and so its own API version) on every send.
func newGraphRequestInformation(
	adapter abstractions.RequestAdapter,
	method abstractions.HttpMethod,
	urlTemplate string,
	apiVersion GraphAPIVersion,
) *abstractions.RequestInformation {
	resolvedTemplate := strings.Replace(urlTemplate, "{+baseurl}", graphBaseURL(adapter, apiVersion), 1)
	return abstractions.NewRequestInformationWithMethodAndUrlTemplateAndPathParameters(method, resolvedTemplate, map[string]string{})
}

// withODataErrorMappings returns errorMappings, or a mapping of every 4xx/5xx status to an ODataError when none are given.
func withODataErrorMappings(errorMappings abstractions.ErrorMappings) abstractions.ErrorMappings {
	if len(errorMappings) > 0 {
		return errorMappings
	}
	return abstractions.ErrorMappings{
		"XXX": odataerrors.CreateODataErrorFromDiscriminatorValue,
	}
}

// endpointUrlTemplate returns "{+baseurl}/<endpoint>", ensuring exactly one slash between the base URL and the endpoint.
func endpointUrlTemplate(endpoint string) string {
	return "{+baseurl}/" + strings.TrimPrefix(endpoint, "/")
}

// ByIDRequestUrlTemplate constructs a URL template for a single resource request using the provided configuration.
// The function combines the endpoint path with a resource ID and optional suffix to create a complete URL template.
// The resource ID is written into the template literally, so it is never expanded as a template variable.
// Supported ID patterns include "('id')", "(id)" and "/{id}". For example, if the config contains:
//   - Endpoint: "/deviceManagement/configurationPolicies"
//   - ResourceIDPattern: "('id')"
//   - ResourceID: "12345"
//   - EndpointSuffix: "/settings"
//
// The resulting template would be: "{+baseurl}/deviceManagement/configurationPolicies('12345')/settings"
// With ResourceIDPattern "/{id}" it would be: "{+baseurl}/deviceManagement/configurationPolicies/12345/settings"
//
// Parameters:
//   - reqConfig: a RequestUrlTemplateConfig containing the endpoint path, resource ID pattern, actual ID, and optional suffix
//
// Returns:
//   - string: The constructed URL template ready for use with the Kiota request adapter
func ByIDRequestUrlTemplate(reqConfig RequestUrlTemplateConfig) string {
	pattern := reqConfig.GetResourceIDPattern()
	var idFormat string
	if strings.Contains(pattern, "{id}") {
		idFormat = strings.ReplaceAll(pattern, "{id}", reqConfig.GetResourceID())
	} else {
		idFormat = strings.ReplaceAll(pattern, "id", reqConfig.GetResourceID())
	}

	return endpointUrlTemplate(reqConfig.GetEndpoint()) + idFormat + reqConfig.GetEndpointSuffix()
}
