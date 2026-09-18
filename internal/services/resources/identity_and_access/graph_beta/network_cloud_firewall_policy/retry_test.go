package graphBetaNetworkCloudFirewallPolicy_test

import (
	"compress/gzip"
	"context"
	kiotahttp "github.com/microsoft/kiota-http-go"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	cloud "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/identity_and_access/graph_beta/network_cloud_firewall_policy"
	rule "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/identity_and_access/graph_beta/network_cloud_firewall_policy_rule"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/microsoft/kiota-abstractions-go/authentication"
	graph "github.com/microsoftgraph/msgraph-beta-sdk-go"
	"github.com/stretchr/testify/require"
)

func TestUnitResourceNetworkCloudFirewall_PostRetries(t *testing.T) {
	for _, child := range []bool{false, true} {
		name := "policy"
		if child {
			name = "rule"
		}
		t.Run(name, func(t *testing.T) {
			for _, tc := range []struct {
				name                       string
				status, failures, attempts int
				maxRetries, delay          int64
			}{
				{"throttled_then_created", 429, 1, 2, 0, 0},
				{"throttling_exhausted", 429, 100, 4, 0, 0},
				{"ambiguous_server_failure", 503, 1, 1, 0, 0},
				{"retry_disabled", 429, 100, 1, 0, 0},
				{"configured_limit", 429, 100, 2, 1, 9},
				{"configured_delay", 429, 1, 2, 1, 5},
				{"configured_no_server_retry", 503, 1, 1, 1, 9},
			} {
				t.Run(tc.name, func(t *testing.T) {
					t.Setenv("TF_ACC", "")
					ctx := context.Background()
					attempts, reads := 0, 0
					var bodies []string
					body := `{"id":"created","name":"probe","settings":{"defaultAction":"allow"},"version":"1.0.0","lastModifiedDateTime":"2026-09-07T00:00:00Z"}`
					if child {
						body = `{"@odata.type":"#microsoft.graph.networkaccess.cloudFirewallRule","id":"created","name":"probe","priority":100,"action":"block","settings":{"status":"disabled"},"matchingConditions":{"sources":null,"destinations":null}}`
					}
					server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
						w.Header().Set("Content-Type", "application/json")
						if req.Method == http.MethodPost {
							attempts++
							var reader io.Reader = req.Body
							if req.Header.Get("Content-Encoding") == "gzip" {
								decoded, err := gzip.NewReader(req.Body)
								if err != nil {
									t.Error(err)
									w.WriteHeader(400)
									return
								}
								defer decoded.Close()
								reader = decoded
							}
							data, err := io.ReadAll(reader)
							if err != nil {
								t.Error(err)
								w.WriteHeader(500)
								return
							}
							bodies = append(bodies, string(data))
							if attempts <= tc.failures {
								if tc.name != "configured_delay" {
									w.Header().Set("Retry-After", "0")
								}
								w.WriteHeader(tc.status)
								_, _ = io.WriteString(w, `{"error":{"code":"TestFailure","message":"retry regression"}}`)
								return
							}
							w.WriteHeader(http.StatusCreated)
						} else {
							reads++
						}
						_, _ = io.WriteString(w, body)
					}))
					defer server.Close()
					httpClient, err := client.ConfigureGraphClientOptions(ctx, &client.ProviderData{ClientOptions: &client.ClientOptions{EnableRetry: tc.name != "retry_disabled", EnableCompression: true, MaxRetries: tc.maxRetries, RetryDelaySeconds: tc.delay}})
					require.NoError(t, err)
					adapter, err := graph.NewGraphRequestAdapterWithParseNodeFactoryAndSerializationWriterFactoryAndHttpClient(&authentication.AnonymousAuthenticationProvider{}, nil, nil, httpClient)
					require.NoError(t, err)
					adapter.SetBaseUrl(server.URL + "/beta")
					var r resource.Resource = cloud.NewNetworkCloudFirewallPolicyResource()
					if child {
						r = rule.NewNetworkCloudFirewallPolicyRuleResource()
					}
					configure := resource.ConfigureResponse{}
					r.(resource.ResourceWithConfigure).Configure(ctx, resource.ConfigureRequest{ProviderData: &client.GraphClients{KiotaGraphBetaClient: graph.NewGraphServiceClient(adapter), RetryOptions: kiotahttp.RetryHandlerOptions{MaxRetries: int(tc.maxRetries), DelaySeconds: int(tc.delay)}}}, &configure)
					require.False(t, configure.Diagnostics.HasError(), configure.Diagnostics)
					schema := resource.SchemaResponse{}
					r.Schema(ctx, resource.SchemaRequest{}, &schema)
					plan := tfsdk.Plan{Schema: schema.Schema, Raw: tftypes.NewValue(schema.Schema.Type().TerraformType(ctx), nil)}
					fields := map[string]any{"id": types.StringUnknown(), "name": "probe"}
					if child {
						fields["policy_id"] = "00000000-0000-0000-0000-000000000001"
						fields["priority"] = int32(100)
						fields["action"] = "block"
						fields["enabled"] = false
						fields["status"] = types.StringUnknown()
					} else {
						fields["default_action"] = "allow"
						fields["version"] = types.StringUnknown()
						fields["last_modified_date_time"] = types.StringUnknown()
					}
					for key, value := range fields {
						require.False(t, plan.SetAttribute(ctx, path.Root(key), value).HasError())
					}
					response := resource.CreateResponse{State: tfsdk.State{Schema: schema.Schema}}
					started := time.Now()
					r.Create(ctx, resource.CreateRequest{Plan: plan}, &response)
					if tc.name == "configured_delay" {
						require.GreaterOrEqual(t, time.Since(started), 5*time.Second)
					}
					require.Equal(t, tc.attempts, attempts)
					for _, sent := range bodies {
						require.Equal(t, bodies[0], sent, "retry must replay the same POST body")
					}
					success := tc.failures < tc.attempts
					require.Equal(t, !success, response.Diagnostics.HasError(), response.Diagnostics)
					if success {
						require.Equal(t, 1, reads)
						var id string
						require.False(t, response.State.GetAttribute(ctx, path.Root("id"), &id).HasError())
						require.Equal(t, "created", id)
					} else {
						require.Zero(t, reads)
					}
				})
			}
		})
	}
}
