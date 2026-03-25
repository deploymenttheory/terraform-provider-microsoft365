package graphBetaAdministrativeUnitMembership

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// MapRemoteStateToTerraform maps the remote administrative unit members to Terraform state
func MapRemoteStateToTerraform(ctx context.Context, data *AdministrativeUnitMembershipResourceModel, memberIDs []string) {
	tflog.Debug(ctx, fmt.Sprintf("Starting to map remote state to Terraform state for %s", ResourceName))

	if len(memberIDs) == 0 {
		data.Members = types.SetNull(types.StringType)
		tflog.Debug(ctx, "No members found in remote state")
		return
	}

	memberValues := make([]attr.Value, 0, len(memberIDs))
	for _, memberID := range memberIDs {
		if memberID != "" {
			memberValues = append(memberValues, types.StringValue(memberID))
			tflog.Trace(ctx, fmt.Sprintf("Mapped member ID: %s", memberID))
		}
	}

	if len(memberValues) > 0 {
		data.Members = types.SetValueMust(types.StringType, memberValues)
		tflog.Debug(ctx, fmt.Sprintf("Mapped %d members to Terraform state", len(memberValues)))
	} else {
		data.Members = types.SetNull(types.StringType)
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished mapping resource %s", ResourceName))
}
