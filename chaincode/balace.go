package chaincode

type Balance struct {
    ID          string  `json:"id"`
    Owner       string  `json:"owner"`
    Balance     float64 `json:"balance"`
    LastUpdated int     `json:"last_updated"` 
}

