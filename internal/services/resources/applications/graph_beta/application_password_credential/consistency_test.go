package graphBetaApplicationPasswordCredential

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

func credentialTestState(t *testing.T, r *ApplicationPasswordCredentialResource) tfsdk.State {
	t.Helper()
	ctx := context.Background()
	var schema resource.SchemaResponse
	r.Schema(ctx, resource.SchemaRequest{}, &schema)
	state := tfsdk.State{Schema: schema.Schema}
	model := ApplicationPasswordCredentialResourceModel{
		ApplicationID: types.StringValue("11111111-1111-1111-1111-111111111111"), DisplayName: types.StringValue("Password credential regression"),
		StartDateTime: types.StringUnknown(), EndDateTime: types.StringUnknown(), KeyID: types.StringUnknown(),
		SecretText: types.StringUnknown(), Hint: types.StringUnknown(), CustomKeyIdentifier: types.StringUnknown(),
		Timeouts: timeouts.Value{Object: types.ObjectNull(map[string]attr.Type{
			"create": types.StringType, "read": types.StringType, "update": types.StringType, "delete": types.StringType,
		})},
	}
	require.False(t, state.Set(ctx, model).HasError())
	return state
}

func TestUnitResourceApplicationPasswordCredential_04_CreateRejectsPermanentErrors(t *testing.T) {
	for _, status := range []int{http.StatusBadRequest, http.StatusForbidden} {
		t.Run(http.StatusText(status), func(t *testing.T) {
			synctest.Test(t, func(t *testing.T) {
				transport := httpmock.NewMockTransport()
				calls := 0
				transport.RegisterResponder(
					"POST",
					"https://graph.microsoft.com/beta/applications/11111111-1111-1111-1111-111111111111/addPassword",
					func(_ *http.Request) (*http.Response, error) {
						calls++
						return httpmock.NewJsonResponse(status, map[string]any{"error": map[string]string{"code": http.StatusText(status), "message": http.StatusText(status)}})
					},
				)
				r := &ApplicationPasswordCredentialResource{client: providerclient.NewMockGraphClients(&http.Client{Transport: transport}).MockBetaClient}
				state := credentialTestState(t, r)
				response := resource.CreateResponse{State: state}
				r.Create(context.Background(), resource.CreateRequest{Plan: tfsdk.Plan{Schema: state.Schema, Raw: state.Raw}}, &response)
				require.True(t, response.Diagnostics.HasError())
				require.Equal(t, 1, calls, "validation and permission errors must not trigger a second password creation attempt")
			})
		})
	}
}

func TestUnitResourceApplicationPasswordCredential_05_CreateHonorsDeadline(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		transport := httpmock.NewMockTransport()
		calls := 0
		transport.RegisterResponder(
			"POST",
			"https://graph.microsoft.com/beta/applications/11111111-1111-1111-1111-111111111111/addPassword",
			func(_ *http.Request) (*http.Response, error) {
				calls++
				return httpmock.NewJsonResponse(
					http.StatusNotFound,
					map[string]any{"error": map[string]string{"code": "Request_ResourceNotFound", "message": "Parent application not found"}},
				)
			},
		)
		r := &ApplicationPasswordCredentialResource{client: providerclient.NewMockGraphClients(&http.Client{Transport: transport}).MockBetaClient}
		state := credentialTestState(t, r)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		response := resource.CreateResponse{State: state}
		r.Create(ctx, resource.CreateRequest{Plan: tfsdk.Plan{Schema: state.Schema, Raw: state.Raw}}, &response)
		require.True(t, response.Diagnostics.HasError())
		require.Contains(t, fmt.Sprint(response.Diagnostics.Errors()), context.DeadlineExceeded.Error())
		require.Equal(t, 1, calls)
	})
}

func TestUnitResourceApplicationPasswordCredential_03_CreateAfterTransientParentNotFound(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		transport := httpmock.NewMockTransport()
		const endpoint = "https://graph.microsoft.com/beta/applications/11111111-1111-1111-1111-111111111111"
		const keyID = "22222222-2222-2222-2222-222222222222"
		calls := 0
		transport.RegisterResponder("POST", endpoint+"/addPassword", func(_ *http.Request) (*http.Response, error) {
			calls++
			if calls == 1 {
				return httpmock.NewJsonResponse(404, map[string]any{"error": map[string]string{"code": "Request_ResourceNotFound", "message": "Resource does not exist"}})
			}
			return httpmock.NewJsonResponse(200, map[string]any{
				"keyId": keyID, "displayName": "Password credential regression", "secretText": "unit-only-secret", "hint": "uni",
				"startDateTime": "2027-01-01T00:00:00Z", "endDateTime": "2028-01-01T00:00:00Z",
			})
		})
		transport.RegisterResponder(
			"GET",
			endpoint,
			func(_ *http.Request) (*http.Response, error) {
				response := httpmock.NewStringResponse(
					200,
					`{"id":"11111111-1111-1111-1111-111111111111","passwordCredentials":[{"keyId":"22222222-2222-2222-2222-222222222222","displayName":"Password credential regression","hint":"uni","startDateTime":"2027-01-01T00:00:00Z","endDateTime":"2028-01-01T00:00:00Z"}]}`,
				)
				response.Header.Set("Content-Type", "application/json")
				return response, nil
			},
		)
		r := &ApplicationPasswordCredentialResource{client: providerclient.NewMockGraphClients(&http.Client{Transport: transport}).MockBetaClient}
		state := credentialTestState(t, r)
		response := resource.CreateResponse{State: state}
		r.Create(context.Background(), resource.CreateRequest{Plan: tfsdk.Plan{Schema: state.Schema, Raw: state.Raw}}, &response)
		require.False(t, response.Diagnostics.HasError(), "%v", response.Diagnostics)
		require.Equal(t, 2, calls, "only the rejected parent-not-found request may be retried")
		var created ApplicationPasswordCredentialResourceModel
		require.False(t, response.State.Get(context.Background(), &created).HasError())
		require.Equal(t, keyID, created.KeyID.ValueString())
		require.Equal(t, "unit-only-secret", created.SecretText.ValueString(), "a metadata GET cannot replace the creation-only secret")
	})
}
