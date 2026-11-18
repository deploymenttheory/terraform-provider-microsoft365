package graphBetaIdentityAndAccessConnectedOrganization

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
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/require"

	providerclient "github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
)

func deleteTestState(t *testing.T, r *ConnectedOrganizationResource) tfsdk.State {
	t.Helper()
	ctx := context.Background()
	var schema resource.SchemaResponse
	r.Schema(ctx, resource.SchemaRequest{}, &schema)
	model := minimalModel()
	model.ID = types.StringValue("11111111-1111-1111-1111-111111111111")
	model.Timeouts = timeouts.Value{Object: types.ObjectNull(map[string]attr.Type{
		"create": types.StringType, "read": types.StringType, "update": types.StringType, "delete": types.StringType,
	})}
	state := tfsdk.State{Schema: schema.Schema}
	require.False(t, state.Set(ctx, model).HasError())
	return state
}

func TestUnitResourceConnectedOrganization_12_DeleteOutcomes(t *testing.T) {
	for _, test := range []struct {
		name      string
		statuses  []int
		wantCalls int
		wantError bool
		deadline  time.Duration
	}{
		{name: "transient service unavailable", statuses: []int{503, 503, 204}, wantCalls: 3},
		{name: "already absent", statuses: []int{404}, wantCalls: 1},
		{name: "permission denied", statuses: []int{403}, wantCalls: 1, wantError: true},
		{name: "invalid request", statuses: []int{400}, wantCalls: 1, wantError: true},
		{name: "deadline", statuses: []int{503}, wantCalls: 1, wantError: true, deadline: time.Second},
		{name: "retry limit", statuses: []int{503}, wantCalls: 11, wantError: true},
	} {
		t.Run(test.name, func(t *testing.T) {
			synctest.Test(t, func(t *testing.T) {
				transport := httpmock.NewMockTransport()
				calls := 0
				transport.RegisterResponder(
					"DELETE",
					"https://graph.microsoft.com/beta/identityGovernance/entitlementManagement/connectedOrganizations/11111111-1111-1111-1111-111111111111",
					func(req *http.Request) (*http.Response, error) {
						require.Equal(t, "/beta/identityGovernance/entitlementManagement/connectedOrganizations/11111111-1111-1111-1111-111111111111", req.URL.Path)
						status := test.statuses[min(calls, len(test.statuses)-1)]
						calls++
						if status == http.StatusNoContent {
							return httpmock.NewStringResponse(status, ""), nil
						}
						return httpmock.NewJsonResponse(status, map[string]any{"error": map[string]string{
							"code": http.StatusText(status), "message": http.StatusText(status),
						}})
					},
				)
				r := NewConnectedOrganizationResource().(*ConnectedOrganizationResource)
				r.client = providerclient.NewMockGraphClients(&http.Client{Transport: transport}).MockBetaClient
				state := deleteTestState(t, r)
				ctx := context.Background()
				if test.deadline > 0 {
					var cancel context.CancelFunc
					ctx, cancel = context.WithTimeout(ctx, test.deadline)
					defer cancel()
				}
				response := resource.DeleteResponse{State: state}
				r.Delete(ctx, resource.DeleteRequest{State: state}, &response)
				require.Equal(t, test.wantError, response.Diagnostics.HasError(), "%v", response.Diagnostics)
				require.Equal(t, test.wantCalls, calls)
				if test.wantError {
					require.True(t, response.State.Raw.Equal(state.Raw), "an unconfirmed deletion must retain the complete resource state")
				} else {
					require.True(t, response.State.Raw.IsNull(), "a confirmed deletion must remove the resource from state")
				}
				if test.deadline > 0 {
					require.Contains(t, fmt.Sprint(response.Diagnostics.Errors()), context.DeadlineExceeded.Error())
				}
			})
		})
	}
}
