package destroy

import (
	"context"
	"fmt"
	"strings"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/types"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// checkResourceDestroyed verifies if a single resource has been destroyed
func checkResourceDestroyed(
	ctx context.Context,
	testResource types.TestResource,
	resourceState *terraform.ResourceState,
	displayName string,
) resourceResult {
	result := resourceResult{
		name: displayName,
		id:   resourceState.Primary.ID,
	}

	resourceExists, err := testResource.Exists(ctx, nil, resourceState.Primary)

	if err != nil {
		result.status = "error"
		result.err = err
	} else if resourceExists != nil && *resourceExists {
		result.status = "still_exists"
	} else {
		result.status = "destroyed"
	}

	return result
}

// getDisplayName extracts a display name from resource state, falling back to resource name
func getDisplayName(resourceState *terraform.ResourceState, resourceKey string) string {
	displayName := resourceState.Primary.Attributes["display_name"]
	if displayName == "" {
		displayName = getResourceName(resourceKey)
	}
	return displayName
}

// getResourceName extracts the resource name from the state resource key
func getResourceName(resourceKey string) string {
	parts := strings.Split(resourceKey, ".")
	if len(parts) > 1 {
		return parts[len(parts)-1]
	}
	return resourceKey
}

// printResultsAndCheckError prints results and returns an error if any resources failed to destroy
func printResultsAndCheckError(results []resourceResult) error {
	destroyedCount := 0
	failedCount := 0

	for _, r := range results {
		switch r.status {
		case "destroyed":
			destroyedCount++
			fmt.Printf("    --- PASS: CheckDestroy/%s (id:%s)\n", r.name, r.id)
		case "still_exists":
			failedCount++
			fmt.Printf("    --- FAIL: CheckDestroy/%s (id:%s) - still exists\n", r.name, r.id)
		case "error":
			failedCount++
			fmt.Printf("    --- FAIL: CheckDestroy/%s (id:%s) - %v\n", r.name, r.id, r.err)
		}
	}

	if failedCount > 0 {
		fmt.Printf("--- FAIL: CheckDestroy (%d destroyed, %d failed)\n", destroyedCount, failedCount)
		return fmt.Errorf("cleanup verification failed: %d resources not properly destroyed", failedCount)
	}

	fmt.Printf("--- PASS: CheckDestroy (%d resources)\n", len(results))
	return nil
}

// printResultsByType prints results grouped by type and returns an error if any resources failed to destroy
func printResultsByType(resourcesByType map[string][]resourceResult, processOrder []string) error {
	totalResources := 0
	destroyedCount := 0
	failedCount := 0

	for _, resourceType := range processOrder {
		results := resourcesByType[resourceType]
		if len(results) == 0 {
			continue
		}

		fmt.Printf("=== RUN   CheckDestroy/%s\n", resourceType)

		for _, result := range results {
			totalResources++

			switch result.status {
			case "destroyed":
				destroyedCount++
				fmt.Printf("    --- PASS: CheckDestroy/%s/%s (id:%s)\n", resourceType, result.name, result.id)
			case "still_exists":
				failedCount++
				fmt.Printf("    --- FAIL: CheckDestroy/%s/%s (id:%s) - still exists\n", resourceType, result.name, result.id)
			case "error":
				failedCount++
				fmt.Printf("    --- FAIL: CheckDestroy/%s/%s (id:%s) - %v\n", resourceType, result.name, result.id, result.err)
			}
		}
	}

	if failedCount > 0 {
		fmt.Printf("--- FAIL: CheckDestroy (%d destroyed, %d failed)\n", destroyedCount, failedCount)
		return fmt.Errorf("cleanup verification failed: %d resources not properly destroyed", failedCount)
	}

	fmt.Printf("--- PASS: CheckDestroy (%d resources)\n", totalResources)
	return nil
}

// deriveShortName intelligently derives a human-readable short name from a resource type
func deriveShortName(fullType string) string {
	parts := strings.Split(fullType, "_")

	if len(parts) < 3 {
		return fullType
	}

	// Skip provider prefix (microsoft365)
	idx := 1

	// Skip API version (graph_beta, graph_v1_0, graph, powershell)
	if idx < len(parts) {
		api := parts[idx]
		if api == "graph" || api == "powershell" {
			idx++
			if idx < len(parts) && (parts[idx] == "beta" || parts[idx] == "v1") {
				idx++
				if idx < len(parts) && parts[idx-1] == "v1" && parts[idx] == "0" {
					idx++
				}
			}
		}
	}

	// Known multi-word services
	knownServices := map[string]int{
		"device_and_app_management": 4,
		"identity_and_access":       3,
		"device_management":         2,
		"windows_365":               2,
		"m365_admin":                2,
		"groups":                    1,
		"users":                     1,
		"applications":              1,
	}

	// Check for known services
	serviceLength := 1
	for serviceName, length := range knownServices {
		serviceKey := strings.Join(parts[idx:min(idx+length, len(parts))], "_")
		if serviceKey == serviceName {
			serviceLength = length
			break
		}
	}

	// Skip past the service to get resource name
	resourceStartIdx := idx + serviceLength
	if resourceStartIdx >= len(parts) {
		return fullType
	}

	// Everything from resourceStartIdx onwards is the resource name
	resourceParts := parts[resourceStartIdx:]

	// Convert to PascalCase
	var result strings.Builder
	for _, part := range resourceParts {
		if len(part) > 0 {
			result.WriteString(strings.ToUpper(string(part[0])))
			if len(part) > 1 {
				result.WriteString(part[1:])
			}
		}
	}

	shortName := result.String()
	if shortName == "" {
		return fullType
	}
	return shortName
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
