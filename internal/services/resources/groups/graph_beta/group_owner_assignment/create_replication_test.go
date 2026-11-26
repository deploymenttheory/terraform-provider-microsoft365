package graphBetaGroupOwnerAssignment

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"testing/synctest"
	"time"

	providerclient "github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/require"
)

func TestUnitGroupOwnerAssignmentCreateReplication(t *testing.T) {
	for _, test := range []struct {
		name                       string
		misses, staleReads, status int
		deadline                   time.Duration
		wantCalls                  int
		wantError                  bool
		groupMisses                int
	}{
		{"validation group replication", 0, 0, 404, 0, 1, false, 1},
		{"new group", 1, 0, 404, 0, 2, false, 0},
		{"owner reference settling", 0, 2, 404, 0, 1, false, 0},
		{"persistent missing group", 99, 0, 404, 0, 16, true, 0},
		{"forbidden", 1, 0, 403, 0, 1, true, 0},
		{"parent deadline", 99, 0, 404, time.Second, 1, true, 0},
	} {
		t.Run(test.name, func(t *testing.T) {
			synctest.Test(t, func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
				defer cancel()
				if test.deadline > 0 {
					var cancel context.CancelFunc
					ctx, cancel = context.WithTimeout(ctx, test.deadline)
					defer cancel()
				}
				groupID := "11111111-1111-4111-8111-111111111111"
				ownerID := "22222222-2222-4222-8222-222222222222"
				base := "https://graph.microsoft.com/beta"
				transport := httpmock.NewMockTransport()
				groupReads := 0
				transport.RegisterResponder("GET", base+"/groups/"+groupID, func(req *http.Request) (*http.Response, error) {
					groupReads++
					if groupReads <= test.groupMisses {
						return httpmock.NewJsonResponse(404, map[string]any{"error": map[string]any{"code": "Request_ResourceNotFound", "message": "Group validation not yet replicated"}})
					}
					return httpmock.NewJsonResponse(200, map[string]any{"id": groupID, "mailEnabled": false, "securityEnabled": true, "groupTypes": []string{}})
				})
				owner := map[string]any{"@odata.type": "#microsoft.graph.servicePrincipal", "id": ownerID, "displayName": "owned fixture principal"}
				transport.RegisterResponder("GET", base+"/servicePrincipals/"+ownerID, httpmock.NewJsonResponderOrPanic(200, owner))
				calls, reads := 0, 0
				transport.RegisterResponder("POST", base+"/groups/"+groupID+"/owners/$ref", func(req *http.Request) (*http.Response, error) {
					calls++
					var body map[string]any
					require.NoError(t, json.NewDecoder(req.Body).Decode(&body))
					require.Equal(t, base+"/servicePrincipals/"+ownerID, body["@odata.id"], "every retry must retain the owner reference")
					if calls <= test.misses {
						return httpmock.NewJsonResponse(test.status, map[string]any{"error": map[string]any{"code": "Request_ResourceNotFound", "message": "Group reference not yet available"}})
					}
					return httpmock.NewStringResponse(204, ""), nil
				})
				transport.RegisterResponder("GET", base+"/groups/"+groupID+"/owners", func(req *http.Request) (*http.Response, error) {
					reads++
					owners := []any{owner}
					if reads <= test.staleReads {
						owners = []any{}
					}
					return httpmock.NewJsonResponse(200, map[string]any{"value": owners})
				})
				r := NewGroupOwnerAssignmentResource().(*GroupOwnerAssignmentResource)
				r.client = providerclient.NewMockGraphClients(&http.Client{Transport: transport}).MockBetaClient
				var schema resource.SchemaResponse
				r.Schema(ctx, resource.SchemaRequest{}, &schema)
				model := GroupOwnerAssignmentResourceModel{GroupID: types.StringValue(groupID), OwnerID: types.StringValue(ownerID), OwnerObjectType: types.StringValue("ServicePrincipal"), Timeouts: timeouts.Value{Object: types.ObjectNull(map[string]attr.Type{"create": types.StringType, "read": types.StringType, "update": types.StringType, "delete": types.StringType})}}
				state := tfsdk.State{Schema: schema.Schema}
				diags := state.Set(ctx, &model)
				require.False(t, diags.HasError(), diags)
				response := resource.CreateResponse{State: tfsdk.State{Schema: schema.Schema}}
				r.Create(ctx, resource.CreateRequest{Plan: tfsdk.Plan{Schema: state.Schema, Raw: state.Raw}}, &response)
				require.Equal(t, test.wantError, response.Diagnostics.HasError(), response.Diagnostics)
				require.Equal(t, test.wantCalls, calls)
				if !test.wantError {
					var actual GroupOwnerAssignmentResourceModel
					require.False(t, response.State.Get(ctx, &actual).HasError())
					require.Equal(t, types.StringValue("ServicePrincipal"), actual.OwnerType)
					require.Equal(t, types.StringValue(groupID+"/"+ownerID), actual.ID)
				}
			})
		})
	}
}
