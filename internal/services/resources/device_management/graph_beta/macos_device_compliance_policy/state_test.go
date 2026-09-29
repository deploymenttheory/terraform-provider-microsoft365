package graphBetaMacosDeviceCompliancePolicy

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
	"github.com/stretchr/testify/require"
)

// TestUnitMapScheduledActionsForRuleToState covers null API names, import without
// prior state, legacy state, and an authoritative name if Graph returns one.
func TestUnitMapScheduledActionsForRuleToState(t *testing.T) {
	ctx := context.Background()
	action := graphmodels.NewDeviceComplianceScheduledActionForRule()
	action.SetScheduledActionConfigurations([]graphmodels.DeviceComplianceActionItemable{})
	empty := types.ListNull(types.ObjectType{AttrTypes: map[string]attr.Type{}})
	initial, err := mapScheduledActionsForRuleToState(
		ctx,
		[]graphmodels.DeviceComplianceScheduledActionForRuleable{action},
		empty,
	)
	require.NoError(t, err)
	elementType := initial.ElementType(ctx).(types.ObjectType)
	cases := []struct {
		name   string
		prior  types.String
		remote *string
		want   string
	}{
		{name: "unknown plan", prior: types.StringUnknown(), want: "PasswordRequired"},
		{name: "null plan", prior: types.StringNull(), want: "PasswordRequired"},
		{
			name:  "explicit rule",
			prior: types.StringValue("PasswordRequired"),
			want:  "PasswordRequired",
		},
		{name: "legacy rule", prior: types.StringValue("unavailable"), want: "unavailable"},
		{
			name:   "empty remote rule",
			prior:  types.StringValue("unavailable"),
			remote: new(""),
			want:   "unavailable",
		},
		{
			name:   "returned remote rule",
			prior:  types.StringValue("unavailable"),
			remote: new("PasswordRequired"),
			want:   "PasswordRequired",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			attrs := initial.Elements()[0].(types.Object).Attributes()
			attrs["rule_name"] = tc.prior
			prior := types.ListValueMust(
				elementType,
				[]attr.Value{types.ObjectValueMust(elementType.AttrTypes, attrs)},
			)
			action.SetRuleName(tc.remote)
			got, err := mapScheduledActionsForRuleToState(
				ctx,
				[]graphmodels.DeviceComplianceScheduledActionForRuleable{action},
				prior,
			)
			require.NoError(t, err)
			require.Equal(
				t,
				types.StringValue(tc.want),
				got.Elements()[0].(types.Object).Attributes()["rule_name"],
			)
		})
	}
	t.Run("import", func(t *testing.T) {
		action.SetRuleName(nil)
		got, err := mapScheduledActionsForRuleToState(
			ctx,
			[]graphmodels.DeviceComplianceScheduledActionForRuleable{action},
			types.ListNull(elementType),
		)
		require.NoError(t, err)
		require.Equal(
			t,
			types.StringValue("PasswordRequired"),
			got.Elements()[0].(types.Object).Attributes()["rule_name"],
		)
	})
}

// TestUnitMapScheduledActionConfigurationsToState verifies list typing and values
// survive nested set conversion, including non-empty notification recipients.
func TestUnitMapScheduledActionConfigurationsToState(t *testing.T) {
	ctx := context.Background()
	action := graphmodels.NewDeviceComplianceActionItem()
	action.SetActionType(new(graphmodels.NOTIFICATION_DEVICECOMPLIANCEACTIONTYPE))
	action.SetGracePeriodHours(new(int32(24)))
	action.SetNotificationTemplateId(new("11111111-1111-1111-1111-111111111111"))
	for _, recipients := range [][]string{nil, {}, {"22222222-2222-2222-2222-222222222222", "33333333-3333-3333-3333-333333333333"}} {
		action.SetNotificationMessageCCList(recipients)
		got, err := mapScheduledActionConfigurationsToState(
			ctx,
			[]graphmodels.DeviceComplianceActionItemable{action},
			nil,
		)
		require.NoError(t, err)
		var models []ScheduledActionConfigurationModel
		require.False(t, got.ElementsAs(ctx, &models, false).HasError())
		require.Len(t, models, 1)
		require.Equal(t, "notification", models[0].ActionType.ValueString())
		require.Equal(t, int32(24), models[0].GracePeriodHours.ValueInt32())
		require.Len(t, models[0].NotificationMessageCcList.Elements(), len(recipients))
		for i, recipient := range recipients {
			require.Equal(
				t,
				types.StringValue(recipient),
				models[0].NotificationMessageCcList.Elements()[i],
			)
		}
	}
}
