package graphBetaIdentityAndAccessRoleDefinition

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// mapOptionalString maps an optional Graph string to state.
// Microsoft Entra returns an empty string for some unset optional properties (for example description
// on a role created without one), so empty is treated as null to match an omitted configuration value.
func mapOptionalString(value *string) types.String {
	if value == nil || *value == "" {
		return types.StringNull()
	}
	return types.StringValue(*value)
}
