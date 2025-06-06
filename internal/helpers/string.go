package helpers

import (
	"fmt"
	"regexp"
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
