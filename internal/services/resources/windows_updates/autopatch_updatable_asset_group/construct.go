package graphBetaWindowsUpdatesAutopatchUpdatableAssetGroup

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphadmin "github.com/microsoftgraph/msgraph-beta-sdk-go/admin"
	graphmodelswindowsupdates "github.com/microsoftgraph/msgraph-beta-sdk-go/models/windowsupdates"
)

func constructResource(ctx context.Context, _ *WindowsUpdatesAutopatchUpdatableAssetGroupResourceModel) (graphmodelswindowsupdates.UpdatableAssetGroupable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s resource from model", ResourceName))

	requestBody := graphmodelswindowsupdates.NewUpdatableAssetGroup()

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]any{
			"error": err.Error(),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing %s resource", ResourceName))
	return requestBody, nil
}

func constructAddMembersRequest(ctx context.Context, deviceObjectIDs []string) (graphadmin.WindowsUpdatesUpdatableAssetsItemMicrosoftGraphWindowsUpdatesAddMembersByIdAddMembersByIdPostRequestBodyable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing addMembersById request for %s resource", ResourceName))

	requestBody := graphadmin.NewWindowsUpdatesUpdatableAssetsItemMicrosoftGraphWindowsUpdatesAddMembersByIdAddMembersByIdPostRequestBody()
	memberEntityType := "#microsoft.graph.windowsUpdates.azureADDevice"
	requestBody.SetMemberEntityType(&memberEntityType)
	requestBody.SetIds(deviceObjectIDs)

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s (addMembersById)", ResourceName), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]any{
			"error": err.Error(),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing addMembersById request for %s resource", ResourceName))
	return requestBody, nil
}

func constructRemoveMembersRequest(ctx context.Context, deviceObjectIDs []string) (graphadmin.WindowsUpdatesUpdatableAssetsItemMicrosoftGraphWindowsUpdatesRemoveMembersByIdRemoveMembersByIdPostRequestBodyable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing removeMembersById request for %s resource", ResourceName))

	requestBody := graphadmin.NewWindowsUpdatesUpdatableAssetsItemMicrosoftGraphWindowsUpdatesRemoveMembersByIdRemoveMembersByIdPostRequestBody()
	memberEntityType := "#microsoft.graph.windowsUpdates.azureADDevice"
	requestBody.SetMemberEntityType(&memberEntityType)
	requestBody.SetIds(deviceObjectIDs)

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s (removeMembersById)", ResourceName), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]any{
			"error": err.Error(),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing removeMembersById request for %s resource", ResourceName))
	return requestBody, nil
}

func extractDeviceObjectIDs(model *WindowsUpdatesAutopatchUpdatableAssetGroupResourceModel) []string {
	if model.EntraDeviceObjectIds.IsNull() || model.EntraDeviceObjectIds.IsUnknown() {
		return nil
	}
	elements := model.EntraDeviceObjectIds.Elements()
	ids := make([]string, 0, len(elements))
	for _, elem := range elements {
		if strVal, ok := elem.(types.String); ok {
			ids = append(ids, strVal.ValueString())
		}
	}
	return ids
}
