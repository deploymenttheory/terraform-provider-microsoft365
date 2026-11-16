package graphBetaIdentityAndAccessSecurityDefaultsEnforcementPolicy

import (
	"context"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
)

// enableConflictSummary is the diagnostic summary added when Graph rejects enabling security defaults.
const enableConflictSummary = "Security defaults could not be enabled"

// enableConflictDetail explains the most common reason Graph rejects enabling security defaults.
const enableConflictDetail = "Microsoft Entra ID rejects enabling security defaults while Conditional Access policies are enabled in the tenant. " +
	"Security defaults and Conditional Access are mutually exclusive: disable or delete the enabled Conditional Access policies first, " +
	"or keep is_enabled = false and protect the tenant with Conditional Access instead. See the Graph error above for the exact reason."

// addEnableConflictDiagnostic adds an explanatory error after HandleKiotaGraphError when a PATCH that enables
// security defaults is rejected with 400 Bad Request or 409 Conflict. Other failures (for example 403) are already
// explained by the shared Graph error handling and get no extra diagnostic.
func addEnableConflictDiagnostic(ctx context.Context, err error, isEnabled types.Bool, diagnostics *diag.Diagnostics) {
	if err == nil || !isEnabled.ValueBool() {
		return
	}
	statusCode := errors.GraphError(ctx, err).StatusCode
	if statusCode != http.StatusBadRequest && statusCode != http.StatusConflict {
		return
	}
	diagnostics.AddError(enableConflictSummary, enableConflictDetail)
}
