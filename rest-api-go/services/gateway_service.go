package services

import "github.com/hyperledger/fabric-gateway/pkg/client"

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
	network := g.GetNetwork("default") // Replace with appropriate channel ID if needed
	return network.GetContract(chainCodeName)
}
