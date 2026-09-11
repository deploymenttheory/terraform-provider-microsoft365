package client

import (
	"bytes"
	"compress/gzip"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	abstractions "github.com/microsoft/kiota-abstractions-go"
	"github.com/microsoft/kiota-abstractions-go/authentication"
	khttp "github.com/microsoft/kiota-http-go"
	msgraphsdk "github.com/microsoftgraph/msgraph-sdk-go"
	"github.com/stretchr/testify/require"
)

func TestUnit_EnsureCompressionPrecedesRetry_Permutations(t *testing.T) {
	retry := khttp.NewRetryHandler()
	compression := khttp.NewCompressionHandler()
	redirect := khttp.NewRedirectHandler()
	for _, handlers := range [][]khttp.Middleware{
		{retry, compression, redirect}, {retry, redirect, compression},
		{compression, retry, redirect}, {compression, redirect, retry},
		{redirect, retry, compression}, {redirect, compression, retry},
	} {
		before := append([]khttp.Middleware(nil), handlers...)
		result := ensureCompressionPrecedesRetry(handlers)
		require.Len(t, result, len(before))
		require.ElementsMatch(t, before, result)
		compressionIndex, retryIndex := -1, -1
		var others []khttp.Middleware
		for i, handler := range result {
			if handler == compression {
				compressionIndex = i
			} else {
				others = append(others, handler)
			}
			if handler == retry {
				retryIndex = i
			}
		}
		require.Less(t, compressionIndex, retryIndex)
		var expectedOthers []khttp.Middleware
		for _, handler := range before {
			if handler != compression {
				expectedOthers = append(expectedOthers, handler)
			}
		}
		require.Equal(t, expectedOthers, others)
		require.Equal(t, result, ensureCompressionPrecedesRetry(result))
	}
	withoutRetry := []khttp.Middleware{redirect, compression}
	require.Equal(t, withoutRetry, ensureCompressionPrecedesRetry(withoutRetry))
}

func TestUnit_AddCompressionHandler_ReplacesExisting(t *testing.T) {
	result := addCompressionHandler(context.Background(),
		[]khttp.Middleware{khttp.NewRetryHandler(), khttp.NewCompressionHandler()},
		&ClientOptions{EnableCompression: true})
	require.Len(t, result, 2)
	require.IsType(t, &khttp.CompressionHandler{}, result[1])
}

// Use the same Graph adapter as NewGraphClients: a plain HTTP client does not
// reproduce the request cloning performed by the SDK's tracing middleware.
func TestUnit_ConfigureGraphClientOptions_CompressedRetry(t *testing.T) {
	for _, method := range []abstractions.HttpMethod{abstractions.POST, abstractions.PATCH} {
		for _, scenario := range []struct {
			name     string
			custom   bool
			disable  bool
			fallback bool
		}{
			{name: "disabled options"},
			{name: "custom options", custom: true},
			{name: "request compression disabled", custom: true, disable: true},
			{name: "415 then 429", custom: true, fallback: true},
		} {
			t.Run(method.String()+"/"+scenario.name, func(t *testing.T) {
				t.Setenv("TF_ACC", "")
				const payload = `{"name":"retry probe","description":"preserve the complete payload"}`
				type receivedRequest struct {
					body     []byte
					encoding string
					err      error
				}
				var received []receivedRequest
				server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					body, err := io.ReadAll(r.Body)
					received = append(received, receivedRequest{body, r.Header.Get("Content-Encoding"), err})
					w.Header().Set("Retry-After", "0")
					if scenario.fallback && len(received) == 1 {
						w.WriteHeader(http.StatusUnsupportedMediaType)
					} else if len(received) == 1 || (scenario.fallback && len(received) == 2) {
						w.WriteHeader(http.StatusTooManyRequests)
					} else {
						w.WriteHeader(http.StatusNoContent)
					}
				}))
				defer server.Close()
				options := &ClientOptions{}
				if scenario.custom {
					options = &ClientOptions{EnableRetry: true, MaxRetries: 1, RetryDelaySeconds: 1, EnableCompression: true}
				}
				httpClient, err := ConfigureGraphClientOptions(context.Background(), &ProviderData{ClientOptions: options})
				require.NoError(t, err)
				adapter, err := msgraphsdk.NewGraphRequestAdapterWithParseNodeFactoryAndSerializationWriterFactoryAndHttpClient(
					&authentication.AnonymousAuthenticationProvider{}, nil, nil, httpClient)
				require.NoError(t, err)
				request := abstractions.NewRequestInformation()
				request.UrlTemplate = server.URL
				request.Method = method
				request.Content = []byte(payload)
				request.Headers.Add("Content-Type", "application/json")
				if scenario.disable {
					request.AddRequestOptions([]abstractions.RequestOption{khttp.NewCompressionOptionsReference(false)})
				}
				err = adapter.SendNoContent(context.Background(), request, nil)
				if scenario.custom {
					require.NoError(t, err)
				} else {
					require.Error(t, err, "disabled retry must return the first 429")
				}
				wantAttempts := 2
				if !scenario.custom {
					wantAttempts = 1
				}
				if scenario.fallback {
					wantAttempts = 3
				}
				require.Len(t, received, wantAttempts)
				for i, got := range received {
					require.NoError(t, got.err)
					body := got.body
					if !scenario.custom || scenario.disable || (scenario.fallback && i > 0) {
						require.Empty(t, got.encoding)
					} else {
						require.Equal(t, "gzip", got.encoding)
						reader, err := gzip.NewReader(bytes.NewReader(body))
						require.NoError(t, err, "attempt %d must contain valid gzip", i+1)
						body, err = io.ReadAll(reader)
						require.NoError(t, err)
						require.NoError(t, reader.Close())
					}
					require.Equal(t, payload, string(body), "attempt %d must preserve the payload", i+1)
				}
			})
		}
	}
}
