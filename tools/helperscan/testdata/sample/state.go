package sample

// MapRemoteStateToTerraform is a standard function and is not reported.
func MapRemoteStateToTerraform(values []string) []string {
	return mapWidgetsToState(values)
}

// findWidgetByName calls itself, which does not count as a reference.
func findWidgetByName(names []string, name string) string {
	if len(names) == 0 {
		return ""
	}
	if names[0] == name {
		return name
	}
	return findWidgetByName(names[1:], name)
}
