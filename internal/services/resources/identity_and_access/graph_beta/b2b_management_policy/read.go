package graphBetaIdentityAndAccessB2bManagementPolicy

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud"
	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
)

// Variables so unit tests can shorten them (export_test.go).
var (
	notFoundConfirmationWindow   = 60 * time.Second
	notFoundConfirmationInterval = 3 * time.Second
)

// confirmPolicyNotFound polls after a 404 until the policy becomes readable or the confirmation
// window expires. The caller owns the Graph call and passes the original not-found error.
func confirmPolicyNotFound(
	ctx context.Context,
	id string,
	notFoundErr error,
	get func(context.Context) (graphmodels.B2bManagementPolicyable, error),
) (graphmodels.B2bManagementPolicyable, error) {
	tflog.Debug(ctx, fmt.Sprintf("B2B management policy %s returned 404, confirming it is not a stale Entra replica", id))

	select {
	case <-time.After(notFoundConfirmationInterval):
	case <-ctx.Done():
		return nil, ctx.Err()
	}

	pollCtx, cancel := context.WithTimeout(ctx, notFoundConfirmationWindow)
	defer cancel()

	var policy graphmodels.B2bManagementPolicyable
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
		return nil, notFoundErr
	case ctx.Err() != nil:
		return nil, fmt.Errorf("confirming B2B management policy %s absence: %w", id, ctx.Err())
	default:
		return nil, pollErr
	}
}
