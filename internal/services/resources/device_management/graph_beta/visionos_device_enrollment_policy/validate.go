package graphBetaVisionOSDeviceEnrollmentPolicy

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// intuneProvisioningClientAppID is the App ID for the "Intune Provisioning Client" service
// principal (shown as "Intune Autopilot ConfidentialClient" in some tenants), which must own the
// device_security_group for enrollment time grouping to work.
const intuneProvisioningClientAppID = "f1346770-5b25-470b-88bd-d5744ab7952c"

// validateSecurityGroupOwnership validates that the specified security group has the Intune
// Provisioning Client as an owner, as required for enrollment time grouping.
func validateSecurityGroupOwnership(ctx context.Context, client *msgraphbetasdk.GraphServiceClient, groupID string) diag.Diagnostics {
	var diags diag.Diagnostics

	tflog.Info(ctx, fmt.Sprintf("Validating security group %s has Intune Provisioning Client as owner", groupID))

	owners, err := client.
		Groups().
		ByGroupId(groupID).
		Owners().
		Get(ctx, nil)
	if err != nil {
		diags.AddError(
			"Failed to validate security group ownership",
			fmt.Sprintf("Could not retrieve owners for security group %s: %s", groupID, err.Error()),
		)
		return diags
	}

	for _, owner := range owners.GetValue() {
		if servicePrincipal, ok := owner.(models.ServicePrincipalable); ok {
			if appID := servicePrincipal.GetAppId(); appID != nil && *appID == intuneProvisioningClientAppID {
				return diags
			}
		}
	}

	diags.AddError(
		"Invalid security group ownership",
		fmt.Sprintf(
			"Security group %s must have the Intune Provisioning Client (AppID: %s) set as its owner. In some tenants, this service principal may appear as 'Intune Autopilot ConfidentialClient'.",
			groupID,
			intuneProvisioningClientAppID,
		),
	)
	return diags
}
