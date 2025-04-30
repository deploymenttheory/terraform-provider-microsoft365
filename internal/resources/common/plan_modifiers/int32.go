package planmodifiers

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

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
