package models

type Distribution struct {
	ID          string  `json:"id"  gorm:"primaryKey"`
	Amount      float64 `json:"amount"`
	RecipientID string  `json:"recipient_id"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}
