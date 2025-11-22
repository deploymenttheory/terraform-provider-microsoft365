package graphBetaWindowsEnrollmentStatusPage

import (
	"context"
	"fmt"
	"regexp"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
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

	for _, appId := range appIdStrings {
		appIdValue := appId.ValueString()

		// Query the specific app by ID
		app, err := client.
			DeviceAppManagement().
			MobileApps().
			ByMobileAppId(appIdValue).
			Get(ctx, nil)

		if err != nil {
			tflog.Error(ctx, "Failed to retrieve mobile app for validation", map[string]any{
				"appId": appIdValue,
				"error": err.Error(),
			})
			return fmt.Errorf("supplied app ID '%s' does not match any valid Windows app types. Valid app types include: windowsAppX, windowsMobileMSI, windowsUniversalAppX, officeSuiteApp, windowsMicrosoftEdgeApp, winGetApp, win32LobApp, win32CatalogApp", appIdValue)
		}

		// Validate the app type using SDK type assertions
		isValidType := false
		var appTypeName string

		switch app.(type) {
		case *graphmodels.WindowsAppX:
			isValidType = true
			appTypeName = "windowsAppX"
		case *graphmodels.WindowsMobileMSI:
			isValidType = true
			appTypeName = "windowsMobileMSI"
		case *graphmodels.WindowsUniversalAppX:
			isValidType = true
			appTypeName = "windowsUniversalAppX"
		case *graphmodels.OfficeSuiteApp:
			isValidType = true
			appTypeName = "officeSuiteApp"
		case *graphmodels.WindowsMicrosoftEdgeApp:
			isValidType = true
			appTypeName = "windowsMicrosoftEdgeApp"
		case *graphmodels.WinGetApp:
			isValidType = true
			appTypeName = "winGetApp"
		case *graphmodels.Win32LobApp:
			isValidType = true
			appTypeName = "win32LobApp"
		case *graphmodels.Win32CatalogApp:
			isValidType = true
			appTypeName = "win32CatalogApp"
		default:
			if odataType := app.GetOdataType(); odataType != nil {
				appTypeName = *odataType
			} else {
				appTypeName = "unknown"
			}
		}

		if !isValidType {
			return fmt.Errorf("supplied app ID '%s' has type '%s' which is not a valid Windows app type. Valid app types include: windowsAppX, windowsMobileMSI, windowsUniversalAppX, officeSuiteApp, windowsMicrosoftEdgeApp, winGetApp, win32LobApp, win32CatalogApp", appIdValue, appTypeName)
		}

		displayName := ""
		if app.GetDisplayName() != nil {
			displayName = *app.GetDisplayName()
		}

		tflog.Debug(ctx, "Validated mobile app", map[string]any{
			"appId":       appIdValue,
			"displayName": displayName,
			"appType":     appTypeName,
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Successfully validated %d mobile app IDs", len(appIdStrings)))
	return nil
}
