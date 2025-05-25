package graphBetaResourceOperation

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/constructors"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// constructResource constructs and returns a ResourceOperation
func constructResource(ctx context.Context, data ResourceOperationResourceModel) (graphmodels.ResourceOperationable, error) {
	tflog.Debug(ctx, "Starting resource operation construction")

	resourceOperation := graphmodels.NewResourceOperation()

	constructors.SetStringProperty(data.ResourceName, resourceOperation.SetResourceName)
	constructors.SetStringProperty(data.ActionName, resourceOperation.SetActionName)
	constructors.SetStringProperty(data.Description, resourceOperation.SetDescription)

	if err := constructors.DebugLogGraphObject(ctx, "Constructed resource operation", resourceOperation); err != nil {
		tflog.Error(ctx, "Failed to log resource operation", map[string]interface{}{
			"error": err.Error(),
		})
	}

	return resourceOperation, nil
}
