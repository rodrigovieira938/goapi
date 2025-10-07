package users

import (
	"database/sql"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/lib/pq"
	"github.com/rodrigovieira938/goapi/api/resource/permissions"
	"github.com/rodrigovieira938/goapi/api/router/middleware"
	"github.com/rodrigovieira938/goapi/config"
	"github.com/rodrigovieira938/goapi/util"
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
	rows, err := api.db.Query("SELECT * from \"user\";")
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

	var users []User = make([]User, 0)

	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Username, &user.Email,
			&user.Password); err != nil {
			break
		}
		users = append(users, user)
	}
	json.NewEncoder(w).Encode(users)
}

func (api *API) Post(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		util.JsonError(w, "{\"error\":\"Content-Type must be application/json\"}", http.StatusBadRequest)
		return
	}
	var user User
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&user)
	if err != nil {
		util.JsonError(w, "{\"error\":\"Invalid JSON!\"}", http.StatusBadRequest)
		return
	}
	user.ID = 1 // Make id valid since its ignored
	validate := validator.New(validator.WithRequiredStructEnabled())
	err = validate.Struct(user)
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		if validationErrors != nil {
			util.JsonError(w, "{\"error\":\""+validationErrors.Error()+"\"}", http.StatusBadRequest)
			return
		}
	}
	//TODO: encrypt the password and add salt
	row := api.db.QueryRow(`INSERT INTO "user" (username, email, password) VALUES ($1, $2, $3) RETURNING id`, user.Username, user.Email, user.Password)
	err = row.Scan(&user.ID)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23505" {
				util.JsonError(w, "{\"error\":\"Email or username already in use!\"}", http.StatusConflict)
				return
			}
		} else {
			slog.Error("Error inserting car", "error", err)
		}
		util.JsonError(w, "{\"error\":\"Internal Server Error\"}", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (api *API) Me(w http.ResponseWriter, r *http.Request) {
	token, _ := api.auth.ParseToken(r.Header.Get("Authorization")) //This token was checked by Require
	id, _ := api.auth.GetIDFromToken(token)
	row := api.db.QueryRow("SELECT * FROM \"user\" WHERE id=$1", id)
	if row == nil {
		//Shouldn't happend since AuthMiddleware.Require checked it
		util.JsonError(w, "{\"error\":\"Internal Server Error\"}", http.StatusInternalServerError)
		return
	}
	var user User
	row.Scan(&user.ID, &user.Username, &user.Email, &user.Password)
	json.NewEncoder(w).Encode(map[string]any{"id": user.ID, "username": user.Username, "email": user.Email})
}
func (api *API) Id(w http.ResponseWriter, r *http.Request) {
	userId := mux.Vars(r)["id"]
	row := api.db.QueryRow("SELECT * FROM \"user\" WHERE id=$1", userId)
	var user User
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Password)
	if err != nil {
		util.JsonError(w, "{\"error\":\"User doesn't exist\"}", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(map[string]any{"id": user.ID, "username": user.Username, "email": user.Email})
}

func (api *API) Perms(w http.ResponseWriter, r *http.Request) {
	token, _ := api.auth.ParseToken(r.Header.Get("Authorization")) //This token was checked by Require
	userId, _ := api.auth.GetIDFromToken(token)
	rows, err := api.db.Query(`
		SELECT p.id, p.name, p.description
			FROM permission p
			JOIN user_permission up ON p.id = up.permission_id
			WHERE up.user_id = $1;
	`, userId)
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

	var perms []permissions.Permission = make([]permissions.Permission, 0)

	for rows.Next() {
		var perm permissions.Permission
		if err := rows.Scan(&perm.ID, &perm.Name, &perm.Description); err != nil {
			break
		}
		perms = append(perms, perm)
	}
	json.NewEncoder(w).Encode(perms)
}

func (api *API) PermsId(w http.ResponseWriter, r *http.Request) {
	userId := mux.Vars(r)["id"]
	rows, err := api.db.Query(`
		SELECT p.id, p.name, p.description
			FROM permission p
			JOIN user_permission up ON p.id = up.permission_id
			WHERE up.user_id = $1;
	`, userId)
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

	var perms []permissions.Permission = make([]permissions.Permission, 0)

	for rows.Next() {
		var perm permissions.Permission
		if err := rows.Scan(&perm.ID, &perm.Name, &perm.Description); err != nil {
			break
		}
		perms = append(perms, perm)
	}
	json.NewEncoder(w).Encode(perms)
}
