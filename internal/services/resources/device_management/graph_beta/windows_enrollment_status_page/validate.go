package graphBetaWindowsEnrollmentStatusPage

import (
	"context"
	"fmt"
	"regexp"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/deviceappmanagement"
)

// validateRequest validates the entire request payload
func validateRequest(ctx context.Context, client *msgraphbetasdk.GraphServiceClient, data *WindowsEnrollmentStatusPageResourceModel) error {
	if err := validateSelectedMobileAppIds(ctx, client, data.SelectedMobileAppIds); err != nil {
		return err
	}

	return nil
}

// validateSelectedMobileAppIds validates that the provided mobile app IDs exist and are valid Windows app types
func validateSelectedMobileAppIds(ctx context.Context, client *msgraphbetasdk.GraphServiceClient, appIds types.Set) error {
	if appIds.IsNull() || appIds.IsUnknown() {
		return nil
	}

	var appIdStrings []types.String
	appIds.ElementsAs(ctx, &appIdStrings, false)

	if len(appIdStrings) == 0 {
		return nil
	}

	// First validate UUID format
	guidRegex := regexp.MustCompile(constants.GuidRegex)
	for _, appId := range appIdStrings {
		if !guidRegex.MatchString(appId.ValueString()) {
			return fmt.Errorf("invalid application ID format: %s. Must be a valid UUID", appId.ValueString())
		}
	}

	// Skip API validation if client is nil (for unit tests)
	if client == nil {
		tflog.Debug(ctx, "Skipping API validation (client is nil)")
		return nil
	}

	// Get all Windows app types from Microsoft Graph
	filter := "isof('microsoft.graph.windowsAppX') or isof('microsoft.graph.windowsMobileMSI') or isof('microsoft.graph.windowsUniversalAppX') or isof('microsoft.graph.officeSuiteApp') or isof('microsoft.graph.windowsMicrosoftEdgeApp') or isof('microsoft.graph.winGetApp') or isof('microsoft.graph.win32LobApp') or isof('microsoft.graph.win32CatalogApp')"
	orderby := "displayname"
	top := int32(250)

	requestConfig := &deviceappmanagement.MobileAppsRequestBuilderGetRequestConfiguration{
		QueryParameters: &deviceappmanagement.MobileAppsRequestBuilderGetQueryParameters{
			Filter:  &filter,
			Orderby: []string{orderby},
			Top:     &top,
		},
	}

	mobileApps, err := client.
		DeviceAppManagement().
		MobileApps().
		Get(ctx, requestConfig)

	if err != nil {
		tflog.Error(ctx, "Failed to retrieve mobile apps for validation", map[string]any{
			"error": err.Error(),
		})
		return fmt.Errorf("failed to validate mobile app IDs: unable to retrieve available apps from Microsoft Graph")
	}

	// Create a map of valid app IDs for quick lookup
	validAppIds := make(map[string]string)   // ID -> DisplayName
	validAppTypes := make(map[string]string) // ID -> AppType

	if mobileApps.GetValue() != nil {
		for _, app := range mobileApps.GetValue() {
			if app.GetId() != nil && app.GetDisplayName() != nil && app.GetOdataType() != nil {
				validAppIds[*app.GetId()] = *app.GetDisplayName()
				validAppTypes[*app.GetId()] = *app.GetOdataType()
			}
		}
	}

	// Validate each provided app ID
	for _, appId := range appIdStrings {
		appIdValue := appId.ValueString()
		displayName, exists := validAppIds[appIdValue]

		if !exists {
			return fmt.Errorf("supplied app ID '%s' does not match any valid Windows app types. Valid app types include: windowsAppX, windowsMobileMSI, windowsUniversalAppX, officeSuiteApp, windowsMicrosoftEdgeApp, winGetApp, win32LobApp, win32CatalogApp", appIdValue)
		}

		tflog.Debug(ctx, "Validated mobile app", map[string]any{
			"appId":       appIdValue,
			"displayName": displayName,
			"appType":     validAppTypes[appIdValue],
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Validated %d mobile app IDs against %d available Windows apps", len(appIdStrings), len(validAppIds)))
	return nil
}
