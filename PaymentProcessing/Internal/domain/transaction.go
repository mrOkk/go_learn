package domain

type Transaction struct {
	ID         string  `json:"id"`
	MerchantID string  `json:"merchant_id"`
	Amount     float64 `json:"amount"`
	Status     string  `json:"status"`
	Timestamp  string  `json:"timestamp"`
}
