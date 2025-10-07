package cars

import (
	"database/sql"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/lib/pq"
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

	var cars []Car = make([]Car, 0)

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
func (api *API) Id(w http.ResponseWriter, r *http.Request) {
	carId := mux.Vars(r)["id"]
	row := api.db.QueryRow("SELECT * FROM \"car\" WHERE id=$1", carId)
	var car Car
	err := row.Scan(&car.ID, &car.Model, &car.Brand, &car.Year, &car.Color, &car.Doors, &car.PricePerDay, &car.LicensePlate, &car.BaggageVolume)
	if err != nil {
		util.JsonError(w, "{\"error\":\"Car doesn't exist\"}", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(car)
}
func (api *API) Post(w http.ResponseWriter, r *http.Request) {
	var car Car
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&car)
	if err != nil {
		util.JsonError(w, "{\"error\":\"Invalid JSON! Check if json is valid or if all required fields are present\"}", http.StatusBadRequest)
		return
	}
	car.ID = 1 // Make id valid since its ignored
	validate := validator.New(validator.WithRequiredStructEnabled())
	err = validate.Struct(car)
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		if validationErrors != nil {
			util.JsonError(w, "{\"error\":\""+validationErrors.Error()+"\"}", http.StatusBadRequest)
			return
		}
	}
	row := api.db.QueryRow(
		`INSERT INTO car (model, brand, year, color, doors, price_per_day, license_plate, baggage_volume) 
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		 RETURNING id`,
		car.Model, car.Brand, car.Year, car.Color, car.Doors, car.PricePerDay, car.LicensePlate, car.BaggageVolume)
	err = row.Scan(&car.ID)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23505" {
				if pqErr.Constraint == "car_license_plate_key" {
					util.JsonError(w, "{\"error\":\"Car with the same license plate already exists\"}", http.StatusConflict)
					return
				} else {
					slog.Error("TODO Handle unique constrait", "col", pqErr.Constraint)
				}
			}
		} else {
			slog.Error("Error inserting car", "error", err)
		}
		util.JsonError(w, "{\"error\":\"Internal Server Error\"}", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(car)
}
