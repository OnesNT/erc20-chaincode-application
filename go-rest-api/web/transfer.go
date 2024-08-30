package web

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/hyperledger/fabric-gateway/pkg/client"
)

// Transfer handles chaincode invoke requests for transferring an asset.
func (setup *OrgSetup) Transfer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	fmt.Println("Received Transfer request")
	if err := r.ParseForm(); err != nil {
		http.Error(w, fmt.Sprintf("ParseForm() err: %s", err), http.StatusBadRequest)
		return
	}

	chainCodeName := r.FormValue("chaincodeid")
	channelID := r.FormValue("channelid")
	amount := r.FormValue("amount")
	recipientCN := r.FormValue("recipientCN")

	if chainCodeName == "" || channelID == "" || amount == "" || recipientCN == "" {
		http.Error(w, "Missing required form data", http.StatusBadRequest)
		return
	}

	// Construct recipient identity string
	recipient := fmt.Sprintf("x509::CN=%s,OU=client,O=Hyperledger,ST=North Carolina,C=US::CN=ca.org1.example.com,O=org1.example.com,L=Durham,ST=North Carolina,C=US", recipientCN)
	recipient = strings.TrimSpace(recipient)

	// fmt.Printf("channel: %s, chaincode: %s, function: Transfer, args: [%d, %s]\n", channelID, chainCodeName, amount, recipient)

	// Access the network and contract
	network := setup.Gateway.GetNetwork(channelID)
	if network == nil {
		http.Error(w, fmt.Sprintf("Channel ID %s does not exist or cannot be accessed", channelID), http.StatusBadRequest)
		return
	}

	contract := network.GetContract(chainCodeName)

	// Create and submit the transaction proposal for transferring the asset
	txn_proposal, err := contract.NewProposal("Transfer", client.WithArguments(recipient, amount))
	if err != nil {
		http.Error(w, fmt.Sprintf("Error creating txn proposal: %s", err), http.StatusInternalServerError)
		return
	}

	txn_endorsed, err := txn_proposal.Endorse()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error endorsing txn: %s", err), http.StatusInternalServerError)
		return
	}

	txn_committed, err := txn_endorsed.Submit()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error submitting transaction: %s", err), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Transaction ID : %s Response: %s", txn_committed.TransactionID(), txn_endorsed.Result())
}
