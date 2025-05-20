package graphBetaMobileApp

import (
	"reflect"
	"strings"

	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// getAppTypeFromMobileApp determines the app type from the MobileAppable interface
func getAppTypeFromMobileApp(data graphmodels.MobileAppable) string {
	// First check if the odata.type is available
	if data.GetOdataType() != nil {
		odataType := *data.GetOdataType()

		// Strip the "#microsoft.graph." prefix if present
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

// getODataTypeForAppType maps friendly app type names to OData types
func getODataTypeForAppType(appType string) string {
	switch appType {
	// MacOS Apps
	case "macOSPkgApp":
		return "#microsoft.graph.macOSPkgApp"
	case "macOSDmgApp":
		return "#microsoft.graph.macOSDmgApp"
	case "macOSLobApp":
		return "#microsoft.graph.macOSLobApp"
	case "macOSMicrosoftDefenderApp":
		return "#microsoft.graph.macOSMicrosoftDefenderApp"
	case "macOSMicrosoftEdgeApp":
		return "#microsoft.graph.macOSMicrosoftEdgeApp"
	case "macOSOfficeSuiteApp":
		return "#microsoft.graph.macOSOfficeSuiteApp"
	case "macOsVppApp":
		return "#microsoft.graph.macOsVppApp"
	case "macOSWebClip":
		return "#microsoft.graph.macOSWebClip"

	// Android Apps
	case "androidForWorkApp":
		return "#microsoft.graph.androidForWorkApp"
	case "androidLobApp":
		return "#microsoft.graph.androidLobApp"
	case "androidManagedStoreApp":
		return "#microsoft.graph.androidManagedStoreApp"
	case "androidManagedStoreWebApp":
		return "#microsoft.graph.androidManagedStoreWebApp"
	case "androidStoreApp":
		return "#microsoft.graph.androidStoreApp"
	case "managedAndroidLobApp":
		return "#microsoft.graph.managedAndroidLobApp"
	case "managedAndroidStoreApp":
		return "#microsoft.graph.managedAndroidStoreApp"

	// iOS Apps
	case "iosiPadOSWebClip":
		return "#microsoft.graph.iosiPadOSWebClip"
	case "iosLobApp":
		return "#microsoft.graph.iosLobApp"
	case "iosStoreApp":
		return "#microsoft.graph.iosStoreApp"
	case "iosVppApp":
		return "#microsoft.graph.iosVppApp"
	case "managedIOSLobApp":
		return "#microsoft.graph.managedIOSLobApp"
	case "managedIOSStoreApp":
		return "#microsoft.graph.managedIOSStoreApp"

	// Windows Apps
	case "windowsAppX":
		return "#microsoft.graph.windowsAppX"
	case "windowsMicrosoftEdgeApp":
		return "#microsoft.graph.windowsMicrosoftEdgeApp"
	case "windowsMobileMSI":
		return "#microsoft.graph.windowsMobileMSI"
	case "windowsPhone81AppX":
		return "#microsoft.graph.windowsPhone81AppX"
	case "windowsPhone81AppXBundle":
		return "#microsoft.graph.windowsPhone81AppXBundle"
	case "windowsPhone81StoreApp":
		return "#microsoft.graph.windowsPhone81StoreApp"
	case "windowsPhoneXAP":
		return "#microsoft.graph.windowsPhoneXAP"
	case "windowsStoreApp":
		return "#microsoft.graph.windowsStoreApp"
	case "windowsUniversalAppX":
		return "#microsoft.graph.windowsUniversalAppX"
	case "windowsWebApp":
		return "#microsoft.graph.windowsWebApp"
	case "winGetApp":
		return "#microsoft.graph.winGetApp"

	// Web Apps
	case "webApp":
		return "#microsoft.graph.webApp"

	// Microsoft Store Apps
	case "microsoftStoreForBusinessApp":
		return "#microsoft.graph.microsoftStoreForBusinessApp"

	// Office Apps
	case "officeSuiteApp":
		return "#microsoft.graph.officeSuiteApp"

	// Win32 Apps
	case "win32CatalogApp":
		return "#microsoft.graph.win32CatalogApp"
	case "win32LobApp":
		return "#microsoft.graph.win32LobApp"

	// Other App Types
	case "managedApp":
		return "#microsoft.graph.managedApp"
	case "managedMobileLobApp":
		return "#microsoft.graph.managedMobileLobApp"
	case "mobileLobApp":
		return "#microsoft.graph.mobileLobApp"

	// If none of the above, return empty string
	default:
		return ""
	}
}
