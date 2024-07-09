package db

type Order struct {
	ID          int    `json:"id"`
	OrderCode   string `json:"order_code"`
	Amount      string `json:"amount"`
	TableNumber string `json:"table_number"`
	StatusID    int    `json:"status_id"`
	OrderDate   string `json:"order_date"`
	MenuID      int    `json:"menu_id"`
	CreatedAt   string `json:"created_at,omitempty"`
	UpdatedAt   string `json:"updated_at,omitempty"`
}
