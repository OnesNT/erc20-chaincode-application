package chaincode

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/hyperledger/fabric-contract-api-go/v2/contractapi"
)

// Asset describes basic details of what makes up a simple asset
type Asset struct {
	ID             string `json:"id"`
	Color          string `json:"color"`
	Owner          string `json:"owner"`
	Size           int    `json:"size"`
	AppraisedValue int    `json:"appraised_value"`
}

// getTableName returns the name of the struct type to use as the table name
func (a *Asset) getTableName() string {
	return reflect.TypeOf(a).Elem().Name() // Get the struct name, the value is "Asset"
}

// Save creates a new asset and stores it in the world state
func (a *Asset) Save(ctx contractapi.TransactionContextInterface) error {
	assetJSON, err := json.Marshal(a)
	if err != nil {
		return err
	}

	tableName := a.getTableName() // Automatically get the table name

	// Store the asset with the tableName prefix
	err = ctx.GetStub().PutState(tableName+"||"+a.ID, assetJSON)
	if err != nil {
		return fmt.Errorf("failed to put asset into world state: %v", err)
	}

	return nil
}

// ReadAsset retrieves an asset from the world state based on the asset ID
func ReadAsset(ctx contractapi.TransactionContextInterface, id string) (*Asset, error) {
	asset := &Asset{}
	tableName := asset.getTableName() // Automatically get the table name

	// Retrieve the asset using the tableName prefix
	assetJSON, err := ctx.GetStub().GetState(tableName + "||" + id)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if assetJSON == nil {
		return nil, fmt.Errorf("asset %s does not exist", id)
	}

	err = json.Unmarshal(assetJSON, asset)
	if err != nil {
		return nil, err
	}

	return asset, nil
}

// Update updates the asset in the world state
func (a *Asset) Update(ctx contractapi.TransactionContextInterface) error {
	assetJSON, err := json.Marshal(a)
	if err != nil {
		return err
	}

	tableName := a.getTableName() // Automatically get the table name

	// Update the asset with the tableName prefix
	err = ctx.GetStub().PutState(tableName+"||"+a.ID, assetJSON)
	if err != nil {
		return fmt.Errorf("failed to update asset in world state: %v", err)
	}

	return nil
}

// DeleteAsset deletes an asset from the world state based on its ID
func DeleteAsset(ctx contractapi.TransactionContextInterface, id string) error {
	asset := &Asset{}
	tableName := asset.getTableName() // Automatically get the table name

	// Check if asset exists using the tableName prefix
	assetExists, err := IsAssetExists(ctx, id)
	if err != nil {
		return err
	}
	if !assetExists {
		return fmt.Errorf("the asset %s does not exist", id)
	}

	// Delete the asset with the tableName prefix
	return ctx.GetStub().DelState(tableName + "||" + id)
}

// IsAssetExists checks if an asset with the given ID exists in the world state
func IsAssetExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	asset := &Asset{}
	tableName := asset.getTableName() // Automatically get the table name

	// Check if the asset exists using the tableName prefix
	assetJSON, err := ctx.GetStub().GetState(tableName + "||" + id)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}

	return assetJSON != nil, nil
}
