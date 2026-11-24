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

// Resource ID patterns accepted by ResourceIDPattern. Each caller chooses the format its endpoint needs; the
// choice is deliberate per endpoint family and is never inferred.
const (
	// ResourceIDPatternParentheses addresses the resource with an OData parentheses key: "<endpoint>('<id>')".
	// It is the original format of the custom requests (ShocOne, 1b060e6a, Nov 2024), written for the Intune
	// settings catalog (e.g. deviceManagement/configurationPolicies('<id>')/settings) and used by the Intune
	// policy, inventory policy and enrollment policy callers.
	ResourceIDPatternParentheses = "('id')"

	// ResourceIDPatternSegment addresses the resource with an OData key segment: "<endpoint>/<id>".
	// ShocOne added it with the explicit PATCH switch (dccde7e7, Apr 2025) for device categories, and kenchan0130
	// copied it to DELETE (47ce3efb). It is the form Microsoft Learn documents for directory navigation references,
	// e.g. DELETE /servicePrincipals/{id}/microsoft.graph.agentIdentity/sponsors/{sponsorObjectId}/$ref, and is used
	// by the owner, sponsor, credential and localized notification message callers.
	ResourceIDPatternSegment = "/{id}"

	// ResourceIDPatternNone appends no ID: the endpoint is a singleton or already contains the ID.
	ResourceIDPatternNone = ""
)

// ByIDRequestUrlTemplate constructs a URL template for a single resource request using the provided configuration.
// The endpoint is always joined to the base URL with exactly one slash, then the resource ID is written in the
// format named by ResourceIDPattern, then the optional suffix is appended:
//   - ResourceIDPatternParentheses "('id')": "{+baseurl}/<endpoint>('<id>')<suffix>"
//   - ResourceIDPatternSegment "/{id}":      "{+baseurl}/<endpoint>/<id><suffix>"
//   - ResourceIDPatternNone "":              "{+baseurl}/<endpoint><suffix>"
//   - any other pattern (e.g. "(id)"):       every "id" in the pattern is replaced with the resource ID
//
// This is the PATCH/DELETE switch (dccde7e7, 47ce3efb) in one place, so GET, PATCH and DELETE build identical
// URLs for the same config. The resource ID is written into the template literally (it is not escaped), exactly
// as before, so every existing caller's URL is unchanged.
//
// For example, Endpoint "/deviceManagement/configurationPolicies", ResourceIDPattern "('id')", ResourceID "12345"
// and EndpointSuffix "/settings" give "{+baseurl}/deviceManagement/configurationPolicies('12345')/settings".
//
// Parameters:
//   - reqConfig: a RequestUrlTemplateConfig containing the endpoint path, resource ID pattern, actual ID, and optional suffix
//
// Returns:
//   - string: The constructed URL template ready for use with the Kiota request adapter
func ByIDRequestUrlTemplate(reqConfig RequestUrlTemplateConfig) string {
	urlTemplate := endpointUrlTemplate(reqConfig.GetEndpoint())

	switch reqConfig.GetResourceIDPattern() {
	case ResourceIDPatternSegment:
		urlTemplate += "/" + reqConfig.GetResourceID()
	case ResourceIDPatternParentheses:
		urlTemplate += "('" + reqConfig.GetResourceID() + "')"
	case ResourceIDPatternNone:
	default:
		urlTemplate += strings.ReplaceAll(reqConfig.GetResourceIDPattern(), "id", reqConfig.GetResourceID())
	}

	return urlTemplate + reqConfig.GetEndpointSuffix()
}
