package graphBetaAgentsAgentCollection

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

// validateRequest validates the agent collection request by checking for duplicate display names
func validateRequest(ctx context.Context, client *msgraphbetasdk.GraphServiceClient, data *AgentCollectionResourceModel, currentID string) error {
	tflog.Debug(ctx, "Starting validation of agent collection request")

	if err := validateDisplayName(ctx, client, data.DisplayName.ValueString(), currentID); err != nil {
		return fmt.Errorf("display name validation failed: %w", err)
	}

	tflog.Debug(ctx, "Successfully validated agent collection request")
	return nil
}

// validateDisplayName validates that no other agent collection exists with the same display name
func validateDisplayName(ctx context.Context, client *msgraphbetasdk.GraphServiceClient, displayName string, currentID string) error {
	if displayName == "" {
		return fmt.Errorf("display_name cannot be empty")
	}

	tflog.Debug(ctx, fmt.Sprintf("Validating display name uniqueness: %s", displayName))

	collections, err := client.
		AgentRegistry().
		AgentCollections().
		Get(ctx, nil)

	if err != nil {
		return fmt.Errorf("failed to list agent collections: %w", err)
	}

	if collections == nil || collections.GetValue() == nil {
		tflog.Debug(ctx, "No existing agent collections found")
		return nil
	}

	for _, collection := range collections.GetValue() {
		collectionDisplayName := collection.GetDisplayName()
		collectionID := collection.GetId()

		if collectionDisplayName == nil || collectionID == nil {
			continue
		}

		// Skip the current resource when updating
		if currentID != "" && *collectionID == currentID {
			continue
		}

		if *collectionDisplayName == displayName {
			return fmt.Errorf("an agent collection with display name '%s' already exists (ID: %s)", displayName, *collectionID)
		}
	}

	tflog.Debug(ctx, fmt.Sprintf("Display name '%s' is unique", displayName))
	return nil
}
