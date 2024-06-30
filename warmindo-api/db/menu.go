package db

type Menu struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Image       string `json:"image"`
	Description string `json:"description"`
	Price       int    `json:"price"`
	CategoryID  int    `json:"category_id"`
	CreatedAt   string `json:"created_at,omitempty"`
	UpdatedAt   string `json:"updated_at,omitempty"`
}
