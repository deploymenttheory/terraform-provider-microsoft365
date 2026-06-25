package graphBetaApplicationsServicePrincipalTokenLifetimePolicyAssignment

import (
	"context"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// MapRemoteStateToTerraform maps the assignment state back to Terraform state
func MapRemoteStateToTerraform(ctx context.Context, data *ServicePrincipalTokenLifetimePolicyAssignmentResourceModel, found bool) {
	if !found {
		tflog.Debug(ctx, "Token lifetime policy assignment not found in remote state")
		return
	}

	data.ID.String()
	tflog.Debug(ctx, "Token lifetime policy assignment found in remote state", map[string]any{
		"service_principal_id":    data.ServicePrincipalID.ValueString(),
		"token_lifetime_policy_id": data.TokenLifetimePolicyID.ValueString(),
	})
}
