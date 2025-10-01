package auth

import (
	"database/sql"
	"encoding/json"
	"net/http"
)

type API struct {
	db *sql.DB
}

func New(db *sql.DB) *API {
	return &API{db: db}
}

func (api *API) Login(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(Response{"TODO: Generate JWT"})
}

func (api *API) Refresh(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(Response{"TODO: Generate JWT"})
}
