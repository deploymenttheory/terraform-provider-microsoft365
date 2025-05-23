// File: internal/resources/common/client_setup.go

package common

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
	msgraphsdk "github.com/microsoftgraph/msgraph-sdk-go"
)

// SetGraphStableClientForResource is a helper function to retrieve and validate the Graph Stable client for resources.
func SetGraphStableClientForResource(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse, resourceName string) *msgraphsdk.GraphServiceClient {
	return getClient(ctx, req.ProviderData, resp, resourceName, func(clients *client.GraphClients) *msgraphsdk.GraphServiceClient {
		return clients.StableClient
	})
}

// SetGraphStableClientForDataSource is a helper function to retrieve and validate the Graph Stable client for data sources.
func SetGraphStableClientForDataSource(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse, dataSourceName string) *msgraphsdk.GraphServiceClient {
	return getClient(ctx, req.ProviderData, resp, dataSourceName, func(clients *client.GraphClients) *msgraphsdk.GraphServiceClient {
		return clients.StableClient
	})
}

// SetGraphBetaClientForResource is a helper function to retrieve and validate the Graph Beta client for resources.
func SetGraphBetaClientForResource(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse, resourceName string) *msgraphbetasdk.GraphServiceClient {
	return getClient(ctx, req.ProviderData, resp, resourceName, func(clients *client.GraphClients) *msgraphbetasdk.GraphServiceClient {
		return clients.BetaClient
	})
}

// SetGraphBetaClientForDataSource is a helper function to retrieve and validate the Graph Beta client for data sources.
func SetGraphBetaClientForDataSource(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse, dataSourceName string) *msgraphbetasdk.GraphServiceClient {
	return getClient(ctx, req.ProviderData, resp, dataSourceName, func(clients *client.GraphClients) *msgraphbetasdk.GraphServiceClient {
		return clients.BetaClient
	})
}

// getClient is a helper function to retrieve and validate the appropriate Graph client from provider data.
func getClient[T any, R any](ctx context.Context, providerData any, resp R, name string, getClientFunc func(*client.GraphClients) *T) *T {
	tflog.Debug(ctx, fmt.Sprintf("Configuring %s", name))

	if providerData == nil {
		tflog.Warn(ctx, fmt.Sprintf("Provider data is nil, skipping %s configuration", name))
		return nil
	}

	clients, ok := providerData.(*client.GraphClients)
	if !ok {
		tflog.Error(ctx, "Unexpected Provider Data Type", map[string]interface{}{
			"expected": "*client.GraphClients",
			"actual":   fmt.Sprintf("%T", providerData),
		})

		if respWithDiagnostics, ok := any(resp).(interface{ AddError(string, string) }); ok {
			respWithDiagnostics.AddError(
				"Unexpected Provider Data Type",
				fmt.Sprintf("Expected *client.GraphClients, got: %T. Please report this issue to the provider developers.", providerData),
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
