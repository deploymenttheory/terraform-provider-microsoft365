package graphBetaTargetedManagedAppConfigurations

import (
	"os"
	"strconv"
	"strings"
	"sync"
)

const maxSchemaDepth = 15

// Global configuration cache for dynamic depth analysis
var (
	globalConfigCache = make(map[string]int)
	cacheMutex        sync.RWMutex
	hasConfiguration  bool
)

// ResolveDCV2ConfigurationDepth analyzes HCL configuration and caches the dcv2 depth requirement.
// This is called during resource operations when configuration is available.
//
// ANALYSIS PROCESS:
// 1. Parses HCL configuration to find maximum nesting depth of "children" blocks
// 2. Adds safety buffer and enforces Microsoft specification bounds (1-15 levels)
// 3. Caches the result for optimal performance on subsequent schema constructions
//
// This enables right-sizing the schema to prevent exponential timeout issues while maintaining functionality.
func ResolveDCV2ConfigurationDepth(configName, configContent string) {
	if configContent == "" {
		return
	}

	// Analyze the HCL configuration for maximum nesting depth
	maxDepth := 0
	currentDepth := 0

	lines := strings.Split(configContent, "\n")
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		// Track nesting by looking for children blocks and setting structure
		if strings.Contains(trimmed, "children") && strings.Contains(trimmed, "=") {
			currentDepth++
			if currentDepth > maxDepth {
				maxDepth = currentDepth
			}
		}

		// Reset depth tracking on closing braces at start of line (dedent)
		if strings.HasPrefix(trimmed, "}") && currentDepth > 0 {
			currentDepth--
		}
	}

	// Add safety buffer and enforce bounds
	safeDepth := maxDepth + 1
	if safeDepth < 2 {
		safeDepth = 3 // Minimum viable depth
	}
	if safeDepth > maxSchemaDepth {
		safeDepth = maxSchemaDepth // Cap at Microsoft maximum
	}

	// Cache the analyzed depth for optimal performance
	cacheMutex.Lock()
	globalConfigCache[configName] = safeDepth
	hasConfiguration = true
	cacheMutex.Unlock()
}

// getMaxSchemaDepth returns the maximum schema depth by analyzing the actual HCL configuration.
//
// PROBLEM SUMMARY:
// Microsoft 365 Settings Catalog policies support 15 levels of nested recursion as per Microsoft documentation.
// However, this creates exponential schema combinations during Terraform Plugin Framework initialization:
// - Level 1: ~10 attribute types
// - Level 5: ~100,000 combinations
// - Level 15: 10^15+ combinations (causing 10+ minute timeouts)
//
// The deep recursion occurs in these schema paths:
// - choiceSettingAttributes() → choiceSettingChildAttributes() → choiceSettingAttributes() (recursive)
// - groupSettingCollectionAttributes() → groupSettingCollectionChildAttributes() → groupSettingCollectionAttributes() (recursive)
// - Each level multiplies the schema validation complexity exponentially
//
// TESTING IMPACT:
// - Unit tests using resource.UnitTest() timeout after 10+ minutes
// - Acceptance tests using resource.Test() timeout after 10+ minutes
// - Schema construction becomes CPU and memory intensive
// - Tests fail due to terraform-plugin-testing framework limitations
//
// SOLUTION APPROACH:
// This function provides environment-aware depth limiting for testing scenarios while maintaining
// full Microsoft-compliant 15-level depth in production environments.
//
// USAGE:
// 1. Production (normal operation): Returns 15 (full Microsoft specification compliance)
// 2. Custom testing depth: Set TF_SCHEMA_MAX_DEPTH=N environment variable for specific depth
// 3. Unit testing: GO_TESTING=1 triggers depth limit of 3 levels
// 4. Acceptance testing: TF_ACC=1 triggers depth limit of 3 levels
//
// The 3-level limit for testing provides sufficient schema validation coverage while preventing
// exponential recursion timeouts. Real-world Settings Catalog policies rarely exceed 3-4 levels
// of nesting, making this a practical testing limitation.
//
// TRADE-OFFS:
// - Testing: Faster execution (milliseconds vs minutes), but limited deep nesting validation
// - Production: Full 15-level compliance, slower initialization but comprehensive functionality
// getMaxSchemaDepth determines the optimal schema depth using logical prioritization.
//
// LOGICAL FLOW:
// 1. Testing Environment Detection: If running tests, use reduced depth for performance
// 2. Cached Configuration Analysis: Use actual configuration depth when available
// 3. Production Fallback: Use full Microsoft specification depth (15 levels)
//
// This ensures tests run fast while production gets right-sized schemas based on
// actual configuration requirements, falling back to full depth for maximum compatibility.
func getMaxSchemaDepth() int {
	// Priority 1: Testing environment detection - custom override for specific testing scenarios
	if testDepth := os.Getenv("TF_SCHEMA_MAX_DEPTH"); testDepth != "" {
		if depth, err := strconv.Atoi(testDepth); err == nil && depth > 0 {
			return depth
		}
	}

	// Testing environment - use safe reduced depth for performance
	if os.Getenv("GO_TESTING") == "1" || os.Getenv("TF_ACC") != "" {
		return 4 // Covers most real-world scenarios while preventing timeouts
	}

	// Priority 2: Use cached configuration analysis (automatic right-sizing)
	// When configuration has been analyzed during resource operations
	if cachedDepth := GetResolvedDCV2MaxDepth(); cachedDepth > 0 {
		return cachedDepth
	}

	// Priority 3: Production fallback - full Microsoft specification compliance
	return maxSchemaDepth // 15 levels as per Microsoft documentation requirements
}

// GetResolvedDCV2MaxDepth returns the maximum depth of the dcv2 configuration
func GetResolvedDCV2MaxDepth() int {
	cacheMutex.RLock()
	defer cacheMutex.RUnlock()

	if !hasConfiguration {
		return 0 // No cached configuration available
	}

	maxCachedDepth := 0
	for _, depth := range globalConfigCache {
		if depth > maxCachedDepth {
			maxCachedDepth = depth
		}
	}

	return maxCachedDepth
}
