package helpers

import (
	"context"
	"crypto"
	"crypto/rsa"
	"crypto/x509"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	pkcs12 "software.sslmate.com/src/go-pkcs12"
)

// ParseCertificateData decodes and parses PKCS#12 data, extracting certificates and a private key.
//
// This function attempts to decode PKCS#12 data using the provided password. It extracts
// the certificate chain (including the end-entity certificate and any CA certificates),
// as well as the private key associated with the end-entity certificate.
//
// The function performs several validations:
// - It checks if any certificates are present in the decoded data.
// - It verifies the presence of a private key.
// - It ensures the private key is of RSA type.
//
// The function logs debug, error, and info messages at various stages of the process.
//
// Parameters:
//   - ctx: A context.Context for logging and potential cancellation.
//   - certData: A byte slice containing the PKCS#12 data to be parsed.
//   - password: A byte slice containing the password to decrypt the PKCS#12 data.
//
// Returns:
//   - []*x509.Certificate: A slice of parsed X.509 certificates, with the end-entity
//     certificate as the first element, followed by any CA certificates.
//   - crypto.PrivateKey: The private key associated with the end-entity certificate.
//   - error: An error if any step of the parsing or validation process fails. This will
//     be nil if the function executes successfully.
//
// Possible errors:
//   - Failure to parse PKCS#12 data
//   - No certificates found in the PKCS#12 data
//   - No private key found in the PKCS#12 data
//   - Private key is not of RSA type
//
// Usage example:
//
//	certs, privKey, err := ParseCertificateData(ctx, pkcs12Data, []byte("password"))
//	if err != nil {
//	    // Handle error
//	}
//	// Use certs and privKey as needed
func ParseCertificateData(ctx context.Context, certData []byte, password []byte) ([]*x509.Certificate, crypto.PrivateKey, error) {
	tflog.Debug(ctx, "Attempting to parse PKCS#12 data")

	// Decode the PKCS#12 data
	privateKey, certificate, caCerts, err := pkcs12.DecodeChain(certData, string(password))
	if err != nil {
		tflog.Error(ctx, "Failed to parse PKCS#12 data", map[string]interface{}{
			"error": err,
		})
		return nil, nil, fmt.Errorf("failed to parse PKCS#12 data: %v", err)
	}
	certs := append([]*x509.Certificate{certificate}, caCerts...)

	validCerts := []*x509.Certificate{}
	for _, cert := range certs {
		if cert != nil {
			validCerts = append(validCerts, cert)
		}
	}

	if len(validCerts) == 0 {
		tflog.Error(ctx, "No valid certificates found in PKCS#12 data")
		return nil, nil, errors.New("no valid certificates found in PKCS#12 data")
	}

	if privateKey == nil {
		tflog.Error(ctx, "No private key found in PKCS#12 data")
		return nil, nil, errors.New("no private key found in PKCS#12 data")
	}

	rsaKey, ok := privateKey.(*rsa.PrivateKey)
	if !ok {
		tflog.Error(ctx, "Private key is not of RSA type", map[string]interface{}{
			"actualType": fmt.Sprintf("%T", privateKey),
		})
		return nil, nil, fmt.Errorf("private key is not of RSA type, got %T", privateKey)
	}

	tflog.Info(ctx, "PKCS#12 data parsed successfully", map[string]interface{}{
		"certificateCount": len(validCerts),
		"privateKeyType":   "RSA",
		"privateKeyBits":   rsaKey.N.BitLen(),
	})

	return validCerts, privateKey, nil
}
