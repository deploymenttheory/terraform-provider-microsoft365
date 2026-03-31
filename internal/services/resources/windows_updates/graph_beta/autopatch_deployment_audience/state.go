package graphBetaWindowsUpdatesAutopatchDeploymentAudience

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodelswindowsupdates "github.com/microsoftgraph/msgraph-beta-sdk-go/models/windowsupdates"
)

func MapRemoteStateToTerraform(
	ctx context.Context,
	data *WindowsUpdatesAutopatchDeploymentAudienceResourceModel,
	remoteResource graphmodelswindowsupdates.DeploymentAudienceable,
	members []graphmodelswindowsupdates.UpdatableAssetable,
	exclusions []graphmodelswindowsupdates.UpdatableAssetable,
) {
	if remoteResource == nil {
		tflog.Debug(ctx, "Remote resource is nil")
		return
	}

	tflog.Debug(ctx, "Starting to map remote state to Terraform state", map[string]any{
		"resourceId": remoteResource.GetId(),
	})

	data.ID = types.StringValue(*remoteResource.GetId())

	// Default member_type to "azureADDevice" if not set
	if data.MemberType.IsNull() || data.MemberType.IsUnknown() || data.MemberType.ValueString() == "" {
		data.MemberType = types.StringValue("azureADDevice")
	}

	memberType := data.MemberType.ValueString()

	memberIDs := make([]attr.Value, 0, len(members))
	for _, member := range members {
		if member == nil {
			continue
		}
		if deriveAssetType(member) == memberType {
			if id := member.GetId(); id != nil {
				memberIDs = append(memberIDs, types.StringValue(*id))
			}
		}
	}

	if len(memberIDs) > 0 {
		data.Members = types.SetValueMust(types.StringType, memberIDs)
	} else if !data.Members.IsNull() {
		data.Members = types.SetValueMust(types.StringType, []attr.Value{})
	}

	exclusionIDs := make([]attr.Value, 0, len(exclusions))
	for _, exclusion := range exclusions {
		if exclusion == nil {
			continue
		}
		if deriveAssetType(exclusion) == memberType {
			if id := exclusion.GetId(); id != nil {
				exclusionIDs = append(exclusionIDs, types.StringValue(*id))
			}
		}
	}

	if len(exclusionIDs) > 0 {
		data.Exclusions = types.SetValueMust(types.StringType, exclusionIDs)
	} else if !data.Exclusions.IsNull() {
		data.Exclusions = types.SetValueMust(types.StringType, []attr.Value{})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished mapping remote state to Terraform state: %d members, %d exclusions", len(memberIDs), len(exclusionIDs)))
}

func deriveAssetType(asset graphmodelswindowsupdates.UpdatableAssetable) string {
	odataType := asset.GetOdataType()
	if odataType == nil {
		return "azureADDevice"
	}

	parts := strings.Split(*odataType, ".")
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}

	return "azureADDevice"
}
