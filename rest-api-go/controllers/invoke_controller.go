package controllers

import (
	"fmt"
	"net/http"
	"rest-api-go/services"

	"github.com/hyperledger/fabric-gateway/pkg/client"
)

// InvokeController handles requests for invoke token.
type InvokeController struct {
	Service *services.GatewayService
}

// NewTransferController creates a new TransferController instance.
func NewInvokeController(setup *services.OrgSetup) *InvokeController {
	return &InvokeController{Service: services.NewGatewayService(setup)}
}

// InitializeContract handles initializing the chaincode with token information
func (c *InvokeController) InitializeContract(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse the form values for initialization parameters
	if err := r.ParseForm(); err != nil {
		http.Error(w, fmt.Sprintf("Failed to parse form: %v", err), http.StatusBadRequest)
		return
	}

	chainCodeName := r.FormValue("chaincodeid")
	channelID := r.FormValue("channelid")
	tokenName := r.FormValue("name")
	symbol := r.FormValue("symbol")
	decimals := r.FormValue("decimals")

	if chainCodeName == "" || channelID == "" || tokenName == "" || symbol == "" || decimals == "" {
		http.Error(w, "Missing required fields: chaincodeid, channelid, name, symbol, or decimals", http.StatusBadRequest)
		return
	}

	fmt.Printf("Initializing contract on channel: %s, chaincode: %s, with token name: %s, symbol: %s, decimals: %s\n",
		channelID, chainCodeName, tokenName, symbol, decimals)

	// Retrieve the network and contract
	network := c.Service.GetNetwork(channelID)
	if network == nil {
		http.Error(w, fmt.Sprintf("Network %s does not exist", channelID), http.StatusBadRequest)
		return
	}
	contract := network.GetContract(chainCodeName)
	if contract == nil {
		http.Error(w, fmt.Sprintf("Contract %s does not exist", chainCodeName), http.StatusBadRequest)
		return
	}

	// Call the Initialize function on the chaincode
	txn_proposal, err := contract.NewProposal("Initialize", client.WithArguments(tokenName, symbol, decimals))
	if err != nil {
		http.Error(w, fmt.Sprintf("Error creating transaction proposal: %v", err), http.StatusInternalServerError)
		return
	}

	// Endorse the transaction proposal
	txn_endorsed, err := txn_proposal.Endorse()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error endorsing transaction: %v", err), http.StatusInternalServerError)
		return
	}

	// Submit the endorsed transaction
	txn_committed, err := txn_endorsed.Submit()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error submitting transaction: %v", err), http.StatusInternalServerError)
		return
	}

	// Respond with success message
	fmt.Fprintf(w, "Initialization successful. Transaction ID: %s", txn_committed.TransactionID())
}
