package models

type Wallet struct {
	ID        string  `json:"id" gorm:"primaryKey"`
	UserID    string  `json:"user_id"`
	Balance   float64 `json:"balance"`
	Nonce     string  `json:"nonce" gorm:"column:nonce"`
	Currency  string  `json:"currency"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
}
