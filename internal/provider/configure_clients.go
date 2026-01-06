package provider

import (
	"context"
	"fmt"
	"net/http"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
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
//  3. Converts the provider model to client provider data.
//  4. Configures the Microsoft Graph clients using the client package.
//
// If any errors occur during these steps, appropriate diagnostics are added
// to the response.
func (p *M365Provider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	tflog.Info(ctx, "Configuring Microsoft365 Provider")

	// In unit test mode, use mock clients instead of real sdk ones
	// and pass the mock clients to data sources and resources instead.
	if p.unitTestMode {
		tflog.Info(ctx, "Provider is in unit test mode. Using mock clients.")

		mockClients := client.NewMockGraphClients(http.DefaultClient)
		p.clients = mockClients

		resp.DataSourceData = mockClients
		resp.ResourceData = mockClients
		resp.ActionData = mockClients
		resp.ListResourceData = mockClients
		return
	}

	var config M365ProviderModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		tflog.Error(ctx, "Error getting provider configuration", map[string]any{
			"diagnostics": resp.Diagnostics.ErrorsCount(),
		})
		return
	}

	data, diags := setProviderConfiguration(ctx, config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		tflog.Error(ctx, "Error populating provider data", map[string]any{
			"diagnostics": resp.Diagnostics.ErrorsCount(),
		})
		return
	}

	clientData := convertToClientProviderData(ctx, &data)

	graphClientInterface := client.NewGraphClients(ctx, clientData, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		tflog.Error(ctx, "Error configuring and building Microsoft Graph clients", map[string]any{
			"diagnostics": resp.Diagnostics.ErrorsCount(),
		})
		return
	}

	p.clients = graphClientInterface

	resp.DataSourceData = graphClientInterface
	resp.ResourceData = graphClientInterface
	resp.ActionData = graphClientInterface
	resp.ListResourceData = graphClientInterface

	tflog.Debug(ctx, "Provider configuration completed", map[string]any{
		"graph_client_set":      p.clients.GetKiotaGraphV1Client() != nil,
		"graph_beta_client_set": p.clients.GetKiotaGraphBetaClient() != nil,
		"config":                fmt.Sprintf("%+v", config),
	})
}

// convertToClientProviderData converts the provider model to client provider data
func convertToClientProviderData(ctx context.Context, data *M365ProviderModel) *client.ProviderData {
	var clientData client.ProviderData

	clientData.Cloud = data.Cloud.ValueString()
	clientData.TenantID = data.TenantID.ValueString()
	clientData.AuthMethod = data.AuthMethod.ValueString()
	clientData.TelemetryOptout = data.TelemetryOptout.ValueBool()
	clientData.DebugMode = data.DebugMode.ValueBool()

	var entraIDOptions EntraIDOptionsModel
	data.EntraIDOptions.As(ctx, &entraIDOptions, basetypes.ObjectAsOptions{})

	oidcRequestURL := entraIDOptions.OIDCRequestURL.ValueString()
	oidcRequestToken := entraIDOptions.OIDCRequestToken.ValueString()

	tflog.Info(ctx, "convertToClientProviderData OIDC values", map[string]any{
		"oidc_request_url":       oidcRequestURL,
		"oidc_request_token_set": oidcRequestToken != "",
	})

	clientData.EntraIDOptions = &client.EntraIDOptions{
		ClientID:                   entraIDOptions.ClientID.ValueString(),
		ClientSecret:               entraIDOptions.ClientSecret.ValueString(),
		ClientCertificate:          entraIDOptions.ClientCertificate.ValueString(),
		ClientCertificatePassword:  entraIDOptions.ClientCertificatePassword.ValueString(),
		Username:                   entraIDOptions.Username.ValueString(),
		RedirectUrl:                entraIDOptions.RedirectUrl.ValueString(),
		FederatedTokenFilePath:     entraIDOptions.FederatedTokenFilePath.ValueString(),
		ManagedIdentityClientID:    entraIDOptions.ManagedIdentityID.ValueString(),
		ManagedIdentityResourceID:  "", // Not in the model
		OIDCTokenFilePath:          entraIDOptions.OIDCTokenFilePath.ValueString(),
		OIDCToken:                  "", // Not in the model
		OIDCRequestToken:           oidcRequestToken,
		OIDCRequestURL:             oidcRequestURL,
		DisableInstanceDiscovery:   entraIDOptions.DisableInstanceDiscovery.ValueBool(),
		SendCertificateChain:       entraIDOptions.SendCertificateChain.ValueBool(),
		AdditionallyAllowedTenants: getAdditionallyAllowedTenants(entraIDOptions.AdditionallyAllowedTenants),
	}

	var clientOptionsModel ClientOptionsModel
	data.ClientOptions.As(ctx, &clientOptionsModel, basetypes.ObjectAsOptions{})

	clientData.ClientOptions = &client.ClientOptions{
		EnableRetry:             clientOptionsModel.EnableRetry.ValueBool(),
		MaxRetries:              clientOptionsModel.MaxRetries.ValueInt64(),
		RetryDelaySeconds:       clientOptionsModel.RetryDelaySeconds.ValueInt64(),
		EnableRedirect:          clientOptionsModel.EnableRedirect.ValueBool(),
		MaxRedirects:            clientOptionsModel.MaxRedirects.ValueInt64(),
		EnableCompression:       clientOptionsModel.EnableCompression.ValueBool(),
		CustomUserAgent:         clientOptionsModel.CustomUserAgent.ValueString(),
		EnableHeadersInspection: clientOptionsModel.EnableHeadersInspection.ValueBool(),
		TimeoutSeconds:          clientOptionsModel.TimeoutSeconds.ValueInt64(),
		UseProxy:                clientOptionsModel.UseProxy.ValueBool(),
		ProxyURL:                clientOptionsModel.ProxyURL.ValueString(),
		ProxyUsername:           clientOptionsModel.ProxyUsername.ValueString(),
		ProxyPassword:           clientOptionsModel.ProxyPassword.ValueString(),
		EnableChaos:             clientOptionsModel.EnableChaos.ValueBool(),
		ChaosPercentage:         clientOptionsModel.ChaosPercentage.ValueInt64(),
		ChaosStatusCode:         clientOptionsModel.ChaosStatusCode.ValueInt64(),
		ChaosStatusMessage:      clientOptionsModel.ChaosStatusMessage.ValueString(),
	}

	return &clientData
}

// Helper function to convert types.List to []string for AdditionallyAllowedTenants
func getAdditionallyAllowedTenants(tenants types.List) []string {
	var result []string
	for _, tenant := range tenants.Elements() {
		if strVal, ok := tenant.(types.String); ok {
			result = append(result, strVal.ValueString())
		}
	}
	return result
}
