package client

import (
	"testing"

	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
	msgraphsdk "github.com/microsoftgraph/msgraph-sdk-go"
	"github.com/stretchr/testify/assert"
)

// MockGraphServiceClient is a mock implementation of the GraphServiceClient for testing purposes.
type MockGraphServiceClient struct{}

// Ensure MockGraphServiceClient implements the necessary methods for both stable and beta clients.

func TestGraphClients(t *testing.T) {
	stableClient := &MockGraphServiceClient{}
	betaClient := &MockGraphServiceClient{}

	clients := GraphClients{
		StableClient: stableClient,
		BetaClient:   betaClient,
	}

	assert.NotNil(t, clients.StableClient, "StableClient should not be nil")
	assert.NotNil(t, clients.BetaClient, "BetaClient should not be nil")

	assert.Equal(t, stableClient, clients.StableClient, "StableClient should be set correctly")
	assert.Equal(t, betaClient, clients.BetaClient, "BetaClient should be set correctly")
}

func TestGraphClientsInitialization(t *testing.T) {
	stableClient := &msgraphsdk.GraphServiceClient{}
	betaClient := &msgraphbetasdk.GraphServiceClient{}

	clients := GraphClients{
		StableClient: stableClient,
		BetaClient:   betaClient,
	}

	if clients.StableClient != stableClient {
		t.Errorf("Expected StableClient to be %v, got %v", stableClient, clients.StableClient)
	}

	if clients.BetaClient != betaClient {
		t.Errorf("Expected BetaClient to be %v, got %v", betaClient, clients.BetaClient)
	}
}
