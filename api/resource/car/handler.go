package car

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/rodrigovieira938/goapi/util"
)

type API struct {
	db *sql.DB
}

func New(db *sql.DB) *API {
	return &API{db: db}
}

func (api *API) Get(w http.ResponseWriter, r *http.Request) {
	rows, err := api.db.Query("SELECT * from car;")
	w.Header().Add("Content-Type", "application/json")
	if err != nil {
		if err == sql.ErrNoRows {
			w.Write([]byte("[]"))
			return
		}
		util.JsonError(w, "{\"error\":\"Internal Server Error\"}", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var cars []Car

	for rows.Next() {
		var car Car
		if err := rows.Scan(&car.ID, &car.Model, &car.Brand,
			&car.Year, &car.Color, &car.Doors, &car.PricePerDay, &car.LicensePlate, &car.BaggageVolume); err != nil {
			break
		}
		cars = append(cars, car)
	}
	json.NewEncoder(w).Encode(cars)
}
func (api *API) Post(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		util.JsonError(w, "{\"error\":\"Content-Type must be application/json\"}", http.StatusBadRequest)
		return
	}
	var car Car
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&car)
	if err != nil {
		util.JsonError(w, "{\"error\":\"Invalid JSON! Check if json is valid or if all required fields are present\"}", http.StatusBadRequest)
		return
	}
	row := api.db.QueryRow(
		`INSERT INTO car (model, brand, year, color, doors, price_per_day, license_plate, baggage_volume) 
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		 RETURNING id`,
		car.Model, car.Brand, car.Year, car.Color, car.Doors, car.PricePerDay, car.LicensePlate, car.BaggageVolume)
	err = row.Scan(&car.ID)
	if err != nil {
		util.JsonError(w, "{\"error\":\"Internal Server Error\"}", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(car)
}
