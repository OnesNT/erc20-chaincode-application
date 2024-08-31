// package main

// import (
// 	"log"
// 	"net/http"
// 	"rest-api-go/controllers"
// 	"rest-api-go/services"
// )

// func main() {
// 	cryptoPath := "../../test-network/organizations/peerOrganizations/org1.example.com"
// 	orgConfig := services.OrgSetup{
// 		OrgName:      "Org1",
// 		MSPID:        "Org1MSP",
// 		CertPath:     cryptoPath + "/users/Creator1@org1.example.com/msp/signcerts/cert.pem",
// 		KeyPath:      cryptoPath + "/users/Creator1@org1.example.com/msp/keystore/",
// 		TLSCertPath:  cryptoPath + "/peers/peer0.org1.example.com/tls/ca.crt",
// 		PeerEndpoint: "dns:///localhost:7051",
// 		GatewayPeer:  "peer0.org1.example.com",
// 	}

// 	orgSetup, err := services.Initialize(orgConfig)
// 	if err != nil {
// 		log.Fatalf("Error initializing setup for Org1: %v", err)
// 	}

// 	tokenController := controllers.NewTokenController(orgSetup)

// 	http.HandleFunc("/transfer", tokenController.Transfer)
// 	http.HandleFunc("/balance", tokenController.GetClientAccountBalance)
// 	http.HandleFunc("/invoke", tokenController.InitializeContract)
// 	http.HandleFunc("/mint", tokenController.Mint)

// 	log.Println("Starting server on port 8080")
// 	if err := http.ListenAndServe(":8080", nil); err != nil {
// 		log.Fatalf("Server failed: %v", err)
// 	}
// }

package main

import (
	"log"
	"net/http"
	"rest-api-go/controllers"
	"rest-api-go/services"
)

func main() {
	// Note: We are not initializing orgSetup here anymore.
	// We'll create a blank OrgSetup to be used later in the controller.
	orgConfig := &services.OrgSetup{}

	tokenController := controllers.NewTokenController(orgConfig)

	http.HandleFunc("/transfer", tokenController.Transfer)
	http.HandleFunc("/balance", tokenController.GetClientAccountBalance)
	http.HandleFunc("/invoke", tokenController.InitializeContract)
	http.HandleFunc("/mint", tokenController.Mint)

	log.Println("Starting server on port 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
