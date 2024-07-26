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

	if len(certs) == 0 {
		tflog.Error(ctx, "No certificates found in PKCS#12 data")
		return nil, nil, errors.New("no certificates found in PKCS#12 data")
	}

	if privateKey == nil {
		tflog.Error(ctx, "No private key found in PKCS#12 data")
		return nil, nil, errors.New("no private key found in PKCS#12 data")
	}

	rsaKey, ok := privateKey.(*rsa.PrivateKey)
	if !ok {
		tflog.Error(ctx, "Private key is not of RSA type")
		return nil, nil, errors.New("private key is not of RSA type")
	}

	tflog.Info(ctx, "PKCS#12 data parsed successfully", map[string]interface{}{
		"certificateCount": len(certs),
		"privateKeyType":   "RSA",
		"privateKeyBits":   rsaKey.N.BitLen(),
	})

	return certs, privateKey, nil
}
