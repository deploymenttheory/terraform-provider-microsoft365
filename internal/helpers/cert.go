package helpers

import (
	"context"
	"crypto"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"golang.org/x/crypto/pkcs12"
)

// ParseCertificateData reads and parses the certificate data, extracting the certificate and private key.
// It first tries to parse the data as PEM. If that fails, it assumes PKCS#12 format and tries to decode it.
func ParseCertificateData(ctx context.Context, certData []byte, password []byte) ([]*x509.Certificate, crypto.PrivateKey, error) {
	var certs []*x509.Certificate
	var key crypto.PrivateKey
	var err error
	var certType string

	// Try to parse as PEM
	blocks := []*pem.Block{}
	for {
		var block *pem.Block
		block, certData = pem.Decode(certData)
		if block == nil {
			break
		}
		blocks = append(blocks, block)
	}

	for _, block := range blocks {
		switch block.Type {
		case "CERTIFICATE":
			cert, err := x509.ParseCertificate(block.Bytes)
			if err != nil {
				tflog.Error(ctx, "Failed to parse PEM certificate", map[string]interface{}{
					"error": err,
				})
				return nil, nil, err
			}
			certs = append(certs, cert)
			certType = "PEM"
		case "ENCRYPTED PRIVATE KEY":
			decryptedKey, err := x509.DecryptPEMBlock(block, password)
			if err != nil {
				tflog.Error(ctx, "Failed to decrypt PEM private key", map[string]interface{}{
					"error": err,
				})
				return nil, nil, err
			}
			key, err = x509.ParsePKCS8PrivateKey(decryptedKey)
			if err != nil {
				tflog.Error(ctx, "Failed to parse decrypted PEM private key", map[string]interface{}{
					"error": err,
				})
				return nil, nil, err
			}
		case "PRIVATE KEY":
			key, err = x509.ParsePKCS8PrivateKey(block.Bytes)
			if err != nil {
				key, err = x509.ParsePKCS1PrivateKey(block.Bytes)
			}
			if err != nil {
				tflog.Error(ctx, "Failed to parse PEM private key", map[string]interface{}{
					"error": err,
				})
				return nil, nil, err
			}
		case "RSA PRIVATE KEY":
			key, err = x509.ParsePKCS1PrivateKey(block.Bytes)
			if err != nil {
				tflog.Error(ctx, "Failed to parse PEM RSA private key", map[string]interface{}{
					"error": err,
				})
				return nil, nil, err
			}
		}
	}

	// If PEM parsing failed, try to decode as PKCS#12
	if len(certs) == 0 || key == nil {
		tflog.Debug(ctx, "Attempting to parse as PKCS#12")
		privateKey, certificate, err := pkcs12.Decode(certData, string(password))
		if err != nil {
			tflog.Error(ctx, "Failed to parse PKCS#12 data", map[string]interface{}{
				"error": err,
			})
			return nil, nil, err
		}

		certs = append(certs, certificate)
		key = privateKey
		certType = "PKCS#12"
	}

	if len(certs) == 0 {
		tflog.Error(ctx, "No certificates found")
		return nil, nil, errors.New("no certificates found")
	}
	if key == nil {
		tflog.Error(ctx, "No private key found")
		return nil, nil, errors.New("no private key found")
	}

	// Check that the private key is of the expected RSA type
	if _, ok := key.(*rsa.PrivateKey); !ok {
		tflog.Error(ctx, "Private key is not of RSA type")
		return nil, nil, errors.New("private key is not of RSA type")
	}

	tflog.Info(ctx, "Certificate and private key parsed successfully", map[string]interface{}{
		"certificateType": certType,
	})

	return certs, key, nil
}
