package helpers

// Utility function to convert a bool to *bool
func BoolPtr(b bool) *bool {
	return &b
}

// Utility function to convert a string to *string
func StringPtr(s string) *string {
	return &s
}

// GetStringValue safely extracts a string value from a pointer, returning empty string if nil
func GetStringValue(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}
