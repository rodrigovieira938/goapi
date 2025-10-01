package router

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rodrigovieira938/goapi/api/resource/auth"
	"github.com/rodrigovieira938/goapi/api/resource/cars"
	"github.com/rodrigovieira938/goapi/api/resource/reservations"
	"github.com/rodrigovieira938/goapi/api/resource/users"
	"github.com/rodrigovieira938/goapi/api/router/middleware"
)

func New(db *sql.DB) *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, World!")
	})
	r.Use(middleware.JsonResponse)
	carAPI := cars.New(db)
	r.HandleFunc("/cars", carAPI.Get).Methods("GET")
	r.HandleFunc("/cars", carAPI.Post).Methods("POST")

	userAPI := users.New(db)
	r.HandleFunc("/users", userAPI.Get).Methods("GET")
	r.HandleFunc("/users", userAPI.Post).Methods("POST")

	reservationAPI := reservations.New(db)
	r.HandleFunc("/reservations", reservationAPI.Get).Methods("GET")
	r.HandleFunc("/reservations", reservationAPI.Post).Methods("POST")

	authAPI := auth.New(db)
	r.HandleFunc("/auth/login", authAPI.Login).Methods("POST")
	r.HandleFunc("/auth/refresh", authAPI.Refresh).Methods("POST")

	return r
}
