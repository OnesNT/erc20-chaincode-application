package services

import (
	"fmt"
	"net/http"
	"path/filepath"
	"rest-api-go/utils"
	"strconv"

	"github.com/hyperledger/fabric-gateway/pkg/client"
)

// GatewayService provides access to the blockchain network and contracts.
type GatewayService struct {
	Setup *OrgSetup
}

// NewGatewayService creates a new GatewayService instance.
func NewGatewayService(setup *OrgSetup) *GatewayService {
	return &GatewayService{Setup: setup}
}

// GetNetwork gets a network from the gateway.
func (g *GatewayService) GetNetwork(channelID string) *client.Network {
	return g.Setup.Gateway.GetNetwork(channelID)
}

// GetContract gets a contract from the network.
func (g *GatewayService) GetContract(chainCodeName string) *client.Contract {
	network := g.GetNetwork("default")
	return network.GetContract(chainCodeName)
}

// CallChaincode calls a chaincode function with specified arguments.
func (g *GatewayService) CallChaincode(channelID, chainCodeName, functionChaincode string, args []string) (string, error) {
	// Retrieve the network and contract
	network := g.GetNetwork(channelID)
	if network == nil {
		return "", fmt.Errorf("network %s does not exist", channelID)
	}
	contract := network.GetContract(chainCodeName)
	if contract == nil {
		return "", fmt.Errorf("contract %s does not exist", chainCodeName)
	}

	// Call the specified function on the chaincode
	txn_proposal, err := contract.NewProposal(functionChaincode, client.WithArguments(args...))
	if err != nil {
		return "", fmt.Errorf("error creating transaction proposal: %v", err)
	}

	// Endorse the transaction proposal
	txn_endorsed, err := txn_proposal.Endorse()
	if err != nil {
		return "", fmt.Errorf("error endorsing transaction: %v", err)
	}

	// Submit the endorsed transaction
	txn_committed, err := txn_endorsed.Submit()
	if err != nil {
		return "", fmt.Errorf("error submitting transaction: %v", err)
	}

	return txn_committed.TransactionID(), nil
}

// CallChaincodeGET queries the chaincode and returns the result as an integer.
func (g *GatewayService) CallChaincodeGET(channelID, chainCodeName, functionChaincode string) (int, error) {
	// Retrieve the network and contract
	network := g.GetNetwork(channelID)
	if network == nil {
		return 0, fmt.Errorf("network %s does not exist", channelID)
	}
	contract := network.GetContract(chainCodeName)
	if contract == nil {
		return 0, fmt.Errorf("contract %s does not exist", chainCodeName)
	}

	// Evaluate the transaction
	result, err := contract.EvaluateTransaction(functionChaincode)
	if err != nil {
		return 0, fmt.Errorf("failed to evaluate transaction: %v", err)
	}

	// Convert result to integer
	balance, err := strconv.Atoi(string(result))
	if err != nil {
		return 0, fmt.Errorf("failed to convert result to integer: %v", err)
	}

	return balance, nil
}

// InitializeWithCertAndKey handles extracting the cert and key from the request and reinitializes the service.
func (s *GatewayService) InitializeWithCertAndKey(r *http.Request) error {

	// Extract certificate and key from the request
	certPEM, keyPEM, err := utils.GetCertificateAndPrivateKeyFromForm(r)
	if err != nil {
		return fmt.Errorf("failed to get certificate and key: %w", err)
	}

	// Save the cert and key to temporary files
	certPath := filepath.Join("/tmp", "cert.pem")
	keyPath := filepath.Join("/tmp", "key.pem")

	err = utils.SavePEMToFile(certPath, certPEM)
	if err != nil {
		return fmt.Errorf("failed to save certificate: %w", err)
	}

	err = utils.SavePEMToFile(keyPath, keyPEM)
	if err != nil {
		return fmt.Errorf("failed to save key: %w", err)
	}

	// Define orgconfig
	cryptoPath := "../../test-network/organizations/peerOrganizations/org1.example.com"
	orgConfig := OrgSetup{
		OrgName:      "Org1",
		MSPID:        "Org1MSP",
		CertPath:     certPath,
		KeyPath:      keyPath,
		TLSCertPath:  cryptoPath + "/peers/peer0.org1.example.com/tls/ca.crt",
		PeerEndpoint: "dns:///localhost:7051",
		GatewayPeer:  "peer0.org1.example.com",
	}

	// Update the OrgSetup configuration dynamically
	s.Setup = &orgConfig

	// Reinitialize the service with the new certificate and key
	*s = *NewGatewayService(s.Setup)

	return nil
}
