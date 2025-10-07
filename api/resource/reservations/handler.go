package reservations

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/rodrigovieira938/goapi/api/router/middleware"
	"github.com/rodrigovieira938/goapi/config"
	"github.com/rodrigovieira938/goapi/util"
	timeformat "github.com/rodrigovieira938/goapi/util/time_format"
)

type API struct {
	db *sql.DB
	//only here to acess functions from the middleware
	auth *middleware.AuthMiddleware
}

func New(db *sql.DB, cfg *config.AuthConfig) *API {
	return &API{db: db, auth: middleware.NewAuthMiddleware(cfg, db)}
}
func (api *API) Get(w http.ResponseWriter, r *http.Request) {
	token, _ := api.auth.ParseToken(r.Header.Get("Authorization")) //This token was checked by Require
	userId, _ := api.auth.GetIDFromToken(token)
	rows, err := api.db.Query("SELECT * from reservation WHERE user_id = $1;", userId)
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

	var reservations []Reservation = make([]Reservation, 0)

	for rows.Next() {
		var reservation Reservation
		if err := rows.Scan(&reservation.ID, &reservation.UserID, &reservation.CarID, &reservation.StartDate, &reservation.EndDate); err != nil {
			break
		}
		reservations = append(reservations, reservation)
	}
	json.NewEncoder(w).Encode(reservations)
}

func (api *API) Id(w http.ResponseWriter, r *http.Request) {
	token, _ := api.auth.ParseToken(r.Header.Get("Authorization")) //This token was checked by Require
	userId, _ := api.auth.GetIDFromToken(token)
	reservationID := mux.Vars(r)["id"]
	row := api.db.QueryRow("SELECT * FROM \"reservation\" WHERE id=$1", reservationID)
	var reservation Reservation
	err := row.Scan(&reservation.ID, reservation.UserID, reservation.CarID, reservation.StartDate, reservation.EndDate)
	if err != nil {
		util.JsonError(w, "{\"error\":\"Reservation doesn't exist\"}", http.StatusNotFound)
		return
	}
	if reservation.UserID != userId && !api.auth.UserHasPerm(userId, "read:perms") {
		util.JsonError(w, "{\"error\":\"Forbiden\"}", http.StatusForbidden)
		return
	}
	json.NewEncoder(w).Encode(reservation)
}

func (api *API) Post(w http.ResponseWriter, r *http.Request) {
	var reservation Reservation
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&reservation)
	if err != nil {
		util.JsonError(w, "{\"error\":\"Invalid JSON!\"}", http.StatusBadRequest)
		return
	}
	reservation.ID = 1 // Make id valid since its ignored
	validate := validator.New(validator.WithRequiredStructEnabled())
	_ = validate.RegisterValidation("ISO8601date", timeformat.IsISO8601Date)
	err = validate.Struct(reservation)
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		if validationErrors != nil {
			util.JsonError(w, "{\"error\":\""+validationErrors.Error()+"\"}", http.StatusBadRequest)
			return
		}
	}
	row := api.db.QueryRow("INSERT into reservation (user_id, car_id, start_date, end_date) VALUES ($1, $2, $3, $4) RETURNING id", reservation.UserID, reservation.CarID, reservation.StartDate, reservation.EndDate)
	err = row.Scan(&reservation.ID)
	if err != nil {
		slog.Error("Error inserting reservation", "error", err)
		util.JsonError(w, "{\"error\":\"Internal Server Error\"}", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reservation)
}

func (api *API) exists(id string) bool {
	var exists int
	err := api.db.QueryRow(`
        SELECT CASE WHEN EXISTS (
            SELECT 1 FROM reservation WHERE id=$1
        ) THEN 1 ELSE 0 END
    `, id).Scan(&exists)
	if err != nil {
		// treat DB error as user not existing
		return false
	}
	return exists == 1
}

func (api *API) Put(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var reservation Reservation
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&reservation)
	if err != nil {
		util.JsonError(w, "{\"error\":\"Invalid JSON!\"}", http.StatusBadRequest)
		return
	}
	reservation.ID = 1 // Make id valid since its ignored
	validate := validator.New(validator.WithRequiredStructEnabled())
	_ = validate.RegisterValidation("ISO8601date", timeformat.IsISO8601Date)
	err = validate.Struct(reservation)
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		if validationErrors != nil {
			util.JsonError(w, "{\"error\":\""+validationErrors.Error()+"\"}", http.StatusBadRequest)
			return
		}
	}
	if !api.exists(id) {
		util.JsonError(w, "{\"error\":\"Reservation doesn't exist\"}", http.StatusNotFound)
		return
	}
	_, err = api.db.Exec(`
        UPDATE car SET user_id=$1, car_id=$2, start_date=$3, end_date=$4 WHERE id=$5
    `, reservation.UserID, reservation.CarID, reservation.StartDate, reservation.EndDate, id)
	if err != nil {
		slog.Error("Put /reservations/{id}", "err", err)
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

	if value, ok := patch["id"].(string); ok {
		setClauses = append(setClauses, fmt.Sprintf("id=$%d", i))
		args = append(args, value)
		i++
	}
	if value, ok := patch["user_id"].(float64); ok {
		setClauses = append(setClauses, fmt.Sprintf("user_id=$%d", i))
		args = append(args, value)
		i++
	}
	if value, ok := patch["car_id"].(float64); ok {
		setClauses = append(setClauses, fmt.Sprintf("car_id=$%d", i))
		args = append(args, value)
		i++
	}
	if value, ok := patch["start_date"].(string); ok {
		setClauses = append(setClauses, fmt.Sprintf("start_date=$%d", i))
		args = append(args, value)
		i++
	}
	if value, ok := patch["end_date"].(float64); ok {
		setClauses = append(setClauses, fmt.Sprintf("start_date=$%d", i))
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
		slog.Error("PATCH /reservations/{id}", "err", err)
		util.JsonError(w, "{\"error\":\"Internal Server Error\"}", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(patch)
}

func (api *API) Delete(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if !api.exists(id) {
		util.JsonError(w, "{\"error\":\"Reservation doesn't exist\"}", http.StatusNotFound)
		return
	}
	api.db.Exec("DELETE FROM reservation WHERE id=$1", id)
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "{}")
}
