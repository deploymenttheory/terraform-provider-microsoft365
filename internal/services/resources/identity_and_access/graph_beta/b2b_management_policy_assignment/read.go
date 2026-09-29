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

// Variables so unit tests can shorten them (export_test.go).
var (
	notFoundConfirmationWindow   = 60 * time.Second
	notFoundConfirmationInterval = 3 * time.Second
)

var errAssignmentNotVisible = stderrors.New("assignment not visible in appliesTo")

// findAppliesTo returns the directory object from appliesTo, or nil if absent (including a 404 for
// the policy). With confirmNotFound, absence is reported only after it persists for
// notFoundConfirmationWindow: stale replicas briefly return an empty appliesTo, and dropping the
// resource from state on one would make the next plan re-create an existing assignment.
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
		return nil, nil
	default:
		return nil, pollErr
	}
}
