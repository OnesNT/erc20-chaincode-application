package web

import (
	"crypto/x509"
	"fmt"
	"log"
	"os"
	"path"
	"time"

	"github.com/hyperledger/fabric-gateway/pkg/client"
	"github.com/hyperledger/fabric-gateway/pkg/identity"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// Initialize the setup for the organization.
func Initialize(setup OrgSetup) (*OrgSetup, error) {
	log.Printf("Initializing connection for %s...\n", setup.OrgName)
	clientConnection := setup.newGrpcConnection()
	id := setup.newIdentity()
	sign := setup.newSign()

	gateway, err := client.Connect(
		id,
		client.WithSign(sign),
		client.WithClientConnection(clientConnection),
		client.WithEvaluateTimeout(5*time.Second),
		client.WithEndorseTimeout(15*time.Second),
		client.WithSubmitTimeout(5*time.Second),
		client.WithCommitStatusTimeout(1*time.Minute),
	)
	if err != nil {
		panic(err)
	}
	setup.Gateway = *gateway
	log.Println("Initialization complete")
	return &setup, nil
}

// newGrpcConnection creates a gRPC connection to the Gateway server.
func (setup OrgSetup) newGrpcConnection() *grpc.ClientConn {
	certificate, err := loadCertificate(setup.TLSCertPath)
	if err != nil {
		panic(err)
	}

	certPool := x509.NewCertPool()
	certPool.AddCert(certificate)
	transportCredentials := credentials.NewClientTLSFromCert(certPool, setup.GatewayPeer)

	connection, err := grpc.NewClient(setup.PeerEndpoint, grpc.WithTransportCredentials(transportCredentials))
	if err != nil {
		panic(fmt.Errorf("failed to create gRPC connection: %w", err))
	}

	return connection
}

// newIdentity creates a client identity for this Gateway connection using an X.509 certificate.
func (setup OrgSetup) newIdentity() *identity.X509Identity {
	certificate, err := loadCertificate(setup.CertPath)
	if err != nil {
		panic(err)
	}

	// Print the certificate content for debugging
	fmt.Printf("Certificate PEM:\n%v\n", certificate)

	id, err := identity.NewX509Identity(setup.MSPID, certificate)
	if err != nil {
		panic(err)
	}

	return id
}

// newSign creates a function that generates a digital signature from a message digest using a private key.
func (setup OrgSetup) newSign() identity.Sign {
	files, err := os.ReadDir(setup.KeyPath)
	if err != nil {
		panic(fmt.Errorf("failed to read private key directory: %w", err))
	}
	privateKeyPEM, err := os.ReadFile(path.Join(setup.KeyPath, files[0].Name()))

	if err != nil {
		panic(fmt.Errorf("failed to read private key file: %w", err))
	}

	// Print the certificate content for debugging
	fmt.Printf("privateKeyPEM:\n%v\n", privateKeyPEM)

	privateKey, err := identity.PrivateKeyFromPEM(privateKeyPEM)
	if err != nil {
		panic(err)
	}

	sign, err := identity.NewPrivateKeySign(privateKey)
	if err != nil {
		panic(err)
	}

	return sign
}

// func (setup OrgSetup) newSign() identity.Sign {
// 	files, err := os.ReadDir(setup.KeyPath)
// 	if err != nil {
// 		panic(fmt.Errorf("failed to read private key directory: %w", err))
// 	}

// 	// Assuming the private key is stored in a PEM file
// 	privateKeyPEM, err := os.ReadFile(path.Join(setup.KeyPath, files[0].Name()))
// 	if err != nil {
// 		panic(fmt.Errorf("failed to read private key file: %w", err))
// 	}

// 	// Decode the PEM block
// 	block, _ := pem.Decode(privateKeyPEM)
// 	if block == nil {
// 		panic(fmt.Errorf("failed to parse PEM block"))
// 	}

// 	// Try to parse the private key as an EC private key
// 	privateKey, err := x509.ParseECPrivateKey(block.Bytes)
// 	if err != nil {
// 		panic(fmt.Errorf("failed to parse EC private key: %w", err))
// 	}

// 	// Use the private key to create a signing function
// 	sign, err := identity.NewPrivateKeySign(privateKey)
// 	if err != nil {
// 		panic(fmt.Errorf("failed to create a sign function from the private key: %w", err))
// 	}

// 	return sign
// }

func loadCertificate(filename string) (*x509.Certificate, error) {
	certificatePEM, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read certificate file: %w", err)
	}
	return identity.CertificateFromPEM(certificatePEM)
}
