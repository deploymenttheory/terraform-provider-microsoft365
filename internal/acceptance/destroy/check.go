package destroy

import (
	"context"
	"fmt"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/types"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

const (
	checkDestroyTimeout = 5 * time.Minute
)

// resourceResult represents the destruction status of a single resource
type resourceResult struct {
	name   string
	id     string
	status string // "destroyed", "still_exists", "error"
	err    error
}

// CheckDestroyedFunc returns a TestCheckFunc which validates a specific resource no longer exists
func CheckDestroyedFunc(testResource types.TestResource, resourceType, resourceName string) func(state *terraform.State) error {
	if testResource == nil {
		panic("testResource cannot be nil")
	}
	if resourceType == "" {
		panic("resourceType cannot be empty")
	}
	if resourceName == "" {
		panic("resourceName cannot be empty")
	}

	return func(state *terraform.State) error {
		ctx, cancel := context.WithTimeout(context.Background(), checkDestroyTimeout)
		defer cancel()

		for label, resourceState := range state.RootModule().Resources {
			if resourceState.Type != resourceType {
				continue
			}

			if label != resourceName {
				continue
			}

			result, err := testResource.Exists(ctx, nil, resourceState.Primary)

			if result == nil && err == nil {
				return fmt.Errorf("should have either an error or a result when checking if %q has been destroyed", resourceName)
			}

			if result != nil && *result {
				return fmt.Errorf("%q still exists", resourceName)
			}
		}

		return nil
	}
}

// CheckDestroyedAllFunc returns a TestCheckFunc which validates all resources of a given type no longer exist
func CheckDestroyedAllFunc(testResource types.TestResource, resourceType string) func(state *terraform.State) error {
	if testResource == nil {
		panic("testResource cannot be nil")
	}
	if resourceType == "" {
		panic("resourceType cannot be empty")
	}

	return func(state *terraform.State) error {
		ctx, cancel := context.WithTimeout(context.Background(), checkDestroyTimeout)
		defer cancel()

		fmt.Println("=== RUN   CheckDestroy")

		var results []resourceResult
		shortType := deriveShortName(resourceType)

		for label, resourceState := range state.RootModule().Resources {
			if resourceState.Type != resourceType {
				continue
			}

			displayName := getDisplayName(resourceState, label)
			result := checkResourceDestroyed(ctx, testResource, resourceState, displayName)
			results = append(results, result)
		}

		if len(results) == 0 {
			return nil
		}

		fmt.Printf("=== TYPE     %s\n", shortType)
		return printResultsAndCheckError(results)
	}
}

// ResourceTypeMapping maps Terraform resource types to their TestResource implementations
type ResourceTypeMapping struct {
	ResourceType string
	TestResource types.TestResource
	ShortName    string // Optional: Override automatic short name generation
}

// CheckDestroyedTypesFunc returns a TestCheckFunc that verifies multiple resource types are destroyed
// Useful for tests that create dependencies (e.g., role scope tags with group assignments)
func CheckDestroyedTypesFunc(waitDuration time.Duration, mappings ...ResourceTypeMapping) func(*terraform.State) error {
	if len(mappings) == 0 {
		panic("at least one ResourceTypeMapping must be provided")
	}

	for i, mapping := range mappings {
		if mapping.TestResource == nil {
			panic(fmt.Sprintf("mappings[%d].TestResource cannot be nil", i))
		}
		if mapping.ResourceType == "" {
			panic(fmt.Sprintf("mappings[%d].ResourceType cannot be empty", i))
		}
	}

	return func(s *terraform.State) error {
		if waitDuration > 0 {
			time.Sleep(waitDuration)
		}

		ctx, cancel := context.WithTimeout(context.Background(), checkDestroyTimeout)
		defer cancel()

		fmt.Println("=== RUN   CheckDestroy")

		resourceMap, shortNameMap := buildResourceMaps(mappings)
		resourcesByType, processOrder := groupResourcesByType(ctx, s, resourceMap, shortNameMap)

		if len(processOrder) == 0 {
			return nil
		}

		return printResultsByType(resourcesByType, processOrder)
	}
}

// buildResourceMaps creates lookup maps for resource types and short names
func buildResourceMaps(mappings []ResourceTypeMapping) (map[string]types.TestResource, map[string]string) {
	resourceMap := make(map[string]types.TestResource)
	shortNameMap := make(map[string]string)

	for _, mapping := range mappings {
		resourceMap[mapping.ResourceType] = mapping.TestResource
		if mapping.ShortName != "" {
			shortNameMap[mapping.ResourceType] = mapping.ShortName
		} else {
			shortNameMap[mapping.ResourceType] = deriveShortName(mapping.ResourceType)
		}
	}

	return resourceMap, shortNameMap
}

// groupResourcesByType checks all resources and groups results by type
func groupResourcesByType(
	ctx context.Context,
	state *terraform.State,
	resourceMap map[string]types.TestResource,
	shortNameMap map[string]string,
) (map[string][]resourceResult, []string) {
	resourcesByType := make(map[string][]resourceResult)
	var processOrder []string

	for key, rs := range state.RootModule().Resources {
		testResource, exists := resourceMap[rs.Type]
		if !exists {
			continue
		}

		shortType := shortNameMap[rs.Type]

		if _, seen := resourcesByType[shortType]; !seen {
			processOrder = append(processOrder, shortType)
			resourcesByType[shortType] = []resourceResult{}
		}

		displayName := getDisplayName(rs, key)
		result := checkResourceDestroyed(ctx, testResource, rs, displayName)
		resourcesByType[shortType] = append(resourcesByType[shortType], result)
	}

	return resourcesByType, processOrder
}
