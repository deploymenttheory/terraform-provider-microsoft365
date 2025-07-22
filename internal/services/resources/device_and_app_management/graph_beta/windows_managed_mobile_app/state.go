package graphBetaDeviceAndAppManagementWindowsManagedMobileApp

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// MapRemoteStateToTerraform maps the base properties of a WindowsManagedMobileAppResourceModel to a Terraform state.
func MapRemoteStateToTerraform(ctx context.Context, data *WindowsManagedMobileAppResourceModel, remoteResource graphmodels.ManagedMobileAppable) {
	if remoteResource == nil {
		tflog.Debug(ctx, "Remote resource is nil")
		return
	}

	tflog.Debug(ctx, "Starting to map remote state to Terraform state", map[string]any{
		"resourceId": convert.GraphToFrameworkString(remoteResource.GetId()).ValueString(),
	})

	data.ID = convert.GraphToFrameworkString(remoteResource.GetId())
	data.Version = convert.GraphToFrameworkString(remoteResource.GetVersion())

	if appIdentifier := remoteResource.GetMobileAppIdentifier(); appIdentifier != nil {
		if windowsId, ok := appIdentifier.(*graphmodels.WindowsAppIdentifier); ok {
			data.MobileAppIdentifier = &WindowsMobileAppIdentifierResourceModel{
				WindowsAppId: convert.GraphToFrameworkString(windowsId.GetWindowsAppId()),
			}
		}
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished stating resource %s with id %s", ResourceName, data.ID.ValueString()))
}