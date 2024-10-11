package chaincode

type SettingCommon struct {
	ID                 	string 	`json:"id"`
	InterestRate		float64	`json:"inteset_rate"`
	CirculatingSupply	float64	`json:"curculating_supply"`
}