package chaincode

import (
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/v2/contractapi"
)

type Payment struct {
    ID            string    `json:"id"`
    LoanID        string    `json:"loan_id"`
    PaymentAmount float64   `json:"payment_amount"`
    PaymentDate   time.Time `json:"payment_date"`
    CreatedAt     time.Time `json:"created_at"`
}

