package customrequests

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
	abstractions "github.com/microsoft/kiota-abstractions-go"
	"github.com/microsoft/kiota-abstractions-go/authentication"
	kiotahttp "github.com/microsoft/kiota-http-go"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const optionsTestCollectionURL = "https://graph.microsoft.com/beta/deviceManagement/deviceCategories"

// newRetryingTestAdapter returns a Graph beta request adapter whose retry middleware retries every response up to
// three times, sending through an isolated httpmock transport.
func newRetryingTestAdapter(t *testing.T) (abstractions.RequestAdapter, *httpmock.MockTransport) {
	t.Helper()

	mockTransport := httpmock.NewMockTransport()
	retryHandler := kiotahttp.NewRetryHandlerWithOptions(kiotahttp.RetryHandlerOptions{MaxRetries: 3, ShouldRetry: func(time.Duration, int, *http.Request, *http.Response) bool { return true }})
	httpClient := &http.Client{Transport: kiotahttp.NewCustomTransportWithParentTransport(mockTransport, retryHandler)}
	adapter, err := msgraphbetasdk.NewGraphRequestAdapterWithParseNodeFactoryAndSerializationWriterFactoryAndHttpClient(&authentication.AnonymousAuthenticationProvider{}, nil, nil, httpClient)
	require.NoError(t, err)
	msgraphbetasdk.NewGraphServiceClient(adapter)

	return adapter, mockTransport
}

// countingUnavailableResponder answers 503 Service Unavailable with Retry-After 0 and counts the calls.
func countingUnavailableResponder(calls *int) httpmock.Responder {
	return func(*http.Request) (*http.Response, error) {
		*calls++
		resp := httpmock.NewStringResponse(http.StatusServiceUnavailable, `{"error":{"code":"ServiceUnavailable","message":"unavailable"}}`)
		resp.Header.Set("Content-Type", "application/json")
		resp.Header.Set("Retry-After", "0")
		return resp, nil
	}
}

// TestRequestOptionsReachMiddleware checks that the RequestOptions of every verb config reach the adapter's
// middleware on every send path: a retry option that retries nothing stops the retry middleware from retrying a 503
// that it otherwise retries. The paged GET cases answer the first page and fail the @odata.nextLink page, so they
// prove the options are sent with every page.
func TestRequestOptionsReachMiddleware(t *testing.T) {
	ctx := context.Background()
	endpoint := "deviceManagement/deviceCategories"
	byIDURL := optionsTestCollectionURL + "/" + testResourceID
	putURL := optionsTestCollectionURL + "('" + testResourceID + "')"
	rawBody := map[string]string{"displayName": "test"}

	verbs := []struct {
		name  string
		verb  string
		url   string
		paged bool
		send  func(adapter abstractions.RequestAdapter, options []abstractions.RequestOption) error
	}{
		{name: "PostRequest", verb: http.MethodPost, url: optionsTestCollectionURL, send: func(adapter abstractions.RequestAdapter, options []abstractions.RequestOption) error {
			var response map[string]any
			_, err := PostRequest(ctx, adapter, PostRequestConfig{APIVersion: GraphAPIBeta, Endpoint: endpoint, RawBody: rawBody, RawResponse: &response, RequestOptions: options}, nil, nil)
			return err
		}},
		{name: "PostRequestNoContent", verb: http.MethodPost, url: optionsTestCollectionURL, send: func(adapter abstractions.RequestAdapter, options []abstractions.RequestOption) error {
			return PostRequestNoContent(ctx, adapter, PostRequestConfig{APIVersion: GraphAPIBeta, Endpoint: endpoint, RawBody: rawBody, RequestOptions: options})
		}},
		{name: "GetRequest parsable", verb: http.MethodGet, url: byIDURL, send: func(adapter abstractions.RequestAdapter, options []abstractions.RequestOption) error {
			_, err := GetRequest(ctx, adapter, GetRequestConfig{APIVersion: GraphAPIBeta, Endpoint: endpoint, ResourceID: testResourceID, RequestOptions: options}, models.CreateDeviceCategoryFromDiscriminatorValue, nil)
			return err
		}},
		{name: "GetRequest raw", verb: http.MethodGet, url: byIDURL, send: func(adapter abstractions.RequestAdapter, options []abstractions.RequestOption) error {
			var response map[string]any
			_, err := GetRequest(ctx, adapter, GetRequestConfig{APIVersion: GraphAPIBeta, Endpoint: endpoint, ResourceID: testResourceID, RawResponse: &response, RequestOptions: options}, nil, nil)
			return err
		}},
		{name: "GetRequest raw paged", verb: http.MethodGet, url: optionsTestCollectionURL, paged: true, send: func(adapter abstractions.RequestAdapter, options []abstractions.RequestOption) error {
			var response map[string]any
			_, err := GetRequest(ctx, adapter, GetRequestConfig{APIVersion: GraphAPIBeta, Endpoint: endpoint, RawResponse: &response, RequestOptions: options}, nil, nil)
			return err
		}},
		{name: "GetRequestByResourceId", verb: http.MethodGet, url: byIDURL, send: func(adapter abstractions.RequestAdapter, options []abstractions.RequestOption) error {
			_, err := GetRequestByResourceId(ctx, adapter, GetRequestConfig{APIVersion: GraphAPIBeta, Endpoint: endpoint, ResourceID: testResourceID, ResourceIDPattern: ResourceIDPatternSegment, RequestOptions: options})
			return err
		}},
		{name: "GetRequestByResourceId paged", verb: http.MethodGet, url: optionsTestCollectionURL, paged: true, send: func(adapter abstractions.RequestAdapter, options []abstractions.RequestOption) error {
			_, err := GetRequestByResourceId(ctx, adapter, GetRequestConfig{APIVersion: GraphAPIBeta, Endpoint: endpoint, ResourceIDPattern: ResourceIDPatternNone, RequestOptions: options})
			return err
		}},
		{name: "PatchRequestByResourceId", verb: http.MethodPatch, url: byIDURL, send: func(adapter abstractions.RequestAdapter, options []abstractions.RequestOption) error {
			return PatchRequestByResourceId(ctx, adapter, PatchRequestConfig{APIVersion: GraphAPIBeta, Endpoint: endpoint, ResourceID: testResourceID, ResourceIDPattern: ResourceIDPatternSegment, RequestBody: newTestBody(), RequestOptions: options})
		}},
		{name: "PatchRequestByResourceId raw", verb: http.MethodPatch, url: byIDURL, send: func(adapter abstractions.RequestAdapter, options []abstractions.RequestOption) error {
			var response map[string]any
			return PatchRequestByResourceId(ctx, adapter, PatchRequestConfig{APIVersion: GraphAPIBeta, Endpoint: endpoint, ResourceID: testResourceID, ResourceIDPattern: ResourceIDPatternSegment, RawBody: rawBody, RawResponse: &response, RequestOptions: options})
		}},
		{name: "PutRequestByResourceId", verb: http.MethodPut, url: putURL, send: func(adapter abstractions.RequestAdapter, options []abstractions.RequestOption) error {
			return PutRequestByResourceId(ctx, adapter, PutRequestConfig{APIVersion: GraphAPIBeta, Endpoint: endpoint, ResourceID: testResourceID, RequestBody: newTestBody(), RequestOptions: options})
		}},
		{name: "DeleteRequestByResourceId", verb: http.MethodDelete, url: byIDURL, send: func(adapter abstractions.RequestAdapter, options []abstractions.RequestOption) error {
			return DeleteRequestByResourceId(ctx, adapter, DeleteRequestConfig{APIVersion: GraphAPIBeta, Endpoint: endpoint, ResourceID: testResourceID, ResourceIDPattern: ResourceIDPatternSegment, RequestOptions: options})
		}},
	}

	noRetry := kiotahttp.RetryHandlerOptions{MaxRetries: 3, ShouldRetry: func(time.Duration, int, *http.Request, *http.Response) bool { return false }}
	optionCases := []struct {
		suffix        string
		options       []abstractions.RequestOption
		expectedCalls int
	}{
		{suffix: " without options retries", expectedCalls: 4},
		{suffix: " with options", options: []abstractions.RequestOption{&noRetry}, expectedCalls: 1},
	}

	for _, v := range verbs {
		for _, oc := range optionCases {
			t.Run(v.name+oc.suffix, func(t *testing.T) {
				adapter, mockTransport := newRetryingTestAdapter(t)

				calls := 0
				if v.paged {
					mockTransport.RegisterResponder(v.verb, v.url, jsonResponder(http.StatusOK, `{"value":[{"id":"a"}],"@odata.nextLink":"`+v.url+`?$skiptoken=2"}`))
					mockTransport.RegisterResponderWithQuery(v.verb, v.url, "$skiptoken=2", countingUnavailableResponder(&calls))
				} else {
					mockTransport.RegisterResponder(v.verb, v.url, countingUnavailableResponder(&calls))
				}

				err := v.send(adapter, oc.options)

				requireODataError(t, err, http.StatusServiceUnavailable, "ServiceUnavailable")
				assert.Equal(t, oc.expectedCalls, calls)
			})
		}
	}
}
