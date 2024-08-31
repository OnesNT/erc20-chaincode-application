package utils

import (
	"crypto/x509"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
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

func getCertificateAndPrivateKeyFromForm(c *gin.Context) ([]byte, []byte, error) {
	certFile, certErr := c.FormFile("cert")
	keyFile, keyErr := c.FormFile("key")
	if certErr != nil || keyErr != nil {
		return nil, nil, fmt.Errorf("cert and key files are required")
	}

	certContent, certErr := certFile.Open()
	if certErr != nil {
		return nil, nil, fmt.Errorf("unable to open cert file")
	}
	defer certContent.Close()

	keyContent, keyErr := keyFile.Open()
	if keyErr != nil {
		return nil, nil, fmt.Errorf("unable to open key file")
	}
	defer keyContent.Close()

	certificatePEM := make([]byte, certFile.Size)
	_, certErr = certContent.Read(certificatePEM)
	if certErr != nil {
		return nil, nil, fmt.Errorf("unable to read cert file")
	}

	privateKeyPEM := make([]byte, keyFile.Size)
	_, keyErr = keyContent.Read(privateKeyPEM)
	if keyErr != nil {
		return nil, nil, fmt.Errorf("unable to read key file")
	}

	return certificatePEM, privateKeyPEM, nil
}
