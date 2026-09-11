package graphBetaWindowsEnrollmentStatusPage

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

var validWindowsAppTypeNames = []string{
	"windowsAppX",
	"windowsMobileMSI",
	"windowsUniversalAppX",
	"officeSuiteApp",
	"windowsMicrosoftEdgeApp",
	"winGetApp",
	"win32LobApp",
	"win32CatalogApp",
}

var validWindowsAppTypes = func() map[string]bool {
	m := make(map[string]bool, len(validWindowsAppTypeNames))
	for _, t := range validWindowsAppTypeNames {
		m[t] = true
	}
	return m
}()

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

	// Each app is fetched by ID rather than matched against a listing, as any listing is
	// page limited and would spuriously reject apps in tenants with large app catalogs.
	for _, appId := range appIdStrings {
		appIdValue := appId.ValueString()

		app, err := client.
			DeviceAppManagement().
			MobileApps().
			ByMobileAppId(appIdValue).
			Get(ctx, nil)

		if err != nil {
			if graphErr := errors.GraphError(ctx, err); graphErr.StatusCode == http.StatusNotFound {
				return fmt.Errorf("supplied app ID '%s' was not found in Intune", appIdValue)
			}

			tflog.Error(ctx, "Failed to retrieve mobile app for validation", map[string]any{
				"appId": appIdValue,
				"error": err.Error(),
			})
			return fmt.Errorf("failed to validate mobile app ID '%s': unable to retrieve app from Microsoft Graph", appIdValue)
		}

		if app == nil || app.GetOdataType() == nil {
			return fmt.Errorf("supplied app ID '%s' returned no app type from Microsoft Graph", appIdValue)
		}

		appType := strings.TrimPrefix(*app.GetOdataType(), "#microsoft.graph.")
		if !validWindowsAppTypes[appType] {
			return fmt.Errorf("supplied app ID '%s' is of type '%s' which is not a valid Windows app type. Valid app types include: %s", appIdValue, appType, strings.Join(validWindowsAppTypeNames, ", "))
		}

		displayName := ""
		if app.GetDisplayName() != nil {
			displayName = *app.GetDisplayName()
		}

		tflog.Debug(ctx, "Validated mobile app", map[string]any{
			"appId":       appIdValue,
			"displayName": displayName,
			"appType":     appType,
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Validated %d mobile app IDs", len(appIdStrings)))
	return nil
}
