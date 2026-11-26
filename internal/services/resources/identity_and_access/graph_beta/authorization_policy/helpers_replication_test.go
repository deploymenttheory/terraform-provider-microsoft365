package graphBetaAuthorizationPolicy

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"testing/synctest"
	"time"

	"github.com/jarcoal/httpmock"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
	"github.com/stretchr/testify/require"

	providerclient "github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
)

func TestUnitResourceAuthorizationPolicy_23_WriteReplication(t *testing.T) {
	for _, test := range []struct {
		name    string
		misses  int
		timeout time.Duration
		calls   int
		failed  bool
	}{
		{"transient", 1, 0, 2, false},
		{"persistent", 100, 0, 16, true},
		{"deadline", 100, time.Second, 1, true},
	} {
		t.Run(test.name, func(t *testing.T) {
			synctest.Test(t, func(t *testing.T) {
				ctx := context.Background()
				if test.timeout > 0 {
					var cancel context.CancelFunc
					ctx, cancel = context.WithTimeout(ctx, test.timeout)
					defer cancel()
				}
				transport := httpmock.NewMockTransport()
				calls := 0
				transport.RegisterResponder("PATCH", "https://graph.microsoft.com/beta/policies/authorizationPolicy/authorizationPolicy", func(request *http.Request) (*http.Response, error) {
					var payload map[string]any
					require.NoError(t, json.NewDecoder(request.Body).Decode(&payload))
					require.Equal(t, false, payload["allowedToUseSSPR"], "every retry must retain the complete write")
					require.Equal(t, false, payload["defaultUserRolePermissions"].(map[string]any)["allowedToCreateApps"])
					calls++
					if calls <= test.misses {
						return httpmock.NewStringResponse(404, `{"error":{"code":"Request_ResourceNotFound","message":"Policy replica unavailable"}}`), nil
					}
					return httpmock.NewStringResponse(204, ""), nil
				})
				r := NewAuthorizationPolicyResource().(*AuthorizationPolicyResource)
				r.client = providerclient.NewMockGraphClients(&http.Client{Transport: transport}).MockBetaClient
				body := graphmodels.NewAuthorizationPolicy()
				value := false
				body.SetAllowedToUseSSPR(&value)
				permissions := graphmodels.NewDefaultUserRolePermissions()
				permissions.SetAllowedToCreateApps(&value)
				body.SetDefaultUserRolePermissions(permissions)
				err := r.patchSettings(ctx, body)
				if test.failed {
					require.Error(t, err)
				} else {
					require.NoError(t, err)
				}
				require.Equal(t, test.calls, calls)
			})
		})
	}
}
