package car

type Car struct {
	ID           int     `json:"id"`
	Model        string  `json:"model"`
	Brand        string  `json:"brand"`
	Year         int     `json:"year"`
	Color        string  `json:"color"`
	Doors        int     `json:"doors"`
	PricePerDay  float64 `json:"price_per_day"`
	LicensePlate string  `json:"license_plate"`
	// Stored in liters
	BaggageVolume float64 `json:"baggage_volume"`
}
