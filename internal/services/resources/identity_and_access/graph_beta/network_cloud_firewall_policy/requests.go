package graphBetaNetworkCloudFirewallPolicy

import (
	"context"
	"fmt"
	"net/http"
	"time"

	abstractions "github.com/microsoft/kiota-abstractions-go"
	graph "github.com/microsoftgraph/msgraph-beta-sdk-go/models/networkaccess"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models/odataerrors"
)

func (r *NetworkCloudFirewallPolicyResource) createPolicy(ctx context.Context, body *policyRequest) (graph.CloudFirewallPolicyable, error) {
	builder := r.client.NetworkAccess().CloudFirewallPolicies()
	info, err := builder.ToPostRequestInformation(ctx, body, nil)
	if err != nil {
		return nil, fmt.Errorf("build cloud firewall create request: %w", err)
	}
	// Retry explicit throttling only: ambiguous server failures can leave an untracked duplicate.
	retryOptions := r.retryOptions
	retryOptions.ShouldRetry = func(_ time.Duration, _ int, _ *http.Request, response *http.Response) bool {
		return response.StatusCode == http.StatusTooManyRequests
	}
	info.AddRequestOptions([]abstractions.RequestOption{&retryOptions})
	result, err := r.client.GetAdapter().Send(ctx, info, graph.CreateCloudFirewallPolicyFromDiscriminatorValue, abstractions.ErrorMappings{"XXX": odataerrors.CreateODataErrorFromDiscriminatorValue})
	if err != nil || result == nil {
		return nil, err //nolint:wrapcheck // Keep the Kiota error type for shared Graph HTTP status handling.
	}
	return result.(graph.CloudFirewallPolicyable), nil
}
