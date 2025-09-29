package graphBetaWindowsAutopilotDevicePreparationPolicy

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

const (
	// The App ID for "Intune Provisioning Client" / "Intune Autopilot ConfidentialClient" service principal
	intuneProvisioningClientAppID = "f1346770-5b25-470b-88bd-d5744ab7952c"
)

// validateSecurityGroupOwnership validates that the specified security group has the Intune Provisioning Client as an owner
func validateSecurityGroupOwnership(ctx context.Context, client *msgraphbetasdk.GraphServiceClient, groupID string) diag.Diagnostics {
	var diags diag.Diagnostics

	tflog.Info(ctx, fmt.Sprintf("Validating security group %s has Intune Provisioning Client as owner", groupID))

	owners, err := client.Groups().
		ByGroupId(groupID).
		Owners().
		Get(ctx, nil)

	if err != nil {
		tflog.Error(ctx, "Failed to get security group owners", map[string]any{
			"group_id": groupID,
			"error":    err.Error(),
		})
		diags.AddError(
			"Failed to validate security group ownership",
			fmt.Sprintf("Could not retrieve owners for security group %s: %s", groupID, err.Error()),
		)
		return diags
	}

	tflog.Debug(ctx, fmt.Sprintf("Retrieved %d owners for security group %s", len(owners.GetValue()), groupID))

	// Check if the Intune Provisioning Client is an owner
	hasIntuneProvisioningClient := false
	for _, owner := range owners.GetValue() {
		servicePrincipal, ok := owner.(models.ServicePrincipalable)
		if ok {
			appID := servicePrincipal.GetAppId()
			if appID != nil && *appID == intuneProvisioningClientAppID {
				tflog.Info(ctx, "Found Intune Provisioning Client as owner of security group", map[string]any{
					"group_id": groupID,
					"app_id":   *appID,
				})
				hasIntuneProvisioningClient = true
				break
			}
		}
	}

	if !hasIntuneProvisioningClient {
		tflog.Error(ctx, "Security group does not have Intune Provisioning Client as owner", map[string]any{
			"group_id":                   groupID,
			"required_service_principal": intuneProvisioningClientAppID,
		})
		diags.AddError(
			"Invalid security group ownership",
			fmt.Sprintf("Security group %s must have the Intune Provisioning Client (AppID: %s) set as its owner. In some tenants, this service principal may appear as 'Intune Autopilot ConfidentialClient'.",
				groupID, intuneProvisioningClientAppID),
		)
	}

	return diags
}
