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
	"github.com/rodrigovieira938/goapi/config"
)

func New(db *sql.DB, cfg *config.Config) *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, World!")
	})

	authMiddleware := middleware.NewAuthMiddleware(&cfg.Auth, db)

	r.Use(middleware.JsonResponse)
	carAPI := cars.New(db)
	r.HandleFunc("/cars", carAPI.Get).Methods("GET")
	r.Handle("/cars", authMiddleware.WithPerms(http.HandlerFunc(carAPI.Get), []string{"write:cars"})).Methods("POST")

	userAPI := users.New(db, &cfg.Auth)
	r.Handle("/users", authMiddleware.Reject(http.HandlerFunc(userAPI.Post))).Methods("POST")
	r.Handle("/users/me", authMiddleware.Require(http.HandlerFunc(userAPI.Me))).Methods("GET")
	//Need read:users
	r.Handle("/users", authMiddleware.WithPerms(http.HandlerFunc(userAPI.Get), []string{"read:users"})).Methods("GET")
	r.Handle("/users/{id}", authMiddleware.WithPerms(http.HandlerFunc(userAPI.Id), []string{"read:users"})).Methods("GET")

	reservationAPI := reservations.New(db)

	//If user has perm read:reservations it reads reservations from all users instead of /users/me
	r.Handle("/reservations", authMiddleware.Require(http.HandlerFunc(reservationAPI.Get))).Methods("GET")
	r.Handle("/reservations", authMiddleware.Require(http.HandlerFunc(reservationAPI.Post))).Methods("POST")

	authAPI := auth.New(db, &cfg.Auth)
	r.HandleFunc("/auth/login", authAPI.Login).Methods("POST")
	return r
}
