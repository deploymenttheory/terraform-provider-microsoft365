package graphBetaDeviceAndAppManagementIOSManagedMobileApp

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// constructResource maps the Terraform schema to the SDK model
func constructResource(ctx context.Context, data *IOSManagedMobileAppResourceModel) (graphmodels.ManagedMobileAppable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s resource from model", ResourceName))

	requestBody := graphmodels.NewManagedMobileApp()

	convert.FrameworkToGraphString(data.Version, requestBody.SetVersion)

	if data.MobileAppIdentifier != nil {
		identifier := graphmodels.NewIosMobileAppIdentifier()
		convert.FrameworkToGraphString(data.MobileAppIdentifier.BundleId, identifier.SetBundleId)
		requestBody.SetMobileAppIdentifier(identifier)
	}

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]any{
			"error": err.Error(),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing %s resource", ResourceName))

	return requestBody, nil
}
