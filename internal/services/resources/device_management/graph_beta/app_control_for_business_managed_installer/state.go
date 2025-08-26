package graphBetaAppControlForBusinessManagedInstaller

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// MapRemoteStateToTerraform maps the properties of a WindowsManagementApp to a Terraform state.
func MapRemoteStateToTerraform(ctx context.Context, data *AppControlForBusinessManagedInstallerResourceModel, remoteResource graphmodels.WindowsManagementAppable) error {
	if remoteResource == nil {
		tflog.Debug(ctx, "Remote resource is nil")
		return fmt.Errorf("remote resource is nil")
	}

	tflog.Debug(ctx, "Starting to map remote state to Terraform state", map[string]any{
		"resourceId": convert.GraphToFrameworkString(remoteResource.GetId()).ValueString(),
	})

	data.ID = convert.GraphToFrameworkString(remoteResource.GetId())
	data.AvailableVersion = convert.GraphToFrameworkString(remoteResource.GetAvailableVersion())
	data.ManagedInstallerConfiguredDateTime = convert.GraphToFrameworkString(remoteResource.GetManagedInstallerConfiguredDateTime())

	// Map managed installer status to Terraform field
	managedInstaller := remoteResource.GetManagedInstaller()
	if managedInstaller != nil {
		switch *managedInstaller {
		case graphmodels.ENABLED_MANAGEDINSTALLERSTATUS:
			data.IntuneManagementExtensionAsManagedInstaller = types.StringValue("Enabled")
		case graphmodels.DISABLED_MANAGEDINSTALLERSTATUS:
			data.IntuneManagementExtensionAsManagedInstaller = types.StringValue("Disabled")
		default:
			return fmt.Errorf("unknown managed installer status: %v", *managedInstaller)
		}
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished stating resource %s with id %s", ResourceName, data.ID.ValueString()))
	return nil
}
