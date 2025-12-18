package planmodifiers

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Int32Modifier defines the interface for int32 plan modifiers
type Int32Modifier interface {
	planmodifier.Int32
	Description(context.Context) string
	MarkdownDescription(context.Context) string
}

type int32Modifier struct {
	description         string
	markdownDescription string
}

func (m int32Modifier) Description(ctx context.Context) string {
	return m.description
}

func (m int32Modifier) MarkdownDescription(ctx context.Context) string {
	return m.markdownDescription
}

// UseStateForUnknown implementation
type useStateForUnknownInt32 struct {
	int32Modifier
}

func (m useStateForUnknownInt32) PlanModifyInt32(ctx context.Context, req planmodifier.Int32Request, resp *planmodifier.Int32Response) {
	if !req.PlanValue.IsUnknown() {
		return
	}

	if req.StateValue.IsNull() {
		return
	}

	resp.PlanValue = req.StateValue
}

// UseStateForUnknownInt32 returns a plan modifier that copies a known prior state int32
// value into a planned unknown value.
func UseStateForUnknownInt32() Int32Modifier {
	return useStateForUnknownInt32{
		int32Modifier: int32Modifier{
			description:         "Use state value if unknown",
			markdownDescription: "Use state value if unknown",
		},
	}
}

// Add new default value modifier
type defaultValueInt32 struct {
	int32Modifier
	defaultValue int32
}

func (m defaultValueInt32) PlanModifyInt32(ctx context.Context, req planmodifier.Int32Request, resp *planmodifier.Int32Response) {
	if !req.PlanValue.IsNull() {
		return
	}

	resp.PlanValue = types.Int32Value(m.defaultValue)
}

// Int32DefaultValue returns a plan modifier that sets a default value if the planned value is null
func Int32DefaultValue(defaultValue int32) Int32Modifier {
	return defaultValueInt32{
		int32Modifier: int32Modifier{
			description:         "Set default value if null",
			markdownDescription: "Set default value if null",
		},
		defaultValue: defaultValue,
	}
}

// RequiresOtherAttributeEnabledInt32 returns a plan modifier that ensures an Int32 attribute
// can only be used when another specified attribute is enabled (set to true).
func RequiresOtherAttributeEnabledInt32(dependencyPath path.Path) planmodifier.Int32 {
	return &requiresOtherAttributeEnabledInt32Modifier{
		dependencyPath: dependencyPath,
	}
}

type requiresOtherAttributeEnabledInt32Modifier struct {
	dependencyPath path.Path
}

func (m *requiresOtherAttributeEnabledInt32Modifier) Description(ctx context.Context) string {
	return fmt.Sprintf("Ensures this attribute is only used when %s is enabled", m.dependencyPath)
}

func (m *requiresOtherAttributeEnabledInt32Modifier) MarkdownDescription(ctx context.Context) string {
	return fmt.Sprintf("Ensures this attribute is only used when `%s` is enabled", m.dependencyPath)
}

func (m *requiresOtherAttributeEnabledInt32Modifier) PlanModifyInt32(ctx context.Context, req planmodifier.Int32Request, resp *planmodifier.Int32Response) {
	if req.PlanValue.IsNull() {
		return
	}

	var dependencyValue types.Bool
	diags := req.Plan.GetAttribute(ctx, m.dependencyPath, &dependencyValue)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !dependencyValue.IsNull() && !dependencyValue.IsUnknown() && !dependencyValue.ValueBool() {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid attribute usage",
			fmt.Sprintf("This attribute can only be used when %s is enabled (true)", m.dependencyPath),
		)
	}
}

// RequiresOtherAttributeValueInt32 returns a plan modifier that ensures an Int32 attribute
// can only be used when another specified attribute has a specific string value.
func RequiresOtherAttributeValueInt32(dependencyPath path.Path, requiredValue string) planmodifier.Int32 {
	return &requiresOtherAttributeValueInt32Modifier{
		dependencyPath: dependencyPath,
		requiredValue:  requiredValue,
	}
}

type requiresOtherAttributeValueInt32Modifier struct {
	dependencyPath path.Path
	requiredValue  string
}

func (m *requiresOtherAttributeValueInt32Modifier) Description(ctx context.Context) string {
	return fmt.Sprintf("Ensures this attribute is only used when %s is set to %q", m.dependencyPath, m.requiredValue)
}

func (m *requiresOtherAttributeValueInt32Modifier) MarkdownDescription(ctx context.Context) string {
	return fmt.Sprintf("Ensures this attribute is only used when `%s` is set to `%s`", m.dependencyPath, m.requiredValue)
}

func (m *requiresOtherAttributeValueInt32Modifier) PlanModifyInt32(ctx context.Context, req planmodifier.Int32Request, resp *planmodifier.Int32Response) {
	if req.PlanValue.IsNull() {
		return
	}

	var dependencyValue types.String
	diags := req.Plan.GetAttribute(ctx, m.dependencyPath, &dependencyValue)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !dependencyValue.IsNull() && !dependencyValue.IsUnknown() && dependencyValue.ValueString() != m.requiredValue {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid attribute usage",
			fmt.Sprintf("This attribute can only be used when %s is set to %q", m.dependencyPath, m.requiredValue),
		)
	}
}
