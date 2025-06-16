package sharedStater

import (
	"context"

	construct "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors/graph_beta/device_and_app_management"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// MapMobileAppCategoriesStateToTerraform converts API categories to a set of display names (for built-ins)
// or UUIDs (for custom categories) for Terraform state.
func MapMobileAppCategoriesStateToTerraform(ctx context.Context, categories []graphmodels.MobileAppCategoryable) types.Set {
	if len(categories) == 0 {
		return types.SetNull(types.StringType)
	}

	categoryValues := make([]attr.Value, 0, len(categories))

	for _, category := range categories {
		if category == nil || category.GetId() == nil {
			continue
		}

		categoryID := *category.GetId()
		mapped := false

		// Match ID to built-in name
		for name, builtInID := range construct.BuiltInCategoryMapping {
			if categoryID == builtInID {
				categoryValues = append(categoryValues, types.StringValue(name))
				mapped = true
				break
			}
		}

		if !mapped {
			// Fallback to ID for custom category
			categoryValues = append(categoryValues, types.StringValue(categoryID))
		}
	}

	if len(categoryValues) == 0 {
		return types.SetNull(types.StringType)
	}

	set, diags := types.SetValue(types.StringType, categoryValues)
	if diags.HasError() {
		tflog.Error(ctx, "Failed to build category set", map[string]interface{}{
			"errors": diags.Errors(),
		})
		return types.SetNull(types.StringType)
	}

	return set
}
