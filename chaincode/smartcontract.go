package chaincode

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/v2/contractapi"
)

const chaincodeName = "test"

// BSmartContract provides functions for interacting with Chaincode A
type BSmartContract struct {
	contractapi.Contract
}


func (s *BSmartContract) InitLedger(ctx contractapi.TransactionContextInterface, chaincodeAName string) {
	blance := []Balance{
		{ID:"chaincode1", Owner: "Chill", Balance: 1000, LastUpdated: 23/12/2023}
		{ID:"chaincode2", Owner: "Bee", Balance : 2000, LastUpdated: 20/11/2024}
		{ID:"user1", Owner:"Quang", Balance: 1000, LastUpdated: 11/11/2021}
		{ID:"user2", Owner:"Huy", Balance: 2000, LastUpdated: 20/2/2022}

	}
	loan := []Loan{
		{ID: "loan1", Owner: "Chill", PrincipalAmount: 200.23, InterestRate: 6, StartDate: 1/1/2020, EndDate: 1/1/2022, RemainingPrincipal: 0, LoanStatus:"closed"},
		{ID: "loan2", Owner: "Chill", PrincipalAmount: 100, InterestRate: 6, StartDate: 1/3/2023, EndDate: 1/3/2025, RemainingPrincipal: 90, LoanStatus:"openning"},
		{ID: "loan3", Owner: "Chill", PrincipalAmount: 105, InterestRate: 6, StartDate: 1/5/2023, EndDate: 1/1/2025, RemainingPrincipal: 105, LoanStatus:"openning"},
	}
	payment := []Payment{
		{ID: "payment1", LoanID: "loan1", PaymentAmount: 200.23,PaymentDate: 1/1/2020, CreatedAt: 1/1/2020}
		{ID: "payment2", LoanID: "loan2", PaymentAmount: 10,PaymentDate: 10/10/2023, CreatedAt: 10/10/2023}
	}
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
	
