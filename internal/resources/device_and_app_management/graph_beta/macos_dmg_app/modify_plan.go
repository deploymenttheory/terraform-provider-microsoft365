package graphBetaMacOSDmgApp

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// ModifyPlan handles plan modification/diff suppression for macOS DMG app resources.
func (r *MacOSDmgAppResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	// Currently, no specific plan modifications are needed for macOS DMG apps.
	// This method is implemented to satisfy the ResourceWithModifyPlan interface.
	// Future plan modifications can be added here as needed.
}
