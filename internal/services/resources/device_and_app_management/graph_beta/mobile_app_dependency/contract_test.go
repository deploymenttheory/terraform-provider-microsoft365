package graphBetaMobileAppDependency

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/microsoft/kiota-abstractions-go/authentication"
	kiotahttp "github.com/microsoft/kiota-http-go"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
	"github.com/stretchr/testify/require"
)

const (
	testSource = "11111111-1111-1111-1111-111111111111"
	testTarget = "22222222-2222-2222-2222-222222222222"
	testID     = testSource + "_" + testTarget
)

func dependencyTestModel() MobileAppDependencyResourceModel {
	return MobileAppDependencyResourceModel{
		ID: types.StringValue(testID), SourceID: types.StringValue(testSource),
		TargetID: types.StringValue(testTarget), DependencyType: types.StringValue("detect"),
		Timeouts: timeouts.Value{Object: types.ObjectNull(map[string]attr.Type{
			"create": types.StringType,
			"read":   types.StringType,
			"update": types.StringType,
			"delete": types.StringType,
		})},
	}
}

func dependencyTestState(
	t *testing.T,
	r *MobileAppDependencyResource,
	model MobileAppDependencyResourceModel,
) tfsdk.State {
	t.Helper()
	var schema resource.SchemaResponse
	r.Schema(context.Background(), resource.SchemaRequest{}, &schema)
	state := tfsdk.State{Schema: schema.Schema}
	require.False(t, state.Set(context.Background(), model).HasError())
	return state
}

func dependencyTestClient(t *testing.T, handler http.HandlerFunc) *MobileAppDependencyResource {
	t.Helper()
	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)
	adapter, err := kiotahttp.NewNetHttpRequestAdapter(
		&authentication.AnonymousAuthenticationProvider{},
	)
	require.NoError(t, err)
	adapter.SetBaseUrl(server.URL)
	r := NewMobileAppDependencyResource().(*MobileAppDependencyResource)
	r.client = msgraphbetasdk.NewGraphServiceClient(adapter)
	return r
}

func TestUnitDependencyReadUsesCollectionWithSavedOrImportedID(t *testing.T) {
	for _, imported := range []bool{false, true} {
		t.Run(fmt.Sprint(imported), func(t *testing.T) {
			badRoute, err := os.ReadFile(
				"tests/responses/validate_get/get_dependency_bad_route.json",
			)
			require.NoError(t, err)
			r := dependencyTestClient(t, func(w http.ResponseWriter, req *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				if req.URL.Path != "/deviceAppManagement/mobileApps/"+testSource+"/relationships" {
					w.WriteHeader(http.StatusBadRequest)
					_, _ = w.Write(badRoute)
					return
				}
				_, _ = fmt.Fprintf(
					w,
					`{"value":[{"@odata.type":"#microsoft.graph.mobileAppDependency","id":%q,"sourceId":%q,"targetId":%q,"targetType":"child","dependencyType":"detect"}]}`,
					testID,
					testSource,
					testTarget,
				)
			})
			model := dependencyTestModel()
			if imported {
				model.SourceID, model.TargetID = types.StringNull(), types.StringNull()
			}
			state := dependencyTestState(t, r, model)
			resp := resource.ReadResponse{State: state}
			r.Read(context.Background(), resource.ReadRequest{State: state}, &resp)
			require.False(t, resp.Diagnostics.HasError(), "%v", resp.Diagnostics)
			require.False(t, resp.State.Raw.IsNull(), "refresh must retain the relationship")
			var actual MobileAppDependencyResourceModel
			require.False(t, resp.State.Get(context.Background(), &actual).HasError())
			require.Equal(t, testID, actual.ID.ValueString())
			require.Equal(t, testSource, actual.SourceID.ValueString())
		})
	}
}

func TestUnitDependencyReadFailuresRetainState(t *testing.T) {
	for _, status := range []int{400, 403} {
		t.Run(fmt.Sprint(status), func(t *testing.T) {
			r := dependencyTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(status)
				_, _ = fmt.Fprint(
					w,
					`{"error":{"code":"TestError","message":"Relationship request failed"}}`,
				)
			})
			state := dependencyTestState(t, r, dependencyTestModel())
			resp := resource.ReadResponse{State: state}
			r.Read(context.Background(), resource.ReadRequest{State: state}, &resp)
			require.True(t, resp.Diagnostics.HasError())
			require.True(t, resp.State.Raw.Equal(state.Raw))
		})
	}
}

func TestUnitDependencyMissingReadVersusCreateReadback(t *testing.T) {
	r := dependencyTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprint(w, `{"value":[]}`)
	})
	state := dependencyTestState(t, r, dependencyTestModel())
	resp := resource.ReadResponse{State: state}
	r.Read(context.Background(), resource.ReadRequest{State: state}, &resp)
	require.False(t, resp.Diagnostics.HasError(), "%v", resp.Diagnostics)
	require.True(t, resp.State.Raw.IsNull())

	resp = resource.ReadResponse{State: state}
	ctx := context.WithValue(
		context.Background(),
		"retry_operation",
		"Create",
	) //nolint:staticcheck // Shared CRUD retry context convention.
	r.Read(ctx, resource.ReadRequest{State: state}, &resp)
	require.True(t, resp.Diagnostics.HasError())
	require.True(t, resp.State.Raw.Equal(state.Raw))
}
