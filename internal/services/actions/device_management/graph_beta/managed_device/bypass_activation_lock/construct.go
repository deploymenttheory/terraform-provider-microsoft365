package graphBetaBypassActivationLockManagedDevice

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

// constructBypassActivationLockRequest validates devices and constructs the request
// The Bypass Activation Lock endpoint requires no request body parameters
// This function performs API validation before the bypass operation and returns validation results
func constructBypassActivationLockRequest(ctx context.Context, client *msgraphbetasdk.GraphServiceClient, deviceIDs []string) (*ValidationResult, error) {
	tflog.Debug(ctx, "Constructing Activation Lock bypass request")

	// Perform API validation of devices
	validation, err := validateRequest(ctx, client, deviceIDs)
	if err != nil {
		tflog.Error(ctx, "Failed to validate devices via API", map[string]any{"error": err.Error()})
		return nil, err
	}

	// Debug log the request (no body for this endpoint)
	if err := constructors.DebugLogGraphObject(ctx, "Bypass Activation Lock request (no body required)", nil); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]any{"error": err.Error()})
	}

	tflog.Debug(ctx, "Request construction and validation completed successfully")
	return validation, nil
}
