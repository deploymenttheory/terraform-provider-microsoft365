package graphBetaWindowsPlatformScript

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// listResourceAssignments checks if a script has any assignments by querying the assignments endpoint.
// This is more reliable than using the isAssigned field which can be stale/incorrect in the API.
func (r *WindowsPlatformScriptListResource) listResourceAssignments(ctx context.Context, scriptId string) (bool, error) {
	tflog.Debug(ctx, fmt.Sprintf("Checking assignments for script: %s", scriptId))

	assignmentsResponse, err := r.client.
		DeviceManagement().
		DeviceManagementScripts().
		ByDeviceManagementScriptId(scriptId).
		Assignments().
		Get(ctx, nil)

	if err != nil {
		tflog.Warn(ctx, fmt.Sprintf("Error checking assignments for script %s: %v", scriptId, err))
		return false, err
	}

	if assignmentsResponse == nil {
		tflog.Debug(ctx, fmt.Sprintf("Script %s has nil assignments response", scriptId))
		return false, nil
	}

	assignments := assignmentsResponse.GetValue()
	hasAssignments := len(assignments) > 0
	tflog.Debug(ctx, fmt.Sprintf("Script %s has %d assignments (hasAssignments=%t)", scriptId, len(assignments), hasAssignments))
	return hasAssignments, nil
}
