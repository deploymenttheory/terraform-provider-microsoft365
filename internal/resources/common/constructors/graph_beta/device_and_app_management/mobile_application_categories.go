package sharedConstructors

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	abstractions "github.com/microsoft/kiota-abstractions-go"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

var builtInCategoryMapping = map[string]string{
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

// assignMobileAppCategories performs a full update of mobile app categories.
func AssignMobileAppCategories(
	ctx context.Context,
	client *msgraphbetasdk.GraphServiceClient,
	appID string,
	categoryValues []string,
	readPermissions []string,
) error {
	tflog.Debug(ctx, fmt.Sprintf("Associating app %s with %d categories", appID, len(categoryValues)))

	if len(categoryValues) == 0 {
		return nil
	}

	// Fetch existing categories (for user-defined category name lookup)
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

	// Resolve category names/IDs to IDs
	var resolvedIDs []string
	for _, categoryValue := range categoryValues {
		var categoryID string
		if _, err := uuid.Parse(categoryValue); err == nil {
			categoryID = categoryValue
		} else if id, ok := builtInCategoryMapping[categoryValue]; ok {
			categoryID = id
		} else if id, ok := dynamicCategoryMap[categoryValue]; ok {
			categoryID = id
		} else {
			return fmt.Errorf("category with name '%s' not found", categoryValue)
		}
		resolvedIDs = append(resolvedIDs, categoryID)
	}

	// Remove all current mobile app categories
	for _, categoryID := range dynamicCategoryMap {
		err := mobileAppCategoryDeleteRequest(ctx, client, appID, categoryID)
		if err != nil {
			tflog.Debug(ctx, fmt.Sprintf("Skipping delete for category %s: %v", categoryID, err))
		}
	}

	// Assign
	for _, categoryID := range resolvedIDs {
		reqInfo := buildMobileAppCategoryRequest(appID, categoryID)
		err := mobileAppCategoryPostRequest(ctx, client, reqInfo)
		if err != nil {
			return fmt.Errorf("failed to associate mobile app with category ID '%s': %w", categoryID, err)
		}
		tflog.Debug(ctx, fmt.Sprintf("Associated mobile app with category ID: %s", categoryID))
	}

	return nil
}

// buildMobileAppCategoryRequest constructs the POST request to associate a category.
func buildMobileAppCategoryRequest(appID, categoryID string) *abstractions.RequestInformation {
	referenceURL := fmt.Sprintf("https://graph.microsoft.com/beta/deviceAppManagement/mobileAppCategories/%s", categoryID)

	requestInfo := abstractions.NewRequestInformation()
	requestInfo.Method = abstractions.POST
	requestInfo.UrlTemplate = "{+baseurl}/deviceAppManagement/mobileApps/" + appID + "/categories/$ref"
	requestInfo.PathParameters = map[string]string{
		"baseurl": "https://graph.microsoft.com/beta",
	}

	jsonContent := fmt.Sprintf(`{"@odata.id": "%s"}`, referenceURL)
	requestInfo.SetStreamContentAndContentType([]byte(jsonContent), "application/json")

	return requestInfo
}

// mobileAppCategoryPostRequest executes a POST to associate a category.
func mobileAppCategoryPostRequest(ctx context.Context, client *msgraphbetasdk.GraphServiceClient, req *abstractions.RequestInformation) error {
	return client.GetAdapter().SendNoContent(ctx, req, nil)
}

// mobileAppCategoryDeleteRequest executes a DELETE to remove a category association.
func mobileAppCategoryDeleteRequest(ctx context.Context, client *msgraphbetasdk.GraphServiceClient, appID, categoryID string) error {
	requestInfo := abstractions.NewRequestInformation()
	requestInfo.Method = abstractions.DELETE
	requestInfo.UrlTemplate = "{+baseurl}/deviceAppManagement/mobileApps/" + appID + "/categories/" + categoryID + "/$ref"
	requestInfo.PathParameters = map[string]string{
		"baseurl": "https://graph.microsoft.com/beta",
	}

	return client.GetAdapter().SendNoContent(ctx, requestInfo, nil)
}
