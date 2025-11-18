package testlog

import (
	"fmt"
	"strings"
	"time"
)

// StepAction logs a test step action with consistent formatting
func StepAction(resourceType, action string) {
	resourceName := deriveResourceName(resourceType)
	fmt.Printf("--- %s %s...\n", action, resourceName)
}

// WaitForConsistency logs a wait message for eventual consistency
func WaitForConsistency(service string, duration time.Duration) {
	fmt.Printf("--- Waiting %v for resource %s to achieve eventual consistency...\n", duration, service)
}

// deriveResourceName extracts a human-readable resource name from a Terraform resource type
// Example: "microsoft365_graph_beta_device_management_role_scope_tag" -> "role scope tag"
func deriveResourceName(resourceType string) string {
	parts := strings.Split(resourceType, "_")

	if len(parts) < 3 {
		return resourceType
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

	// Known multi-word services to skip
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
		return resourceType
	}

	// Everything from resourceStartIdx onwards is the resource name
	resourceParts := parts[resourceStartIdx:]

	// Join with spaces for human-readable name
	return strings.Join(resourceParts, " ")
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
