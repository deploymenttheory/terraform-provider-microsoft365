package graphBetaAdministrativeUnitRoleAssignment

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

// validateRequest validates the role assignment request before creation.
// Verifies that the role member ID corresponds to a valid directory object.
func validateRequest(ctx context.Context, client *msgraphbetasdk.GraphServiceClient, data *AdministrativeUnitRoleAssignmentResourceModel, diagnostics *diag.Diagnostics) bool {
	memberID := data.RoleMemberID.ValueString()

	tflog.Debug(ctx, fmt.Sprintf("Validating role member directory object: %s", memberID))

	directoryObj, err := client.
		DirectoryObjects().
		ByDirectoryObjectId(memberID).
		Get(ctx, nil)

	if err != nil {
		diagnostics.AddError(
			"Invalid role_member_id",
			fmt.Sprintf("The directory object with ID '%s' could not be found or accessed: %s", memberID, err.Error()),
		)
		return false
	}

	if directoryObj == nil || directoryObj.GetId() == nil {
		diagnostics.AddError(
			"Invalid role_member_id",
			fmt.Sprintf("The directory object with ID '%s' does not exist.", memberID),
		)
		return false
	}

	tflog.Debug(ctx, fmt.Sprintf("Successfully validated role member directory object: %s", memberID))
	return true
}
