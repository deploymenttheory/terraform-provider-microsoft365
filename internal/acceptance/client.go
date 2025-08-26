package acceptance

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

// TestGraphClient creates a Graph client for acceptance tests using environment variables
func TestGraphClient() (*msgraphbetasdk.GraphServiceClient, error) {
	// Create provider data similar to how the provider does it
	providerData := &client.ProviderData{
		AuthMethod: os.Getenv("M365_AUTH_METHOD"),
		Cloud:      os.Getenv("M365_CLOUD"),
		TenantID:   os.Getenv("M365_TENANT_ID"),
		EntraIDOptions: &client.EntraIDOptions{
			ClientID:     os.Getenv("M365_CLIENT_ID"),
			ClientSecret: os.Getenv("M365_CLIENT_SECRET"),
		},
		// Initialize ClientOptions with reasonable defaults to prevent nil pointer dereference
		ClientOptions: &client.ClientOptions{
			EnableRetry:             true,
			MaxRetries:              3,
			RetryDelaySeconds:       1,
			EnableRedirect:          true,
			MaxRedirects:            5,
			EnableCompression:       true,
			CustomUserAgent:         "terraform-provider-microsoft365/acceptance-tests",
			EnableHeadersInspection: false,
			TimeoutSeconds:          60,
			UseProxy:                false,
			ProxyURL:                "",
			ProxyUsername:           "",
			ProxyPassword:           "",
			EnableChaos:             false,
			ChaosPercentage:         0,
			ChaosStatusCode:         0,
			ChaosStatusMessage:      "",
		},
	}

	// Use the same client building logic as the provider
	ctx := context.Background()
	var diags diag.Diagnostics

	graphClients := client.NewGraphClients(ctx, providerData, &diags)
	if diags.HasError() {
		// Convert diagnostics to error message
		var errMsg strings.Builder
		for _, d := range diags.Errors() {
			errMsg.WriteString(d.Summary())
			errMsg.WriteString(": ")
			errMsg.WriteString(d.Detail())
			errMsg.WriteString("; ")
		}
		return nil, fmt.Errorf("failed to build Graph clients: %s", errMsg.String())
	}

	if graphClients == nil {
		return nil, fmt.Errorf("graph clients is nil")
	}

	betaClient := graphClients.GetKiotaGraphBetaClient()
	if betaClient == nil {
		return nil, fmt.Errorf("beta client is nil")
	}

	return betaClient, nil
}

// TestHTTPClient creates an authenticated HTTP client for acceptance tests using environment variables
func TestHTTPClient() (*client.AuthenticatedHTTPClient, error) {
	// Create provider data similar to how the provider does it
	providerData := &client.ProviderData{
		AuthMethod: os.Getenv("M365_AUTH_METHOD"),
		Cloud:      os.Getenv("M365_CLOUD"),
		TenantID:   os.Getenv("M365_TENANT_ID"),
		EntraIDOptions: &client.EntraIDOptions{
			ClientID:     os.Getenv("M365_CLIENT_ID"),
			ClientSecret: os.Getenv("M365_CLIENT_SECRET"),
		},
		// Initialize ClientOptions with reasonable defaults to prevent nil pointer dereference
		ClientOptions: &client.ClientOptions{
			EnableRetry:             true,
			MaxRetries:              3,
			RetryDelaySeconds:       1,
			EnableRedirect:          true,
			MaxRedirects:            5,
			EnableCompression:       true,
			CustomUserAgent:         "terraform-provider-microsoft365/acceptance-tests",
			EnableHeadersInspection: false,
			TimeoutSeconds:          60,
			UseProxy:                false,
			ProxyURL:                "",
			ProxyUsername:           "",
			ProxyPassword:           "",
			EnableChaos:             false,
			ChaosPercentage:         0,
			ChaosStatusCode:         0,
			ChaosStatusMessage:      "",
		},
	}

	// Use the same client building logic as the provider
	ctx := context.Background()
	var diags diag.Diagnostics

	graphClients := client.NewGraphClients(ctx, providerData, &diags)
	if diags.HasError() {
		// Convert diagnostics to error message
		var errMsg strings.Builder
		for _, d := range diags.Errors() {
			errMsg.WriteString(d.Summary())
			errMsg.WriteString(": ")
			errMsg.WriteString(d.Detail())
			errMsg.WriteString("; ")
		}
		return nil, fmt.Errorf("failed to build Graph clients: %s", errMsg.String())
	}

	if graphClients == nil {
		return nil, fmt.Errorf("graph clients is nil")
	}

	betaHTTPClient := graphClients.GetGraphBetaClient()
	if betaHTTPClient == nil {
		return nil, fmt.Errorf("beta HTTP client is nil")
	}

	return betaHTTPClient, nil
}
