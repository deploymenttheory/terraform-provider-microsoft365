package graphBetaWindowsUpdatesAutopatchDeploymentAudience

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphadmin "github.com/microsoftgraph/msgraph-beta-sdk-go/admin"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models/windowsupdates"
)

func constructAudienceResource(ctx context.Context) (windowsupdates.DeploymentAudienceable, error) {
	tflog.Debug(ctx, "Constructing deployment audience resource")

	requestBody := windowsupdates.NewDeploymentAudience()

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]any{
			"error": err.Error(),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing deployment audience for %s resource", ResourceName))
	return requestBody, nil
}

func createUpdatableAsset(id string, memberType string) windowsupdates.UpdatableAssetable {
	odataType := "#microsoft.graph.windowsUpdates." + memberType

	switch memberType {
	case "azureADDevice":
		asset := windowsupdates.NewAzureADDevice()
		asset.SetOdataType(&odataType)
		asset.SetId(&id)
		return asset
	case "updatableAssetGroup":
		asset := windowsupdates.NewUpdatableAssetGroup()
		asset.SetOdataType(&odataType)
		asset.SetId(&id)
		return asset
	default:
		asset := windowsupdates.NewUpdatableAsset()
		asset.SetOdataType(&odataType)
		asset.SetId(&id)
		return asset
	}
}

func constructAddMembersRequest(ctx context.Context, plan *WindowsUpdatesAutopatchDeploymentAudienceResourceModel) (graphadmin.WindowsUpdatesDeploymentAudiencesItemMicrosoftGraphWindowsUpdatesUpdateAudienceUpdateAudiencePostRequestBodyable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing add members request for %s resource", ResourceName))

	requestBody := graphadmin.NewWindowsUpdatesDeploymentAudiencesItemMicrosoftGraphWindowsUpdatesUpdateAudienceUpdateAudiencePostRequestBody()
	memberType := plan.MemberType.ValueString()

	var addMembers []windowsupdates.UpdatableAssetable
	if !plan.Members.IsNull() && !plan.Members.IsUnknown() {
		elements := plan.Members.Elements()
		addMembers = make([]windowsupdates.UpdatableAssetable, 0, len(elements))
		for _, elem := range elements {
			if strVal, ok := elem.(types.String); ok {
				asset := createUpdatableAsset(strVal.ValueString(), memberType)
				addMembers = append(addMembers, asset)
			}
		}
	}

	var addExclusions []windowsupdates.UpdatableAssetable
	if !plan.Exclusions.IsNull() && !plan.Exclusions.IsUnknown() {
		elements := plan.Exclusions.Elements()
		addExclusions = make([]windowsupdates.UpdatableAssetable, 0, len(elements))
		for _, elem := range elements {
			if strVal, ok := elem.(types.String); ok {
				asset := createUpdatableAsset(strVal.ValueString(), memberType)
				addExclusions = append(addExclusions, asset)
			}
		}
	}

	if len(addMembers) > 0 {
		requestBody.SetAddMembers(addMembers)
	}
	if len(addExclusions) > 0 {
		requestBody.SetAddExclusions(addExclusions)
	}

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s (add members)", ResourceName), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]any{
			"error": err.Error(),
		})
	}

	return requestBody, nil
}

func constructUpdateMembersRequest(ctx context.Context, plan *WindowsUpdatesAutopatchDeploymentAudienceResourceModel, state *WindowsUpdatesAutopatchDeploymentAudienceResourceModel) (graphadmin.WindowsUpdatesDeploymentAudiencesItemMicrosoftGraphWindowsUpdatesUpdateAudienceUpdateAudiencePostRequestBodyable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing update members request for %s resource", ResourceName))

	requestBody := graphadmin.NewWindowsUpdatesDeploymentAudiencesItemMicrosoftGraphWindowsUpdatesUpdateAudienceUpdateAudiencePostRequestBody()
	memberType := plan.MemberType.ValueString()

	planMemberIDs := make(map[string]bool)
	if !plan.Members.IsNull() && !plan.Members.IsUnknown() {
		for _, elem := range plan.Members.Elements() {
			if strVal, ok := elem.(types.String); ok {
				planMemberIDs[strVal.ValueString()] = true
			}
		}
	}

	stateMemberIDs := make(map[string]bool)
	if !state.Members.IsNull() && !state.Members.IsUnknown() {
		for _, elem := range state.Members.Elements() {
			if strVal, ok := elem.(types.String); ok {
				stateMemberIDs[strVal.ValueString()] = true
			}
		}
	}

	var addMembers []windowsupdates.UpdatableAssetable
	for id := range planMemberIDs {
		if !stateMemberIDs[id] {
			addMembers = append(addMembers, createUpdatableAsset(id, memberType))
		}
	}

	var removeMembers []windowsupdates.UpdatableAssetable
	for id := range stateMemberIDs {
		if !planMemberIDs[id] {
			removeMembers = append(removeMembers, createUpdatableAsset(id, memberType))
		}
	}

	planExclusionIDs := make(map[string]bool)
	if !plan.Exclusions.IsNull() && !plan.Exclusions.IsUnknown() {
		for _, elem := range plan.Exclusions.Elements() {
			if strVal, ok := elem.(types.String); ok {
				planExclusionIDs[strVal.ValueString()] = true
			}
		}
	}

	stateExclusionIDs := make(map[string]bool)
	if !state.Exclusions.IsNull() && !state.Exclusions.IsUnknown() {
		for _, elem := range state.Exclusions.Elements() {
			if strVal, ok := elem.(types.String); ok {
				stateExclusionIDs[strVal.ValueString()] = true
			}
		}
	}

	var addExclusions []windowsupdates.UpdatableAssetable
	for id := range planExclusionIDs {
		if !stateExclusionIDs[id] {
			addExclusions = append(addExclusions, createUpdatableAsset(id, memberType))
		}
	}

	var removeExclusions []windowsupdates.UpdatableAssetable
	for id := range stateExclusionIDs {
		if !planExclusionIDs[id] {
			removeExclusions = append(removeExclusions, createUpdatableAsset(id, memberType))
		}
	}

	if len(addMembers) > 0 {
		requestBody.SetAddMembers(addMembers)
	}
	if len(removeMembers) > 0 {
		requestBody.SetRemoveMembers(removeMembers)
	}
	if len(addExclusions) > 0 {
		requestBody.SetAddExclusions(addExclusions)
	}
	if len(removeExclusions) > 0 {
		requestBody.SetRemoveExclusions(removeExclusions)
	}

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s (update members)", ResourceName), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]any{
			"error": err.Error(),
		})
	}

	return requestBody, nil
}

func constructRemoveAllMembersRequest(ctx context.Context, state *WindowsUpdatesAutopatchDeploymentAudienceResourceModel) (graphadmin.WindowsUpdatesDeploymentAudiencesItemMicrosoftGraphWindowsUpdatesUpdateAudienceUpdateAudiencePostRequestBodyable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing remove all members request for %s resource", ResourceName))

	requestBody := graphadmin.NewWindowsUpdatesDeploymentAudiencesItemMicrosoftGraphWindowsUpdatesUpdateAudienceUpdateAudiencePostRequestBody()
	memberType := state.MemberType.ValueString()

	var removeMembers []windowsupdates.UpdatableAssetable
	if !state.Members.IsNull() && !state.Members.IsUnknown() {
		for _, elem := range state.Members.Elements() {
			if strVal, ok := elem.(types.String); ok {
				removeMembers = append(removeMembers, createUpdatableAsset(strVal.ValueString(), memberType))
			}
		}
	}

	var removeExclusions []windowsupdates.UpdatableAssetable
	if !state.Exclusions.IsNull() && !state.Exclusions.IsUnknown() {
		for _, elem := range state.Exclusions.Elements() {
			if strVal, ok := elem.(types.String); ok {
				removeExclusions = append(removeExclusions, createUpdatableAsset(strVal.ValueString(), memberType))
			}
		}
	}

	if len(removeMembers) > 0 {
		requestBody.SetRemoveMembers(removeMembers)
	}
	if len(removeExclusions) > 0 {
		requestBody.SetRemoveExclusions(removeExclusions)
	}

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s (remove all members)", ResourceName), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]any{
			"error": err.Error(),
		})
	}

	return requestBody, nil
}
