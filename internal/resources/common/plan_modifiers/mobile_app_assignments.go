package planmodifiers

import (
	"context"
	"reflect"
	"sort"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func MobileAppAssignmentsListModifier() planmodifier.List {
	return &mobileAppAssignmentsListPlanModifier{}
}

type mobileAppAssignmentsListPlanModifier struct{}

func (m *mobileAppAssignmentsListPlanModifier) Description(ctx context.Context) string {
	return "Ensures correct ordering of mobile app assignments by matching content rather than relying on list position"
}

func (m *mobileAppAssignmentsListPlanModifier) MarkdownDescription(ctx context.Context) string {
	return "Ensures correct ordering of mobile app assignments by matching content rather than relying on list position"
}

func (m *mobileAppAssignmentsListPlanModifier) PlanModifyList(ctx context.Context, req planmodifier.ListRequest, resp *planmodifier.ListResponse) {
	tflog.Debug(ctx, "Entered MobileAppAssignmentsListPlanModifier method")

	if req.StateValue.IsNull() || req.PlanValue.IsNull() {
		tflog.Debug(ctx, "State or plan value is null, skipping assignment modification")
		return
	}

	var planAssignments, stateAssignments []types.Object

	diags := req.PlanValue.ElementsAs(ctx, &planAssignments, false)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	diags = req.StateValue.ElementsAs(ctx, &stateAssignments, false)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	// Sort assignments by the same logic as the stating logic
	// which is the same as the resp from intune.
	sortMobileAppAssignments(planAssignments)
	sortMobileAppAssignments(stateAssignments)

	// Strip assignment IDs as they are computed in the schema and so will never be in the plan.
	stripAssignmentIds(ctx, req, planAssignments)
	stripAssignmentIds(ctx, req, stateAssignments)

	tflog.Debug(ctx, "Sorted plan assignments", map[string]interface{}{"assignments": planAssignments})
	tflog.Debug(ctx, "Sorted state assignments", map[string]interface{}{"assignments": stateAssignments})

	if reflect.DeepEqual(planAssignments, stateAssignments) {
		resp.PlanValue = req.StateValue
		tflog.Debug(ctx, "Assignments match after sorting; reusing state ordering")
	} else {
		tflog.Debug(ctx, "Assignments differ after sorting; using plan value")
	}
}

// sortMobileAppAssignments sorts a slice of mobile app assignments
// The sort order is as follows:
// 1. First tier: Sort by intent alphabetically
// 2. Second tier: Within same intent, sort by target_type alphabetically
// 3. Third tier: Within same target_type, sort by group_id alphabetically
func sortMobileAppAssignments(assignments []types.Object) {
	sort.SliceStable(assignments, func(i, j int) bool {
		// First tier: Sort by intent alphabetically
		intentI := getStringAttr(assignments[i], "intent")
		intentJ := getStringAttr(assignments[j], "intent")
		if intentI != intentJ {
			return intentI < intentJ
		}

		// Second tier: Sort by target_type alphabetically
		targetTypeI := getNestedStringAttr(assignments[i], "target", "target_type")
		targetTypeJ := getNestedStringAttr(assignments[j], "target", "target_type")
		if targetTypeI != targetTypeJ {
			return targetTypeI < targetTypeJ
		}

		// Third tier: Sort by group_id alphabetically if neither are null
		groupIDI := getNestedStringAttr(assignments[i], "target", "group_id")
		groupIDJ := getNestedStringAttr(assignments[j], "target", "group_id")

		// Handle potential empty strings from nulls uniformly
		if groupIDI != groupIDJ {
			return groupIDI < groupIDJ
		}

		// Final fallback to ensure stability
		return i < j
	})
}

func stripAssignmentIds(ctx context.Context, req planmodifier.ListRequest, assignments []types.Object) {
	for i := range assignments {
		attrs := assignments[i].Attributes()
		attrs["id"] = types.StringNull()
		assignments[i] = types.ObjectValueMust(req.PlanValue.ElementType(ctx).(types.ObjectType).AttrTypes, attrs)
	}
}

func getStringAttr(obj types.Object, attrName string) string {
	attr, ok := obj.Attributes()[attrName]
	if !ok || attr.IsNull() || attr.IsUnknown() {
		return ""
	}

	val, ok := attr.(types.String)
	if !ok {
		return ""
	}

	return val.ValueString()
}

func getNestedStringAttr(obj types.Object, nestedName, attrName string) string {
	nestedAttr, ok := obj.Attributes()[nestedName]
	if !ok || nestedAttr.IsNull() || nestedAttr.IsUnknown() {
		return ""
	}

	nestedObj, ok := nestedAttr.(types.Object)
	if !ok {
		return ""
	}

	return getStringAttr(nestedObj, attrName)
}
