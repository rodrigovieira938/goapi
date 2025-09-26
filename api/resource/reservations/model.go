package reservations

type Reservation struct {
	ID        int    `json:"id"`
	UserID    int    `json:"user_id"`
	CarID     int    `json:"car_id"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}
