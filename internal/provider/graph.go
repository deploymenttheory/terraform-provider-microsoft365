package provider

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	khttp "github.com/microsoft/kiota-http-go"
	msgraphsdk "github.com/microsoftgraph/msgraph-sdk-go"
	msgraphgocore "github.com/microsoftgraph/msgraph-sdk-go-core"
)

// configureGraphClientOptions sets up the client options for the Microsoft Graph SDK.
func configureGraphClientOptions(ctx context.Context, useProxy bool, proxyURL string, enableChaos bool) (*http.Client, error) {
	tflog.Debug(ctx, "Configuring Graph client options")

	// Get default client options
	defaultClientOptions := msgraphsdk.GetDefaultClientOptions()
	defaultMiddleware := msgraphgocore.GetDefaultMiddlewaresWithOptions(&defaultClientOptions)

	// Add chaos handler if enabled
	if enableChaos {
		chaosHandler := khttp.NewChaosHandler()
		defaultMiddleware = append(defaultMiddleware, chaosHandler)
	}

	// Configure HTTP client with or without proxy settings
	var httpClient *http.Client
	var err error
	if useProxy && proxyURL != "" {
		httpClient, err = khttp.GetClientWithProxySettings(proxyURL, defaultMiddleware...)
		if err != nil {
			return nil, fmt.Errorf("unable to create HTTP client with proxy settings: %w", err)
		}
	} else {
		httpClient = khttp.GetDefaultClient(defaultMiddleware...)
	}

	tflog.Debug(ctx, "Configured Graph client options")
	return httpClient, nil
}
