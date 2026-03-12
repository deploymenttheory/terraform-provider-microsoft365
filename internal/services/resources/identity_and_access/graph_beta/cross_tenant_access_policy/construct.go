package graphBetaCrossTenantAccessPolicy

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// constructResource converts the Terraform resource model to a CrossTenantAccessPolicy SDK request body.
func constructResource(ctx context.Context, data *CrossTenantAccessPolicyResourceModel) (graphmodels.CrossTenantAccessPolicyable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s resource from model", ResourceName))

	requestBody := graphmodels.NewCrossTenantAccessPolicy()

	if !data.AllowedCloudEndpoints.IsNull() && !data.AllowedCloudEndpoints.IsUnknown() {
		if err := convert.FrameworkToGraphStringSet(ctx, data.AllowedCloudEndpoints, requestBody.SetAllowedCloudEndpoints); err != nil {
			return nil, fmt.Errorf("failed to set allowed_cloud_endpoints: %w", err)
		}
	}

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]any{
			"error": err.Error(),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing %s resource", ResourceName))

	return requestBody, nil
}

// constructRestoreDefaultsBody builds the PATCH body used when restore_defaults_on_destroy is true.
// It sets allowed_cloud_endpoints to the default endpoints, restoring the service default.
func constructRestoreDefaultsBody() graphmodels.CrossTenantAccessPolicyable {
	requestBody := graphmodels.NewCrossTenantAccessPolicy()
	requestBody.SetAllowedCloudEndpoints([]string{
		"microsoftonline.us",
		"partner.microsoftonline.cn",
	})
	return requestBody
}
