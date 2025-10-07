package cars

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

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
func (api *API) exists(id string) bool {
	var exists int
	err := api.db.QueryRow(`
        SELECT CASE WHEN EXISTS (
            SELECT 1 FROM car WHERE id=$1
        ) THEN 1 ELSE 0 END
    `, id).Scan(&exists)
	if err != nil {
		// treat DB error as user not existing
		return false
	}
	return exists == 1
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
		util.JsonError(w, "{\"error\":\"Invalid JSON!\"}", http.StatusBadRequest)
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
func (api *API) Put(w http.ResponseWriter, r *http.Request) {
	carId := mux.Vars(r)["id"]
	var car Car
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&car)
	if err != nil {
		util.JsonError(w, "{\"error\":\"Invalid JSON!\"}", http.StatusBadRequest)
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
	if !api.exists(carId) {
		util.JsonError(w, "{\"error\":\"Car doesn't exist\"}", http.StatusNotFound)
		return
	}
	_, err = api.db.Exec(`
        UPDATE car SET model=$1, brand=$2, year=$3, color=$4, doors=$5, price_per_day=$6, license_plate=$7, baggage_volume=$8 WHERE id=$9
    `, car.Model, car.Brand, car.Year, car.Color, car.Doors, car.PricePerDay, car.LicensePlate, car.BaggageVolume, carId)
	if err != nil {
		slog.Error("Put /cars/{id}", "err", err)
		http.Error(w, "DB error", http.StatusInternalServerError)
		return
	}
	w.Write([]byte("{}"))
}
func (api *API) Patch(w http.ResponseWriter, r *http.Request) {
	var patch map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&patch); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	setClauses := []string{}
	args := []interface{}{}
	i := 1

	if value, ok := patch["model"].(string); ok {
		setClauses = append(setClauses, fmt.Sprintf("model=$%d", i))
		args = append(args, value)
		i++
	}
	if value, ok := patch["branch"].(string); ok {
		setClauses = append(setClauses, fmt.Sprintf("branch=$%d", i))
		args = append(args, value)
		i++
	}
	if value, ok := patch["year"].(float64); ok {
		setClauses = append(setClauses, fmt.Sprintf("year=$%d", i))
		args = append(args, value)
		i++
	}
	if value, ok := patch["color"].(string); ok {
		setClauses = append(setClauses, fmt.Sprintf("color=$%d", i))
		args = append(args, value)
		i++
	}
	if value, ok := patch["doors"].(float64); ok {
		setClauses = append(setClauses, fmt.Sprintf("doors=$%d", i))
		args = append(args, value)
		i++
	}
	if value, ok := patch["price_per_day"].(float64); ok {
		setClauses = append(setClauses, fmt.Sprintf("price_per_day=$%d", i))
		args = append(args, value)
		i++
	}
	if value, ok := patch["license_plate"].(string); ok {
		setClauses = append(setClauses, fmt.Sprintf("license_plate=$%d", i))
		args = append(args, value)
		i++
	}
	if value, ok := patch["baggage_volume"].(float64); ok {
		setClauses = append(setClauses, fmt.Sprintf("price_per_day=$%d", i))
		args = append(args, value)
		i++
	}
	carId := mux.Vars(r)["id"]
	args = append(args, carId)
	query := fmt.Sprintf("UPDATE car SET %s WHERE id=$%d", strings.Join(setClauses, ","), i)
	if i == 1 {
		util.JsonError(w, "{\"error\":\"Empty request\"}", http.StatusBadRequest)
		return
	}

	if !api.exists(carId) {
		util.JsonError(w, "{\"error\":\"Car doesn't exist\"}", http.StatusNotFound)
		return
	}
	_, err := api.db.Exec(query, args...)
	if err != nil {
		slog.Error("PATCH /cars/{id}", "err", err)
		util.JsonError(w, "{\"error\":\"Internal Server Error\"}", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(patch)
}
func (api *API) Delete(w http.ResponseWriter, r *http.Request) {
	carId := mux.Vars(r)["id"]
	if !api.exists(carId) {
		util.JsonError(w, "{\"error\":\"Car doesn't exist\"}", http.StatusNotFound)
		return
	}
	api.db.Exec("DELETE FROM car WHERE id=$1", carId)
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "{}")
}
