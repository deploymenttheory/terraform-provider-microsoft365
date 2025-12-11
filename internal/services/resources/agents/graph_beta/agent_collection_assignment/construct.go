package graphBetaAgentsAgentCollectionAssignment

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// constructResource builds the request body for adding an agent instance to a collection.
func constructResource(ctx context.Context, client *msgraphbetasdk.GraphServiceClient, data *AgentCollectionAssignmentResourceModel) (graphmodels.ReferenceCreateable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s resource from model", ResourceName))

	if err := validateRequest(ctx, client, data); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	requestBody := graphmodels.NewReferenceCreate()

	agentInstanceID := data.AgentInstanceID.ValueString()
	odataId := fmt.Sprintf("https://graph.microsoft.com/beta/agentRegistry/agentInstances('%s')", agentInstanceID)
	requestBody.SetOdataId(&odataId)

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]any{
			"error": err.Error(),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing %s resource", ResourceName))

	return requestBody, nil
}
