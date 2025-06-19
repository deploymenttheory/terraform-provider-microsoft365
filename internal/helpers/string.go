package helpers

import (
	"fmt"
	"regexp"
	"strings"
)

// ExtractGUIDFromString extracts a GUID from the beginning of a string.
func ExtractGUIDFromString(stringToExtractFrom string) (string, error) {
	// Regex to match a GUID
	re := regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}`)
	guid := re.FindString(stringToExtractFrom)
	if guid == "" {
		return "", fmt.Errorf("could not find a valid GUID in the provided string: %s", stringToExtractFrom)
	}
	return guid, nil
}

// StringToInt converts a string to an integer based on a provided map
func StringToInt(str string, mapping map[string]int) (int, error) {
	if val, exists := mapping[str]; exists {
		return val, nil
	}
	return -1, fmt.Errorf("invalid string: %s. Supported strings: %v", str, mapping)
}

// SplitCommaSeparatedString splits a comma-separated string into a slice of strings
func SplitCommaSeparatedString(s string) []string {
	if s == "" {
		return []string{}
	}

	// Split the string by commas and return the result
	return strings.Split(s, ",")
}
