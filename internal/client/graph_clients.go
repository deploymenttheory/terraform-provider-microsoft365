package client

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
	msgraphsdk "github.com/microsoftgraph/msgraph-sdk-go"
)

// GraphClients encapsulates both the stable and beta GraphServiceClients
// provided by the Microsoft Graph SDKs. These clients are used to interact
// with the Microsoft Graph API and its beta endpoints, respectively.
//
// The stable client (V1Client) is used for making API calls to the
// stable Microsoft Graph endpoints, which are generally considered
// production-ready and have a higher level of reliability and support.
// The v1.0 endpoint of Microsoft Graph provides a stable and reliable API
// that is fully supported by Microsoft, ensuring that applications built
// on this endpoint have a solid foundation and offer the best possible
// user experience.
//
// The beta client (BetaClient) is used for making API calls to the
// beta Microsoft Graph endpoints, which allow developers to test and
// experiment with newest features in the graph ecosystem.
//
// Microsoft claim,  that the beta endpoint is not intended
// for use in production environments. However, much of the gui uses graph beta
// e.g with intune, conditional access, etc within a production context. I.e
// microsoft use the beta endpoints consistently like it's a production endpoint.
// Despite the beta label. Conversations with microsoft product teams, have explained
// that the reason for this is as follows:
//
// graph v1.0 has a very strict breaking change policy, allowing for one
// breaking change per year. This is to ensure that the api is stable and reliable.
// However, the beta endpoint is not subject to this policy, and allows for more
// frequent breaking changes. This is to allow for new features to be added to the
// graph api without having to wait for a year by microsoft development teams.
//
// Additionally, it's become the norm that for many api endpoints, they never get
// a v1.0 endpoint, ever. Intune is a good example of this, where endpoints for
// are still in 'beta', despite being in production for many years. Microsoft
// have also stated off the record that in many cases they will support the beta
// api like they do the v1.0 api.
//
// Conseqently, depsite the offical line that developers should use the v1.0
// it's not that clear cut.
//
// For these reasons, this provider shall use what the gui uses for a given
// piece of functionality. Typically mapped to whatever graph x-ray
// (https://graphxray.merill.net/) observes during api calls.
//
// Fields:
//
//	StableClient (*msgraphsdk.GraphServiceClient): The client for interacting
//	  with the stable Microsoft Graph API, providing access to well-supported
//	  and reliable endpoints suitable for production use.
//
//	BetaClient (*msgraphbetasdk.GraphServiceClient): The client for interacting
//	  with the beta Microsoft Graph API, providing access to new and experimental
//	  features that are subject to change and should be used with caution in
//	  production environments.
//
// Usage:
// The GraphClients struct is intended to be instantiated and configured by
// the provider during initialization, and then passed to the resources that
// need to interact with the Microsoft Graph API. This separation of stable
// and beta clients allows resources to choose the appropriate client based
// on the API features they require.

// GraphClientInterface defines the interface for GraphClients
type GraphClientInterface interface {
	GetV1Client() *msgraphsdk.GraphServiceClient
	GetBetaClient() *msgraphbetasdk.GraphServiceClient
}

type GraphClients struct {
	V1Client   *msgraphsdk.GraphServiceClient
	BetaClient *msgraphbetasdk.GraphServiceClient
}

// GetStableClient returns the stable client
func (g *GraphClients) GetV1Client() *msgraphsdk.GraphServiceClient {
	return g.V1Client
}

// GetBetaClient returns the beta client
func (g *GraphClients) GetBetaClient() *msgraphbetasdk.GraphServiceClient {
	return g.BetaClient
}

// SetGraphStableClientForResource is a helper function to retrieve and validate the Graph V1.0 client for resources.
func SetGraphStableClientForResource(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse, resourceName string) *msgraphsdk.GraphServiceClient {
	return getClient(ctx, req.ProviderData, resp, resourceName, func(clients GraphClientInterface) *msgraphsdk.GraphServiceClient {
		return clients.GetV1Client()
	})
}

// SetGraphStableClientForDataSource is a helper function to retrieve and validate the Graph V1.0 client for data sources.
func SetGraphStableClientForDataSource(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse, dataSourceName string) *msgraphsdk.GraphServiceClient {
	return getClient(ctx, req.ProviderData, resp, dataSourceName, func(clients GraphClientInterface) *msgraphsdk.GraphServiceClient {
		return clients.GetV1Client()
	})
}

// SetGraphBetaClientForResource is a helper function to retrieve and validate the Graph Beta client for resources.
func SetGraphBetaClientForResource(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse, resourceName string) *msgraphbetasdk.GraphServiceClient {
	return getClient(ctx, req.ProviderData, resp, resourceName, func(clients GraphClientInterface) *msgraphbetasdk.GraphServiceClient {
		return clients.GetBetaClient()
	})
}

// SetGraphBetaClientForDataSource is a helper function to retrieve and validate the Graph Beta client for data sources.
func SetGraphBetaClientForDataSource(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse, dataSourceName string) *msgraphbetasdk.GraphServiceClient {
	return getClient(ctx, req.ProviderData, resp, dataSourceName, func(clients GraphClientInterface) *msgraphbetasdk.GraphServiceClient {
		return clients.GetBetaClient()
	})
}

// getClient is a helper function to retrieve and validate the appropriate Graph client from provider data.
func getClient[T any, R any](ctx context.Context, providerData any, resp R, name string, getClientFunc func(GraphClientInterface) *T) *T {
	tflog.Debug(ctx, fmt.Sprintf("Configuring %s", name))

	if providerData == nil {
		tflog.Warn(ctx, fmt.Sprintf("Provider data is nil, skipping %s configuration", name))
		return nil
	}

	clients, ok := providerData.(GraphClientInterface)
	if !ok {
		tflog.Error(ctx, "Unexpected Provider Data Type", map[string]interface{}{
			"expected": "GraphClientInterface",
			"actual":   fmt.Sprintf("%T", providerData),
		})

		if respWithDiagnostics, ok := any(resp).(interface{ AddError(string, string) }); ok {
			respWithDiagnostics.AddError(
				"Unexpected Provider Data Type",
				fmt.Sprintf("Expected GraphClientInterface, got: %T. Please report this issue to the provider developers.", providerData),
			)
		}
		return nil
	}

	client := getClientFunc(clients)
	if client == nil {
		tflog.Warn(ctx, fmt.Sprintf("%s client is nil, %s may not be fully configured", name, name))
		return nil
	}

	tflog.Debug(ctx, fmt.Sprintf("Initialized %s with Graph Client", name))
	return client
}
