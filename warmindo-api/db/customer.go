package db

type Customer struct {
	ID          int    `json:"id"`
	OrderCode   string `json:"order_code"`
	TableNumber int    `json:"table_number"`
	Active      bool   `json:"active"`
	StartDate   string `json:"start_date,omitempty"`
	EndDate     string `json:"end_date,omitempty"`
}
