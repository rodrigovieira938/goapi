package reservations

import (
	"database/sql"
	"encoding/json"
	"log/slog"
	"net/http"

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
	if r.Header.Get("Content-Type") != "application/json" {
		util.JsonError(w, "{\"error\":\"Content-Type must be application/json\"}", http.StatusBadRequest)
		return
	}
	var reservation Reservation
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&reservation)
	if err != nil {
		util.JsonError(w, "{\"error\":\"Invalid JSON! Check if json is valid or if all required fields are present\"}", http.StatusBadRequest)
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
		//TODO: check for unique email and username
		slog.Error("Error inserting car", "error", err)
		util.JsonError(w, "{\"error\":\"Internal Server Error\"}", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reservation)
}
