package graphBetaWindowsAutopatchDeploymentAudienceMembers

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphadmin "github.com/microsoftgraph/msgraph-beta-sdk-go/admin"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models/windowsupdates"
)

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

func constructCreateRequest(ctx context.Context, plan *WindowsUpdateDeploymentAudienceMembersResourceModel) (graphadmin.WindowsUpdatesDeploymentAudiencesItemMicrosoftGraphWindowsUpdatesUpdateAudienceUpdateAudiencePostRequestBodyable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing create request for %s resource", ResourceName))

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

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s (Create)", ResourceName), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]any{
			"error": err.Error(),
		})
	}

	return requestBody, nil
}

func constructUpdateRequest(ctx context.Context, plan *WindowsUpdateDeploymentAudienceMembersResourceModel, state *WindowsUpdateDeploymentAudienceMembersResourceModel) (graphadmin.WindowsUpdatesDeploymentAudiencesItemMicrosoftGraphWindowsUpdatesUpdateAudienceUpdateAudiencePostRequestBodyable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing update request for %s resource", ResourceName))

	requestBody := graphadmin.NewWindowsUpdatesDeploymentAudiencesItemMicrosoftGraphWindowsUpdatesUpdateAudienceUpdateAudiencePostRequestBody()
	memberType := plan.MemberType.ValueString()

	planMemberIDs := make(map[string]bool)
	if !plan.Members.IsNull() && !plan.Members.IsUnknown() {
		elements := plan.Members.Elements()
		for _, elem := range elements {
			if strVal, ok := elem.(types.String); ok {
				planMemberIDs[strVal.ValueString()] = true
			}
		}
	}

	stateMemberIDs := make(map[string]bool)
	if !state.Members.IsNull() && !state.Members.IsUnknown() {
		elements := state.Members.Elements()
		for _, elem := range elements {
			if strVal, ok := elem.(types.String); ok {
				stateMemberIDs[strVal.ValueString()] = true
			}
		}
	}

	var addMembers []windowsupdates.UpdatableAssetable
	for id := range planMemberIDs {
		if !stateMemberIDs[id] {
			asset := createUpdatableAsset(id, memberType)
			addMembers = append(addMembers, asset)
		}
	}

	var removeMembers []windowsupdates.UpdatableAssetable
	for id := range stateMemberIDs {
		if !planMemberIDs[id] {
			asset := createUpdatableAsset(id, memberType)
			removeMembers = append(removeMembers, asset)
		}
	}

	planExclusionIDs := make(map[string]bool)
	if !plan.Exclusions.IsNull() && !plan.Exclusions.IsUnknown() {
		elements := plan.Exclusions.Elements()
		for _, elem := range elements {
			if strVal, ok := elem.(types.String); ok {
				planExclusionIDs[strVal.ValueString()] = true
			}
		}
	}

	stateExclusionIDs := make(map[string]bool)
	if !state.Exclusions.IsNull() && !state.Exclusions.IsUnknown() {
		elements := state.Exclusions.Elements()
		for _, elem := range elements {
			if strVal, ok := elem.(types.String); ok {
				stateExclusionIDs[strVal.ValueString()] = true
			}
		}
	}

	var addExclusions []windowsupdates.UpdatableAssetable
	for id := range planExclusionIDs {
		if !stateExclusionIDs[id] {
			asset := createUpdatableAsset(id, memberType)
			addExclusions = append(addExclusions, asset)
		}
	}

	var removeExclusions []windowsupdates.UpdatableAssetable
	for id := range stateExclusionIDs {
		if !planExclusionIDs[id] {
			asset := createUpdatableAsset(id, memberType)
			removeExclusions = append(removeExclusions, asset)
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

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s (Update)", ResourceName), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]any{
			"error": err.Error(),
		})
	}

	return requestBody, nil
}

func constructDeleteRequest(ctx context.Context, state *WindowsUpdateDeploymentAudienceMembersResourceModel) (graphadmin.WindowsUpdatesDeploymentAudiencesItemMicrosoftGraphWindowsUpdatesUpdateAudienceUpdateAudiencePostRequestBodyable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing delete request for %s resource", ResourceName))

	requestBody := graphadmin.NewWindowsUpdatesDeploymentAudiencesItemMicrosoftGraphWindowsUpdatesUpdateAudienceUpdateAudiencePostRequestBody()
	memberType := state.MemberType.ValueString()

	var removeMembers []windowsupdates.UpdatableAssetable
	if !state.Members.IsNull() && !state.Members.IsUnknown() {
		elements := state.Members.Elements()
		removeMembers = make([]windowsupdates.UpdatableAssetable, 0, len(elements))
		for _, elem := range elements {
			if strVal, ok := elem.(types.String); ok {
				asset := createUpdatableAsset(strVal.ValueString(), memberType)
				removeMembers = append(removeMembers, asset)
			}
		}
	}

	var removeExclusions []windowsupdates.UpdatableAssetable
	if !state.Exclusions.IsNull() && !state.Exclusions.IsUnknown() {
		elements := state.Exclusions.Elements()
		removeExclusions = make([]windowsupdates.UpdatableAssetable, 0, len(elements))
		for _, elem := range elements {
			if strVal, ok := elem.(types.String); ok {
				asset := createUpdatableAsset(strVal.ValueString(), memberType)
				removeExclusions = append(removeExclusions, asset)
			}
		}
	}

	if len(removeMembers) > 0 {
		requestBody.SetRemoveMembers(removeMembers)
	}
	if len(removeExclusions) > 0 {
		requestBody.SetRemoveExclusions(removeExclusions)
	}

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s (Delete)", ResourceName), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]any{
			"error": err.Error(),
		})
	}

	return requestBody, nil
}
