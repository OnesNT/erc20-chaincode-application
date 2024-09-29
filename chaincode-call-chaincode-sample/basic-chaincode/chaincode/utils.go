package chaincode

import (
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/v2/contractapi"
	"github.com/hyperledger/fabric-protos-go-apiv2/common"
	"github.com/hyperledger/fabric-protos-go-apiv2/peer"
	"google.golang.org/protobuf/proto"
)

// getInvokerChaincodeName extracts the invoking chaincode name from the context
func getInvokerChaincodeName(ctx contractapi.TransactionContextInterface) (string, error) {
	// Extract the invoking chaincode name from the context
	stub := ctx.GetStub()
	signedProposal, err := stub.GetSignedProposal()
	if err != nil {
		return "", fmt.Errorf("failed to get signed proposal: %v", err)
	}
	// fmt.Printf("Debug ProposalBytes: %s\n", signedProposal.ProposalBytes)

	proposal := &peer.Proposal{}
	err = proto.Unmarshal(signedProposal.ProposalBytes, proposal)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal proposal: %v", err)
	}
	// fmt.Printf("Debug proposal: %s\n", proposal)

	// Unmarshal Header from Proposal
	header := &common.Header{}
	err = proto.Unmarshal(proposal.Header, header)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal header: %v", err)
	}
	// fmt.Printf("Debug Header: %+v\n", header)

	// Unmarshal ChannelHeader from Header
	channelHeader := &common.ChannelHeader{}
	err = proto.Unmarshal(header.ChannelHeader, channelHeader)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal channel header: %v", err)
	}
	// fmt.Printf("Debug ChannelHeader: %+v\n", channelHeader)

	// Unmarshal ChaincodeHeaderExtension from ChannelHeader Extension
	chaincodeHeaderExtension := &peer.ChaincodeHeaderExtension{}
	err = proto.Unmarshal(channelHeader.Extension, chaincodeHeaderExtension)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal chaincode header extension: %v", err)
	}
	// fmt.Printf("Invoked by chaincode: %s\n", chaincodeHeaderExtension.ChaincodeId.Name)

	// payload := &peer.ChaincodeProposalPayload{}
	// err = proto.Unmarshal(proposal.Payload, payload)
	// if err != nil {
	// 	return "", fmt.Errorf("failed to unmarshal header: %v", err)
	// }
	// fmt.Printf("Debug payload: %s\n", payload.GetInput())

	return chaincodeHeaderExtension.ChaincodeId.Name, nil
}
