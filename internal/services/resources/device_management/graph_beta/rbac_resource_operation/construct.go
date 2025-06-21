package graphBetaRBACResourceOperation

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// constructResource constructs and returns a RBACResourceOperation
func constructResource(ctx context.Context, data RBACResourceOperationResourceModel) (graphmodels.ResourceOperationable, error) {
	tflog.Debug(ctx, "Starting resource operation construction")

	RBACResourceOperation := graphmodels.NewResourceOperation()

	convert.FrameworkToGraphString(data.ResourceName, RBACResourceOperation.SetResourceName)
	convert.FrameworkToGraphString(data.ActionName, RBACResourceOperation.SetActionName)
	convert.FrameworkToGraphString(data.Description, RBACResourceOperation.SetDescription)

	if err := constructors.DebugLogGraphObject(ctx, "Constructed resource operation", RBACResourceOperation); err != nil {
		tflog.Error(ctx, "Failed to log resource operation", map[string]interface{}{
			"error": err.Error(),
		})
	}

	return RBACResourceOperation, nil
}
