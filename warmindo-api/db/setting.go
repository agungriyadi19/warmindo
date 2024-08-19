package db

type Settings struct {
	TotalTable int     `json:"total_table"`
	Latitude   float64 `json:"latitude"`
	Longitude  float64 `json:"longitude"`
	Radius     float64 `json:"radius"`
}
