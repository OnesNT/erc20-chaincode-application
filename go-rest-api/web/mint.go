package web

import (
	"fmt"
	"net/http"

	"github.com/hyperledger/fabric-gateway/pkg/client"
)

// Mint handles chaincode mint requests.
func (setup *OrgSetup) Mint(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	fmt.Println("Received Mint request")
	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "ParseForm() err: %s", err)
		return
	}

	// Extract chaincode-related information from the request
	chainCodeName := r.FormValue("chaincodeid") // Chaincode name
	channelID := r.FormValue("channelid")       // Channel ID
	amount := r.FormValue("amount")             // Amount to mint

	if chainCodeName == "" || channelID == "" || amount == "" {
		http.Error(w, "Missing required fields: chaincodeid, channelid, or amount", http.StatusBadRequest)
		return
	}

	fmt.Printf("Channel: %s, Chaincode: %s, Amount: %s\n", channelID, chainCodeName, amount)

	// Retrieve the network and contract from the gateway
	network := setup.Gateway.GetNetwork(channelID)
	contract := network.GetContract(chainCodeName)

	// Prepare the arguments for the Mint transaction
	args := []string{amount}

	// Create a transaction proposal for the 'Mint' function
	txn_proposal, err := contract.NewProposal("Mint", client.WithArguments(args...))
	if err != nil {
		fmt.Fprintf(w, "Error creating transaction proposal: %s", err)
		return
	}

	// Endorse the transaction proposal
	txn_endorsed, err := txn_proposal.Endorse()
	if err != nil {
		fmt.Fprintf(w, "Error endorsing transaction: %s", err)
		return
	}

	// Submit the endorsed transaction
	txn_committed, err := txn_endorsed.Submit()
	if err != nil {
		fmt.Fprintf(w, "Error submitting transaction: %s", err)
		return
	}

	// Respond with the transaction details
	fmt.Fprintf(w, "Mint successful. Transaction ID: %s, Response: %s", txn_committed.TransactionID(), txn_endorsed.Result())
}
