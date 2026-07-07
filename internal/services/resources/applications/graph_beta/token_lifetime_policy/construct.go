package graphBetaApplicationsTokenLifetimePolicy

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// constructResource builds a TokenLifetimePolicy SDK object from the Terraform model
func constructResource(ctx context.Context, data *TokenLifetimePolicyResourceModel) (graphmodels.TokenLifetimePolicyable, error) {
	tflog.Debug(ctx, "Constructing token lifetime policy resource")

	requestBody := graphmodels.NewTokenLifetimePolicy()

	convert.FrameworkToGraphString(data.DisplayName, requestBody.SetDisplayName)
	convert.FrameworkToGraphString(data.Description, requestBody.SetDescription)

	if err := convert.FrameworkToGraphStringList(ctx, data.Definition, requestBody.SetDefinition); err != nil {
		return nil, fmt.Errorf("failed to set definition: %w", err)
	}

	convert.FrameworkToGraphBool(data.IsOrganizationDefault, requestBody.SetIsOrganizationDefault)

	if err := constructors.DebugLogGraphObject(ctx, "Constructed token lifetime policy request body", requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log request body", map[string]any{"error": err.Error()})
	}

	return requestBody, nil
}
