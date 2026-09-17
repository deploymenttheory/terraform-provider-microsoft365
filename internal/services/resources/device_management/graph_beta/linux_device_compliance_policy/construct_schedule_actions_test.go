package graphBetaLinuxDeviceCompliancePolicy

import (
	"context"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	jsonserialization "github.com/microsoft/kiota-serialization-json-go"
	"github.com/stretchr/testify/require"
)

func TestUnitLinuxComplianceScheduledActionPayload(t *testing.T) {
	configurationType := types.ObjectType{AttrTypes: map[string]attr.Type{
		"action_type":                  types.StringType,
		"grace_period_hours":           types.Int32Type,
		"notification_template_id":     types.StringType,
		"notification_message_cc_list": types.ListType{ElemType: types.StringType},
	}}
	for _, tc := range []struct {
		name       string
		actionType string
		wantError  bool
	}{
		{name: "block_with_default_rule", actionType: "block"},
		{name: "invalid_action", actionType: "invalid", wantError: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			configuration := types.ObjectValueMust(configurationType.AttrTypes, map[string]attr.Value{
				"action_type":                  types.StringValue(tc.actionType),
				"grace_period_hours":           types.Int32Value(24),
				"notification_template_id":     types.StringNull(),
				"notification_message_cc_list": types.ListNull(types.StringType),
			})
			body, err := constructScheduledActions(context.Background(), ScheduledActionForRuleModel{
				RuleName:                      types.StringNull(),
				ScheduledActionConfigurations: types.SetValueMust(configurationType, []attr.Value{configuration}),
			})
			if tc.wantError {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			writer := jsonserialization.NewJsonSerializationWriter()
			t.Cleanup(func() { require.NoError(t, writer.Close()) })
			require.NoError(t, writer.WriteObjectValue("", body))
			encoded, err := writer.GetSerializedContent()
			require.NoError(t, err)
			expected, err := helpers.ParseJSONFile("tests/responses/validate_create/post_scheduled_actions.json")
			require.NoError(t, err)
			require.JSONEq(t, expected, string(encoded))
		})
	}
}
