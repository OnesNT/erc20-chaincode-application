// package web

// import (
// 	"fmt"
// 	"net/http"

// 	"github.com/hyperledger/fabric-gateway/pkg/client"
// )

// // OrgSetup contains organization's config to interact with the network.
// type OrgSetup struct {
// 	OrgName      string
// 	MSPID        string
// 	CryptoPath   string
// 	CertPath     string
// 	KeyPath      string
// 	TLSCertPath  string
// 	PeerEndpoint string
// 	GatewayPeer  string
// 	Gateway      client.Gateway
// }

// // Serve starts http web server.
//
//	func Serve(setups OrgSetup) {
//		// http.HandleFunc("/query", setups.Query)
//		// http.HandleFunc("/invoke", setups.Invoke)
//		fmt.Println("Listening (http://localhost:3000/)...")
//		if err := http.ListenAndServe(":3000", nil); err != nil {
//			fmt.Println(err)
//		}
//	}
package web

import (
	"fmt"
	"net/http"

	"github.com/hyperledger/fabric-gateway/pkg/client"
)

type OrgSetup struct {
	OrgName      string
	MSPID        string
	CryptoPath   string
	CertPath     string
	KeyPath      string
	TLSCertPath  string
	PeerEndpoint string
	GatewayPeer  string
	Gateway      client.Gateway
}

// Serve starts http web server with check_balance and transfer handlers.
func Serve(setups OrgSetup) {
	http.HandleFunc("/check_balance", setups.CheckBalance)
	http.HandleFunc("/balance", setups.GetClientAccountBalance)
	http.HandleFunc("/query_total_supply", setups.QueryTotalSupply)
	http.HandleFunc("/transfer", setups.Transfer)
	http.HandleFunc("/transfer-wallet", setups.TransferFromAToB)
	http.HandleFunc("/initialize", setups.InitializeContract)
	http.HandleFunc("/mint", setups.Mint) // Add this line for minting
	fmt.Println("Listening (http://localhost:3000/)...")
	if err := http.ListenAndServe(":3000", nil); err != nil {
		fmt.Println(err)
	}
}
