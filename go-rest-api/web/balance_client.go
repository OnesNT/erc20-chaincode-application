package web

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// GetClientAccountBalance handles requests to query the balance of the requesting client's account.
func (setup *OrgSetup) GetClientAccountBalance(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse form values if needed, or set up static values if not using form data
	chainCodeName := r.URL.Query().Get("chaincodeid") // Chaincode name
	channelID := r.URL.Query().Get("channelid")       // Channel ID

	if chainCodeName == "" || channelID == "" {
		http.Error(w, "Missing required fields: chaincodeid or channelid", http.StatusBadRequest)
		return
	}

	// Retrieve the network and contract from the gateway
	network := setup.Gateway.GetNetwork(channelID)
	contract := network.GetContract(chainCodeName)

	// Evaluate the transaction to get the client account balance
	result, err := contract.EvaluateTransaction("ClientAccountBalance")
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to evaluate transaction: %v", err), http.StatusInternalServerError)
		return
	}

	// Send the result back as a JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"balance": string(result),
	})
}
