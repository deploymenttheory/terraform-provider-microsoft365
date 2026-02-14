package graphBetaUsersUser

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// UserListConfigModel represents the configuration for listing users
type UserListConfigModel struct {
	DisplayNameFilter       types.String `tfsdk:"display_name_filter"`
	UserPrincipalNameFilter types.String `tfsdk:"user_principal_name_filter"`
	AccountEnabledFilter    types.Bool   `tfsdk:"account_enabled_filter"`
	UserTypeFilter          types.String `tfsdk:"user_type_filter"`
	ODataFilter             types.String `tfsdk:"odata_filter"`
}
