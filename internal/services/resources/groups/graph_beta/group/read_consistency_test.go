package graphBetaGroup

import (
	"context"
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

func TestUnitGroupReadDirectoryConsistency(t *testing.T) {
	for _, tc := range []struct {
		name      string
		statuses  []int
		reads     int
		wantGone  bool
		wantError bool
		timeout   time.Duration
	}{
		{name: "404 after a successful read is transient", statuses: []int{200, 404, 200}, reads: 2},
		{name: "persistent 404 confirms deletion", statuses: []int{404}, reads: 1, wantGone: true},
		{name: "403 preserves state and reports an error", statuses: []int{403}, reads: 1, wantError: true},
		{name: "cancellation preserves the group", statuses: []int{404}, reads: 1, wantError: true, timeout: time.Second},
	} {
		t.Run(tc.name, func(t *testing.T) {
			synctest.Test(t, func(t *testing.T) {
				ctx := context.Background()
				if tc.timeout > 0 {
					var cancel context.CancelFunc
					ctx, cancel = context.WithTimeout(ctx, tc.timeout)
					defer cancel()
				}
				transport := httpmock.NewMockTransport()
				calls := 0
				transport.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/groups/`, func(_ *http.Request) (*http.Response, error) {
					status := tc.statuses[min(calls, len(tc.statuses)-1)]
					calls++
					if status != http.StatusOK {
						code := "Request_ResourceNotFound"
						if status == http.StatusForbidden {
							code = "Authorization_RequestDenied"
						}
						return httpmock.NewJsonResponse(status, map[string]any{"error": map[string]any{
							"code": code, "message": "Directory replica response",
						}})
					}
					return httpmock.NewJsonResponse(200, map[string]any{
						"id": "11111111-1111-1111-1111-111111111111", "displayName": "Directory consistency",
						"mailNickname": "directory-consistency", "mailEnabled": false, "securityEnabled": true,
					})
				})
				r := NewGroupResource().(*GroupResource)
				r.client = providerclient.NewMockGraphClients(&http.Client{Transport: transport}).MockBetaClient
				var schema resource.SchemaResponse
				r.Schema(ctx, resource.SchemaRequest{}, &schema)
				state := tfsdk.State{Schema: schema.Schema}
				model := GroupResourceModel{
					ID:           types.StringValue("11111111-1111-1111-1111-111111111111"),
					GroupTypes:   types.SetNull(types.StringType),
					GroupOwners:  types.SetNull(types.StringType),
					GroupMembers: types.SetNull(types.StringType),
					Timeouts: timeouts.Value{Object: types.ObjectNull(map[string]attr.Type{
						"create": types.StringType, "read": types.StringType,
						"update": types.StringType, "delete": types.StringType,
					})},
				}
				diags := state.Set(ctx, &model)
				require.False(t, diags.HasError(), diags)
				for range tc.reads {
					resp := resource.ReadResponse{State: state}
					r.Read(ctx, resource.ReadRequest{State: state}, &resp)
					require.Equal(t, tc.wantError, resp.Diagnostics.HasError(), resp.Diagnostics)
					require.Equal(t, tc.wantGone, resp.State.Raw.IsNull(), "an intermittent directory 404 must not recreate a group")
					state = resp.State
				}
				if tc.wantError {
					require.Equal(t, 1, calls, "authorization failures must not be retried")
				}
				if tc.wantGone {
					require.Greater(t, calls, 1, "deletion must be confirmed by repeated not-found responses")
				}
			})
		})
	}
}
