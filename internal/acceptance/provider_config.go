package acceptance

import (
	"fmt"
	"os"
)

// ProviderConfig returns the provider configuration block for acceptance tests
// This configuration uses environment variables for authentication
func ProviderConfig() string {
	return `
provider "microsoft365" {
  # Configuration from environment variables
}
`
}

// ProviderConfigWithExplicitValues returns the provider configuration with explicit values
// for cases where environment variables need to be overridden
func ProviderConfigWithExplicitValues(tenantID, clientID, clientSecret, authMethod, cloud string) string {
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
`, tenantID, authMethod, cloud, clientID, clientSecret)
}

// ProviderConfigForCurrentEnv returns the provider configuration using current environment variables
func ProviderConfigForCurrentEnv() string {
	tenantID := os.Getenv("M365_TENANT_ID")
	clientID := os.Getenv("M365_CLIENT_ID")
	clientSecret := os.Getenv("M365_CLIENT_SECRET")
	authMethod := os.Getenv("M365_AUTH_METHOD")
	cloud := os.Getenv("M365_CLOUD")

	if tenantID == "" || clientID == "" || clientSecret == "" || authMethod == "" || cloud == "" {
		// Fall back to implicit configuration from environment
		return ProviderConfig()
	}

	return ProviderConfigWithExplicitValues(tenantID, clientID, clientSecret, authMethod, cloud)
}

// ConfigWithProvider prefixes any terraform configuration with the provider block
func ConfigWithProvider(config string) string {
	return ProviderConfigForCurrentEnv() + config
}
