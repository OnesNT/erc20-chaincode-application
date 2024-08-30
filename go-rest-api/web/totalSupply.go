package web

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// QueryTotalSupply handles the request to query the total supply of tokens
func (setup *OrgSetup) QueryTotalSupply(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse form data
	if err := r.ParseForm(); err != nil {
		http.Error(w, fmt.Sprintf("ParseForm() err: %v", err), http.StatusBadRequest)
		return
	}

	chainCodeName := r.FormValue("chaincodeid")
	channelID := r.FormValue("channelid")

	if chainCodeName == "" || channelID == "" {
		http.Error(w, "Missing required fields: chaincodeid or channelid", http.StatusBadRequest)
		return
	}

	// Retrieve the network and contract from the gateway
	network := setup.Gateway.GetNetwork(channelID)
	contract := network.GetContract(chainCodeName)

	// Evaluate the transaction to query the total supply
	result, err := contract.EvaluateTransaction("TotalSupply")
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to evaluate transaction: %v", err), http.StatusInternalServerError)
		return
	}

	// Send the result back as a JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"total_supply": string(result),
	})
}
