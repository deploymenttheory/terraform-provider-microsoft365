package graphBetaWindowsAutopatchDeploymentAudienceMembers

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodelswindowsupdates "github.com/microsoftgraph/msgraph-beta-sdk-go/models/windowsupdates"
)

func MapRemoteStateToTerraform(ctx context.Context, data *WindowsUpdateDeploymentAudienceMembersResourceModel, members []graphmodelswindowsupdates.UpdatableAssetable, exclusions []graphmodelswindowsupdates.UpdatableAssetable) {
	tflog.Debug(ctx, "Starting to map remote state to Terraform state")

	memberIDs := make([]attr.Value, 0, len(members))
	for _, member := range members {
		if member == nil {
			continue
		}

		memberType := deriveAssetType(member)
		if memberType == data.MemberType.ValueString() {
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

		exclusionType := deriveAssetType(exclusion)
		if exclusionType == data.MemberType.ValueString() {
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

	data.ID = types.StringValue(fmt.Sprintf("%s_%s", data.AudienceID.ValueString(), data.MemberType.ValueString()))

	tflog.Debug(ctx, "Finished mapping remote state to Terraform state")
}

func deriveAssetType(asset graphmodelswindowsupdates.UpdatableAssetable) string {
	odataType := asset.GetOdataType()
	if odataType == nil {
		return "azureADDevice"
	}

	parts := strings.Split(*odataType, ".")
	if len(parts) > 0 {
		typeName := parts[len(parts)-1]
		return typeName
	}

	return "azureADDevice"
}
