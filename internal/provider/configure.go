package provider

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
	msgraphsdk "github.com/microsoftgraph/msgraph-sdk-go"
	"github.com/microsoftgraph/msgraph-sdk-go-core/authentication"
)

// Configure sets up the Microsoft365 provider with the given configuration.
// It processes the provider schema, retrieves values from the configuration or
// environment variables, sets up authentication, and initializes the Microsoft
// Graph clients.
//
// The function supports various authentication methods, proxy settings, and
// national cloud deployments. It performs the following main steps:
//  1. Extracts and validates the configuration data.
//  2. Sets up logging and context with relevant fields.
//  3. Determines cloud-specific constants and endpoints.
//  4. Configures the Entra ID client options.
//  5. Obtains credentials based on the specified authentication method.
//  6. Creates and configures the Microsoft Graph clients (stable and beta).
//
// If any errors occur during these steps, appropriate diagnostics are added
// to the response.
func (p *M365Provider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	tflog.Info(ctx, "Configuring Microsoft365 Provider")

	if p.testMode {
		tflog.Warn(ctx, "Provider is in test mode. Skipping configuration.")
		return
	}

	var config M365ProviderModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		tflog.Error(ctx, "Error getting provider configuration", map[string]interface{}{
			"diagnostics": resp.Diagnostics.ErrorsCount(),
		})
		return
	}

	data, diags := populateProviderData(ctx, config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		tflog.Error(ctx, "Error populating provider data", map[string]interface{}{
			"diagnostics": resp.Diagnostics.ErrorsCount(),
		})
		return
	}

	authorityURL, apiScope, graphServiceRoot, graphBetaServiceRoot, err := setCloudConstants(data.Cloud.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid Microsoft Cloud Type",
			fmt.Sprintf("An error occurred while attempting to get cloud constants for cloud type '%s'. "+
				"Please ensure the cloud type is valid. Detailed error: %s", data.Cloud.ValueString(), err.Error()),
		)
		return
	}

	clientOptions, err := configureEntraIDClientOptions(ctx, &data, authorityURL)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to configure client options",
			fmt.Sprintf("An error occurred while attempting to configure client options. Detailed error: %s", err.Error()),
		)
		return
	}

	cred, err := obtainCredential(ctx, &data, clientOptions)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to create credentials",
			fmt.Sprintf("An error occurred while attempting to create the credentials: %s", err.Error()),
		)
		return
	}

	authProvider, err := authentication.NewAzureIdentityAuthenticationProviderWithScopes(
		cred,
		[]string{apiScope},
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to create authentication provider",
			fmt.Sprintf("An error occurred while attempting to create the authentication provider using the provided credentials. "+
				"This may be due to misconfigured client options, incorrect credentials, or issues with the underlying authentication library. "+
				"Please verify your client options and credentials configuration. Detailed error: %s", err.Error()),
		)
		return
	}

	httpClient, err := configureGraphClientOptions(ctx, &data)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to configure Graph client options",
			fmt.Sprintf("An error occurred while attempting to configure the Microsoft Graph client options. Detailed error: %s", err.Error()),
		)
		return
	}

	stableAdapter, err := msgraphsdk.NewGraphRequestAdapterWithParseNodeFactoryAndSerializationWriterFactoryAndHttpClient(
		authProvider,
		nil,
		nil,
		httpClient,
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to create Microsoft Graph Stable SDK Adapter",
			fmt.Sprintf("An error occurred while attempting to create the Microsoft Graph Stable SDK adapter. This might be due to issues with the authentication provider, HTTP client setup, or the SDK's internal components. Detailed error: %s", err.Error()),
		)
		return
	}

	betaAdapter, err := msgraphbetasdk.NewGraphRequestAdapterWithParseNodeFactoryAndSerializationWriterFactoryAndHttpClient(
		authProvider,
		nil,
		nil,
		httpClient,
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to create Microsoft Graph Beta SDK Adapter",
			fmt.Sprintf("An error occurred while attempting to create the Microsoft Graph Beta SDK adapter. This might be due to issues with the authentication provider, HTTP client setup, or the SDK's internal components. Detailed error: %s", err.Error()),
		)
		return
	}

	stableAdapter.SetBaseUrl(graphServiceRoot)
	betaAdapter.SetBaseUrl(graphBetaServiceRoot)

	clients := &client.GraphClients{
		StableClient: msgraphsdk.NewGraphServiceClient(stableAdapter),
		BetaClient:   msgraphbetasdk.NewGraphServiceClient(betaAdapter),
	}

	p.clients = clients

	resp.DataSourceData = clients
	resp.ResourceData = clients

	tflog.Debug(ctx, "Provider configuration completed", map[string]interface{}{
		"graph_client_set":      p.clients.StableClient != nil,
		"graph_beta_client_set": p.clients.BetaClient != nil,
		"config":                fmt.Sprintf("%+v", config),
	})
}
