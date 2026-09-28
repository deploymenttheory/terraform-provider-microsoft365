package graphBetaIdentityAndAccessB2bManagementPolicyAssignment

import (
	"context"
	stderrors "errors"
	"fmt"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud"
	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// notFoundConfirmationWindow bounds how long a refresh keeps re-reading an assignment that is
// missing before accepting that it was removed; notFoundConfirmationInterval is the delay between
// those reads. They are variables only so unit tests can shorten them (see export_test.go).
var (
	notFoundConfirmationWindow   = 60 * time.Second
	notFoundConfirmationInterval = 3 * time.Second
)

// errAssignmentNotVisible is the pending condition reported while confirming a missing assignment.
var errAssignmentNotVisible = stderrors.New("assignment not visible in appliesTo")

// findAppliesTo returns the directory object from the policy's appliesTo, or nil when the policy
// does not apply to it (including when the policy itself returns 404). When confirmNotFound is
// set, a missing assignment is only reported once it has stayed missing for
// notFoundConfirmationWindow.
//
// Microsoft Entra replicas lag behind writes: shortly after an assignment or its policy is created,
// a stale replica returns an empty appliesTo or 404 for the policy. Terraform removes a resource
// from state when a refresh reports it missing, so a single stale read would make the next plan
// re-create an assignment that exists. Reads that follow a write are already retried by
// ReadWithRetry and pass confirmNotFound=false.
//
// See: https://devblogs.microsoft.com/identity/designing-for-eventual-consistency-for-microsoft-entra/
func (r *B2bManagementPolicyAssignmentResource) findAppliesTo(ctx context.Context, policyID, directoryObjectID string, confirmNotFound bool) (graphmodels.DirectoryObjectable, error) {
	find := func(ctx context.Context) (graphmodels.DirectoryObjectable, error) {
		directoryObjects, err := r.listAppliesTo(ctx, policyID)
		if err != nil {
			if errors.GraphError(ctx, err).StatusCode == 404 {
				return nil, nil
			}
			return nil, err
		}
		for _, directoryObject := range directoryObjects {
			if directoryObject.GetId() != nil && *directoryObject.GetId() == directoryObjectID {
				return directoryObject, nil
			}
		}
		return nil, nil
	}

	directoryObject, err := find(ctx)
	if err != nil || directoryObject != nil || !confirmNotFound {
		return directoryObject, err
	}

	tflog.Debug(ctx, fmt.Sprintf("Directory object %s not in appliesTo of B2B management policy %s, confirming it is not a stale Entra replica", directoryObjectID, policyID))

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
		found, findErr := find(ctx)
		if findErr != nil {
			return false, &crud.FatalPollError{Err: findErr}
		}
		if found == nil {
			return false, errAssignmentNotVisible
		}
		directoryObject = found
		return true, nil
	})

	switch {
	case pollErr == nil:
		tflog.Debug(ctx, fmt.Sprintf("Directory object %s is in appliesTo again; the miss came from a stale replica", directoryObjectID))
		return directoryObject, nil
	case pollCtx.Err() != nil && ctx.Err() == nil:
		// The confirmation window elapsed without the assignment reappearing: it is gone.
		return nil, nil
	default:
		return nil, pollErr
	}
}
