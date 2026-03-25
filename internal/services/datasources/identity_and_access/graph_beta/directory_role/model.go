// REF: https://learn.microsoft.com/en-us/graph/api/resources/directoryrole?view=graph-rest-beta

package graphBetaDirectoryRole

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type DirectoryRoleDataSourceModel struct {
	RoleID      types.String          `tfsdk:"role_id"`
	DisplayName types.String          `tfsdk:"display_name"`
	ListAll     types.Bool            `tfsdk:"list_all"`
	Items       []DirectoryRoleModel  `tfsdk:"items"`
	Timeouts    timeouts.Value        `tfsdk:"timeouts"`
}

type DirectoryRoleModel struct {
	ID             types.String `tfsdk:"id"`
	DisplayName    types.String `tfsdk:"display_name"`
	Description    types.String `tfsdk:"description"`
	RoleTemplateID types.String `tfsdk:"role_template_id"`
}
