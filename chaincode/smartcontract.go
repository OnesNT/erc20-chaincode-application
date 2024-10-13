package chaincode

import (
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/v2/contractapi"
)

const chaincodeName = "test"

// BSmartContract provides functions for interacting with Chaincode A
type BSmartContract struct {
	contractapi.Contract
}

// GetOwnerNameFromCertificate extracts the owner name (Common Name) from the caller's X.509 certificate
func GetOwnerNameFromCertificate(ctx contractapi.TransactionContextInterface) (string, error) {
	// Get the client certificate (this will return the X.509 certificate string)
	certificate, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		return "", fmt.Errorf("failed to get client certificate: %v", err)
	}

	// Decode the certificate
	block, _ := pem.Decode([]byte(certificate))
	if block == nil {
		return "", fmt.Errorf("failed to parse certificate PEM")
	}

	// Parse the X.509 certificate
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return "", fmt.Errorf("failed to parse X.509 certificate: %v", err)
	}

	// Extract the common name (CN) from the certificate
	ownerName := cert.Subject.CommonName
	if ownerName == "" {
		return "", fmt.Errorf("common name (CN) not found in the certificate")
	}

	return ownerName, nil
}

// CreateBalanceForCaller creates a balance for the caller if it doesn't already exist
func (s *SmartContract) CreateBalanceForCaller(ctx contractapi.TransactionContextInterface, initialBalance float64) error {
    // Retrieve the caller's identity (client ID)
    callerID, err := ctx.GetClientIdentity().GetID()
    if err != nil {
        return fmt.Errorf("failed to get client identity: %v", err)
    }

    // Check if a balance already exists for the caller using IsBalanceExists
    balanceExists, err := IsBalanceExists(ctx, callerID)
    if err != nil {
        return fmt.Errorf("failed to check if balance exists: %v", err)
    }
    if balanceExists {
        return fmt.Errorf("balance for caller already exists")
    }

    // Get the owner name (Common Name) from the certificate
    ownerName, err := GetOwnerNameFromCertificate(ctx)
    if err != nil {
        return fmt.Errorf("failed to retrieve owner name: %v", err)
    }

    // Create a new balance entry for the caller
    balance := &Balance{
        ID:          callerID,              // Use client ID as the balance ID
        Owner:       ownerName,             // Set owner as the extracted common name
        Balance:     initialBalance,        // Set the initial balance
        LastUpdated: int(time.Now().Unix()), // Use current time as LastUpdated (Unix timestamp)
    }

    // Save the balance to the world state
    err = balance.Save(ctx)
    if err != nil {
        return fmt.Errorf("failed to create balance: %v", err)
    }

    return nil
}

// CallMintFromDebtPED invokes the MintFromDebt function of Chaincode A
func (s *BSmartContract) CallMintFromDebtPED(ctx contractapi.TransactionContextInterface, chaincodeAName string, amount_mint float64) error {
	args := [][]byte{[]byte("MintFromDebt"), []byte(fmt.Sprintf("%f", amount_mint))}

	response := ctx.GetStub().InvokeChaincode(chaincodeAName, args, "")
	if response.Status != 200 {
		return fmt.Errorf("failed to invoke Chaincode A: %s", response.Message)
	}

	return nil
}

// CallTransferFromChaincode invokes the TransferFrom function of Chaincode A
func (s *BSmartContract) CallTransferFromChaincode(ctx contractapi.TransactionContextInterface, chaincodeAName string, amount_transfer float64, to_ID string) error {
	args := [][]byte{[]byte("TransferFromChaincode"), []byte(fmt.Sprintf("%f", amount_transfer)), []byte(to_ID)}

	response := ctx.GetStub().InvokeChaincode(chaincodeAName, args, "")
	if response.Status != 200 {
		return fmt.Errorf("failed to invoke Chaincode A: %s", response.Message)
	}

	return nil
}

// CallTransferFromUser invokes the TransferFromUser function of Chaincode A
func (s *BSmartContract) CallTransferFromUser(ctx contractapi.TransactionContextInterface, chaincodeAName string, amount_transfer float64, to_ID string) error {
	args := [][]byte{[]byte("TransferFromUser"), []byte(fmt.Sprintf("%f", amount_transfer)), []byte(to_ID)}

	response := ctx.GetStub().InvokeChaincode(chaincodeAName, args, "")
	if response.Status != 200 {
		return fmt.Errorf("failed to invoke Chaincode A: %s", response.Message)
	}

	return nil
}

// CallPayOff invokes the PayOff function of Chaincode A
func (s *BSmartContract) CallPayOff(ctx contractapi.TransactionContextInterface, chaincodeAName string, amount_payoff float64) error {
	args := [][]byte{[]byte("PayOff"), []byte(fmt.Sprintf("%f", amount_payoff))}

	response := ctx.GetStub().InvokeChaincode(chaincodeAName, args, "")
	if response.Status != 200 {
		return fmt.Errorf("failed to invoke Chaincode A: %s", response.Message)
	}

	return nil
}
	
