package permissions

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rodrigovieira938/goapi/util"
)

type API struct {
	db *sql.DB
}

func New(db *sql.DB) *API {
	return &API{db: db}
}

func (api *API) Get(w http.ResponseWriter, r *http.Request) {
	rows, err := api.db.Query("SELECT * from permission;")
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

	var permissions []Permission = make([]Permission, 0)

	for rows.Next() {
		var perm Permission
		if err := rows.Scan(&perm.ID, &perm.Name, &perm.Description); err != nil {
			break
		}
		permissions = append(permissions, perm)
	}
	json.NewEncoder(w).Encode(permissions)
}
func (api *API) Id(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	row := api.db.QueryRow("SELECT * FROM \"permission\" WHERE id=$1", id)
	var perm Permission
	err := row.Scan(&perm.ID, &perm.Name, &perm.Description)
	if err != nil {
		util.JsonError(w, "{\"error\":\"Permission doesn't exist\"}", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(perm)
}
