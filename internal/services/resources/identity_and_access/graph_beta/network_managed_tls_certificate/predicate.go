package graphBetaNetworkManagedTLSCertificate

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud"
	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
)

const managedTLSCertificateStatusPollInterval = 5 * time.Second

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
		if enabled {
			switch status {
			case "active", "expiring":
				return true, nil
			case "expired", "revoked":
				return false, &crud.FatalPollError{Err: fmt.Errorf("certificate entered terminal status %q while enabling", status)}
			default:
				return false, fmt.Errorf("certificate activation is still pending with status %q", status)
			}
		}

		// unknownfuturevalue is how Graph reports disabled if the Prefer header is dropped.
		if status == "disabled" || status == "unknownfuturevalue" {
			return true, nil
		}
		if status == "expired" || status == "revoked" {
			return false, &crud.FatalPollError{Err: fmt.Errorf("certificate entered terminal status %q while disabling", status)}
		}
		return false, fmt.Errorf("certificate deactivation is still pending with status %q", status)
	})
	if err != nil {
		// latest is returned with the error so callers can report the last status seen.
		return latest, fmt.Errorf("wait for certificate status: %w", err)
	}
	return latest, nil
}

// stillActivating reports whether a failed activation wait only ran out of time
// while Graph was still enrolling the certificate.
func stillActivating(ctx context.Context, latest *managedTLSCertificateResponse) bool {
	return ctx.Err() != nil && latest != nil && latest.status != nil &&
		strings.EqualFold(*latest.status, "enrolling")
}
