package graphBetaDeviceEnrollmentNotification

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// constructBrandingOptions creates a NotificationMessageTemplate for PATCH request to update branding options
func constructBrandingOptions(ctx context.Context, data *DeviceEnrollmentNotificationResourceModel) (graphmodels.NotificationMessageTemplateable, error) {
	tflog.Debug(ctx, "Constructing branding options PATCH request")

	if data.BrandingOptions.IsNull() || data.BrandingOptions.IsUnknown() {
		return nil, nil
	}

	template := graphmodels.NewNotificationMessageTemplate()

	// Convert the set of strings to a comma-separated string for the bitmask enum parser
	brandingValues := make([]string, 0)
	for _, brandingOption := range data.BrandingOptions.Elements() {
		brandingStr := brandingOption.String()
		// Remove quotes from the string value
		brandingStr = brandingStr[1 : len(brandingStr)-1]
		brandingValues = append(brandingValues, brandingStr)
	}

	// Join the values with commas to create the bitmask format
	brandingOptionsStr := ""
	if len(brandingValues) > 0 {
		brandingOptionsStr = brandingValues[0]
		for i := 1; i < len(brandingValues); i++ {
			brandingOptionsStr += "," + brandingValues[i]
		}
	}

	// Parse the comma-separated string into the proper enum bitmask
	if brandingOptionsStr != "" {
		brandingEnum, err := graphmodels.ParseNotificationTemplateBrandingOptions(brandingOptionsStr)
		if err != nil {
			return nil, fmt.Errorf("failed to parse branding options: %s", err)
		}
		if brandingEnumTyped, ok := brandingEnum.(*graphmodels.NotificationTemplateBrandingOptions); ok {
			template.SetBrandingOptions(brandingEnumTyped)
		} else {
			return nil, fmt.Errorf("failed to cast branding options to correct type")
		}
	}

	tflog.Debug(ctx, fmt.Sprintf("Constructed branding options PATCH with: %s", brandingOptionsStr))

	return template, nil
}
