package customrequests

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	abstractions "github.com/microsoft/kiota-abstractions-go"
	"github.com/microsoft/kiota-abstractions-go/authentication"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models"
	"github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testResourceID = "d557c813-b8e5-4efc-b00e-9c0bd5fd10df"
	usGovV1Root    = "https://graph.microsoft.us/v1.0"
	usGovBetaRoot  = "https://graph.microsoft.us/beta"
)

// newTestAdapter returns a Graph beta request adapter whose HTTP client sends through an isolated httpmock transport.
func newTestAdapter(t *testing.T, baseURL string) (abstractions.RequestAdapter, *httpmock.MockTransport) {
	t.Helper()

	transport := httpmock.NewMockTransport()
	adapter, err := msgraphbetasdk.NewGraphRequestAdapterWithParseNodeFactoryAndSerializationWriterFactoryAndHttpClient(&authentication.AnonymousAuthenticationProvider{}, nil, nil, &http.Client{Transport: transport})
	require.NoError(t, err)

	// Creating a service client registers the default JSON serializers, as the provider does.
	msgraphbetasdk.NewGraphServiceClient(adapter)
	adapter.SetBaseUrl(baseURL)

	return adapter, transport
}

func newTestBody() *models.DeviceCategory {
	body := models.NewDeviceCategory()
	name := "test"
	body.SetDisplayName(&name)
	return body
}

func odataErrorResponder(status int, code string) httpmock.Responder {
	return httpmock.NewStringResponder(status, `{"error":{"code":"`+code+`","message":"`+code+` message"}}`).HeaderSet(http.Header{"Content-Type": []string{"application/json"}})
}

func jsonResponder(status int, body string) httpmock.Responder {
	return httpmock.NewStringResponder(status, body).HeaderSet(http.Header{"Content-Type": []string{"application/json"}})
}

func requireODataError(t *testing.T, err error, status int, code string) {
	t.Helper()

	require.Error(t, err)
	var odataErr *odataerrors.ODataError
	require.True(t, errors.As(err, &odataErr), "expected an ODataError, got %T: %v", err, err)
	assert.Equal(t, status, odataErr.ResponseStatusCode)
	require.NotNil(t, odataErr.GetErrorEscaped())
	assert.Equal(t, code, *odataErr.GetErrorEscaped().GetCode())
}

func TestGraphBaseURL(t *testing.T) {
	tests := []struct {
		name       string
		baseURL    string
		apiVersion GraphAPIVersion
		expected   string
	}{
		{name: "public v1 adapter to beta", baseURL: "https://graph.microsoft.com/v1.0", apiVersion: GraphAPIBeta, expected: "https://graph.microsoft.com/beta"},
		{name: "public beta adapter to v1", baseURL: "https://graph.microsoft.com/beta", apiVersion: GraphAPIV1, expected: "https://graph.microsoft.com/v1.0"},
		{name: "us gov adapter keeps host", baseURL: usGovV1Root, apiVersion: GraphAPIBeta, expected: usGovBetaRoot},
		{name: "china adapter keeps host", baseURL: "https://microsoftgraph.chinacloudapi.cn/beta", apiVersion: GraphAPIBeta, expected: "https://microsoftgraph.chinacloudapi.cn/beta"},
		{name: "trailing slash", baseURL: "https://graph.microsoft.us/beta/", apiVersion: GraphAPIV1, expected: usGovV1Root},
		{name: "base without version", baseURL: "https://graph.microsoft.us", apiVersion: GraphAPIBeta, expected: usGovBetaRoot},
		{name: "empty version keeps adapter version", baseURL: usGovBetaRoot, apiVersion: "", expected: usGovBetaRoot},
		{name: "empty adapter base falls back to public cloud", baseURL: "", apiVersion: GraphAPIBeta, expected: "https://graph.microsoft.com/beta"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			adapter, _ := newTestAdapter(t, tt.baseURL)
			assert.Equal(t, tt.expected, graphBaseURL(adapter, tt.apiVersion))
		})
	}
}

// TestConfigFunctionsHonourAdapterBaseURL proves every Config function sends to the adapter's host (a national cloud here) with the configured API version.
func TestConfigFunctionsHonourAdapterBaseURL(t *testing.T) {
	ctx := context.Background()
	itemURL := usGovBetaRoot + "/deviceManagement/deviceCategories/" + testResourceID
	collectionURL := usGovBetaRoot + "/deviceManagement/deviceCategories"
	categoryJSON := `{"id":"` + testResourceID + `","displayName":"test"}`

	tests := []struct {
		name   string
		method string
		url    string
		status int
		body   string
		call   func(adapter abstractions.RequestAdapter) error
	}{
		{
			name: "PostRequest", method: http.MethodPost, url: collectionURL, status: http.StatusCreated, body: categoryJSON,
			call: func(adapter abstractions.RequestAdapter) error {
				result, err := PostRequest(ctx, adapter, PostRequestConfig{APIVersion: GraphAPIBeta, Endpoint: "/deviceManagement/deviceCategories", RequestBody: newTestBody()}, models.CreateDeviceCategoryFromDiscriminatorValue, nil)
				if err == nil {
					assert.Equal(t, testResourceID, *result.(models.DeviceCategoryable).GetId())
				}
				return err
			},
		},
		{
			name: "PostRequestNoContent", method: http.MethodPost, url: collectionURL, status: http.StatusNoContent,
			call: func(adapter abstractions.RequestAdapter) error {
				return PostRequestNoContent(ctx, adapter, PostRequestConfig{APIVersion: GraphAPIBeta, Endpoint: "deviceManagement/deviceCategories", RequestBody: newTestBody()})
			},
		},
		{
			name: "GetRequest", method: http.MethodGet, url: itemURL, status: http.StatusOK, body: categoryJSON,
			call: func(adapter abstractions.RequestAdapter) error {
				_, err := GetRequest(ctx, adapter, GetRequestConfig{APIVersion: GraphAPIBeta, Endpoint: "deviceManagement/deviceCategories", ResourceID: testResourceID}, models.CreateDeviceCategoryFromDiscriminatorValue, nil)
				return err
			},
		},
		{
			name: "GetRequestByResourceId", method: http.MethodGet, url: itemURL, status: http.StatusOK, body: categoryJSON,
			call: func(adapter abstractions.RequestAdapter) error {
				body, err := GetRequestByResourceId(ctx, adapter, GetRequestConfig{APIVersion: GraphAPIBeta, Endpoint: "/deviceManagement/deviceCategories", ResourceIDPattern: "/{id}", ResourceID: testResourceID})
				if err == nil {
					assert.JSONEq(t, categoryJSON, string(body))
				}
				return err
			},
		},
		{
			name: "PatchRequestByResourceId", method: http.MethodPatch, url: itemURL, status: http.StatusNoContent,
			call: func(adapter abstractions.RequestAdapter) error {
				return PatchRequestByResourceId(ctx, adapter, PatchRequestConfig{APIVersion: GraphAPIBeta, Endpoint: "/deviceManagement/deviceCategories", ResourceIDPattern: "/{id}", ResourceID: testResourceID, RequestBody: newTestBody()})
			},
		},
		{
			name: "PutRequestByResourceId", method: http.MethodPut, url: collectionURL + "('" + testResourceID + "')", status: http.StatusNoContent,
			call: func(adapter abstractions.RequestAdapter) error {
				return PutRequestByResourceId(ctx, adapter, PutRequestConfig{APIVersion: GraphAPIBeta, Endpoint: "/deviceManagement/deviceCategories", ResourceID: testResourceID, RequestBody: newTestBody()})
			},
		},
		{
			name: "DeleteRequestByResourceId", method: http.MethodDelete, url: itemURL, status: http.StatusNoContent,
			call: func(adapter abstractions.RequestAdapter) error {
				return DeleteRequestByResourceId(ctx, adapter, DeleteRequestConfig{APIVersion: GraphAPIBeta, Endpoint: "/deviceManagement/deviceCategories", ResourceIDPattern: "/{id}", ResourceID: testResourceID})
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// The adapter is configured for v1.0 in the US Gov cloud; the Config asks for beta.
			adapter, transport := newTestAdapter(t, usGovV1Root)
			transport.RegisterResponder(tt.method, tt.url, jsonResponder(tt.status, tt.body))

			require.NoError(t, tt.call(adapter))
			assert.Equal(t, 1, transport.GetCallCountInfo()[tt.method+" "+tt.url], "request must go to the adapter's host, calls: %v", transport.GetCallCountInfo())
			assert.Equal(t, 1, transport.GetTotalCallCount())
		})
	}
}

// TestConfigFunctionsReturnODataErrors proves a 4xx or 5xx response returns an error wrapping an ODataError for every verb.
func TestConfigFunctionsReturnODataErrors(t *testing.T) {
	ctx := context.Background()

	calls := map[string]func(adapter abstractions.RequestAdapter) error{
		"PostRequest": func(adapter abstractions.RequestAdapter) error {
			_, err := PostRequest(ctx, adapter, PostRequestConfig{APIVersion: GraphAPIBeta, Endpoint: "deviceManagement/deviceCategories", RequestBody: newTestBody()}, models.CreateDeviceCategoryFromDiscriminatorValue, nil)
			return err
		},
		"PostRequestNoContent": func(adapter abstractions.RequestAdapter) error {
			return PostRequestNoContent(ctx, adapter, PostRequestConfig{APIVersion: GraphAPIBeta, Endpoint: "deviceManagement/deviceCategories", RequestBody: newTestBody()})
		},
		"GetRequest": func(adapter abstractions.RequestAdapter) error {
			_, err := GetRequest(ctx, adapter, GetRequestConfig{APIVersion: GraphAPIBeta, Endpoint: "deviceManagement/deviceCategories", ResourceID: testResourceID}, models.CreateDeviceCategoryFromDiscriminatorValue, nil)
			return err
		},
		"GetRequestByResourceId": func(adapter abstractions.RequestAdapter) error {
			_, err := GetRequestByResourceId(ctx, adapter, GetRequestConfig{APIVersion: GraphAPIBeta, Endpoint: "/deviceManagement/deviceCategories", ResourceIDPattern: "('id')", ResourceID: testResourceID})
			return err
		},
		"PatchRequestByResourceId": func(adapter abstractions.RequestAdapter) error {
			return PatchRequestByResourceId(ctx, adapter, PatchRequestConfig{APIVersion: GraphAPIBeta, Endpoint: "deviceManagement/deviceCategories", ResourceIDPattern: "/{id}", ResourceID: testResourceID, RequestBody: newTestBody()})
		},
		"PutRequestByResourceId": func(adapter abstractions.RequestAdapter) error {
			return PutRequestByResourceId(ctx, adapter, PutRequestConfig{APIVersion: GraphAPIBeta, Endpoint: "/deviceManagement/deviceCategories", ResourceID: testResourceID, RequestBody: newTestBody()})
		},
		"DeleteRequestByResourceId": func(adapter abstractions.RequestAdapter) error {
			return DeleteRequestByResourceId(ctx, adapter, DeleteRequestConfig{APIVersion: GraphAPIBeta, Endpoint: "deviceManagement/deviceCategories", ResourceIDPattern: "('id')", ResourceID: testResourceID})
		},
	}

	statuses := []struct {
		status int
		code   string
	}{
		{status: http.StatusBadRequest, code: "BadRequest"},
		{status: http.StatusNotFound, code: "ResourceNotFound"},
		{status: http.StatusInternalServerError, code: "InternalServerError"},
	}

	for name, call := range calls {
		for _, st := range statuses {
			t.Run(name+"/"+st.code, func(t *testing.T) {
				adapter, transport := newTestAdapter(t, "https://graph.microsoft.com/beta")
				transport.RegisterNoResponder(odataErrorResponder(st.status, st.code))

				requireODataError(t, call(adapter), st.status, st.code)
				assert.Equal(t, 1, transport.GetTotalCallCount())
			})
		}
	}
}

// TestGetRequestByResourceIdPaging proves paged responses are still followed through @odata.nextLink and combined into one value array.
func TestGetRequestByResourceIdPaging(t *testing.T) {
	ctx := context.Background()
	adapter, transport := newTestAdapter(t, usGovBetaRoot)

	firstURL := usGovBetaRoot + "/applications/" + testResourceID + "/microsoft.graph.agentIdentityBlueprint/sponsors"
	nextLink := firstURL + "?$skiptoken=page%272%27"
	transport.RegisterResponder(http.MethodGet, firstURL, jsonResponder(http.StatusOK, `{"value":[{"id":"a"},{"id":"b"}],"@odata.nextLink":"`+nextLink+`"}`))
	transport.RegisterResponderWithQuery(http.MethodGet, firstURL, "$skiptoken=page%272%27", jsonResponder(http.StatusOK, `{"value":[{"id":"c"}]}`))

	body, err := GetRequestByResourceId(ctx, adapter, GetRequestConfig{APIVersion: GraphAPIBeta, Endpoint: "/applications/" + testResourceID + "/microsoft.graph.agentIdentityBlueprint/sponsors"})
	require.NoError(t, err)

	var combined struct {
		Value []struct {
			ID string `json:"id"`
		} `json:"value"`
	}
	require.NoError(t, json.Unmarshal(body, &combined))
	require.Len(t, combined.Value, 3)
	assert.Equal(t, "a", combined.Value[0].ID)
	assert.Equal(t, "c", combined.Value[2].ID)
	assert.Equal(t, 2, transport.GetTotalCallCount())
}

func TestGetRequestByResourceIdPagingError(t *testing.T) {
	ctx := context.Background()
	adapter, transport := newTestAdapter(t, "https://graph.microsoft.com/beta")

	firstURL := "https://graph.microsoft.com/beta/deviceManagement/deviceCategories"
	transport.RegisterResponder(http.MethodGet, firstURL, jsonResponder(http.StatusOK, `{"value":[{"id":"a"}],"@odata.nextLink":"`+firstURL+`?$skiptoken=2"}`))
	transport.RegisterResponderWithQuery(http.MethodGet, firstURL, "$skiptoken=2", odataErrorResponder(http.StatusNotFound, "ResourceNotFound"))

	_, err := GetRequestByResourceId(ctx, adapter, GetRequestConfig{APIVersion: GraphAPIBeta, Endpoint: "/deviceManagement/deviceCategories"})
	requireODataError(t, err, http.StatusNotFound, "ResourceNotFound")
}

func TestGetRequestByResourceIdNoContent(t *testing.T) {
	ctx := context.Background()
	adapter, transport := newTestAdapter(t, "https://graph.microsoft.com/beta")
	transport.RegisterNoResponder(httpmock.NewStringResponder(http.StatusNoContent, ""))

	body, err := GetRequestByResourceId(ctx, adapter, GetRequestConfig{APIVersion: GraphAPIBeta, Endpoint: "/deviceManagement/deviceCategories", ResourceIDPattern: "('id')", ResourceID: testResourceID})
	require.NoError(t, err)
	assert.Empty(t, body)
}

// TestByIDPatternsBuildRequestURLs proves both "/{id}" and "('id')" patterns reach the expected URL, including an endpoint suffix.
func TestByIDPatternsBuildRequestURLs(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name    string
		pattern string
		suffix  string
		url     string
	}{
		{name: "slash pattern", pattern: "/{id}", url: "https://graph.microsoft.com/beta/applications/" + testResourceID},
		{name: "slash pattern with suffix", pattern: "/{id}", suffix: "/$ref", url: "https://graph.microsoft.com/beta/applications/" + testResourceID + "/$ref"},
		{name: "quoted pattern", pattern: "('id')", url: "https://graph.microsoft.com/beta/applications('" + testResourceID + "')"},
		{name: "quoted pattern with suffix", pattern: "('id')", suffix: "/owners", url: "https://graph.microsoft.com/beta/applications('" + testResourceID + "')/owners"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			adapter, transport := newTestAdapter(t, "https://graph.microsoft.com/v1.0")
			var requested []string
			transport.RegisterNoResponder(func(req *http.Request) (*http.Response, error) {
				requested = append(requested, req.Method+" "+req.URL.String())
				if req.Body != nil {
					_, _ = io.Copy(io.Discard, req.Body)
				}
				return httpmock.NewStringResponse(http.StatusNoContent, ""), nil
			})

			require.NoError(t, DeleteRequestByResourceId(ctx, adapter, DeleteRequestConfig{APIVersion: GraphAPIBeta, Endpoint: "/applications", ResourceIDPattern: tt.pattern, ResourceID: testResourceID, EndpointSuffix: tt.suffix}))
			require.NoError(t, PatchRequestByResourceId(ctx, adapter, PatchRequestConfig{APIVersion: GraphAPIBeta, Endpoint: "applications", ResourceIDPattern: tt.pattern, ResourceID: testResourceID, EndpointSuffix: tt.suffix, RequestBody: newTestBody()}))
			_, err := GetRequestByResourceId(ctx, adapter, GetRequestConfig{APIVersion: GraphAPIBeta, Endpoint: "/applications", ResourceIDPattern: tt.pattern, ResourceID: testResourceID, EndpointSuffix: tt.suffix})
			require.NoError(t, err)

			assert.Equal(t, []string{http.MethodDelete + " " + tt.url, http.MethodPatch + " " + tt.url, http.MethodGet + " " + tt.url}, requested)
		})
	}
}
