package graphBetaWin32App

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/microsoft/kiota-abstractions-go/authentication"
	kiotahttp "github.com/microsoft/kiota-http-go"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
	"github.com/stretchr/testify/require"
)

func TestUnitResourceWin32App_CreateFailureRetainsDeletableApplication(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	const appID = "11111111-1111-4111-8111-111111111111"
	created, contentFailed, deleted := false, false, false
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch {
		case req.Method == http.MethodPost && req.URL.Path == "/deviceAppManagement/mobileApps":
			created = true
			w.WriteHeader(http.StatusCreated)
			require.NoError(t, json.NewEncoder(w).Encode(map[string]any{"@odata.type": "#microsoft.graph.win32LobApp", "id": appID}))
		case req.Method == http.MethodPost && req.URL.Path == "/deviceAppManagement/mobileApps/"+appID+"/graph.win32LobApp/contentVersions":
			contentFailed = true
			w.WriteHeader(http.StatusForbidden)
			require.NoError(t, json.NewEncoder(w).Encode(map[string]any{"error": map[string]any{"code": "Authorization_RequestDenied", "message": "injected content failure"}}))
		case req.Method == http.MethodDelete && req.URL.Path == "/deviceAppManagement/mobileApps/"+appID:
			deleted = true
			w.WriteHeader(http.StatusNoContent)
		default:
			t.Errorf("unexpected request %s %s", req.Method, req.URL.Path)
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()
	adapter, err := kiotahttp.NewNetHttpRequestAdapter(&authentication.AnonymousAuthenticationProvider{})
	require.NoError(t, err)
	adapter.SetBaseUrl(server.URL)
	r := NewWin32LobAppResource().(*Win32LobAppResource)
	r.client = msgraphbetasdk.NewGraphServiceClient(adapter)
	var schema resource.SchemaResponse
	r.Schema(ctx, resource.SchemaRequest{}, &schema)
	objectType := schema.Schema.Type().TerraformType(ctx).(tftypes.Object)
	values := map[string]tftypes.Value{}
	for name, kind := range objectType.AttributeTypes {
		values[name] = tftypes.NewValue(kind, nil)
		if attribute, ok := schema.Schema.Attributes[name]; ok && attribute.IsComputed() {
			values[name] = tftypes.NewValue(kind, tftypes.UnknownValue)
		}
	}
	values["display_name"] = tftypes.NewValue(tftypes.String, "owned failed-create fixture")
	source := writePlainZipFixture(t, "payload.zip", map[string][]byte{"setup.cmd": []byte("exit /b 0")})
	sourceType := objectType.AttributeTypes["app_installer_zip"].(tftypes.Object)
	sourceValues := map[string]tftypes.Value{}
	for name, kind := range sourceType.AttributeTypes {
		sourceValues[name] = tftypes.NewValue(kind, nil)
	}
	sourceValues["installer_file_path_source"] = tftypes.NewValue(tftypes.String, source)
	values["app_installer_zip"] = tftypes.NewValue(sourceType, sourceValues)
	response := resource.CreateResponse{State: tfsdk.State{Schema: schema.Schema}}
	r.Create(ctx, resource.CreateRequest{Plan: tfsdk.Plan{Schema: schema.Schema, Raw: tftypes.NewValue(objectType, values)}}, &response)
	require.True(t, created, "base application must have been created")
	require.True(t, contentFailed, "failure must occur after Graph allocated the application")
	require.True(t, response.Diagnostics.HasError(), response.Diagnostics)
	require.True(t, response.State.Raw.IsFullyKnown(), "failed Create state must contain no unknown planned values")
	var retainedID string
	require.False(t, response.State.GetAttribute(ctx, path.Root("id"), &retainedID).HasError(), "failed Create must retain its application ID")
	require.Equal(t, appID, retainedID)
	deletion := resource.DeleteResponse{State: response.State}
	r.Delete(ctx, resource.DeleteRequest{State: response.State}, &deletion)
	require.False(t, deletion.Diagnostics.HasError(), deletion.Diagnostics)
	require.True(t, deleted, "retained state must permit normal Delete")
	require.True(t, deletion.State.Raw.IsNull())
}
