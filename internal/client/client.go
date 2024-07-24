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
// The v1.0 endpoint of Microsoft Graph provides a stable and reliable API
// that is fully supported by Microsoft, ensuring that applications built
// on this endpoint have a solid foundation and offer the best possible
// user experience.
//
// The beta client (BetaClient) is used for making API calls to the
// beta Microsoft Graph endpoints, which allow developers to test and
// experiment with new features before they are released to the general public.
// However, it is important to note that the beta endpoint is not intended
// for use in production environments. APIs and functionalities available
// in the beta endpoint are subject to change, and features might be modified
// or removed without notice, potentially causing disruptions or breaking
// changes to applications. Additionally, the beta endpoint might not have
// the same level of support, reliability, or performance as the v1.0 endpoint,
// and can be unexpectedly unavailable or have slower response times.
//
// For these reasons, developers are strongly recommended to use the v1.0
// endpoint when building production applications to ensure stability,
// reliability, and a better user experience.
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
type GraphClients struct {
	StableClient *msgraphsdk.GraphServiceClient
	BetaClient   *msgraphbetasdk.GraphServiceClient
}
