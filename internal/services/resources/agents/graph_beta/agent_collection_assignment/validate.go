package graphBetaAgentsAgentCollectionAssignment

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

// validateRequest validates the agent collection assignment request
func validateRequest(ctx context.Context, client *msgraphbetasdk.GraphServiceClient, data *AgentCollectionAssignmentResourceModel) error {
	tflog.Debug(ctx, "Starting validation of agent collection assignment request")

	if err := validateIsAgentInstance(ctx, client, data.AgentInstanceID.ValueString()); err != nil {
		return fmt.Errorf("agent instance validation failed: %w", err)
	}

	tflog.Debug(ctx, "Successfully validated agent collection assignment request")
	return nil
}

// validateIsAgentInstance validates that the provided ID is a valid agent instance
func validateIsAgentInstance(ctx context.Context, client *msgraphbetasdk.GraphServiceClient, agentInstanceID string) error {
	if agentInstanceID == "" {
		return fmt.Errorf("agent_instance_id cannot be empty")
	}

	tflog.Debug(ctx, fmt.Sprintf("Validating agent instance exists: %s", agentInstanceID))

	instances, err := client.
		AgentRegistry().
		AgentInstances().
		Get(ctx, nil)

	if err != nil {
		return fmt.Errorf("failed to list agent instances: %w", err)
	}

	if instances == nil || instances.GetValue() == nil {
		return fmt.Errorf("no agent instances found, cannot validate agent_instance_id '%s'", agentInstanceID)
	}

	for _, instance := range instances.GetValue() {
		instanceID := instance.GetId()

		if instanceID == nil {
			continue
		}

		if *instanceID == agentInstanceID {
			tflog.Debug(ctx, fmt.Sprintf("Agent instance '%s' exists", agentInstanceID))
			return nil
		}
	}

	return fmt.Errorf("agent instance with ID '%s' does not exist", agentInstanceID)
}
