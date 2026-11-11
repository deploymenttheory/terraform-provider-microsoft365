package graphBetaMacosDeviceCompliancePolicy

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/microsoft/kiota-abstractions-go/store"
	jsonserialization "github.com/microsoft/kiota-serialization-json-go"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
	"github.com/stretchr/testify/require"
)

func TestUnitConstructScheduledActionCcList(t *testing.T) {
	ctx := context.Background()
	action := graphmodels.NewDeviceComplianceActionItem()
	action.SetActionType(new(graphmodels.BLOCK_DEVICECOMPLIANCEACTIONTYPE))
	action.SetGracePeriodHours(new(int32(72)))
	base, err := mapScheduledActionConfigurationsToState(ctx, []graphmodels.DeviceComplianceActionItemable{action}, nil)
	require.NoError(t, err)
	objectType := base.ElementType(ctx).(types.ObjectType)
	for _, tc := range []struct {
		name  string
		value types.List
		want  []string
	}{
		{"omitted", types.ListNull(types.StringType), nil},
		{"unknown", types.ListUnknown(types.StringType), nil},
		{"empty", types.ListValueMust(types.StringType, []attr.Value{}), []string{}},
		{"populated", types.ListValueMust(types.StringType, []attr.Value{types.StringValue("11111111-1111-4111-8111-111111111111"), types.StringValue("22222222-2222-4222-8222-222222222222")}), []string{"11111111-1111-4111-8111-111111111111", "22222222-2222-4222-8222-222222222222"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			values := base.Elements()[0].(types.Object).Attributes()
			values["notification_message_cc_list"] = tc.value
			configs := types.SetValueMust(objectType, []attr.Value{types.ObjectValueMust(objectType.AttrTypes, values)})
			items, err := constructScheduledActionItem(ctx, configs)
			require.NoError(t, err)
			factory := store.NewBackingStoreSerializationWriterProxyFactory(jsonserialization.NewJsonSerializationWriterFactory())
			writer, err := factory.GetSerializationWriter("application/json")
			require.NoError(t, err)
			defer writer.Close()
			require.NoError(t, writer.WriteObjectValue("", items[0]))
			content, err := writer.GetSerializedContent()
			require.NoError(t, err)
			var payload map[string]json.RawMessage
			require.NoError(t, json.Unmarshal(content, &payload))
			if tc.want == nil {
				require.NotContains(t, payload, "notificationMessageCCList")
			} else {
				var recipients []string
				require.NoError(t, json.Unmarshal(payload["notificationMessageCCList"], &recipients))
				require.Equal(t, tc.want, recipients)
			}
		})
	}
}

// TestUnitConstructScheduledActions verifies create and update send a usable rule
// name even when the optional/computed Terraform value is unknown during planning.
func TestUnitConstructScheduledActions(t *testing.T) {
	ctx := context.Background()
	action := graphmodels.NewDeviceComplianceActionItem()
	action.SetActionType(new(graphmodels.BLOCK_DEVICECOMPLIANCEACTIONTYPE))
	action.SetGracePeriodHours(new(int32(72)))
	configs, err := mapScheduledActionConfigurationsToState(
		ctx,
		[]graphmodels.DeviceComplianceActionItemable{action},
		nil,
	)
	require.NoError(t, err)
	for _, tc := range []struct {
		name  string
		value types.String
		want  string
	}{
		{"omitted", types.StringUnknown(), "PasswordRequired"},
		{"null", types.StringNull(), "PasswordRequired"},
		{"explicit", types.StringValue("PasswordRequired"), "PasswordRequired"},
		{"legacy", types.StringValue("unavailable"), "unavailable"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			data := ScheduledActionForRuleModel{
				RuleName:                      tc.value,
				ScheduledActionConfigurations: configs,
			}
			created, err := constructScheduledActionsForPolicyCreation(ctx, data)
			require.NoError(t, err)
			updated, err := constructDeviceComplianceScheduledActionForRulesWithPatchMethod(
				ctx,
				data,
			)
			require.NoError(t, err)
			for _, rules := range [][]graphmodels.DeviceComplianceScheduledActionForRuleable{created, updated.GetDeviceComplianceScheduledActionForRules()} {
				require.Len(t, rules, 1)
				require.Equal(t, tc.want, *rules[0].GetRuleName())
				require.Len(t, rules[0].GetScheduledActionConfigurations(), 1)
				require.Equal(
					t,
					int32(72),
					*rules[0].GetScheduledActionConfigurations()[0].GetGracePeriodHours(),
				)
			}
		})
	}
}
