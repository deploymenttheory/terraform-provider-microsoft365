package graphBetaSettingsCatalogConfigurationPolicy

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

// resetGlobalCache clears globals to ensure tests are isolated.
func resetGlobalCache(t *testing.T) {
	t.Helper()
	cacheMutex.Lock()
	defer cacheMutex.Unlock()
	for k := range globalConfigCache {
		delete(globalConfigCache, k)
	}
	hasConfiguration = false
}

func TestGetResolvedDCV2MaxDepth_WhenNoConfiguration_ReturnsZero(t *testing.T) {
	resetGlobalCache(t)
	require.Equal(t, 0, GetResolvedDCV2MaxDepth())
}

func TestResolveDCV2ConfigurationDepth_EmptyContent_NoCacheChange(t *testing.T) {
	resetGlobalCache(t)
	ResolveDCV2ConfigurationDepth("cfg-empty", "")
	require.Equal(t, 0, GetResolvedDCV2MaxDepth())
}

func TestResolveDCV2ConfigurationDepth_NoChildren_MinDepthIsThree(t *testing.T) {
	resetGlobalCache(t)
	// No occurrences of "children =" should yield safeDepth = 3
	cfg := `setting_instance { odata_type = "#type" setting_definition_id = "id" }`
	ResolveDCV2ConfigurationDepth("cfg-no-children", cfg)
	require.Equal(t, 3, GetResolvedDCV2MaxDepth())
}

func TestResolveDCV2ConfigurationDepth_SingleLevelChildren_DepthTwo(t *testing.T) {
	resetGlobalCache(t)
	// One level of children increases maxDepth to 1 → safeDepth = 2
	cfg := strings.Join([]string{
		"children = [",
		"  {",
		"    odata_type = \"#type\"",
		"    setting_definition_id = \"id\"",
		"  }",
		"}", // closing brace at start of line reduces currentDepth
	}, "\n")
	ResolveDCV2ConfigurationDepth("cfg-single", cfg)
	require.Equal(t, 2, GetResolvedDCV2MaxDepth())
}

func TestResolveDCV2ConfigurationDepth_DeepChildren_CappedAtMax(t *testing.T) {
	resetGlobalCache(t)
	// Build content that drives currentDepth well above max; then close with matching braces
	var b strings.Builder
	levels := 20
	for i := 0; i < levels; i++ {
		b.WriteString("children = [\n")
	}
	for i := 0; i < levels; i++ {
		b.WriteString("}\n")
	}
	ResolveDCV2ConfigurationDepth("cfg-deep", b.String())
	require.Equal(t, maxSchemaDepth, GetResolvedDCV2MaxDepth())
}

func TestGetMaxSchemaDepth_Priority_EnvOverrideWins(t *testing.T) {
	resetGlobalCache(t)
	// Cache some other value that should be ignored by explicit env override
	ResolveDCV2ConfigurationDepth("cfg", "children = [\n}\n") // yields depth 2
	t.Setenv("TF_SCHEMA_MAX_DEPTH", "9")
	require.Equal(t, 9, getMaxSchemaDepth())
}

func TestGetMaxSchemaDepth_Priority_TestingFlags(t *testing.T) {
	resetGlobalCache(t)
	// Even with cached depth, testing flag should force 4
	ResolveDCV2ConfigurationDepth("cfg", "children = [\n}\n")
	t.Setenv("GO_TESTING", "1")
	require.Equal(t, 4, getMaxSchemaDepth())
}

func TestGetMaxSchemaDepth_Priority_TFACCFlag(t *testing.T) {
	resetGlobalCache(t)
	ResolveDCV2ConfigurationDepth("cfg", "children = [\n}\n")
	t.Setenv("TF_ACC", "1")
	require.Equal(t, 4, getMaxSchemaDepth())
}

func TestGetMaxSchemaDepth_UsesCachedDepthWhenAvailable(t *testing.T) {
	resetGlobalCache(t)
	// Build 4 nested children lines (maxDepth=4) → safeDepth=5
	var b strings.Builder
	for i := 0; i < 4; i++ {
		b.WriteString("children = [\n")
	}
	for i := 0; i < 4; i++ {
		b.WriteString("}\n")
	}
	ResolveDCV2ConfigurationDepth("cfg-cached", b.String())
	// Ensure no testing env flags interfere
	t.Setenv("GO_TESTING", "")
	t.Setenv("TF_ACC", "")
	t.Setenv("TF_SCHEMA_MAX_DEPTH", "")
	require.Equal(t, 5, getMaxSchemaDepth())
}

func TestGetMaxSchemaDepth_FallbackToProductionMax(t *testing.T) {
	resetGlobalCache(t)
	// No env, no cache
	t.Setenv("GO_TESTING", "")
	t.Setenv("TF_ACC", "")
	t.Setenv("TF_SCHEMA_MAX_DEPTH", "")
	require.Equal(t, maxSchemaDepth, getMaxSchemaDepth())
}

func TestResolveDCV2ConfigurationDepth_MultipleConfigs_TracksMaxAcrossCache(t *testing.T) {
	resetGlobalCache(t)
	ResolveDCV2ConfigurationDepth("cfg-a", "children = [\n}\n") // safeDepth 2
	// Create content that yields larger safeDepth
	var b strings.Builder
	for i := 0; i < 3; i++ { // maxDepth 3 → safeDepth 4
		b.WriteString("children = [\n")
	}
	for i := 0; i < 3; i++ {
		b.WriteString("}\n")
	}
	ResolveDCV2ConfigurationDepth("cfg-b", b.String())
	require.Equal(t, 4, GetResolvedDCV2MaxDepth())
}

// Guard against accidental changes to constant.
func TestMaxSchemaDepth_Constant(t *testing.T) {
	require.Equal(t, 15, maxSchemaDepth)
}

// Sanity check: safeDepth calculation boundaries.
func TestResolveDCV2ConfigurationDepth_SafeDepthLowerBound(t *testing.T) {
	resetGlobalCache(t)
	ResolveDCV2ConfigurationDepth("cfg", "") // no-op
	// Use content that doesn't trigger children; should yield 3
	ResolveDCV2ConfigurationDepth("cfg", fmt.Sprintf("%d", 123))
	require.Equal(t, 3, GetResolvedDCV2MaxDepth())
}
