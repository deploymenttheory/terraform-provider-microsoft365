package graphBetaGroupLifecyclePolicy

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// constructResource maps the Terraform schema to the SDK model
func constructResource(ctx context.Context, data *GroupLifecyclePolicyResourceModel) (graphmodels.GroupLifecyclePolicyable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s resource from model", ResourceName))

	requestBody := graphmodels.NewGroupLifecyclePolicy()

	convert.FrameworkToGraphString(data.AlternateNotificationEmails, requestBody.SetAlternateNotificationEmails)
	convert.FrameworkToGraphInt32(data.GroupLifetimeInDays, requestBody.SetGroupLifetimeInDays)
	convert.FrameworkToGraphString(data.ManagedGroupTypes, requestBody.SetManagedGroupTypes)

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]interface{}{
			"error": err.Error(),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing %s resource", ResourceName))

	return requestBody, nil
}
