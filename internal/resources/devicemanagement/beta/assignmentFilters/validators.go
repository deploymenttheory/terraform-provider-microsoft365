package assignmentFilter

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Define supported platform types
var supportedPlatformTypes = map[string]int{
	"IOS":     1, // Assuming 1 is the integer value for IOS_DEVICEPLATFORMTYPE
	"Android": 2, // Assuming 2 is the integer value for ANDROID_DEVICEPLATFORMTYPE
	"Windows": 3, // Assuming 3 is the integer value for WINDOWS_DEVICEPLATFORMTYPE
	// Add other supported platform types here with their corresponding int values
}

// Custom validator for platform types
// Custom validator for platform types
func validatePlatform() validator.String {
	return validator.StringFunc(func(value types.String) (warns []string, errs []error) {
		if value.IsUnknown() || value.IsNull() {
			return
		}

		if _, exists := supportedPlatformTypes[value.ValueString()]; !exists {
			errs = append(errs, fmt.Errorf("invalid device platform type: %s. Supported types: %v", value.ValueString(), supportedPlatformTypes))
		}
		return
	})
}

// DevicePlatformType is the type for supported platform types.
type DevicePlatformType int

// StringToDevicePlatformType converts a string to DevicePlatformType based on a provided map.
func StringToDevicePlatformType(str string, mapping map[string]DevicePlatformType) (DevicePlatformType, error) {
	if val, exists := mapping[str]; exists {
		return val, nil
	}
	return -1, fmt.Errorf("invalid string: %s. Supported strings: %v", str, mapping)
}
