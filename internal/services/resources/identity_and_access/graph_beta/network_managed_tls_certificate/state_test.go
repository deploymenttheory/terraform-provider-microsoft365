package graphBetaNetworkManagedTLSCertificate

import (
	"context"
	"testing"
)

func TestMapRemoteStateToTerraformDerivesEnabledFromLifecycleStatus(t *testing.T) {
	for _, test := range []struct {
		status  string
		enabled bool
	}{
		{status: "disabled", enabled: false},
		{status: "unknownFutureValue", enabled: false},
		{status: "creating", enabled: false},
		{status: "enrolling", enabled: true},
		{status: "active", enabled: true},
		{status: "expiring", enabled: true},
		{status: "expired", enabled: false},
		{status: "revoked", enabled: false},
	} {
		t.Run(test.status, func(t *testing.T) {
			model := &NetworkManagedTLSCertificateResourceModel{}
			MapRemoteStateToTerraform(context.Background(), model, &managedTLSCertificateResponse{status: &test.status})
			if model.Status.ValueString() != test.status {
				t.Fatalf("status = %q, want %q", model.Status.ValueString(), test.status)
			}
			if model.Enabled.ValueBool() != test.enabled {
				t.Fatalf("enabled = %t for status %q, want %t", model.Enabled.ValueBool(), test.status, test.enabled)
			}
		})
	}
}
