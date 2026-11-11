package graphBetaCloudPcProvisioningPolicy

import (
	"context"
	"testing"

	"encoding/json"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/microsoft/kiota-abstractions-go/authentication"
	serialization "github.com/microsoft/kiota-serialization-json-go"
	msgraph "github.com/microsoftgraph/msgraph-beta-sdk-go"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
)

func TestUnitResourceCloudPcProvisioningPolicy_DeveloperReady(t *testing.T) {
	ctx := context.Background()
	const imageID = "microsoftwindowsdesktop_windows-ent-cpc_win11-26h2-ent-cpc-devready"
	var schemaResponse resource.SchemaResponse
	NewCloudPcProvisioningPolicyResource().Schema(ctx, resource.SchemaRequest{}, &schemaResponse)
	attribute := schemaResponse.Schema.Attributes["image_id"].(schema.StringAttribute)
	var validation validator.StringResponse
	for _, v := range attribute.Validators {
		v.ValidateString(ctx, validator.StringRequest{Path: path.Root("image_id"), ConfigValue: types.StringValue(imageID)}, &validation)
	}
	require.False(t, validation.Diagnostics.HasError(), validation.Diagnostics)
	for _, tc := range []struct {
		name, response string
		status         int
		wantError      bool
	}{
		{"metadata_name", `{"id":"` + imageID + `","displayName":"Gallery supplied display name","status":"supported"}`, 200, false},
		{"warning", `{"id":"` + imageID + `","displayName":"Gallery warning display name","status":"supportedWithWarning"}`, 200, false},
		{"unavailable", `{"error":{"code":"ResourceNotFound","message":"Image unavailable"}}`, 404, true},
		{"unsupported", `{"id":"` + imageID + `","displayName":"Unsupported image","status":"notSupported"}`, 200, true},
		{"missing_name", `{"id":"` + imageID + `","status":"supported"}`, 200, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			calls := 0
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				calls++
				require.Equal(t, http.MethodGet, r.Method)
				require.Equal(t, "/beta/deviceManagement/virtualEndpoint/galleryImages/"+imageID, r.URL.Path)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tc.status)
				_, _ = w.Write([]byte(tc.response))
			}))
			defer server.Close()
			adapter, err := msgraph.NewGraphRequestAdapter(&authentication.AnonymousAuthenticationProvider{})
			require.NoError(t, err)
			adapter.SetBaseUrl(server.URL + "/beta")
			body, err := constructResource(ctx, &CloudPcProvisioningPolicyResourceModel{ID: types.StringValue("existing-policy"), ImageId: types.StringValue(imageID), ScopeIds: types.SetNull(types.StringType)}, msgraph.NewGraphServiceClient(adapter))
			require.Equal(t, 1, calls)
			if tc.wantError {
				require.Error(t, err)
				require.Nil(t, body)
				return
			}
			require.NoError(t, err)
			var metadata map[string]any
			require.NoError(t, json.Unmarshal([]byte(tc.response), &metadata))
			require.Equal(t, metadata["displayName"], *body.GetImageDisplayName())
			writer := serialization.NewJsonSerializationWriter()
			require.NoError(t, writer.WriteObjectValue("", body))
			encoded, err := writer.GetSerializedContent()
			require.NoError(t, err)
			var payload map[string]any
			require.NoError(t, json.Unmarshal(encoded, &payload))
			require.Equal(t, imageID, payload["imageId"])
			require.Equal(t, metadata["displayName"], payload["imageDisplayName"])
			var state CloudPcProvisioningPolicyResourceModel
			MapRemoteStateToTerraform(ctx, &state, body)
			require.Equal(t, imageID, state.ImageId.ValueString())
			require.Equal(t, metadata["displayName"], state.ImageDisplayName.ValueString())
		})
	}
}
