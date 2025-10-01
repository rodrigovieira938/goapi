package auth

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/rodrigovieira938/goapi/config"
)

type API struct {
	db  *sql.DB
	cfg *config.AuthConfig
}

func New(db *sql.DB, cfg *config.AuthConfig) *API {
	return &API{db: db, cfg: cfg}
}

func (api *API) Login(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(Response{"TODO: Generate JWT"})
}

func (api *API) Refresh(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(Response{"TODO: Generate JWT"})
}
