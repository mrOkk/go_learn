package domain

type Transaction struct {
	ID         int64   `json:"id"`
	MerchantID int64   `json:"merchant_id"`
	Amount     float64 `json:"amount"`
	Status     string  `json:"status"`
	Timestamp  string  `json:"timestamp"`
}
