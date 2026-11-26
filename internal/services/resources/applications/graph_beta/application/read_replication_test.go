package graphBetaApplication

import (
	"context"
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
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
	"github.com/stretchr/testify/require"
)

func TestUnitApplicationReadReplication(t *testing.T) {
	for _, test := range []struct {
		name                             string
		parentMissing, ownersMissing     int
		deadline, wantMissing, wantError bool
		wantParent, wantOwners           int
	}{
		{"parent_replica", 1, 0, false, false, false, 2, 1},
		{"owners_replica", 0, 1, false, false, false, 1, 2},
		{"confirmed_parent_absence", 100, 0, false, true, false, 16, 0},
		{"unavailable_owners_preserves_parent", 0, 100, false, false, true, 1, 16},
		{"parent_timeout", 100, 0, true, false, true, 1, 0},
		{"owners_timeout", 0, 100, true, false, true, 1, 1},
	} {
		t.Run(test.name, func(t *testing.T) {
			synctest.Test(t, func(t *testing.T) {
				const id = "11111111-1111-4111-8111-111111111111"
				const appID = "22222222-2222-4222-8222-222222222222"
				const endpoint = "https://graph.microsoft.com/beta/applications/" + id
				transport := httpmock.NewMockTransport()
				parentGets, ownerGets := 0, 0
				missing := map[string]any{"error": map[string]any{"code": "Request_ResourceNotFound", "message": "Resource missing from this replica"}}
				transport.RegisterResponder("GET", endpoint, func(_ *http.Request) (*http.Response, error) {
					parentGets++
					if parentGets <= test.parentMissing {
						return httpmock.NewJsonResponse(404, missing)
					}
					return httpmock.NewJsonResponse(200, map[string]any{"id": id, "appId": appID, "displayName": "Read replication"})
				})
				transport.RegisterResponder("GET", endpoint+"/owners", func(_ *http.Request) (*http.Response, error) {
					ownerGets++
					if ownerGets <= test.ownersMissing {
						return httpmock.NewJsonResponse(404, missing)
					}
					return httpmock.NewJsonResponse(200, map[string]any{"value": []any{}})
				})
				r := NewApplicationResource().(*ApplicationResource)
				r.client = providerclient.NewMockGraphClients(&http.Client{Transport: transport}).MockBetaClient
				var schema resource.SchemaResponse
				r.Schema(context.Background(), resource.SchemaRequest{}, &schema)
				remote := graphmodels.NewApplication()
				remote.SetId(new(id))
				remote.SetAppId(new(appID))
				var model ApplicationResourceModel
				MapRemoteStateToTerraform(context.Background(), &model, remote)
				model.OwnerUserIds = types.SetNull(types.StringType)
				model.Timeouts = timeouts.Value{Object: types.ObjectNull(map[string]attr.Type{"create": types.StringType, "read": types.StringType, "update": types.StringType, "delete": types.StringType})}
				state := tfsdk.State{Schema: schema.Schema}
				require.False(t, state.Set(context.Background(), &model).HasError())
				ctx := context.Background()
				if test.deadline {
					var cancel context.CancelFunc
					ctx, cancel = context.WithTimeout(ctx, time.Second)
					defer cancel()
				}
				response := resource.ReadResponse{State: state}
				r.Read(ctx, resource.ReadRequest{State: state}, &response)
				require.Equal(t, test.wantError, response.Diagnostics.HasError(), response.Diagnostics)
				require.Equal(t, test.wantMissing, response.State.Raw.IsNull())
				require.Equal(t, test.wantParent, parentGets)
				require.Equal(t, test.wantOwners, ownerGets)
			})
		})
	}
}
