package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"rest-api-go/services"
	"strings"
)

// InvokeController handles requests for invoke token.
type TokenController struct {
	Service *services.GatewayService
}

// NewTransferController creates a new TransferController instance.
func NewTokenController(setup *services.OrgSetup) *TokenController {
	return &TokenController{Service: services.NewGatewayService(setup)}
}

// InitializeContract handles initializing the chaincode with token information.
func (c *TokenController) InitializeContract(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Initialize the service with the cert and key from the request
	// err := c.Service.InitializeWithCertAndKey(r)
	// if err != nil {
	// 	http.Error(w, fmt.Sprintf("Failed to initialize with certificate and key: %v", err), http.StatusBadRequest)
	// 	return
	// }

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

	// Call the service to initialize the contract
	transactionID, err := c.Service.CallChaincode(channelID, chainCodeName, "Initialize", []string{tokenName, symbol, decimals})
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to initialize contract: %v", err), http.StatusInternalServerError)
		return
	}

	// Respond with success message
	fmt.Fprintf(w, "Initialization successful. Transaction ID: %s", transactionID)
}

// Mint handles chaincode mint requests.
func (c *TokenController) Mint(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Initialize the service with the cert and key from the request
	// err := c.Service.InitializeWithCertAndKey(r)
	// if err != nil {
	// 	http.Error(w, fmt.Sprintf("Failed to initialize with certificate and key: %v", err), http.StatusBadRequest)
	// 	return
	// }

	fmt.Println("Received Mint request")
	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "ParseForm() err: %s", err)
		return
	}

	// Extract chaincode-related information from the request
	chainCodeName := r.FormValue("chaincodeid")
	channelID := r.FormValue("channelid")
	amount := r.FormValue("amount")

	if chainCodeName == "" || channelID == "" || amount == "" {
		http.Error(w, "Missing required fields: chaincodeid, channelid, or amount", http.StatusBadRequest)
		return
	}

	fmt.Printf("Channel: %s, Chaincode: %s, Amount: %s\n", channelID, chainCodeName, amount)

	// Call the service to initialize the contract
	transactionID, err := c.Service.CallChaincode(channelID, chainCodeName, "Mint", []string{amount})
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to initialize contract: %v", err), http.StatusInternalServerError)
		return
	}
	// Respond with success message
	fmt.Fprintf(w, "Minting successful. Transaction ID: %s", transactionID)
}

// GetClientAccountBalance handles the request to get client account balance.
func (c *TokenController) GetClientAccountBalance(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Initialize the service with the cert and key from the request
	// err := c.Service.InitializeWithCertAndKey(r)
	// if err != nil {
	// 	http.Error(w, fmt.Sprintf("Failed to initialize with certificate and key: %v", err), http.StatusBadRequest)
	// 	return
	// }

	chainCodeName := r.URL.Query().Get("chaincodeid")
	channelID := r.URL.Query().Get("channelid")

	if chainCodeName == "" || channelID == "" {
		http.Error(w, "Missing required fields: chaincodeid or channelid", http.StatusBadRequest)
		return
	}

	// Call the service to get the client account balance
	result, err := c.Service.CallChaincodeGET(channelID, chainCodeName, "ClientAccountBalance")
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get client account balance: %v", err), http.StatusInternalServerError)
		return
	}

	// Respond with the balance
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"balance": result,
	})
}

// Transfer handles chaincode invoke requests for transferring an asset.
func (c *TokenController) Transfer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Initialize the service with the cert and key from the request
	// err := c.Service.InitializeWithCertAndKey(r)
	// if err != nil {
	// 	http.Error(w, fmt.Sprintf("Failed to initialize with certificate and key: %v", err), http.StatusBadRequest)
	// 	return
	// }

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

	// Call the service to initialize the contract
	transactionID, err := c.Service.CallChaincode(channelID, chainCodeName, "Transfer", []string{recipient, amount})
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to initialize contract: %v", err), http.StatusInternalServerError)
		return
	}

	// Respond with success message
	fmt.Fprintf(w, "Initialization successful. Transaction ID: %s", transactionID)
}
