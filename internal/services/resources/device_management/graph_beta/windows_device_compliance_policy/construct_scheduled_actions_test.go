package graphBetaWindowsDeviceCompliancePolicy

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/microsoft/kiota-abstractions-go/store"
	jsonserialization "github.com/microsoft/kiota-serialization-json-go"
	"github.com/stretchr/testify/require"
)

func TestUnitConstructScheduledActionCcList(t *testing.T) {
	ctx := context.Background()
	objectType := types.ObjectType{AttrTypes: map[string]attr.Type{
		"action_type":                  types.StringType,
		"grace_period_hours":           types.Int32Type,
		"notification_template_id":     types.StringType,
		"notification_message_cc_list": types.ListType{ElemType: types.StringType},
	}}
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
			configs, diags := types.SetValueFrom(ctx, objectType, []ScheduledActionConfigurationModel{{
				ActionType: types.StringValue("block"), GracePeriodHours: types.Int32Value(72),
				NotificationTemplateId: types.StringNull(), NotificationMessageCcList: tc.value,
			}})
			require.False(t, diags.HasError(), "%s", diags)
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
