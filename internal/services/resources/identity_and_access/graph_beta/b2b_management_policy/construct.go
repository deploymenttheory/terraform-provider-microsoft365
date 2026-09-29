package graphBetaIdentityAndAccessB2bManagementPolicy

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// constructResource builds the request body. description must never be set (see the model).
func constructResource(ctx context.Context, data *B2bManagementPolicyResourceModel) (graphmodels.B2bManagementPolicyable, error) {
	tflog.Debug(ctx, "Constructing B2B management policy resource")

	requestBody := graphmodels.NewB2bManagementPolicy()

	convert.FrameworkToGraphString(data.DisplayName, requestBody.SetDisplayName)

	if err := convert.FrameworkToGraphStringList(ctx, data.Definition, requestBody.SetDefinition); err != nil {
		return nil, fmt.Errorf("failed to set definition: %w", err)
	}

	convert.FrameworkToGraphBool(data.IsOrganizationDefault, requestBody.SetIsOrganizationDefault)

	if err := constructors.DebugLogGraphObject(ctx, "Constructed B2B management policy request body", requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log request body", map[string]any{"error": err.Error()})
	}

	return requestBody, nil
}
