package graphBetaAppControlForBusinessManagedInstaller

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// constructResource maps the Terraform schema to the API request
// For this resource, we don't construct a request body since we use POST to setAsManagedInstaller endpoint
func constructResource(ctx context.Context, data *AppControlForBusinessManagedInstallerResourceModel) error {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s resource from model", ResourceName))

	// Validate the input value
	managedInstallerValue := data.IntuneManagementExtensionAsManagedInstaller.ValueString()
	if managedInstallerValue != "Enabled" && managedInstallerValue != "Disabled" {
		return fmt.Errorf("invalid value for intune_management_extension_as_managed_installer: %s. Must be 'Enabled' or 'Disabled'", managedInstallerValue)
	}

	tflog.Debug(ctx, fmt.Sprintf("Validated configuration for resource %s", ResourceName), map[string]interface{}{
		"intune_management_extension_as_managed_installer": managedInstallerValue,
	})

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing %s resource", ResourceName))

	return nil
}