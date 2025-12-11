package graphBetaAgentsAgentCollectionAssignment

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// SetCompositeID sets the composite ID for the assignment resource.
func SetCompositeID(data *AgentCollectionAssignmentResourceModel) {
	agentInstanceID := data.AgentInstanceID.ValueString()
	agentCollectionID := data.AgentCollectionID.ValueString()
	data.ID = types.StringValue(fmt.Sprintf("%s/%s", agentInstanceID, agentCollectionID))
}

// MapRemoteResourceStateToTerraform checks if the agent instance is a member of the collection
// and updates the Terraform state accordingly.
func MapRemoteResourceStateToTerraform(ctx context.Context, data *AgentCollectionAssignmentResourceModel, membersResponse graphmodels.AgentInstanceCollectionResponseable) bool {
	if membersResponse == nil {
		tflog.Debug(ctx, "Members response is nil")
		return false
	}

	agentInstanceID := data.AgentInstanceID.ValueString()
	agentCollectionID := data.AgentCollectionID.ValueString()

	tflog.Debug(ctx, fmt.Sprintf("Checking if agent instance %s is a member of collection %s", agentInstanceID, agentCollectionID))

	// Check if the agent instance is a member of the collection
	found := false
	for _, member := range membersResponse.GetValue() {
		if member.GetId() != nil && *member.GetId() == agentInstanceID {
			found = true
			break
		}
	}

	if found {
		SetCompositeID(data)
		tflog.Debug(ctx, fmt.Sprintf("Agent instance %s is a member of collection %s", agentInstanceID, agentCollectionID))
	} else {
		tflog.Debug(ctx, fmt.Sprintf("Agent instance %s is not a member of collection %s", agentInstanceID, agentCollectionID))
	}

	return found
}
