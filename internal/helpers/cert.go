package helpers

import (
	"crypto/x509"
	"encoding/base64"
	"errors"
	"io"
	"os"

	pkcs12 "software.sslmate.com/src/go-pkcs12"
)

// GetCertificatesAndKeyFromCertOrFilePath takes either a base64-encoded certificate or a file path to a PKCS#12 file,
// decodes it, and returns the certificates and private key.
func GetCertificatesAndKeyFromCertOrFilePath(certOrFilePath string, password string) ([]*x509.Certificate, interface{}, error) {
	certData, err := base64.StdEncoding.DecodeString(certOrFilePath)
	if err == nil {
		key, cert, err := pkcs12.Decode(certData, password)
		if err == nil {
			return []*x509.Certificate{cert}, key, nil
		}
	}

	file, err := os.Open(certOrFilePath)
	if err != nil {
		return nil, nil, errors.New("could not open file or decode base64 input")
	}
	defer file.Close()

	pfxData, err := io.ReadAll(file)
	if err != nil {
		return nil, nil, errors.New("could not read file content")
	}

	key, cert, err := pkcs12.Decode(pfxData, password)
	if err != nil {
		return nil, nil, err
	}

	return []*x509.Certificate{cert}, key, nil
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
