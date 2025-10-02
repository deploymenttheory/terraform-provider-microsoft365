package graphBetaAuthenticationContext

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

// constructResource uses the Microsoft Graph SDK models directly instead of raw HTTP calls
func constructResource(ctx context.Context, client *msgraphbetasdk.GraphServiceClient, data *AuthenticationContextResourceModel, isCreate bool) (graphmodels.AuthenticationContextClassReferenceable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s resource from model", ResourceName))

	if err := validateRequest(ctx, client, data, isCreate); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	requestBody := graphmodels.NewAuthenticationContextClassReference()

	convert.FrameworkToGraphString(data.ID, requestBody.SetId)
	convert.FrameworkToGraphString(data.DisplayName, requestBody.SetDisplayName)
	convert.FrameworkToGraphString(data.Description, requestBody.SetDescription)
	convert.FrameworkToGraphBool(data.IsAvailable, requestBody.SetIsAvailable)

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]any{
			"error": err.Error(),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing %s resource using Graph SDK", ResourceName))

	return requestBody, nil
}
