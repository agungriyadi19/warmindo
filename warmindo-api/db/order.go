package db

type Order struct {
	ID          int    `json:"id"`
	Amount      string `json:"amount"`
	TableNumber string `json:"table_number"`
	StatusID    string `json:"status_id"`
	OrderDate   string `json:"order_date"`
	MenuID      string `json:"menu_id"`
	CreatedAt   string `json:"created_at,omitempty"`
	UpdatedAt   string `json:"updated_at,omitempty"`
}
