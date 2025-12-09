package graphBetaAgentIdentity

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// ResolveSponsorAndOwnerChanges compares the current state with the desired plan
// and returns which sponsors and owners need to be added or removed.
func ResolveSponsorAndOwnerChanges(ctx context.Context, currentState, plan *AgentIdentityResourceModel) (sponsorsToAdd, sponsorsToRemove, ownersToAdd, ownersToRemove []string) {
	tflog.Debug(ctx, "Calculating sponsors and owners to add/remove")

	sponsorsToAdd, sponsorsToRemove = calculateSetDifferences(ctx, "sponsors", currentState.SponsorIds, plan.SponsorIds)
	ownersToAdd, ownersToRemove = calculateSetDifferences(ctx, "owners", currentState.OwnerIds, plan.OwnerIds)

	tflog.Debug(ctx, "Calculated changes", map[string]any{
		"sponsorsToAdd":    len(sponsorsToAdd),
		"sponsorsToRemove": len(sponsorsToRemove),
		"ownersToAdd":      len(ownersToAdd),
		"ownersToRemove":   len(ownersToRemove),
	})

	return sponsorsToAdd, sponsorsToRemove, ownersToAdd, ownersToRemove
}

// calculateSetDifferences compares current and planned sets to determine what to add and remove
func calculateSetDifferences(ctx context.Context, resourceType string, current, plan types.Set) (toAdd, toRemove []string) {
	if current.IsNull() && plan.IsNull() {
		return nil, nil
	}

	var currentItems, planItems []string

	if !current.IsNull() && !current.IsUnknown() {
		current.ElementsAs(ctx, &currentItems, false)
	}

	if !plan.IsNull() && !plan.IsUnknown() {
		plan.ElementsAs(ctx, &planItems, false)
	}

	currentMap := make(map[string]bool)
	for _, item := range currentItems {
		currentMap[item] = true
	}

	planMap := make(map[string]bool)
	for _, item := range planItems {
		planMap[item] = true
	}

	for _, item := range planItems {
		if !currentMap[item] {
			toAdd = append(toAdd, item)
		}
	}

	for _, item := range currentItems {
		if !planMap[item] {
			toRemove = append(toRemove, item)
		}
	}

	tflog.Debug(ctx, fmt.Sprintf("Calculated %s differences", resourceType), map[string]any{
		"toAdd":    len(toAdd),
		"toRemove": len(toRemove),
	})

	return toAdd, toRemove
}
