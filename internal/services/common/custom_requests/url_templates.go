package customrequests

import (
	"fmt"
	"sort"
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

// PutRequestConfig implements RequestUrlTemplateConfig interface
func (p PutRequestConfig) GetResourceIDPattern() string { return p.ResourceIDPattern }
func (p PutRequestConfig) GetResourceID() string        { return p.ResourceID }
func (p PutRequestConfig) GetEndpoint() string          { return p.Endpoint }
func (p PutRequestConfig) GetEndpointSuffix() string    { return p.EndpointSuffix }

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
// queryParameters are appended as a URI template query expansion so Kiota puts them on the wire (see
// queryExpansionTemplate), and headers are added to the request. Both may be nil.
func newGraphRequestInformation(
	adapter abstractions.RequestAdapter,
	method abstractions.HttpMethod,
	urlTemplate string,
	apiVersion GraphAPIVersion,
	queryParameters map[string]string,
	headers map[string]string,
) *abstractions.RequestInformation {
	resolvedTemplate := strings.Replace(urlTemplate, "{+baseurl}", graphBaseURL(adapter, apiVersion), 1)
	queryTemplate, queryVariables := queryExpansionTemplate(resolvedTemplate, queryParameters)

	requestInfo := abstractions.NewRequestInformationWithMethodAndUrlTemplateAndPathParameters(method, resolvedTemplate+queryTemplate, map[string]string{})
	for name, value := range queryVariables {
		requestInfo.QueryParameters[name] = value
	}
	addRequestHeaders(requestInfo, headers)

	return requestInfo
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

// queryExpansionTemplate returns the RFC 6570 form-style query expansion (e.g. "{?%24expand,%24select}") for
// queryParameters, and the template variables to set on the RequestInformation. Kiota only puts query parameters
// on the wire when the URL template names them, so the expansion is what makes QueryParameters reach the URL.
// Template variable names may only hold letters, digits and "_", so every other byte of a parameter name is
// percent-encoded (e.g. "$expand" becomes "%24expand"), exactly as the generated SDK request builders name OData
// query options. Keys are sorted so the URL is deterministic. When urlTemplate already holds a literal query string,
// a form-style continuation ("{&...}") is used instead.
func queryExpansionTemplate(urlTemplate string, queryParameters map[string]string) (string, map[string]string) {
	if len(queryParameters) == 0 {
		return "", nil
	}

	keys := make([]string, 0, len(queryParameters))
	for key := range queryParameters {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	names := make([]string, 0, len(keys))
	variables := make(map[string]string, len(keys))
	for _, key := range keys {
		name := uriTemplateVariableName(key)
		names = append(names, name)
		variables[name] = queryParameters[key]
	}

	operator := "?"
	if strings.Contains(urlTemplate, "?") {
		operator = "&"
	}

	return "{" + operator + strings.Join(names, ",") + "}", variables
}

// uriTemplateVariableName percent-encodes every byte of name that is not a letter, digit or "_".
func uriTemplateVariableName(name string) string {
	var builder strings.Builder
	for i := 0; i < len(name); i++ {
		c := name[i]
		if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '_' {
			builder.WriteByte(c)
			continue
		}
		fmt.Fprintf(&builder, "%%%02X", c)
	}
	return builder.String()
}
