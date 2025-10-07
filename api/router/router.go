package router

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rodrigovieira938/goapi/api/resource/auth"
	"github.com/rodrigovieira938/goapi/api/resource/cars"
	"github.com/rodrigovieira938/goapi/api/resource/permissions"
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
	r.HandleFunc("/cars/{id}", carAPI.Id).Methods("GET")
	//Need write:cars
	r.Handle("/cars", authMiddleware.WithPerms(http.HandlerFunc(carAPI.Post), []string{"write:cars"})).Methods("POST")
	r.Handle("/cars/{id}", authMiddleware.WithPerms(http.HandlerFunc(carAPI.Put), []string{"write:cars"})).Methods("PUT")
	r.Handle("/cars/{id}", authMiddleware.WithPerms(http.HandlerFunc(carAPI.Patch), []string{"write:cars"})).Methods("PATCH")
	r.Handle("/cars/{id}", authMiddleware.WithPerms(http.HandlerFunc(carAPI.Delete), []string{"write:cars"})).Methods("DELETE")

	userAPI := users.New(db, &cfg.Auth)
	r.Handle("/users", authMiddleware.Reject(http.HandlerFunc(userAPI.Post))).Methods("POST")
	r.Handle("/users/me", authMiddleware.Require(http.HandlerFunc(userAPI.Me))).Methods("GET")
	//Need read:users
	r.Handle("/users", authMiddleware.WithPerms(http.HandlerFunc(userAPI.Get), []string{"read:users"})).Methods("GET")
	r.Handle("/users/{id}", authMiddleware.WithPerms(http.HandlerFunc(userAPI.Id), []string{"read:users"})).Methods("GET")

	reservationAPI := reservations.New(db, &cfg.Auth)

	//Get reservations of /users/me
	r.Handle("/reservations", authMiddleware.Require(http.HandlerFunc(reservationAPI.Get))).Methods("GET")
	//Get reservations: if reservations is owned by user return; else demand read:perms
	r.Handle("/reservations/{id}", authMiddleware.Require(http.HandlerFunc(reservationAPI.Id))).Methods("GET")
	//Need write:perms
	r.Handle("/reservations", authMiddleware.WithPerms(http.HandlerFunc(reservationAPI.Post), []string{"write:reservations"})).Methods("POST")
	r.Handle("/reservations/{id}", authMiddleware.WithPerms(http.HandlerFunc(reservationAPI.Put), []string{"write:reservations"})).Methods("PUT")
	r.Handle("/reservations/{id}", authMiddleware.WithPerms(http.HandlerFunc(reservationAPI.Patch), []string{"write:reservations"})).Methods("PATCH")
	r.Handle("/reservations/{id}", authMiddleware.WithPerms(http.HandlerFunc(reservationAPI.Delete), []string{"write:reservations"})).Methods("DELETE")

	permissionAPI := permissions.New(db)

	r.HandleFunc("/permissions", permissionAPI.Get).Methods("GET")
	r.HandleFunc("/permissions/{id}", permissionAPI.Id).Methods("GET")

	authAPI := auth.New(db, &cfg.Auth)
	r.HandleFunc("/auth/login", authAPI.Login).Methods("POST")
	return r
}
