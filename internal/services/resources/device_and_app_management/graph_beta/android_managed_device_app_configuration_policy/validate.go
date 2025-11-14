package graphBetaAndroidManagedDeviceAppConfigurationPolicy

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/deviceappmanagement"
)

// -----------------------------------------------------------------------------
// Validators performed at the attribute level of the Terraform configuration
// -----------------------------------------------------------------------------

// validateIOSMobileAppIds validates that all provided app IDs exist in Intune as valid iOS apps
func validateAndroidMobileAppIds(ctx context.Context, client *msgraphbetasdk.GraphServiceClient, appIds []string) error {
	if len(appIds) == 0 {
		return nil
	}

	tflog.Debug(ctx, fmt.Sprintf("Validating %d Android mobile app IDs", len(appIds)))

	filter := "isof('microsoft.graph.androidManagedStoreApp') and microsoft.graph.androidManagedStoreApp/packageId ne 'com.microsoft.windowsintune.companyportal' and microsoft.graph.androidManagedStoreApp/packageId ne 'com.microsoft.intune' and microsoft.graph.androidManagedStoreApp/supportsOemConfig eq false"
	orderby := []string{"displayName"}

	requestConfig := &deviceappmanagement.MobileAppsRequestBuilderGetRequestConfiguration{
		QueryParameters: &deviceappmanagement.MobileAppsRequestBuilderGetQueryParameters{
			Filter:  &filter,
			Orderby: orderby,
		},
	}

	// Query the API for available Android apps
	response, err := client.
		DeviceAppManagement().
		MobileApps().
		Get(ctx, requestConfig)

	if err != nil {
		return fmt.Errorf("failed to query Android mobile apps from Intune: %w", err)
	}

	if response == nil {
		return fmt.Errorf("received nil response when querying Android mobile apps")
	}

	validAppIds := make(map[string]bool)
	if apps := response.GetValue(); apps != nil {
		for _, app := range apps {
			if appId := app.GetId(); appId != nil {
				validAppIds[strings.ToLower(*appId)] = true
			}
		}
	}

	tflog.Debug(ctx, fmt.Sprintf("Found %d valid Android apps in Intune", len(validAppIds)))

	// Validate each provided app ID
	var invalidIds []string
	for _, appId := range appIds {
		if !validAppIds[strings.ToLower(appId)] {
			invalidIds = append(invalidIds, appId)
		}
	}

	if len(invalidIds) > 0 {
		return fmt.Errorf("the following app IDs are not valid Android Managed Store apps in your Intune tenant: %s. Valid apps must be androidManagedStoreApp type with supportsOemConfig=false, excluding Company Portal and Intune apps", strings.Join(invalidIds, ", "))
	}

	tflog.Debug(ctx, "All provided Android mobile app IDs are valid")
	return nil
}

// -----------------------------------------------------------------------------
// Validators performed at the resource level of the Terraform configuration
// -----------------------------------------------------------------------------

// ConfigValidators implements the resource.ResourceWithConfigValidators interface.
// This method returns resource-level validators that perform cross-field validation across
// the entire resource configuration, as opposed to attribute-level validators that validate
// individual fields in isolation.
func (r *AndroidManagedDeviceAppConfigurationPolicyResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	// Android Managed Store App Configuration doesn't have mutually exclusive fields like iOS
	return []resource.ConfigValidator{}
}
