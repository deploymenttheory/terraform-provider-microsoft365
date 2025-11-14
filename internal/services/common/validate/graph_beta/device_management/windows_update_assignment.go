package sharedValidators

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// SingleIncludeExcludeAssignmentValidator ensures only one include and one exclude assignment block exist
type SingleIncludeExcludeAssignmentValidator struct{}

// Description of the validator
func (v SingleIncludeExcludeAssignmentValidator) Description(ctx context.Context) string {
	return "Ensure that only one assignment block with target 'include' and one with 'exclude' exist."
}

// MarkdownDescription of the validator
func (v SingleIncludeExcludeAssignmentValidator) MarkdownDescription(ctx context.Context) string {
	return "Ensure that only one assignment block with target `include` and one with `exclude` exist."
}

// ValidateList performs the actual validation
func (v SingleIncludeExcludeAssignmentValidator) ValidateList(ctx context.Context, req validator.ListRequest, resp *validator.ListResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	var assignments []types.Object
	diags := req.ConfigValue.ElementsAs(ctx, &assignments, false)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var includeCount, excludeCount int

	for _, assignment := range assignments {
		targetAttr := assignment.Attributes()["target"]

		if targetAttr.IsNull() || targetAttr.IsUnknown() {
			continue
		}

		targetValue := targetAttr.(types.String).ValueString()

		switch targetValue {
		case "include":
			includeCount++
		case "exclude":
			excludeCount++
		default:
			// ignore invalid targets (covered by separate validators)
		}
	}

	if includeCount > 1 {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Too Many 'include' Assignments",
			fmt.Sprintf("Only one assignment block with target 'include' is allowed, found %d.", includeCount),
		)
	}
	if excludeCount > 1 {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Too Many 'exclude' Assignments",
			fmt.Sprintf("Only one assignment block with target 'exclude' is allowed, found %d.", excludeCount),
		)
	}
}

// Factory function
func SingleIncludeExcludeAssignment() validator.List {
	return SingleIncludeExcludeAssignmentValidator{}
}
