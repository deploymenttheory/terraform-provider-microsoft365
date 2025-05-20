package graphBetaMobileApp

import (
	"reflect"
	"strings"

	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// getAppTypeFromMobileApp determines the app type from the MobileAppable interface
func getAppTypeFromMobileApp(data graphmodels.MobileAppable) string {
	if data.GetOdataType() != nil {
		odataType := *data.GetOdataType()

		const prefix = "#microsoft.graph."
		if len(odataType) > len(prefix) && odataType[:len(prefix)] == prefix {
			appType := odataType[len(prefix):]
			return appType
		}

		// If we have an OData type but it doesn't have the expected prefix,
		// still return it without the # symbol
		if len(odataType) > 0 && odataType[0] == '#' {
			return odataType[1:]
		}
	}

	// If odata.type isn't available or is empty, fall back to type checking
	// This is a more reliable approach than the complex type switch
	// We'll use reflection to get the actual type name
	typeName := reflect.TypeOf(data).String()

	// Extract just the type name from the package path
	parts := strings.Split(typeName, ".")
	if len(parts) > 1 {
		// Remove the "able" suffix if present
		name := strings.TrimSuffix(parts[len(parts)-1], "able")

		// Convert from PascalCase to camelCase for consistency with OData types
		if len(name) > 1 {
			return strings.ToLower(name[:1]) + name[1:]
		}
		return strings.ToLower(name)
	}

	// If all else fails
	return "unknown"
}
