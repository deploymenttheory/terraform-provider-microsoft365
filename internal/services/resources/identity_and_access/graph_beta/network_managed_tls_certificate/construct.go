package graphBetaNetworkManagedTLSCertificate

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	s "github.com/microsoft/kiota-abstractions-go/serialization"
)

const (
	managedTLSCertificateNamePrefix = "M-TLSi-"
	managedTLSCertificateNameChars  = "0123456789abcdefghijklmnopqrstuvwxyz"
)

// constructResource builds the payload observed from the Entra Global Secure
// Access Microsoft-managed TLS certificate blade.
func constructResource(ctx context.Context, data *NetworkManagedTLSCertificateResourceModel) (s.Parsable, error) {
	if data.Name.IsNull() || data.Name.IsUnknown() {
		name, err := generateManagedTLSCertificateName()
		if err != nil {
			return nil, fmt.Errorf("generate Microsoft-managed TLS certificate name: %w", err)
		}
		data.Name = types.StringValue(name)
	}

	requestBody := &managedTLSCertificateRequestBody{
		name:             data.Name.ValueStringPointer(),
		commonName:       data.CommonName.ValueStringPointer(),
		organizationName: data.OrganizationName.ValueStringPointer(),
		validityMonths:   data.ValidityMonths.ValueInt32Pointer(),
	}

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]any{"error": err.Error()})
	}

	return requestBody, nil
}

func constructStatusUpdate(ctx context.Context, enabled bool) (s.Parsable, error) {
	status := "disabled"
	if enabled {
		status = "enabled"
	}
	requestBody := &managedTLSCertificateRequestBody{status: &status}
	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]any{"error": err.Error()})
	}
	return requestBody, nil
}

func generateManagedTLSCertificateName() (string, error) {
	name := make([]byte, 5)
	max := big.NewInt(int64(len(managedTLSCertificateNameChars)))
	for i := range name {
		index, err := rand.Int(rand.Reader, max)
		if err != nil {
			return "", err
		}
		name[i] = managedTLSCertificateNameChars[index.Int64()]
	}
	return managedTLSCertificateNamePrefix + string(name), nil
}

type managedTLSCertificateRequestBody struct {
	name             *string
	commonName       *string
	organizationName *string
	validityMonths   *int32
	status           *string
}

func (b *managedTLSCertificateRequestBody) Serialize(writer s.SerializationWriter) error {
	if err := writer.WriteStringValue("name", b.name); err != nil {
		return err
	}
	if err := writer.WriteStringValue("commonName", b.commonName); err != nil {
		return err
	}
	if err := writer.WriteStringValue("organizationName", b.organizationName); err != nil {
		return err
	}
	if err := writer.WriteInt32Value("validityMonths", b.validityMonths); err != nil {
		return err
	}
	return writer.WriteStringValue("status", b.status)
}

func (b *managedTLSCertificateRequestBody) GetFieldDeserializers() map[string]func(s.ParseNode) error {
	return map[string]func(s.ParseNode) error{}
}
