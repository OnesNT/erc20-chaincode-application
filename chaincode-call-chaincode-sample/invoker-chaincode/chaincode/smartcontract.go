package chaincode

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/v2/contractapi"
)

const chaincodeName = "invoker"

// BSmartContract provides functions for interacting with Chaincode A
type BSmartContract struct {
	contractapi.Contract
}

// ReadAssetFromA invokes the ReadAsset function of Chaincode A to get an asset by its ID
func (s *BSmartContract) ReadAssetFromA(ctx contractapi.TransactionContextInterface, chaincodeAName, assetID string) (*Asset, error) {
	// Construct the arguments for invoking Chaincode A
	args := [][]byte{[]byte("ReadAsset"), []byte(assetID)}

	// Call Chaincode A's ReadAsset function using InvokeChaincode
	response := ctx.GetStub().InvokeChaincode(chaincodeAName, args, "")
	if response.Status != 200 {
		return nil, fmt.Errorf("failed to invoke Chaincode A: %s", response.Message)
	}

	// Unmarshal the asset from the response payload
	var asset Asset
	if err := json.Unmarshal(response.Payload, &asset); err != nil {
		return nil, fmt.Errorf("failed to unmarshal asset from Chaincode A: %v", err)
	}

	return &asset, nil
}

func (s *BSmartContract) ReadUserFromA(ctx contractapi.TransactionContextInterface, chaincodeAName, userID string) (*User, error) {
	// Construct the arguments for invoking Chaincode A
	args := [][]byte{[]byte("ReadUser"), []byte(userID)}

	// Call Chaincode A's ReadAsset function using InvokeChaincode
	response := ctx.GetStub().InvokeChaincode(chaincodeAName, args, "")
	if response.Status != 200 {
		return nil, fmt.Errorf("failed to invoke Chaincode A: %s", response.Message)
	}

	// Unmarshal the asset from the response payload
	var user User
	if err := json.Unmarshal(response.Payload, &user); err != nil {
		return nil, fmt.Errorf("failed to unmarshal user from Chaincode A: %v", err)
	}

	return &user, nil
}

// Asset describes basic details of what makes up a simple asset
type Asset struct {
	AppraisedValue int    `json:"appraised_value"`
	Color          string `json:"color"`
	ID             string `json:"id"`
	Owner          string `json:"owner"`
	Size           int    `json:"size"`
}

type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
	Sex  string `json:"sex"`
}
