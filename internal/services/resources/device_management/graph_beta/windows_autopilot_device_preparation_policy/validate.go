package graphBetaWindowsAutopilotDevicePreparationPolicy

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/deviceappmanagement"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/devicemanagement"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models"
	graphcore "github.com/microsoftgraph/msgraph-sdk-go-core"
)

const (
	// The App ID for "Intune Provisioning Client" / "Intune Autopilot ConfidentialClient" service principal
	intuneProvisioningClientAppID = "f1346770-5b25-470b-88bd-d5744ab7952c"

	// allowedAppsFilter is the OData filter for the Windows app types that are valid for
	// Windows Autopilot Device Preparation Policy allowed_apps.
	allowedAppsFilter = "isof('microsoft.graph.windowsAppX')" +
		" or isof('microsoft.graph.windowsMobileMSI')" +
		" or isof('microsoft.graph.windowsUniversalAppX')" +
		" or isof('microsoft.graph.officeSuiteApp')" +
		" or isof('microsoft.graph.windowsMicrosoftEdgeApp')" +
		" or isof('microsoft.graph.winGetApp')" +
		" or isof('microsoft.graph.win32LobApp')" +
		" or isof('microsoft.graph.win32CatalogApp')"
)

// validateSecurityGroupOwnership validates that the specified security group has the Intune Provisioning Client as an owner
func validateSecurityGroupOwnership(ctx context.Context, client *msgraphbetasdk.GraphServiceClient, groupID string) diag.Diagnostics {
	var diags diag.Diagnostics

	tflog.Info(ctx, fmt.Sprintf("Validating security group %s has Intune Provisioning Client as owner", groupID))

	owners, err := client.
		Groups().
		ByGroupId(groupID).
		Owners().
		Get(ctx, nil)

	if err != nil {
		tflog.Error(ctx, "Failed to get security group owners", map[string]any{
			"group_id": groupID,
			"error":    err.Error(),
		})
		diags.AddError(
			"Failed to validate security group ownership",
			fmt.Sprintf(
				"Could not retrieve owners for security group %s: %s",
				groupID,
				err.Error(),
			),
		)
		return diags
	}

	tflog.Debug(ctx, fmt.Sprintf("Retrieved %d owners for security group %s", len(owners.GetValue()), groupID))

	// Check if the Intune Provisioning Client is an owner
	hasIntuneProvisioningClient := false
	for _, owner := range owners.GetValue() {
		servicePrincipal, ok := owner.(models.ServicePrincipalable)
		if ok {
			appID := servicePrincipal.GetAppId()
			if appID != nil && *appID == intuneProvisioningClientAppID {
				tflog.Info(
					ctx,
					"Found Intune Provisioning Client as owner of security group",
					map[string]any{
						"group_id": groupID,
						"app_id":   *appID,
					},
				)
				hasIntuneProvisioningClient = true
				break
			}
		}
	}

	if !hasIntuneProvisioningClient {
		tflog.Error(
			ctx,
			"Security group does not have Intune Provisioning Client as owner",
			map[string]any{
				"group_id":                   groupID,
				"required_service_principal": intuneProvisioningClientAppID,
			},
		)
		diags.AddError(
			"Invalid security group ownership",
			fmt.Sprintf(
				"Security group %s must have the Intune Provisioning Client (AppID: %s) set as its owner. In some tenants, this service principal may appear as 'Intune Autopilot ConfidentialClient'.",
				groupID,
				intuneProvisioningClientAppID,
			),
		)
	}

	return diags
}

// validateRequest validates the allowed_apps and allowed_scripts fields by verifying
// each supplied ID exists in the tenant and is of the expected type.
//
// For allowed_apps it queries:
//
//	GET /beta/deviceAppManagement/mobileApps?$filter=<windows-app-types>&$select=id&$orderby=displayName
//
// and validates that each app_id exists and its @odata.type matches the supplied app_type.
//
// For allowed_scripts it queries:
//
//	GET /beta/deviceManagement/deviceManagementScripts?$select=displayName,id,description
//
// and validates that each script ID exists.
func validateRequest(ctx context.Context, client *msgraphbetasdk.GraphServiceClient, data *WindowsAutopilotDevicePreparationPolicyResourceModel) diag.Diagnostics {
	var diags diag.Diagnostics

	if len(data.AllowedApps) > 0 {
		diags.Append(validateAllowedApps(ctx, client, data.AllowedApps)...)
	}

	if len(data.AllowedScripts) > 0 {
		diags.Append(validateAllowedScripts(ctx, client, data.AllowedScripts)...)
	}

	return diags
}

// validateAllowedApps queries the tenant's Windows mobile apps and verifies that each
// configured app_id exists and its app_type matches the actual type in the tenant.
func validateAllowedApps(
	ctx context.Context,
	client *msgraphbetasdk.GraphServiceClient,
	allowedApps []AllowedAppModel,
) diag.Diagnostics {
	var diags diag.Diagnostics

	tflog.Debug(ctx, "Validating allowed_apps against tenant mobile apps", map[string]any{
		"count": len(allowedApps),
	})

	filter := allowedAppsFilter
	orderby := "displayName"
	requestParams := &deviceappmanagement.MobileAppsRequestBuilderGetRequestConfiguration{
		QueryParameters: &deviceappmanagement.MobileAppsRequestBuilderGetQueryParameters{
			Filter:  &filter,
			Orderby: []string{orderby},
			Select:  []string{"id", "displayName"},
		},
	}

	appsResponse, err := client.DeviceAppManagement().MobileApps().Get(ctx, requestParams)
	if err != nil {
		diags.AddError(
			"Failed to retrieve mobile apps for validation",
			fmt.Sprintf("Could not query tenant mobile apps: %s", err.Error()),
		)
		return diags
	}

	// Build id -> @odata.type map using page iterator to handle pagination.
	appTypeByID := make(map[string]string)

	pageIterator, err := graphcore.NewPageIterator[models.MobileAppable](
		appsResponse,
		client.GetAdapter(),
		models.CreateMobileAppCollectionResponseFromDiscriminatorValue,
	)
	if err != nil {
		diags.AddError(
			"Failed to create page iterator for mobile apps validation",
			err.Error(),
		)
		return diags
	}

	err = pageIterator.Iterate(ctx, func(item models.MobileAppable) bool {
		if item != nil && item.GetId() != nil {
			odataType := ""
			if item.GetOdataType() != nil {
				odataType = *item.GetOdataType()
			}
			appTypeByID[*item.GetId()] = odataType
		}
		return true
	})
	if err != nil {
		diags.AddError(
			"Failed to iterate mobile apps pages during validation",
			err.Error(),
		)
		return diags
	}

	tflog.Debug(ctx, fmt.Sprintf("Retrieved %d Windows mobile apps for validation", len(appTypeByID)))

	// Validate each configured app.
	for _, app := range allowedApps {
		if app.AppID.IsNull() || app.AppID.IsUnknown() {
			continue
		}

		appID := app.AppID.ValueString()
		configuredType := app.AppType.ValueString()
		expectedOdataType := "#microsoft.graph." + configuredType

		actualOdataType, exists := appTypeByID[appID]
		if !exists {
			diags.AddError(
				"Invalid allowed_apps entry: app not found",
				fmt.Sprintf(
					"App ID '%s' (configured type: '%s') was not found in the tenant's Windows app list. "+
						"Ensure the app exists and is one of the supported types: "+
						"windowsAppX, windowsMobileMSI, windowsUniversalAppX, officeSuiteApp, "+
						"windowsMicrosoftEdgeApp, winGetApp, win32LobApp, win32CatalogApp.",
					appID, configuredType,
				),
			)
			continue
		}

		if actualOdataType != expectedOdataType {
			actualShortType := strings.TrimPrefix(actualOdataType, "#microsoft.graph.")
			diags.AddError(
				"Invalid allowed_apps entry: app_type mismatch",
				fmt.Sprintf(
					"App ID '%s' has type '%s' in the tenant but '%s' was specified. "+
						"Please correct the app_type value.",
					appID, actualShortType, configuredType,
				),
			)
		}
	}

	return diags
}

// validateAllowedScripts queries the tenant's device management scripts and verifies
// that each configured script ID exists.
func validateAllowedScripts(ctx context.Context, client *msgraphbetasdk.GraphServiceClient, allowedScripts []types.String) diag.Diagnostics {
	var diags diag.Diagnostics

	tflog.Debug(ctx, "Validating allowed_scripts against tenant device management scripts", map[string]any{
		"count": len(allowedScripts),
	})

	requestParams := &devicemanagement.DeviceManagementScriptsRequestBuilderGetRequestConfiguration{
		QueryParameters: &devicemanagement.DeviceManagementScriptsRequestBuilderGetQueryParameters{
			Select: []string{"id", "displayName", "description"},
		},
	}

	scriptsResponse, err := client.
		DeviceManagement().
		DeviceManagementScripts().
		Get(ctx, requestParams)

	if err != nil {
		diags.AddError(
			"Failed to retrieve device management scripts for validation",
			fmt.Sprintf("Could not query tenant device management scripts: %s", err.Error()),
		)
		return diags
	}

	validScriptIDs := make(map[string]bool)
	for _, script := range scriptsResponse.GetValue() {
		if script.GetId() != nil {
			validScriptIDs[*script.GetId()] = true
		}
	}

	tflog.Debug(ctx, fmt.Sprintf("Retrieved %d device management scripts for validation", len(validScriptIDs)))

	for _, scriptID := range allowedScripts {
		if scriptID.IsNull() || scriptID.IsUnknown() {
			continue
		}

		id := scriptID.ValueString()
		if !validScriptIDs[id] {
			diags.AddError(
				"Invalid allowed_scripts entry: script not found",
				fmt.Sprintf(
					"Device management script with ID '%s' was not found in the tenant. "+
						"Ensure the script exists before referencing it in allowed_scripts.",
					id,
				),
			)
		}
	}

	return diags
}
