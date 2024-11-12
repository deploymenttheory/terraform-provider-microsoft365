// planmodifiers/int64.go

package planmodifiers

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Int64Modifier defines the interface for int64 plan modifiers
type Int64Modifier interface {
	planmodifier.Int64
	Description(context.Context) string
	MarkdownDescription(context.Context) string
}

type int64Modifier struct {
	description         string
	markdownDescription string
}

func (m int64Modifier) Description(ctx context.Context) string {
	return m.description
}

func (m int64Modifier) MarkdownDescription(ctx context.Context) string {
	return m.markdownDescription
}

// UseStateForUnknown implementation
type useStateForUnknownInt64 struct {
	int64Modifier
}

func (m useStateForUnknownInt64) PlanModifyInt64(ctx context.Context, req planmodifier.Int64Request, resp *planmodifier.Int64Response) {
	if !req.PlanValue.IsUnknown() {
		return
	}

	if req.StateValue.IsNull() {
		return
	}

	resp.PlanValue = req.StateValue
}

// UseStateForUnknownInt64 returns a plan modifier that copies a known prior state int64
// value into a planned unknown value.
func UseStateForUnknownInt64() Int64Modifier {
	return useStateForUnknownInt64{
		int64Modifier: int64Modifier{
			description:         "Use state value if unknown",
			markdownDescription: "Use state value if unknown",
		},
	}
}

// RequiresReplace implementation
type requiresReplaceInt64 struct {
	int64Modifier
}

func (m requiresReplaceInt64) PlanModifyInt64(ctx context.Context, req planmodifier.Int64Request, resp *planmodifier.Int64Response) {
	if req.PlanValue.IsUnknown() {
		return
	}

	if req.StateValue.IsNull() {
		return
	}

	if req.PlanValue.Equal(req.StateValue) {
		return
	}

	resp.RequiresReplace = true
}

// RequiresReplaceInt64 returns a plan modifier that requires resource replacement
// if the value changes.
func RequiresReplaceInt64() Int64Modifier {
	return requiresReplaceInt64{
		int64Modifier: int64Modifier{
			description:         "Requires resource replacement if value changes",
			markdownDescription: "Requires resource replacement if value changes",
		},
	}
}

// DefaultValue implementation
type defaultValueInt64 struct {
	int64Modifier
	defaultValue types.Int64
}

func (m defaultValueInt64) PlanModifyInt64(ctx context.Context, req planmodifier.Int64Request, resp *planmodifier.Int64Response) {
	if !req.PlanValue.IsNull() {
		return
	}

	resp.PlanValue = m.defaultValue
}

// DefaultValueInt64 returns a plan modifier that sets a default value if the planned
// value is null.
func DefaultValueInt64(defaultValue int64) Int64Modifier {
	return defaultValueInt64{
		int64Modifier: int64Modifier{
			description:         fmt.Sprintf("Default value set to %d", defaultValue),
			markdownDescription: fmt.Sprintf("Default value set to `%d`", defaultValue),
		},
		defaultValue: types.Int64Value(defaultValue),
	}
}

// RequiresReplaceIf implementation
type requiresReplaceIfInt64 struct {
	int64Modifier
	predicate func(context.Context, planmodifier.Int64Request) bool
}

func (m requiresReplaceIfInt64) PlanModifyInt64(ctx context.Context, req planmodifier.Int64Request, resp *planmodifier.Int64Response) {
	if req.PlanValue.IsUnknown() {
		return
	}

	if req.StateValue.IsNull() {
		return
	}

	if req.PlanValue.Equal(req.StateValue) {
		return
	}

	if !m.predicate(ctx, req) {
		return
	}

	resp.RequiresReplace = true
}

// RequiresReplaceIfInt64 returns a plan modifier that requires resource replacement
// if the value changes and the given predicate returns true.
func RequiresReplaceIfInt64(predicate func(context.Context, planmodifier.Int64Request) bool) Int64Modifier {
	return requiresReplaceIfInt64{
		int64Modifier: int64Modifier{
			description:         "Requires resource replacement if value changes and condition is met",
			markdownDescription: "Requires resource replacement if value changes and condition is met",
		},
		predicate: predicate,
	}
}
