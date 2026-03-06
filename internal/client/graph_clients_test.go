package client

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
	msgraphsdk "github.com/microsoftgraph/msgraph-sdk-go"
	"github.com/stretchr/testify/assert"
)

// TestGraphClients_Getters validates all getter methods return the correct clients.
func TestGraphClients_Getters(t *testing.T) {
	v1Client := &msgraphsdk.GraphServiceClient{}
	betaClient := &msgraphbetasdk.GraphServiceClient{}
	v1HTTPClient := &AuthenticatedHTTPClient{baseURL: "https://graph.microsoft.com/v1.0"}
	betaHTTPClient := &AuthenticatedHTTPClient{baseURL: "https://graph.microsoft.com/beta"}

	clients := &GraphClients{
		KiotaGraphV1Client:   v1Client,
		KiotaGraphBetaClient: betaClient,
		GraphV1Client:        v1HTTPClient,
		GraphBetaClient:      betaHTTPClient,
	}

	t.Run("GetKiotaGraphV1Client", func(t *testing.T) {
		result := clients.GetKiotaGraphV1Client()
		assert.Equal(t, v1Client, result)
	})

	t.Run("GetKiotaGraphBetaClient", func(t *testing.T) {
		result := clients.GetKiotaGraphBetaClient()
		assert.Equal(t, betaClient, result)
	})

	t.Run("GetGraphV1Client", func(t *testing.T) {
		result := clients.GetGraphV1Client()
		assert.Equal(t, v1HTTPClient, result)
		assert.Equal(t, "https://graph.microsoft.com/v1.0", result.GetBaseURL())
	})

	t.Run("GetGraphBetaClient", func(t *testing.T) {
		result := clients.GetGraphBetaClient()
		assert.Equal(t, betaHTTPClient, result)
		assert.Equal(t, "https://graph.microsoft.com/beta", result.GetBaseURL())
	})
}

// TestSetGraphStableClientForResource validates the SetGraphStableClientForResource helper.
func TestSetGraphStableClientForResource(t *testing.T) {
	ctx := context.Background()

	t.Run("Valid provider data", func(t *testing.T) {
		mockClients := NewMockGraphClients(nil)
		req := resource.ConfigureRequest{
			ProviderData: mockClients,
		}
		resp := &resource.ConfigureResponse{}

		client := SetGraphStableClientForResource(ctx, req, resp, "test_resource")
		assert.NotNil(t, client)
		assert.False(t, resp.Diagnostics.HasError())
	})

	t.Run("Nil provider data", func(t *testing.T) {
		req := resource.ConfigureRequest{
			ProviderData: nil,
		}
		resp := &resource.ConfigureResponse{}

		client := SetGraphStableClientForResource(ctx, req, resp, "test_resource")
		assert.Nil(t, client)
	})

	t.Run("Invalid provider data type", func(t *testing.T) {
		req := resource.ConfigureRequest{
			ProviderData: "invalid-type",
		}
		resp := &resource.ConfigureResponse{}

		client := SetGraphStableClientForResource(ctx, req, resp, "test_resource")
		assert.Nil(t, client)
	})
}

// TestSetGraphStableClientForDataSource validates the SetGraphStableClientForDataSource helper.
func TestSetGraphStableClientForDataSource(t *testing.T) {
	ctx := context.Background()

	t.Run("Valid provider data", func(t *testing.T) {
		mockClients := NewMockGraphClients(nil)
		req := datasource.ConfigureRequest{
			ProviderData: mockClients,
		}
		resp := &datasource.ConfigureResponse{}

		client := SetGraphStableClientForDataSource(ctx, req, resp, "test_datasource")
		assert.NotNil(t, client)
		assert.False(t, resp.Diagnostics.HasError())
	})

	t.Run("Nil provider data", func(t *testing.T) {
		req := datasource.ConfigureRequest{
			ProviderData: nil,
		}
		resp := &datasource.ConfigureResponse{}

		client := SetGraphStableClientForDataSource(ctx, req, resp, "test_datasource")
		assert.Nil(t, client)
	})
}

// TestSetGraphBetaClientForResource validates the SetGraphBetaClientForResource helper.
func TestSetGraphBetaClientForResource(t *testing.T) {
	ctx := context.Background()

	t.Run("Valid provider data", func(t *testing.T) {
		mockClients := NewMockGraphClients(nil)
		req := resource.ConfigureRequest{
			ProviderData: mockClients,
		}
		resp := &resource.ConfigureResponse{}

		client := SetGraphBetaClientForResource(ctx, req, resp, "test_resource")
		assert.NotNil(t, client)
		assert.False(t, resp.Diagnostics.HasError())
	})
}

// TestSetGraphBetaClientForDataSource validates the SetGraphBetaClientForDataSource helper.
func TestSetGraphBetaClientForDataSource(t *testing.T) {
	ctx := context.Background()

	t.Run("Valid provider data", func(t *testing.T) {
		mockClients := NewMockGraphClients(nil)
		req := datasource.ConfigureRequest{
			ProviderData: mockClients,
		}
		resp := &datasource.ConfigureResponse{}

		client := SetGraphBetaClientForDataSource(ctx, req, resp, "test_datasource")
		assert.NotNil(t, client)
		assert.False(t, resp.Diagnostics.HasError())
	})
}

// TestSetGraphBetaClientForEphemeralResource validates the SetGraphBetaClientForEphemeralResource helper.
func TestSetGraphBetaClientForEphemeralResource(t *testing.T) {
	ctx := context.Background()

	t.Run("Valid provider data", func(t *testing.T) {
		mockClients := NewMockGraphClients(nil)
		req := ephemeral.ConfigureRequest{
			ProviderData: mockClients,
		}
		resp := &ephemeral.ConfigureResponse{}

		client := SetGraphBetaClientForEphemeralResource(ctx, req, resp, "test_ephemeral")
		assert.NotNil(t, client)
		assert.False(t, resp.Diagnostics.HasError())
	})
}

// TestSetGraphBetaClientForAction validates the SetGraphBetaClientForAction helper.
func TestSetGraphBetaClientForAction(t *testing.T) {
	ctx := context.Background()

	t.Run("Valid provider data", func(t *testing.T) {
		mockClients := NewMockGraphClients(nil)
		req := action.ConfigureRequest{
			ProviderData: mockClients,
		}
		resp := &action.ConfigureResponse{}

		client := SetGraphBetaClientForAction(ctx, req, resp, "test_action")
		assert.NotNil(t, client)
		assert.False(t, resp.Diagnostics.HasError())
	})
}

// TestSetGraphBetaClientForListResource validates the SetGraphBetaClientForListResource helper.
func TestSetGraphBetaClientForListResource(t *testing.T) {
	ctx := context.Background()

	t.Run("Valid provider data", func(t *testing.T) {
		mockClients := NewMockGraphClients(nil)
		req := list.ConfigureRequest{
			ProviderData: mockClients,
		}
		resp := &list.ConfigureResponse{}

		client := SetGraphBetaClientForListResource(ctx, req, resp, "test_list")
		assert.NotNil(t, client)
		assert.False(t, resp.Diagnostics.HasError())
	})
}

// TestSetGraphV1HTTPClientForResource validates the SetGraphV1HTTPClientForResource helper.
func TestSetGraphV1HTTPClientForResource(t *testing.T) {
	ctx := context.Background()

	t.Run("Valid provider data", func(t *testing.T) {
		mockClients := NewMockGraphClients(nil)
		req := resource.ConfigureRequest{
			ProviderData: mockClients,
		}
		resp := &resource.ConfigureResponse{}

		client := SetGraphV1HTTPClientForResource(ctx, req, resp, "test_resource")
		assert.NotNil(t, client)
		assert.False(t, resp.Diagnostics.HasError())
	})
}

// TestSetGraphV1HTTPClientForDataSource validates the SetGraphV1HTTPClientForDataSource helper.
func TestSetGraphV1HTTPClientForDataSource(t *testing.T) {
	ctx := context.Background()

	t.Run("Valid provider data", func(t *testing.T) {
		mockClients := NewMockGraphClients(nil)
		req := datasource.ConfigureRequest{
			ProviderData: mockClients,
		}
		resp := &datasource.ConfigureResponse{}

		client := SetGraphV1HTTPClientForDataSource(ctx, req, resp, "test_datasource")
		assert.NotNil(t, client)
		assert.False(t, resp.Diagnostics.HasError())
	})
}

// TestSetGraphBetaHTTPClientForResource validates the SetGraphBetaHTTPClientForResource helper.
func TestSetGraphBetaHTTPClientForResource(t *testing.T) {
	ctx := context.Background()

	t.Run("Valid provider data", func(t *testing.T) {
		mockClients := NewMockGraphClients(nil)
		req := resource.ConfigureRequest{
			ProviderData: mockClients,
		}
		resp := &resource.ConfigureResponse{}

		client := SetGraphBetaHTTPClientForResource(ctx, req, resp, "test_resource")
		assert.NotNil(t, client)
		assert.False(t, resp.Diagnostics.HasError())
	})
}

// TestSetGraphBetaHTTPClientForDataSource validates the SetGraphBetaHTTPClientForDataSource helper.
func TestSetGraphBetaHTTPClientForDataSource(t *testing.T) {
	ctx := context.Background()

	t.Run("Valid provider data", func(t *testing.T) {
		mockClients := NewMockGraphClients(nil)
		req := datasource.ConfigureRequest{
			ProviderData: mockClients,
		}
		resp := &datasource.ConfigureResponse{}

		client := SetGraphBetaHTTPClientForDataSource(ctx, req, resp, "test_datasource")
		assert.NotNil(t, client)
		assert.False(t, resp.Diagnostics.HasError())
	})
}
