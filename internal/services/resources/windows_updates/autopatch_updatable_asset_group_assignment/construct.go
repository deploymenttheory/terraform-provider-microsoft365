package graphBetaWindowsUpdatesAutopatchUpdatableAssetGroupAssignment

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphadmin "github.com/microsoftgraph/msgraph-beta-sdk-go/admin"
)

func constructAddMembersRequest(ctx context.Context, deviceIDs []string) (graphadmin.WindowsUpdatesUpdatableAssetsItemMicrosoftGraphWindowsUpdatesAddMembersByIdAddMembersByIdPostRequestBodyable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing addMembersById request for %s resource", ResourceName))

	requestBody := graphadmin.NewWindowsUpdatesUpdatableAssetsItemMicrosoftGraphWindowsUpdatesAddMembersByIdAddMembersByIdPostRequestBody()
	memberEntityType := "#microsoft.graph.windowsUpdates.azureADDevice"
	requestBody.SetMemberEntityType(&memberEntityType)
	requestBody.SetIds(deviceIDs)

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s (AddMembersById)", ResourceName), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]any{
			"error": err.Error(),
		})
	}

	return requestBody, nil
}

func constructRemoveMembersRequest(ctx context.Context, deviceIDs []string) (graphadmin.WindowsUpdatesUpdatableAssetsItemMicrosoftGraphWindowsUpdatesRemoveMembersByIdRemoveMembersByIdPostRequestBodyable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing removeMembersById request for %s resource", ResourceName))

	requestBody := graphadmin.NewWindowsUpdatesUpdatableAssetsItemMicrosoftGraphWindowsUpdatesRemoveMembersByIdRemoveMembersByIdPostRequestBody()
	memberEntityType := "#microsoft.graph.windowsUpdates.azureADDevice"
	requestBody.SetMemberEntityType(&memberEntityType)
	requestBody.SetIds(deviceIDs)

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s (RemoveMembersById)", ResourceName), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]any{
			"error": err.Error(),
		})
	}

	return requestBody, nil
}

func extractDeviceIDs(model *WindowsUpdatesAutopatchUpdatableAssetGroupAssignmentResourceModel) []string {
	if model.EntraDeviceIds.IsNull() || model.EntraDeviceIds.IsUnknown() {
		return nil
	}
	elements := model.EntraDeviceIds.Elements()
	ids := make([]string, 0, len(elements))
	for _, elem := range elements {
		if strVal, ok := elem.(types.String); ok {
			ids = append(ids, strVal.ValueString())
		}
	}
	return ids
}
