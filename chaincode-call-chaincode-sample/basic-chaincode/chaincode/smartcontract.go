package chaincode

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/v2/contractapi"
)

const chaincodeName = "basic"

// SmartContract provides functions for managing an Asset
type SmartContract struct {
	contractapi.Contract
}

// InitLedger adds a base set of assets to the ledger
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	assets := []Asset{
		{ID: "asset1", Color: "blue", Size: 5, Owner: "Tomoko", AppraisedValue: 300},
		{ID: "asset2", Color: "red", Size: 5, Owner: "Brad", AppraisedValue: 400},
		{ID: "asset3", Color: "green", Size: 10, Owner: "Jin Soo", AppraisedValue: 500},
		{ID: "asset4", Color: "yellow", Size: 10, Owner: "Max", AppraisedValue: 600},
		{ID: "asset5", Color: "black", Size: 15, Owner: "Adriana", AppraisedValue: 700},
		{ID: "asset6", Color: "white", Size: 15, Owner: "Michel", AppraisedValue: 800},
	}

	users := []User{
		{ID: "user1", Name: "Quang", Age: 22, Sex: "Male"},
		{ID: "user2", Name: "Huy", Age: 30, Sex: "Male"},
		{ID: "user3", Name: "Teo", Age: 21, Sex: "Male"},
		{ID: "user4", Name: "Thuy", Age: 18, Sex: "Female"},
		{ID: "user5", Name: "Ha", Age: 12, Sex: "Female"},
		{ID: "user6", Name: "Hue", Age: 29, Sex: "Female"},
	}

	for _, asset := range assets {
		err := asset.Save(ctx)
		if err != nil {
			return fmt.Errorf("failed to put asset into world state: %v", err)
		}
	}

	for _, user := range users {
		err := user.Save(ctx)
		if err != nil {
			return fmt.Errorf("failed to put user into world state: %v", err)
		}
	}

	return nil
}

// CreateAsset issues a new asset to the world state with given details.
func (s *SmartContract) CreateAsset(ctx contractapi.TransactionContextInterface, id string, color string, size int, owner string, appraisedValue int) error {
	asset := Asset{
		ID:             id,
		Color:          color,
		Size:           size,
		Owner:          owner,
		AppraisedValue: appraisedValue,
	}

	exists, err := IsAssetExists(ctx, id)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("the asset %s already exists", id)
	}

	return asset.Save(ctx)
}

func (s *SmartContract) CreateUser(ctx contractapi.TransactionContextInterface, id string, name string, age int, sex string) error {
	user := User{
		ID:   id,
		Name: name,
		Age:  age,
		Sex:  sex,
	}

	exists, err := IsUserExists(ctx, id)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("the user %s already exists", id)
	}

	return user.Save(ctx)
}

// ReadAsset returns the asset stored in the world state with given id.
func (s *SmartContract) ReadAsset(ctx contractapi.TransactionContextInterface, id string) (*Asset, error) {
	// Extract the invoking chaincode name
	ccName, err := getInvokerChaincodeName(ctx)
	if err != nil {
		return nil, err
	}
	if ccName != chaincodeName {
		fmt.Printf("Invoked by chaincode: %s\n", ccName)
	}

	asset, err := ReadAsset(ctx, id)
	if err != nil {
		return nil, err
	}

	return asset, nil
}

func (s *SmartContract) ReadUser(ctx contractapi.TransactionContextInterface, id string) (*User, error) {
	// Extract the invoking chaincode name
	ccName, err := getInvokerChaincodeName(ctx)
	if err != nil {
		return nil, err
	}
	if ccName != chaincodeName {
		fmt.Printf("Invoked by chaincode: %s\n", ccName)
	}

	user, err := ReadUser(ctx, id)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// UpdateAsset updates an existing asset in the world state with provided parameters.
func (s *SmartContract) UpdateAsset(ctx contractapi.TransactionContextInterface, id string, color string, size int, owner string, appraisedValue int) error {
	asset := Asset{
		ID:             id,
		Color:          color,
		Size:           size,
		Owner:          owner,
		AppraisedValue: appraisedValue,
	}

	exists, err := IsAssetExists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the asset %s does not exist", id)
	}

	return asset.Save(ctx)
}

func (s *SmartContract) UpdateUser(ctx contractapi.TransactionContextInterface, id string, name string, age int, sex string) error {
	user := User{
		ID:   id,
		Name: name,
		Age:  age,
		Sex:  sex,
	}

	exists, err := IsUserExists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the user %s does not exist", id)
	}

	return user.Save(ctx)
}

// DeleteAsset deletes an asset from the world state.
func (s *SmartContract) DeleteAsset(ctx contractapi.TransactionContextInterface, id string) error {
	return DeleteAsset(ctx, id)
}

func (s *SmartContract) DeleteUser(ctx contractapi.TransactionContextInterface, id string) error {
	return DeleteUser(ctx, id)
}

// TransferAsset updates the owner field of an asset with the given id in the world state, and returns the old owner.
func (s *SmartContract) TransferAsset(ctx contractapi.TransactionContextInterface, id string, newOwner string) (string, error) {
	asset, err := ReadAsset(ctx, id)
	if err != nil {
		return "", err
	}

	oldOwner := asset.Owner
	asset.Owner = newOwner

	err = asset.Save(ctx)
	if err != nil {
		return "", err
	}

	return oldOwner, nil
}

// // GetAllAssets returns all assets found in the world state
// func (s *SmartContract) GetAllAssets(ctx contractapi.TransactionContextInterface) ([]*Asset, error) {
// 	// range query with empty string for startKey and endKey does an
// 	// open-ended query of all assets in the chaincode namespace.
// 	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer resultsIterator.Close()

// 	var assets []*Asset
// 	for resultsIterator.HasNext() {
// 		queryResponse, err := resultsIterator.Next()
// 		if err != nil {
// 			return nil, err
// 		}

// 		var asset Asset
// 		err = json.Unmarshal(queryResponse.Value, &asset)
// 		if err != nil {
// 			return nil, err
// 		}
// 		assets = append(assets, &asset)
// 	}

// 	return assets, nil
// }

func (s *SmartContract) GetAllAssets(ctx contractapi.TransactionContextInterface) ([]*Asset, error) {
	// Query only keys that start with "Asset||" to filter out non-asset records
	resultsIterator, err := ctx.GetStub().GetStateByRange("Asset||", "Asset||\ufff0")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var assets []*Asset
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var asset Asset
		err = json.Unmarshal(queryResponse.Value, &asset)
		if err != nil {
			return nil, err
		}
		assets = append(assets, &asset)
	}

	return assets, nil
}

// GetAllAssets returns all assets found in the world state
func (s *SmartContract) GetAllUsers(ctx contractapi.TransactionContextInterface) ([]*User, error) {
	// range query with empty string for startKey and endKey does an
	// open-ended query of all assets in the chaincode namespace.
	resultsIterator, err := ctx.GetStub().GetStateByRange("User||", "User||\ufff0")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var users []*User
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var user User
		err = json.Unmarshal(queryResponse.Value, &user)
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	return users, nil
}
