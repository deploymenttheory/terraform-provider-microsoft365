package graphBetaDeviceAndAppManagementAndroidManagedMobileApp

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/deviceappmanagement"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
)

// constructResource maps the Terraform schema to the SDK model
func constructResource(
	ctx context.Context,
	data *AndroidManagedMobileAppResourceModel,
) (graphmodels.ManagedMobileAppable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s resource from model", ResourceName))

	requestBody := graphmodels.NewManagedMobileApp()

	// version is assigned by Intune, so it is never sent.
	if data.MobileAppIdentifier != nil {
		identifier := graphmodels.NewAndroidMobileAppIdentifier()
		convert.FrameworkToGraphString(data.MobileAppIdentifier.PackageId, identifier.SetPackageId)
		requestBody.SetMobileAppIdentifier(identifier)
	}

	if err := constructors.DebugLogGraphObject(
		ctx,
		fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName),
		requestBody,
	); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]any{
			"error": err.Error(),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing %s resource", ResourceName))

	return requestBody, nil
}

// constructTargetAppsRequest builds the targetApps replacement request. Graph assigns each
// entry's id and version itself, so Android entries are rebuilt with only their package ID.
func constructTargetAppsRequest(apps []graphmodels.ManagedMobileAppable) *deviceappmanagement.ManagedAppPoliciesItemTargetAppsPostRequestBody {
	entries := make([]graphmodels.ManagedMobileAppable, 0, len(apps))
	for _, app := range apps {
		identifier, ok := app.GetMobileAppIdentifier().(graphmodels.AndroidMobileAppIdentifierable)
		if !ok || identifier.GetPackageId() == nil {
			entries = append(entries, app)
			continue
		}
		packageID := graphmodels.NewAndroidMobileAppIdentifier()
		packageID.SetPackageId(identifier.GetPackageId())
		entry := graphmodels.NewManagedMobileApp()
		entry.SetMobileAppIdentifier(packageID)
		entries = append(entries, entry)
	}
	body := deviceappmanagement.NewManagedAppPoliciesItemTargetAppsPostRequestBody()
	body.SetApps(entries)
	return body
}
