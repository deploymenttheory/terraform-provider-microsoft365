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

// packageID returns the Android package of an app entry, or "" for other identifiers.
func packageID(app graphmodels.ManagedMobileAppable) string {
	if app == nil {
		return ""
	}
	identifier, ok := app.GetMobileAppIdentifier().(graphmodels.AndroidMobileAppIdentifierable)
	if !ok || identifier.GetPackageId() == nil {
		return ""
	}
	return *identifier.GetPackageId()
}

// findByPackage returns the entry targeting the package, matched case-insensitively.
func findByPackage(
	apps []graphmodels.ManagedMobileAppable,
	pkg string,
) graphmodels.ManagedMobileAppable {
	for _, app := range apps {
		if strings.EqualFold(packageID(app), pkg) {
			return app
		}
	}
	return nil
}

// findByID returns the entry with the given id.
func findByID(apps []graphmodels.ManagedMobileAppable, id string) graphmodels.ManagedMobileAppable {
	for _, app := range apps {
		if app != nil && app.GetId() != nil && strings.EqualFold(*app.GetId(), id) {
			return app
		}
	}
	return nil
}
