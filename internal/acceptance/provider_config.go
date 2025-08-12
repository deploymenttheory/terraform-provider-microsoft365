package acceptance

import (
	"fmt"
	"os"
	"strings"
)

// ProviderConfig returns the provider configuration block for acceptance tests
// This configuration uses environment variables for authentication
func ProviderConfig() string {
	return `
provider "microsoft365" {
  # Configuration from environment variables
}

provider "random" {
  # Random provider for generating unique resource names in tests
}
`
}

// M365ProviderBlock returns the provider configuration with explicit values
// for cases where environment variables need to be overridden
func M365ProviderBlock(tenantID, clientID, clientSecret, authMethod, cloud string) string {
	return fmt.Sprintf(`
provider "microsoft365" {
  tenant_id    = "%s"
  auth_method  = "%s"
  cloud        = "%s"
  entra_id_options = {
    client_id     = "%s"
    client_secret = "%s"
  }
}

provider "random" {
  # Random provider for generating unique resource names in tests
}
`, tenantID, authMethod, cloud, clientID, clientSecret)
}

// M365ProviderBlockValueInjection returns the provider configuration using current environment variables
func M365ProviderBlockValueInjection() string {
	tenantID := os.Getenv("M365_TENANT_ID")
	clientID := os.Getenv("M365_CLIENT_ID")
	clientSecret := os.Getenv("M365_CLIENT_SECRET")
	authMethod := os.Getenv("M365_AUTH_METHOD")
	cloud := os.Getenv("M365_CLOUD")

	if tenantID == "" || clientID == "" || clientSecret == "" || authMethod == "" || cloud == "" {
		// Fall back to implicit configuration from environment
		return ProviderConfig()
	}

	return M365ProviderBlock(tenantID, clientID, clientSecret, authMethod, cloud)
}

// ConfiguredM365ProviderBlock prefixes any terraform configuration with a configured
// M365 provider block only if the config doesn't already contain a provider block.
func ConfiguredM365ProviderBlock(config string) string {
	// If the config already contains a provider block, just return it as-is
	if strings.Contains(config, `provider "microsoft365"`) {
		return config
	}
	return M365ProviderBlockValueInjection() + config
}
