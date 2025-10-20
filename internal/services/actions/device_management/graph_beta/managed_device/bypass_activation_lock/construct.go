package graphBetaBypassActivationLockManagedDevice

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func constructBypassActivationLockRequest(ctx context.Context) error {
	// Bypass Activation Lock endpoint requires no request body parameters
	// This function is included for consistency with action pattern architecture
	if err := constructors.DebugLogGraphObject(ctx, "Bypass Activation Lock request", nil); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]any{"error": err.Error()})
		return err
	}

	return nil
}
