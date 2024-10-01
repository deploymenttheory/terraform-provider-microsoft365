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
// It configures the HTTP client with the default middlewares and optionally adds a chaos handler.
// If useProxy is true, it configures the HTTP client with the provided proxy URL.
func configureGraphClientOptions(ctx context.Context, useProxy bool, proxyURL string, enableChaos bool) (*http.Client, error) {
	tflog.Debug(ctx, "Configuring Graph client options")

	defaultClientOptions := msgraphsdk.GetDefaultClientOptions()
	defaultMiddleware := msgraphgocore.GetDefaultMiddlewaresWithOptions(&defaultClientOptions)

	if enableChaos {
		chaosHandler := khttp.NewChaosHandler()
		defaultMiddleware = append(defaultMiddleware, chaosHandler)
	}

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
