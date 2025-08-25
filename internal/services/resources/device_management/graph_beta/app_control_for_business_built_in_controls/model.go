package graphBetaAppControlForBusinessBuiltInControls

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// AppControlForBusinessResourceBuiltInControlsModel represents the Terraform resource schema for App Control for Business configuration policies
type AppControlForBusinessResourceBuiltInControlsModel struct {
	ID                             types.String   `tfsdk:"id"`
	Name                           types.String   `tfsdk:"name"`
	Description                    types.String   `tfsdk:"description"`
	RoleScopeTagIds                types.Set      `tfsdk:"role_scope_tag_ids"`
	EnableAppControl               types.String   `tfsdk:"enable_app_control"`
	AdditionalRulesForTrustingApps types.Set      `tfsdk:"additional_rules_for_trusting_apps"`
	Assignments                    types.Set      `tfsdk:"assignments"`
	Timeouts                       timeouts.Value `tfsdk:"timeouts"`
}
