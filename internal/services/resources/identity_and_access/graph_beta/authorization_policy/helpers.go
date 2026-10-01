package graphBetaAuthorizationPolicy

import (
	"context"
	"strings"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// mapPermissionGrantPolicyIDs preserves configured prefix casing when Graph returns
// ManagePermissionGrants instead of managePermissionGrants. Policy IDs remain case-sensitive,
// and remote additions or removals are retained so Terraform can detect drift.
func mapPermissionGrantPolicyIDs(ctx context.Context, remote []string, configured types.Set) types.Set {
	configuredValues := make(map[string]string, len(configured.Elements()))
	for _, element := range configured.Elements() {
		if value, ok := element.(types.String); ok && !value.IsNull() && !value.IsUnknown() {
			configuredValues[permissionGrantPolicyKey(value.ValueString())] = value.ValueString()
		}
	}

	values := make([]string, 0, len(remote))
	for _, value := range remote {
		if configuredValue, ok := configuredValues[permissionGrantPolicyKey(value)]; ok {
			value = configuredValue
		}
		values = append(values, value)
	}
	return convert.GraphToFrameworkStringSetPreserveEmpty(ctx, values)
}

func permissionGrantPolicyKey(value string) string {
	prefix, id, found := strings.Cut(value, ".")
	if !found {
		return value
	}
	return strings.ToLower(prefix) + "." + id
}
