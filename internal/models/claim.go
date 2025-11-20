package models

type Claim struct {
	ID         string  `json:"id" gorm:"primaryKey"`
	UserID     string  `json:"user_id"`
	PropertyID string  `json:"property_id"`
	Amount     float64 `json:"amount"`
	Status     string  `json:"status"`
	CreatedAt  string  `json:"created_at"`
	UpdatedAt  string  `json:"updated_at"`
}
