package mapper

import (
	"fmt"
	"os"
	"strings"
)

// Exporter is an interface for exporting environment variables.
type Exporter interface {
	Export(key, value string) error
	HasCapability() bool
}

// Mapper maps service credentials to environment variables.
type Mapper struct {
	exporter Exporter
}

// New creates a new credential mapper with the given exporter.
func New(exporter Exporter) *Mapper {
	return &Mapper{
		exporter: exporter,
	}
}

// Map maps the credentials for the given service to environment variables.
func (m *Mapper) Map(service string) error {
	cred, ok := GetServiceCredential(service)
	if !ok {
		return fmt.Errorf("unknown service: %s (supported services: %v)",
			service, strings.Join(SupportedServices(), ", "))
	}

	clientID := os.Getenv(cred.ClientIDVar)
	clientSecret := os.Getenv(cred.ClientSecretVar)

	if !m.exporter.HasCapability() {
		// Fallback to stdout for local testing
		fmt.Printf("export %s=%s\n", EnvM365ClientID, clientID)
		fmt.Printf("export %s=%s\n", EnvM365ClientSecret, clientSecret)

		if clientID != "" && clientSecret != "" {
			fmt.Printf("✅ Credentials configured for %s\n", service)
		} else {
			fmt.Printf("⚠️  No credentials found for %s\n", service)
		}
		return nil
	}

	if err := m.exporter.Export(EnvM365ClientID, clientID); err != nil {
		return fmt.Errorf("failed to export client ID: %w", err)
	}

	if err := m.exporter.Export(EnvM365ClientSecret, clientSecret); err != nil {
		return fmt.Errorf("failed to export client secret: %w", err)
	}

	if clientID != "" && clientSecret != "" {
		fmt.Printf("✅ Credentials configured for %s\n", service)
	} else {
		fmt.Printf("⚠️  No credentials found for %s\n", service)
	}

	return nil
}

// EnvExporter exports variables to GitHub Actions environment.
type EnvExporter struct {
	githubEnvPath string
}

// NewEnvExporter creates a new GitHub Actions environment exporter.
func NewEnvExporter() *EnvExporter {
	return &EnvExporter{
		githubEnvPath: os.Getenv("GITHUB_ENV"),
	}
}

// HasCapability returns true if the exporter can export to GITHUB_ENV.
func (e *EnvExporter) HasCapability() bool {
	return e.githubEnvPath != ""
}

// Export writes a key-value pair to the GITHUB_ENV file.
func (e *EnvExporter) Export(key, value string) error {
	if !e.HasCapability() {
		return fmt.Errorf("GITHUB_ENV not set")
	}

	f, err := os.OpenFile(e.githubEnvPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open GITHUB_ENV file: %w", err)
	}
	defer f.Close()

	if _, err := fmt.Fprintf(f, "%s=%s\n", key, value); err != nil {
		return fmt.Errorf("failed to write to GITHUB_ENV: %w", err)
	}

	return nil
}
