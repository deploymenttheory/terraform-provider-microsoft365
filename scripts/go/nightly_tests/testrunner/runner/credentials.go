package runner

import (
	"os"
	"strings"
)

// Required environment variables for authentication.
const (
	EnvClientID     = "M365_CLIENT_ID"
	EnvClientSecret = "M365_CLIENT_SECRET"
	EnvTenantID     = "M365_TENANT_ID"
	EnvAuthMethod   = "M365_AUTH_METHOD"
	EnvCloud        = "M365_CLOUD"
)

// Credentials represents authentication credentials.
type Credentials struct {
	ClientID     string
	ClientSecret string
	TenantID     string
	AuthMethod   string
	Cloud        string
}

// IsValid returns true if the credentials are complete and non-empty.
func (c *Credentials) IsValid() bool {
	return strings.TrimSpace(c.ClientID) != "" &&
		strings.TrimSpace(c.ClientSecret) != ""
}

// CredentialProvider is an interface for retrieving credentials.
type CredentialProvider interface {
	GetCredentials() (*Credentials, error)
}

// EnvCredentialProvider retrieves credentials from environment variables.
type EnvCredentialProvider struct{}

// NewEnvCredentialProvider creates a new environment-based credential provider.
func NewEnvCredentialProvider() *EnvCredentialProvider {
	return &EnvCredentialProvider{}
}

// GetCredentials retrieves credentials from environment variables.
func (p *EnvCredentialProvider) GetCredentials() (*Credentials, error) {
	return &Credentials{
		ClientID:     os.Getenv(EnvClientID),
		ClientSecret: os.Getenv(EnvClientSecret),
		TenantID:     os.Getenv(EnvTenantID),
		AuthMethod:   os.Getenv(EnvAuthMethod),
		Cloud:        os.Getenv(EnvCloud),
	}, nil
}
