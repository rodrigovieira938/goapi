package reservations

type Reservation struct {
	ID        int    `json:"id" validate:"gte=1"`
	UserID    int    `json:"user_id" validate:"gte=1"`
	CarID     int    `json:"car_id" validate:"gte=1"`
	StartDate string `json:"start_date" validate:"ISO8601date,required"`
	EndDate   string `json:"end_date" validate:"ISO8601date,required"`
}
