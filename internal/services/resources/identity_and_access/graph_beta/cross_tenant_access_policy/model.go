// REF: https://learn.microsoft.com/en-us/graph/api/resources/crosstenantaccesspolicy?view=graph-rest-beta
package graphBetaCrossTenantAccessPolicy

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// CrossTenantAccessPolicyResourceModel represents the schema for the Cross Tenant Access Policy resource.
// This is a singleton resource — one policy exists per tenant and cannot be created or deleted via the API.
type CrossTenantAccessPolicyResourceModel struct {
	ID                       types.String   `tfsdk:"id"`
	DisplayName              types.String   `tfsdk:"display_name"`
	AllowedCloudEndpoints    types.Set      `tfsdk:"allowed_cloud_endpoints"`
	RestoreDefaultsOnDestroy types.Bool     `tfsdk:"restore_defaults_on_destroy"`
	Timeouts                 timeouts.Value `tfsdk:"timeouts"`
}
