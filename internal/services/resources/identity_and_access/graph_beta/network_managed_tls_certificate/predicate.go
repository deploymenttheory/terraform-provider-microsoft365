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

		if status == "disabled" {
			return true, nil
		}
		if status == "expired" || status == "revoked" {
			return false, &crud.FatalPollError{Err: fmt.Errorf("certificate entered terminal status %q while disabling", status)}
		}
		return false, fmt.Errorf("certificate deactivation is still pending with status %q", status)
	})
	if err != nil {
		return nil, err
	}
	return latest, nil
}
