package graphBetaNetworkManagedTLSCertificate

import (
	"context"
	stderrors "errors"
	"fmt"
	"strings"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud"
	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
)

const managedTLSCertificateStatusPollInterval = 5 * time.Second

var (
	// errTerminalStatus marks a certificate that can no longer reach the requested status.
	errTerminalStatus = stderrors.New("certificate entered a terminal status")
	errPendingEnable  = stderrors.New("certificate activation is still pending")
	errPendingDisable = stderrors.New("certificate deactivation is still pending")
)

func (r *NetworkManagedTLSCertificateResource) waitForManagedTLSCertificateStatus(ctx context.Context, certificateID string, enabled bool) (*managedTLSCertificateResponse, error) {
	var latest *managedTLSCertificateResponse

	err := crud.PollUntil(ctx, managedTLSCertificateStatusPollInterval, func(ctx context.Context) (bool, error) {
		certificate, err := r.getManagedTLSCertificate(ctx, certificateID)
		if err != nil {
			info := errors.GraphError(ctx, err)
			if info.StatusCode == 400 || info.StatusCode == 401 || info.StatusCode == 403 {
				return false, &crud.FatalPollError{Err: err}
			}
			return false, err
		}
		latest = certificate
		if certificate.status == nil {
			return false, fmt.Errorf("Microsoft Graph returned the certificate without a lifecycle status")
		}

		status := strings.ToLower(*certificate.status)
		switch {
		case status == "expired" || status == "revoked":
			return false, &crud.FatalPollError{Err: fmt.Errorf("%w: %q", errTerminalStatus, status)}
		case enabled && (status == "active" || status == "expiring"):
			return true, nil
		case !enabled && status == "disabled":
			return true, nil
		case enabled:
			return false, fmt.Errorf("%w with status %q", errPendingEnable, status)
		default:
			return false, fmt.Errorf("%w with status %q", errPendingDisable, status)
		}
	})
	if err != nil {
		// latest is returned with the error so callers can report the last status seen.
		return latest, fmt.Errorf("wait for certificate status: %w", err)
	}
	return latest, nil
}
