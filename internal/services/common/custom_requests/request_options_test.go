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
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestPostRequestOptionsReachMiddleware checks that PostRequestConfig.RequestOptions reach the adapter's middleware:
// a retry option that retries nothing stops the retry middleware from retrying a 503 that it otherwise retries.
func TestPostRequestOptionsReachMiddleware(t *testing.T) {
	noRetry := kiotahttp.RetryHandlerOptions{MaxRetries: 3, ShouldRetry: func(time.Duration, int, *http.Request, *http.Response) bool { return false }}

	tests := []struct {
		name          string
		options       []abstractions.RequestOption
		noContent     bool
		expectedCalls int
	}{
		{name: "PostRequest without options retries", expectedCalls: 4},
		{name: "PostRequest with options", options: []abstractions.RequestOption{&noRetry}, expectedCalls: 1},
		{name: "PostRequestNoContent without options retries", noContent: true, expectedCalls: 4},
		{name: "PostRequestNoContent with options", options: []abstractions.RequestOption{&noRetry}, noContent: true, expectedCalls: 1},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockTransport := httpmock.NewMockTransport()
			retryHandler := kiotahttp.NewRetryHandlerWithOptions(kiotahttp.RetryHandlerOptions{MaxRetries: 3, ShouldRetry: func(time.Duration, int, *http.Request, *http.Response) bool { return true }})
			httpClient := &http.Client{Transport: kiotahttp.NewCustomTransportWithParentTransport(mockTransport, retryHandler)}
			adapter, err := msgraphbetasdk.NewGraphRequestAdapterWithParseNodeFactoryAndSerializationWriterFactoryAndHttpClient(&authentication.AnonymousAuthenticationProvider{}, nil, nil, httpClient)
			require.NoError(t, err)
			msgraphbetasdk.NewGraphServiceClient(adapter)

			calls := 0
			mockTransport.RegisterResponder(http.MethodPost, "https://graph.microsoft.com/beta/deviceManagement/deviceCategories", func(req *http.Request) (*http.Response, error) {
				calls++
				resp := httpmock.NewStringResponse(http.StatusServiceUnavailable, `{"error":{"code":"ServiceUnavailable","message":"unavailable"}}`)
				resp.Header.Set("Content-Type", "application/json")
				resp.Header.Set("Retry-After", "0")
				return resp, nil
			})

			config := PostRequestConfig{APIVersion: GraphAPIBeta, Endpoint: "deviceManagement/deviceCategories", RawBody: map[string]string{"displayName": "test"}, RequestOptions: tc.options}
			if tc.noContent {
				err = PostRequestNoContent(context.Background(), adapter, config)
			} else {
				var response map[string]any
				config.RawResponse = &response
				_, err = PostRequest(context.Background(), adapter, config, nil, nil)
			}

			requireODataError(t, err, http.StatusServiceUnavailable, "ServiceUnavailable")
			assert.Equal(t, tc.expectedCalls, calls)
		})
	}
}
