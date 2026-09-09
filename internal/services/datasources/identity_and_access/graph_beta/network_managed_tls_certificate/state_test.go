package graphBetaNetworkManagedTLSCertificate

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestMapRemoteStateToDataSourcePreservesConfiguredID(t *testing.T) {
	configuredID := "AAAAAAAA-BBBB-CCCC-DDDD-EEEEEEEEEEEE"
	remoteID := "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"
	data := &NetworkManagedTLSCertificateDataSourceModel{
		ID: types.StringValue(configuredID),
	}

	mapRemoteStateToDataSource(data, &managedTLSCertificateResponse{id: &remoteID})

	if data.ID.ValueString() != configuredID {
		t.Fatalf("id = %q, want configured value %q", data.ID.ValueString(), configuredID)
	}
}
