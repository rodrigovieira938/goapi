package car

type Car struct {
	ID           int     `json:"id" validate:"gte=1"`
	Model        string  `json:"model" validate:"required"`
	Brand        string  `json:"brand" validate:"required"`
	Year         int     `json:"year" validate:"required"`
	Color        string  `json:"color" validate:"iscolor"`
	Doors        int     `json:"doors" validate:"required"`
	PricePerDay  float64 `json:"price_per_day" validate:"gt=0"`
	LicensePlate string  `json:"license_plate" validate:"required"`
	// Stored in liters
	BaggageVolume float64 `json:"baggage_volume" validate:"required"`
}
