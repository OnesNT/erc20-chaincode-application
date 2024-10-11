package chaincode

type Loan struct {
	ID                 string    `json:"id"`
	Owner              string    `json:"owner"`
	PrincipalAmount    float64   `json:"principal_amount"`
	InterestRate       float64   `json:"interest_rate"`
	StartDate          time.Time `json:"start_date"`
	EndDate            time.Time `json:"end_date"`
	RemainingPrincipal float64   `json:"remaining_principal"`
	LoanStatus         string    `json:"loan_status"`
}