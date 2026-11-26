package graphBetaAuthorizationPolicy

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"testing/synctest"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/require"

	providerclient "github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
)

func TestUnitResourceAuthorizationPolicy_23_WriteReplication(t *testing.T) {
	for _, operation := range []string{"create", "update"} {
		t.Run(operation, func(t *testing.T) {
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
						ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
						defer cancel()
						if test.timeout > 0 {
							var cancel context.CancelFunc
							ctx, cancel = context.WithTimeout(ctx, test.timeout)
							defer cancel()
						}
						transport := httpmock.NewMockTransport()
						calls := 0
						remote := map[string]any{}
						transport.RegisterResponder("GET", "https://graph.microsoft.com/beta/policies/authorizationPolicy", func(request *http.Request) (*http.Response, error) {
							return httpmock.NewJsonResponse(200, map[string]any{"value": []any{remote}})
						})
						transport.RegisterResponder("PATCH", "https://graph.microsoft.com/beta/policies/authorizationPolicy/authorizationPolicy", func(request *http.Request) (*http.Response, error) {
							var payload map[string]any
							require.NoError(t, json.NewDecoder(request.Body).Decode(&payload))
							require.Equal(t, false, payload["allowedToUseSSPR"], "every retry must retain the complete write")
							require.Equal(t, false, payload["defaultUserRolePermissions"].(map[string]any)["allowedToCreateApps"])
							remote = payload
							remote["id"] = singletonID
							calls++
							if calls <= test.misses {
								return httpmock.NewStringResponse(404, `{"error":{"code":"Request_ResourceNotFound","message":"Policy replica unavailable"}}`), nil
							}
							return httpmock.NewStringResponse(204, ""), nil
						})
						r := NewAuthorizationPolicyResource().(*AuthorizationPolicyResource)
						r.client = providerclient.NewMockGraphClients(&http.Client{Transport: transport}).MockBetaClient
						model := testModel()
						model.AllowedToUseSSPR = types.BoolValue(false)
						model.AllowInvitesFrom = types.StringUnknown()
						state := testState(t, r, model)
						var failed bool
						if operation == "create" {
							response := resource.CreateResponse{State: tfsdk.State{Schema: state.Schema}}
							r.Create(ctx, resource.CreateRequest{Plan: tfsdk.Plan{Schema: state.Schema, Raw: state.Raw}}, &response)
							failed = response.Diagnostics.HasError()
							require.Equal(t, test.failed, failed, response.Diagnostics)
						} else {
							response := resource.UpdateResponse{State: state}
							r.Update(ctx, resource.UpdateRequest{Plan: tfsdk.Plan{Schema: state.Schema, Raw: state.Raw}, State: state}, &response)
							failed = response.Diagnostics.HasError()
							require.Equal(t, test.failed, failed, response.Diagnostics)
						}
						require.Equal(t, test.failed, failed)
						require.Equal(t, test.calls, calls)
					})
				})
			}
		})
	}
}
