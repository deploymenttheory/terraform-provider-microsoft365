package normalize

import (
	"encoding/json"
	"sort"
)

// JSONAlphabetically normalizes the JSON structure by sorting all keys alphabetically recursively
func JSONAlphabetically(input string) (string, error) {
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(input), &data); err != nil {
		return "", err
	}

	var normalize func(interface{}) interface{}
	normalize = func(v interface{}) interface{} {
		switch v := v.(type) {
		case map[string]interface{}:
			sorted := make(map[string]interface{})
			keys := make([]string, 0, len(v))
			for k := range v {
				keys = append(keys, k)
			}
			sort.Strings(keys)
			for _, k := range keys {
				sorted[k] = normalize(v[k])
			}
			return sorted

		case []interface{}:
			sortedArray := make([]interface{}, len(v))
			for i, val := range v {
				sortedArray[i] = normalize(val)
			}
			return sortedArray

		default:
			return v
		}
	}

	normalizedJSON, err := json.Marshal(normalize(data))
	if err != nil {
		return "", err
	}

	return string(normalizedJSON), nil
}
