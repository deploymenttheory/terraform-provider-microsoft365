package graphBetaLinuxDeviceCompliancePolicy

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/microsoft/kiota-abstractions-go/authentication"
	msgraph "github.com/microsoftgraph/msgraph-beta-sdk-go"
	"github.com/stretchr/testify/require"
)

func TestUnitLinuxComplianceDeleteRetriesTransientErrors(t *testing.T) {
	for _, test := range []struct {
		name          string
		firstStatus   int
		expectedCalls int
		expectError   bool
	}{
		{name: "transient_internal_error", firstStatus: 500, expectedCalls: 2},
		{name: "forbidden", firstStatus: 403, expectedCalls: 1, expectError: true},
		{name: "bad_request", firstStatus: 400, expectedCalls: 1, expectError: true},
	} {
		t.Run(test.name, func(t *testing.T) {
			ctx := context.Background()
			var calls atomic.Int32
			server := httptest.NewServer(
				http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
					call := calls.Add(1)
					if req.Method != http.MethodDelete ||
						req.URL.Path != "/beta/deviceManagement/compliancePolicies/policy-id" {
						t.Errorf("unexpected request: %s %s", req.Method, req.URL)
						w.WriteHeader(http.StatusBadRequest)
						return
					}
					if call == 1 {
						w.Header().Set("Content-Type", "application/json")
						w.WriteHeader(test.firstStatus)
						fmt.Fprint(
							w,
							`{"error":{"code":"RequestFailed","message":"delete failed"}}`,
						)
						return
					}
					w.WriteHeader(http.StatusNoContent)
				}),
			)
			t.Cleanup(server.Close)

			adapter, err := msgraph.NewGraphRequestAdapter(
				&authentication.AnonymousAuthenticationProvider{},
			)
			require.NoError(t, err)
			adapter.SetBaseUrl(server.URL + "/beta")
			r := &LinuxDeviceCompliancePolicyResource{
				client: msgraph.NewGraphServiceClient(adapter),
			}
			var schema resource.SchemaResponse
			r.Schema(ctx, resource.SchemaRequest{}, &schema)
			stateType := schema.Schema.Type().TerraformType(ctx).(tftypes.Object)
			attributes := make(map[string]tftypes.Value, len(stateType.AttributeTypes))
			for name, attributeType := range stateType.AttributeTypes {
				attributes[name] = tftypes.NewValue(attributeType, nil)
			}
			attributes["id"] = tftypes.NewValue(tftypes.String, "policy-id")
			attributes["name"] = tftypes.NewValue(tftypes.String, "Policy to delete")
			state := tfsdk.State{
				Schema: schema.Schema,
				Raw:    tftypes.NewValue(stateType, attributes),
			}
			response := resource.DeleteResponse{State: state}

			r.Delete(ctx, resource.DeleteRequest{State: state}, &response)

			require.Equal(t, int32(test.expectedCalls), calls.Load())
			require.Equal(
				t,
				test.expectError,
				response.Diagnostics.HasError(),
				"%v",
				response.Diagnostics,
			)
			if test.expectError {
				require.True(
					t,
					state.Raw.Equal(response.State.Raw),
					"failed deletion must preserve state",
				)
			} else {
				require.True(
					t,
					response.State.Raw.IsNull(),
					"successful deletion must remove state",
				)
			}
		})
	}
}
