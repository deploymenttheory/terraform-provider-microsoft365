package graphBetaWindowsUpdatesAutopatchDeviceRegistration

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphadmin "github.com/microsoftgraph/msgraph-beta-sdk-go/admin"
	graphmodelswindowsupdates "github.com/microsoftgraph/msgraph-beta-sdk-go/models/windowsupdates"
)

func constructEnrollRequest(ctx context.Context, data *WindowsUpdatesAutopatchDeviceRegistrationResourceModel) (graphadmin.WindowsUpdatesUpdatableAssetsMicrosoftGraphWindowsUpdatesEnrollAssetsByIdEnrollAssetsByIdPostRequestBodyable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing enroll request for %s resource", ResourceName))

	requestBody := graphadmin.NewWindowsUpdatesUpdatableAssetsMicrosoftGraphWindowsUpdatesEnrollAssetsByIdEnrollAssetsByIdPostRequestBody()

	updateCategoryRaw, err := graphmodelswindowsupdates.ParseUpdateCategory(data.UpdateCategory.ValueString())
	if err != nil {
		return nil, fmt.Errorf("failed to parse update category: %w", err)
	}
	updateCategory := updateCategoryRaw.(*graphmodelswindowsupdates.UpdateCategory)
	requestBody.SetUpdateCategory(updateCategory)

	memberEntityType := "#microsoft.graph.windowsUpdates.azureADDevice"
	requestBody.SetMemberEntityType(&memberEntityType)

	var deviceIDs []string
	if !data.DeviceIds.IsNull() && !data.DeviceIds.IsUnknown() {
		elements := data.DeviceIds.Elements()
		deviceIDs = make([]string, 0, len(elements))
		for _, elem := range elements {
			if strVal, ok := elem.(types.String); ok {
				deviceIDs = append(deviceIDs, strVal.ValueString())
			}
		}
	}

	if len(deviceIDs) > 0 {
		requestBody.SetIds(deviceIDs)
	}

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s (Enroll)", ResourceName), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]any{
			"error": err.Error(),
		})
	}

	return requestBody, nil
}

func constructUnenrollRequest(ctx context.Context, data *WindowsUpdatesAutopatchDeviceRegistrationResourceModel) (graphadmin.WindowsUpdatesUpdatableAssetsMicrosoftGraphWindowsUpdatesUnenrollAssetsByIdUnenrollAssetsByIdPostRequestBodyable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing unenroll request for %s resource", ResourceName))

	requestBody := graphadmin.NewWindowsUpdatesUpdatableAssetsMicrosoftGraphWindowsUpdatesUnenrollAssetsByIdUnenrollAssetsByIdPostRequestBody()

	updateCategoryRaw, err := graphmodelswindowsupdates.ParseUpdateCategory(data.UpdateCategory.ValueString())
	if err != nil {
		return nil, fmt.Errorf("failed to parse update category: %w", err)
	}
	updateCategory := updateCategoryRaw.(*graphmodelswindowsupdates.UpdateCategory)
	requestBody.SetUpdateCategory(updateCategory)

	memberEntityType := "#microsoft.graph.windowsUpdates.azureADDevice"
	requestBody.SetMemberEntityType(&memberEntityType)

	var deviceIDs []string
	if !data.DeviceIds.IsNull() && !data.DeviceIds.IsUnknown() {
		elements := data.DeviceIds.Elements()
		deviceIDs = make([]string, 0, len(elements))
		for _, elem := range elements {
			if strVal, ok := elem.(types.String); ok {
				deviceIDs = append(deviceIDs, strVal.ValueString())
			}
		}
	}

	if len(deviceIDs) > 0 {
		requestBody.SetIds(deviceIDs)
	}

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s (Unenroll)", ResourceName), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]any{
			"error": err.Error(),
		})
	}

	return requestBody, nil
}
