package graphBetaApplicationCertificateCredential

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/google/uuid"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// FindKeyCredentialByDisplayName searches through a slice of KeyCredentials and returns
// the keyId of the first credential matching the specified display name.
// Returns nil if no matching credential is found.
func FindKeyCredentialByDisplayName(credentials []graphmodels.KeyCredentialable, displayName string) *uuid.UUID {
	for _, cred := range credentials {
		if cred.GetDisplayName() != nil && *cred.GetDisplayName() == displayName {
			return cred.GetKeyId()
		}
	}
	return nil
}

// FindKeyCredentialByKeyID searches through a slice of KeyCredentials and returns
// the credential matching the specified keyId.
// Returns nil if no matching credential is found.
func FindKeyCredentialByKeyID(credentials []graphmodels.KeyCredentialable, keyID string) graphmodels.KeyCredentialable {
	for _, cred := range credentials {
		if cred.GetKeyId() != nil && cred.GetKeyId().String() == keyID {
			return cred
		}
	}
	return nil
}

// decodeCertificateKey decodes the certificate key value based on the specified encoding.
// Supported encodings:
//   - pem: PEM format with -----BEGIN/END CERTIFICATE----- headers. Extracts and decodes the base64 content.
//   - base64: Raw base64-encoded DER certificate. Decodes directly.
//   - hex: Hex-encoded certificate. Decodes hex to raw bytes.
//
// Returns the raw certificate bytes for the Graph API.
func decodeCertificateKey(key string, encoding string) ([]byte, error) {
	switch encoding {
	case "pem":
		return decodePEM(key)
	case "base64":
		return base64.StdEncoding.DecodeString(key)
	case "hex":
		return hex.DecodeString(key)
	default:
		return nil, fmt.Errorf("unsupported encoding: %s", encoding)
	}
}

// decodePEM extracts and decodes the base64 content from a PEM-formatted certificate.
// PEM format: -----BEGIN CERTIFICATE-----\n<base64 content>\n-----END CERTIFICATE-----
func decodePEM(pem string) ([]byte, error) {
	// Remove PEM headers
	content := pem
	content = strings.ReplaceAll(content, "-----BEGIN CERTIFICATE-----", "")
	content = strings.ReplaceAll(content, "-----END CERTIFICATE-----", "")
	content = strings.ReplaceAll(content, "\n", "")
	content = strings.ReplaceAll(content, "\r", "")
	content = strings.TrimSpace(content)

	if content == "" {
		return nil, fmt.Errorf("PEM certificate content is empty after removing headers")
	}

	return base64.StdEncoding.DecodeString(content)
}
