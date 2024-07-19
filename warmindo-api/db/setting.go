package db

type Settings struct {
	Title     string  `json:"title"`
	Subtitle  string  `json:"subtitle"`
	Seats     int     `json:"seats"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Radius    float64 `json:"radius"`
}
