package mocks

// EnsureField ensures that a field exists in a data map
// If the field doesn't exist, it will be initialized with the provided default value
// This is particularly useful for ensuring collection fields are initialized in mock responses
func EnsureField(data map[string]interface{}, fieldName string, defaultValue interface{}) {
	if data[fieldName] == nil {
		data[fieldName] = defaultValue
	}
}
