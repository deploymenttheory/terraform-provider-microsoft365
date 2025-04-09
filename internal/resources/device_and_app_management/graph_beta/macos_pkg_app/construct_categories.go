package graphBetaMacOSPKGApp

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// constructCategories creates an array of MobileAppCategoryable objects
// from an array of display names. The function maps each display name to its
// corresponding ID based on predefined mappings or treats it as a user-defined category.
func constructCategories(ctx context.Context, displayNames []string) []graphmodels.MobileAppCategoryable {
	// Built in categories
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

	tflog.Debug(ctx, fmt.Sprintf("Constructing mobile app categories for display names: %v", displayNames))
	categories := make([]graphmodels.MobileAppCategoryable, 0, len(displayNames))

	builtInCount := 0
	userDefinedCount := 0

	for _, name := range displayNames {
		category := graphmodels.NewMobileAppCategory()
		displayNameCopy := name
		category.SetDisplayName(&displayNameCopy)

		if id, exists := categoryMapping[name]; exists {
			idCopy := id
			category.SetId(&idCopy)
			tflog.Debug(ctx, fmt.Sprintf("Using built-in category '%s' with ID '%s'", name, id))
			builtInCount++
		} else {
			tflog.Debug(ctx, fmt.Sprintf("Using user-defined category '%s'", name))
			userDefinedCount++
		}

		categories = append(categories, category)
	}

	tflog.Debug(ctx, fmt.Sprintf("Constructed %d total categories (%d built-in, %d user-defined)",
		len(categories), builtInCount, userDefinedCount))
	return categories
}
