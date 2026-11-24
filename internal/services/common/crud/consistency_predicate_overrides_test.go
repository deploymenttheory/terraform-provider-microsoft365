package crud

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
)

func TestAttributesMatchPredicateOverrides(t *testing.T) {
	ctx := context.Background()
	upperRules := func(m *predicateTestModel) {
		m.Rules = rulesValue(t, settingsValue(t, types.StringValue("ALLOW"), types.BoolValue(true)), settingsValue(t, types.StringValue("DENY"), types.BoolValue(false)))
	}
	nameCaseInsensitive := []AttributeMatchOverride{{Path: path.MatchRoot("name"), CaseInsensitiveStrings: new(true)}}
	nameCaseSensitive := []AttributeMatchOverride{{Path: path.MatchRoot("name"), CaseInsensitiveStrings: new(false)}}

	tests := []struct {
		name     string
		opts     AttributeMatchOptions
		expected func(*predicateTestModel)
		actual   func(*predicateTestModel)
		want     bool
	}{
		{name: "root override ignores case on its attribute", opts: AttributeMatchOptions{Overrides: nameCaseInsensitive}, actual: func(m *predicateTestModel) { m.Name = types.StringValue("EXAMPLE") }, want: true},
		{name: "root override leaves other attributes case-sensitive", opts: AttributeMatchOptions{Overrides: nameCaseInsensitive}, actual: func(m *predicateTestModel) { m.Tags = stringSetValue(t, "A", "B", "C") }},
		{name: "root override takes precedence over the global option", opts: AttributeMatchOptions{CaseInsensitiveStrings: true, Overrides: nameCaseSensitive}, actual: func(m *predicateTestModel) { m.Name = types.StringValue("EXAMPLE") }},
		{name: "global option still applies outside the overridden attribute", opts: AttributeMatchOptions{CaseInsensitiveStrings: true, Overrides: nameCaseSensitive}, actual: func(m *predicateTestModel) {
			m.Tags = stringSetValue(t, "A", "B", "C")
		}, want: true},
		{name: "later override wins", opts: AttributeMatchOptions{Overrides: append(append([]AttributeMatchOverride{}, nameCaseInsensitive...), nameCaseSensitive...)}, actual: func(m *predicateTestModel) {
			m.Name = types.StringValue("EXAMPLE")
		}},
		{name: "override without the setting keeps the global option", opts: AttributeMatchOptions{CaseInsensitiveStrings: true, Overrides: []AttributeMatchOverride{{Path: path.MatchRoot("name"), UnorderedLists: new(true)}}}, actual: func(m *predicateTestModel) {
			m.Name = types.StringValue("EXAMPLE")
		}, want: true},
		{name: "nested object override ignores case on its field", opts: AttributeMatchOptions{Overrides: []AttributeMatchOverride{{Path: path.MatchRoot("settings").AtName("mode"), CaseInsensitiveStrings: new(true)}}}, actual: func(m *predicateTestModel) {
			m.Settings = settingsValue(t, types.StringValue("BLOCK"), types.BoolValue(true))
		}, want: true},
		{name: "nested object override does not reach another attribute", opts: AttributeMatchOptions{Overrides: []AttributeMatchOverride{{Path: path.MatchRoot("settings").AtName("mode"), CaseInsensitiveStrings: new(true)}}}, actual: upperRules},
		{name: "override inside set elements ignores case on the field", opts: AttributeMatchOptions{Overrides: []AttributeMatchOverride{{Path: path.MatchRoot("rules").AtAnySetValue().AtName("mode"), CaseInsensitiveStrings: new(true)}}}, actual: upperRules, want: true},
		{name: "override on a parent applies to its nested fields", opts: AttributeMatchOptions{Overrides: []AttributeMatchOverride{{Path: path.MatchRoot("rules"), CaseInsensitiveStrings: new(true)}}}, actual: upperRules, want: true},
		{name: "nested override takes precedence over its parent override", opts: AttributeMatchOptions{Overrides: []AttributeMatchOverride{{Path: path.MatchRoot("rules"), CaseInsensitiveStrings: new(true)}, {Path: path.MatchRoot("rules").AtAnySetValue().AtName("mode"), CaseInsensitiveStrings: new(false)}}}, actual: upperRules},
		{name: "override inside list elements ignores case", opts: AttributeMatchOptions{Overrides: []AttributeMatchOverride{{Path: path.MatchRoot("items").AtAnyListIndex(), CaseInsensitiveStrings: new(true)}}}, actual: func(m *predicateTestModel) {
			m.Items = stringListValue(t, "X", "Y")
		}, want: true},
		{name: "override inside unordered list elements ignores case", opts: AttributeMatchOptions{UnorderedLists: true, Overrides: []AttributeMatchOverride{{Path: path.MatchRoot("items").AtAnyListIndex(), CaseInsensitiveStrings: new(true)}}}, actual: func(m *predicateTestModel) {
			m.Items = stringListValue(t, "Y", "X")
		}, want: true},
		{name: "override inside map values ignores case", opts: AttributeMatchOptions{Overrides: []AttributeMatchOverride{{Path: path.MatchRoot("labels").AtAnyMapKey(), CaseInsensitiveStrings: new(true)}}}, actual: func(m *predicateTestModel) {
			m.Labels = types.MapValueMust(types.StringType, map[string]attr.Value{"k": types.StringValue("V")})
		}, want: true},
		{name: "unordered lists on one attribute", opts: AttributeMatchOptions{Overrides: []AttributeMatchOverride{{Path: path.MatchRoot("items"), UnorderedLists: new(true)}}}, actual: func(m *predicateTestModel) {
			m.Items = stringListValue(t, "y", "x")
		}, want: true},
		{name: "ordered lists on one attribute override the global option", opts: AttributeMatchOptions{UnorderedLists: true, Overrides: []AttributeMatchOverride{{Path: path.MatchRoot("items"), UnorderedLists: new(false)}}}, actual: func(m *predicateTestModel) {
			m.Items = stringListValue(t, "y", "x")
		}},
		{name: "null expected ignored on one attribute", opts: AttributeMatchOptions{Overrides: []AttributeMatchOverride{{Path: path.MatchRoot("name"), IgnoreNullExpected: new(true)}}}, expected: func(m *predicateTestModel) {
			m.Name = types.StringNull()
		}, actual: func(m *predicateTestModel) { m.Name = types.StringValue("Graph default") }, want: true},
		{name: "null expected still compared outside the overridden attribute", opts: AttributeMatchOptions{Overrides: []AttributeMatchOverride{{Path: path.MatchRoot("name"), IgnoreNullExpected: new(true)}}}, expected: func(m *predicateTestModel) {
			m.Settings = settingsValue(t, types.StringNull(), types.BoolValue(true))
		}},
		{name: "null expected compared on one nested field despite the global option", opts: AttributeMatchOptions{IgnoreNullExpected: true, Overrides: []AttributeMatchOverride{{Path: path.MatchRoot("settings").AtName("mode"), IgnoreNullExpected: new(false)}}}, expected: func(m *predicateTestModel) {
			m.Settings = settingsValue(t, types.StringNull(), types.BoolValue(true))
		}},
		{name: "null matches empty on one attribute", opts: AttributeMatchOptions{Overrides: []AttributeMatchOverride{{Path: path.MatchRoot("items"), NullMatchesEmpty: new(true)}}}, expected: func(m *predicateTestModel) {
			m.Items = types.ListNull(types.StringType)
		}, actual: func(m *predicateTestModel) { m.Items = stringListValue(t) }, want: true},
		{name: "null does not match empty outside the overridden attribute", opts: AttributeMatchOptions{Overrides: []AttributeMatchOverride{{Path: path.MatchRoot("items"), NullMatchesEmpty: new(true)}}}, expected: func(m *predicateTestModel) {
			m.Tags = types.SetNull(types.StringType)
		}, actual: func(m *predicateTestModel) { m.Tags = stringSetValue(t) }},
		{name: "null matches empty on one nested field", opts: AttributeMatchOptions{Overrides: []AttributeMatchOverride{{Path: path.MatchRoot("rules").AtAnySetValue().AtName("mode"), NullMatchesEmpty: new(true)}}}, expected: func(m *predicateTestModel) {
			m.Rules = rulesValue(t, settingsValue(t, types.StringNull(), types.BoolValue(true)))
		}, actual: func(m *predicateTestModel) { m.Rules = rulesValue(t, settingsValue(t, types.StringValue(""), types.BoolValue(true))) }, want: true},
		{name: "required override fails an unknown expected attribute that reads back null", opts: AttributeMatchOptions{Overrides: []AttributeMatchOverride{{Path: path.MatchRoot("name"), RequireKnown: new(true)}}}, expected: func(m *predicateTestModel) {
			m.Name = types.StringUnknown()
		}, actual: func(m *predicateTestModel) { m.Name = types.StringNull() }},
		{name: "required override holds when an unknown expected attribute reads back populated", opts: AttributeMatchOptions{Overrides: []AttributeMatchOverride{{Path: path.MatchRoot("name"), RequireKnown: new(true)}}}, expected: func(m *predicateTestModel) {
			m.Name = types.StringUnknown()
		}, actual: func(m *predicateTestModel) { m.Name = types.StringValue("Graph default") }, want: true},
		{name: "required override fails a nested set field that reads back empty", opts: AttributeMatchOptions{NullMatchesEmpty: true, Overrides: []AttributeMatchOverride{{Path: path.MatchRoot("rules").AtAnySetValue().AtName("mode"), RequireKnown: new(true)}}}, expected: func(m *predicateTestModel) {
			m.Rules = rulesValue(t, settingsValue(t, types.StringNull(), types.BoolValue(true)))
		}, actual: func(m *predicateTestModel) { m.Rules = rulesValue(t, settingsValue(t, types.StringValue(""), types.BoolValue(true))) }},
		{name: "required override does not apply to nested values", opts: AttributeMatchOptions{Overrides: []AttributeMatchOverride{{Path: path.MatchRoot("settings"), RequireKnown: new(true)}}}, expected: func(m *predicateTestModel) {
			m.Settings = settingsValue(t, types.StringNull(), types.BoolValue(true))
		}, actual: func(m *predicateTestModel) { m.Settings = settingsValue(t, types.StringNull(), types.BoolValue(true)) }, want: true},
		{name: "required override false removes a globally required attribute", opts: AttributeMatchOptions{Paths: RootPaths("name"), RequireKnown: RootPaths("id", "name"), Overrides: []AttributeMatchOverride{{Path: path.MatchRoot("id"), RequireKnown: new(false)}}}, actual: func(m *predicateTestModel) {
			m.ID = types.StringNull()
		}, want: true},
		{name: "required override false leaves other required attributes", opts: AttributeMatchOptions{Paths: RootPaths("name"), RequireKnown: RootPaths("id", "count"), Overrides: []AttributeMatchOverride{{Path: path.MatchRoot("id"), RequireKnown: new(false)}}}, actual: func(m *predicateTestModel) {
			m.Count = types.Int64Null()
		}},
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

func TestAttributesMatchPredicateOverrideOnNestedPath(t *testing.T) {
	ctx := context.Background()
	expected := baseModel(t)
	actual := baseModel(t)
	actual.Settings = settingsValue(t, types.StringValue("BLOCK"), types.BoolValue(true))
	paths := []path.Path{path.Root("settings").AtName("mode")}

	pred := AttributesMatchPredicate(planFromModel(t, expected), AttributeMatchOptions{Paths: paths})
	assert.False(t, pred(ctx, stateFromModel(t, actual)))

	pred = AttributesMatchPredicate(planFromModel(t, expected), AttributeMatchOptions{Paths: paths, Overrides: []AttributeMatchOverride{{Path: path.MatchRoot("settings"), CaseInsensitiveStrings: new(true)}}})
	assert.True(t, pred(ctx, stateFromModel(t, actual)), "an override on a parent applies to a nested path in Paths")

	pred = AttributesMatchPredicate(planFromModel(t, expected), AttributeMatchOptions{Paths: paths, Overrides: []AttributeMatchOverride{{Path: path.MatchRoot("settings").AtName("mode"), CaseInsensitiveStrings: new(true)}}})
	assert.True(t, pred(ctx, stateFromModel(t, actual)), "an override on the nested path itself applies")
}

func TestAttributesMatchPredicateRejectValues(t *testing.T) {
	ctx := context.Background()
	ownerType := path.Root("name")
	rejectUnknown := []RejectedAttributeValues{RejectStrings(ownerType, "Unknown")}
	rejectSeveral := []RejectedAttributeValues{RejectStrings(ownerType, "Unknown", "Pending")}

	tests := []struct {
		name   string
		opts   AttributeMatchOptions
		actual func(*predicateTestModel)
		want   bool
	}{
		{name: "populated value that is not rejected holds", opts: AttributeMatchOptions{RejectValues: rejectUnknown}, want: true},
		{name: "single rejected value fails", opts: AttributeMatchOptions{RejectValues: rejectUnknown}, actual: func(m *predicateTestModel) { m.Name = types.StringValue("Unknown") }},
		{name: "first of several rejected values fails", opts: AttributeMatchOptions{RejectValues: rejectSeveral}, actual: func(m *predicateTestModel) { m.Name = types.StringValue("Unknown") }},
		{name: "second of several rejected values fails", opts: AttributeMatchOptions{RejectValues: rejectSeveral}, actual: func(m *predicateTestModel) { m.Name = types.StringValue("Pending") }},
		{name: "value not among several rejected values holds", opts: AttributeMatchOptions{RejectValues: rejectSeveral}, actual: func(m *predicateTestModel) { m.Name = types.StringValue("User") }, want: true},
		{name: "unknown read-back fails", opts: AttributeMatchOptions{RejectValues: rejectUnknown}, actual: func(m *predicateTestModel) { m.Name = types.StringUnknown() }},
		{name: "null read-back fails", opts: AttributeMatchOptions{RejectValues: rejectUnknown}, actual: func(m *predicateTestModel) { m.Name = types.StringNull() }},
		{name: "empty string read-back fails", opts: AttributeMatchOptions{RejectValues: rejectUnknown}, actual: func(m *predicateTestModel) { m.Name = types.StringValue("") }},
		{name: "rejected value compares case-sensitively by default", opts: AttributeMatchOptions{RejectValues: rejectUnknown}, actual: func(m *predicateTestModel) { m.Name = types.StringValue("unknown") }, want: true},
		{name: "rejected value follows the override for its path", opts: AttributeMatchOptions{RejectValues: rejectUnknown, Overrides: []AttributeMatchOverride{{Path: path.MatchRoot("name"), CaseInsensitiveStrings: new(true)}}}, actual: func(m *predicateTestModel) {
			m.Name = types.StringValue("unknown")
		}},
		{name: "rejected nested value fails", opts: AttributeMatchOptions{RejectValues: []RejectedAttributeValues{RejectStrings(path.Root("settings").AtName("mode"), "block")}}},
		{name: "rejected non-string value fails", opts: AttributeMatchOptions{RejectValues: []RejectedAttributeValues{{Path: path.Root("count"), Values: []attr.Value{types.Int64Value(2)}}}}},
		{name: "rejected value on a missing path fails", opts: AttributeMatchOptions{RejectValues: []RejectedAttributeValues{RejectStrings(path.Root("missing"), "Unknown")}}},
		{name: "rejected values combine with matched paths", opts: AttributeMatchOptions{Paths: RootPaths("count"), RejectValues: rejectUnknown}, actual: func(m *predicateTestModel) { m.Count = types.Int64Value(1) }},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := baseModel(t)
			if tt.actual != nil {
				tt.actual(&actual)
			}
			pred := AttributesMatchPredicate(planFromModel(t, baseModel(t)), tt.opts)
			assert.Equal(t, tt.want, pred(ctx, stateFromModel(t, actual)))
		})
	}
}

func TestRejectStrings(t *testing.T) {
	want := RejectedAttributeValues{Path: path.Root("owner_type"), Values: []attr.Value{types.StringValue("Unknown"), types.StringValue("")}}
	assert.Equal(t, want, RejectStrings(path.Root("owner_type"), "Unknown", ""))
	assert.Empty(t, RejectStrings(path.Root("owner_type")).Values)
}
