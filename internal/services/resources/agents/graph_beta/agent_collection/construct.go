package graphBetaAgentCollection

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// constructResource builds the request body for creating or updating an agent collection.
// currentID should be empty string for create, or the resource ID for update.
func constructResource(ctx context.Context, client *msgraphbetasdk.GraphServiceClient, data *AgentCollectionResourceModel, currentID string) (graphmodels.AgentCollectionable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s resource from model", ResourceName))

	if err := validateRequest(ctx, client, data, currentID); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	requestBody := graphmodels.NewAgentCollection()

	convert.FrameworkToGraphString(data.DisplayName, requestBody.SetDisplayName)
	convert.FrameworkToGraphString(data.Description, requestBody.SetDescription)
	convert.FrameworkToGraphString(data.ManagedBy, requestBody.SetManagedBy)
	convert.FrameworkToGraphString(data.OriginatingStore, requestBody.SetOriginatingStore)

	if err := convert.FrameworkToGraphStringSet(ctx, data.OwnerIds, requestBody.SetOwnerIds); err != nil {
		return nil, fmt.Errorf("failed to set owner_ids: %w", err)
	}

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]any{
			"error": err.Error(),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing %s resource", ResourceName))

	return requestBody, nil
}
