package graphBetaNetworkManagedTLSCertificate

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const (
	managedTLSCertificateNamePrefix = "M-TLSi-"
	managedTLSCertificateNameChars  = "0123456789abcdefghijklmnopqrstuvwxyz"
)

// constructResource builds the payload observed from the Entra Global Secure
// Access Microsoft-managed TLS certificate blade. A name the configuration
// leaves unset is generated in the portal's format and written back to data.
func constructResource(ctx context.Context, data *NetworkManagedTLSCertificateResourceModel) (managedTLSCertificateRequest, error) {
	if data.Name.IsNull() || data.Name.IsUnknown() {
		name := make([]byte, 5)
		maxIndex := big.NewInt(int64(len(managedTLSCertificateNameChars)))
		for i := range name {
			index, err := rand.Int(rand.Reader, maxIndex)
			if err != nil {
				return managedTLSCertificateRequest{}, fmt.Errorf("generate Microsoft-managed TLS certificate name: %w", err)
			}
			name[i] = managedTLSCertificateNameChars[index.Int64()]
		}
		data.Name = types.StringValue(managedTLSCertificateNamePrefix + string(name))
	}

	requestBody := managedTLSCertificateRequest{
		Name:             data.Name.ValueStringPointer(),
		CommonName:       data.CommonName.ValueStringPointer(),
		OrganizationName: data.OrganizationName.ValueStringPointer(),
		ValidityMonths:   data.ValidityMonths.ValueInt32Pointer(),
	}

	tflog.Debug(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName), map[string]any{"request_body": requestBody})

	return requestBody, nil
}
