package provider

// Define supported platform types
var supportedPlatformTypes = []string{
	models.IOS_DEVICEPLATFORMTYPE,
	models.ANDROID_DEVICEPLATFORMTYPE,
	models.WINDOWS_DEVICEPLATFORMTYPE,
	// Add other supported platform types here
}

// Custom validator for platform types
func validatePlatform() validator.String {
	return validator.StringFunc(func(value types.String) (warns []string, errs []error) {
		if value.IsUnknown() || value.IsNull() {
			return
		}

		for _, v := range supportedPlatformTypes {
			if value.ValueString() == v {
				return
			}
		}
		errs = append(errs, fmt.Errorf("invalid deviplatform type: %s. Supported types: %v", value.ValueString(), supportedPlatformTypes))
		return
	})
}