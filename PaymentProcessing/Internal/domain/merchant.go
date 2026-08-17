package domain

type Merchant struct {
	Id       string  `json:"id"`
	Name     string  `json:"name"`
	IsActive bool    `json:"is_active"`
	Balance  float64 `json:"balance"`
}
