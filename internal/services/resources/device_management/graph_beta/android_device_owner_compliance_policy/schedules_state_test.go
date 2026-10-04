package graphBetaAndroidDeviceOwnerCompliancePolicy

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
	"github.com/stretchr/testify/require"
)

func TestUnitAndroidComplianceScheduledRuleNames(t *testing.T) {
	ctx := context.Background()
	rule := graphmodels.NewDeviceComplianceScheduledActionForRule()
	rule.SetScheduledActionConfigurations([]graphmodels.DeviceComplianceActionItemable{})
	policy := graphmodels.NewAndroidDeviceOwnerCompliancePolicy()
	policy.SetScheduledActionsForRule([]graphmodels.DeviceComplianceScheduledActionForRuleable{rule})
	data := DeviceCompliancePolicyResourceModel{}
	MapRemoteStateToTerraform(ctx, &data, policy)
	elementType := data.ScheduledActionsForRule.ElementType(ctx).(types.ObjectType)
	initial := data.ScheduledActionsForRule
	for _, tc := range []struct {
		name   string
		prior  types.String
		remote *string
		want   string
	}{
		{name: "null", prior: types.StringNull(), want: "PasswordRequired"},
		{name: "unknown", prior: types.StringUnknown(), want: "PasswordRequired"},
		{name: "explicit", prior: types.StringValue("PasswordRequired"), want: "PasswordRequired"},
		{name: "legacy", prior: types.StringValue("unavailable"), want: "unavailable"},
		{name: "empty_remote", prior: types.StringValue("unavailable"), remote: new(""), want: "unavailable"},
		{name: "authoritative_remote", prior: types.StringValue("unavailable"), remote: new("PasswordRequired"), want: "PasswordRequired"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			attrs := initial.Elements()[0].(types.Object).Attributes()
			attrs["rule_name"] = tc.prior
			data.ScheduledActionsForRule = types.ListValueMust(elementType, []attr.Value{types.ObjectValueMust(elementType.AttrTypes, attrs)})
			rule.SetRuleName(tc.remote)
			MapRemoteStateToTerraform(ctx, &data, policy)
			var rules []ScheduledActionForRuleModel
			require.False(t, data.ScheduledActionsForRule.ElementsAs(ctx, &rules, false).HasError())
			require.Equal(t, types.StringValue(tc.want), rules[0].RuleName)
		})
	}
}

func TestUnitAndroidComplianceScheduledActionDrift(t *testing.T) {
	ctx := context.Background()
	const zeroGUID = "00000000-0000-0000-0000-000000000000"
	const realGUID = "11111111-1111-1111-1111-111111111111"
	action := func(hours int32) *graphmodels.DeviceComplianceActionItem {
		item := graphmodels.NewDeviceComplianceActionItem()
		item.SetActionType(new(graphmodels.NOTIFICATION_DEVICECOMPLIANCEACTIONTYPE))
		item.SetGracePeriodHours(&hours)
		item.SetNotificationTemplateId(new(zeroGUID))
		item.SetNotificationMessageCCList([]string{})
		return item
	}
	first, second := action(24), action(72)
	rule := graphmodels.NewDeviceComplianceScheduledActionForRule()
	rule.SetScheduledActionConfigurations([]graphmodels.DeviceComplianceActionItemable{first, second})
	policy := graphmodels.NewAndroidDeviceOwnerCompliancePolicy()
	policy.SetScheduledActionsForRule([]graphmodels.DeviceComplianceScheduledActionForRuleable{rule})
	data := DeviceCompliancePolicyResourceModel{}
	MapRemoteStateToTerraform(ctx, &data, policy)
	var rules []ScheduledActionForRuleModel
	require.False(t, data.ScheduledActionsForRule.ElementsAs(ctx, &rules, false).HasError())
	var configs []ScheduledActionConfigurationModel
	require.False(t, rules[0].ScheduledActionConfigurations.ElementsAs(ctx, &configs, false).HasError())
	for i := range configs {
		if configs[i].GracePeriodHours.ValueInt32() == 24 {
			configs[i].NotificationTemplateId = types.StringValue("")
		}
	}
	priorConfigs, diags := types.SetValueFrom(ctx, rules[0].ScheduledActionConfigurations.ElementType(ctx), configs)
	require.False(t, diags.HasError())
	rules[0].ScheduledActionConfigurations = priorConfigs
	prior, diags := types.ListValueFrom(ctx, data.ScheduledActionsForRule.ElementType(ctx), rules)
	require.False(t, diags.HasError())

	for _, tc := range []struct {
		name   string
		mutate func(*graphmodels.DeviceComplianceActionItem)
		drift  bool
	}{
		{name: "reordered_set_preserves_only_equivalent_empty_template"},
		{name: "grace_period", drift: true, mutate: func(item *graphmodels.DeviceComplianceActionItem) { item.SetGracePeriodHours(new(int32(48))) }},
		{name: "action_type", drift: true, mutate: func(item *graphmodels.DeviceComplianceActionItem) {
			item.SetActionType(new(graphmodels.BLOCK_DEVICECOMPLIANCEACTIONTYPE))
		}},
		{name: "template", drift: true, mutate: func(item *graphmodels.DeviceComplianceActionItem) { item.SetNotificationTemplateId(new(realGUID)) }},
		{name: "recipients", drift: true, mutate: func(item *graphmodels.DeviceComplianceActionItem) {
			item.SetNotificationMessageCCList([]string{realGUID})
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			remote := action(24)
			if tc.mutate != nil {
				tc.mutate(remote)
			}
			rule.SetScheduledActionConfigurations([]graphmodels.DeviceComplianceActionItemable{second, remote})
			data.ScheduledActionsForRule = prior
			MapRemoteStateToTerraform(ctx, &data, policy)
			require.Equal(t, !tc.drift, prior.Equal(data.ScheduledActionsForRule), "only equivalent representations should preserve state")
			var actualRules []ScheduledActionForRuleModel
			require.False(t, data.ScheduledActionsForRule.ElementsAs(ctx, &actualRules, false).HasError())
			var actual []ScheduledActionConfigurationModel
			require.False(t, actualRules[0].ScheduledActionConfigurations.ElementsAs(ctx, &actual, false).HasError())
			require.Len(t, actual, 2)
			if tc.drift {
				found := false
				for _, item := range actual {
					if item.GracePeriodHours.ValueInt32() == *remote.GetGracePeriodHours() && item.ActionType.ValueString() == remote.GetActionType().String() {
						require.Equal(t, *remote.GetNotificationTemplateId(), item.NotificationTemplateId.ValueString())
						require.Len(t, item.NotificationMessageCcList.Elements(), len(remote.GetNotificationMessageCCList()))
						found = true
					}
				}
				require.True(t, found, "remote action values must be reflected")
			}
		})
	}
}

// A failed conversion must surface a diagnostic and preserve the saved
// schedule, rather than silently replacing it with remote or default values.
func TestUnitAndroidComplianceInvalidScheduledActionState(t *testing.T) {
	ctx := context.Background()
	invalidConfigurations := types.SetValueMust(types.StringType, []attr.Value{types.StringValue("invalid action")})
	ruleType := types.ObjectType{AttrTypes: map[string]attr.Type{
		"rule_name":                       types.StringType,
		"scheduled_action_configurations": invalidConfigurations.Type(ctx),
	}}
	invalidRule := types.ObjectValueMust(ruleType.AttrTypes, map[string]attr.Value{
		"rule_name":                       types.StringValue("PasswordRequired"),
		"scheduled_action_configurations": invalidConfigurations,
	})
	for _, tc := range []struct {
		name   string
		prior  types.List
		detail string
	}{
		{
			name:   "invalid_rule_state",
			prior:  types.ListValueMust(types.StringType, []attr.Value{types.StringValue("invalid rule")}),
			detail: "failed to read prior scheduled actions",
		},
		{
			name:   "invalid_action_state",
			prior:  types.ListValueMust(ruleType, []attr.Value{invalidRule}),
			detail: "failed to read prior action configurations",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			policy := graphmodels.NewAndroidDeviceOwnerCompliancePolicy()
			policy.SetScheduledActionsForRule([]graphmodels.DeviceComplianceScheduledActionForRuleable{
				graphmodels.NewDeviceComplianceScheduledActionForRule(),
			})
			data := DeviceCompliancePolicyResourceModel{ScheduledActionsForRule: tc.prior}
			diags := MapRemoteStateToTerraform(ctx, &data, policy)
			require.Len(t, diags.Errors(), 1)
			require.Equal(t, "Error mapping scheduled actions for rule", diags.Errors()[0].Summary())
			require.Contains(t, diags.Errors()[0].Detail(), tc.detail)
			require.True(t, tc.prior.Equal(data.ScheduledActionsForRule), "failed conversion must preserve saved schedules")
		})
	}
}
