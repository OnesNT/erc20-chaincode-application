package models

type OrgSetup struct {
	OrgName      string
	MSPID        string
	CertPath     string
	KeyPath      string
	TLSCertPath  string
	PeerEndpoint string
	GatewayPeer  string
	Gateway      Gateway
}
