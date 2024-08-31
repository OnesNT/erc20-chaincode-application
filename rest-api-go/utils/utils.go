package utils

import (
	"crypto/x509"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/hyperledger/fabric-gateway/pkg/identity"
)

// LoadCertificate loads an X.509 certificate from a PEM file.
func LoadCertificate(filename string) (*x509.Certificate, error) {
	certificatePEM, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return identity.CertificateFromPEM(certificatePEM)
}

// GetCertificateAndPrivateKeyFromForm extracts certificate and private key from the form data.
func GetCertificateAndPrivateKeyFromForm(r *http.Request) ([]byte, []byte, error) {
	err := r.ParseMultipartForm(32 << 20) // 32 MB limit
	if err != nil {
		return nil, nil, fmt.Errorf("unable to parse form: %v", err)
	}

	certFile, _, err := r.FormFile("cert")
	if err != nil {
		return nil, nil, fmt.Errorf("unable to get cert file: %v", err)
	}
	defer certFile.Close()

	keyFile, _, err := r.FormFile("key")
	if err != nil {
		return nil, nil, fmt.Errorf("unable to get key file: %v", err)
	}
	defer keyFile.Close()

	certificatePEM, err := io.ReadAll(certFile)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to read cert file: %v", err)
	}

	privateKeyPEM, err := io.ReadAll(keyFile)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to read key file: %v", err)
	}

	return certificatePEM, privateKeyPEM, nil
}

// SavePEMToFile saves PEM data to a specified file.
func SavePEMToFile(filePath string, pemData []byte) error {
	// Create or truncate the file at the given path
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write the PEM data to the file
	_, err = file.Write(pemData)
	if err != nil {
		return err
	}

	return nil
}
