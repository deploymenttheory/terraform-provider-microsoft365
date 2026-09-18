package graphBetaNetworkCloudFirewallPolicyRule

import (
	"context"
	"fmt"
	"net/http"
	"time"

	abstractions "github.com/microsoft/kiota-abstractions-go"
	graph "github.com/microsoftgraph/msgraph-beta-sdk-go/models/networkaccess"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models/odataerrors"
)

func (r *NetworkCloudFirewallPolicyRuleResource) createRule(ctx context.Context, data *NetworkCloudFirewallPolicyRuleResourceModel) (graph.PolicyRuleable, error) {
	body, err := constructRule(ctx, data, nil)
	if err != nil {
		return nil, fmt.Errorf("build cloud firewall create request: %w", err)
	}
	builder := r.client.NetworkAccess().CloudFirewallPolicies().ByCloudFirewallPolicyId(data.PolicyID.ValueString()).PolicyRules()
	info, err := builder.ToPostRequestInformation(ctx, body, nil)
	if err != nil {
		return nil, err //nolint:wrapcheck // Keep the Kiota error type for shared Graph HTTP status handling.
	}
	// Retry explicit throttling only: ambiguous server failures can leave an untracked duplicate.
	retryOptions := r.retryOptions
	retryOptions.ShouldRetry = func(_ time.Duration, _ int, _ *http.Request, response *http.Response) bool {
		return response.StatusCode == http.StatusTooManyRequests
	}
	info.AddRequestOptions([]abstractions.RequestOption{&retryOptions})
	result, err := r.client.GetAdapter().Send(ctx, info, graph.CreatePolicyRuleFromDiscriminatorValue, abstractions.ErrorMappings{"XXX": odataerrors.CreateODataErrorFromDiscriminatorValue})
	if err != nil || result == nil {
		return nil, err //nolint:wrapcheck // Keep the Kiota error type for shared Graph HTTP status handling.
	}
	return result.(graph.PolicyRuleable), nil
}
