// File: internal/client/mock_graph_clients.go

package client

import (
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
	msgraphsdk "github.com/microsoftgraph/msgraph-sdk-go"
	"github.com/stretchr/testify/mock"
)

// MockStableGraphServiceClient is a flexible mock of the stable GraphServiceClient
type MockStableGraphServiceClient struct {
	mock.Mock
}

// MockBetaGraphServiceClient is a flexible mock of the beta GraphServiceClient
type MockBetaGraphServiceClient struct {
	mock.Mock
}

// MockGraphClients implements the GraphClientInterface
type MockGraphClients struct {
	mock.Mock
	MockStableClient *MockStableGraphServiceClient
	MockBetaClient   *MockBetaGraphServiceClient
}

// GetStableClient returns the mock stable client
func (m *MockGraphClients) GetStableClient() *msgraphsdk.GraphServiceClient {
	args := m.Called()
	return args.Get(0).(*msgraphsdk.GraphServiceClient)
}

// GetBetaClient returns the mock beta client
func (m *MockGraphClients) GetBetaClient() *msgraphbetasdk.GraphServiceClient {
	args := m.Called()
	return args.Get(0).(*msgraphbetasdk.GraphServiceClient)
}

// NewMockGraphClients creates a new instance of MockGraphClients
func NewMockGraphClients() *MockGraphClients {
	return &MockGraphClients{
		MockStableClient: new(MockStableGraphServiceClient),
		MockBetaClient:   new(MockBetaGraphServiceClient),
	}
}
