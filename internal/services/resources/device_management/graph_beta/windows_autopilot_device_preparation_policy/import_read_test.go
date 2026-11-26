package graphBetaWindowsAutopilotDevicePreparationPolicy

import (
	"context"
	"net/http"
	"testing"

	providerclient "github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	builders "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors/graph_beta/device_management"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/jarcoal/httpmock"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models"
	"github.com/stretchr/testify/require"
)

// A fresh import must recover configured values from Graph, without a prior plan.
func TestUnitAutopilotPreparationImportOOBESettings(t *testing.T) {
	for _, enabled := range []bool{false, true} {
		t.Run(types.BoolValue(enabled).String(), func(t *testing.T) {
			message := "Contact the acceptance test support desk."
			response := models.NewDeviceManagementConfigurationSettingCollectionResponse()
			response.SetValue([]models.DeviceManagementConfigurationSettingable{
				builders.ConstructStringSimpleSettingInstance(SettingDefCustomErrorMessage, message, SettingInstanceTemplateCustomErrorMessage, SettingValueTemplateCustomErrorMessage),
				builders.ConstructBoolChoiceSettingInstance(SettingDefAllowSkip, enabled, SettingInstanceTemplateAllowSkip, SettingValueTemplateAllowSkip),
				builders.ConstructBoolChoiceSettingInstance(SettingDefAllowDiagnostics, enabled, SettingInstanceTemplateAllowDiagnostics, SettingValueTemplateAllowDiagnostics),
			})
			state := WindowsAutopilotDevicePreparationPolicyResourceModel{TemplateId: types.StringValue(TemplateIDUserDriven)}
			require.NoError(t, mapSettingsToState(context.Background(), &state, response))
			require.Equal(t, types.StringValue(message), state.OOBESettings.CustomErrorMessage)
			require.Equal(t, types.BoolValue(enabled), state.OOBESettings.AllowSkip)
			require.Equal(t, types.BoolValue(enabled), state.OOBESettings.AllowDiagnostics)
		})
	}
}

func TestUnitAutopilotPreparationReadEnrollmentTarget(t *testing.T) {
	for _, test := range []struct {
		name      string
		status    int
		target    string
		wantError bool
	}{
		{"import", 200, "22222222-2222-4222-8222-222222222222", false},
		{"cleared target", 200, "", false},
		{"missing child", 404, "", true},
		{"forbidden child", 403, "", true},
	} {
		t.Run(test.name, func(t *testing.T) {
			ctx := context.Background()
			transport := httpmock.NewMockTransport()
			id := "11111111-1111-4111-8111-111111111111"
			groupID := "22222222-2222-4222-8222-222222222222"
			base := "https://graph.microsoft.com/beta/deviceManagement/configurationPolicies/" + id
			transport.RegisterResponder("GET", base, httpmock.NewJsonResponderOrPanic(200, map[string]any{
				"id": id, "name": "imported-policy", "templateReference": map[string]any{"templateId": TemplateIDAutomatic},
			}))
			transport.RegisterResponder("GET", base+"/settings", httpmock.NewJsonResponderOrPanic(200, map[string]any{"value": []any{}}))
			transport.RegisterResponder("GET", base+"/assignments", httpmock.NewJsonResponderOrPanic(200, map[string]any{"value": []any{}}))
			targetURL := "https://graph.microsoft.com/beta/deviceManagement/configurationPolicies('" + id + "')/retrieveEnrollmentTimeDeviceMembershipTarget"
			calls := 0
			transport.RegisterResponder("GET", targetURL, func(req *http.Request) (*http.Response, error) {
				calls++
				if test.status != 200 {
					return httpmock.NewJsonResponse(test.status, map[string]any{"error": map[string]any{"code": "Request_ResourceNotFound", "message": "Child action unavailable"}})
				}
				return httpmock.NewJsonResponse(200, map[string]any{"enrollmentTimeDeviceMembershipTargetValidationStatuses": []any{map[string]any{"targetId": test.target}}})
			})
			r := NewWindowsAutopilotDevicePreparationPolicyResource().(*WindowsAutopilotDevicePreparationPolicyResource)
			r.client = providerclient.NewMockGraphClients(&http.Client{Transport: transport}).MockBetaClient
			var schema resource.SchemaResponse
			r.Schema(ctx, resource.SchemaRequest{}, &schema)
			model := WindowsAutopilotDevicePreparationPolicyResourceModel{ID: types.StringValue(id), DeviceSecurityGroup: types.StringNull(), Timeouts: timeouts.Value{Object: types.ObjectNull(map[string]attr.Type{"create": types.StringType, "read": types.StringType, "update": types.StringType, "delete": types.StringType})}}
			model.RoleScopeTagIds = types.SetNull(types.StringType)
			model.Assignments = types.SetNull(schema.Schema.Attributes["assignments"].GetType().(types.SetType).ElemType)
			state := tfsdk.State{Schema: schema.Schema}
			diags := state.Set(ctx, &model)
			require.False(t, diags.HasError(), diags)
			resp := resource.ReadResponse{State: state}
			r.Read(ctx, resource.ReadRequest{State: state}, &resp)
			require.Equal(t, test.wantError, resp.Diagnostics.HasError(), resp.Diagnostics)
			require.False(t, resp.State.Raw.IsNull(), "a child failure must preserve the parent policy")
			if test.wantError {
				require.Equal(t, 1, calls)
				return
			}
			var actual WindowsAutopilotDevicePreparationPolicyResourceModel
			require.False(t, resp.State.Get(ctx, &actual).HasError())
			if test.target == "" {
				require.True(t, actual.DeviceSecurityGroup.IsNull())
			} else {
				require.Equal(t, types.StringValue(groupID), actual.DeviceSecurityGroup)
			}
			require.Equal(t, 1, calls)
		})
	}
}
