package helpers

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"math/big"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	pkcs12 "software.sslmate.com/src/go-pkcs12"
)

func TestParseCertificateData(t *testing.T) {
	ctx := context.Background()

	t.Run("Valid PFX with password", func(t *testing.T) {
		pfxData, password, err := generatePFXWithPassword()
		require.NoError(t, err)

		certs, privKey, err := ParseCertificateData(ctx, pfxData, []byte(password))
		assert.NoError(t, err)
		assert.Len(t, certs, 1)
		assert.NotNil(t, privKey)
		_, ok := privKey.(*rsa.PrivateKey)
		assert.True(t, ok)
	})

	t.Run("Valid PFX without password", func(t *testing.T) {
		pfxData, err := generatePFXWithoutPassword()
		require.NoError(t, err)

		certs, privKey, err := ParseCertificateData(ctx, pfxData, []byte(""))
		assert.NoError(t, err)
		assert.Len(t, certs, 1)
		assert.NotNil(t, privKey)
		_, ok := privKey.(*rsa.PrivateKey)
		assert.True(t, ok)
	})

	t.Run("PFX with non-RSA key", func(t *testing.T) {
		pfxData, password, err := generatePFXWithNonRSAKey()
		require.NoError(t, err)

		_, _, err = ParseCertificateData(ctx, pfxData, []byte(password))
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "private key is not of RSA type")
	})

	t.Run("Invalid PFX data", func(t *testing.T) {
		_, _, err := ParseCertificateData(ctx, []byte("invalid data"), []byte("password"))
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to parse PKCS#12 data")
	})

	t.Run("Incorrect password", func(t *testing.T) {
		pfxData, _, err := generatePFXWithPassword()
		require.NoError(t, err)

		_, _, err = ParseCertificateData(ctx, pfxData, []byte("wrongpassword"))
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to parse PKCS#12 data")
	})
}

// generatePFXWithPassword creates a PKCS#12 (PFX) certificate with an RSA private key and password
func generatePFXWithPassword() (pfxData []byte, password string, err error) {
	// Generate RSA private key
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, "", err
	}

	// Create certificate template
	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			CommonName: "Test Cert",
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(24 * time.Hour),
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	// Create certificate
	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &privateKey.PublicKey, privateKey)
	if err != nil {
		return nil, "", err
	}

	cert, err := x509.ParseCertificate(certDER)
	if err != nil {
		return nil, "", err
	}

	// Encode to PKCS#12
	password = "testpassword"
	pfxData, err = pkcs12.Modern.Encode(privateKey, cert, nil, password)
	if err != nil {
		return nil, "", err
	}

	return pfxData, password, nil
}

// generatePFXWithoutPassword creates a PKCS#12 (PFX) certificate with an RSA private key and no password
func generatePFXWithoutPassword() (pfxData []byte, err error) {
	// Generate RSA private key
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}

	// Create certificate template
	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			CommonName: "Test Cert No Password",
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(24 * time.Hour),
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	// Create certificate
	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &privateKey.PublicKey, privateKey)
	if err != nil {
		return nil, err
	}

	cert, err := x509.ParseCertificate(certDER)
	if err != nil {
		return nil, err
	}

	// Encode to PKCS#12 without password
	pfxData, err = pkcs12.Modern.Encode(privateKey, cert, nil, "")
	if err != nil {
		return nil, err
	}

	return pfxData, nil
}

// generatePFXWithNonRSAKey creates a PKCS#12 (PFX) certificate with an ECDSA private key (non-RSA)
func generatePFXWithNonRSAKey() (pfxData []byte, password string, err error) {
	// Generate ECDSA private key
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, "", err
	}

	// Create certificate template
	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			CommonName: "Test Cert ECDSA",
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(24 * time.Hour),
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	// Create certificate
	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &privateKey.PublicKey, privateKey)
	if err != nil {
		return nil, "", err
	}

	cert, err := x509.ParseCertificate(certDER)
	if err != nil {
		return nil, "", err
	}

	// Encode to PKCS#12
	password = "testpassword"
	pfxData, err = pkcs12.Modern.Encode(privateKey, cert, nil, password)
	if err != nil {
		return nil, "", err
	}

	return pfxData, password, nil
}
