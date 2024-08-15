package common

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
	msgraphsdk "github.com/microsoftgraph/msgraph-sdk-go"
)

// SetGraphStableClient is a helper function to retrieve and validate the Graph Stable client from provider data.
func SetGraphStableClient(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse, resourceName string) *msgraphsdk.GraphServiceClient {
	return getClient(ctx, req, resp, resourceName, func(clients *client.GraphClients) *msgraphsdk.GraphServiceClient {
		return clients.StableClient
	})
}

// SetGraphBetaClient is a helper function to retrieve and validate the Graph Beta client from provider data.
func SetGraphBetaClient(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse, resourceName string) *msgraphbetasdk.GraphServiceClient {
	return getClient(ctx, req, resp, resourceName, func(clients *client.GraphClients) *msgraphbetasdk.GraphServiceClient {
		return clients.BetaClient
	})
}

// getClient is a helper function to retrieve and validate the appropriate Graph client from provider data.
func getClient[T any](ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse, resourceName string, getClientFunc func(*client.GraphClients) *T) *T {
	tflog.Debug(ctx, fmt.Sprintf("Configuring %s Resource", resourceName))

	if req.ProviderData == nil {
		tflog.Warn(ctx, fmt.Sprintf("Provider data is nil, skipping %s resource configuration", resourceName))
		return nil
	}

	clients, ok := req.ProviderData.(*client.GraphClients)
	if !ok {
		tflog.Error(ctx, "Unexpected Provider Data Type", map[string]interface{}{
			"expected": "*client.GraphClients",
			"actual":   fmt.Sprintf("%T", req.ProviderData),
		})
		resp.Diagnostics.AddError(
			"Unexpected Provider Data Type",
			fmt.Sprintf("Expected *client.GraphClients, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return nil
	}

	client := getClientFunc(clients)
	if client == nil {
		tflog.Warn(ctx, fmt.Sprintf("%s is nil, %s resource may not be fully configured", resourceName, resourceName))
		return nil
	}

	tflog.Debug(ctx, fmt.Sprintf("Initialized %s Resource with Graph Client", resourceName))
	return client
}
