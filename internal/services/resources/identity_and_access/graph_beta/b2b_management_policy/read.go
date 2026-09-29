package graphBetaIdentityAndAccessB2bManagementPolicy

import (
	"context"
	"fmt"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud"
	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// Variables so unit tests can shorten them (export_test.go).
var (
	notFoundConfirmationWindow   = 60 * time.Second
	notFoundConfirmationInterval = 3 * time.Second
)

// getPolicy reads a policy. With confirmNotFound, a 404 is returned only after it persists for
// notFoundConfirmationWindow: stale Entra replicas return 404 for recently created policies, and
// dropping the resource from state on one would make the next plan create a duplicate.
func (r *B2bManagementPolicyResource) getPolicy(ctx context.Context, id string, confirmNotFound bool) (graphmodels.B2bManagementPolicyable, error) {
	get := func(ctx context.Context) (graphmodels.B2bManagementPolicyable, error) {
		return r.client.
			Policies().
			B2bManagementPolicies().
			ByB2bManagementPolicyId(id).
			Get(ctx, nil)
	}

	policy, err := get(ctx)
	if err == nil || !confirmNotFound || errors.GraphError(ctx, err).StatusCode != 404 {
		return policy, err
	}

	tflog.Debug(ctx, fmt.Sprintf("B2B management policy %s returned 404, confirming it is not a stale Entra replica", id))

	select {
	case <-time.After(notFoundConfirmationInterval):
	case <-ctx.Done():
		return nil, ctx.Err()
	}

	pollCtx, cancel := context.WithTimeout(ctx, notFoundConfirmationWindow)
	defer cancel()

	pollErr := crud.PollUntil(pollCtx, notFoundConfirmationInterval, func(ctx context.Context) (bool, error) {
		found, getErr := get(ctx)
		if getErr == nil {
			policy = found
			return true, nil
		}
		if errors.GraphError(ctx, getErr).StatusCode != 404 {
			return false, &crud.FatalPollError{Err: getErr}
		}
		return false, getErr
	})

	switch {
	case pollErr == nil:
		tflog.Debug(ctx, fmt.Sprintf("B2B management policy %s is readable again; the 404 came from a stale replica", id))
		return policy, nil
	case pollCtx.Err() != nil && ctx.Err() == nil:
		return nil, err
	default:
		return nil, pollErr
	}
}
