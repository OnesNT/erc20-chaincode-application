package chaincode

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/hyperledger/fabric-contract-api-go/v2/contractapi"
)

type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
	Sex  string `json:"sex"`
}

func (u *User) getTableName() string {
	return reflect.TypeOf(*u).Name()
}

func (u *User) Save(ctx contractapi.TransactionContextInterface) error {
	userJson, err := json.Marshal(u)
	if err != nil {
		return fmt.Errorf("failed to marshal user: %v", err)
	}

	tableName := u.getTableName()

	err = ctx.GetStub().PutState(tableName+"||"+u.ID, userJson)
	if err != nil {
		return fmt.Errorf("failed to put user into world state: %v", err)
	}

	return nil
}

func ReadUser(ctx contractapi.TransactionContextInterface, id string) (*User, error) {
	user := &User{}
	tableName := user.getTableName()

	userJson, err := ctx.GetStub().GetState(tableName + "||" + id)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}

	if userJson == nil {
		return nil, fmt.Errorf("user %s does not exist", id)
	}

	err = json.Unmarshal(userJson, user)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal user data: %v", err)
	}

	return user, nil
}

// Update updates the user in the world state
func (u *User) Update(ctx contractapi.TransactionContextInterface) error {
	exists, err := IsUserExists(ctx, u.ID)
	if err != nil {
		return fmt.Errorf("failed to check if user exists: %v", err)
	}
	if !exists {
		return fmt.Errorf("user %s does not exist", u.ID)
	}

	userJson, err := json.Marshal(u)
	if err != nil {
		return fmt.Errorf("failed to marshal user: %v", err)
	}

	tableName := u.getTableName()

	err = ctx.GetStub().PutState(tableName+"||"+u.ID, userJson)
	if err != nil {
		return fmt.Errorf("failed to update user in world state: %v", err)
	}

	return nil
}

func DeleteUser(ctx contractapi.TransactionContextInterface, id string) error {
	userExists, err := IsUserExists(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to check if user exists: %v", err)
	}
	if !userExists {
		return fmt.Errorf("the user %s does not exist", id)
	}

	user := &User{}
	tableName := user.getTableName()

	// Delete the user with the tableName prefix
	err = ctx.GetStub().DelState(tableName + "||" + id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %v", err)
	}

	return nil
}

func IsUserExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	user := &User{}
	tableName := user.getTableName()

	// Check if the user exists in the world state
	assetJSON, err := ctx.GetStub().GetState(tableName + "||" + id)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}

	return assetJSON != nil, nil
}
