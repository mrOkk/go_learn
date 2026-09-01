package domain

type Merchant struct {
	Id       int64   `json:"id"`
	Name     string  `json:"name"`
	IsActive bool    `json:"is_active"`
	Balance  float64 `json:"balance"`
}
