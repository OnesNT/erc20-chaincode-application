package main

import (
	"log"
	"net/http"
	"rest-api-go/controllers"
	"rest-api-go/services"
)

func main() {
	cryptoPath := "../../test-network/organizations/peerOrganizations/org1.example.com"
	orgConfig := services.OrgSetup{
		OrgName:      "Org1",
		MSPID:        "Org1MSP",
		CertPath:     cryptoPath + "/users/Creator1@org1.example.com/msp/signcerts/cert.pem",
		KeyPath:      cryptoPath + "/users/Creator1@org1.example.com/msp/keystore/",
		TLSCertPath:  cryptoPath + "/peers/peer0.org1.example.com/tls/ca.crt",
		PeerEndpoint: "dns:///localhost:7051",
		GatewayPeer:  "peer0.org1.example.com",
	}

	orgSetup, err := services.Initialize(orgConfig)
	if err != nil {
		log.Fatalf("Error initializing setup for Org1: %v", err)
	}

	transferController := controllers.NewTransferController(orgSetup)
	balanceController := controllers.NewBalanceController(orgSetup)
	InvokeController := controllers.NewInvokeController(orgSetup)
	MintController := controllers.NewMinController(orgSetup)

	http.HandleFunc("/transfer", transferController.Transfer)
	http.HandleFunc("/balance", balanceController.GetClientAccountBalance)
	http.HandleFunc("/invoke", InvokeController.InitializeContract)
	http.HandleFunc("/mint", MintController.Mint)

	log.Println("Starting server on port 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
