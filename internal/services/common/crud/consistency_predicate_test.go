package crud

import (
	"context"
	"slices"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var predicateSettingsAttrTypes = map[string]attr.Type{
	"mode":    types.StringType,
	"enabled": types.BoolType,
}

type predicateTestModel struct {
	ID       types.String `tfsdk:"id"`
	Name     types.String `tfsdk:"name"`
	Count    types.Int64  `tfsdk:"count"`
	Tags     types.Set    `tfsdk:"tags"`
	Items    types.List   `tfsdk:"items"`
	Labels   types.Map    `tfsdk:"labels"`
	Settings types.Object `tfsdk:"settings"`
	Rules    types.Set    `tfsdk:"rules"`
}

func predicateTestSchema() schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id":     schema.StringAttribute{Computed: true},
			"name":   schema.StringAttribute{Optional: true},
			"count":  schema.Int64Attribute{Optional: true},
			"tags":   schema.SetAttribute{Optional: true, ElementType: types.StringType},
			"items":  schema.ListAttribute{Optional: true, ElementType: types.StringType},
			"labels": schema.MapAttribute{Optional: true, ElementType: types.StringType},
			"settings": schema.SingleNestedAttribute{Optional: true, Attributes: map[string]schema.Attribute{
				"mode":    schema.StringAttribute{Optional: true},
				"enabled": schema.BoolAttribute{Optional: true},
			}},
			"rules": schema.SetNestedAttribute{Optional: true, NestedObject: schema.NestedAttributeObject{Attributes: map[string]schema.Attribute{
				"mode":    schema.StringAttribute{Optional: true},
				"enabled": schema.BoolAttribute{Optional: true},
			}}},
		},
	}
}

// baseModel returns a fully known model that the tests change one attribute at a time.
func baseModel(t *testing.T) predicateTestModel {
	t.Helper()
	return predicateTestModel{
		ID:       types.StringValue("00000000-0000-0000-0000-000000000001"),
		Name:     types.StringValue("Example"),
		Count:    types.Int64Value(2),
		Tags:     stringSetValue(t, "a", "b", "c"),
		Items:    stringListValue(t, "x", "y"),
		Labels:   types.MapValueMust(types.StringType, map[string]attr.Value{"k": types.StringValue("v")}),
		Settings: settingsValue(t, types.StringValue("block"), types.BoolValue(true)),
		Rules:    rulesValue(t, settingsValue(t, types.StringValue("allow"), types.BoolValue(true)), settingsValue(t, types.StringValue("deny"), types.BoolValue(false))),
	}
}

func stringSetValue(t *testing.T, values ...string) types.Set {
	t.Helper()
	elems := make([]attr.Value, 0, len(values))
	for _, v := range values {
		elems = append(elems, types.StringValue(v))
	}
	s, d := types.SetValue(types.StringType, elems)
	require.False(t, d.HasError())
	return s
}

func stringListValue(t *testing.T, values ...string) types.List {
	t.Helper()
	elems := make([]attr.Value, 0, len(values))
	for _, v := range values {
		elems = append(elems, types.StringValue(v))
	}
	l, d := types.ListValue(types.StringType, elems)
	require.False(t, d.HasError())
	return l
}

func settingsValue(t *testing.T, mode types.String, enabled types.Bool) types.Object {
	t.Helper()
	o, d := types.ObjectValue(predicateSettingsAttrTypes, map[string]attr.Value{"mode": mode, "enabled": enabled})
	require.False(t, d.HasError())
	return o
}

func rulesValue(t *testing.T, rules ...attr.Value) types.Set {
	t.Helper()
	s, d := types.SetValue(types.ObjectType{AttrTypes: predicateSettingsAttrTypes}, rules)
	require.False(t, d.HasError())
	return s
}

func stateFromModel(t *testing.T, m predicateTestModel) tfsdk.State {
	t.Helper()
	s := predicateTestSchema()
	state := tfsdk.State{Schema: s, Raw: tftypes.NewValue(s.Type().TerraformType(context.Background()), nil)}
	require.False(t, state.Set(context.Background(), m).HasError())
	return state
}

func planFromModel(t *testing.T, m predicateTestModel) tfsdk.Plan {
	t.Helper()
	state := stateFromModel(t, m)
	return tfsdk.Plan{Schema: state.Schema, Raw: state.Raw}
}

var allPredicatePaths = RootPaths("id", "name", "count", "tags", "items", "labels", "settings", "rules")

func TestRequireStableReads(t *testing.T) {
	reads := []bool{true, false, true, true, true}
	i := 0
	pred := RequireStableReads(func(context.Context, tfsdk.State) bool {
		ok := reads[i]
		i++
		return ok
	})
	var got []bool
	for range reads {
		got = append(got, pred(context.Background(), tfsdk.State{}))
	}
	// A stale read resets the count; only the third match in a row passes.
	want := []bool{false, false, false, false, true}
	if !slices.Equal(got, want) {
		t.Fatalf("results = %v, expected %v", got, want)
	}
}

func TestRequireConsecutiveReads(t *testing.T) {
	reads := []bool{true, true, true, true, false, true, true, true, true, true}
	i := 0
	pred := RequireConsecutiveReads(5, func(context.Context, tfsdk.State) bool {
		ok := reads[i]
		i++
		return ok
	})
	var got []bool
	for range reads {
		got = append(got, pred(context.Background(), tfsdk.State{}))
	}
	want := []bool{false, false, false, false, false, false, false, false, false, true}
	if !slices.Equal(got, want) {
		t.Fatalf("results = %v, expected %v", got, want)
	}
}

func TestRootPaths(t *testing.T) {
	assert.Equal(t, []path.Path{path.Root("a"), path.Root("b")}, RootPaths("a", "b"))
	assert.Empty(t, RootPaths())
}

func TestAttributesMatchPredicate(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name     string
		opts     AttributeMatchOptions
		expected func(*predicateTestModel)
		actual   func(*predicateTestModel)
		want     bool
	}{
		{name: "identical values match", want: true},
		{name: "changed string does not match", actual: func(m *predicateTestModel) { m.Name = types.StringValue("Stale") }},
		{name: "changed int does not match", actual: func(m *predicateTestModel) { m.Count = types.Int64Value(1) }},
		{name: "set in another order matches", actual: func(m *predicateTestModel) { m.Tags = stringSetValue(t, "c", "a", "b") }, want: true},
		{name: "set with another element does not match", actual: func(m *predicateTestModel) { m.Tags = stringSetValue(t, "a", "b", "d") }},
		{name: "set with fewer elements does not match", actual: func(m *predicateTestModel) { m.Tags = stringSetValue(t, "a", "b") }},
		{name: "nested set in another order matches", actual: func(m *predicateTestModel) {
			m.Rules = rulesValue(t, settingsValue(t, types.StringValue("deny"), types.BoolValue(false)), settingsValue(t, types.StringValue("allow"), types.BoolValue(true)))
		}, want: true},
		{name: "list in another order does not match", actual: func(m *predicateTestModel) { m.Items = stringListValue(t, "y", "x") }},
		{name: "list in another order matches when unordered", opts: AttributeMatchOptions{UnorderedLists: true}, actual: func(m *predicateTestModel) { m.Items = stringListValue(t, "y", "x") }, want: true},
		{name: "unordered list with another element does not match", opts: AttributeMatchOptions{UnorderedLists: true}, actual: func(m *predicateTestModel) { m.Items = stringListValue(t, "y", "z") }},
		{name: "list with fewer elements does not match", actual: func(m *predicateTestModel) { m.Items = stringListValue(t, "x") }},
		{name: "changed map value does not match", actual: func(m *predicateTestModel) {
			m.Labels = types.MapValueMust(types.StringType, map[string]attr.Value{"k": types.StringValue("w")})
		}},
		{name: "map with another key does not match", actual: func(m *predicateTestModel) {
			m.Labels = types.MapValueMust(types.StringType, map[string]attr.Value{"j": types.StringValue("v")})
		}},
		{name: "map with extra key does not match", actual: func(m *predicateTestModel) {
			m.Labels = types.MapValueMust(types.StringType, map[string]attr.Value{"k": types.StringValue("v"), "j": types.StringValue("v")})
		}},
		{name: "changed nested attribute does not match", actual: func(m *predicateTestModel) { m.Settings = settingsValue(t, types.StringValue("block"), types.BoolValue(false)) }},
		{name: "unknown expected value is skipped", expected: func(m *predicateTestModel) { m.Name = types.StringUnknown() }, actual: func(m *predicateTestModel) { m.Name = types.StringValue("Graph default") }, want: true},
		{name: "unknown expected nested attribute is skipped", expected: func(m *predicateTestModel) { m.Settings = settingsValue(t, types.StringValue("block"), types.BoolUnknown()) }, actual: func(m *predicateTestModel) {
			m.Settings = settingsValue(t, types.StringValue("block"), types.BoolValue(false))
		}, want: true},
		{name: "unknown expected set element matches any remaining element", expected: func(m *predicateTestModel) {
			m.Tags = types.SetValueMust(types.StringType, []attr.Value{types.StringUnknown(), types.StringValue("a"), types.StringValue("b")})
		}, actual: func(m *predicateTestModel) { m.Tags = stringSetValue(t, "a", "b", "z") }, want: true},
		{name: "unknown expected set element still needs the same count", expected: func(m *predicateTestModel) {
			m.Tags = types.SetValueMust(types.StringType, []attr.Value{types.StringUnknown(), types.StringValue("a")})
		}, actual: func(m *predicateTestModel) { m.Tags = stringSetValue(t, "a") }},
		{name: "unknown expected object fields do not claim an element another needs", expected: func(m *predicateTestModel) {
			m.Rules = rulesValue(t, settingsValue(t, types.StringUnknown(), types.BoolValue(true)), settingsValue(t, types.StringValue("allow"), types.BoolValue(true)))
		}, actual: func(m *predicateTestModel) {
			m.Rules = rulesValue(t, settingsValue(t, types.StringValue("allow"), types.BoolValue(true)), settingsValue(t, types.StringValue("other"), types.BoolValue(true)))
		}, want: true},
		{name: "null expected must read back null", expected: func(m *predicateTestModel) { m.Name = types.StringNull() }, actual: func(m *predicateTestModel) { m.Name = types.StringValue("Graph default") }},
		{name: "null expected matches null", expected: func(m *predicateTestModel) { m.Name = types.StringNull() }, actual: func(m *predicateTestModel) { m.Name = types.StringNull() }, want: true},
		{name: "null expected is skipped when ignored", opts: AttributeMatchOptions{IgnoreNullExpected: true}, expected: func(m *predicateTestModel) { m.Name = types.StringNull() }, actual: func(m *predicateTestModel) {
			m.Name = types.StringValue("Graph default")
		}, want: true},
		{name: "null expected nested attribute is skipped when ignored", opts: AttributeMatchOptions{IgnoreNullExpected: true}, expected: func(m *predicateTestModel) {
			m.Settings = settingsValue(t, types.StringNull(), types.BoolValue(true))
		}, actual: func(m *predicateTestModel) { m.Settings = settingsValue(t, types.StringValue("block"), types.BoolValue(true)) }, want: true},
		{name: "null read-back does not match a known value", actual: func(m *predicateTestModel) { m.Tags = types.SetNull(types.StringType) }},
		{name: "null read-back does not match a value when null expected is ignored", opts: AttributeMatchOptions{IgnoreNullExpected: true}, actual: func(m *predicateTestModel) { m.Name = types.StringNull() }},
		{name: "null set does not match empty set by default", expected: func(m *predicateTestModel) { m.Tags = types.SetNull(types.StringType) }, actual: func(m *predicateTestModel) { m.Tags = stringSetValue(t) }},
		{name: "null set matches empty set when enabled", opts: AttributeMatchOptions{NullMatchesEmpty: true}, expected: func(m *predicateTestModel) { m.Tags = types.SetNull(types.StringType) }, actual: func(m *predicateTestModel) {
			m.Tags = stringSetValue(t)
		}, want: true},
		{name: "empty list matches null list when enabled", opts: AttributeMatchOptions{NullMatchesEmpty: true}, expected: func(m *predicateTestModel) { m.Items = stringListValue(t) }, actual: func(m *predicateTestModel) {
			m.Items = types.ListNull(types.StringType)
		}, want: true},
		{name: "empty map matches null map when enabled", opts: AttributeMatchOptions{NullMatchesEmpty: true}, expected: func(m *predicateTestModel) { m.Labels = types.MapValueMust(types.StringType, map[string]attr.Value{}) }, actual: func(m *predicateTestModel) {
			m.Labels = types.MapNull(types.StringType)
		}, want: true},
		{name: "empty string matches null string when enabled", opts: AttributeMatchOptions{NullMatchesEmpty: true}, expected: func(m *predicateTestModel) { m.Name = types.StringValue("") }, actual: func(m *predicateTestModel) {
			m.Name = types.StringNull()
		}, want: true},
		{name: "null set does not match non-empty set when enabled", opts: AttributeMatchOptions{NullMatchesEmpty: true}, expected: func(m *predicateTestModel) { m.Tags = types.SetNull(types.StringType) }, actual: func(m *predicateTestModel) {
			m.Tags = stringSetValue(t, "a")
		}},
		{name: "null int does not match zero when enabled", opts: AttributeMatchOptions{NullMatchesEmpty: true}, expected: func(m *predicateTestModel) { m.Count = types.Int64Null() }, actual: func(m *predicateTestModel) { m.Count = types.Int64Value(0) }},
		{name: "null object does not match an object when enabled", opts: AttributeMatchOptions{NullMatchesEmpty: true}, expected: func(m *predicateTestModel) { m.Settings = types.ObjectNull(predicateSettingsAttrTypes) }, actual: func(m *predicateTestModel) {
			m.Settings = settingsValue(t, types.StringNull(), types.BoolNull())
		}},
		{name: "string case differs by default", actual: func(m *predicateTestModel) { m.ID = types.StringValue("00000000-0000-0000-0000-00000000000A") }, expected: func(m *predicateTestModel) {
			m.ID = types.StringValue("00000000-0000-0000-0000-00000000000a")
		}},
		{name: "string case ignored when enabled", opts: AttributeMatchOptions{CaseInsensitiveStrings: true}, expected: func(m *predicateTestModel) {
			m.Tags = stringSetValue(t, "A", "B", "C")
		}, want: true},
		{name: "case-insensitive strings still compare values", opts: AttributeMatchOptions{CaseInsensitiveStrings: true}, actual: func(m *predicateTestModel) { m.Name = types.StringValue("Other") }},
		{name: "required attribute populated", opts: AttributeMatchOptions{RequireKnown: RootPaths("id", "count")}, want: true},
		{name: "required attribute empty string", opts: AttributeMatchOptions{RequireKnown: RootPaths("id")}, actual: func(m *predicateTestModel) { m.ID = types.StringValue("") }, expected: func(m *predicateTestModel) {
			m.ID = types.StringValue("")
		}},
		{name: "required attribute null", opts: AttributeMatchOptions{RequireKnown: RootPaths("count")}, actual: func(m *predicateTestModel) { m.Count = types.Int64Null() }, expected: func(m *predicateTestModel) {
			m.Count = types.Int64Unknown()
		}},
		{name: "required attribute unknown", opts: AttributeMatchOptions{RequireKnown: RootPaths("id")}, actual: func(m *predicateTestModel) { m.ID = types.StringUnknown() }, expected: func(m *predicateTestModel) {
			m.ID = types.StringUnknown()
		}},
		{name: "unknown read-back never matches", actual: func(m *predicateTestModel) { m.Name = types.StringUnknown() }},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expected := baseModel(t)
			if tt.expected != nil {
				tt.expected(&expected)
			}
			actual := baseModel(t)
			if tt.actual != nil {
				tt.actual(&actual)
			}
			opts := tt.opts
			if opts.Paths == nil {
				opts.Paths = allPredicatePaths
			}
			pred := AttributesMatchPredicate(planFromModel(t, expected), opts)
			assert.Equal(t, tt.want, pred(ctx, stateFromModel(t, actual)))
		})
	}
}

func TestAttributesMatchPredicateNullState(t *testing.T) {
	s := predicateTestSchema()
	state := tfsdk.State{Schema: s, Raw: tftypes.NewValue(s.Type().TerraformType(context.Background()), nil)}
	pred := AttributesMatchPredicate(planFromModel(t, baseModel(t)), AttributeMatchOptions{Paths: allPredicatePaths})
	assert.False(t, pred(context.Background(), state))
}

func TestAttributesMatchPredicateExpectedFromState(t *testing.T) {
	m := baseModel(t)
	pred := AttributesMatchPredicate(stateFromModel(t, m), AttributeMatchOptions{Paths: allPredicatePaths})
	assert.True(t, pred(context.Background(), stateFromModel(t, m)))
}

func TestAttributesMatchPredicateNestedPath(t *testing.T) {
	expected := baseModel(t)
	actual := baseModel(t)
	actual.Settings = settingsValue(t, types.StringValue("block"), types.BoolValue(false))
	pred := AttributesMatchPredicate(planFromModel(t, expected), AttributeMatchOptions{Paths: []path.Path{path.Root("settings").AtName("mode")}})
	assert.True(t, pred(context.Background(), stateFromModel(t, actual)))
	pred = AttributesMatchPredicate(planFromModel(t, expected), AttributeMatchOptions{Paths: []path.Path{path.Root("settings").AtName("enabled")}})
	assert.False(t, pred(context.Background(), stateFromModel(t, actual)))
}

func TestAttributesMatchPredicateInvalidPath(t *testing.T) {
	m := baseModel(t)
	pred := AttributesMatchPredicate(planFromModel(t, m), AttributeMatchOptions{Paths: RootPaths("missing")})
	assert.False(t, pred(context.Background(), stateFromModel(t, m)), "a path missing from the expected schema never matches")

	pred = AttributesMatchPredicate(planFromModel(t, m), AttributeMatchOptions{RequireKnown: RootPaths("missing")})
	assert.False(t, pred(context.Background(), stateFromModel(t, m)), "a required path missing from the read-back schema never matches")

	other := schema.Schema{Attributes: map[string]schema.Attribute{"other": schema.StringAttribute{Optional: true}}}
	otherState := tfsdk.State{Schema: other, Raw: tftypes.NewValue(other.Type().TerraformType(context.Background()), map[string]tftypes.Value{"other": tftypes.NewValue(tftypes.String, "x")})}
	pred = AttributesMatchPredicate(planFromModel(t, m), AttributeMatchOptions{Paths: RootPaths("name")})
	assert.False(t, pred(context.Background(), otherState), "a path missing from the read-back schema never matches")
}

func TestValueMatchesTypeMismatch(t *testing.T) {
	ctx := context.Background()
	opts := AttributeMatchOptions{CaseInsensitiveStrings: true}
	assert.False(t, opts.valueMatches(ctx, settingsValue(t, types.StringValue("a"), types.BoolValue(true)), types.StringValue("a")))
	assert.False(t, opts.valueMatches(ctx, stringSetValue(t, "a"), stringListValue(t, "a")))
	assert.False(t, opts.valueMatches(ctx, stringListValue(t, "a"), stringSetValue(t, "a")))
	assert.False(t, opts.valueMatches(ctx, types.MapValueMust(types.StringType, map[string]attr.Value{}), stringSetValue(t)))
	assert.False(t, opts.valueMatches(ctx, types.StringValue("1"), types.Int64Value(1)))
	assert.False(t, opts.valueMatches(ctx, types.StringValue("a"), nil))
	assert.True(t, opts.valueMatches(ctx, nil, types.StringValue("a")))
	assert.True(t, isPopulated(ctx, types.BoolValue(false)))
	assert.False(t, isPopulated(ctx, nil))
}
