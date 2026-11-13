package graphBetaDeviceAndAppManagementAndroidManagedMobileApp

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
)

// MapRemoteStateToTerraform maps the base properties of an AndroidManagedMobileAppResourceModel to a Terraform state.
func MapRemoteStateToTerraform(
	ctx context.Context,
	data *AndroidManagedMobileAppResourceModel,
	remoteResource graphmodels.ManagedMobileAppable,
) {
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
		if androidId, ok := appIdentifier.(graphmodels.AndroidMobileAppIdentifierable); ok {
			remotePackage := convert.GraphToFrameworkString(androidId.GetPackageId())
			// Keep the configured casing when Graph returns the same package differently cased.
			if data.MobileAppIdentifier == nil ||
				!strings.EqualFold(
					data.MobileAppIdentifier.PackageId.ValueString(),
					remotePackage.ValueString(),
				) {
				data.MobileAppIdentifier = &AndroidMobileAppIdentifierResourceModel{
					PackageId: remotePackage,
				}
			}
		}
	}

	tflog.Debug(
		ctx,
		fmt.Sprintf("Finished stating resource %s with id %s", ResourceName, data.ID.ValueString()),
	)
}
