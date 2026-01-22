package sharedStater

import (
	"context"
	"testing"

	construct "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors/graph_beta/device_and_app_management"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// TestMapMobileAppCategoriesStateToTerraform tests the MapMobileAppCategoriesStateToTerraform function
func TestMapMobileAppCategoriesStateToTerraform(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name       string
		categories []graphmodels.MobileAppCategoryable
		validate   func(t *testing.T, result types.Set)
	}{
		{
			name:       "Empty categories list",
			categories: []graphmodels.MobileAppCategoryable{},
			validate: func(t *testing.T, result types.Set) {
				assert.True(t, result.IsNull())
			},
		},
		{
			name:       "Nil categories list",
			categories: nil,
			validate: func(t *testing.T, result types.Set) {
				assert.True(t, result.IsNull())
			},
		},
		{
			name: "Built-in category mapping",
			categories: func() []graphmodels.MobileAppCategoryable {
				// Use the first built-in category ID for testing
				businessID := construct.BuiltInCategoryMapping["Business"]
				category := graphmodels.NewMobileAppCategory()
				category.SetId(&businessID)
				return []graphmodels.MobileAppCategoryable{category}
			}(),
			validate: func(t *testing.T, result types.Set) {
				assert.False(t, result.IsNull())
				assert.False(t, result.IsUnknown())
				elements := result.Elements()
				require.Len(t, elements, 1)
			},
		},
		{
			name: "Custom category UUID",
			categories: func() []graphmodels.MobileAppCategoryable {
				customID := "12345678-1234-1234-1234-123456789012"
				category := graphmodels.NewMobileAppCategory()
				category.SetId(&customID)
				return []graphmodels.MobileAppCategoryable{category}
			}(),
			validate: func(t *testing.T, result types.Set) {
				assert.False(t, result.IsNull())
				elements := result.Elements()
				require.Len(t, elements, 1)
			},
		},
		{
			name: "Multiple categories",
			categories: func() []graphmodels.MobileAppCategoryable {
				businessID := construct.BuiltInCategoryMapping["Business"]
				customID := "12345678-1234-1234-1234-123456789012"
				
				cat1 := graphmodels.NewMobileAppCategory()
				cat1.SetId(&businessID)
				
				cat2 := graphmodels.NewMobileAppCategory()
				cat2.SetId(&customID)
				
				return []graphmodels.MobileAppCategoryable{cat1, cat2}
			}(),
			validate: func(t *testing.T, result types.Set) {
				assert.False(t, result.IsNull())
				elements := result.Elements()
				require.Len(t, elements, 2)
			},
		},
		{
			name: "Nil category in list",
			categories: func() []graphmodels.MobileAppCategoryable {
				businessID := construct.BuiltInCategoryMapping["Business"]
				cat1 := graphmodels.NewMobileAppCategory()
				cat1.SetId(&businessID)
				
				return []graphmodels.MobileAppCategoryable{cat1, nil}
			}(),
			validate: func(t *testing.T, result types.Set) {
				assert.False(t, result.IsNull())
				elements := result.Elements()
				// Should only include the non-nil category
				require.Len(t, elements, 1)
			},
		},
		{
			name: "Category with nil ID",
			categories: func() []graphmodels.MobileAppCategoryable {
				cat := graphmodels.NewMobileAppCategory()
				cat.SetId(nil)
				return []graphmodels.MobileAppCategoryable{cat}
			}(),
			validate: func(t *testing.T, result types.Set) {
				// Should skip categories with nil IDs
				assert.True(t, result.IsNull())
			},
		},
		{
			name: "Mixed valid and invalid categories",
			categories: func() []graphmodels.MobileAppCategoryable {
				businessID := construct.BuiltInCategoryMapping["Business"]
				validCat := graphmodels.NewMobileAppCategory()
				validCat.SetId(&businessID)
				
				invalidCat := graphmodels.NewMobileAppCategory()
				invalidCat.SetId(nil)
				
				return []graphmodels.MobileAppCategoryable{invalidCat, validCat, nil}
			}(),
			validate: func(t *testing.T, result types.Set) {
				assert.False(t, result.IsNull())
				elements := result.Elements()
				// Should only include the valid category
				require.Len(t, elements, 1)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MapMobileAppCategoriesStateToTerraform(ctx, tt.categories)
			if tt.validate != nil {
				tt.validate(t, result)
			}
		})
	}
}
