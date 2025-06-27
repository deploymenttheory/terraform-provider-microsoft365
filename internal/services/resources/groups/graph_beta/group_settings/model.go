// REF: https://learn.microsoft.com/en-us/graph/api/resources/directorysetting?view=graph-rest-beta
package graphBetaGroupSettings

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type GroupSettingsResourceModel struct {
	ID          types.String   `tfsdk:"id"`
	GroupID     types.String   `tfsdk:"group_id"`
	DisplayName types.String   `tfsdk:"display_name"`
	TemplateID  types.String   `tfsdk:"template_id"`
	Values      types.Set      `tfsdk:"values"`
	Timeouts    timeouts.Value `tfsdk:"timeouts"`
}

type SettingValueModel struct {
	Name  types.String `tfsdk:"name"`
	Value types.String `tfsdk:"value"`
}
