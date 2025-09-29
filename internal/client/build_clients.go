package client

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
	msgraphsdk "github.com/microsoftgraph/msgraph-sdk-go"
	"github.com/microsoftgraph/msgraph-sdk-go-core/authentication"
)

// NewGraphClients sets up the Microsoft Graph clients with the given configuration.
// It processes the provider data, sets up authentication, and initializes the Microsoft
// Graph clients (stable and beta).
//
// The function supports various authentication methods, proxy settings, and
// national cloud deployments. It performs the following main steps:
//  1. Determines cloud-specific constants and endpoints.
//  2. Configures the Entra ID client options.
//  3. Obtains credentials based on the specified authentication method.
//  4. Creates and configures the Microsoft Graph clients (stable and beta).
//
// If any errors occur during these steps, appropriate diagnostics are added
// to the diagnostics collection.
func NewGraphClients(ctx context.Context, data *ProviderData, diags *diag.Diagnostics) GraphClientInterface {
	tflog.Info(ctx, "Configuring Microsoft Graph Clients")

	authorityURL, apiScope, graphServiceRoot, graphBetaServiceRoot, err := SetCloudConstants(data.Cloud)
	if err != nil {
		diags.AddError(
			"Invalid Microsoft Cloud Type",
			fmt.Sprintf("An error occurred while attempting to get cloud constants for cloud type '%s'. "+
				"Please ensure the cloud type is valid. Detailed error: %s", data.Cloud, err.Error()),
		)
		return nil
	}

	clientOptions, err := ConfigureEntraIDClientOptions(ctx, data, authorityURL)
	if err != nil {
		diags.AddError(
			"Unable to configure client options",
			fmt.Sprintf("An error occurred while attempting to configure client options. Detailed error: %s", err.Error()),
		)
		return nil
	}

	cred, err := ObtainCredential(ctx, data, clientOptions)
	if err != nil {
		diags.AddError(
			"Unable to create credentials",
			fmt.Sprintf("An error occurred while attempting to create the credentials: %s", err.Error()),
		)
		return nil
	}

	authProvider, err := authentication.NewAzureIdentityAuthenticationProviderWithScopes(
		cred,
		[]string{apiScope},
	)
	if err != nil {
		diags.AddError(
			"Unable to create authentication provider",
			fmt.Sprintf("An error occurred while attempting to create the authentication provider using the provided credentials. "+
				"This may be due to misconfigured client options, incorrect credentials, or issues with the underlying authentication library. "+
				"Please verify your client options and credentials configuration. Detailed error: %s", err.Error()),
		)
		return nil
	}

	httpClient, err := ConfigureGraphClientOptions(ctx, data)
	if err != nil {
		diags.AddError(
			"Unable to configure Graph client options",
			fmt.Sprintf("An error occurred while attempting to configure the Microsoft Graph client options. Detailed error: %s", err.Error()),
		)
		return nil
	}

	graphV1Adapter, err := msgraphsdk.NewGraphRequestAdapterWithParseNodeFactoryAndSerializationWriterFactoryAndHttpClient(
		authProvider,
		nil,
		nil,
		httpClient,
	)
	if err != nil {
		diags.AddError(
			"Failed to create Microsoft Graph Stable SDK Adapter",
			fmt.Sprintf("An error occurred while attempting to create the Microsoft Graph Stable SDK adapter. This might be due to issues with the authentication provider, HTTP client setup, or the SDK's internal components. Detailed error: %s", err.Error()),
		)
		return nil
	}

	graphBetaAdapter, err := msgraphbetasdk.NewGraphRequestAdapterWithParseNodeFactoryAndSerializationWriterFactoryAndHttpClient(
		authProvider,
		nil,
		nil,
		httpClient,
	)
	if err != nil {
		diags.AddError(
			"Failed to create Microsoft Graph Beta SDK Adapter",
			fmt.Sprintf("An error occurred while attempting to create the Microsoft Graph Beta SDK adapter. This might be due to issues with the authentication provider, HTTP client setup, or the SDK's internal components. Detailed error: %s", err.Error()),
		)
		return nil
	}

	graphV1Adapter.SetBaseUrl(graphServiceRoot)
	graphBetaAdapter.SetBaseUrl(graphBetaServiceRoot)

	// Create HTTP clients for raw JSON calls
	graphV1Client := NewAuthenticatedHTTPClient(httpClient, cred, apiScope, graphServiceRoot)
	graphBetaClient := NewAuthenticatedHTTPClient(httpClient, cred, apiScope, graphBetaServiceRoot)

	clients := &GraphClients{
		KiotaGraphV1Client:   msgraphsdk.NewGraphServiceClient(graphV1Adapter),
		KiotaGraphBetaClient: msgraphbetasdk.NewGraphServiceClient(graphBetaAdapter),
		GraphV1Client:        graphV1Client,
		GraphBetaClient:      graphBetaClient,
	}

	tflog.Debug(ctx, "Graph clients configuration completed", map[string]any{
		"graph_client_set":           clients.GetKiotaGraphV1Client() != nil,
		"graph_beta_client_set":      clients.GetKiotaGraphBetaClient() != nil,
		"graph_http_client_set":      clients.GetGraphV1Client() != nil,
		"graph_beta_http_client_set": clients.GetGraphBetaClient() != nil,
	})

	return clients
}
