package normalize

import (
	"encoding/json"
	"sort"
)

// JSONAlphabetically normalizes the JSON structure by recursively sorting the keys
// of all objects (maps) in alphabetical order. It ensures that:
//  1. All keys within objects are sorted consistently across all nesting levels.
//  2. Arrays (slices) retain their original order and are not sorted. However, any
//     nested objects within arrays are normalized by sorting their keys alphabetically.
//
// This approach preserves the semantics of JSON arrays while ensuring consistent
// ordering of keys in objects for easier comparison or processing.
func JSONAlphabetically(input string) (string, error) {
	var data any
	if err := json.Unmarshal([]byte(input), &data); err != nil {
		return "", err
	}

	var normalize func(any) any
	normalize = func(v any) any {
		switch v := v.(type) {
		case map[string]any:
			// Sort keys in maps
			sorted := make(map[string]any)
			keys := make([]string, 0, len(v))
			for k := range v {
				keys = append(keys, k)
			}
			sort.Strings(keys)
			for _, k := range keys {
				sorted[k] = normalize(v[k]) // Normalize recursively
			}
			return sorted

		case []any:
			// Retain array order but normalize nested objects within the array
			normalizedArray := make([]any, len(v))
			for i, val := range v {
				normalizedArray[i] = normalize(val) // Normalize each element
			}
			return normalizedArray

		default:
			// Return primitive types as-is
			return v
		}
	}

	normalizedJSON, err := json.Marshal(normalize(data))
	if err != nil {
		return "", err
	}

	return string(normalizedJSON), nil
}
