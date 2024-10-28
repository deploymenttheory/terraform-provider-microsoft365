package graphBetaWindowsSettingsCatalog

import (
	"context"
	"sort"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ModifyPlan is the main entry point for diff suppression
func (r *WindowsSettingsCatalogResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	// Skip if creating or deleting
	if req.State.Raw.IsNull() || req.Plan.Raw.IsNull() {
		return
	}

	var plan, state WindowsSettingsCatalogProfileResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Call individual suppressors
	suppressAssignmentDiffs(&plan, &state)
	// Add other suppressors as needed:
	// suppressSettingsDiffs(&plan, &state)
	// suppressMetadataDiffs(&plan, &state)

	resp.Plan.Set(ctx, &plan)
}

// suppressAssignmentDiffs handles all assignment-related diff suppression
func suppressAssignmentDiffs(plan, state *WindowsSettingsCatalogProfileResourceModel) {
	if plan.Assignments == nil || state.Assignments == nil {
		return
	}

	// Handle exclude_group_ids ordering
	if hasSameGroupIds(plan.Assignments.ExcludeGroupIds, state.Assignments.ExcludeGroupIds) {
		plan.Assignments.ExcludeGroupIds = state.Assignments.ExcludeGroupIds
	}

	// Handle include_groups ordering and content
	if hasSameIncludeGroups(plan.Assignments.IncludeGroups, state.Assignments.IncludeGroups) {
		plan.Assignments.IncludeGroups = state.Assignments.IncludeGroups
	}
}

// hasSameGroupIds checks if two slices of group IDs contain the same elements regardless of order
func hasSameGroupIds(plan, state []types.String) bool {
	if len(plan) != len(state) {
		return false
	}

	// Create maps to compare content regardless of order
	planMap := make(map[string]struct{})
	stateMap := make(map[string]struct{})

	for _, p := range plan {
		planMap[p.ValueString()] = struct{}{}
	}

	for _, s := range state {
		stateMap[s.ValueString()] = struct{}{}
	}

	// Compare maps
	if len(planMap) != len(stateMap) {
		return false
	}

	for k := range planMap {
		if _, ok := stateMap[k]; !ok {
			return false
		}
	}

	return true
}

// hasSameIncludeGroups checks if two slices of include groups contain the same elements regardless of order
func hasSameIncludeGroups(plan, state []IncludeGroup) bool {
	if len(plan) != len(state) {
		return false
	}

	// Create maps of group configurations using a composite key
	planMap := make(map[string]IncludeGroup)
	stateMap := make(map[string]IncludeGroup)

	// Build maps for both plan and state
	for _, p := range plan {
		// Create composite key from all fields
		key := createIncludeGroupKey(p)
		planMap[key] = p
	}

	for _, s := range state {
		key := createIncludeGroupKey(s)
		stateMap[key] = s
	}

	// Check if all configurations exist in both maps
	for key := range planMap {
		if _, exists := stateMap[key]; !exists {
			return false
		}
	}

	return true
}

// createIncludeGroupKey creates a unique key for an IncludeGroup by combining all its values
func createIncludeGroupKey(group IncludeGroup) string {
	// Sort the values to ensure consistent key generation
	values := []string{
		group.GroupId.ValueString(),
		group.IncludeGroupsFilterId.ValueString(),
		group.IncludeGroupsFilterType.ValueString(),
	}
	sort.Strings(values)
	return strings.Join(values, "|")
}
