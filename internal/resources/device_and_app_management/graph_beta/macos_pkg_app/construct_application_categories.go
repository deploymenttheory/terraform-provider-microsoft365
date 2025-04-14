package graphBetaMacOSPKGApp

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	abstractions "github.com/microsoft/kiota-abstractions-go"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

// assignCategoriesToMobileApplication creates an associations between a mobile app and categories.
// It supports both category IDs (UUIDs) and category names as inputs, looking up the
// appropriate IDs as needed.
func assignCategoriesToMobileApplication(
	ctx context.Context,
	client *msgraphbetasdk.GraphServiceClient,
	appID string,
	categoryValues []string,
	readPermissions []string) error {

	tflog.Debug(ctx, fmt.Sprintf("Associating app %s with %d categories", appID, len(categoryValues)))

	if len(categoryValues) == 0 {
		return nil
	}

	// Handle Built-in categories mapping
	categoryMapping := map[string]string{
		"Other apps":             "0720a99e-562b-4a77-83f0-9a7523fcf13e",
		"Books & Reference":      "f1fc9fe2-728d-4867-9a72-a61e18f8c606",
		"Data management":        "046e0b16-76ce-4b49-bf1b-1cc5bd94fb47",
		"Productivity":           "ed899483-3019-425e-a470-28e901b9790e",
		"Business":               "2b73ae71-12c8-49be-b462-3dae769ccd9d",
		"Development & Design":   "79bc98d4-7ddf-4841-9bc1-5c84a26d7ee8",
		"Photos & Media":         "5dcd7a90-0306-4f09-a75d-6b97a243f04e",
		"Collaboration & Social": "f79135dc-8e41-48c1-9a59-ab9a7259c38e",
		"Computer management":    "981deed8-6857-4e78-a50e-c3f61d312737",
	}

	allCategories, err := client.
		DeviceAppManagement().
		MobileAppCategories().
		Get(ctx, nil)

	if err != nil {
		return fmt.Errorf("failed to list application categories: %w", err)
	}

	dynamicCategoryMap := make(map[string]string)
	for _, cat := range allCategories.GetValue() {
		if cat.GetDisplayName() != nil && cat.GetId() != nil {
			dynamicCategoryMap[*cat.GetDisplayName()] = *cat.GetId()
		}
	}

	for _, categoryValue := range categoryValues {
		var categoryID string

		if _, err := uuid.Parse(categoryValue); err == nil {
			categoryID = categoryValue
			tflog.Debug(ctx, fmt.Sprintf("Using category ID: %s", categoryID))
		} else {
			if id, exists := categoryMapping[categoryValue]; exists {
				categoryID = id
				tflog.Debug(ctx, fmt.Sprintf("Found built-in category '%s' with ID '%s'", categoryValue, categoryID))
			} else if id, exists := dynamicCategoryMap[categoryValue]; exists {
				categoryID = id
				tflog.Debug(ctx, fmt.Sprintf("Found user-defined category '%s' with ID '%s'", categoryValue, categoryID))
			} else {
				return fmt.Errorf("category with name '%s' not found", categoryValue)
			}
		}

		referenceURL := fmt.Sprintf("https://graph.microsoft.com/beta/deviceAppManagement/mobileAppCategories/%s", categoryID)
		requestInfo := abstractions.NewRequestInformation()
		requestInfo.Method = abstractions.POST
		endpoint := fmt.Sprintf("deviceAppManagement/mobileApps/%s/categories/$ref", appID)
		requestInfo.UrlTemplate = "{+baseurl}/" + endpoint
		requestInfo.PathParameters = map[string]string{
			"baseurl": "https://graph.microsoft.com/beta",
		}

		jsonContent := fmt.Sprintf(`{"@odata.id": "%s"}`, referenceURL)
		requestInfo.SetStreamContentAndContentType([]byte(jsonContent), "application/json")

		err = client.GetAdapter().SendNoContent(ctx, requestInfo, nil)
		if err != nil {
			return fmt.Errorf("failed to associate app with category: %w", err)
		}

		tflog.Debug(ctx, fmt.Sprintf("Successfully associated category %s with app", categoryID))
	}

	return nil
}
