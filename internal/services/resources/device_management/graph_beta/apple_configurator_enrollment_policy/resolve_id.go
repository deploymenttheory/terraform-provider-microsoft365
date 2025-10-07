package graphBetaAppleConfiguratorEnrollmentPolicy

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

// resolveDepOnboardingSettingsId determines the depOnboardingSetting id to use.
// this id is the 'intuneAccountId' in the /deviceManagement endpoint.
func (r *AppleConfiguratorEnrollmentPolicyResource) resolveDepOnboardingSettingsId(ctx context.Context, provided types.String) (string, error) {
	if !provided.IsNull() && !provided.IsUnknown() && provided.ValueString() != "" {
		return provided.ValueString(), nil
	}

	dm, err := r.client.
		DeviceManagement().
		Get(ctx, nil)

	if err != nil {
		return "", fmt.Errorf("failed to GET /deviceManagement to resolve intuneAccountId: %w", err)
	}

	if dm == nil {
		return "", fmt.Errorf("deviceManagement response is nil")
	}

	intuneAccountId := dm.GetIntuneAccountId()
	if intuneAccountId == nil {
		return "", fmt.Errorf("intuneAccountId is nil in deviceManagement response")
	}

	return intuneAccountId.String(), nil
}
