package acceptance

import (
	"fmt"
	"os"
	"strings"
)

// ProviderConfigBuilder helps build provider configuration blocks
type ProviderConfigBuilder struct {
	cloud            string
	authMethod       string
	tenantID         string
	clientID         string
	clientSecret     string
	certificate      string
	certificatePass  string
	redirectURL      string
	debugMode        bool
	telemetryOptout  bool
	useProxy         bool
	proxyURL         string
	envVarOverrides  bool
}

// NewProviderConfigBuilder creates a new provider configuration builder
func NewProviderConfigBuilder() *ProviderConfigBuilder {
	return &ProviderConfigBuilder{
		cloud:           "public",
		authMethod:      "device_code", 
		debugMode:       false,
		telemetryOptout: false,
		useProxy:        false,
		envVarOverrides: false,
	}
}

// WithCloud sets the cloud environment
func (p *ProviderConfigBuilder) WithCloud(cloud string) *ProviderConfigBuilder {
	p.cloud = cloud
	return p
}

// WithAuthMethod sets the authentication method  
func (p *ProviderConfigBuilder) WithAuthMethod(method string) *ProviderConfigBuilder {
	p.authMethod = method
	return p
}

// WithClientSecret sets up client secret authentication
func (p *ProviderConfigBuilder) WithClientSecret(clientID, clientSecret string) *ProviderConfigBuilder {
	p.authMethod = "client_secret"
	p.clientID = clientID
	p.clientSecret = clientSecret
	return p
}

// WithClientCertificate sets up client certificate authentication
func (p *ProviderConfigBuilder) WithClientCertificate(clientID, certPath, certPass string) *ProviderConfigBuilder {
	p.authMethod = "client_certificate"
	p.clientID = clientID
	p.certificate = certPath
	p.certificatePass = certPass
	return p
}

// WithProxy enables proxy configuration
func (p *ProviderConfigBuilder) WithProxy(proxyURL string) *ProviderConfigBuilder {
	p.useProxy = true
	p.proxyURL = proxyURL
	return p
}

// WithDebugMode enables debug mode
func (p *ProviderConfigBuilder) WithDebugMode(enabled bool) *ProviderConfigBuilder {
	p.debugMode = enabled
	return p
}

// WithTelemetryOptout sets telemetry optout preference
func (p *ProviderConfigBuilder) WithTelemetryOptout(optout bool) *ProviderConfigBuilder {
	p.telemetryOptout = optout
	return p
}

// WithEnvironmentVariables enables environment variable precedence
func (p *ProviderConfigBuilder) WithEnvironmentVariables() *ProviderConfigBuilder {
	p.envVarOverrides = true
	return p
}

// Build generates the Terraform provider configuration string
func (p *ProviderConfigBuilder) Build() string {
	var config strings.Builder
	
	config.WriteString("provider \"microsoft365\" {\n")
	
	if p.cloud != "" && p.cloud != "public" {
		config.WriteString(fmt.Sprintf("  cloud = \"%s\"\n", p.cloud))
	}
	
	if p.authMethod != "" {
		config.WriteString(fmt.Sprintf("  auth_method = \"%s\"\n", p.authMethod))
	}
	
	if p.tenantID != "" {
		config.WriteString(fmt.Sprintf("  tenant_id = \"%s\"\n", p.tenantID))
	}
	
	if p.debugMode {
		config.WriteString("  debug_mode = true\n")
	}
	
	if p.telemetryOptout {
		config.WriteString("  telemetry_optout = true\n")
	}
	
	// Add Entra ID options if needed
	if p.clientID != "" || p.clientSecret != "" || p.certificate != "" || p.redirectURL != "" {
		config.WriteString("  entra_id_options = {\n")
		
		if p.clientID != "" {
			config.WriteString(fmt.Sprintf("    client_id = \"%s\"\n", p.clientID))
		}
		if p.clientSecret != "" {
			config.WriteString(fmt.Sprintf("    client_secret = \"%s\"\n", p.clientSecret))
		}
		if p.certificate != "" {
			config.WriteString(fmt.Sprintf("    client_certificate = \"%s\"\n", p.certificate))
		}
		if p.certificatePass != "" {
			config.WriteString(fmt.Sprintf("    client_certificate_password = \"%s\"\n", p.certificatePass))
		}
		if p.redirectURL != "" {
			config.WriteString(fmt.Sprintf("    redirect_url = \"%s\"\n", p.redirectURL))
		}
		
		config.WriteString("  }\n")
	}
	
	// Add client options if needed
	if p.useProxy {
		config.WriteString("  client_options = {\n")
		config.WriteString("    use_proxy = true\n")
		if p.proxyURL != "" {
			config.WriteString(fmt.Sprintf("    proxy_url = \"%s\"\n", p.proxyURL))
		}
		config.WriteString("  }\n")
	}
	
	config.WriteString("}\n")
	
	return config.String()
}

// ProviderConfigWithAuthMethod creates a basic provider config with the specified auth method
func ProviderConfigWithAuthMethod(authMethod string) string {
	return NewProviderConfigBuilder().
		WithAuthMethod(authMethod).
		Build()
}

// ProviderConfigWithCloud creates a basic provider config with the specified cloud
func ProviderConfigWithCloud(cloud string) string {
	return NewProviderConfigBuilder().
		WithCloud(cloud).
		Build()
}

// ProviderConfigForClientSecret creates a provider config for client secret authentication
func ProviderConfigForClientSecret(clientID, clientSecret string) string {
	return NewProviderConfigBuilder().
		WithClientSecret(clientID, clientSecret).
		Build()
}

// ProviderConfigForClientCertificate creates a provider config for client certificate authentication
func ProviderConfigForClientCertificate(clientID, certPath, certPass string) string {
	return NewProviderConfigBuilder().
		WithClientCertificate(clientID, certPath, certPass).
		Build()
}

// RequiredEnvVars returns the list of environment variables required for acceptance tests
func RequiredEnvVars() []string {
	return []string{
		"M365_TENANT_ID",
		"M365_CLIENT_ID",
		"M365_CLIENT_SECRET",
		"M365_AUTH_METHOD",
		"M365_CLOUD",
	}
}

// OptionalEnvVars returns the list of optional environment variables for acceptance tests  
func OptionalEnvVars() []string {
	return []string{
		"M365_REDIRECT_URL",
		"M365_USE_PROXY",
		"M365_PROXY_URL",
		"M365_ENABLE_CHAOS",
		"M365_TELEMETRY_OPTOUT",
		"M365_DEBUG_MODE",
	}
}

// CheckRequiredEnvVars validates that required environment variables are set
func CheckRequiredEnvVars() error {
	var missing []string
	
	for _, envVar := range RequiredEnvVars() {
		if v := os.Getenv(envVar); v == "" {
			missing = append(missing, envVar)
		}
	}
	
	if len(missing) > 0 {
		return fmt.Errorf("required environment variables are missing: %s", strings.Join(missing, ", "))
	}
	
	return nil
}