// Package predicateequivalence compares a resource's former local consistency predicate with
// the shared crud.AttributesMatchPredicate expression that replaces it. It is a temporary
// migration harness for CPR-02 (#4949) and is removed once the migration is verified.
package predicateequivalence

import (
	"context"
	"fmt"
	"math/big"
	"os"
	"slices"
	"sort"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// Factory returns a fresh predicate for the expected values in expected.
type Factory func(expected tfsdk.State) func(context.Context, tfsdk.State) bool

// Options configures Compare.
type Options struct {
	// Reads is how many times each predicate is called on the same read-back state. Use
	// more than one for predicates that count consecutive matching reads.
	Reads int
	// Allow reports whether a difference for the named case is an intended one.
	Allow func(name string) bool
	// Base adjusts the generated base value, for attributes that need specific values.
	Base func(attrs map[string]tftypes.Value)
	// Known lists top-level attributes that Create and Update always set before the
	// read-back, so the expected value is never unknown or null.
	Known    []string
	MaxDepth int
}

type flags struct {
	computed, required bool
}

type variant struct {
	name string
	val  tftypes.Value
}

// Compare feeds the same expected/actual pairs to legacy and shared and fails t when the
// sequences of results differ for a case that opts.Allow does not accept. It returns the
// number of pairs compared and the names of the allowed differences.
func Compare(t *testing.T, s schema.Schema, legacy, shared Factory, opts Options) (int, []string) {
	t.Helper()
	ctx := context.Background()
	if opts.Reads == 0 {
		opts.Reads = 1
	}
	if opts.MaxDepth == 0 {
		opts.MaxDepth = 4
	}
	objType, ok := s.Type().TerraformType(ctx).(tftypes.Object)
	if !ok {
		t.Fatalf("schema type is not an object")
	}
	fl := map[string]flags{}
	walkAttributes("", s.Attributes, fl)
	walkBlocks("", s.Blocks, fl)

	base := fill(objType, "a", fl, "", true)
	if opts.Base != nil {
		attrs := copyAttributes(base)
		opts.Base(attrs)
		base = tftypes.NewValue(objType, attrs)
	}
	if os.Getenv("PE_DEBUG") != "" {
		t.Logf("base %s", base.String())
	}
	for _, name := range opts.Known {
		f := fl[name]
		f.computed = false
		fl[name] = f
	}
	plan := unknownComputed(objType, base, fl, "")
	sparse := nullOptional(objType, base, fl, nil)
	sparseExpected := nullOptional(objType, base, fl, opts.Known)

	type pair struct {
		name          string
		expected, got tftypes.Value
	}
	pairs := []pair{
		{"base/base", base, base},
		{"plan/base", plan, base},
		{"sparse/base", sparseExpected, base},
		{"base/sparse", base, sparse},
		{"sparse/sparse", sparseExpected, sparse},
		{"base/nullstate", base, tftypes.NewValue(objType, nil)},
	}
	for _, v := range variants(objType, base, fl, "", true, opts.MaxDepth) {
		pairs = append(pairs, pair{"base/" + v.name, base, v.val}, pair{"plan/" + v.name, plan, v.val})
	}
	for _, v := range variants(objType, base, fl, "", false, opts.MaxDepth) {
		pairs = append(pairs, pair{v.name + "/base", v.val, base})
	}
	for _, v := range variants(objType, plan, fl, "", false, opts.MaxDepth) {
		pairs = append(pairs, pair{"plan:" + v.name + "/base", v.val, base})
	}

	var allowed []string
	failures := 0
	for _, p := range pairs {
		if os.Getenv("PE_DEBUG") != "" {
			t.Logf("pair %s", p.name)
		}
		exp := tfsdk.State{Schema: s, Raw: p.expected}
		got := tfsdk.State{Schema: s, Raw: p.got}
		l := run(ctx, legacy, exp, got, opts.Reads)
		sh := run(ctx, shared, exp, got, opts.Reads)
		if fmt.Sprint(l) == fmt.Sprint(sh) {
			continue
		}
		msg := fmt.Sprintf("%s: legacy=%v shared=%v", p.name, l, sh)
		if opts.Allow != nil && opts.Allow(p.name) {
			allowed = append(allowed, msg)
			continue
		}
		failures++
		t.Errorf("difference %s", msg)
	}
	sort.Strings(allowed)
	for _, a := range allowed {
		t.Logf("allowed difference %s", a)
	}
	t.Logf("compared %d expected/actual pairs over %d attributes, %d unexpected differences, %d allowed differences", len(pairs), len(objType.AttributeTypes), failures, len(allowed))
	return len(pairs), allowed
}

func run(ctx context.Context, f Factory, exp, got tfsdk.State, reads int) (out []bool) {
	defer func() {
		if r := recover(); r != nil {
			out = append(out, false)
		}
	}()
	p := f(exp)
	for range reads {
		out = append(out, p(ctx, got))
	}
	return out
}

type attributeLike interface {
	IsComputed() bool
	IsRequired() bool
}

type nestedAttributeLike interface {
	GetNestedObject() schema.NestedAttributeObject
}

func walkAttributes(prefix string, attrs map[string]schema.Attribute, fl map[string]flags) {
	for name, a := range attrs {
		key := prefix + name
		fl[key] = flags{computed: a.IsComputed(), required: a.IsRequired()}
		if n, ok := a.(nestedAttributeLike); ok {
			walkAttributes(key+".", n.GetNestedObject().Attributes, fl)
		}
	}
}

func walkBlocks(prefix string, blocks map[string]schema.Block, fl map[string]flags) {
	for name, b := range blocks {
		key := prefix + name
		fl[key] = flags{}
		switch bb := b.(type) {
		case schema.SingleNestedBlock:
			walkAttributes(key+".", bb.Attributes, fl)
			walkBlocks(key+".", bb.Blocks, fl)
		case schema.ListNestedBlock:
			walkAttributes(key+".", bb.NestedObject.Attributes, fl)
			walkBlocks(key+".", bb.NestedObject.Blocks, fl)
		case schema.SetNestedBlock:
			walkAttributes(key+".", bb.NestedObject.Attributes, fl)
			walkBlocks(key+".", bb.NestedObject.Blocks, fl)
		}
	}
}

// fill returns a populated value of type typ. Strings are seed-name so that two fills with
// different seeds never collide.
func fill(typ tftypes.Type, seed string, fl map[string]flags, key string, top bool) tftypes.Value {
	switch {
	case typ.Is(tftypes.String):
		name := key
		if i := strings.LastIndex(name, "."); i >= 0 {
			name = name[i+1:]
		}
		return tftypes.NewValue(typ, seed+"-"+name)
	case typ.Is(tftypes.Bool):
		return tftypes.NewValue(typ, true)
	case typ.Is(tftypes.Number):
		return tftypes.NewValue(typ, big.NewFloat(float64(len(seed)+1)))
	case typ.Is(tftypes.DynamicPseudoType):
		return tftypes.NewValue(tftypes.String, seed)
	}
	switch tt := typ.(type) {
	case tftypes.Object:
		attrs := map[string]tftypes.Value{}
		for name, at := range tt.AttributeTypes {
			attrs[name] = fill(at, seed, fl, join(key, name), false)
		}
		return tftypes.NewValue(tt, attrs)
	case tftypes.List:
		return tftypes.NewValue(tt, []tftypes.Value{fill(tt.ElementType, seed, fl, key, false)})
	case tftypes.Set:
		return tftypes.NewValue(tt, []tftypes.Value{fill(tt.ElementType, seed, fl, key, false)})
	case tftypes.Map:
		return tftypes.NewValue(tt, map[string]tftypes.Value{"k": fill(tt.ElementType, seed, fl, key, false)})
	case tftypes.Tuple:
		elems := make([]tftypes.Value, len(tt.ElementTypes))
		for i, et := range tt.ElementTypes {
			elems[i] = fill(et, seed, fl, key, false)
		}
		return tftypes.NewValue(tt, elems)
	}
	return tftypes.NewValue(typ, nil)
}

func join(prefix, name string) string {
	if prefix == "" {
		return name
	}
	return prefix + "." + name
}

// unknownComputed returns base with each top-level Computed attribute that is not Required set to unknown,
// as in a Create plan.
func unknownComputed(objType tftypes.Object, base tftypes.Value, fl map[string]flags, _ string) tftypes.Value {
	attrs := copyAttributes(base)
	for name, at := range objType.AttributeTypes {
		if f := fl[name]; f.computed && !f.required {
			attrs[name] = tftypes.NewValue(at, tftypes.UnknownValue)
		}
	}
	return tftypes.NewValue(objType, attrs)
}

// nullOptional returns base with each top-level attribute that is not Required or kept set to null.
func nullOptional(objType tftypes.Object, base tftypes.Value, fl map[string]flags, keep []string) tftypes.Value {
	attrs := copyAttributes(base)
	for name, at := range objType.AttributeTypes {
		if !fl[name].required && !slices.Contains(keep, name) {
			attrs[name] = tftypes.NewValue(at, nil)
		}
	}
	return tftypes.NewValue(objType, attrs)
}

// variants returns single mutations of v. actualSide allows mutations Graph could return
// (any null); the expected side only nulls non-Required attributes and only makes Computed
// attributes unknown, as a plan can.
func variants(typ tftypes.Type, v tftypes.Value, fl map[string]flags, key string, actualSide bool, depth int) []variant {
	var out []variant
	if !v.IsKnown() || v.IsNull() || depth < 0 {
		return out
	}
	f := fl[key]
	if key != "" {
		if actualSide || !f.required {
			out = append(out, variant{key + "=null", tftypes.NewValue(typ, nil)})
		}
		if !actualSide && f.computed {
			out = append(out, variant{key + "=unknown", tftypes.NewValue(typ, tftypes.UnknownValue)})
		}
	}
	switch {
	case typ.Is(tftypes.String):
		var s string
		_ = v.As(&s)
		out = append(out,
			variant{key + "=changed", tftypes.NewValue(typ, s+"-x")},
			variant{key + "=upper", tftypes.NewValue(typ, strings.ToUpper(s))},
			variant{key + "=empty", tftypes.NewValue(typ, "")},
		)
		return out
	case typ.Is(tftypes.Bool):
		var b bool
		_ = v.As(&b)
		return append(out, variant{key + "=flipped", tftypes.NewValue(typ, !b)})
	case typ.Is(tftypes.Number):
		var n big.Float
		_ = v.As(&n)
		return append(out, variant{key + "=changed", tftypes.NewValue(typ, new(big.Float).Add(&n, big.NewFloat(1)))})
	}
	switch tt := typ.(type) {
	case tftypes.Object:
		attrs := map[string]tftypes.Value{}
		_ = v.As(&attrs)
		names := make([]string, 0, len(attrs))
		for name := range attrs {
			names = append(names, name)
		}
		sort.Strings(names)
		for _, name := range names {
			for _, cv := range variants(tt.AttributeTypes[name], attrs[name], fl, join(key, name), actualSide, depth-1) {
				copyAttrs := map[string]tftypes.Value{}
				for k, val := range attrs {
					copyAttrs[k] = val
				}
				copyAttrs[name] = cv.val
				out = append(out, variant{cv.name, tftypes.NewValue(tt, copyAttrs)})
			}
		}
	case tftypes.List, tftypes.Set:
		var elems []tftypes.Value
		_ = v.As(&elems)
		var et tftypes.Type
		if l, ok := tt.(tftypes.List); ok {
			et = l.ElementType
		} else {
			et = tt.(tftypes.Set).ElementType
		}
		out = append(out,
			variant{key + "=emptycollection", tftypes.NewValue(typ, []tftypes.Value{})},
			variant{key + "=extraelement", tftypes.NewValue(typ, append(append([]tftypes.Value{}, elems...), fill(et, "b", fl, key, false)))},
		)
		if len(elems) > 0 {
			for _, cv := range variants(et, elems[0], fl, key, actualSide, depth-1) {
				// A collection element is never null in a plan or a read, and a set element is never unknown.
				if cv.name == key+"=null" {
					continue
				}
				if _, isSet := tt.(tftypes.Set); isSet && !cv.val.IsKnown() {
					continue
				}
				replaced := append([]tftypes.Value{cv.val}, elems[1:]...)
				out = append(out, variant{cv.name + "[elem]", tftypes.NewValue(typ, replaced)})
			}
		}
	case tftypes.Map:
		elems := map[string]tftypes.Value{}
		_ = v.As(&elems)
		out = append(out, variant{key + "=emptycollection", tftypes.NewValue(typ, map[string]tftypes.Value{})})
		if e, ok := elems["k"]; ok {
			for _, cv := range variants(tt.ElementType, e, fl, key, actualSide, depth-1) {
				if cv.name == key+"=null" {
					continue
				}
				out = append(out, variant{cv.name + "[k]", tftypes.NewValue(typ, map[string]tftypes.Value{"k": cv.val})})
			}
		}
	}
	return out
}

// copyAttributes returns a copy of an object value's attributes. Value.As shares the
// value's own map, so it must not be modified.
func copyAttributes(v tftypes.Value) map[string]tftypes.Value {
	shared := map[string]tftypes.Value{}
	_ = v.As(&shared)
	out := make(map[string]tftypes.Value, len(shared))
	for k, val := range shared {
		out[k] = val
	}
	return out
}
