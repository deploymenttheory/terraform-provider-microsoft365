// REF: https://www.ctrlshiftenter.cloud/2025/03/17/mastering-app-control-for-business-part-2-policy-templates-rule-options/
// REF: https://github.com/HotCakeX/Harden-Windows-Security/wiki/AppControl-Manager

package graphBetaAppControlForBusinessPolicy

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// AppControlForBusinessPolicyResourceModel represents the Terraform resource schema for App Control for Business configuration policies with XML content
type AppControlForBusinessPolicyResourceModel struct {
	ID              types.String   `tfsdk:"id"`
	Name            types.String   `tfsdk:"name"`
	Description     types.String   `tfsdk:"description"`
	RoleScopeTagIds types.Set      `tfsdk:"role_scope_tag_ids"`
	PolicyXML       types.String   `tfsdk:"policy_xml"`
	Assignments     types.Set      `tfsdk:"assignments"`
	Timeouts        timeouts.Value `tfsdk:"timeouts"`
}
