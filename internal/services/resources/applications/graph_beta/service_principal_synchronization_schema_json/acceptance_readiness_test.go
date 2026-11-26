package graphBetaServicePrincipalSynchronizationSchemaJson_test

import (
	"context"
	"net/http"
	"sync/atomic"
	"testing"
	"time"

	providerclient "github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/require"
)

func TestUnitSynchronizationSchemaFixtureReadiness(t *testing.T) {
	for _, test := range []struct {
		name      string
		statuses  []int
		timeout   time.Duration
		wantCalls int32
		wantError bool
	}{
		{name: "fresh_then_unauthorized_resets_readiness", statuses: []int{200, 200, 401, 200, 200, 200}, timeout: time.Minute, wantCalls: 6},
		{name: "deadline_before_stable_reads", statuses: []int{200}, timeout: 100 * time.Millisecond, wantCalls: 1, wantError: true},
		{name: "non_retryable_error_after_fresh_read", statuses: []int{200, 400}, timeout: time.Minute, wantCalls: 2, wantError: true},
	} {
		t.Run(test.name, func(t *testing.T) {
			var calls atomic.Int32
			transport := httpmock.NewMockTransport()
			transport.RegisterNoResponder(func(req *http.Request) (*http.Response, error) {
				require.Equal(t, http.MethodGet, req.Method)
				index := int(calls.Add(1)) - 1
				if index >= len(test.statuses) {
					index = len(test.statuses) - 1
				}
				status := test.statuses[index]
				body := `{"synchronizationRules":[{"id":"rule"}],"directories":[{"id":"directory"}],"version":"1"}`
				if status != http.StatusOK {
					body = `{"error":{"code":"UnknownError","message":"An unknown error has occurred."}}`
				}
				response := httpmock.NewStringResponse(status, body)
				response.Header.Set("Content-Type", "application/json")
				return response, nil
			})
			client := providerclient.NewMockGraphClients(&http.Client{Transport: transport}).MockBetaClient
			ctx, cancel := context.WithTimeout(context.Background(), test.timeout)
			defer cancel()
			schema, err := waitForRawSchema(ctx, client, "11111111-1111-1111-1111-111111111111", "scim.test")
			if test.wantError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.NotEmpty(t, schema["synchronizationRules"])
			}
			require.Equal(t, test.wantCalls, calls.Load(), "fixture must not accept an isolated successful replica")
		})
	}
}
