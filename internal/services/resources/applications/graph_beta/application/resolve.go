package graphBetaApplication

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// ResolveOwnerChanges compares the current state with the desired plan
// and returns which owners need to be added or removed.
func ResolveOwnerChanges(ctx context.Context, currentState, plan *ApplicationResourceModel) (ownersToAdd, ownersToRemove []string) {
	tflog.Debug(ctx, "Calculating owners to add/remove")

	ownersToAdd, ownersToRemove = calculateSetDifferences(ctx, "owners", currentState.OwnerUserIds, plan.OwnerUserIds)

	tflog.Debug(ctx, "Calculated changes", map[string]any{
		"ownersToAdd":    len(ownersToAdd),
		"ownersToRemove": len(ownersToRemove),
	})

	return ownersToAdd, ownersToRemove
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
