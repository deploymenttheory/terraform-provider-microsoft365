package graphBetaAuthorizationPolicy

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud"
	grapherrors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
)

// completeInvitationWrite waits for the settings write to become visible before
// applying invitations. A successful PATCH can still be queued by Graph.
func (r *AuthorizationPolicyResource) completeInvitationWrite(
	ctx context.Context,
	body graphmodels.AuthorizationPolicyable,
	desired tfsdk.State,
) error {
	if body.GetAllowInvitesFrom() == nil || body.GetDefaultUserRolePermissions() == nil {
		return nil
	}
	if err := r.waitForPolicySettings(ctx, desired); err != nil {
		return fmt.Errorf("wait for authorization settings before invitations: %w", err)
	}
	invitations := graphmodels.NewAuthorizationPolicy()
	invitations.SetAllowInvitesFrom(body.GetAllowInvitesFrom())
	if err := r.patchSettings(ctx, invitations); err != nil {
		return fmt.Errorf("update invitation settings: %w", err)
	}
	return nil
}

func (r *AuthorizationPolicyResource) waitForPolicySettings(ctx context.Context, desired tfsdk.State) error {
	var planned AuthorizationPolicyResourceModel
	if diags := desired.Get(ctx, &planned); diags.HasError() {
		return fmt.Errorf("%w: read planned settings: %v", errInvalidPolicyState, diags.Errors())
	}
	names := make([]string, 0, len(readbackAttributes))
	for _, name := range readbackAttributes {
		if name != "allow_invites_from" {
			names = append(names, name)
		}
	}
	match := readbackMatch
	match.Paths = crud.RootPaths(names...)
	ready := crud.RequireStableReads(crud.AttributesMatchPredicate(desired, match))
	err := crud.PollUntil(ctx, 2*time.Second, func(ctx context.Context) (bool, error) {
		response, err := r.client.Policies().AuthorizationPolicy().Get(ctx, nil)
		if err != nil {
			wrapped := fmt.Errorf("read pending authorization settings: %w", err)
			if grapherrors.GraphError(ctx, err).StatusCode == 404 {
				return false, wrapped
			}
			return false, &crud.FatalPollError{Err: wrapped}
		}
		current := planned
		if err := mapRemoteState(ctx, &current, response); err != nil {
			return false, &crud.FatalPollError{Err: fmt.Errorf("map pending authorization settings: %w", err)}
		}
		state := tfsdk.State{Schema: desired.Schema}
		if diags := state.Set(ctx, &current); diags.HasError() {
			return false, &crud.FatalPollError{Err: fmt.Errorf("%w: project pending settings: %v", errInvalidPolicyState, diags.Errors())}
		}
		return ready(ctx, state), nil
	})
	if err != nil {
		return fmt.Errorf("confirm authorization settings: %w", err)
	}
	return nil
}

// The policy is a permanent singleton. Back-to-back updates can temporarily return
// 404 while the preceding change replicates; other errors remain immediately visible.
func (r *AuthorizationPolicyResource) patchSettings(ctx context.Context, body graphmodels.AuthorizationPolicyable) error {
	policy := r.client.Policies().AuthorizationPolicy().ByAuthorizationPolicyId(singletonID)
	for attempt := 0; ; attempt++ {
		_, err := policy.Patch(ctx, body, nil)
		if err == nil {
			return nil
		}
		if grapherrors.GraphError(ctx, err).StatusCode != 404 || attempt >= 15 {
			return fmt.Errorf("patch authorization policy: %w", err)
		}
		select {
		case <-ctx.Done():
			return fmt.Errorf("wait for authorization policy replication: %w", ctx.Err())
		case <-time.After(2 * time.Second):
		}
	}
}
