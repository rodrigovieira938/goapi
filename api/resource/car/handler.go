package car

import (
	"database/sql"
	"fmt"
	"net/http"
)

type API struct {
	db *sql.DB
}

func New(db *sql.DB) *API {
	return &API{db: db}
}

func (api *API) Get(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "{\"test\":\"Car GET endpoint\"}")
}
func (api *API) Post(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "{\"test\":\"Car POST endpoint\"}")
}
