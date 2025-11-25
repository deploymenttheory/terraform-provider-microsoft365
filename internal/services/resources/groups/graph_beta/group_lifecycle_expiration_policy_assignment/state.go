package graphBetaGroupLifecycleExpirationPolicyAssignment

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// MapRemoteStateToTerraform maps the group assignment state to Terraform.
// Since this is a simple association resource, we only need to ensure the ID and GroupID are set.
func MapRemoteStateToTerraform(ctx context.Context, data *GroupLifecycleExpirationPolicyAssignmentResourceModel, groupID string) {
	tflog.Debug(ctx, "Starting to map remote state to Terraform state", map[string]any{
		"groupId": groupID,
	})

	data.ID = types.StringValue(groupID)
	data.GroupID = types.StringValue(groupID)

	tflog.Debug(ctx, "Finished mapping remote state to Terraform state", map[string]any{
		"id": data.ID.ValueString(),
	})
}
