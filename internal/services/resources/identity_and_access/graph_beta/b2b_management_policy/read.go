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

// notFoundConfirmationWindow bounds how long a refresh keeps re-reading a policy that returned
// 404 before accepting that it was deleted; notFoundConfirmationInterval is the delay between
// those reads. They are variables only so unit tests can shorten them (see export_test.go).
var (
	notFoundConfirmationWindow   = 60 * time.Second
	notFoundConfirmationInterval = 3 * time.Second
)

// getPolicy reads a B2B management policy. When confirmNotFound is set, a 404 is only returned
// once it has persisted for notFoundConfirmationWindow.
//
// Microsoft Entra replicas lag behind writes: seconds after a policy is created, a GET served by
// a stale replica returns 404 Directory_ObjectNotFound even though the policy exists. Terraform
// removes a resource from state when a refresh reports it missing, so a single stale 404 would
// make the next plan create a duplicate policy and orphan the existing one. Reads that follow a
// write are already retried by ReadWithRetry and pass confirmNotFound=false.
//
// See: https://devblogs.microsoft.com/identity/designing-for-eventual-consistency-for-microsoft-entra/
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

	// Re-reading immediately would most likely hit the same stale replica, so wait one interval
	// before the first confirmation read.
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
		// The confirmation window elapsed with every read returning 404: the policy is gone.
		return nil, err
	default:
		return nil, pollErr
	}
}
