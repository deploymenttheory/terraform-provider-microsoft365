package sharedmodels

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/stretchr/testify/require"
)

func stringSet(t *testing.T, values ...string) types.Set {
	t.Helper()
	elems := make([]attr.Value, 0, len(values))
	for _, v := range values {
		elems = append(elems, types.StringValue(v))
	}
	s, d := types.SetValue(types.StringType, elems)
	require.False(t, d.HasError())
	return s
}

func rulesObject(t *testing.T, m AssignmentRulesModel) types.Object {
	t.Helper()
	o, d := types.ObjectValueFrom(context.Background(), AssignmentRulesAttrTypes, m)
	require.False(t, d.HasError())
	return o
}

func TestSettingsPatchEnumeratedSendsAllMemberLists(t *testing.T) {
	ctx := context.Background()
	null := types.SetNull(types.StringType)
	plan := ForwardingProfileSettings{
		State:                    types.StringValue("enabled"),
		ClientFallbackAction:     types.StringUnknown(),
		ForwardingSourceMatching: types.ObjectUnknown(ForwardingSourceMatchingAttrTypes),
		LocalNetworkSettings:     types.ObjectUnknown(LocalNetworkSettingsAttrTypes),
		AssignmentRules: rulesObject(t, AssignmentRulesModel{
			Membership: types.StringValue(MembershipEnumerated),
			UserIDs:    null,
			GroupIDs:   stringSet(t, "b", "a"),
			DeviceIDs:  null,
			Platforms:  stringSet(t, "iOS", "windows", "macOS"),
		}),
	}
	body := map[string]any{}
	require.False(t, SettingsPatch(ctx, body, plan, nil).HasError())

	raw, err := json.Marshal(body)
	require.NoError(t, err)
	require.JSONEq(t, `{
	  "state": "enabled",
	  "assignmentRules": {
	    "assignments": {
	      "@odata.type": "#microsoft.graph.networkaccess.forwardingProfileAssignmentsEnumeratedMembership",
	      "users": [],
	      "groups": [{"id": "a"}, {"id": "b"}],
	      "devices": []
	    },
	    "devicePropertyFilters": {"platforms": "windows,macOS,iOS"}
	  }
	}`, string(raw))
}

func TestSettingsPatchSendsOnlyChangedValues(t *testing.T) {
	ctx := context.Background()
	source := func(v string) types.Object {
		o, d := types.ObjectValueFrom(
			ctx,
			ForwardingSourceMatchingAttrTypes,
			ForwardingSourceMatchingModel{SessionType: types.StringValue(v)},
		)
		require.False(t, d.HasError())
		return o
	}
	local, d := types.ObjectValueFrom(
		ctx,
		LocalNetworkSettingsAttrTypes,
		LocalNetworkSettingsModel{AllowLocalSubnetBypass: types.BoolValue(true)},
	)
	require.False(t, d.HasError())
	rules := rulesObject(t, AssignmentRulesModel{
		Membership: types.StringValue(MembershipAll),
		UserIDs:    types.SetNull(types.StringType),
		GroupIDs:   types.SetNull(types.StringType),
		DeviceIDs:  types.SetNull(types.StringType),
		Platforms:  stringSet(t, "windows"),
	})
	state := ForwardingProfileSettings{
		State:                    types.StringValue("disabled"),
		ClientFallbackAction:     types.StringValue("bypass"),
		ForwardingSourceMatching: source("any"),
		LocalNetworkSettings:     local,
		AssignmentRules:          rules,
	}
	plan := state
	plan.ForwardingSourceMatching = source("agent")

	body := map[string]any{}
	require.False(t, SettingsPatch(ctx, body, plan, &state).HasError())
	require.Equal(
		t,
		map[string]any{"forwardingSourceMatching": map[string]any{"sessionType": "agent"}},
		body,
	)
}

func TestSettingsFromWire(t *testing.T) {
	ctx := context.Background()
	raw := `{
	  "state": "disabled",
	  "clientFallbackAction": "bypass",
	  "localNetworkSettings": null,
	  "forwardingSourceMatching": {"sessionType": "any"},
	  "assignmentRules": {
	    "assignments": {
	      "@odata.type": "#microsoft.graph.networkaccess.forwardingProfileAssignmentsEnumeratedMembership",
	      "users": [{"id": "u1"}], "groups": [], "devices": [{"id": "d1"}]
	    },
	    "devicePropertyFilters": {"platforms": "Windows, macos,IOS"}
	  }
	}`
	var w ForwardingProfileWire
	require.NoError(t, json.Unmarshal([]byte(raw), &w))
	s, d := SettingsFromWire(ctx, &w)
	require.False(t, d.HasError())
	require.True(t, s.LocalNetworkSettings.IsNull())

	var rules AssignmentRulesModel
	require.False(t, s.AssignmentRules.As(ctx, &rules, basetypes.ObjectAsOptions{}).HasError())
	require.Equal(t, MembershipEnumerated, rules.Membership.ValueString())
	require.True(t, rules.GroupIDs.IsNull(), "empty member lists map to null")
	require.True(t, rules.UserIDs.Equal(stringSet(t, "u1")))
	require.True(t, rules.DeviceIDs.Equal(stringSet(t, "d1")))
	require.True(
		t,
		rules.Platforms.Equal(stringSet(t, "windows", "macOS", "iOS")),
		"platforms are canonicalized",
	)
}

func TestSettingsFromWirePlatformsNone(t *testing.T) {
	ctx := context.Background()
	raw := `{"assignmentRules": {
	  "assignments": {"@odata.type": "#microsoft.graph.networkaccess.forwardingProfileAssignmentsNoMembership"},
	  "devicePropertyFilters": {"platforms": "none"}
	}}`
	var w ForwardingProfileWire
	require.NoError(t, json.Unmarshal([]byte(raw), &w))
	s, d := SettingsFromWire(ctx, &w)
	require.False(t, d.HasError())
	var rules AssignmentRulesModel
	require.False(t, s.AssignmentRules.As(ctx, &rules, basetypes.ObjectAsOptions{}).HasError())
	require.Equal(t, MembershipNone, rules.Membership.ValueString())
	require.Empty(t, rules.Platforms.Elements())
}

func TestSettingsFromWireUnknownAssignmentType(t *testing.T) {
	w := ForwardingProfileWire{AssignmentRules: &ForwardingProfileAssignmentsRW{
		Assignments: &ForwardingProfileAssignmentsWire{
			ODataType: "#microsoft.graph.networkaccess.somethingNew",
		},
	}}
	_, d := SettingsFromWire(context.Background(), &w)
	require.True(t, d.HasError())
}
