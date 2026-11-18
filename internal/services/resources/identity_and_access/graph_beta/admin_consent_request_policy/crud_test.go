package graphBetaAdminConsentRequestPolicy

import (
	"context"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/require"

	providerclient "github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
)

func policyTestState(t *testing.T, r *AdminConsentRequestPolicyResource) tfsdk.State {
	t.Helper()
	var schema resource.SchemaResponse
	r.Schema(context.Background(), resource.SchemaRequest{}, &schema)
	state := tfsdk.State{Schema: schema.Schema}
	data := AdminConsentRequestPolicyResourceModel{
		ID:                    types.StringValue(singletonID),
		IsEnabled:             types.BoolValue(false),
		NotifyReviewers:       types.BoolValue(false),
		RemindersEnabled:      types.BoolValue(false),
		RequestDurationInDays: types.Int32Value(0),
		Reviewers:             types.SetValueMust(types.ObjectType{AttrTypes: reviewerAttributeTypes}, nil),
		Version:               types.Int32Value(0),
		Timeouts: timeouts.Value{
			Object: types.ObjectNull(map[string]attr.Type{"create": types.StringType, "read": types.StringType, "update": types.StringType, "delete": types.StringType}),
		},
	}
	require.False(t, state.Set(context.Background(), &data).HasError())
	return state
}

func TestUnitAdminConsentPolicyReadErrorsPreserveSingleton(t *testing.T) {
	for _, status := range []int{400, 403, 404, 204} {
		t.Run(http.StatusText(status), func(t *testing.T) {
			transport := httpmock.NewMockTransport()
			transport.RegisterResponder("GET", "https://graph.microsoft.com/beta/policies/adminConsentRequestPolicy", func(_ *http.Request) (*http.Response, error) {
				if status == 204 {
					return httpmock.NewStringResponse(204, ""), nil
				}
				return httpmock.NewJsonResponse(status, map[string]any{"error": map[string]any{"code": "Request_BadRequest", "message": "probe error"}})
			})
			r := NewAdminConsentRequestPolicyResource().(*AdminConsentRequestPolicyResource)
			r.client = providerclient.NewMockGraphClients(&http.Client{Transport: transport}).MockBetaClient
			state := policyTestState(t, r)
			resp := resource.ReadResponse{State: state}
			r.Read(context.Background(), resource.ReadRequest{State: state}, &resp)
			require.True(t, resp.Diagnostics.HasError(), resp.Diagnostics)
			require.False(t, resp.State.Raw.IsNull())
		})
	}
}

func TestUnitAdminConsentPolicyImportAndDelete(t *testing.T) {
	r := NewAdminConsentRequestPolicyResource().(*AdminConsentRequestPolicyResource)
	var imported resource.ImportStateResponse
	r.ImportState(context.Background(), resource.ImportStateRequest{ID: "wrong"}, &imported)
	require.True(t, imported.Diagnostics.HasError())
	state := policyTestState(t, r)
	resp := resource.DeleteResponse{State: state}
	r.Delete(context.Background(), resource.DeleteRequest{State: state}, &resp)
	require.False(t, resp.Diagnostics.HasError())
	require.True(t, resp.State.Raw.IsNull(), "destroy must release state without requiring a Graph client")
}

func TestUnitAdminConsentPolicyRejectsUnknownReviewers(t *testing.T) {
	r := NewAdminConsentRequestPolicyResource().(*AdminConsentRequestPolicyResource)
	state := policyTestState(t, r)
	var data AdminConsentRequestPolicyResourceModel
	require.False(t, state.Get(context.Background(), &data).HasError())
	data.Reviewers = types.SetUnknown(types.ObjectType{AttrTypes: reviewerAttributeTypes})
	_, err := constructResource(context.Background(), &data)
	require.ErrorIs(t, err, errReviewers)
}
