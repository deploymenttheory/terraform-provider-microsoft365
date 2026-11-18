package graphBetaGroup

import (
	"context"
	"net/http"
	"testing"
	"testing/synctest"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/require"

	providerclient "github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
)

func TestUnitGroupUpdateDirectoryConsistency(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		ctx := context.Background()
		transport := httpmock.NewMockTransport()
		transport.RegisterResponder("PATCH", "https://graph.microsoft.com/beta/groups/11111111-1111-1111-1111-111111111111", httpmock.NewStringResponder(http.StatusNoContent, ""))
		reads := 0
		transport.RegisterResponder("GET", "https://graph.microsoft.com/beta/groups/11111111-1111-1111-1111-111111111111", func(_ *http.Request) (*http.Response, error) {
			reads++
			description := "Old description"
			if reads > 1 {
				description = "Updated description"
			}
			return httpmock.NewJsonResponse(http.StatusOK, map[string]any{
				"id":          "11111111-1111-1111-1111-111111111111",
				"displayName": "Directory consistency", "mailNickname": "directory-consistency",
				"mailEnabled": false, "securityEnabled": true, "description": description,
			})
		})
		r := NewGroupResource().(*GroupResource)
		r.client = providerclient.NewMockGraphClients(&http.Client{Transport: transport}).MockBetaClient
		var schema resource.SchemaResponse
		r.Schema(ctx, resource.SchemaRequest{}, &schema)
		model := GroupResourceModel{
			ID:           types.StringValue("11111111-1111-1111-1111-111111111111"),
			DisplayName:  types.StringValue("Directory consistency"),
			MailNickname: types.StringValue("directory-consistency"),
			MailEnabled:  types.BoolValue(false), SecurityEnabled: types.BoolValue(true),
			Description: types.StringValue("Old description"),
			GroupTypes:  types.SetNull(types.StringType), GroupOwners: types.SetNull(types.StringType),
			GroupMembers: types.SetNull(types.StringType),
			Timeouts: timeouts.Value{Object: types.ObjectNull(map[string]attr.Type{
				"create": types.StringType, "read": types.StringType, "update": types.StringType, "delete": types.StringType,
			})},
		}
		state := tfsdk.State{Schema: schema.Schema}
		require.False(t, state.Set(ctx, &model).HasError())
		model.Description = types.StringValue("Updated description")
		plan := tfsdk.Plan{Schema: schema.Schema}
		require.False(t, plan.Set(ctx, &model).HasError())
		response := resource.UpdateResponse{State: state}
		r.Update(ctx, resource.UpdateRequest{State: state, Plan: plan}, &response)
		require.False(t, response.Diagnostics.HasError(), response.Diagnostics)
		var actual GroupResourceModel
		require.False(t, response.State.Get(ctx, &actual).HasError())
		require.Equal(t, model.ID, actual.ID)
		require.Equal(t, model.Description, actual.Description, "a stale replica must not become the final update state")
		require.GreaterOrEqual(t, reads, 2)
	})
}
