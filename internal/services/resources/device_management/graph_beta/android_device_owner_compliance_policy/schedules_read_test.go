package graphBetaAndroidDeviceOwnerCompliancePolicy

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/microsoft/kiota-abstractions-go/authentication"
	msgraph "github.com/microsoftgraph/msgraph-beta-sdk-go"
	"github.com/stretchr/testify/require"
)

func TestUnitAndroidComplianceReadScheduledActions(t *testing.T) {
	for _, tc := range []struct {
		name      string
		empty     bool
		forbidden bool
	}{
		{name: "import_existing_actions"},
		{name: "empty_actions", empty: true},
		{name: "read_failure_preserves_state", forbidden: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			calls := 0
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				calls++
				w.Header().Set("Content-Type", "application/json")
				if req.Method != http.MethodGet || req.URL.Path != "/beta/deviceManagement/deviceCompliancePolicies/policy-id" {
					t.Errorf("unexpected request: %s %s", req.Method, req.URL)
					w.WriteHeader(http.StatusBadRequest)
					return
				}
				if tc.forbidden {
					w.WriteHeader(http.StatusForbidden)
					fmt.Fprint(w, `{"error":{"code":"Forbidden","message":"Cannot read policy"}}`)
					return
				}
				// Graph omits navigation properties unless explicitly expanded.
				expand := req.URL.Query().Get("$expand")
				if !strings.Contains(expand, "assignments") {
					t.Error("assignment expansion must be retained")
				}
				schedules := ""
				if strings.Contains(expand, "scheduledActionsForRule($expand=scheduledActionConfigurations)") {
					schedules = `,"scheduledActionsForRule":[]`
					if !tc.empty {
						schedules = `,"scheduledActionsForRule":[{"id":"rule-id","ruleName":null,"scheduledActionConfigurations":[{"id":"block-id","actionType":"block","gracePeriodHours":72,"notificationTemplateId":"00000000-0000-0000-0000-000000000000","notificationMessageCCList":[]},{"id":"push-id","actionType":"pushNotification","gracePeriodHours":0,"notificationTemplateId":"00000000-0000-0000-0000-000000000000","notificationMessageCCList":[]}]}]`
					}
				}
				fmt.Fprintf(w, `{"@odata.type":"#microsoft.graph.androidDeviceOwnerCompliancePolicy","id":"policy-id","displayName":"Example Android policy","roleScopeTagIds":["0"],"assignments":[]%s}`, schedules)
			}))
			t.Cleanup(server.Close)
			adapter, err := msgraph.NewGraphRequestAdapter(&authentication.AnonymousAuthenticationProvider{})
			require.NoError(t, err)
			adapter.SetBaseUrl(server.URL + "/beta")
			r := &AndroidDeviceOwnerCompliancePolicyResource{client: msgraph.NewGraphServiceClient(adapter)}
			var schema resource.SchemaResponse
			r.Schema(ctx, resource.SchemaRequest{}, &schema)
			stateType := schema.Schema.Type().TerraformType(ctx).(tftypes.Object)
			attributes := make(map[string]tftypes.Value, len(stateType.AttributeTypes))
			for name, attributeType := range stateType.AttributeTypes {
				attributes[name] = tftypes.NewValue(attributeType, nil)
			}
			attributes["id"] = tftypes.NewValue(tftypes.String, "policy-id")
			state := tfsdk.State{Schema: schema.Schema, Raw: tftypes.NewValue(stateType, attributes)}
			response := resource.ReadResponse{State: state}
			r.Read(ctx, resource.ReadRequest{State: state}, &response)
			require.Equal(t, tc.forbidden, response.Diagnostics.HasError(), "%v", response.Diagnostics)
			if tc.forbidden {
				require.True(t, state.Raw.Equal(response.State.Raw))
				return
			}
			var actual DeviceCompliancePolicyResourceModel
			require.False(t, response.State.Get(ctx, &actual).HasError())
			if tc.empty {
				require.True(t, actual.ScheduledActionsForRule.IsNull())
				return
			}
			var rules []ScheduledActionForRuleModel
			require.False(t, actual.ScheduledActionsForRule.IsNull(), "existing schedules must be read during import")
			require.False(t, actual.ScheduledActionsForRule.ElementsAs(ctx, &rules, false).HasError())
			require.Len(t, rules, 1)
			require.True(t, rules[0].RuleName.IsNull())
			var actions []ScheduledActionConfigurationModel
			require.False(t, rules[0].ScheduledActionConfigurations.ElementsAs(ctx, &actions, false).HasError())
			require.Len(t, actions, 2)
			gracePeriods := map[string]int32{}
			for _, action := range actions {
				gracePeriods[action.ActionType.ValueString()] = action.GracePeriodHours.ValueInt32()
				require.Equal(t, "00000000-0000-0000-0000-000000000000", action.NotificationTemplateId.ValueString())
				require.Empty(t, action.NotificationMessageCcList.Elements())
			}
			require.Equal(t, map[string]int32{"block": 72, "pushNotification": 0}, gracePeriods)
			refreshed := resource.ReadResponse{State: response.State}
			r.Read(ctx, resource.ReadRequest{State: response.State}, &refreshed)
			require.False(t, refreshed.Diagnostics.HasError(), "%v", refreshed.Diagnostics)
			require.True(t, response.State.Raw.Equal(refreshed.State.Raw), "subsequent refresh must preserve imported schedules")
			require.Equal(t, 2, calls)
		})
	}
}
