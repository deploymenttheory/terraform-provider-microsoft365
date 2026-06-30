package graphBetaMacOSDepEnrollmentProfile

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	errDeviceManagementNil = errors.New("deviceManagement response is nil")
	errIntuneAccountIdNil  = errors.New("intuneAccountId is nil in deviceManagement response")
)

// resolveDepOnboardingSettingsId determines the depOnboardingSetting id to use.
// This id is the 'intuneAccountId' in the /deviceManagement endpoint.
func (r *MacOSDepEnrollmentProfileResource) resolveDepOnboardingSettingsId(
	ctx context.Context,
	provided types.String,
) (string, error) {
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
		return "", errDeviceManagementNil
	}

	intuneAccountId := dm.GetIntuneAccountId()
	if intuneAccountId == nil {
		return "", errIntuneAccountIdNil
	}

	return intuneAccountId.String(), nil
}
