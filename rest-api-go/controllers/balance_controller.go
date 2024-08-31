package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"rest-api-go/services"
)

// BalanceController handles requests for querying client account balance.
type BalanceController struct {
	Service *services.GatewayService
}

// NewBalanceController creates a new BalanceController instance.
func NewBalanceController(setup *services.OrgSetup) *BalanceController {
	return &BalanceController{Service: services.NewGatewayService(setup)}
}

// GetClientAccountBalance handles the request to get client account balance.
func (c *BalanceController) GetClientAccountBalance(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	certificatePEM, privateKeyPEM, err := utils.getCertificateAndPrivateKeyFromForm(c)

	chainCodeName := r.URL.Query().Get("chaincodeid")
	channelID := r.URL.Query().Get("channelid")

	if chainCodeName == "" || channelID == "" {
		http.Error(w, "Missing required fields: chaincodeid or channelid", http.StatusBadRequest)
		return
	}

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

	result, err := contract.EvaluateTransaction("ClientAccountBalance")
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to evaluate transaction: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"balance": string(result),
	})
}
