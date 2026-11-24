package graphBetaIdentityAndAccessNetworkRemoteNetwork

import (
	"context"
	"errors"
	"fmt"
)

var (
	errNoForwardingProfiles        = errors.New("at least one forwarding profile ID is required")
	errInvalidForwardingProfileIDs = errors.New("forwarding_profile_ids could not be read")
)

// constructResource builds the raw JSON POST body for a remote network. Graph requires at least one forwarding profile.
func constructResource(ctx context.Context, plan NetworkRemoteNetworkResourceModel) (map[string]any, error) {
	var profileIDs []string
	if !plan.ForwardingProfileIDs.IsNull() && !plan.ForwardingProfileIDs.IsUnknown() {
		if diags := plan.ForwardingProfileIDs.ElementsAs(ctx, &profileIDs, false); diags.HasError() {
			return nil, fmt.Errorf("%w: %v", errInvalidForwardingProfileIDs, diags.Errors())
		}
	}
	if len(profileIDs) == 0 {
		return nil, errNoForwardingProfiles
	}
	refs := make([]map[string]string, 0, len(profileIDs))
	for _, id := range profileIDs {
		refs = append(refs, map[string]string{"id": id})
	}
	return map[string]any{
		"name":               plan.Name.ValueString(),
		"region":             plan.Region.ValueString(),
		"forwardingProfiles": refs,
	}, nil
}
