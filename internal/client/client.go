package client

import (
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
	msgraphsdk "github.com/microsoftgraph/msgraph-sdk-go"
)

// GraphClients encapsulates both the stable and beta GraphServiceClients
// provided by the Microsoft Graph SDKs. These clients are used to interact
// with the Microsoft Graph API and its beta endpoints, respectively.
//
// The stable client (StableClient) is used for making API calls to the
// stable Microsoft Graph endpoints, which are generally considered
// production-ready and have a higher level of reliability and support.
//
// The beta client (BetaClient) is used for making API calls to the
// beta Microsoft Graph endpoints, which might include new or experimental
// features that are not yet available in the stable API. These endpoints
// are subject to change and should be used with caution in production environments.
//
// Fields:
//
//	StableClient (*msgraphsdk.GraphServiceClient): The client for interacting
//	  with the stable Microsoft Graph API. This client provides access to
//	  the well-supported and stable endpoints for production use.
//
//	BetaClient (*msgraphbetasdk.GraphServiceClient): The client for interacting
//	  with the beta Microsoft Graph API. This client provides access to
//	  the beta endpoints, which may include new and experimental features
//	  that are not yet available in the stable API and may be subject to
//	  changes or deprecation.
//
// Usage:
// The GraphClients struct is intended to be instantiated and configured by
// the provider during initialization, and then passed to the resources that
// need to interact with the Microsoft Graph API. This separation of stable
// and beta clients allows resources to choose the appropriate client based
// on the API features they require.
type GraphClients struct {
	StableClient *msgraphsdk.GraphServiceClient
	BetaClient   *msgraphbetasdk.GraphServiceClient
}
