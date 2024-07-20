package helpers

import (
	"crypto/x509"
	"encoding/base64"
	"errors"
	"io"
	"os"

	pkcs12 "software.sslmate.com/src/go-pkcs12"
)

// GetRawCertificateFromCertOrFilePath takes either a DER-encoded certificate
// or a file path to a DER-encoded PKCS#12 file, decodes it, and returns the raw certificate.
func GetRawCertificateFromCertOrFilePath(certOrFilePath string, password string) (*x509.Certificate, error) {
	certData, err := base64.StdEncoding.DecodeString(certOrFilePath)
	if err == nil {
		cert, err := x509.ParseCertificate(certData)
		if err == nil {
			return cert, nil
		}
	}

	file, err := os.Open(certOrFilePath)
	if err != nil {
		return nil, errors.New("could not open file or decode base64 input")
	}
	defer file.Close()

	pfxData, err := io.ReadAll(file)
	if err != nil {
		return nil, errors.New("could not read file content")
	}

	_, cert, err := pkcs12.Decode(pfxData, password)
	if err != nil {
		return nil, err
	}

	return cert, nil
}

// ConvertBase64ToCert takes a base64 encoded PKCS#12 file, decodes it, and returns the certificate.
func ConvertBase64ToCert(base64PfxData string, password string) (*x509.Certificate, error) {
	pfxData, err := base64.StdEncoding.DecodeString(base64PfxData)
	if err != nil {
		return nil, err
	}

	_, cert, err := pkcs12.Decode(pfxData, password)
	if err != nil {
		return nil, err
	}

	return cert, nil
}
