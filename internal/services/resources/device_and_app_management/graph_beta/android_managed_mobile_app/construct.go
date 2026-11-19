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

// writableApp rebuilds an entry without the server-assigned id and version, which
// targetApps derives itself.
func writableApp(app graphmodels.ManagedMobileAppable) graphmodels.ManagedMobileAppable {
	pkg := packageID(app)
	if pkg == "" {
		return app
	}
	identifier := graphmodels.NewAndroidMobileAppIdentifier()
	identifier.SetPackageId(&pkg)
	entry := graphmodels.NewManagedMobileApp()
	entry.SetMobileAppIdentifier(identifier)
	return entry
}

// constructPolicyApps builds the complete targetApps replacement request.
func constructPolicyApps(apps []graphmodels.ManagedMobileAppable) *deviceappmanagement.ManagedAppPoliciesItemTargetAppsPostRequestBody {
	body := deviceappmanagement.NewManagedAppPoliciesItemTargetAppsPostRequestBody()
	body.SetApps(apps)
	return body
}
