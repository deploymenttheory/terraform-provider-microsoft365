package graphBetaSettingsCatalogConfigurationPolicy

import (
	"context"
	"fmt"
	"strings"
)

// buildNameFilter constructs an OData filter for partial name matching
func buildNameFilter(data *SettingsCatalogListConfigModel) string {
	if data.NameFilter.IsNull() || data.NameFilter.IsUnknown() {
		return ""
	}
	return fmt.Sprintf("contains(name,'%s')", data.NameFilter.ValueString())
}

// buildPlatformFilter constructs an OData filter for platform(s) using OR logic
func buildPlatformFilter(ctx context.Context, data *SettingsCatalogListConfigModel) string {
	if data.PlatformFilter.IsNull() || data.PlatformFilter.IsUnknown() {
		return ""
	}

	var platforms []string
	diags := data.PlatformFilter.ElementsAs(ctx, &platforms, false)
	if diags.HasError() || len(platforms) == 0 {
		return ""
	}

	var platformFilters []string
	for _, platform := range platforms {
		platformFilters = append(platformFilters, fmt.Sprintf("platforms eq '%s'", platform))
	}

	if len(platformFilters) == 0 {
		return ""
	}

	return fmt.Sprintf("(%s)", strings.Join(platformFilters, " or "))
}

// buildTemplateFamilyFilter constructs an OData filter for template family
func buildTemplateFamilyFilter(data *SettingsCatalogListConfigModel) string {
	if data.TemplateFamilyFilter.IsNull() || data.TemplateFamilyFilter.IsUnknown() {
		return ""
	}
	return fmt.Sprintf("templateReference/templateFamily eq '%s'", data.TemplateFamilyFilter.ValueString())
}

// combineFilters joins non-empty filter parts with AND logic
func combineFilters(filters ...string) string {
	var parts []string
	for _, filter := range filters {
		if filter != "" {
			parts = append(parts, filter)
		}
	}
	if len(parts) == 0 {
		return ""
	}
	return strings.Join(parts, " and ")
}
