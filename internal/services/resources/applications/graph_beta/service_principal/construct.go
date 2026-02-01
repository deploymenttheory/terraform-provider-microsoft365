package graphBetaServicePrincipal

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// constructResource constructs a ServicePrincipal resource from the Terraform model
func constructResource(ctx context.Context, data *ServicePrincipalResourceModel) (graphmodels.ServicePrincipalable, error) {
	requestBody := graphmodels.NewServicePrincipal()

	// Required field: appId
	appId := data.AppID.ValueString()
	requestBody.SetAppId(&appId)

	convert.FrameworkToGraphBool(data.AccountEnabled, requestBody.SetAccountEnabled)
	convert.FrameworkToGraphBool(data.AppRoleAssignmentRequired, requestBody.SetAppRoleAssignmentRequired)

	if err := convert.FrameworkToGraphStringSet(ctx, data.Tags, requestBody.SetTags); err != nil {
		return nil, fmt.Errorf("failed to set tags: %w", err)
	}

	return requestBody, nil
}
