package router

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rodrigovieira938/goapi/api/resource/car"
	"github.com/rodrigovieira938/goapi/api/router/middleware"
)

func New(db *sql.DB) *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, World!")
	})
	r.Use(middleware.JsonResponse)
	carAPI := car.New(db)
	r.HandleFunc("/car", carAPI.Get).Methods("GET")
	r.HandleFunc("/car", carAPI.Post).Methods("POST")

	return r
}
