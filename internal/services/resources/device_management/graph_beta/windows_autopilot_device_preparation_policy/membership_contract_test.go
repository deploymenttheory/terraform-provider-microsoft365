package graphBetaWindowsAutopilotDevicePreparationPolicy

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/microsoft/kiota-abstractions-go/authentication"
	kiotahttp "github.com/microsoft/kiota-http-go"
	msgraph "github.com/microsoftgraph/msgraph-beta-sdk-go"
	"github.com/stretchr/testify/require"
)

func TestUnitResourceWindowsAutopilotDevicePreparationPolicy_MembershipReadContract(t *testing.T) {
	for _, status := range []int{http.StatusOK, http.StatusNotFound, http.StatusBadRequest} {
		t.Run(fmt.Sprint(status), func(t *testing.T) {
			calls := 0
			const id = "00000000-0000-0000-0000-000000000001"
			const target = "00000000-0000-0000-0000-000000000002"
			server := httptest.NewServer(
				http.HandlerFunc(func(w http.ResponseWriter, q *http.Request) {
					w.Header().Set("Content-Type", "application/json")
					var body string
					switch q.URL.Path {
					case "/beta/deviceManagement/configurationPolicies/" + id:
						body = `{"id":"` + id + `","name":"contract policy","templateReference":{"templateId":"` + TemplateIDAutomatic + `"}}`
					case "/beta/deviceManagement/configurationPolicies/" + id + "/settings",
						"/beta/deviceManagement/configurationPolicies/" + id + "/assignments":
						body = `{"value":[]}`
					case "/beta/deviceManagement/configurationPolicies('" + id + "')/retrieveEnrollmentTimeDeviceMembershipTarget":
						calls++
						if status != http.StatusOK {
							w.WriteHeader(status)
							body = `{"error":{"code":"ResourceNotFound","message":"Membership read failed"}}`
						} else {
							body = `{"validationSucceeded":true,"enrollmentTimeDeviceMembershipTargetValidationStatuses":[{"targetId":"` + target + `","validationSucceeded":true}]}`
						}
					default:
						t.Errorf("unexpected request URI: %s", q.URL.RequestURI())
						w.WriteHeader(http.StatusNotFound)
						body = `{"error":{"code":"ResourceNotFound","message":"Invalid route"}}`
					}
					_, err := fmt.Fprint(w, body)
					require.NoError(t, err)
				}),
			)
			t.Cleanup(server.Close)
			adapter, err := kiotahttp.NewNetHttpRequestAdapter(
				&authentication.AnonymousAuthenticationProvider{},
			)
			require.NoError(t, err)
			adapter.SetBaseUrl(server.URL + "/beta")
			r := NewWindowsAutopilotDevicePreparationPolicyResource().(*WindowsAutopilotDevicePreparationPolicyResource)
			r.client = msgraph.NewGraphServiceClient(adapter)
			var schema resource.SchemaResponse
			r.Schema(t.Context(), resource.SchemaRequest{}, &schema)
			state := tfsdk.State{Schema: schema.Schema}
			model := WindowsAutopilotDevicePreparationPolicyResourceModel{
				ID: types.StringValue(
					id,
				),
				Name:                types.StringValue("contract policy"),
				DeviceSecurityGroup: types.StringValue(target),
				RoleScopeTagIds:     types.SetNull(types.StringType),
				Assignments: types.SetNull(
					schema.Schema.Attributes["assignments"].GetType().(types.SetType).ElemType,
				),
				Timeouts: timeouts.Value{
					Object: types.ObjectNull(
						map[string]attr.Type{
							"create": types.StringType,
							"read":   types.StringType,
							"update": types.StringType,
							"delete": types.StringType,
						},
					),
				},
			}
			require.False(t, state.Set(t.Context(), &model).HasError())
			response := resource.ReadResponse{State: state}
			r.Read(t.Context(), resource.ReadRequest{State: state}, &response)
			require.Equal(t, 1, calls, "membership request must use the exact beta route")
			if status != http.StatusOK {
				require.True(t, response.Diagnostics.HasError(), "%v", response.Diagnostics)
				require.True(
					t,
					state.Raw.Equal(response.State.Raw),
					"failed membership read must preserve existing state",
				)
				require.True(
					t,
					strings.Contains(fmt.Sprint(response.Diagnostics), "Membership read failed"),
				)
			} else {
				require.False(t, response.Diagnostics.HasError(), "%v", response.Diagnostics)
				var actual WindowsAutopilotDevicePreparationPolicyResourceModel
				require.False(t, response.State.Get(t.Context(), &actual).HasError())
				require.Equal(t, target, actual.DeviceSecurityGroup.ValueString())
			}
		})
	}
}
