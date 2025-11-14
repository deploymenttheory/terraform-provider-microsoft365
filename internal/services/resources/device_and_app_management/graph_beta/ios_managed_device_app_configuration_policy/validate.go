package graphBetaIOSManagedDeviceAppConfigurationPolicy

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/deviceappmanagement"

	resourceLevel "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/validate/resource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// -----------------------------------------------------------------------------
// Validators performed at the attribute level of the Terraform configuration
// -----------------------------------------------------------------------------

// validateIOSMobileAppIds validates that all provided app IDs exist in Intune as valid iOS apps
func validateIOSMobileAppIds(ctx context.Context, client *msgraphbetasdk.GraphServiceClient, appIds []string) error {
	if len(appIds) == 0 {
		return nil
	}

	tflog.Debug(ctx, fmt.Sprintf("Validating %d iOS mobile app IDs", len(appIds)))

	filter := "(microsoft.graph.managedApp/appAvailability eq null or microsoft.graph.managedApp/appAvailability eq 'lineOfBusiness' or isAssigned eq true) and (isof('microsoft.graph.iosLobApp') or isof('microsoft.graph.iosStoreApp') or isof('microsoft.graph.iosVppApp') or isof('microsoft.graph.managedIOSStoreApp') or isof('microsoft.graph.managedIOSLobApp'))"
	orderby := []string{"displayName"}

	requestConfig := &deviceappmanagement.MobileAppsRequestBuilderGetRequestConfiguration{
		QueryParameters: &deviceappmanagement.MobileAppsRequestBuilderGetQueryParameters{
			Filter:  &filter,
			Orderby: orderby,
		},
	}

	// Query the API for available iOS apps
	response, err := client.
		DeviceAppManagement().
		MobileApps().
		Get(ctx, requestConfig)

	if err != nil {
		return fmt.Errorf("failed to query iOS mobile apps from Intune: %w", err)
	}

	if response == nil {
		return fmt.Errorf("received nil response when querying iOS mobile apps")
	}

	validAppIds := make(map[string]bool)
	if apps := response.GetValue(); apps != nil {
		for _, app := range apps {
			if appId := app.GetId(); appId != nil {
				validAppIds[strings.ToLower(*appId)] = true
			}
		}
	}

	tflog.Debug(ctx, fmt.Sprintf("Found %d valid iOS apps in Intune", len(validAppIds)))

	// Validate each provided app ID
	var invalidIds []string
	for _, appId := range appIds {
		if !validAppIds[strings.ToLower(appId)] {
			invalidIds = append(invalidIds, appId)
		}
	}

	if len(invalidIds) > 0 {
		return fmt.Errorf("the following app IDs are not valid iOS apps in your Intune tenant: %s. Valid iOS app types include: iosLobApp, iosStoreApp, iosVppApp, managedIOSStoreApp, managedIOSLobApp", strings.Join(invalidIds, ", "))
	}

	tflog.Debug(ctx, "All provided iOS mobile app IDs are valid")
	return nil
}

// -----------------------------------------------------------------------------
// Validators performed at the resource level of the Terraform configuration
// -----------------------------------------------------------------------------

// ConfigValidators implements the resource.ResourceWithConfigValidators interface.
// This method returns resource-level validators that perform cross-field validation across
// the entire resource configuration, as opposed to attribute-level validators that validate
// individual fields in isolation.
//
// For iOS Mobile App Configuration, this validates that encoded_setting_xml and settings
// are mutually exclusive - only one configuration method can be used at a time.
func (r *IOSManagedDeviceAppConfigurationPolicyResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{
		resourceLevel.MutuallyExclusiveAttributes(
			[]path.Path{
				path.Root("encoded_setting_xml"),
				path.Root("settings"),
			},
			[]string{
				"encoded_setting_xml",
				"settings",
			},
		),
	}
}
