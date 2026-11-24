package crud

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// StableReads is how many matching reads in a row confirm a write. A stale read
// can follow a fresh one.
const StableReads = 3

// RequireStableReads holds only after StableReads matching reads in a row.
func RequireStableReads(
	predicate func(context.Context, tfsdk.State) bool,
) func(context.Context, tfsdk.State) bool {
	return RequireConsecutiveReads(StableReads, predicate)
}

// RequireConsecutiveReads holds only after count matching reads in a row. A read that
// does not match resets the count. Use it instead of RequireStableReads when a resource
// has shown stale reads after more than StableReads fresh ones.
func RequireConsecutiveReads(count int, predicate func(context.Context, tfsdk.State) bool) func(context.Context, tfsdk.State) bool {
	matched := 0
	return func(ctx context.Context, s tfsdk.State) bool {
		if !predicate(ctx, s) {
			matched = 0
			return false
		}
		matched++
		return matched >= count
	}
}

// AttributeGetter reads one attribute by path. tfsdk.Plan, tfsdk.State and tfsdk.Config
// all satisfy it, so the expected values can come from the request plan, or from a
// response state that Create or Update set from the planned model before reading back.
type AttributeGetter interface {
	GetAttribute(ctx context.Context, p path.Path, target any) diag.Diagnostics
}

// AttributeMatchOptions configures AttributesMatchPredicate.
type AttributeMatchOptions struct {
	// Paths lists the attributes whose expected value must be reflected by the
	// read-back state. An expected value that is unknown, at any depth, is skipped
	// because Graph computes it. Sets, and lists when UnorderedLists is true, match
	// regardless of element order.
	Paths []path.Path
	// RequireKnown lists attributes that must be known and not null in the read-back
	// state, and not empty when they are strings, such as a Computed id that only a
	// real API read populates.
	RequireKnown []path.Path
	// IgnoreNullExpected skips an expected value that is null, at any depth. Use it
	// when Graph fills unconfigured Optional attributes with its own default. By
	// default a null expected value must read back as null.
	IgnoreNullExpected bool
	// NullMatchesEmpty treats null and an empty string, list, set or map as equal.
	NullMatchesEmpty bool
	// CaseInsensitiveStrings compares strings with strings.EqualFold, for values such
	// as GUIDs that Graph returns in a different case from the configuration.
	CaseInsensitiveStrings bool
	// UnorderedLists matches list elements regardless of order, for lists that Graph
	// returns in its own order.
	UnorderedLists bool
}

// RootPaths returns a path.Root for each top-level attribute name.
func RootPaths(names ...string) []path.Path {
	paths := make([]path.Path, 0, len(names))
	for _, name := range names {
		paths = append(paths, path.Root(name))
	}
	return paths
}

// AttributesMatchPredicate returns a ConsistencyPredicate for ReadWithRetry that holds
// when the read-back state reflects the expected value of every path in opts.Paths and
// every path in opts.RequireKnown is populated. A read that removed the resource from
// state never holds. Wrap the result in RequireStableReads when Graph can serve a stale
// read after a fresh one. See ADR-002.
func AttributesMatchPredicate(expected AttributeGetter, opts AttributeMatchOptions) func(ctx context.Context, state tfsdk.State) bool {
	return func(ctx context.Context, state tfsdk.State) bool {
		if state.Raw.IsNull() {
			tflog.Debug(ctx, "Consistency predicate: read-back state is null")
			return false
		}
		for _, p := range opts.RequireKnown {
			var actual attr.Value
			if diags := state.GetAttribute(ctx, p, &actual); diags.HasError() || !isPopulated(ctx, actual) {
				tflog.Debug(ctx, "Consistency predicate: required attribute is not populated", map[string]any{"path": p.String()})
				return false
			}
		}
		for _, p := range opts.Paths {
			var want, got attr.Value
			if diags := expected.GetAttribute(ctx, p, &want); diags.HasError() {
				tflog.Debug(ctx, "Consistency predicate: expected attribute cannot be read", map[string]any{"path": p.String()})
				return false
			}
			if diags := state.GetAttribute(ctx, p, &got); diags.HasError() {
				tflog.Debug(ctx, "Consistency predicate: read-back attribute cannot be read", map[string]any{"path": p.String()})
				return false
			}
			if !opts.valueMatches(ctx, want, got) {
				tflog.Debug(ctx, "Consistency predicate: read-back attribute does not match the expected value", map[string]any{"path": p.String()})
				return false
			}
		}
		return true
	}
}

// isPopulated reports whether v is known, not null and, for a string, not empty.
func isPopulated(ctx context.Context, v attr.Value) bool {
	if v == nil || v.IsNull() || v.IsUnknown() {
		return false
	}
	if s, ok := v.(basetypes.StringValuable); ok {
		str, diags := s.ToStringValue(ctx)
		return !diags.HasError() && str.ValueString() != ""
	}
	return true
}

// valueMatches reports whether the read-back value got reflects the expected value want.
func (opts AttributeMatchOptions) valueMatches(ctx context.Context, want, got attr.Value) bool {
	if want == nil || want.IsUnknown() {
		return true
	}
	if got == nil || got.IsUnknown() {
		return false
	}
	if want.IsNull() || got.IsNull() {
		if want.IsNull() && opts.IgnoreNullExpected {
			return true
		}
		if want.IsNull() && got.IsNull() {
			return true
		}
		return opts.NullMatchesEmpty && isEmpty(ctx, want) && isEmpty(ctx, got)
	}

	switch w := want.(type) {
	case basetypes.ObjectValuable:
		return opts.objectMatches(ctx, w, got)
	case basetypes.SetValuable:
		return opts.setMatches(ctx, w, got)
	case basetypes.ListValuable:
		return opts.listMatches(ctx, w, got)
	case basetypes.MapValuable:
		return opts.mapMatches(ctx, w, got)
	case basetypes.StringValuable:
		return opts.stringMatches(ctx, w, got)
	default:
		return want.Equal(got)
	}
}

// objectMatches compares each expected object attribute with the read-back attribute of the same name.
func (opts AttributeMatchOptions) objectMatches(ctx context.Context, want basetypes.ObjectValuable, got attr.Value) bool {
	g, ok := got.(basetypes.ObjectValuable)
	if !ok {
		return false
	}
	wantObject, wd := want.ToObjectValue(ctx)
	gotObject, gd := g.ToObjectValue(ctx)
	if wd.HasError() || gd.HasError() {
		return false
	}
	gotAttrs := gotObject.Attributes()
	for name, value := range wantObject.Attributes() {
		if !opts.valueMatches(ctx, value, gotAttrs[name]) {
			return false
		}
	}
	return true
}

// setMatches compares set elements regardless of order.
func (opts AttributeMatchOptions) setMatches(ctx context.Context, want basetypes.SetValuable, got attr.Value) bool {
	g, ok := got.(basetypes.SetValuable)
	if !ok {
		return false
	}
	wantSet, wd := want.ToSetValue(ctx)
	gotSet, gd := g.ToSetValue(ctx)
	if wd.HasError() || gd.HasError() {
		return false
	}
	return opts.elementsMatchUnordered(ctx, wantSet.Elements(), gotSet.Elements())
}

// listMatches compares list elements in order, or regardless of order when UnorderedLists is set.
func (opts AttributeMatchOptions) listMatches(ctx context.Context, want basetypes.ListValuable, got attr.Value) bool {
	g, ok := got.(basetypes.ListValuable)
	if !ok {
		return false
	}
	wantList, wd := want.ToListValue(ctx)
	gotList, gd := g.ToListValue(ctx)
	if wd.HasError() || gd.HasError() {
		return false
	}
	if opts.UnorderedLists {
		return opts.elementsMatchUnordered(ctx, wantList.Elements(), gotList.Elements())
	}
	return opts.elementsMatchOrdered(ctx, wantList.Elements(), gotList.Elements())
}

// mapMatches requires the same keys and compares the value under each key.
func (opts AttributeMatchOptions) mapMatches(ctx context.Context, want basetypes.MapValuable, got attr.Value) bool {
	g, ok := got.(basetypes.MapValuable)
	if !ok {
		return false
	}
	wantMap, wd := want.ToMapValue(ctx)
	gotMap, gd := g.ToMapValue(ctx)
	if wd.HasError() || gd.HasError() {
		return false
	}
	gotElems := gotMap.Elements()
	if len(wantMap.Elements()) != len(gotElems) {
		return false
	}
	for key, value := range wantMap.Elements() {
		actual, ok := gotElems[key]
		if !ok || !opts.valueMatches(ctx, value, actual) {
			return false
		}
	}
	return true
}

// stringMatches compares strings exactly, or with strings.EqualFold when CaseInsensitiveStrings is set.
func (opts AttributeMatchOptions) stringMatches(ctx context.Context, want basetypes.StringValuable, got attr.Value) bool {
	if !opts.CaseInsensitiveStrings {
		return want.Equal(got)
	}
	g, ok := got.(basetypes.StringValuable)
	if !ok {
		return false
	}
	wantString, wd := want.ToStringValue(ctx)
	gotString, gd := g.ToStringValue(ctx)
	return !wd.HasError() && !gd.HasError() && strings.EqualFold(wantString.ValueString(), gotString.ValueString())
}

// elementsMatchOrdered matches each expected element with the read-back element at the same index.
func (opts AttributeMatchOptions) elementsMatchOrdered(ctx context.Context, want, got []attr.Value) bool {
	if len(want) != len(got) {
		return false
	}
	for i := range want {
		if !opts.valueMatches(ctx, want[i], got[i]) {
			return false
		}
	}
	return true
}

// elementsMatchUnordered pairs each expected element with a distinct read-back element in
// any order. It uses augmenting paths, so an expected element with unknown parts cannot
// claim a read-back element that another expected element needs.
func (opts AttributeMatchOptions) elementsMatchUnordered(ctx context.Context, want, got []attr.Value) bool {
	if len(want) != len(got) {
		return false
	}
	candidates := make([][]int, len(want))
	for i, w := range want {
		for j, g := range got {
			if opts.valueMatches(ctx, w, g) {
				candidates[i] = append(candidates[i], j)
			}
		}
	}
	owner := make([]int, len(got))
	for j := range owner {
		owner[j] = -1
	}
	var assign func(i int, seen []bool) bool
	assign = func(i int, seen []bool) bool {
		for _, j := range candidates[i] {
			if seen[j] {
				continue
			}
			seen[j] = true
			if owner[j] < 0 || assign(owner[j], seen) {
				owner[j] = i
				return true
			}
		}
		return false
	}
	for i := range want {
		if !assign(i, make([]bool, len(got))) {
			return false
		}
	}
	return true
}

// isEmpty reports whether v is null, or an empty string, list, set or map.
func isEmpty(ctx context.Context, v attr.Value) bool {
	if v.IsNull() {
		return true
	}
	switch t := v.(type) {
	case basetypes.StringValuable:
		s, d := t.ToStringValue(ctx)
		return !d.HasError() && s.ValueString() == ""
	case basetypes.ListValuable:
		l, d := t.ToListValue(ctx)
		return !d.HasError() && len(l.Elements()) == 0
	case basetypes.SetValuable:
		s, d := t.ToSetValue(ctx)
		return !d.HasError() && len(s.Elements()) == 0
	case basetypes.MapValuable:
		m, d := t.ToMapValue(ctx)
		return !d.HasError() && len(m.Elements()) == 0
	default:
		return false
	}
}
