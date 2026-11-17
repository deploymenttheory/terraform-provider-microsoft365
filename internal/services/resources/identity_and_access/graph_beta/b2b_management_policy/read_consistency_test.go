package graphBetaIdentityAndAccessB2bManagementPolicy

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"testing/synctest"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/require"

	providerclient "github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
)

func TestUnitResourceB2bManagementPolicy_25_ExistsAfterTransientNotFound(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		transport := httpmock.NewMockTransport()
		calls := 0
		statuses := []int{http.StatusOK, http.StatusNotFound, http.StatusOK}
		transport.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/policies/b2bManagementPolicies/`, func(_ *http.Request) (*http.Response, error) {
			status := statuses[min(calls, len(statuses)-1)]
			calls++
			if status == http.StatusNotFound {
				return httpmock.NewJsonResponse(status, map[string]any{"error": map[string]string{
					"code": "Directory_ObjectNotFound", "message": "Unable to read the company information from the directory.",
				}})
			}
			return httpmock.NewJsonResponse(status, map[string]any{
				"id": "11111111-1111-1111-1111-111111111111", "displayName": "Inactive B2B policy", "isOrganizationDefault": false,
			})
		})
		client := providerclient.NewMockGraphClients(&http.Client{Transport: transport}).MockBetaClient
		state := &terraform.InstanceState{ID: "11111111-1111-1111-1111-111111111111"}
		require.NoError(t, b2bManagementPolicyExists(client, context.Background(), state))
		require.NoError(t, b2bManagementPolicyExists(client, context.Background(), state), "a replica 404 must not report a healthy policy as absent")
		require.Equal(t, 3, calls)
	})
}

func TestUnitResourceB2bManagementPolicy_27_DeleteOutcomes(t *testing.T) {
	for _, test := range []struct {
		name         string
		deleteStatus int
		getStatus    int
		deadline     time.Duration
		wantStatus   int
	}{
		{name: "confirmed absent", deleteStatus: http.StatusNotFound, getStatus: http.StatusNotFound},
		{name: "delete forbidden", deleteStatus: http.StatusForbidden, wantStatus: http.StatusForbidden},
		{name: "confirmation forbidden", deleteStatus: http.StatusNotFound, getStatus: http.StatusForbidden, wantStatus: http.StatusForbidden},
		{name: "confirmation cancelled", deleteStatus: http.StatusNotFound, getStatus: http.StatusNotFound, deadline: time.Second},
		{name: "retry cancelled", deleteStatus: http.StatusNotFound, getStatus: http.StatusOK, deadline: time.Second},
	} {
		t.Run(test.name, func(t *testing.T) {
			synctest.Test(t, func(t *testing.T) {
				transport := httpmock.NewMockTransport()
				deleteCalls, getCalls := 0, 0
				const endpoint = "https://graph.microsoft.com/beta/policies/b2bManagementPolicies/11111111-1111-1111-1111-111111111111"
				respond := func(status int) (*http.Response, error) {
					if status == http.StatusOK {
						return httpmock.NewJsonResponse(status, map[string]any{"id": "11111111-1111-1111-1111-111111111111"})
					}
					return httpmock.NewJsonResponse(status, map[string]any{"error": map[string]string{
						"code": http.StatusText(status), "message": http.StatusText(status),
					}})
				}
				transport.RegisterResponder("DELETE", endpoint, func(_ *http.Request) (*http.Response, error) {
					deleteCalls++
					return respond(test.deleteStatus)
				})
				transport.RegisterResponder("GET", endpoint, func(_ *http.Request) (*http.Response, error) {
					getCalls++
					return respond(test.getStatus)
				})
				client := providerclient.NewMockGraphClients(&http.Client{Transport: transport}).MockBetaClient
				reader := B2bManagementPolicyResource{client: client}
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
				defer cancel()
				if test.deadline > 0 {
					var stop context.CancelFunc
					ctx, stop = context.WithTimeout(ctx, test.deadline)
					defer stop()
				}
				model := B2bManagementPolicyResourceModel{
					ID: types.StringValue("11111111-1111-1111-1111-111111111111"), DisplayName: types.StringValue("Inactive B2B policy"),
					Definition:            types.ListValueMust(types.StringType, []attr.Value{types.StringValue(`{"B2BManagementPolicy":{}}`)}),
					IsOrganizationDefault: types.BoolValue(false),
					Timeouts: timeouts.Value{Object: types.ObjectNull(map[string]attr.Type{
						"create": types.StringType, "read": types.StringType, "update": types.StringType, "delete": types.StringType,
					})},
				}
				var schema resource.SchemaResponse
				reader.Schema(ctx, resource.SchemaRequest{}, &schema)
				state := tfsdk.State{Schema: schema.Schema}
				require.False(t, state.Set(ctx, model).HasError())
				response := resource.DeleteResponse{State: state}
				reader.Delete(ctx, resource.DeleteRequest{State: state}, &response)
				switch {
				case test.deadline > 0:
					require.True(t, response.Diagnostics.HasError())
					require.Contains(t, fmt.Sprint(response.Diagnostics.Errors()), context.DeadlineExceeded.Error())
				case test.wantStatus > 0:
					require.True(t, response.Diagnostics.HasError())
				default:
					require.False(t, response.Diagnostics.HasError(), "%v", response.Diagnostics)
					require.True(t, response.State.Raw.IsNull())
					require.Greater(t, getCalls, 1, "one GET 404 is insufficient to confirm absence")
				}
				if test.wantStatus > 0 || test.deadline > 0 {
					require.True(t, response.State.Raw.Equal(state.Raw), "failed deletion must preserve the complete resource state")
				}
				require.Equal(t, 1, deleteCalls, "authorization failures and cancellation must not trigger another DELETE")
				if test.deleteStatus == http.StatusForbidden {
					require.Zero(t, getCalls)
				}
			})
		})
	}
}
