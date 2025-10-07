package graphBetaAppleConfiguratorEnrollmentPolicy

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

func constructEnrollmentProfile(ctx context.Context, data *AppleConfiguratorEnrollmentPolicyResourceModel) (graphmodels.EnrollmentProfileable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s resource from model", ResourceName))

	requestBody := graphmodels.NewEnrollmentProfile()
	convert.FrameworkToGraphString(data.DisplayName, requestBody.SetDisplayName)
	convert.FrameworkToGraphString(data.Description, requestBody.SetDescription)
	convert.FrameworkToGraphBool(data.RequiresUserAuthentication, requestBody.SetRequiresUserAuthentication)
	convert.FrameworkToGraphBool(data.EnableAuthenticationViaCompanyPortal, requestBody.SetEnableAuthenticationViaCompanyPortal)
	convert.FrameworkToGraphBool(data.RequireCompanyPortalOnSetupAssistantEnrolledDevices, requestBody.SetRequireCompanyPortalOnSetupAssistantEnrolledDevices)

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]any{
			"error": err.Error(),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing %s resource", ResourceName))

	return requestBody, nil
}
